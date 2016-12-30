package gimage

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

func VipsForeignFindLoad(filename string) (string, error) {
	c_filename := C.CString(filename)
	defer freeCString(c_filename)

	c_operationName := C.vips_foreign_find_load(c_filename)
	if c_operationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(c_operationName), nil
}

func VipsForeignFindLoadBuffer(bytes []byte) (string, error) {
	c_operationName := C.vips_foreign_find_load_buffer(
		cPtr(bytes),
		C.size_t(len(bytes)))
	if c_operationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(c_operationName), nil
}

func VipsInterpolateNew(name string) (*C.VipsInterpolate, error) {
	c_name := C.CString(name)
	defer freeCString(c_name)

	interp := C.vips_interpolate_new(c_name)
	if interp == nil {
		return nil, ErrInvalidInterpolator
	}
	return interp, nil
}

func VipsOperationNew(name string) *C.VipsOperation {
	c_name := C.CString(name)
	defer freeCString(c_name)
	return C.vips_operation_new(c_name)
}
