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

const headerSize = 20

type Source struct {
	lock    sync.Mutex
	header  []byte // 20 byte header of source
	reader  io.Reader
	vipsSrc *C.VipsSource
}

type Target struct {
	lock       sync.Mutex
	writer     io.Writer
	vipsTarget *C.VipsTarget
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
		return nil, handleVipsError()
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
		return nil, handleVipsError()
	}

	bufReader := bufio.NewReader(reader)
	header, err := bufReader.Peek(headerSize)
	if err != nil {
		return nil, err
	}

	source.header = header
	source.reader = bufReader
	source.vipsSrc = vipsSrc

	runtime.SetFinalizer(source, finalizeSource)

	return source, nil
}

// NewSourceFromFile creates a Source from a file path
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

	header := make([]byte, headerSize)
	io.ReadFull(file, header)
	file.Seek(0, 0)

	source.header = header
	source.reader = file
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

func clearTarget(ref *C.VipsTarget) {
	C.clear_target(&ref)
}
func clearSource(ref *C.VipsSource) {
	C.clear_source(&ref)
}
