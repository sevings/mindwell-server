package helper

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/lib/database"
	"github.com/zpatrick/go-config"
)

type validPathSet struct {
	albumPaths         map[string]bool // Normalized album paths (e.g., "W/WY/TNPxFiq")
	avatarPaths        map[string]bool // Normalized avatar paths (e.g., "W/a2qF6")
	coverPaths         map[string]bool // Normalized cover paths (e.g., "l/EXNeC")
	animatedAlbumPaths map[string]bool // Which albums have PNG previews
}

type cleanupStats struct {
	albumsFound     int
	albumsOrphaned  int
	avatarsFound    int
	avatarsOrphaned int
	coversFound     int
	coversOrphaned  int
	totalSize       int64
	orphanedFiles   []string
	errors          []string
}

type unusedImage struct {
	id               int64
	path             string
	extension        string
	previewExtension string
	createdAt        string
	files            []string // Full paths to files on disk
	totalSize        int64
}

var (
	albumPathRe  = regexp.MustCompile(`^albums/(thumbnails|small|medium|large)/(.+)$`)
	avatarPathRe = regexp.MustCompile(`^avatars/(124|92|42)/(.+)$`)
	coverPathRe  = regexp.MustCompile(`^covers/(1920|318)/(.+)$`)
)

func CleanupOrphanedImages(tx *database.AutoTx, cfg *config.Config) {
	log.Println("Starting image cleanup...")

	// Parse command-line flags
	deleteMode := false
	verbose := false
	checkUnused := false
	for _, arg := range os.Args {
		if arg == "--delete" {
			deleteMode = true
		}
		if arg == "--verbose" {
			verbose = true
		}
		if arg == "--unused" {
			checkUnused = true
		}
	}

	// Load image folder path from config
	imageFolder, err := cfg.String("images.folder")
	if err != nil || imageFolder == "" {
		log.Println("Error: Could not read images.folder from config")
		return
	}

	log.Printf("Image folder: %s\n", imageFolder)
	if deleteMode {
		log.Println("Mode: DELETE (will remove files after confirmation)")
	} else {
		log.Println("Mode: DRY-RUN (no files will be deleted)")
	}

	// Check for unused database images if --unused flag is set
	if checkUnused {
		log.Println("\n" + strings.Repeat("=", 70))
		log.Println("CHECKING FOR UNUSED DATABASE IMAGES")
		log.Println(strings.Repeat("=", 70))
		log.Println("\nFinding images in database not attached to any entry (older than 6 months)...")

		unusedImages := findUnusedDatabaseImages(tx, imageFolder)
		if tx.HasQueryError() {
			log.Printf("Error finding unused images: %v\n", tx.Error())
			return
		}

		printUnusedImageStats(unusedImages)

		if deleteMode && len(unusedImages) > 0 {
			if confirmUnusedDeletion(unusedImages) {
				deleteUnusedImages(tx, unusedImages)
			} else {
				log.Println("\nDeletion cancelled by user.")
			}
		} else if len(unusedImages) > 0 {
			log.Println("\nRun with --delete flag to remove unused images.")
		}

		if tx.HasQueryError() {
			log.Printf("Error during deletion: %v\n", tx.Error())
			return
		}
	}

	// Regular orphaned file scan
	log.Println("\n" + strings.Repeat("=", 70))
	log.Println("CHECKING FOR ORPHANED FILES")
	log.Println(strings.Repeat("=", 70))

	// Load valid paths from database
	log.Println("\nLoading valid paths from database...")
	validPaths := loadValidPaths(tx)
	if tx.HasQueryError() {
		log.Printf("Error loading paths from database: %v\n", tx.Error())
		return
	}

	log.Printf("  - Loaded %d album images\n", len(validPaths.albumPaths))
	log.Printf("  - Loaded %d avatars\n", len(validPaths.avatarPaths))
	log.Printf("  - Loaded %d covers\n", len(validPaths.coverPaths))
	log.Printf("  - Loaded %d animated albums\n", len(validPaths.animatedAlbumPaths))

	// Scan filesystem for orphaned files
	log.Printf("\nScanning %s for orphaned files...\n", imageFolder)
	stats := walkAndIdentifyOrphans(imageFolder, validPaths, verbose)

	// Display statistics
	printStats(stats)

	// If delete mode, confirm and delete
	if deleteMode && len(stats.orphanedFiles) > 0 {
		if confirmDeletion(stats) {
			deleteOrphanedFiles(stats)
		} else {
			log.Println("\nDeletion cancelled by user.")
		}
	} else if len(stats.orphanedFiles) > 0 {
		log.Println("\nRun with --delete flag to remove orphaned files.")
	}

	if len(stats.errors) > 0 {
		log.Printf("\n%d errors occurred during scanning:\n", len(stats.errors))
		for i, err := range stats.errors {
			if i < 10 {
				log.Printf("  %d. %s\n", i+1, err)
			}
		}
		if len(stats.errors) > 10 {
			log.Printf("  ... and %d more errors\n", len(stats.errors)-10)
		}
	}
}

func loadValidPaths(tx *database.AutoTx) validPathSet {
	validPaths := validPathSet{
		albumPaths:         make(map[string]bool),
		avatarPaths:        make(map[string]bool),
		coverPaths:         make(map[string]bool),
		animatedAlbumPaths: make(map[string]bool),
	}

	// Load album images
	q := sqlf.Select("path, extension, preview_extension").
		From("images").
		Where("processing = false")

	tx.QueryStmt(q)
	for {
		var path, extension, previewExtension string
		if !tx.Scan(&path, &extension, &previewExtension) {
			break
		}
		validPaths.albumPaths[path] = true
		if previewExtension != "" {
			validPaths.animatedAlbumPaths[path] = true
		}
	}

	// Load avatar paths
	q = sqlf.Select("DISTINCT avatar").
		From("users").
		Where("avatar IS NOT NULL").
		Where("avatar <> ''")

	tx.QueryStmt(q)
	for {
		var avatar string
		if !tx.Scan(&avatar) {
			break
		}
		// Strip extension to get normalized path
		normalized := stripExtension(avatar)
		validPaths.avatarPaths[normalized] = true
	}

	// Load cover paths
	q = sqlf.Select("DISTINCT cover").
		From("users").
		Where("cover IS NOT NULL").
		Where("cover <> ''")

	tx.QueryStmt(q)
	for {
		var cover string
		if !tx.Scan(&cover) {
			break
		}
		// Strip extension to get normalized path
		normalized := stripExtension(cover)
		validPaths.coverPaths[normalized] = true
	}

	return validPaths
}

func stripExtension(path string) string {
	ext := filepath.Ext(path)
	if ext != "" {
		return path[:len(path)-len(ext)]
	}
	return path
}

func isPlaceholderFile(filename string) bool {
	return filename == "placeholder.png"
}

func walkAndIdentifyOrphans(basePath string, validPaths validPathSet, verbose bool) cleanupStats {
	stats := cleanupStats{
		orphanedFiles: make([]string, 0),
		errors:        make([]string, 0),
	}

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			stats.errors = append(stats.errors, fmt.Sprintf("Error accessing %s: %v", path, err))
			return nil // Continue walking
		}

		// Skip directories and the badges folder
		if info.IsDir() {
			if strings.Contains(path, "/badges") {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip placeholder files used during processing
		if isPlaceholderFile(info.Name()) {
			if verbose {
				log.Printf("Skipping placeholder file: %s\n", path)
			}
			return nil
		}

		// Get relative path from base
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			stats.errors = append(stats.errors, fmt.Sprintf("Error getting relative path for %s: %v", path, err))
			return nil
		}

		// Determine image type and check if orphaned
		isOrphaned := false
		imageType := ""

		if matches := albumPathRe.FindStringSubmatch(relPath); matches != nil {
			stats.albumsFound++
			imageType = "album"
			// Extract normalized path (remove size prefix and extension)
			normalizedPath := stripExtension(matches[2])

			// Check if it's a PNG preview file
			if filepath.Ext(relPath) == ".png" {
				// Only valid if this album has a preview extension
				if !validPaths.animatedAlbumPaths[normalizedPath] {
					isOrphaned = true
				}
			} else {
				// Regular album image
				if !validPaths.albumPaths[normalizedPath] {
					isOrphaned = true
				}
			}
		} else if matches := avatarPathRe.FindStringSubmatch(relPath); matches != nil {
			stats.avatarsFound++
			imageType = "avatar"
			normalizedPath := stripExtension(matches[2])
			if !validPaths.avatarPaths[normalizedPath] {
				isOrphaned = true
			}
		} else if matches := coverPathRe.FindStringSubmatch(relPath); matches != nil {
			stats.coversFound++
			imageType = "cover"
			normalizedPath := stripExtension(matches[2])
			if !validPaths.coverPaths[normalizedPath] {
				isOrphaned = true
			}
		} else {
			// Unknown file type, log if verbose
			if verbose {
				log.Printf("Skipping unknown file: %s\n", relPath)
			}
			return nil
		}

		if isOrphaned {
			if verbose {
				log.Printf("Orphaned %s: %s (%d bytes)\n", imageType, relPath, info.Size())
			}

			stats.orphanedFiles = append(stats.orphanedFiles, path)
			stats.totalSize += info.Size()

			switch imageType {
			case "album":
				stats.albumsOrphaned++
			case "avatar":
				stats.avatarsOrphaned++
			case "cover":
				stats.coversOrphaned++
			}
		}

		return nil
	})

	if err != nil {
		stats.errors = append(stats.errors, fmt.Sprintf("Error walking directory: %v", err))
	}

	return stats
}

func printStats(stats cleanupStats) {
	log.Println("\n" + strings.Repeat("=", 70))
	log.Println("CLEANUP STATISTICS")
	log.Println(strings.Repeat("=", 70))

	log.Printf("\nAlbums:\n")
	log.Printf("  Files scanned:  %d\n", stats.albumsFound)
	log.Printf("  Orphaned:       %d (%.1f%%)\n", stats.albumsOrphaned, percentage(stats.albumsOrphaned, stats.albumsFound))

	log.Printf("\nAvatars:\n")
	log.Printf("  Files scanned:  %d\n", stats.avatarsFound)
	log.Printf("  Orphaned:       %d (%.1f%%)\n", stats.avatarsOrphaned, percentage(stats.avatarsOrphaned, stats.avatarsFound))

	log.Printf("\nCovers:\n")
	log.Printf("  Files scanned:  %d\n", stats.coversFound)
	log.Printf("  Orphaned:       %d (%.1f%%)\n", stats.coversOrphaned, percentage(stats.coversOrphaned, stats.coversFound))

	log.Println("\n" + strings.Repeat("-", 70))
	totalFound := stats.albumsFound + stats.avatarsFound + stats.coversFound
	totalOrphaned := stats.albumsOrphaned + stats.avatarsOrphaned + stats.coversOrphaned
	log.Printf("TOTAL:\n")
	log.Printf("  Files scanned:  %d\n", totalFound)
	log.Printf("  Orphaned:       %d (%.1f%%)\n", totalOrphaned, percentage(totalOrphaned, totalFound))
	log.Printf("  Size:           %s\n", formatBytes(stats.totalSize))
	log.Println(strings.Repeat("=", 70))

	if len(stats.orphanedFiles) > 0 {
		log.Printf("\nSample orphaned files (showing up to 20):\n")
		for i, file := range stats.orphanedFiles {
			if i >= 20 {
				log.Printf("  ... and %d more files\n", len(stats.orphanedFiles)-20)
				break
			}
			info, err := os.Stat(file)
			if err == nil {
				log.Printf("  %d. %s (%s)\n", i+1, file, formatBytes(info.Size()))
			} else {
				log.Printf("  %d. %s\n", i+1, file)
			}
		}
	}
}

func confirmDeletion(stats cleanupStats) bool {
	log.Println("\n" + strings.Repeat("!", 70))
	log.Printf("WARNING: About to delete %d files (%s)\n", len(stats.orphanedFiles), formatBytes(stats.totalSize))
	log.Println(strings.Repeat("!", 70))
	log.Print("\nDelete these files? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func deleteOrphanedFiles(stats cleanupStats) {
	log.Printf("\nDeleting %d orphaned files...\n", len(stats.orphanedFiles))

	deletedCount := 0
	deletedSize := int64(0)
	errorCount := 0

	for i, file := range stats.orphanedFiles {
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("[%d/%d] Error stating %s: %v\n", i+1, len(stats.orphanedFiles), file, err)
			errorCount++
			continue
		}

		err = os.Remove(file)
		if err != nil {
			log.Printf("[%d/%d] Error deleting %s: %v\n", i+1, len(stats.orphanedFiles), file, err)
			errorCount++
		} else {
			deletedCount++
			deletedSize += info.Size()
			if (i+1)%100 == 0 || i+1 == len(stats.orphanedFiles) {
				log.Printf("[%d/%d] Deleted %s\n", i+1, len(stats.orphanedFiles), file)
			}
		}
	}

	log.Println("\n" + strings.Repeat("=", 70))
	log.Println("CLEANUP COMPLETE")
	log.Println(strings.Repeat("=", 70))
	log.Printf("Files deleted:   %d / %d\n", deletedCount, len(stats.orphanedFiles))
	log.Printf("Space freed:     %s\n", formatBytes(deletedSize))
	log.Printf("Errors:          %d\n", errorCount)
	log.Println(strings.Repeat("=", 70))
}

func percentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func findUnusedDatabaseImages(tx *database.AutoTx, imageFolder string) []unusedImage {
	// Find images that are:
	// 1. In the images table
	// 2. NOT in entry_images table (not attached to any entry)
	// 3. Created more than 6 months ago
	// 4. Not currently being processed
	q := sqlf.Select("images.id, images.path, images.extension, images.preview_extension, images.created_at").
		From("images").
		LeftJoin("entry_images", "entry_images.image_id = images.id").
		Where("entry_images.image_id IS NULL").
		Where("images.processing = false").
		Where("images.created_at < NOW() - INTERVAL '6 months'").
		OrderBy("images.created_at")

	tx.QueryStmt(q)

	var unused []unusedImage
	for {
		var img unusedImage
		if !tx.Scan(&img.id, &img.path, &img.extension, &img.previewExtension, &img.createdAt) {
			break
		}

		// Build list of files for this image (all size variants)
		img.files = buildImageFilePaths(imageFolder, img.path, img.extension, img.previewExtension)

		// Calculate total size
		for _, file := range img.files {
			info, err := os.Stat(file)
			if err == nil {
				img.totalSize += info.Size()
			}
		}

		unused = append(unused, img)
	}

	return unused
}

func buildImageFilePaths(baseFolder, path, extension, previewExtension string) []string {
	sizes := []string{"thumbnails", "small", "medium", "large"}
	var files []string

	for _, size := range sizes {
		// Regular image file
		file := filepath.Join(baseFolder, "albums", size, path+"."+extension)
		files = append(files, file)

		// Preview file if it exists (for animated images)
		if previewExtension != "" {
			previewFile := filepath.Join(baseFolder, "albums", size, path+"."+previewExtension)
			files = append(files, previewFile)
		}
	}

	return files
}

func printUnusedImageStats(unused []unusedImage) {
	if len(unused) == 0 {
		log.Println("\n✓ No unused images found!")
		return
	}

	totalSize := int64(0)
	totalFiles := 0
	for _, img := range unused {
		totalSize += img.totalSize
		totalFiles += len(img.files)
	}

	log.Println("\n" + strings.Repeat("=", 70))
	log.Println("UNUSED DATABASE IMAGES")
	log.Println(strings.Repeat("=", 70))
	log.Printf("\nFound %d unused images (created > 6 months ago, not in any entry)\n", len(unused))
	log.Printf("Total files: %d\n", totalFiles)
	log.Printf("Total size:  %s\n", formatBytes(totalSize))

	log.Printf("\nSample unused images (showing up to 20):\n")
	for i, img := range unused {
		if i >= 20 {
			log.Printf("  ... and %d more images\n", len(unused)-20)
			break
		}
		log.Printf("  %d. ID=%d, path=%s.%s, created=%s (%s, %d files)\n",
			i+1, img.id, img.path, img.extension, img.createdAt,
			formatBytes(img.totalSize), len(img.files))
	}
	log.Println(strings.Repeat("=", 70))
}

func confirmUnusedDeletion(unused []unusedImage) bool {
	totalSize := int64(0)
	totalFiles := 0
	for _, img := range unused {
		totalSize += img.totalSize
		totalFiles += len(img.files)
	}

	log.Println("\n" + strings.Repeat("!", 70))
	log.Printf("WARNING: About to delete %d database records and %d files (%s)\n",
		len(unused), totalFiles, formatBytes(totalSize))
	log.Println("This will delete data from BOTH disk AND database!")
	log.Println(strings.Repeat("!", 70))
	log.Print("\nDelete these images? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func deleteUnusedImages(tx *database.AutoTx, unused []unusedImage) {
	log.Printf("\nDeleting %d unused images...\n", len(unused))

	deletedFiles := 0
	deletedRecords := 0
	deletedSize := int64(0)
	errorCount := 0

	for i, img := range unused {
		// Delete files from disk
		filesDeleted := 0
		for _, file := range img.files {
			info, err := os.Stat(file)
			if err != nil {
				// File doesn't exist, skip silently
				continue
			}

			err = os.Remove(file)
			if err != nil {
				log.Printf("[%d/%d] Error deleting file %s: %v\n", i+1, len(unused), file, err)
				errorCount++
			} else {
				filesDeleted++
				deletedFiles++
				deletedSize += info.Size()
			}
		}

		// Delete from database
		q := sqlf.DeleteFrom("images").
			Where("id = ?", img.id)

		tx.ExecStmt(q)
		if tx.HasQueryError() {
			log.Printf("[%d/%d] Error deleting image ID=%d from database: %v\n",
				i+1, len(unused), img.id, tx.Error())
			errorCount++
		} else {
			deletedRecords++
		}

		if (i+1)%10 == 0 || i+1 == len(unused) {
			log.Printf("[%d/%d] Deleted image ID=%d (%d files, %s)\n",
				i+1, len(unused), img.id, filesDeleted, formatBytes(img.totalSize))
		}
	}

	log.Println("\n" + strings.Repeat("=", 70))
	log.Println("UNUSED IMAGE CLEANUP COMPLETE")
	log.Println(strings.Repeat("=", 70))
	log.Printf("Database records deleted: %d / %d\n", deletedRecords, len(unused))
	log.Printf("Files deleted:            %d\n", deletedFiles)
	log.Printf("Space freed:              %s\n", formatBytes(deletedSize))
	log.Printf("Errors:                   %d\n", errorCount)
	log.Println(strings.Repeat("=", 70))
}
