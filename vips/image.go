package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"errors"
	"unsafe"
)

type Image struct {
	image       *C.VipsImage
	imageType   ImageType
	sourceBytes []byte
}

func newImage(image *C.VipsImage, imageType ImageType, sourceBytes []byte) *Image {
	return &Image{
		image:       image,
		imageType:   ImageTypeJpeg,
		sourceBytes: sourceBytes,
	}
}

func NewImage(bytes []byte) (*Image, error) {
	// TODO(d): Load from vips
	image, imageType, err := loadImage(bytes)
	if err != nil {
		return nil, err
	}
	return newImage(image, imageType, bytes), nil
}

func NewJpegImage(bytes []byte, shrinkFactor int) (*Image, error) {
	image, err := loadJpegImage(bytes, shrinkFactor)
	if err != nil {
		return nil, err
	}
	return newImage(image, ImageTypeJpeg, bytes), nil
}

func NewWebpImage(bytes []byte, shrinkFactor int) (*Image, error) {
	image, err := loadWebpImage(bytes, shrinkFactor)
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
	return int(C.vips_image_get_width(i.image))
}

func (i Image) Height() int {
	return int(C.vips_image_get_height(i.image))
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

func loadImage(buf []byte) (*C.VipsImage, ImageType, error) {
	var image *C.VipsImage
	imageType := determineImageType(buf)

	if imageType == ImageTypeUnknown {
		return nil, ImageTypeUnknown, errors.New("Unsupported image format")
	}

	imageBuf := unsafe.Pointer(&buf[0])
	length := C.size_t(len(buf))

	err := C.init_image(imageBuf, length, C.int(imageType), &image)
	if err != 0 {
		return nil, ImageTypeUnknown, handleVipsError()
	}

	return image, imageType, nil
}

func loadJpegImage(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}

func loadWebpImage(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}
