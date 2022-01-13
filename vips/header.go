package vips

// #include "header.h"
import "C"
import (
	"strings"
	"unsafe"
)

func vipsHasICCProfile(in *C.VipsImage) bool {
	return int(C.has_icc_profile(in)) != 0
}

func vipsRemoveICCProfile(in *C.VipsImage) bool {
	return fromGboolean(C.remove_icc_profile(in))
}

func vipsHasIPTC(in *C.VipsImage) bool {
	return int(C.has_iptc(in)) != 0
}

func vipsImageGetFields(in *C.VipsImage) (fields []string) {
	const maxFields = 256

	rawFields := C.image_get_fields(in)
	defer C.g_strfreev(rawFields)

	cFields := (*[maxFields]*C.char)(unsafe.Pointer(rawFields))[:maxFields:maxFields]

	for _, field := range cFields {
		if field == nil {
			break
		}
		fields = append(fields, C.GoString(field))
	}
	return
}

func vipsRemoveMetadata(in *C.VipsImage) {
	C.remove_metadata(in)
}

func vipsGetMetaOrientation(in *C.VipsImage) int {
	return int(C.get_meta_orientation(in))
}

func vipsRemoveMetaOrientation(in *C.VipsImage) {
	C.remove_meta_orientation(in)
}

func vipsSetMetaOrientation(in *C.VipsImage, orientation int) {
	C.set_meta_orientation(in, C.int(orientation))
}

func vipsGetImageGetNPages(in *C.VipsImage) int {
	return int(C.get_image_get_n_pages(in))
}

func vipsGetPageHeight(in *C.VipsImage) int {
	return int(C.get_page_height(in))
}

func vipsSetPageHeight(in *C.VipsImage, height int) {
	C.set_page_height(in, C.int(height))
}

func vipsImageGetMetaLoader(in *C.VipsImage) (string, bool) {
	var out *C.char
	defer gFreePointer(unsafe.Pointer(out))
	code := int(C.get_meta_loader(in, &out))
	return C.GoString(out), code == 0
}

// vipsDetermineImageTypeFromMetaLoader determine the image type from vips-loader metadata
func vipsDetermineImageTypeFromMetaLoader(in *C.VipsImage) ImageType {
	vipsLoader, ok := vipsImageGetMetaLoader(in)
	if vipsLoader == "" || !ok {
		return ImageTypeUnknown
	}
	if strings.HasPrefix(vipsLoader, "jpeg") {
		return ImageTypeJPEG
	}
	if strings.HasPrefix(vipsLoader, "png") {
		return ImageTypePNG
	}
	if strings.HasPrefix(vipsLoader, "gif") {
		return ImageTypeGIF
	}
	if strings.HasPrefix(vipsLoader, "svg") {
		return ImageTypeSVG
	}
	if strings.HasPrefix(vipsLoader, "webp") {
		return ImageTypeWEBP
	}
	if strings.HasPrefix(vipsLoader, "jp2k") {
		return ImageTypeJP2K
	}
	if strings.HasPrefix(vipsLoader, "magick") {
		return ImageTypeMagick
	}
	if strings.HasPrefix(vipsLoader, "tiff") {
		return ImageTypeTIFF
	}
	if strings.HasPrefix(vipsLoader, "heif") {
		return ImageTypeHEIF
	}
	if strings.HasPrefix(vipsLoader, "pdf") {
		return ImageTypePDF
	}
	return ImageTypeUnknown
}
