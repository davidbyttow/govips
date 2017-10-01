package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

// // Interpolator represents an interpolator when resizing images
// type Interpolator struct {
// 	interp *C.VipsInterpolate
// }
//
// // NewInterpolator creates a new Interpolator from the given name
// func NewInterpolator(name string) (*Interpolator, error) {
// 	interp, err := vipsInterpolateNew(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	out := &Interpolator{
// 		interp: interp,
// 	}
// 	runtime.SetFinalizer(out, finalizeInterpolator)
// 	return out, nil
// }
//
// func (i *Interpolator) Close() {
// 	if i.interp != nil {
// 		C.g_object_unref(C.gpointer(i.interp))
// 	}
// 	i.interp = nil
// }
//
// func finalizeInterpolator(i *Interpolator) {
// 	i.Close()
// }
