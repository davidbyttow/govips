package main

// #cgo pkg-config: vips
// #include "introspect.h"
import "C"

import (
	"fmt"
	"sort"
	"unsafe"
)

// Introspect uses libvips GObject introspection to discover all operations
// and their argument signatures.
func Introspect() ([]OpDef, error) {
	var result C.IntrospectResult

	if ret := C.vipsgen_introspect(&result); ret != 0 {
		return nil, fmt.Errorf("vips introspection failed")
	}

	ops := make([]OpDef, 0, int(result.n_ops))
	for i := 0; i < int(result.n_ops); i++ {
		cop := &result.ops[i]
		op := OpDef{
			Name:        C.GoString(&cop.name[0]),
			Description: C.GoString(&cop.description[0]),
			Category:    C.GoString(&cop.category[0]),
		}

		for j := 0; j < int(cop.n_args); j++ {
			carg := &cop.args[j]
			arg := ArgDef{
				Name:     C.GoString(&carg.name[0]),
				Type:     ArgType(carg._type),
				Flags:    ArgFlags(carg.flags),
				Priority: int(carg.priority),
				Default:  float64(carg.defval),
				Min:      float64(carg.min),
				Max:      float64(carg.max),
				EnumType: C.GoString(&carg.enum_type[0]),
			}
			op.Args = append(op.Args, arg)
		}

		// Sort args by priority for consistent ordering.
		sort.Slice(op.Args, func(a, b int) bool {
			return op.Args[a].Priority < op.Args[b].Priority
		})

		ops = append(ops, op)
	}

	// Sort operations by name for deterministic output.
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].Name < ops[j].Name
	})

	// Deduplicate operations with the same name (can happen when multiple
	// GTypes share the same vips nickname).
	deduped := make([]OpDef, 0, len(ops))
	seen := make(map[string]bool)
	for _, op := range ops {
		if !seen[op.Name] {
			seen[op.Name] = true
			deduped = append(deduped, op)
		}
	}

	return deduped, nil
}

// IntrospectEnum discovers enum values for a given GType name.
func IntrospectEnum(typeName string) (*EnumDef, error) {
	cName := C.CString(typeName)
	defer C.free(unsafe.Pointer(cName))

	var result C.EnumInfo
	if ret := C.vipsgen_introspect_enum(cName, &result); ret != 0 {
		return nil, fmt.Errorf("failed to introspect enum %s", typeName)
	}

	def := &EnumDef{
		CName: C.GoString(&result.c_name[0]),
	}

	for i := 0; i < int(result.n_values); i++ {
		ev := &result.values[i]
		def.Values = append(def.Values, EnumValue{
			CName: C.GoString(&ev.c_name[0]),
			Nick:  C.GoString(&ev.nick[0]),
			Value: int(ev.value),
		})
	}

	return def, nil
}

// IntrospectEnums discovers enum values for all given GType names.
func IntrospectEnums(typeNames []string) ([]EnumDef, error) {
	if len(typeNames) == 0 {
		return nil, nil
	}

	cNames := make([]*C.char, len(typeNames))
	for i, name := range typeNames {
		cNames[i] = C.CString(name)
		defer C.free(unsafe.Pointer(cNames[i]))
	}

	results := make([]C.EnumInfo, len(typeNames))
	if ret := C.vipsgen_introspect_enums(
		(**C.char)(unsafe.Pointer(&cNames[0])),
		C.int(len(typeNames)),
		&results[0],
	); ret != 0 {
		return nil, fmt.Errorf("failed to introspect enums")
	}

	defs := make([]EnumDef, len(typeNames))
	for i := range typeNames {
		r := &results[i]
		defs[i].CName = C.GoString(&r.c_name[0])
		for j := 0; j < int(r.n_values); j++ {
			ev := &r.values[j]
			defs[i].Values = append(defs[i].Values, EnumValue{
				CName: C.GoString(&ev.c_name[0]),
				Nick:  C.GoString(&ev.nick[0]),
				Value: int(ev.value),
			})
		}
	}

	return defs, nil
}

// CollectEnumTypes returns a deduplicated, sorted list of all enum type names
// referenced by the given operations.
func CollectEnumTypes(ops []OpDef) []string {
	seen := make(map[string]bool)
	for _, op := range ops {
		for _, arg := range op.Args {
			if arg.EnumType != "" {
				seen[arg.EnumType] = true
			}
		}
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
