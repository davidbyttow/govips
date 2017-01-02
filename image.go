package gimage

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
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

func NewImageFromFile(path string, options *Options) (*Image, error) {
	fileName, optionString := vipsFilenameSplit8(path)

	operationName, err := vipsForeignFindLoad(fileName)
	if err != nil {
		return nil, ErrUnsupportedImageFormat
	}

	var out *Image
	if options == nil {
		options = NewOptions().
			SetString("filename", fileName).
			SetImageOut("out", &out)
	}

	if err := CallOperation(operationName, options, optionString); err != nil {
		return nil, err
	}
	return out, nil
}

func NewImageFromBuffer(bytes []byte, options *Options) (*Image, error) {
	operationName, err := vipsForeignFindLoadBuffer(bytes)
	if err != nil {
		return nil, err
	}

	var out *Image
	blob := NewBlob(bytes)
	if options == nil {
		options = NewOptions().
			SetBlob("buffer", blob).
			SetImageOut("out", &out)
	}

	if err := CallOperation(operationName, options, ""); err != nil {
		return nil, err
	}

	return out, nil
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

func (i *Image) ToBytes() ([]byte, error) {
	var size C.size_t
	c_data := C.vips_image_write_to_memory(i.image, &size)
	if c_data == nil {
		return nil, errors.New("Failed to write image to memory")
	}
	defer C.free(c_data)

	bytes := C.GoBytes(unsafe.Pointer(c_data), C.int(size))
	return bytes, nil
}

func (i *Image) WriteToBuffer(suffix string, options *Options) ([]byte, error) {
	fileName, optionString := vipsFilenameSplit8(suffix)
	operationName, err := vipsForeignFindSaveBuffer(fileName)
	if err != nil {
		return nil, err
	}
	var blob *Blob
	if options == nil {
		options = NewOptions().
			SetImage("in", i).
			SetBlobOut("buffer", &blob)
	}
	err = CallOperation(operationName, options, optionString)
	if err != nil {
		return nil, err
	}
	if blob != nil {
		return blob.ToBytes(), nil
	}
	return nil, nil
}

func (i *Image) WriteToFile(path string, options *Options) error {
	fileName, optionString := vipsFilenameSplit8(path)
	operationName, err := vipsForeignFindSave(fileName)
	if err != nil {
		return err
	}
	if options == nil {
		options = NewOptions().
			SetImage("in", i).
			SetString("filename", fileName)
	}
	return CallOperation(operationName, options, optionString)
}
