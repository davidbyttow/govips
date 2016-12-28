package gimage

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import "runtime"

type Blob struct {
	blob *C.VipsBlob
}

func newBlob(blob *C.VipsBlob) *Blob {
	b := &Blob{
		blob: blob,
	}
	runtime.SetFinalizer(b, finalizer)
	return b
}

func NewBlob(buf []byte) *Blob {
	blob := C.vips_blob_new(
		nil,
		cPtr(buf),
		C.size_t(len(buf)))
	return newBlob(blob)
}

func finalizer(b *Blob) {
	// TODO(d): Finalize this properly
}
