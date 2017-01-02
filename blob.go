package govips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"runtime"
	"unsafe"
)

type Blob struct {
	c_blob *C.VipsBlob
}

func newBlob(c_blob *C.VipsBlob) *Blob {
	b := &Blob{
		c_blob: c_blob,
	}
	runtime.SetFinalizer(b, finalizer)
	return b
}

func NewBlob(buf []byte) *Blob {
	c_blob := C.vips_blob_new(
		nil,
		byteArrayPointer(buf),
		C.size_t(len(buf)))
	return newBlob(c_blob)
}

func (t *Blob) ToBytes() []byte {
	c_area := t.CArea()
	return C.GoBytes(unsafe.Pointer(c_area), C.int(c_area.length))
}

func (t *Blob) Length() int {
	return int(t.CArea().length)
}

func (t *Blob) CArea() *C.VipsArea {
	return (*C.VipsArea)(unsafe.Pointer(t.c_blob))
}

func finalizer(b *Blob) {
	C.vips_area_unref(b.CArea())
}
