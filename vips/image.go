package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

type Image struct {
	image     *C.VipsImage
	imageType ImageType
}

func newImage(image *C.VipsImage, imageType ImageType) Image {
	return Image{
		image:     image,
		imageType: ImageTypeJpeg,
	}
}

func NewImage(bytes []byte) (Image, error) {
	// TODO(d): Load from vips
	image, imageType, err := loadImage(bytes)
	if err != nil {
		return Image{}, err
	}
	return newImage(image, imageType), nil
}

func NewJpegImage(bytes []byte, shrinkFactor int) (Image, error) {
	image, err := loadJpegImage(bytes, shrinkFactor)
	if err != nil {
		return Image{}, err
	}
	return newImage(image, ImageTypeJpeg), nil
}

func NewWebpImage(bytes []byte, shrinkFactor int) (Image, error) {
	image, err := loadWebpImage(bytes, shrinkFactor)
	if err != nil {
		return Image{}, err
	}
	return newImage(image, ImageTypeWebp), nil
}

func (i Image) Type() ImageType {
	return i.imageType
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

func (i Image) ShrinkH(shrink int) Image {
	return newImage(i.image, i.imageType)
}

func (i Image) ShrinkV(shrink int) Image {
	return newImage(i.image, i.imageType)
}

// TODO(d): BandsFormat

// TODO(d): Coding

// TODO(d): Interpretation

// TODO(d): GuessInterpretation
