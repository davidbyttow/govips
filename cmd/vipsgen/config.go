package main

// excludeOps lists operations that should not be generated because they
// are hand-written in foreign.c/foreign.go or have complex custom logic.
var excludeOps = map[string]bool{
	// Foreign load operations (hand-written in foreign.c).
	"csvload":           true,
	"csvload_source":    true,
	"csvsave":           true,
	"csvsave_target":    true,
	"fitsload":          true,
	"fitsload_source":   true,
	"fitssave":          true,
	"gifload":           true,
	"gifload_buffer":    true,
	"gifload_source":    true,
	"gifsave":           true,
	"gifsave_buffer":    true,
	"gifsave_target":    true,
	"heifload":          true,
	"heifload_buffer":   true,
	"heifload_source":   true,
	"heifsave":          true,
	"heifsave_buffer":   true,
	"heifsave_target":   true,
	"jp2kload":          true,
	"jp2kload_buffer":   true,
	"jp2kload_source":   true,
	"jp2ksave":          true,
	"jp2ksave_buffer":   true,
	"jp2ksave_target":   true,
	"jpegload":          true,
	"jpegload_buffer":   true,
	"jpegload_source":   true,
	"jpegsave":          true,
	"jpegsave_buffer":   true,
	"jpegsave_mime":     true,
	"jpegsave_target":   true,
	"jxlload":           true,
	"jxlload_buffer":    true,
	"jxlload_source":    true,
	"jxlsave":           true,
	"jxlsave_buffer":    true,
	"jxlsave_target":    true,
	"magickload":        true,
	"magickload_buffer": true,
	"magicksave":        true,
	"magicksave_buffer": true,
	"matload":           true,
	"matrixload":        true,
	"matrixload_source": true,
	"matrixprint":       true,
	"matrixsave":        true,
	"matrixsave_target": true,
	"niftiload":         true,
	"niftiload_source":  true,
	"niftisave":         true,
	"openexrload":       true,
	"openslideload":     true,
	"openslideload_source": true,
	"pdfload":           true,
	"pdfload_buffer":    true,
	"pdfload_source":    true,
	"pngload":           true,
	"pngload_buffer":    true,
	"pngload_source":    true,
	"pngsave":           true,
	"pngsave_buffer":    true,
	"pngsave_target":    true,
	"ppmload":           true,
	"ppmload_source":    true,
	"ppmsave":           true,
	"ppmsave_target":    true,
	"radload":           true,
	"radload_buffer":    true,
	"radload_source":    true,
	"radsave":           true,
	"radsave_buffer":    true,
	"radsave_target":    true,
	"rawload":           true,
	"rawsave":           true,
	"rawsave_fd":        true,
	"svgload":           true,
	"svgload_buffer":    true,
	"svgload_source":    true,
	"tiffload":          true,
	"tiffload_buffer":   true,
	"tiffload_source":   true,
	"tiffsave":          true,
	"tiffsave_buffer":   true,
	"tiffsave_target":   true,
	"vipsload":          true,
	"vipsload_source":   true,
	"vipssave":          true,
	"vipssave_target":   true,
	"webpload":          true,
	"webpload_buffer":   true,
	"webpload_source":   true,
	"webpsave":          true,
	"webpsave_buffer":   true,
	"webpsave_target":   true,

	// Thumbnail from file/buffer (complex fallback logic in resample.go).
	"thumbnail":        true,
	"thumbnail_buffer": true,
	"thumbnail_source": true,

	// Other special-cased operations.
	"switch":       true, // reserved word in Go
	"case":         true, // reserved word in Go
	"system":       true, // shell execution, not an image op
	"profile_load": true, // internal

	// Operations with complex arg types not yet supported by the generator.
	"composite": true, // array image + array enum + x/y offset arrays

	// Hand-written with custom logic in other .go files.
	"icc_transform":   true, // color.go
	"colourspace":     true, // color.go
	"resize":          true, // resample.go (KernelAuto logic)
	"thumbnail_image": true, // resample.go (complex logic)
	"text":            true, // create.go (complex struct)
	"find_trim":       true, // arithmetic.go (custom logic)
	"getpoint":        true, // arithmetic.go (custom logic)

	// Complex output params not supported by the generator.
	"max":     true, // value+x+y+array outputs
	"min":     true, // value+x+y+array outputs
	"measure": true, // returns matrix

	// Operations requiring libvips 8.16+ (not registered as VipsOperation in older versions).
	"sdf":      true, // VipsSdfShape enum added in 8.16
	"addalpha": true, // only a VipsOperation class since 8.16

	// Hand-written multi-page handling.
	"embed": true,
	"crop":  true,

	// Draw operations (in-place mutation, array double ink).
	"draw_rect":   true,
	"draw_image":  true,
	"draw_mask":   true,
	"draw_flood":  true,
	"draw_line":   true,
	"draw_circle": true,
	"draw_smudge": true,
}

// opCategoryOverride maps operation names to their category for operations
// that are direct children of VipsOperation (no intermediate abstract class).
var opCategoryOverride = map[string]string{
	// Convolution family.
	"gaussblur": "convolution",
	"sharpen":   "convolution",
	"canny":     "convolution",
	"sobel":     "convolution",
	"prewitt":   "convolution",
	"scharr":    "convolution",

	// Arithmetic/statistics family.
	"find_trim":  "arithmetic",
	"getpoint":   "arithmetic",
	"measure":    "arithmetic",
	"case":       "conversion",
	"switch":     "conversion",

	// Colour family.
	"colourspace": "colour",
	"CMYK2XYZ":    "colour",
	"XYZ2CMYK":    "colour",

	// Resample family.
	"thumbnail_image":  "resample",
	"thumbnail":        "resample",
	"thumbnail_buffer": "resample",
	"thumbnail_source": "resample",
	"globalbalance":    "resample",
	"match":            "resample",
	"merge":            "resample",
	"mosaic":           "resample",
	"mosaic1":          "resample",
	"remosaic":         "resample",

	// Histogram family.
	"hist_equal":       "histogram",
	"hist_entropy":     "histogram",
	"hist_ismonotonic": "histogram",
	"hist_local":       "histogram",
	"hist_norm":        "histogram",
	"hist_plot":        "histogram",
	"maplut":           "histogram",
	"percent":          "histogram",
	"stdif":            "histogram",

	// Correlation operations.
	"fastcor": "convolution",
	"spcor":   "convolution",

	// Misc operations that are direct children of VipsOperation.
	"matrixinvert":   "arithmetic",
	"matrixmultiply": "arithmetic",
	"system":         "create",
	"profile_load":   "create",
}

// categoryMap normalizes category strings.
var categoryMap = map[string]string{
	"VipsArithmetic":  "arithmetic",
	"arithmetic":      "arithmetic",
	"VipsConversion":  "conversion",
	"conversion":      "conversion",
	"VipsResample":    "resample",
	"resample":        "resample",
	"VipsConvolution": "convolution",
	"convolution":     "convolution",
	"VipsColour":      "colour",
	"colour":          "colour",
	"VipsCreate":      "create",
	"create":          "create",
	"VipsDraw":        "draw",
	"draw":            "draw",
	"VipsMorphology":  "morphology",
	"morphology":      "morphology",
	"VipsForeign":     "foreign",
	"foreign":         "foreign",
	"VipsFreqfilt":    "freqfilt",
	"freqfilt":        "freqfilt",
	"VipsHistogram":   "histogram",
	"histogram":       "histogram",
}

// normalizeCategory maps an introspected category to our standard name.
func normalizeCategory(cat string) string {
	if mapped, ok := categoryMap[cat]; ok {
		return mapped
	}
	return cat
}

// normalizeCategoryForOp returns the category for an operation, checking
// the op-level override first, then the category map.
func normalizeCategoryForOp(opName, cat string) string {
	if override, ok := opCategoryOverride[opName]; ok {
		return override
	}
	return normalizeCategory(cat)
}

// enumGoName maps a C enum type name to the Go type name we use.
var enumGoName = map[string]string{
	"VipsKernel":              "Kernel",
	"VipsSize":                "Size",
	"VipsDirection":           "Direction",
	"VipsAngle":               "Angle",
	"VipsAngle45":             "Angle45",
	"VipsBandFormat":          "BandFormat",
	"VipsBlendMode":           "BlendMode",
	"VipsCoding":              "Coding",
	"VipsCompassDirection":    "Gravity",
	"VipsExtend":              "ExtendStrategy",
	"VipsInteresting":         "Interesting",
	"VipsInterpretation":      "Interpretation",
	"VipsIntent":              "Intent",
	"VipsPrecision":           "Precision",
	"VipsAlign":               "Align",
	"VipsTextWrap":            "TextWrap",
	"VipsCombineMode":         "CombineMode",
	"VipsCombine":             "Combine",
	"VipsOperationBoolean":    "OperationBoolean",
	"VipsOperationMath":       "OperationMath",
	"VipsOperationMath2":      "OperationMath2",
	"VipsOperationComplex":    "OperationComplex",
	"VipsOperationComplex2":   "OperationComplex2",
	"VipsOperationComplexget": "OperationComplexget",
	"VipsOperationRelational": "OperationRelational",
	"VipsOperationRound":      "OperationRound",
	"VipsOperationMorphology": "OperationMorphology",
	"VipsForeignDzLayout":     "ForeignDzLayout",
	"VipsForeignDzDepth":      "ForeignDzDepth",
	"VipsForeignDzContainer":  "ForeignDzContainer",
	"VipsForeignTiffCompression": "TiffCompression",
	"VipsForeignTiffPredictor":   "TiffPredictor",
	"VipsForeignPngFilter":       "PngFilter",
	"VipsForeignSubsample":       "SubsampleMode",
	"VipsForeignHeifCompression": "HeifCompression",
	"VipsRegionShrink":           "RegionShrink",
	"VipsFalseColour":            "FalseColour",
	// VipsAccess is defined as bare int constants in foreign.go, not a named type.
	// Use int for now; will become a proper type in Phase 4.
	// "VipsAccess":                 "Access",
	"VipsDemandStyle":            "DemandStyle",
	"VipsFailOn":                 "FailOn",
	"VipsForeignKeep":            "ForeignKeep",
}

// goEnumName returns the Go type name for a C enum type.
func goEnumName(cName string) string {
	if name, ok := enumGoName[cName]; ok {
		return name
	}
	return ""
}

// versionIntroduced maps operation names to the libvips version that
// introduced them. Only operations added after 8.10 need to be listed.
var versionIntroduced = map[string]struct{ Major, Minor int }{
	"jp2kload":        {8, 11},
	"jp2kload_buffer": {8, 11},
	"jp2kload_source": {8, 11},
	"jp2ksave":        {8, 11},
	"jp2ksave_buffer": {8, 11},
	"jp2ksave_target": {8, 11},
	"transpose3d":     {8, 14},
}
