package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"errors"
	"io"
	"io/ioutil"
	"runtime"
	"unsafe"
)

// ImageRef contains a libvips image and manages its lifecycle. You should
// close an image when done or it will leak until the next GC
type ImageRef struct {
	image  *C.VipsImage
	format ImageType

	// NOTE(d): We keep a reference to this so that the input buffer is
	// never garbage collected during processing. Some image loaders use random
	// access transcoding and therefore need the original buffer to be in memory.
	buf []byte
}

// LoadImage loads an ImageRef from the given reader
func LoadImage(r io.Reader) (*ImageRef, error) {
	startupIfNeeded()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewImageFromBuffer(buf)
}

// NewImageFromFile loads an image from file and creates a new ImageRef
func NewImageFromFile(file string) (*ImageRef, error) {
	startupIfNeeded()

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewImageFromBuffer(buf)
}

// NewImageFromBuffer loads an image buffer and creates a new Image
func NewImageFromBuffer(buf []byte) (*ImageRef, error) {
	startupIfNeeded()

	image, format, err := vipsLoadFromBuffer(buf)
	if err != nil {
		return nil, err
	}

	ref := NewImageRef(image, format)
	ref.buf = buf
	return ref, nil
}

func NewImageRef(vipsImage *C.VipsImage, format ImageType) *ImageRef {
	stream := &ImageRef{
		image:  vipsImage,
		format: format,
	}
	runtime.SetFinalizer(stream, finalizeImage)
	return stream
}

func finalizeImage(ref *ImageRef) {
	ref.Close()
}

// SetImage resets the image for this image and frees the previous one
func (ref *ImageRef) SetImage(image *C.VipsImage) {
	if ref.image != nil {
		defer C.g_object_unref(C.gpointer(ref.image))
	}
	ref.image = image
}

// Format returns the initial format of the vips image when loaded
func (ref *ImageRef) Format() ImageType {
	return ref.format
}

// Close closes an image and frees internal memory associated with it
func (ref *ImageRef) Close() {
	ref.SetImage(nil)
	ref.buf = nil
}

// Image returns a handle to the internal vips image, just in case
func (ref *ImageRef) Image() *C.VipsImage {
	return ref.image
}

// Width returns the width of this image
func (ref *ImageRef) Width() int {
	return int(ref.image.Xsize)
}

// Height returns the height of this iamge
func (ref *ImageRef) Height() int {
	return int(ref.image.Ysize)
}

// Bands returns the number of bands for this image
func (ref *ImageRef) Bands() int {
	return int(ref.image.Bands)
}

// ResX returns the X resolution
func (ref *ImageRef) ResX() float64 {
	return float64(ref.image.Xres)
}

// ResY returns the Y resolution
func (ref *ImageRef) ResY() float64 {
	return float64(ref.image.Yres)
}

// OffsetX returns the X offset
func (ref *ImageRef) OffsetX() int {
	return int(ref.image.Xoffset)
}

// OffsetY returns the Y offset
func (ref *ImageRef) OffsetY() int {
	return int(ref.image.Yoffset)
}

// BandFormat returns the current band format
func (ref *ImageRef) BandFormat() BandFormat {
	return BandFormat(int(ref.image.BandFmt))
}

// Coding returns the image coding
func (ref *ImageRef) Coding() Coding {
	return Coding(int(ref.image.Coding))
}

// Interpretation returns the current interpretation
func (ref *ImageRef) Interpretation() Interpretation {
	return Interpretation(int(ref.image.Type))
}

// Composite overlays the given image over this one
func (ref *ImageRef) Composite(overlay *ImageRef, mode BlendMode) error {
	out, err := vipsComposite([]*C.VipsImage{ref.image, overlay.image}, mode)
	if err != nil {
		return err
	}
	ref.SetImage(out)
	return nil
}

// ToBytes writes the image to memory in VIPs format and returns the raw bytes, useful for storage.
func (ref *ImageRef) ToBytes() ([]byte, error) {
	var cSize C.size_t
	cData := C.vips_image_write_to_memory(ref.image, &cSize)
	if cData == nil {
		return nil, errors.New("Failed to write image to memory")
	}
	defer C.free(cData)

	bytes := C.GoBytes(unsafe.Pointer(cData), C.int(cSize))
	return bytes, nil
}
