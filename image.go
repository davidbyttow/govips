package govips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

// Immutable Image structure that represents an image in memory.
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

// Loads an image buffer from disk and creates a new Image
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

// Loads an image buffer and creates a new Image
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

// Returns the width of this image
func (i *Image) Width() int {
	return int(i.image.Xsize)
}

// Returns the height of this iamge
func (i *Image) Height() int {
	return int(i.image.Ysize)
}

// Returns the number of bands for this image
func (i *Image) Bands() int {
	return int(i.image.Bands)
}

// Returns the X resolution
func (i *Image) ResX() float64 {
	return float64(i.image.Xres)
}

// Returns the Y resolution
func (i *Image) ResY() float64 {
	return float64(i.image.Yres)
}

// Returns the X offset
func (i *Image) OffsetX() int {
	return int(i.image.Xoffset)
}

// Returns the Y offset
func (i *Image) OffsetY() int {
	return int(i.image.Yoffset)
}

// Returns the current band format
func (i *Image) BandFormat() BandFormat {
	return BandFormat(int(i.image.BandFmt))
}

// Returns the image coding
func (i *Image) Coding() Coding {
	return Coding(int(i.image.Coding))
}

// Returns the current interpretation
func (i *Image) Interpretation() Interpretation {
	return Interpretation(int(i.image.Type))
}

// Writes the image to memory in VIPs format and returns the raw bytes, useful for storage.
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

// Writes the image to a buffer in a format represented by the given suffix (e.g., .jpeg)
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

// Writes the image to a file on disk based on the format specified in the path
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
