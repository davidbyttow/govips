package gimage

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import "runtime"

var STRING_BUFFER = fixedString(4096)

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
	c_path := C.CString(path)
	defer freeCString(c_path)

	c_filename := C.CString(STRING_BUFFER)
	defer freeCString(c_filename)

	c_optionString := C.CString(STRING_BUFFER)
	defer freeCString(c_optionString)

	C.filename_split8(c_path, c_filename, c_optionString)

	c_operationName := C.vips_foreign_find_load(c_filename)
	if c_operationName == nil {
		return nil, ErrUnsupportedImageFormat
	}

	var out *Image
	if options == nil {
		options = NewOptions().
			SetString("filename", C.GoString(c_filename)).
			SetImageOut("out", &out)
	}

	operationName := C.GoString(c_operationName)
	if err := Call(operationName, options); err != nil {
		return nil, err
	}

	return out, nil
}

func NewImageFromBuffer(bytes []byte, options *Options) (*Image, error) {
	c_operationName := C.vips_foreign_find_load_buffer(
		cPtr(bytes),
		C.size_t(len(bytes)))
	if c_operationName == nil {
		return nil, ErrUnsupportedImageFormat
	}

	var out *Image
	blob := NewBlob(bytes)
	if options == nil {
		options = NewOptions().
			SetBlob("buffer", blob).
			SetImageOut("out", &out)
	}

	operationName := C.GoString(c_operationName)
	if err := Call(operationName, options); err != nil {
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
