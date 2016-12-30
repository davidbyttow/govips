package gimage

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

import (
	"runtime"
)

type Image struct {
	image *C.VipsImage
}

func newImage(vipsImage *C.VipsImage) *Image {
	image := &Image{
		image: vipsImage,
	}
	runtime.SetFinalizer(image, finalizeImage)
	return image
}

func finalizeImage(i *Image) {
	C.g_object_unref(C.gpointer(i.image))
}

func LoadBuffer(bytes []byte) (*Image, error) {
	blob := NewBlob(bytes)
	image := JpegloadBuffer(blob, nil)
	return image, nil
}

func (i *Image) Width() int {
	return int(i.image.Xsize)
}

func (i *Image) Height() int {
	return int(i.image.Ysize)
}

func (i *Image) Bands() int {
	return int(i.image.Bands)
}

func (i *Image) ResX() float64 {
	return float64(i.image.Xres)
}

func (i *Image) ResY() float64 {
	return float64(i.image.Yres)
}

func (i *Image) OffsetX() int {
	return int(i.image.Xoffset)
}

func (i *Image) OffsetY() int {
	return int(i.image.Yoffset)
}

func (i *Image) BandFormat() BandFormat {
	return BandFormat(int(i.image.BandFmt))
}

func (i *Image) Coding() Coding {
	return Coding(int(i.image.Coding))
}

func (i *Image) Interpretation() Interpretation {
	return Interpretation(int(i.image.Type))
}
