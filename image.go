package gimage

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

import (
	"runtime"
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
	runtime.SetFinalizer(image, finalizeImage)
	return image
}

func finalizeImage(i *Image) {
	C.g_object_unref(C.gpointer(i.image))
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

func (i Image) SetImage(vipsImage *C.VipsImage) *Image {
	return newImage(vipsImage, i.Type(), i.SourceBytes())
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
	return int(i.image.Bands)
}

func (i Image) ResX() float64 {
	return float64(i.image.Xres)
}

func (i Image) ResY() float64 {
	return float64(i.image.Yres)
}

func (i Image) OffsetX() int {
	return int(i.image.Xoffset)
}

func (i Image) OffsetY() int {
	return int(i.image.Yoffset)
}

// TODO(d): BandsFormat

// TODO(d): Coding

// TODO(d): Interpretation

// TODO(d): GuessInterpretation
