package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

func (i Image) Shrink(hShrink, vShrink float64) (*Image, error) {
	defer C.g_object_unref(C.gpointer(i.image))

	var out *C.VipsImage

	ret := C.shrink(
		i.image,
		&out,
		C.double(float64(hShrink)),
		C.double(float64(vShrink)))

	if ret != 0 {
		return nil, handleVipsError()
	}

	return i.SetImage(out), nil
}
