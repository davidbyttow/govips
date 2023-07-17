package vips

// #cgo pkg-config: vips
// #include "stream.h"
import "C"
import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

type Source struct {
	lock    sync.Mutex
	reader  io.Reader
	vipsSrc *C.VipsSourceCustom
}

type Target struct {
	lock       sync.Mutex
	writer     io.Writer
	vipsTarget *C.VipsTargetCustom
}

func (r *Source) Close() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.vipsSrc != nil {
		clearSource(r.vipsSrc)
		r.vipsSrc = nil
	}
}

func (r *Target) Close() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.vipsTarget != nil {
		clearTarget(r.vipsTarget)
		r.vipsTarget = nil
	}
}

//export goSourceRead
func goSourceRead(sourcePtr unsafe.Pointer, buffer unsafe.Pointer, length C.longlong) (read C.longlong) {

	source := (*Source)(unsafe.Pointer(sourcePtr))

	// https://stackoverflow.com/questions/51187973/how-to-create-an-array-or-a-slice-from-an-array-unsafe-pointer-in-golang
	sh := &reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(length),
		Cap:  int(length),
	}

	buf := *(*[]byte)(unsafe.Pointer(sh))

	n, err := source.reader.Read(buf)

	switch {
	case errors.Is(err, io.EOF):
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goSourceRead: EOF [read %d]", n))
		return C.longlong(n)
	case err != nil:
		govipsLog("govips", LogLevelError, fmt.Sprintf("goSourceRead: Error: %v [read %d]", err, n))
		return -1
	default:
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goSourceRead: OK [read %d]", n))
		return C.longlong(n)
	}
}

//export goSourceSeek
func goSourceSeek(sourcePtr unsafe.Pointer, offset C.longlong, whence int) (newOffset C.longlong) {

	source := (*Source)(unsafe.Pointer(sourcePtr))

	skr, ok := source.reader.(io.ReadSeeker)
	if !ok {
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goSourceRead: Seek not supported"))
		return -1 // Unsupported!
	}

	switch whence {
	case io.SeekStart, io.SeekCurrent, io.SeekEnd:
	default:
		govipsLog("govips", LogLevelError, fmt.Sprintf("goSourceSeek: Invalid whence value [%d]", whence))
		return -1
	}

	n, err := skr.Seek(int64(offset), whence)

	if err != nil {
		govipsLog("govips", LogLevelError, fmt.Sprintf("goSourceSeek: Error: %v [offset %d | whence %d]", err, n, whence))
		return -1
	}
	govipsLog("govips", LogLevelDebug, fmt.Sprintf("goSourceSeek: OK [seek %d | whence %d]", n, whence))
	return C.longlong(n)

}

//export goTargetWrite
func goTargetWrite(targetPtr unsafe.Pointer, buffer unsafe.Pointer, length C.longlong) (write C.longlong) {
	target := (*Target)(unsafe.Pointer(targetPtr))

	sh := &reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(length),
		Cap:  int(length),
	}

	buf := *(*[]byte)(unsafe.Pointer(sh))

	n, err := target.writer.Write(buf)

	switch {
	case err != nil:
		govipsLog("govips", LogLevelError, fmt.Sprintf("goTargetWrite: Error: %v [read %d]", err, n))
		return -1
	default:
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goTargetWrite: OK [read %d]", n))
		return C.longlong(n)
	}
}

// NewTargetToWriter creates a Target using the passed in writer
func NewTargetToWriter(writer io.Writer) (*Target, error) {

	target := &Target{}

	targetPtr := unsafe.Pointer(target)

	vipsTarget := C.create_go_custom_target(targetPtr)

	if vipsTarget == nil {
		return nil, fmt.Errorf("error creating target from source: nil source")
	}

	target.writer = writer
	target.vipsTarget = vipsTarget

	runtime.SetFinalizer(target, finalizeTarget)

	return target, nil

}

// NewSourceFromReader creates a Source using the passed in reader
func NewSourceFromReader(reader io.Reader) (*Source, error) {

	source := &Source{}

	sourcePtr := unsafe.Pointer(source)

	vipsSrc := C.create_go_custom_source(sourcePtr)

	if vipsSrc == nil {
		return nil, fmt.Errorf("error creating source from reader: nil source")
	}

	source.reader = reader
	source.vipsSrc = vipsSrc

	runtime.SetFinalizer(source, finalizeSource)

	return source, nil
}

func finalizeSource(ref *Source) {
	govipsLog("govips", LogLevelDebug, fmt.Sprintf("closing source %p", ref))
	ref.Close()
}

func finalizeTarget(ref *Target) {
	govipsLog("govips", LogLevelDebug, fmt.Sprintf("closing target %p", ref))
	ref.Close()
}

func clearTarget(ref *C.VipsTargetCustom) {
	C.clear_target(&ref)
}
func clearSource(ref *C.VipsSourceCustom) {
	C.clear_source(&ref)
}
