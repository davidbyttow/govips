package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// Generate creates C wrapper, C header, and Go bridge files from the
// discovered operations, outputting a single generated.{c,h,go} triplet.
func Generate(ops []OpDef, outputDir string) error {
	// Clean up old per-category generated files.
	cleanOldGenFiles(outputDir)

	// Filter to generatable, non-foreign operations and sort alphabetically.
	var allOps []OpDef
	for _, op := range ops {
		if excludeOps[op.Name] {
			continue
		}
		cat := normalizeCategoryForOp(op.Name, op.Category)
		if cat == "foreign" {
			continue
		}
		allOps = append(allOps, op)
	}
	sort.Slice(allOps, func(i, j int) bool {
		return allOps[i].Name < allOps[j].Name
	})

	cCode := genCSource(allOps)
	hCode := genCHeader(allOps)
	goCode := genGoBridge(allOps)

	cPath := filepath.Join(outputDir, "generated.c")
	hPath := filepath.Join(outputDir, "generated.h")
	goPath := filepath.Join(outputDir, "generated.go")

	if err := writeIfChanged(cPath, cCode); err != nil {
		return err
	}
	if err := writeIfChanged(hPath, hCode); err != nil {
		return err
	}
	if err := writeIfChanged(goPath, goCode); err != nil {
		return err
	}

	fmt.Printf("  generated.{c,h,go}: %d operations\n", len(allOps))
	return nil
}

// cleanOldGenFiles removes legacy per-category gen_*.{c,h,go} and gen_helpers.go
// files, but preserves hand-written files like gen_enum_extras.go.
func cleanOldGenFiles(outputDir string) {
	preserve := map[string]bool{
		"gen_enum_extras.go": true,
	}
	for _, ext := range []string{"*.c", "*.h", "*.go"} {
		matches, _ := filepath.Glob(filepath.Join(outputDir, "gen_"+ext))
		for _, m := range matches {
			if preserve[filepath.Base(m)] {
				continue
			}
			os.Remove(m)
		}
	}
}

func writeIfChanged(path string, content []byte) error {
	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}
	return os.WriteFile(path, content, 0644)
}

// --- Name conversion helpers ---

func goName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	parts := strings.Split(name, "_")
	var result strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		result.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return result.String()
}

func goArgName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	parts := strings.Split(name, "_")
	var result strings.Builder
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			result.WriteString(p)
		} else {
			result.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	s := result.String()
	// Avoid Go keywords.
	switch s {
	case "type":
		return "typ"
	case "func":
		return "fn"
	case "map":
		return "mapVal"
	case "range":
		return "rangeVal"
	case "in":
		return "input"
	}
	return s
}

func goExportedArgName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	parts := strings.Split(name, "_")
	var result strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		result.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	s := result.String()
	// Handle Go keywords.
	switch s {
	case "Type":
		return "Typ"
	case "Func":
		return "Fn"
	case "Map":
		return "MapVal"
	case "Range":
		return "RangeVal"
	case "In":
		return "Input"
	}
	return s
}

func cFuncName(opName string) string {
	return "gen_vips_" + opName
}

func goFuncName(opName string) string {
	return "vipsGen" + goName(opName)
}

func cStructName(opName string) string {
	return "Gen" + goName(opName) + "Opts"
}

func goOptsTypeName(opName string) string {
	return goName(opName) + "Options"
}

func isFirstUpper(s string) bool {
	for _, r := range s {
		return unicode.IsUpper(r)
	}
	return false
}

// --- Go type mapping ---

func goTypeName(arg ArgDef) string {
	switch arg.Type {
	case ArgTypeImage:
		return "*C.VipsImage"
	case ArgTypeDouble:
		return "float64"
	case ArgTypeInt:
		return "int"
	case ArgTypeBool:
		return "bool"
	case ArgTypeString:
		return "string"
	case ArgTypeEnum:
		if name := goEnumName(arg.EnumType); name != "" {
			return name
		}
		return "int"
	case ArgTypeFlags:
		return "int"
	case ArgTypeArrayDouble:
		return "[]float64"
	case ArgTypeArrayInt:
		return "[]int"
	case ArgTypeArrayImage:
		return "[]*C.VipsImage"
	case ArgTypeInterpolate:
		return "*C.VipsInterpolate"
	case ArgTypeBlob:
		return "[]byte"
	default:
		return "interface{}"
	}
}

func goOptTypeName(arg ArgDef) string {
	switch arg.Type {
	case ArgTypeDouble:
		return "*float64"
	case ArgTypeInt:
		return "*int"
	case ArgTypeBool:
		return "*bool"
	case ArgTypeString:
		return "*string"
	case ArgTypeEnum:
		if name := goEnumName(arg.EnumType); name != "" {
			return "*" + name
		}
		return "*int"
	case ArgTypeFlags:
		return "*int"
	case ArgTypeArrayDouble:
		return "[]float64"
	case ArgTypeArrayInt:
		return "[]int"
	case ArgTypeArrayImage:
		return "[]*C.VipsImage"
	case ArgTypeImage:
		return "*C.VipsImage"
	case ArgTypeInterpolate:
		return "*C.VipsInterpolate"
	case ArgTypeBlob:
		return "[]byte"
	default:
		return "interface{}"
	}
}

func cTypeName(arg ArgDef) string {
	switch arg.Type {
	case ArgTypeImage:
		return "VipsImage *"
	case ArgTypeDouble:
		return "double"
	case ArgTypeInt:
		return "int"
	case ArgTypeBool:
		return "int"
	case ArgTypeString:
		return "const char *"
	case ArgTypeEnum:
		return arg.EnumType
	case ArgTypeFlags:
		return "int"
	case ArgTypeArrayDouble:
		return "double *"
	case ArgTypeArrayInt:
		return "int *"
	case ArgTypeArrayImage:
		return "VipsArrayImage *"
	case ArgTypeInterpolate:
		return "VipsInterpolate *"
	case ArgTypeBlob:
		return "void *"
	default:
		return "void *"
	}
}

func cOutputTypeName(arg ArgDef) string {
	switch arg.Type {
	case ArgTypeImage:
		return "VipsImage **"
	case ArgTypeDouble:
		return "double *"
	case ArgTypeInt:
		return "int *"
	case ArgTypeBool:
		return "int *"
	case ArgTypeEnum:
		return "int *"
	default:
		return "void *"
	}
}

// isNilCheckType returns true if the opt arg uses nil check (pointer/slice).
func isNilCheckType(arg ArgDef) bool {
	return arg.Type == ArgTypeArrayDouble || arg.Type == ArgTypeArrayInt ||
		arg.Type == ArgTypeArrayImage || arg.Type == ArgTypeBlob ||
		arg.Type == ArgTypeImage || arg.Type == ArgTypeInterpolate
}

// hasOutputs returns true if the operation has output args.
func hasOutputs(op OpDef) bool {
	return len(op.Outputs()) > 0
}

// primaryImageOutput returns the first image output.
func primaryImageOutput(op OpDef) *ArgDef {
	for _, a := range op.Outputs() {
		if a.Type == ArgTypeImage {
			return &a
		}
	}
	return nil
}

// --- C Source Generation ---

func genCSource(ops []OpDef) []byte {
	var w bytes.Buffer
	fmt.Fprintf(&w, "// Code generated by vipsgen. DO NOT EDIT.\n")
	fmt.Fprintf(&w, "#include \"generated.h\"\n\n")

	for _, op := range ops {
		genCFunc(&w, op)
	}
	return w.Bytes()
}

func genCFunc(w *bytes.Buffer, op OpDef) {
	reqInputs := op.RequiredInputs()
	optInputs := op.OptionalInputs()
	outputs := op.Outputs()
	hasOpts := len(optInputs) > 0

	// Function signature.
	fmt.Fprintf(w, "int %s(", cFuncName(op.Name))

	params := []string{}
	for _, a := range reqInputs {
		switch a.Type {
		case ArgTypeArrayDouble, ArgTypeArrayInt:
			params = append(params, fmt.Sprintf("%s %s, int %s_n", cTypeName(a), goArgName(a.Name), goArgName(a.Name)))
		case ArgTypeArrayImage:
			params = append(params, fmt.Sprintf("VipsImage **%s, int %s_n", goArgName(a.Name), goArgName(a.Name)))
		default:
			params = append(params, fmt.Sprintf("%s %s", cTypeName(a), goArgName(a.Name)))
		}
	}
	for _, a := range outputs {
		params = append(params, fmt.Sprintf("%s out_%s", cOutputTypeName(a), goArgName(a.Name)))
	}
	if hasOpts {
		params = append(params, fmt.Sprintf("%s *opts", cStructName(op.Name)))
	}
	fmt.Fprintf(w, "%s) {\n", strings.Join(params, ", "))

	// Body.
	fmt.Fprintf(w, "    VipsOperation *op = vips_operation_new(\"%s\");\n", op.Name)
	fmt.Fprintf(w, "    if (!op) return -1;\n\n")

	// Set required inputs.
	for _, a := range reqInputs {
		genCSetArg(w, a, goArgName(a.Name), false)
	}

	// Set optional inputs.
	if hasOpts {
		fmt.Fprintf(w, "\n    if (opts) {\n")
		for _, a := range optInputs {
			argName := goArgName(a.Name)
			fmt.Fprintf(w, "        if (opts->has_%s) {\n", argName)
			genCSetOptArg(w, a, "opts->"+argName)
			fmt.Fprintf(w, "        }\n")
		}
		fmt.Fprintf(w, "    }\n")
	}

	// Build.
	fmt.Fprintf(w, "\n    if (vips_cache_operation_buildp(&op)) goto error;\n\n")

	// Extract outputs.
	for _, a := range outputs {
		fmt.Fprintf(w, "    g_object_get(VIPS_OBJECT(op), \"%s\", out_%s, NULL);\n",
			a.Name, goArgName(a.Name))
	}

	// Cleanup.
	fmt.Fprintf(w, "\n    vips_object_unref_outputs(VIPS_OBJECT(op));\n")
	fmt.Fprintf(w, "    g_object_unref(op);\n")
	fmt.Fprintf(w, "    return 0;\n\n")
	fmt.Fprintf(w, "error:\n")
	fmt.Fprintf(w, "    vips_object_unref_outputs(VIPS_OBJECT(op));\n")
	fmt.Fprintf(w, "    g_object_unref(op);\n")
	fmt.Fprintf(w, "    return -1;\n")
	fmt.Fprintf(w, "}\n\n")
}

func genCSetArg(w *bytes.Buffer, a ArgDef, varName string, isOpt bool) {
	switch a.Type {
	case ArgTypeArrayDouble:
		fmt.Fprintf(w, "    {\n")
		fmt.Fprintf(w, "        VipsArrayDouble *arr = vips_array_double_new(%s, %s_n);\n", varName, varName)
		fmt.Fprintf(w, "        int ret = vips_object_set(VIPS_OBJECT(op), \"%s\", arr, NULL);\n", a.Name)
		fmt.Fprintf(w, "        vips_area_unref(VIPS_AREA(arr));\n")
		fmt.Fprintf(w, "        if (ret) goto error;\n")
		fmt.Fprintf(w, "    }\n")
	case ArgTypeArrayInt:
		fmt.Fprintf(w, "    {\n")
		fmt.Fprintf(w, "        VipsArrayInt *arr = vips_array_int_new(%s, %s_n);\n", varName, varName)
		fmt.Fprintf(w, "        int ret = vips_object_set(VIPS_OBJECT(op), \"%s\", arr, NULL);\n", a.Name)
		fmt.Fprintf(w, "        vips_area_unref(VIPS_AREA(arr));\n")
		fmt.Fprintf(w, "        if (ret) goto error;\n")
		fmt.Fprintf(w, "    }\n")
	case ArgTypeArrayImage:
		fmt.Fprintf(w, "    {\n")
		fmt.Fprintf(w, "        VipsArrayImage *arr = vips_array_image_new(%s, %s_n);\n", varName, varName)
		fmt.Fprintf(w, "        int ret = vips_object_set(VIPS_OBJECT(op), \"%s\", arr, NULL);\n", a.Name)
		fmt.Fprintf(w, "        vips_area_unref(VIPS_AREA(arr));\n")
		fmt.Fprintf(w, "        if (ret) goto error;\n")
		fmt.Fprintf(w, "    }\n")
	case ArgTypeBool:
		fmt.Fprintf(w, "    if (vips_object_set(VIPS_OBJECT(op), \"%s\", (gboolean)%s, NULL)) goto error;\n", a.Name, varName)
	case ArgTypeEnum:
		fmt.Fprintf(w, "    if (vips_object_set(VIPS_OBJECT(op), \"%s\", (int)%s, NULL)) goto error;\n", a.Name, varName)
	default:
		fmt.Fprintf(w, "    if (vips_object_set(VIPS_OBJECT(op), \"%s\", %s, NULL)) goto error;\n", a.Name, varName)
	}
}

func genCSetOptArg(w *bytes.Buffer, a ArgDef, varExpr string) {
	switch a.Type {
	case ArgTypeArrayDouble:
		fmt.Fprintf(w, "            {\n")
		fmt.Fprintf(w, "                VipsArrayDouble *arr = vips_array_double_new(%s, %s_n);\n", varExpr, varExpr)
		fmt.Fprintf(w, "                vips_object_set(VIPS_OBJECT(op), \"%s\", arr, NULL);\n", a.Name)
		fmt.Fprintf(w, "                vips_area_unref(VIPS_AREA(arr));\n")
		fmt.Fprintf(w, "            }\n")
	case ArgTypeArrayInt:
		fmt.Fprintf(w, "            {\n")
		fmt.Fprintf(w, "                VipsArrayInt *arr = vips_array_int_new(%s, %s_n);\n", varExpr, varExpr)
		fmt.Fprintf(w, "                vips_object_set(VIPS_OBJECT(op), \"%s\", arr, NULL);\n", a.Name)
		fmt.Fprintf(w, "                vips_area_unref(VIPS_AREA(arr));\n")
		fmt.Fprintf(w, "            }\n")
	case ArgTypeBool:
		fmt.Fprintf(w, "            vips_object_set(VIPS_OBJECT(op), \"%s\", (gboolean)%s, NULL);\n", a.Name, varExpr)
	case ArgTypeEnum:
		fmt.Fprintf(w, "            vips_object_set(VIPS_OBJECT(op), \"%s\", (int)%s, NULL);\n", a.Name, varExpr)
	default:
		fmt.Fprintf(w, "            vips_object_set(VIPS_OBJECT(op), \"%s\", %s, NULL);\n", a.Name, varExpr)
	}
}

// --- C Header Generation ---

func genCHeader(ops []OpDef) []byte {
	var w bytes.Buffer
	fmt.Fprintf(&w, "// Code generated by vipsgen. DO NOT EDIT.\n")
	fmt.Fprintf(&w, "#ifndef GENERATED_H\n")
	fmt.Fprintf(&w, "#define GENERATED_H\n\n")
	fmt.Fprintf(&w, "#include <stdlib.h>\n")
	fmt.Fprintf(&w, "#include <vips/vips.h>\n\n")

	for _, op := range ops {
		genCHeaderOp(&w, op)
	}

	fmt.Fprintf(&w, "#endif\n")
	return w.Bytes()
}

func genCHeaderOp(w *bytes.Buffer, op OpDef) {
	reqInputs := op.RequiredInputs()
	optInputs := op.OptionalInputs()
	outputs := op.Outputs()
	hasOpts := len(optInputs) > 0

	// Generate opts struct if needed.
	if hasOpts {
		fmt.Fprintf(w, "typedef struct {\n")
		for _, a := range optInputs {
			argName := goArgName(a.Name)
			fmt.Fprintf(w, "    int has_%s;\n", argName)
			switch a.Type {
			case ArgTypeImage:
				fmt.Fprintf(w, "    VipsImage *%s;\n", argName)
			case ArgTypeDouble:
				fmt.Fprintf(w, "    double %s;\n", argName)
			case ArgTypeInt:
				fmt.Fprintf(w, "    int %s;\n", argName)
			case ArgTypeBool:
				fmt.Fprintf(w, "    int %s;\n", argName)
			case ArgTypeString:
				fmt.Fprintf(w, "    const char *%s;\n", argName)
			case ArgTypeEnum:
				fmt.Fprintf(w, "    %s %s;\n", a.EnumType, argName)
			case ArgTypeArrayDouble:
				fmt.Fprintf(w, "    double *%s; int %s_n;\n", argName, argName)
			case ArgTypeArrayInt:
				fmt.Fprintf(w, "    int *%s; int %s_n;\n", argName, argName)
			case ArgTypeInterpolate:
				fmt.Fprintf(w, "    VipsInterpolate *%s;\n", argName)
			default:
				fmt.Fprintf(w, "    int %s;\n", argName)
			}
		}
		fmt.Fprintf(w, "} %s;\n\n", cStructName(op.Name))
	}

	// Function prototype.
	params := []string{}
	for _, a := range reqInputs {
		switch a.Type {
		case ArgTypeArrayDouble, ArgTypeArrayInt:
			params = append(params, fmt.Sprintf("%s %s, int %s_n", cTypeName(a), goArgName(a.Name), goArgName(a.Name)))
		case ArgTypeArrayImage:
			params = append(params, fmt.Sprintf("VipsImage **%s, int %s_n", goArgName(a.Name), goArgName(a.Name)))
		default:
			params = append(params, fmt.Sprintf("%s %s", cTypeName(a), goArgName(a.Name)))
		}
	}
	for _, a := range outputs {
		params = append(params, fmt.Sprintf("%s out_%s", cOutputTypeName(a), goArgName(a.Name)))
	}
	if hasOpts {
		params = append(params, fmt.Sprintf("%s *opts", cStructName(op.Name)))
	}
	if len(params) == 0 {
		params = append(params, "void")
	}
	fmt.Fprintf(w, "int %s(%s);\n\n", cFuncName(op.Name), strings.Join(params, ", "))
}

// --- Go Bridge Generation ---

// opHasArrayOpts returns true if the operation has any optional array params.
func opHasArrayOpts(op OpDef) bool {
	for _, a := range op.OptionalInputs() {
		if a.Type == ArgTypeArrayDouble || a.Type == ArgTypeArrayInt {
			return true
		}
	}
	return false
}

func genGoBridge(ops []OpDef) []byte {
	var w bytes.Buffer
	fmt.Fprintf(&w, "// Code generated by vipsgen. DO NOT EDIT.\n")
	fmt.Fprintf(&w, "package vips\n\n")
	fmt.Fprintf(&w, "// #include \"generated.h\"\n")
	fmt.Fprintf(&w, "import \"C\"\n\n")

	// Check if any op in this category needs runtime.Pinner for array opts.
	needsRuntime := false
	for _, op := range ops {
		if opHasArrayOpts(op) {
			needsRuntime = true
			break
		}
	}
	if needsRuntime {
		fmt.Fprintf(&w, "import (\n\t\"runtime\"\n\t\"unsafe\"\n)\n\n")
	} else {
		fmt.Fprintf(&w, "import \"unsafe\"\n\n")
	}

	fmt.Fprintf(&w, "// Ensure imports are used.\n")
	fmt.Fprintf(&w, "var _ = unsafe.Pointer(nil)\n\n")

	for _, op := range ops {
		genGoFunc(&w, op)
	}
	return w.Bytes()
}

func genGoFunc(w *bytes.Buffer, op OpDef) {
	reqInputs := op.RequiredInputs()
	optInputs := op.OptionalInputs()
	outputs := op.Outputs()
	hasOpts := len(optInputs) > 0
	hasOutputs := len(outputs) > 0

	// Generate options struct.
	if hasOpts {
		fmt.Fprintf(w, "// %s are optional parameters for %s.\n", goOptsTypeName(op.Name), op.Name)
		fmt.Fprintf(w, "type %s struct {\n", goOptsTypeName(op.Name))
		for _, a := range optInputs {
			fmt.Fprintf(w, "\t%s %s\n", goExportedArgName(a.Name), goOptTypeName(a))
		}
		fmt.Fprintf(w, "}\n\n")
	}

	// Function signature.
	funcName := goFuncName(op.Name)
	fmt.Fprintf(w, "// %s calls the vips %s operation.\n", funcName, op.Name)
	fmt.Fprintf(w, "// %s\n", op.Description)
	fmt.Fprintf(w, "func %s(", funcName)

	// Parameters.
	goParams := []string{}
	for _, a := range reqInputs {
		goParams = append(goParams, fmt.Sprintf("%s %s", goArgName(a.Name), goTypeName(a)))
	}
	if hasOpts {
		goParams = append(goParams, fmt.Sprintf("opts *%s", goOptsTypeName(op.Name)))
	}
	fmt.Fprintf(w, "%s) (", strings.Join(goParams, ", "))

	// Return types.
	if hasOutputs {
		retTypes := []string{}
		for _, a := range outputs {
			retTypes = append(retTypes, goTypeName(a))
		}
		retTypes = append(retTypes, "error")
		fmt.Fprintf(w, "%s", strings.Join(retTypes, ", "))
	} else {
		fmt.Fprintf(w, "error")
	}
	fmt.Fprintf(w, ") {\n")

	// Body.
	fmt.Fprintf(w, "\tincOpCounter(\"%s\")\n\n", op.Name)

	// Declare string temp vars for required string inputs.
	for _, a := range reqInputs {
		if a.Type == ArgTypeString {
			fmt.Fprintf(w, "\tcStr_%s := C.CString(%s)\n", goArgName(a.Name), goArgName(a.Name))
			fmt.Fprintf(w, "\tdefer C.free(unsafe.Pointer(cStr_%s))\n", goArgName(a.Name))
		}
	}

	// Declare output vars.
	for _, a := range outputs {
		argName := goArgName(a.Name)
		switch a.Type {
		case ArgTypeImage:
			fmt.Fprintf(w, "\tvar out_%s *C.VipsImage\n", argName)
		case ArgTypeDouble:
			fmt.Fprintf(w, "\tvar out_%s C.double\n", argName)
		case ArgTypeInt:
			fmt.Fprintf(w, "\tvar out_%s C.int\n", argName)
		case ArgTypeBool:
			fmt.Fprintf(w, "\tvar out_%s C.int\n", argName)
		case ArgTypeEnum:
			fmt.Fprintf(w, "\tvar out_%s C.int\n", argName)
		}
	}

	// Populate opts struct.
	needsPinner := opHasArrayOpts(op)
	if hasOpts {
		fmt.Fprintf(w, "\n\tvar cOpts C.%s\n", cStructName(op.Name))
		if needsPinner {
			fmt.Fprintf(w, "\tvar pinner runtime.Pinner\n")
			fmt.Fprintf(w, "\tdefer pinner.Unpin()\n")
		}
		fmt.Fprintf(w, "\tif opts != nil {\n")
		for _, a := range optInputs {
			argName := goArgName(a.Name)
			exportedName := goExportedArgName(a.Name)

			if isNilCheckType(a) {
				fmt.Fprintf(w, "\t\tif opts.%s != nil {\n", exportedName)
				fmt.Fprintf(w, "\t\t\tcOpts.has_%s = 1\n", argName)
				switch a.Type {
				case ArgTypeArrayDouble:
					fmt.Fprintf(w, "\t\t\tpinner.Pin(&opts.%s[0])\n", exportedName)
					fmt.Fprintf(w, "\t\t\tcOpts.%s = (*C.double)(unsafe.Pointer(&opts.%s[0]))\n", argName, exportedName)
					fmt.Fprintf(w, "\t\t\tcOpts.%s_n = C.int(len(opts.%s))\n", argName, exportedName)
				case ArgTypeArrayInt:
					fmt.Fprintf(w, "\t\t\tpinner.Pin(&opts.%s[0])\n", exportedName)
					fmt.Fprintf(w, "\t\t\tcOpts.%s = (*C.int)(unsafe.Pointer(&opts.%s[0]))\n", argName, exportedName)
					fmt.Fprintf(w, "\t\t\tcOpts.%s_n = C.int(len(opts.%s))\n", argName, exportedName)
				case ArgTypeImage:
					fmt.Fprintf(w, "\t\t\tcOpts.%s = opts.%s\n", argName, exportedName)
				case ArgTypeInterpolate:
					fmt.Fprintf(w, "\t\t\tcOpts.%s = opts.%s\n", argName, exportedName)
				}
				fmt.Fprintf(w, "\t\t}\n")
			} else {
				fmt.Fprintf(w, "\t\tif opts.%s != nil {\n", exportedName)
				fmt.Fprintf(w, "\t\t\tcOpts.has_%s = 1\n", argName)
				switch a.Type {
				case ArgTypeDouble:
					fmt.Fprintf(w, "\t\t\tcOpts.%s = C.double(*opts.%s)\n", argName, exportedName)
				case ArgTypeInt:
					fmt.Fprintf(w, "\t\t\tcOpts.%s = C.int(*opts.%s)\n", argName, exportedName)
				case ArgTypeBool:
					fmt.Fprintf(w, "\t\t\tcOpts.%s = C.int(boolToInt(*opts.%s))\n", argName, exportedName)
				case ArgTypeString:
					fmt.Fprintf(w, "\t\t\ttmp_%s := C.CString(*opts.%s)\n", argName, exportedName)
					fmt.Fprintf(w, "\t\t\tdefer C.free(unsafe.Pointer(tmp_%s))\n", argName)
					fmt.Fprintf(w, "\t\t\tcOpts.%s = tmp_%s\n", argName, argName)
				case ArgTypeEnum:
					fmt.Fprintf(w, "\t\t\tcOpts.%s = C.%s(*opts.%s)\n", argName, a.EnumType, exportedName)
				default:
					fmt.Fprintf(w, "\t\t\tcOpts.%s = C.int(*opts.%s)\n", argName, exportedName)
				}
				fmt.Fprintf(w, "\t\t}\n")
			}
		}
		fmt.Fprintf(w, "\t}\n")
	}

	// C function call.
	fmt.Fprintf(w, "\n\tret := C.%s(", cFuncName(op.Name))
	cArgs := []string{}
	for _, a := range reqInputs {
		argName := goArgName(a.Name)
		switch a.Type {
		case ArgTypeImage:
			cArgs = append(cArgs, argName)
		case ArgTypeDouble:
			cArgs = append(cArgs, fmt.Sprintf("C.double(%s)", argName))
		case ArgTypeInt:
			cArgs = append(cArgs, fmt.Sprintf("C.int(%s)", argName))
		case ArgTypeBool:
			cArgs = append(cArgs, fmt.Sprintf("C.int(boolToInt(%s))", argName))
		case ArgTypeString:
			cArgs = append(cArgs, fmt.Sprintf("cStr_%s", argName))
		case ArgTypeEnum:
			cArgs = append(cArgs, fmt.Sprintf("C.%s(%s)", a.EnumType, argName))
		case ArgTypeArrayDouble:
			cArgs = append(cArgs, fmt.Sprintf("(*C.double)(unsafe.Pointer(&%s[0]))", argName))
			cArgs = append(cArgs, fmt.Sprintf("C.int(len(%s))", argName))
		case ArgTypeArrayInt:
			cArgs = append(cArgs, fmt.Sprintf("(*C.int)(unsafe.Pointer(&%s[0]))", argName))
			cArgs = append(cArgs, fmt.Sprintf("C.int(len(%s))", argName))
		case ArgTypeArrayImage:
			cArgs = append(cArgs, fmt.Sprintf("(**C.VipsImage)(unsafe.Pointer(&%s[0]))", argName))
			cArgs = append(cArgs, fmt.Sprintf("C.int(len(%s))", argName))
		default:
			cArgs = append(cArgs, argName)
		}
	}
	for _, a := range outputs {
		cArgs = append(cArgs, fmt.Sprintf("&out_%s", goArgName(a.Name)))
	}
	if hasOpts {
		cArgs = append(cArgs, "&cOpts")
	}
	fmt.Fprintf(w, "%s)\n", strings.Join(cArgs, ", "))

	// Error handling.
	fmt.Fprintf(w, "\tif ret != 0 {\n")
	if hasOutputs {
		errVals := []string{}
		for _, a := range outputs {
			errVals = append(errVals, goZeroValue(a))
		}
		imgOut := primaryImageOutput(op)
		if imgOut != nil {
			errVals = append(errVals, fmt.Sprintf("handleImageError(out_%s)", goArgName(imgOut.Name)))
		} else {
			errVals = append(errVals, "handleVipsError()")
		}
		fmt.Fprintf(w, "\t\treturn %s\n", strings.Join(errVals, ", "))
	} else {
		fmt.Fprintf(w, "\t\treturn handleVipsError()\n")
	}
	fmt.Fprintf(w, "\t}\n\n")

	// Return.
	if hasOutputs {
		retVals := []string{}
		for _, a := range outputs {
			argName := goArgName(a.Name)
			switch a.Type {
			case ArgTypeImage:
				retVals = append(retVals, fmt.Sprintf("out_%s", argName))
			case ArgTypeDouble:
				retVals = append(retVals, fmt.Sprintf("float64(out_%s)", argName))
			case ArgTypeInt:
				retVals = append(retVals, fmt.Sprintf("int(out_%s)", argName))
			case ArgTypeBool:
				retVals = append(retVals, fmt.Sprintf("out_%s != 0", argName))
			case ArgTypeEnum:
				goEnum := goEnumName(a.EnumType)
				if goEnum == "" {
					goEnum = "int"
				}
				retVals = append(retVals, fmt.Sprintf("%s(out_%s)", goEnum, argName))
			default:
				retVals = append(retVals, fmt.Sprintf("out_%s", argName))
			}
		}
		retVals = append(retVals, "nil")
		fmt.Fprintf(w, "\treturn %s\n", strings.Join(retVals, ", "))
	} else {
		fmt.Fprintf(w, "\treturn nil\n")
	}
	fmt.Fprintf(w, "}\n\n")
}

func goZeroValue(a ArgDef) string {
	switch a.Type {
	case ArgTypeImage:
		return "nil"
	case ArgTypeDouble:
		return "0"
	case ArgTypeInt:
		return "0"
	case ArgTypeBool:
		return "false"
	case ArgTypeEnum:
		return "0"
	default:
		return "nil"
	}
}

