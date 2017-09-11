package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import "unsafe"

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
	defer freeCString(cOptionString)

	C.vips__filename_split8(cFile, cFilename, cOptionString)

	fileName := C.GoString(cFilename)
	optionString := C.GoString(cOptionString)
	return fileName, optionString
}

func vipsCallString(name string, options *Options, optionString string) error {
	operation := vipsOperationNew(name)
	//defer C.g_object_unref(C.gpointer(operation))

	if optionString != "" {
		cOptionString := C.CString(optionString)
		defer freeCString(cOptionString)

		if C.vips_object_set_from_string(
			(*C.VipsObject)(unsafe.Pointer(operation)),
			cOptionString) != 0 {
			return handleVipsError()
		}
	}
	return vipsCallOperation(operation, options)
}

func vipsCall(name string, options *Options) error {
	operation := vipsOperationNew(name)
	//defer C.g_object_unref(C.gpointer(operation))

	return vipsCallOperation(operation, options)
}

func vipsCallOperation(operation *C.VipsOperation, options *Options) error {
	// TODO(d): Unref the outputs
	if options != nil {
		for _, option := range options.options {
			if option.isOutput {
				continue
			}

			cName := C.CString(option.name)
			defer freeCString(cName)

			C.SetProperty(
				(*C.VipsObject)(unsafe.Pointer(operation)),
				cName,
				&option.gvalue)
		}
	}

	if ret := C.vips_cache_operation_buildp(&operation); ret != 0 {
		C.g_object_unref(C.gpointer(unsafe.Pointer(operation)))
		return handleVipsError()
	}

	// We defer this here because the pointer may have changed.
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(operation)))

	if options != nil {
		for _, option := range options.options {
			if !option.isOutput {
				continue
			}

			cName := C.CString(option.name)
			defer freeCString(cName)

			C.g_object_get_property(
				(*C.GObject)(unsafe.Pointer(operation)),
				(*C.gchar)(cName),
				&option.gvalue)
			option.Deserialize()
		}
	}

	return nil
}
