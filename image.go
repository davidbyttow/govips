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

type OperationStream struct {
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

// OpenFromBuffer loads an image buffer and creates a new Image
func OpenFromBuffer(bytes []byte, opts ...OptionFunc) (*OperationStream, error) {
	startupIfNeeded()

	image, _, err := vipsLoadFromBuffer(bytes)
	if err != nil {
		return nil, err
	}

	return newStream(image), nil
}

func newStream(vipsImage *C.VipsImage) *OperationStream {
	stream := &OperationStream{
		image: vipsImage,
	}
	runtime.SetFinalizer(stream, finalizeStream)
	return stream
}

func finalizeStream(os *OperationStream) {
	os.Close()
}

func (os *OperationStream) SetImage(image *C.VipsImage) {
	if os.image != nil {
		C.g_object_unref(C.gpointer(os.image))
	}
	os.image = image
}

func (os *OperationStream) Close() {
	os.SetImage(nil)
}

// Width returns the width of this image
func (os *OperationStream) Width() int {
	return int(os.image.Xsize)
}

// Height returns the height of this iamge
func (os *OperationStream) Height() int {
	return int(os.image.Ysize)
}

// Bands returns the number of bands for this image
func (os *OperationStream) Bands() int {
	return int(os.image.Bands)
}

// ResX returns the X resolution
func (os *OperationStream) ResX() float64 {
	return float64(os.image.Xres)
}

// ResY returns the Y resolution
func (os *OperationStream) ResY() float64 {
	return float64(os.image.Yres)
}

// OffsetX returns the X offset
func (os *OperationStream) OffsetX() int {
	return int(os.image.Xoffset)
}

// OffsetY returns the Y offset
func (os *OperationStream) OffsetY() int {
	return int(os.image.Yoffset)
}

// BandFormat returns the current band format
func (os *OperationStream) BandFormat() BandFormat {
	return BandFormat(int(os.image.BandFmt))
}

// Coding returns the image coding
func (os *OperationStream) Coding() Coding {
	return Coding(int(os.image.Coding))
}

// Interpretation returns the current interpretation
func (os *OperationStream) Interpretation() Interpretation {
	return Interpretation(int(os.image.Type))
}

// ToBytes writes the image to memory in VIPs format and returns the raw bytes, useful for storage.
func (os *OperationStream) ToBytes() ([]byte, error) {
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
func (os *OperationStream) Export(options ExportOptions) ([]byte, error) {
	return vipsExportBuffer(os.image, &options)
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

func (os *OperationStream) CopyEvents(events []*CallEvent) {
	if len(events) > 0 {
		os.callEvents = append(os.callEvents, events...)
	}
}

func (os *OperationStream) LogCallEvent(name string, options *Options) {
	os.callEvents = append(os.callEvents, &CallEvent{name, options})
}

func (os *OperationStream) CallEventLog() []*CallEvent {
	return os.callEvents
}
