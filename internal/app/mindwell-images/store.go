package images

import (
	"io"
	"math"
	"os"
	"strconv"

	"github.com/sevings/mindwell-server/utils"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/sevings/mindwell-server/models"
)

const imageExtensionJpg = "jpg"
const imageExtensionGif = "gif"

type pageList []*vips.ImageRef

func newPageList(img *vips.ImageRef) pageList {
	return []*vips.ImageRef{img}
}

func (pl *pageList) size() int {
	return len(*pl)
}

func (pl *pageList) width() int {
	return pl.first().Width()
}

func (pl *pageList) height() int {
	return pl.first().Height()
}

func (pl *pageList) first() *vips.ImageRef {
	return (*pl)[0]
}

func (pl *pageList) split() error {
	img := pl.first()

	pages := img.Pages()
	if pages == 1 {
		return nil
	}

	width := img.Width()
	height := img.PageHeight()

	err := img.SetPages(1)
	if err != nil {
		return err
	}

	err = img.SetPageHeight(img.Height())
	if err != nil {
		return err
	}

	*pl = make([]*vips.ImageRef, 0, img.Pages())
	for i := 0; i < pages; i++ {
		var page *vips.ImageRef
		page, err = img.Copy()
		if err != nil {
			return err
		}

		err = page.ExtractArea(0, i*height, width, height)
		if err != nil {
			return err
		}

		*pl = append(*pl, page)
	}

	return nil
}

func (pl *pageList) join() error {
	if pl.size() == 1 {
		return nil
	}

	img := pl.first()

	err := img.ArrayJoin((*pl)[1:], 1)
	if err != nil {
		return err
	}

	err = img.SetPages(pl.size())
	if err != nil {
		return err
	}

	err = img.SetPageHeight(img.Height() / pl.size())
	if err != nil {
		return err
	}

	return nil
}

func (pl *pageList) copy() (pageList, error) {
	pages := make([]*vips.ImageRef, 0, pl.size())

	for _, p := range *pl {
		page, err := p.Copy()
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func (pl *pageList) close() {
	for _, p := range *pl {
		p.Close()
	}
}

func (pl *pageList) thumbnail(width, height int) error {
	crop := vips.InterestingAttention
	if pl.size() > 1 {
		crop = vips.InterestingCentre
	}

	for _, page := range *pl {
		err := page.ThumbnailWithSize(width, height, crop, vips.SizeDown)
		if err != nil {
			return err
		}
	}

	return nil
}

type imageStore struct {
	savePath  string
	saveName  string
	extension string
	pages     pageList
	mi        *MindwellImages
	err       error
}

type storeError string

func (se storeError) Error() string {
	return string(se)
}

func newImageStore(mi *MindwellImages) *imageStore {
	name := utils.GenerateString(8)
	path := name[:1] + "/" + name[:2] + "/"

	return &imageStore{
		savePath: path,
		saveName: name[2:],
		mi:       mi,
	}
}

func (is *imageStore) Destroy() {
	is.pages.close()
}

func (is *imageStore) Error() error {
	return is.err
}

func (is *imageStore) Folder() string {
	return is.mi.Folder()
}

func (is *imageStore) FileName() string {
	return is.savePath + is.saveName
}

func (is *imageStore) FileExtension() string {
	return is.extension
}

func (is *imageStore) IsAnimated() bool {
	return is.pages.size() > 1
}

func (is *imageStore) ReadImage(r io.ReadCloser) {
	defer r.Close()

	blob := make([]byte, 5*1024*1024)
	_, is.err = r.Read(blob)
	if is.err != nil {
		return
	}

	params := &vips.ImportParams{}
	params.AutoRotate.Set(true)
	params.NumPages.Set(-1)

	var img *vips.ImageRef
	img, is.err = vips.LoadImageFromBuffer(blob, params)
	if is.err != nil {
		return
	}

	is.pages = newPageList(img)

	if img.Pages() == 1 {
		is.extension = imageExtensionJpg
	} else {
		is.extension = imageExtensionGif
	}
}

func (is *imageStore) SetID(id int64) {
	is.saveName = is.saveName + strconv.FormatInt(id, 32)
}

func (is *imageStore) PrepareImage() {
	is.err = is.pages.split()
}

func (is *imageStore) Fill(size int, folder string) *models.ImageSize {
	return is.FillRect(size, size, folder)
}

func (is *imageStore) FillRect(width, height int, folder string) *models.ImageSize {
	if is.err != nil {
		return nil
	}

	originWidth := is.pages.width()
	originHeight := is.pages.height()

	ratio := float64(width) / float64(height)
	originRatio := float64(originWidth) / float64(originHeight)

	cropWidth, cropHeight := originWidth, originHeight

	if ratio < originRatio {
		cropWidth = int(math.RoundToEven(float64(originHeight) * ratio))
	} else {
		cropHeight = int(math.RoundToEven(float64(originWidth) / ratio))
	}

	if width > originWidth || height > originHeight {
		width, height = cropWidth, cropHeight
	}

	var pages pageList
	pages, is.err = is.pages.copy()
	if is.err != nil {
		return nil
	}

	defer pages.close()

	is.err = pages.thumbnail(width, height)
	if is.err != nil {
		return nil
	}

	return is.saveImageSize(pages, folder)
}

func (is *imageStore) Fit(size int, folder string) *models.ImageSize {
	return is.FitRect(size, size, folder)
}

func (is *imageStore) FitRect(width, height int, folder string) *models.ImageSize {
	if is.err != nil {
		return nil
	}

	var pages pageList
	pages, is.err = is.pages.copy()
	if is.err != nil {
		return nil
	}

	defer pages.close()

	originHeight := is.pages.height()
	originWidth := is.pages.width()

	if originHeight < height && originWidth < width {
		return is.saveImageSize(pages, folder)
	}

	ratio := float64(width) / float64(height)
	originRatio := float64(originWidth) / float64(originHeight)

	if ratio > originRatio {
		width = int(math.RoundToEven(float64(height) * originRatio))
	} else {
		height = int(math.RoundToEven(float64(width) / originRatio))
	}

	is.err = pages.thumbnail(width, height)
	if is.err != nil {
		return nil
	}

	return is.saveImageSize(pages, folder)
}

func (is *imageStore) FolderRemove(folder, path string) {
	if is.err != nil {
		return
	}

	if len(path) == 0 {
		return
	}

	is.err = os.Remove(is.Folder() + folder + "/" + path)
}

func (is *imageStore) saveImageSize(pages pageList, folder string) *models.ImageSize {
	imgSize := &models.ImageSize{
		Width:  int64(pages.width()),
		Height: int64(pages.height()),
	}

	if is.IsAnimated() {
		imgSize.Preview = is.saveImage(pages.first(), folder, imageExtensionJpg)

		is.err = pages.join()
		if is.err != nil {
			return nil
		}
	}

	imgSize.URL = is.saveImage(pages.first(), folder, is.extension)

	return imgSize
}

func (is *imageStore) saveImage(img *vips.ImageRef, folder, extension string) string {
	path := folder + "/" + is.savePath
	is.err = os.MkdirAll(is.Folder()+path, 0777)
	if is.err != nil {
		return ""
	}

	fileName := path + is.saveName + "." + extension

	if extension == imageExtensionJpg {
		is.saveJpeg(img, fileName)
	} else {
		is.saveGif(img, fileName)
	}

	if is.err != nil {
		return ""
	}

	return is.mi.BaseURL() + fileName
}

func (is *imageStore) saveJpeg(img *vips.ImageRef, fileName string) {
	params := vips.NewJpegExportParams()
	params.StripMetadata = true
	params.Quality = 80
	params.Interlace = true

	var buf []byte
	buf, _, is.err = img.ExportJpeg(params)
	if is.err != nil {
		return
	}

	is.err = os.WriteFile(is.Folder()+fileName, buf, 0644)
}

func (is *imageStore) saveGif(img *vips.ImageRef, fileName string) {
	params := vips.NewGifExportParams()
	params.StripMetadata = true

	var buf []byte
	buf, _, is.err = img.ExportGIF(params)
	if is.err != nil {
		return
	}

	is.err = os.WriteFile(is.Folder()+fileName, buf, 0644)
}
