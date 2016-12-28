package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"runtime"
	"unsafe"
)

func Call(name string, options *Options) error {
	operation := newOperation(name)

	if operation == nil {
		return handleVipsError()
	}

	// TODO(d): Unref the outputs

	operation.applyOptions(options)

	if err := operation.build(); err != nil {
		return err
	}

	if options != nil {
		operation.writeOutputs(options)
	}

	return nil
}

type Operation struct {
	name      string
	operation *C.VipsOperation
}

func (o Operation) Name() string {
	return o.name
}

func newOperation(name string) *Operation {
	op := C.vips_operation_new(C.CString(name))
	operation := &Operation{
		operation: op,
	}
	runtime.SetFinalizer(operation, finalizeOperation)
	return operation
}

func finalizeOperation(o *Operation) {
	C.g_object_unref(C.gpointer(o))
}

func (o Operation) applyOptions(options *Options) {
	for _, option := range options.options {
		if option.isOutput {
			continue
		}
		C.set_property(
			(*C.VipsObject)(unsafe.Pointer(o.operation)),
			C.CString(option.name),
			&option.gvalue)
	}
}

func (o Operation) build() error {
	if ret := C.vips_cache_operation_buildp(&o.operation); ret != 0 {
		return handleVipsError()
	}
	return nil
}

func (o Operation) writeOutputs(options *Options) {
	for _, option := range options.options {
		if !option.isOutput {
			continue
		}
		C.g_object_get_property(
			(*C.GObject)(unsafe.Pointer(o.operation)),
			(*C.gchar)(C.CString(option.name)),
			&option.gvalue)
		option.Deserialize()
	}
}
