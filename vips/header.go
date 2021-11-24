package vips

// #include "header.h"
import "C"
import "unsafe"

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

func vipsImageGetPages(in *C.VipsImage) int {
	return int(C.get_image_get_n_pages(in))
}
