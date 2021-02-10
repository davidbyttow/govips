package vips

// #include "resample.h"
import "C"
import "unsafe"

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

func vipsThumbnail(in *C.VipsImage, width, height int, crop Interesting) (*C.VipsImage, error) {
	incOpCounter("thumbnail")
	var out *C.VipsImage

	if err := C.thumbnail_image(in, &out, C.int(width), C.int(height), C.int(crop)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-resample.html#vips-mapim
func vipsMapim(in *C.VipsImage, index *C.VipsImage, interpolate *C.VipsInterpolate) (*C.VipsImage, error) {
	incOpCounter("mapim")
	var out *C.VipsImage

	if interpolate == nil {
		govipsLog("govips", LogLevelWarning, "could not find interpolator, defaulting to bilinear")
		interpolate = C.interpolate_bilinear_static()
	}

	if err := C.mapim(in, &out, index, interpolate); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-histogram.html#vips-maplut
func vipsMaplut(in *C.VipsImage, lut *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("maplut")
	var out *C.VipsImage

	if err := C.maplut(in, &out, lut); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

type Interpolator string

const (
	Nearest  = "nearest"  // nearest-neighbour interpolation
	Bilinear = "bilinear" // bilinear interpolation
	Bicubic  = "bicubic"  // bicubic interpolation (Catmull-Rom)
	Lbb      = "lbb"      // reduced halo bicubic
	Nohalo   = "nohalo"   // edge sharpening resampler with halo reduction
	Vsqbs    = "vsqbs"    // B-splines with anti-aliasing
)

func vipsInterpolateNew(interpolator Interpolator) *C.VipsInterpolate {
	nickname := C.CString(string(interpolator))
	C.free(unsafe.Pointer(nickname))
	return C.interpolate_new(nickname)
}
