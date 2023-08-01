package vips

// #cgo pkg-config: vips
// #include "stream.h"
import "C"
import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

// Streaming support. Provides a go wrapper to the libvips streaming facilities.
// https://www.libvips.org/API/current/VipsTargetCustom.html

const headerSizeBytes = 1024

// Source is the go analog to a VipsSource. It is used for streaming input sources.
// https://www.libvips.org/2019/11/29/True-streaming-for-libvips.html
type Source struct {
	lock    sync.Mutex
	header  []byte // first headerSizeBytes of source reader
	reader  io.Reader
	vipsSrc *C.VipsSource
}

// Target is the go analog to a VipsTarget. It is used for streaming input sources.
// https://www.libvips.org/2019/11/29/True-streaming-for-libvips.html
type Target struct {
	lock       sync.Mutex
	writer     io.Writer
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
	if file, ok := r.reader.(*os.File); ok {
		file.Close()
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
	if file, ok := r.writer.(*os.File); ok {
		file.Close()
	}
}

// NewTargetToWriter creates a Target using the passed in writer
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-target-custom-new
func NewTargetToWriter(writer io.Writer) (*Target, error) {

	if writer == nil {
		return nil, ErrWriterInvalid
	}

	target := &Target{}

	targetPtr := unsafe.Pointer(target)

	vipsTarget := C.create_go_custom_target(targetPtr)

	if vipsTarget == nil {
		return nil, handleVipsError()
	}

	target.writer = writer
	target.vipsTarget = vipsTarget

	runtime.SetFinalizer(target, finalizeTarget)

	return target, nil

}

// NewTargetToWriter create a target attached to a file
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-target-new-to-file
func NewTargetToFile(path string) (*Target, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	target := &Target{}

	cpath := C.CString(path)
	defer freeCString(cpath)

	vipsTarget := C.vips_target_new_to_file(cpath)

	if vipsTarget == nil {
		return nil, handleVipsError()
	}

	target.vipsTarget = vipsTarget
	target.writer = file

	runtime.SetFinalizer(target, finalizeTarget)

	return target, nil

}

// NewSourceFromReader creates a Source using the passed in reader.
// If reader is nil a nil source and io.EOF is returned.
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-source-custom-new
func NewSourceFromReader(reader io.Reader) (*Source, error) {

	if reader == nil {
		return nil, io.EOF
	}

	source := &Source{}

	sourcePtr := unsafe.Pointer(source)

	vipsSrc := C.create_go_custom_source(sourcePtr)

	if vipsSrc == nil {
		return nil, handleVipsError()
	}

	bufReader := bufio.NewReader(reader)
	header, _ := bufReader.Peek(headerSizeBytes)

	source.header = header
	source.reader = bufReader
	source.vipsSrc = vipsSrc

	runtime.SetFinalizer(source, finalizeSource)

	return source, nil
}

// NewSourceFromFile creates a Source from a file path
// https://www.libvips.org/API/current/VipsTargetCustom.html#vips-source-new-from-file
func NewSourceFromFile(path string) (*Source, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	source := &Source{}

	cpath := C.CString(path)
	defer freeCString(cpath)

	vipsSrc := C.vips_source_new_from_file(cpath)

	if vipsSrc == nil {
		return nil, handleVipsError()
	}

	header := make([]byte, headerSizeBytes)
	io.ReadFull(file, header)
	file.Seek(0, 0)

	source.header = header
	source.reader = file
	source.vipsSrc = vipsSrc

	runtime.SetFinalizer(source, finalizeSource)

	return source, nil
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
		return C.longlong(-1)
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
		return C.longlong(-1)
	}

	n, err := skr.Seek(int64(offset), whence)

	if err != nil {
		govipsLog("govips", LogLevelError, fmt.Sprintf("goSourceSeek: Error: %v [offset %d | whence %d]", err, n, whence))
		return C.longlong(-1)
	}
	govipsLog("govips", LogLevelDebug, fmt.Sprintf("goSourceSeek: OK [seek %d | whence %d]", n, whence))
	return C.longlong(n)

}

//export goTargetRead
func goTargetRead(targetPtr unsafe.Pointer, buffer unsafe.Pointer, length C.longlong) (read C.longlong) {

	target := (*Target)(unsafe.Pointer(targetPtr))

	// https://stackoverflow.com/questions/51187973/how-to-create-an-array-or-a-slice-from-an-array-unsafe-pointer-in-golang
	sh := &reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(length),
		Cap:  int(length),
	}

	buf := *(*[]byte)(unsafe.Pointer(sh))

	reader, ok := target.writer.(io.Reader)

	if !ok {
		return C.longlong(-1)
	}
	n, err := reader.Read(buf)

	switch {
	case errors.Is(err, io.EOF):
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goTargetRead: EOF [read %d]", n))
		return C.longlong(n)
	case err != nil:
		govipsLog("govips", LogLevelError, fmt.Sprintf("goTargetRead: Error: %v [read %d]", err, n))
		return C.longlong(-1)
	default:
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goTargetRead: OK [read %d]", n))
		return C.longlong(n)
	}
}

//export goTargetSeek
func goTargetSeek(targetPtr unsafe.Pointer, offset C.longlong, whence int) (newOffset C.longlong) {

	target := (*Target)(unsafe.Pointer(targetPtr))

	skr, ok := target.writer.(io.ReadSeeker)
	if !ok {
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goTargetRead: Seek not supported"))
		return -1 // Unsupported!
	}

	switch whence {
	case io.SeekStart, io.SeekCurrent, io.SeekEnd:
	default:
		govipsLog("govips", LogLevelError, fmt.Sprintf("goTargetSeek: Invalid whence value [%d]", whence))
		return C.longlong(-1)
	}

	n, err := skr.Seek(int64(offset), whence)

	if err != nil {
		govipsLog("govips", LogLevelError, fmt.Sprintf("goTargetSeek: Error: %v [offset %d | whence %d]", err, n, whence))
		return C.longlong(-1)
	}
	govipsLog("govips", LogLevelDebug, fmt.Sprintf("goTargetSeek: OK [seek %d | whence %d]", n, whence))
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
		govipsLog("govips", LogLevelError, fmt.Sprintf("goTargetWrite: Error: %v [wrote %d]", err, n))
		return C.longlong(-1)
	default:
		govipsLog("govips", LogLevelDebug, fmt.Sprintf("goTargetWrite: OK [wrote %d]", n))
		return C.longlong(n)
	}
}

//export goTargetEnd
func goTargetEnd(targetPtr unsafe.Pointer) (write C.longlong) {
	// target := (*Target)(unsafe.Pointer(targetPtr))

	return C.longlong(0)

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
