package vips

// #cgo CFLAGS: -std=c99
// #include "conversion.h"
import "C"

import "errors"

// BandFormat represents VIPS_FORMAT type
type BandFormat int

// BandFormat enum
const (
	BandFormatNotSet    BandFormat = C.VIPS_FORMAT_NOTSET
	BandFormatUchar     BandFormat = C.VIPS_FORMAT_UCHAR
	BandFormatChar      BandFormat = C.VIPS_FORMAT_CHAR
	BandFormatUshort    BandFormat = C.VIPS_FORMAT_USHORT
	BandFormatShort     BandFormat = C.VIPS_FORMAT_SHORT
	BandFormatUint      BandFormat = C.VIPS_FORMAT_UINT
	BandFormatInt       BandFormat = C.VIPS_FORMAT_INT
	BandFormatFloat     BandFormat = C.VIPS_FORMAT_FLOAT
	BandFormatComplex   BandFormat = C.VIPS_FORMAT_COMPLEX
	BandFormatDouble    BandFormat = C.VIPS_FORMAT_DOUBLE
	BandFormatDpComplex BandFormat = C.VIPS_FORMAT_DPCOMPLEX
)

// BlendMode gives the various Porter-Duff and PDF blend modes.
// See https://libvips.github.io/libvips/API/current/libvips-conversion.html#VipsBlendMode
type BlendMode int

// Constants define the various Porter-Duff and PDF blend modes.
// See https://libvips.github.io/libvips/API/current/libvips-conversion.html#VipsBlendMode
const (
	BlendModeClear      BlendMode = C.VIPS_BLEND_MODE_CLEAR
	BlendModeSource     BlendMode = C.VIPS_BLEND_MODE_SOURCE
	BlendModeOver       BlendMode = C.VIPS_BLEND_MODE_OVER
	BlendModeIn         BlendMode = C.VIPS_BLEND_MODE_IN
	BlendModeOut        BlendMode = C.VIPS_BLEND_MODE_OUT
	BlendModeAtop       BlendMode = C.VIPS_BLEND_MODE_ATOP
	BlendModeDest       BlendMode = C.VIPS_BLEND_MODE_DEST
	BlendModeDestOver   BlendMode = C.VIPS_BLEND_MODE_DEST_OVER
	BlendModeDestIn     BlendMode = C.VIPS_BLEND_MODE_DEST_IN
	BlendModeDestOut    BlendMode = C.VIPS_BLEND_MODE_DEST_OUT
	BlendModeDestAtop   BlendMode = C.VIPS_BLEND_MODE_DEST_ATOP
	BlendModeXOR        BlendMode = C.VIPS_BLEND_MODE_XOR
	BlendModeAdd        BlendMode = C.VIPS_BLEND_MODE_ADD
	BlendModeSaturate   BlendMode = C.VIPS_BLEND_MODE_SATURATE
	BlendModeMultiply   BlendMode = C.VIPS_BLEND_MODE_MULTIPLY
	BlendModeScreen     BlendMode = C.VIPS_BLEND_MODE_SCREEN
	BlendModeOverlay    BlendMode = C.VIPS_BLEND_MODE_OVERLAY
	BlendModeDarken     BlendMode = C.VIPS_BLEND_MODE_DARKEN
	BlendModeLighten    BlendMode = C.VIPS_BLEND_MODE_LIGHTEN
	BlendModeColorDodge BlendMode = C.VIPS_BLEND_MODE_COLOUR_DODGE
	BlendModeColorBurn  BlendMode = C.VIPS_BLEND_MODE_COLOUR_BURN
	BlendModeHardLight  BlendMode = C.VIPS_BLEND_MODE_HARD_LIGHT
	BlendModeSoftLight  BlendMode = C.VIPS_BLEND_MODE_SOFT_LIGHT
	BlendModeDifference BlendMode = C.VIPS_BLEND_MODE_DIFFERENCE
	BlendModeExclusion  BlendMode = C.VIPS_BLEND_MODE_EXCLUSION
)

// Gravity represents VIPS_GRAVITY type
type Gravity int

// Gravity enum
const (
	GravityCentre    Gravity = C.VIPS_COMPASS_DIRECTION_CENTRE
	GravityNorth     Gravity = C.VIPS_COMPASS_DIRECTION_NORTH
	GravityEast      Gravity = C.VIPS_COMPASS_DIRECTION_EAST
	GravitySouth     Gravity = C.VIPS_COMPASS_DIRECTION_SOUTH
	GravityWest      Gravity = C.VIPS_COMPASS_DIRECTION_WEST
	GravityNorthEast Gravity = C.VIPS_COMPASS_DIRECTION_NORTH_EAST
	GravityNorthWest Gravity = C.VIPS_COMPASS_DIRECTION_NORTH_WEST
	GravitySouthEast Gravity = C.VIPS_COMPASS_DIRECTION_SOUTH_EAST
	GravitySouthWest Gravity = C.VIPS_COMPASS_DIRECTION_SOUTH_WEST
)

// Direction represents VIPS_DIRECTION type
type Direction int

// Direction enum
const (
	DirectionHorizontal Direction = C.VIPS_DIRECTION_HORIZONTAL
	DirectionVertical   Direction = C.VIPS_DIRECTION_VERTICAL
)

// Angle represents VIPS_ANGLE type
type Angle int

// Angle enum
const (
	Angle0   Angle = C.VIPS_ANGLE_D0
	Angle90  Angle = C.VIPS_ANGLE_D90
	Angle180 Angle = C.VIPS_ANGLE_D180
	Angle270 Angle = C.VIPS_ANGLE_D270
)

// Angle45 represents VIPS_ANGLE45 type
type Angle45 int

// Angle45 enum
const (
	Angle45_0   Angle45 = C.VIPS_ANGLE45_D0
	Angle45_45  Angle45 = C.VIPS_ANGLE45_D45
	Angle45_90  Angle45 = C.VIPS_ANGLE45_D90
	Angle45_135 Angle45 = C.VIPS_ANGLE45_D135
	Angle45_180 Angle45 = C.VIPS_ANGLE45_D180
	Angle45_225 Angle45 = C.VIPS_ANGLE45_D225
	Angle45_270 Angle45 = C.VIPS_ANGLE45_D270
	Angle45_315 Angle45 = C.VIPS_ANGLE45_D315
)

// ExtendStrategy represents VIPS_EXTEND type
type ExtendStrategy int

// ExtendStrategy enum
const (
	ExtendBlack      ExtendStrategy = C.VIPS_EXTEND_BLACK
	ExtendCopy       ExtendStrategy = C.VIPS_EXTEND_COPY
	ExtendRepeat     ExtendStrategy = C.VIPS_EXTEND_REPEAT
	ExtendMirror     ExtendStrategy = C.VIPS_EXTEND_MIRROR
	ExtendWhite      ExtendStrategy = C.VIPS_EXTEND_WHITE
	ExtendBackground ExtendStrategy = C.VIPS_EXTEND_BACKGROUND
)

// Interesting represents VIPS_INTERESTING type
// https://libvips.github.io/libvips/API/current/libvips-conversion.html#VipsInteresting
type Interesting int

// Interesting constants represent areas of interest which smart cropping will crop based on.
const (
	InterestingNone      Interesting = C.VIPS_INTERESTING_NONE
	InterestingCentre    Interesting = C.VIPS_INTERESTING_CENTRE
	InterestingEntropy   Interesting = C.VIPS_INTERESTING_ENTROPY
	InterestingAttention Interesting = C.VIPS_INTERESTING_ATTENTION
	InterestingLow       Interesting = C.VIPS_INTERESTING_LOW
	InterestingHigh      Interesting = C.VIPS_INTERESTING_HIGH
	InterestingAll       Interesting = C.VIPS_INTERESTING_ALL
	InterestingLast      Interesting = C.VIPS_INTERESTING_LAST
)

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-embed
func vipsEmbed(in *C.VipsImage, left, top, width, height int, extend ExtendStrategy) (*C.VipsImage, error) {
	incOpCounter("embed")
	var out *C.VipsImage

	if err := C.embed_image(in, &out, C.int(left), C.int(top), C.int(width), C.int(height), C.int(extend)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-embed
func vipsEmbedBackground(in *C.VipsImage, left, top, width, height int, backgroundColor *ColorRGBA) (*C.VipsImage, error) {
	incOpCounter("embed")
	var out *C.VipsImage

	if err := C.embed_image_background(in, &out, C.int(left), C.int(top), C.int(width),
		C.int(height), C.double(backgroundColor.R),
		C.double(backgroundColor.G), C.double(backgroundColor.B), C.double(backgroundColor.A)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

func vipsEmbedMultiPage(in *C.VipsImage, left, top, width, height int, extend ExtendStrategy) (*C.VipsImage, error) {
	incOpCounter("embedMultiPage")
	var out *C.VipsImage

	if err := C.embed_multi_page_image(in, &out, C.int(left), C.int(top), C.int(width), C.int(height), C.int(extend)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

func vipsEmbedMultiPageBackground(in *C.VipsImage, left, top, width, height int, backgroundColor *ColorRGBA) (*C.VipsImage, error) {
	incOpCounter("embedMultiPageBackground")
	var out *C.VipsImage

	if err := C.embed_multi_page_image_background(in, &out, C.int(left), C.int(top), C.int(width),
		C.int(height), C.double(backgroundColor.R),
		C.double(backgroundColor.G), C.double(backgroundColor.B), C.double(backgroundColor.A)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-flip
func vipsFlip(in *C.VipsImage, direction Direction) (*C.VipsImage, error) {
	return vipsGenFlip(in, direction)
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-extract-area
func vipsExtractArea(in *C.VipsImage, left, top, width, height int) (*C.VipsImage, error) {
	return vipsGenExtractArea(in, left, top, width, height)
}

func vipsExtractAreaMultiPage(in *C.VipsImage, left, top, width, height int) (*C.VipsImage, error) {
	incOpCounter("extractAreaMultiPage")

	pageHeight := vipsGetPageHeight(in)
	nPages := int(in.Ysize) / pageHeight

	pages := make([]*C.VipsImage, nPages)
	for i := 0; i < nPages; i++ {
		page, err := vipsGenExtractArea(in, left, pageHeight*i+top, width, height)
		if err != nil {
			for j := 0; j < i; j++ {
				clearImage(pages[j])
			}
			return nil, err
		}
		pages[i] = page
	}

	across := 1
	joined, err := vipsGenArrayjoin(pages, &ArrayjoinOptions{Across: &across})
	for _, p := range pages {
		clearImage(p)
	}
	if err != nil {
		return nil, err
	}

	out, err := vipsGenCopy(joined, nil)
	clearImage(joined)
	if err != nil {
		return nil, err
	}

	vipsSetPageHeight(out, height)
	return out, nil
}

// http://libvips.github.io/libvips/API/current/libvips-resample.html#vips-similarity
func vipsSimilarity(in *C.VipsImage, scale float64, angle float64, color *ColorRGBA,
	idx float64, idy float64, odx float64, ody float64) (*C.VipsImage, error) {
	incOpCounter("similarity")
	var out *C.VipsImage

	if err := C.similarity(in, &out, C.double(scale), C.double(angle),
		C.double(color.R), C.double(color.G), C.double(color.B), C.double(color.A),
		C.double(idx), C.double(idy), C.double(odx), C.double(ody)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// http://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-smartcrop
func vipsSmartCrop(in *C.VipsImage, width int, height int, interesting Interesting) (*C.VipsImage, error) {
	_, out, _, err := vipsGenSmartcrop(in, width, height, &SmartcropOptions{Interesting: &interesting})
	return out, err
}

// http://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-crop
func vipsCrop(in *C.VipsImage, left int, top int, width int, height int) (*C.VipsImage, error) {
	incOpCounter("crop")
	var out *C.VipsImage

	if err := C.crop(in, &out, C.int(left), C.int(top), C.int(width), C.int(height)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-composite
func vipsComposite(ins []*C.VipsImage, modes []C.int, xs, ys []C.int) (*C.VipsImage, error) {
	if len(ins) == 0 || len(modes) == 0 || len(xs) == 0 || len(ys) == 0 {
		return nil, errors.New("vipsComposite: empty input slice")
	}
	incOpCounter("composite_multi")
	var out *C.VipsImage

	if err := C.composite_image(&ins[0], &out, C.int(len(ins)), &modes[0], &xs[0], &ys[0]); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-join
func vipsJoin(input1 *C.VipsImage, input2 *C.VipsImage, dir Direction) (*C.VipsImage, error) {
	incOpCounter("join")
	var out *C.VipsImage

	defer C.g_object_unref(C.gpointer(input1))
	defer C.g_object_unref(C.gpointer(input2))
	if err := C.join(input1, input2, &out, C.int(dir)); err != 0 {
		return nil, handleVipsError()
	}
	return out, nil
}

func vipsAddAlpha(in *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("addalpha")
	var out *C.VipsImage

	if err := C.add_alpha(in, &out); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

