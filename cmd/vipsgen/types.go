package main

// ArgFlags describes argument properties from vips introspection.
type ArgFlags int

const (
	ArgInput    ArgFlags = 1 << 0
	ArgOutput   ArgFlags = 1 << 1
	ArgRequired ArgFlags = 1 << 2
	ArgModify   ArgFlags = 1 << 3
)

// ArgType represents the GType category of an argument.
type ArgType int

const (
	ArgTypeUnknown ArgType = iota
	ArgTypeImage
	ArgTypeDouble
	ArgTypeInt
	ArgTypeBool
	ArgTypeString
	ArgTypeEnum
	ArgTypeFlags
	ArgTypeArrayDouble
	ArgTypeArrayInt
	ArgTypeArrayImage
	ArgTypeBlob
	ArgTypeInterpolate
	ArgTypeSource
	ArgTypeTarget
)

// ArgDef describes a single argument to a vips operation.
type ArgDef struct {
	Name     string  // vips argument name (e.g. "in", "out", "sigma")
	Type     ArgType // argument type category
	Flags    ArgFlags
	Priority int     // vips argument ordering priority
	Default  float64 // default value for numeric types
	Min      float64
	Max      float64
	EnumType string // GType name for enum args (e.g. "VipsKernel")
}

// IsInput returns true if this argument is an input.
func (a *ArgDef) IsInput() bool {
	return a.Flags&ArgInput != 0
}

// IsOutput returns true if this argument is an output.
func (a *ArgDef) IsOutput() bool {
	return a.Flags&ArgOutput != 0
}

// IsRequired returns true if this argument is required.
func (a *ArgDef) IsRequired() bool {
	return a.Flags&ArgRequired != 0
}

// OpDef describes a single vips operation discovered by introspection.
type OpDef struct {
	Name        string   // vips operation name (e.g. "gaussblur", "resize")
	Description string   // human-readable description
	Category    string   // category (e.g. "resample", "arithmetic")
	Args        []ArgDef // all arguments
}

// RequiredInputs returns all required input arguments.
func (op *OpDef) RequiredInputs() []ArgDef {
	var result []ArgDef
	for _, a := range op.Args {
		if a.IsInput() && a.IsRequired() {
			result = append(result, a)
		}
	}
	return result
}

// OptionalInputs returns optional input arguments.
func (op *OpDef) OptionalInputs() []ArgDef {
	var result []ArgDef
	for _, a := range op.Args {
		if a.IsInput() && !a.IsRequired() {
			result = append(result, a)
		}
	}
	return result
}

// Outputs returns output arguments.
func (op *OpDef) Outputs() []ArgDef {
	var result []ArgDef
	for _, a := range op.Args {
		if a.IsOutput() {
			result = append(result, a)
		}
	}
	return result
}

// EnumValue represents a single enum member.
type EnumValue struct {
	CName string // e.g. "VIPS_KERNEL_LANCZOS3"
	Nick  string // e.g. "lanczos3"
	Value int
}

// EnumDef describes a vips enum type.
type EnumDef struct {
	CName  string      // e.g. "VipsKernel"
	Values []EnumValue // enum members
}
