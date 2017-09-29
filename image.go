package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

// Image is an immutable structure that represents an image in memory
type Image struct {
	image      *C.VipsImage
	callEvents []*CallEvent
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

// NewImageFromMemory wraps an image around a memory area. The memory area must be a simple
// array (e.g., RGBRGBRGB), left-to-right, top-to-bottom.
func NewImageFromMemory(bytes []byte, width, height, bands int, format BandFormat) (*Image, error) {
	startupIfNeeded()

	vipsImage := C.vips_image_new_from_memory_copy(
		byteArrayPointer(bytes),
		C.size_t(len(bytes)),
		C.int(width),
		C.int(height),
		C.int(bands),
		C.VipsBandFormat(format))

	return newImage(vipsImage), nil
}

// NewImageFromBuffer loads an image buffer and creates a new Image
func NewImageFromBuffer(bytes []byte, opts ...OptionFunc) (*Image, error) {
	startupIfNeeded()

	image, _, err := vipsLoadFromBuffer(bytes)
	if err != nil {
		return nil, err
	}

	return newImage(image), nil
}

func NewThumbnailFromBuffer(bytes []byte, width int, opts ...OptionFunc) (*Image, error) {
	startupIfNeeded()
	blob := NewBlob(bytes)
	out := ThumbnailBuffer(blob, width, opts...)
	return out, nil
}

func newImage(vipsImage *C.VipsImage) *Image {
	image := &Image{
		image: vipsImage,
	}
	runtime.SetFinalizer(image, finalizeImage)
	return image
}

func finalizeImage(i *Image) {
	i.Close()
}

func (i *Image) Close() {
	if i.image != nil {
		C.g_object_unref(C.gpointer(i.image))
		i.image = nil
	}
}

// Width returns the width of this image
func (i *Image) Width() int {
	return int(i.image.Xsize)
}

// Height returns the height of this iamge
func (i *Image) Height() int {
	return int(i.image.Ysize)
}

// Bands returns the number of bands for this image
func (i *Image) Bands() int {
	return int(i.image.Bands)
}

// ResX returns the X resolution
func (i *Image) ResX() float64 {
	return float64(i.image.Xres)
}

// ResY returns the Y resolution
func (i *Image) ResY() float64 {
	return float64(i.image.Yres)
}

// OffsetX returns the X offset
func (i *Image) OffsetX() int {
	return int(i.image.Xoffset)
}

// OffsetY returns the Y offset
func (i *Image) OffsetY() int {
	return int(i.image.Yoffset)
}

// BandFormat returns the current band format
func (i *Image) BandFormat() BandFormat {
	return BandFormat(int(i.image.BandFmt))
}

// Coding returns the image coding
func (i *Image) Coding() Coding {
	return Coding(int(i.image.Coding))
}

// Interpretation returns the current interpretation
func (i *Image) Interpretation() Interpretation {
	return Interpretation(int(i.image.Type))
}

// ToBytes writes the image to memory in VIPs format and returns the raw bytes, useful for storage.
func (i *Image) ToBytes() ([]byte, error) {
	var cSize C.size_t
	cData := C.vips_image_write_to_memory(i.image, &cSize)
	if cData == nil {
		return nil, errors.New("Failed to write image to memory")
	}
	defer C.free(cData)

	bytes := C.GoBytes(unsafe.Pointer(cData), C.int(cSize))
	return bytes, nil
}

// WriteToBuffer writes the image to a buffer in a format represented by the given suffix (e.g., .jpeg)
func (i *Image) Export(options ExportOptions) ([]byte, error) {
	return vipsExportBuffer(i.image, &options)
}

type CallEvent struct {
	Name    string
	Options *Options
}

func (c CallEvent) String() string {
	var args []string
	for _, o := range c.Options.Options {
		args = append(args, o.String())
	}
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(args, ", "))
}

func (i *Image) CopyEvents(events []*CallEvent) {
	if len(events) > 0 {
		i.callEvents = append(i.callEvents, events...)
	}
}

func (i *Image) LogCallEvent(name string, options *Options) {
	i.callEvents = append(i.callEvents, &CallEvent{name, options})
}

func (i *Image) CallEventLog() []*CallEvent {
	return i.callEvents
}
