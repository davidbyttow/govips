package govips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// Blob wraps the internal libvips VipsBlob type
type Blob struct {
	cBlob *C.VipsBlob
}

func newBlob(cBlob *C.VipsBlob) *Blob {
	b := &Blob{
		cBlob: cBlob,
	}
	runtime.SetFinalizer(b, finalizer)
	return b
}

// NewBlob constructs a new blob from the given buffer, does not take ownership of the bytes
func NewBlob(buf []byte) *Blob {
	cBlob := C.vips_blob_new(
		nil,
		byteArrayPointer(buf),
		C.size_t(len(buf)))
	return newBlob(cBlob)
}

// ToBytes creates a copy of the blob's byte array
func (t *Blob) ToBytes() []byte {
	area := t.cArea()
	return C.GoBytes(unsafe.Pointer(area.data), C.int(area.length))
}

// Length returns the length of the byte array
func (t *Blob) Length() int {
	return int(t.cArea().length)
}

// CArea returns the internal representation of the VipsArea for this blob
func (t *Blob) cArea() *C.VipsArea {
	return (*C.VipsArea)(unsafe.Pointer(t.cBlob))
}

func finalizer(b *Blob) {
	C.vips_area_unref(b.cArea())
}
