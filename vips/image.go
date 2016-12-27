package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type Image struct {
	image       *C.VipsImage
	imageType   ImageType
	sourceBytes []byte
}

func newImage(vipsImage *C.VipsImage, imageType ImageType, sourceBytes []byte) *Image {
	image := &Image{
		image:       vipsImage,
		imageType:   imageType,
		sourceBytes: sourceBytes,
	}
	runtime.SetFinalizer(image, freeImage)
	return image
}

func freeImage(image *Image) {
	C.free(unsafe.Pointer(image.image))
}

func LoadBuffer(bytes []byte) (*Image, error) {
	image, imageType, err := loadBuffer(bytes)
	if err != nil {
		return nil, err
	}
	return newImage(image, imageType, bytes), nil
}

func LoadJpegBuffer(bytes []byte, shrinkFactor int) (*Image, error) {
	image, err := loadJpegBuffer(bytes, shrinkFactor)
	if err != nil {
		return nil, err
	}
	return newImage(image, ImageTypeJpeg, bytes), nil
}

func LoadWebpBuffer(bytes []byte, shrinkFactor int) (*Image, error) {
	image, err := loadWebpBuffer(bytes, shrinkFactor)
	if err != nil {
		return nil, err
	}
	return newImage(image, ImageTypeWebp, bytes), nil
}

func (i Image) Type() ImageType {
	return i.imageType
}

func (i Image) SourceBytes() []byte {
	return i.sourceBytes
}

func (i Image) Width() int {
	return int(i.image.Xsize)
}

func (i Image) Height() int {
	return int(i.image.Ysize)
}

func (i Image) Bands() int {
	return int(C.vips_image_get_bands(i.image))
}

func (i Image) ResX() float64 {
	return float64(C.vips_image_get_xres(i.image))
}

func (i Image) ResY() float64 {
	return float64(C.vips_image_get_yres(i.image))
}

func (i Image) OffsetX() float64 {
	return float64(C.vips_image_get_xoffset(i.image))
}

func (i Image) OffsetY() float64 {
	return float64(C.vips_image_get_yoffset(i.image))
}

// TODO(d): BandsFormat

// TODO(d): Coding

// TODO(d): Interpretation

// TODO(d): GuessInterpretation

func loadBuffer(buf []byte) (*C.VipsImage, ImageType, error) {
	imageType := DetermineImageType(buf)

	if imageType == ImageTypeUnknown {
		return nil, ImageTypeUnknown, errors.New("Unsupported image format")
	}

	var image *C.VipsImage
	err := C.init_image(cPtr(buf),
		C.size_t(len(buf)),
		C.int(imageType),
		&image)

	if err != 0 {
		return nil, ImageTypeUnknown, handleVipsError()
	}
	return image, imageType, nil
}

func loadJpegBuffer(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}

func loadWebpBuffer(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}
