package vips

// #include "resample.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

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
func vipsResizeWithVScale(in *C.VipsImage, hscale, vscale float64, kernel Kernel) (*C.VipsImage, error) {
	incOpCounter("resize")
	var out *C.VipsImage

	// libvips recommends Lanczos3 as the default kernel
	if kernel == KernelAuto {
		kernel = KernelLanczos3
	}

	if err := C.resize_image(in, &out, C.double(hscale), C.double(vscale), C.int(kernel)); err != 0 {
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

// https://www.libvips.org/API/current/libvips-resample.html#vips-thumbnail-buffer
func vipsThumbnailFromBuffer(buf []byte, width, height int, crop Interesting) (*C.VipsImage, ImageType, error) {
	src := buf
	// Reference src here so it's not garbage collected during image initialization.
	defer runtime.KeepAlive(src)

	var err error

	imageType := DetermineImageType(src)

	if imageType == ImageTypeBMP {
		src, err = bmpToPNG(src)
		if err != nil {
			return nil, ImageTypeUnknown, err
		}
		imageType = ImageTypePNG
	}

	if !IsTypeSupported(imageType) {
		govipsLog("govips", LogLevelInfo, fmt.Sprintf("failed to understand image format size=%d", len(src)))
		return nil, ImageTypeUnknown, ErrUnsupportedImageFormat
	}

	var out *C.VipsImage

	if err := C.thumbnail_buffer(unsafe.Pointer(&src[0]), C.size_t(len(src)), &out, C.int(width), C.int(height), C.int(crop)); err != 0 {
		return nil, ImageTypeUnknown, handleImageError(out)
	}

	return out, imageType, nil
}

// https://libvips.github.io/libvips/API/current/libvips-resample.html#vips-mapim
func vipsMapim(in *C.VipsImage, index *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("mapim")
	var out *C.VipsImage

	if err := C.mapim(in, &out, index); err != 0 {
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
