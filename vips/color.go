package vips

// #cgo pkg-config: vips
// #include "color.h"
import "C"
import "unsafe"

// Color represents an RGB
type Color struct {
	R, G, B uint8
}

// Interpretation represents VIPS_INTERPRETATION type
type Interpretation int

// Interpretation enum
const (
	InterpretationError     Interpretation = C.VIPS_INTERPRETATION_ERROR
	InterpretationMultiband Interpretation = C.VIPS_INTERPRETATION_MULTIBAND
	InterpretationBW        Interpretation = C.VIPS_INTERPRETATION_B_W
	InterpretationHistogram Interpretation = C.VIPS_INTERPRETATION_HISTOGRAM
	InterpretationXYZ       Interpretation = C.VIPS_INTERPRETATION_XYZ
	InterpretationLAB       Interpretation = C.VIPS_INTERPRETATION_LAB
	InterpretationCMYK      Interpretation = C.VIPS_INTERPRETATION_CMYK
	InterpretationLABQ      Interpretation = C.VIPS_INTERPRETATION_LABQ
	InterpretationRGB       Interpretation = C.VIPS_INTERPRETATION_RGB
	InterpretationRGB16     Interpretation = C.VIPS_INTERPRETATION_RGB16
	InterpretationCMC       Interpretation = C.VIPS_INTERPRETATION_CMC
	InterpretationLCH       Interpretation = C.VIPS_INTERPRETATION_LCH
	InterpretationLABS      Interpretation = C.VIPS_INTERPRETATION_LABS
	InterpretationSRGB      Interpretation = C.VIPS_INTERPRETATION_sRGB
	InterpretationYXY       Interpretation = C.VIPS_INTERPRETATION_YXY
	InterpretationFourier   Interpretation = C.VIPS_INTERPRETATION_FOURIER
	InterpretationGB16      Interpretation = C.VIPS_INTERPRETATION_RGB16
	InterpretationGrey16    Interpretation = C.VIPS_INTERPRETATION_GREY16
	InterpretationMatrix    Interpretation = C.VIPS_INTERPRETATION_MATRIX
	InterpretationScRGB     Interpretation = C.VIPS_INTERPRETATION_scRGB
	InterpretationHSV       Interpretation = C.VIPS_INTERPRETATION_HSV
)

func vipsIsColorSpaceSupported(in *C.VipsImage) bool {
	return C.is_colorspace_supported(in) == 1
}

// https://libvips.github.io/libvips/API/current/libvips-colour.html#vips-colourspace
func vipsToColorSpace(in *C.VipsImage, interpretation Interpretation) (*C.VipsImage, error) {
	incOpCounter("to_colorspace")
	var out *C.VipsImage

	inter := C.VipsInterpretation(interpretation)

	if err := C.to_colorspace(in, &out, inter); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

func vipsOptimizeICCProfile(in *C.VipsImage, inputProfile string) (*C.VipsImage, error) {
	var out *C.VipsImage
	var profile *C.char

	if inputProfile != "" {
		profile = C.CString(inputProfile)
		defer C.free(unsafe.Pointer(profile))
	}

	if res := int(C.optimize_icc_profile(in, &out, profile)); res != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}
