package govips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

var stringBuffer4096 = fixedString(4096)

func vipsForeignFindLoad(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_load(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindLoadBuffer(bytes []byte) (string, error) {
	cOperationName := C.vips_foreign_find_load_buffer(
		byteArrayPointer(bytes),
		C.size_t(len(bytes)))
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	debug("Found foreign load for buffer: %s", C.GoString(cOperationName))
	return C.GoString(cOperationName), nil
}

func vipsForeignFindSave(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_save(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindSaveBuffer(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_save_buffer(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	debug("Found foreign save for buffer: %s", C.GoString(cOperationName))
	return C.GoString(cOperationName), nil
}

func vipsInterpolateNew(name string) (*C.VipsInterpolate, error) {
	cName := C.CString(name)
	defer freeCString(cName)

	interp := C.vips_interpolate_new(cName)
	if interp == nil {
		return nil, ErrInvalidInterpolator
	}
	return interp, nil
}

func vipsOperationNew(name string) *C.VipsOperation {
	cName := C.CString(name)
	defer freeCString(cName)
	return C.vips_operation_new(cName)
}

func vipsFilenameSplit8(file string) (string, string) {
	cFile := C.CString(file)
	defer freeCString(cFile)

	cFilename := C.CString(stringBuffer4096)
	defer freeCString(cFilename)

	cOptionString := C.CString(stringBuffer4096)
	defer freeCString(c_optionString)

	C.vips__filename_split8(cFile, cFilename, cOptionString)

	fileName := C.GoString(cFilename)
	optionString := C.GoString(cOptionString)
	return fileName, optionString
}

func vipsColorspaceIsSupported(image *C.VipsImage) bool {
	return fromGboolean(C.vips_colourspace_issupported(image))
}
