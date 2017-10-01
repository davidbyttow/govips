package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type ImageRef struct {
	image *C.VipsImage
}

type ExportOptions struct {
	Type           ImageType
	Quality        int
	Compression    int
	Interlaced     bool
	StripProfile   bool
	StripMetadata  bool
	Interpretation Interpretation
}

func NewStreamFromBuffer(bytes []byte, fn func(*ImageRef) error) error {
	defer ShutdownThread()
	image, err := OpenFromBuffer(bytes)
	if err != nil {
		return err
	}
	defer image.Close()
	return fn(image)
}

// OpenFromBuffer loads an image buffer and creates a new Image
func OpenFromBuffer(bytes []byte) (*ImageRef, error) {
	startupIfNeeded()

	image, _, err := vipsLoadFromBuffer(bytes)
	if err != nil {
		return nil, err
	}

	return newImageRef(image), nil
}

func newImageRef(vipsImage *C.VipsImage) *ImageRef {
	stream := &ImageRef{
		image: vipsImage,
	}
	runtime.SetFinalizer(stream, finalizeStream)
	return stream
}

func finalizeStream(os *ImageRef) {
	os.Close()
}

func (os *ImageRef) SetImage(image *C.VipsImage) {
	if os.image != nil {
		C.g_object_unref(C.gpointer(os.image))
	}
	os.image = image
}

func (os *ImageRef) Close() {
	os.SetImage(nil)
}

// Width returns the width of this image
func (os *ImageRef) Width() int {
	return int(os.image.Xsize)
}

// Height returns the height of this iamge
func (os *ImageRef) Height() int {
	return int(os.image.Ysize)
}

// Bands returns the number of bands for this image
func (os *ImageRef) Bands() int {
	return int(os.image.Bands)
}

// ResX returns the X resolution
func (os *ImageRef) ResX() float64 {
	return float64(os.image.Xres)
}

// ResY returns the Y resolution
func (os *ImageRef) ResY() float64 {
	return float64(os.image.Yres)
}

// OffsetX returns the X offset
func (os *ImageRef) OffsetX() int {
	return int(os.image.Xoffset)
}

// OffsetY returns the Y offset
func (os *ImageRef) OffsetY() int {
	return int(os.image.Yoffset)
}

// BandFormat returns the current band format
func (os *ImageRef) BandFormat() BandFormat {
	return BandFormat(int(os.image.BandFmt))
}

// Coding returns the image coding
func (os *ImageRef) Coding() Coding {
	return Coding(int(os.image.Coding))
}

// Interpretation returns the current interpretation
func (os *ImageRef) Interpretation() Interpretation {
	return Interpretation(int(os.image.Type))
}

// ToBytes writes the image to memory in VIPs format and returns the raw bytes, useful for storage.
func (os *ImageRef) ToBytes() ([]byte, error) {
	var cSize C.size_t
	cData := C.vips_image_write_to_memory(os.image, &cSize)
	if cData == nil {
		return nil, errors.New("Failed to write image to memory")
	}
	defer C.free(cData)

	bytes := C.GoBytes(unsafe.Pointer(cData), C.int(cSize))
	return bytes, nil
}

// WriteToBuffer writes the image to a buffer in a format represented by the given suffix (e.g., .jpeg)
func (os *ImageRef) Export(options ExportOptions) ([]byte, error) {
	return vipsExportBuffer(os.image, &options)
}
