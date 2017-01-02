package govips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"runtime"
	"unsafe"
)

type Operation struct {
	name      string
	operation *C.VipsOperation
}

func (o *Operation) Name() string {
	return o.name
}

func CallOperation(name string, options *Options, operationString string) error {
	operation := newOperation(name)
	if operationString != "" {
		c_operationString := C.CString(operationString)
		defer freeCString(c_operationString)
		if C.vips_object_set_from_string(
			(*C.VipsObject)(unsafe.Pointer(operation)),
			c_operationString) != 0 {
			return handleVipsError()
		}
	}
	return call(operation, options)
}

func call(operation *Operation, options *Options) error {
	// TODO(d): Unref the outputs

	if options != nil {
		operation.applyOptions(options)
	}

	if err := operation.build(); err != nil {
		return err
	}

	if options != nil {
		operation.writeOutputs(options)
	}

	return nil
}

func newOperation(name string) *Operation {
	o := &Operation{
		operation: vipsOperationNew(name),
	}
	runtime.SetFinalizer(o, finalizeOperation)
	return o
}

func finalizeOperation(o *Operation) {
	C.g_object_unref(C.gpointer(o.operation))
}

func (o *Operation) applyOptions(options *Options) {
	for _, option := range options.options {
		if option.isOutput {
			continue
		}
		c_name := C.CString(option.name)
		defer freeCString(c_name)
		C.SetProperty(
			(*C.VipsObject)(unsafe.Pointer(o.operation)),
			c_name,
			&option.gvalue)
	}
}

func (o *Operation) build() error {
	if ret := C.vips_cache_operation_buildp(&o.operation); ret != 0 {
		return handleVipsError()
	}
	return nil
}

func (o *Operation) writeOutputs(options *Options) {
	for _, option := range options.options {
		if !option.isOutput {
			continue
		}
		c_name := C.CString(option.name)
		defer freeCString(c_name)
		C.g_object_get_property(
			(*C.GObject)(unsafe.Pointer(o.operation)),
			(*C.gchar)(c_name),
			&option.gvalue)
		option.Deserialize()
	}
}
