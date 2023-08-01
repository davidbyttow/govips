package vips

// #cgo pkg-config: vips
// #include "stream.h"
import "C"
import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

// Streaming support. Provides a go wrapper to the libvips streaming facilities.
// https://www.libvips.org/API/current/VipsTargetCustom.html

const headerSizeBytes = 1024

// Source is the go analog to a VipsSource. It is used for streaming input sources.
// https://www.libvips.org/2019/11/29/True-streaming-for-libvips.html
type Source struct {
	lock    sync.Mutex
	header  []byte // first 1024 bytes of source
	file    *os.File
	vipsSrc *C.VipsSource
}

// Target is the go analog to a VipsTarget. It is used for streaming input sources.
// https://www.libvips.org/2019/11/29/True-streaming-for-libvips.html
type Target struct {
	lock       sync.Mutex
	vipsTarget *C.VipsTarget
}

// Close will release the libvips resources held by Source,
// after which it cannot be used further
func (r *Source) Close() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.vipsSrc != nil {
		clearSource(r.vipsSrc)
		r.vipsSrc = nil
	}
	if r.file != nil {
		r.file.Close()
	}
}

// Close will release the libvips resources held by Target,
// after which it cannot be used further
func (r *Target) Close() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.vipsTarget != nil {
		clearTarget(r.vipsTarget)
		r.vipsTarget = nil
	}
}

// NewTargetToPipe creates a Target from a file descriptor
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-target-new-to-descriptor
func NewTargetToPipe(file *os.File) (*Target, error) {
	if file == nil {
		return nil, ErrWriterInvalid
	}

	descriptor := file.Fd()
	vipsTarget := C.vips_target_new_to_descriptor(C.int(descriptor))

	return newTarget(vipsTarget)
}

// NewTargetToFile create a target attached to a file
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-target-new-to-file
func NewTargetToFile(path string) (*Target, error) {

	cpath := C.CString(path)
	defer freeCString(cpath)

	vipsTarget := C.vips_target_new_to_file(cpath)

	return newTarget(vipsTarget)

}

// newTarget initializes and returns a new Target
func newTarget(vipsTarget *C.VipsTarget) (*Target, error) {
	if vipsTarget == nil {
		return nil, handleVipsError()
	}

	target := &Target{
		vipsTarget: vipsTarget,
	}

	runtime.SetFinalizer(target, finalizeTarget)

	return target, nil
}

// NewSourceFromPipe creates a Source from an os.Pipe file descriptor
// If reader is nil a nil source and io.EOF is returned.
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-source-new-from-descriptor
func NewSourceFromPipe(file *os.File) (*Source, error) {

	if file == nil {
		return nil, io.EOF
	}

	descriptor := file.Fd()
	vipsSrc := C.vips_source_new_from_descriptor(C.int(descriptor))

	return newSource(file, vipsSrc)
}

// NewSourceFromFile creates a Source from a file path
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-source-new-from-file
func NewSourceFromFile(path string) (*Source, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cpath := C.CString(path)
	defer freeCString(cpath)

	vipsSrc := C.vips_source_new_from_file(cpath)

	return newSource(file, vipsSrc)
}

// newSource initializes and returns a new Source
func newSource(file *os.File, vipsSrc *C.VipsSource) (*Source, error) {

	if vipsSrc == nil {
		return nil, handleVipsError()
	}

	header := make([]byte, headerSizeBytes)
	io.ReadFull(file, header)
	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	source := &Source{
		header:  header,
		file:    file,
		vipsSrc: vipsSrc,
	}

	runtime.SetFinalizer(source, finalizeSource)
	return source, err
}

func finalizeSource(ref *Source) {
	govipsLog("govips", LogLevelDebug, fmt.Sprintf("closing source %p", ref))
	ref.Close()
}

func finalizeTarget(ref *Target) {
	govipsLog("govips", LogLevelDebug, fmt.Sprintf("closing target %p", ref))
	ref.Close()
}

func clearTarget(ref *C.VipsTarget) {
	C.clear_target(&ref)
}
func clearSource(ref *C.VipsSource) {
	C.clear_source(&ref)
}
