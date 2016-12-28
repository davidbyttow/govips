package gimage

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

func (in Image) Shrink(shrinkH, shrinkY float64) *Image {
	out := in.SetImage(nil)
	Call("shrink", NewOptions().
		SetImage("in", in).
		SetImageOut("out", out).
		SetDouble("hshrink", shrinkH).
		SetDouble("vshrink", shrinkY))
	return out
}
