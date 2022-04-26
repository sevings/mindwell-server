package images

import (
	"io"
	"math"
	"os"
	"strconv"

	"github.com/sevings/mindwell-server/utils"

	"github.com/sevings/mindwell-server/models"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const imageExtensionJpg = "jpg"
const imageExtensionGif = "gif"

type imageStore struct {
	savePath  string
	saveName  string
	extension string
	mw        *imagick.MagickWand
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
		mw:       imagick.NewMagickWand(),
		mi:       mi,
	}
}

func (is *imageStore) Destroy() {
	is.mw.Destroy()
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
	return is.extension == imageExtensionGif
}

func (is *imageStore) ReadImage(r io.ReadCloser) {
	defer r.Close()

	blob := make([]byte, 5*1024*1024)
	_, is.err = r.Read(blob)
	if is.err != nil {
		return
	}

	is.err = is.mw.ReadImageBlob(blob)
	if is.err != nil {
		return
	}

	if is.mw.GetNumberImages() == 1 {
		is.extension = imageExtensionJpg
	} else {
		is.extension = imageExtensionGif
	}
}

func (is *imageStore) SetID(id int64) {
	is.saveName = is.saveName + strconv.FormatInt(id, 32)
}

func (is *imageStore) PrepareImage() {
	if is.extension == imageExtensionGif {
		is.prepareGif()
	} else {
		is.prepareJpeg()
	}

	is.err = is.mw.StripImage()
}

func (is *imageStore) prepareGif() {
	wand := is.mw.CoalesceImages()
	is.mw.Destroy()
	is.mw = wand
}

func (is *imageStore) prepareJpeg() {
	orient := is.mw.GetImageOrientation()
	if orient == imagick.ORIENTATION_TOP_LEFT {
		return
	}

	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	switch orient {
	case imagick.ORIENTATION_TOP_RIGHT:
		is.mw.FlopImage()
	case imagick.ORIENTATION_BOTTOM_RIGHT:
		is.mw.RotateImage(pw, 180)
	case imagick.ORIENTATION_BOTTOM_LEFT:
		is.mw.FlopImage()
		is.mw.RotateImage(pw, 180)
	case imagick.ORIENTATION_LEFT_TOP:
		is.mw.FlopImage()
		is.mw.RotateImage(pw, -90)
	case imagick.ORIENTATION_RIGHT_TOP:
		is.mw.RotateImage(pw, 90)
	case imagick.ORIENTATION_RIGHT_BOTTOM:
		is.mw.FlopImage()
		is.mw.RotateImage(pw, 90)
	case imagick.ORIENTATION_LEFT_BOTTOM:
		is.mw.RotateImage(pw, -90)
	}
}

func (is *imageStore) Fill(size uint, folder string) *models.ImageSize {
	return is.FillRect(size, size, folder)
}

func (is *imageStore) FillRect(width, height uint, folder string) *models.ImageSize {
	if is.err != nil {
		return nil
	}

	originWidth := is.mw.GetImageWidth()
	originHeight := is.mw.GetImageHeight()

	ratio := float64(width) / float64(height)
	originRatio := float64(originWidth) / float64(originHeight)

	crop := math.Abs(ratio-originRatio) > 0.01

	cropWidth, cropHeight := originWidth, originHeight

	if ratio < originRatio {
		cropWidth = uint(float64(originHeight) * ratio)
	} else {
		cropHeight = uint(float64(originWidth) / ratio)
	}

	if width > originWidth || height > originHeight {
		width, height = cropWidth, cropHeight
	}

	x := int(originWidth-cropWidth) / 2
	y := int(originHeight-cropHeight) / 2

	wand := is.mw.Clone()
	defer wand.Destroy()

	wand.ResetIterator()
	for wand.NextImage() {
		if crop {
			is.err = wand.CropImage(cropWidth, cropHeight, x, y)
			if is.err != nil {
				return nil
			}
		}

		is.err = wand.ThumbnailImage(width, height)
		if is.err != nil {
			return nil
		}
	}

	return is.saveImageSize(wand, folder, width, height)
}

func (is *imageStore) Fit(size uint, folder string) *models.ImageSize {
	return is.FitRect(size, size, folder)
}

func (is *imageStore) FitRect(width, height uint, folder string) *models.ImageSize {
	if is.err != nil {
		return nil
	}

	wand := is.mw.Clone()
	defer wand.Destroy()

	originHeight := is.mw.GetImageHeight()
	originWidth := is.mw.GetImageWidth()

	if originHeight < height && originWidth < width {
		return is.saveImageSize(wand, folder, width, height)
	}

	ratio := float64(width) / float64(height)
	originRatio := float64(originWidth) / float64(originHeight)

	if ratio > originRatio {
		width = uint(float64(height) * originRatio)
	} else {
		height = uint(float64(width) / originRatio)
	}

	wand.ResetIterator()
	for wand.NextImage() {
		is.err = wand.ResizeImage(width, height, imagick.FILTER_CUBIC, 0.5)
		if is.err != nil {
			return nil
		}
	}

	return is.saveImageSize(wand, folder, width, height)
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

func (is *imageStore) saveImageSize(wand *imagick.MagickWand, folder string, width, height uint) *models.ImageSize {
	img := &models.ImageSize{
		Width:  int64(width),
		Height: int64(height),
		URL:    is.saveImage(wand, folder, is.extension),
	}

	if is.extension == imageExtensionGif {
		img.Preview = is.saveImage(wand, folder, imageExtensionJpg)
	}

	return img
}

func (is *imageStore) saveImage(wand *imagick.MagickWand, folder, extension string) string {
	path := folder + "/" + is.savePath
	is.err = os.MkdirAll(is.Folder()+path, 0777)
	if is.err != nil {
		return ""
	}

	fileName := path + is.saveName + "." + extension

	if extension == imageExtensionJpg {
		is.saveJpeg(wand, fileName)
	} else {
		is.saveGif(wand, fileName)
	}

	if is.err != nil {
		return ""
	}

	return is.mi.BaseURL() + fileName
}

func (is *imageStore) saveJpeg(wand *imagick.MagickWand, fileName string) {
	is.err = wand.SetImageCompression(imagick.COMPRESSION_JPEG)
	if is.err != nil {
		return
	}

	is.err = wand.SetImageInterlaceScheme(imagick.INTERLACE_JPEG)
	if is.err != nil {
		return
	}

	is.err = wand.SetImageCompressionQuality(80)
	if is.err != nil {
		return
	}

	is.err = wand.WriteImage(is.Folder() + fileName)
}

func (is *imageStore) saveGif(wand *imagick.MagickWand, fileName string) {
	wand = wand.DeconstructImages()
	defer wand.Destroy()

	is.err = wand.OptimizeImageTransparency()
	if is.err != nil {
		return
	}

	is.err = wand.SetImageCompression(imagick.COMPRESSION_LZW)
	if is.err != nil {
		return
	}

	is.err = wand.SetImageInterlaceScheme(imagick.INTERLACE_GIF)
	if is.err != nil {
		return
	}

	is.err = wand.WriteImages(is.Folder()+fileName, true)
}
