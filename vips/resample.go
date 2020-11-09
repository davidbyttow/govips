package vips

// #cgo pkg-config: vips
// #include "resample.h"
import "C"

// Kernel represents VipsKernel type
type Kernel int

// Kernel enum
const (
	KernelAuto     Kernel = -1
	KernelNearest  Kernel = C.VIPS_KERNEL_NEAREST
	KernelLinear   Kernel = C.VIPS_KERNEL_LINEAR
	KernelCubic    Kernel = C.VIPS_KERNEL_CUBIC
	KernelLanczos2 Kernel = C.VIPS_KERNEL_LANCZOS2
	KernelLanczos3 Kernel = C.VIPS_KERNEL_LANCZOS3
	KernelMitchell Kernel = C.VIPS_KERNEL_MITCHELL
)

// https://libvips.github.io/libvips/API/current/libvips-resample.html#vips-resize
func vipsResize(in *C.VipsImage, scale float64, kernel Kernel) (*C.VipsImage, error) {
	incOpCounter("resize")
	var out *C.VipsImage

	// libvips recommends Lanczos3 as the default kernel
	if kernel == KernelAuto {
		kernel = KernelLanczos3
	}

	if err := C.resize_image(in, &out, C.double(scale), C.double(-1), C.int(kernel)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-resample.html#vips-resize
func vipsResizeWithVScale(in *C.VipsImage, scale, vscale float64, kernel Kernel) (*C.VipsImage, error) {
	incOpCounter("resize")
	var out *C.VipsImage

	if err := C.resize_image(in, &out, C.double(scale), C.gdouble(vscale), C.int(kernel)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}
