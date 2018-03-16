package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"io"
	"log"
	"strings"
	"sync"
)

// ResizeStrategy is the strategy to use when resizing an image
type ResizeStrategy int

// ResizeStrategy enum
const (
	ResizeStrategyAuto ResizeStrategy = iota
	ResizeStrategyEmbed
	ResizeStrategyCrop
	ResizeStrategyStretch
)

// ExportParams are options when exporting an image to file or buffer
type ExportParams struct {
	OutputFile      string
	Writer          io.Writer
	Format          ImageType
	Quality         int
	Compression     int
	Interlaced      bool
	Lossless        bool
	StripProfile    bool
	StripMetadata   bool
	Interpretation  Interpretation
	BackgroundColor *Color
}

// Color represents an RGB
type Color struct {
	R, G, B uint8
}

// ColorBlack is shorthand for black RGB
var ColorBlack = Color{0, 0, 0}

// Anchor represents the an anchor for cropping and other image operations
type Anchor int

// Anchor enum
const (
	AnchorAuto Anchor = iota
	AnchorCenter
	AnchorTop
	AnchorTopRight
	AnchorTopLeft
	AnchorRight
	AnchorBottom
	AnchorBottomLeft
	AnchorBottomRight
	AnchorLeft
)

// FlipDirection represents the direction to flip
type FlipDirection int

// Flip enum
const (
	FlipNone FlipDirection = iota
	FlipHorizontal
	FlipVertical
	FlipBoth
)

// ImageType represents an image type
type ImageType int

// ImageType enum
const (
	ImageTypeUnknown ImageType = C.UNKNOWN
	ImageTypeGIF     ImageType = C.GIF
	ImageTypeJPEG    ImageType = C.JPEG
	ImageTypeMagick  ImageType = C.MAGICK
	ImageTypePDF     ImageType = C.PDF
	ImageTypePNG     ImageType = C.PNG
	ImageTypeSVG     ImageType = C.SVG
	ImageTypeTIFF    ImageType = C.TIFF
	ImageTypeWEBP    ImageType = C.WEBP
)

var imageTypeExtensionMap = map[ImageType]string{
	ImageTypeGIF:    ".gif",
	ImageTypeJPEG:   ".jpeg",
	ImageTypeMagick: ".magick",
	ImageTypePDF:    ".pdf",
	ImageTypePNG:    ".png",
	ImageTypeSVG:    ".svg",
	ImageTypeTIFF:   ".tiff",
	ImageTypeWEBP:   ".webp",
}

// OutputExt returns the canonical extension for the ImageType
func (i ImageType) OutputExt() string {
	if ext, ok := imageTypeExtensionMap[i]; ok {
		return ext
	}
	return ""
}

// Kernel represents VipsKernel type
type Kernel int

// Kernel enum
const (
	KernelAuto     Kernel = -1
	KernelNearest  Kernel = C.VIPS_KERNEL_NEAREST
	KernelLinear   Kernel = C.VIPS_KERNEL_LINEAR
	KernelCubic    Kernel = C.VIPS_KERNEL_CUBIC
	KernelLanczos2 Kernel = C.VIPS_KERNEL_LANCZOS2
	KernelLanczos3 Kernel = C.VIPS_KERNEL_LANCZOS3
)

// Interpolator represents the vips interpolator types
type Interpolator string

// Interpolator enum
const (
	InterpolateBicubic  Interpolator = "bicubic"
	InterpolateBilinear Interpolator = "bilinear"
	InterpolateNoHalo   Interpolator = "nohalo"
)

// String returns the canonical name of the interpolator
func (i Interpolator) String() string {
	return string(i)
}

// OperationMath represents VIPS_OPERATION_MATH type
type OperationMath int

// OperationMath enum
const (
	OperationMathSin   OperationMath = C.VIPS_OPERATION_MATH_SIN
	OperationMathCos   OperationMath = C.VIPS_OPERATION_MATH_COS
	OperationMathTan   OperationMath = C.VIPS_OPERATION_MATH_TAN
	OperationMathAsin  OperationMath = C.VIPS_OPERATION_MATH_ASIN
	OperationMathAcos  OperationMath = C.VIPS_OPERATION_MATH_ACOS
	OperationMathAtan  OperationMath = C.VIPS_OPERATION_MATH_ATAN
	OperationMathLog   OperationMath = C.VIPS_OPERATION_MATH_LOG
	OperationMathLog10 OperationMath = C.VIPS_OPERATION_MATH_LOG10
	OperationMathExp   OperationMath = C.VIPS_OPERATION_MATH_EXP
	OperationMathExp10 OperationMath = C.VIPS_OPERATION_MATH_EXP10
)

// OperationMath2 represents VIPS_OPERATION_MATH2 type
type OperationMath2 int

// OperationMath2 enum
const (
	OperationMath2Pow OperationMath2 = C.VIPS_OPERATION_MATH2_POW
	OperationMath2Wop OperationMath2 = C.VIPS_OPERATION_MATH2_WOP
)

// OperationRound represents VIPS_OPERATION_ROUND type
type OperationRound int

// OperationRound enum
const (
	OperationRoundRint  OperationRound = C.VIPS_OPERATION_ROUND_RINT
	OperationRoundCeil  OperationRound = C.VIPS_OPERATION_ROUND_CEIL
	OperationRoundFloor OperationRound = C.VIPS_OPERATION_ROUND_FLOOR
)

// OperationRelational represents VIPS_OPERATION_RELATIONAL type
type OperationRelational int

// OperationRelational enum
const (
	OperationRelationalEqual  OperationRelational = C.VIPS_OPERATION_RELATIONAL_EQUAL
	OperationRelationalNotEq  OperationRelational = C.VIPS_OPERATION_RELATIONAL_NOTEQ
	OperationRelationalLess   OperationRelational = C.VIPS_OPERATION_RELATIONAL_LESS
	OperationRelationalLessEq OperationRelational = C.VIPS_OPERATION_RELATIONAL_LESSEQ
	OperationRelationalMore   OperationRelational = C.VIPS_OPERATION_RELATIONAL_MORE
	OperationRelationalMoreEq OperationRelational = C.VIPS_OPERATION_RELATIONAL_MOREEQ
)

// OperationBoolean represents VIPS_OPERATION_BOOLEAN type
type OperationBoolean int

// OperationBoolean enum
const (
	OperationBooleanAnd    OperationBoolean = C.VIPS_OPERATION_BOOLEAN_AND
	OperationBooleanOr     OperationBoolean = C.VIPS_OPERATION_BOOLEAN_OR
	OperationBooleanEOr    OperationBoolean = C.VIPS_OPERATION_BOOLEAN_EOR
	OperationBooleanLShift OperationBoolean = C.VIPS_OPERATION_BOOLEAN_LSHIFT
	OperationBooleanRShift OperationBoolean = C.VIPS_OPERATION_BOOLEAN_RSHIFT
)

// OperationComplex represents VIPS_OPERATION_COMPLEX type
type OperationComplex int

// OperationComplex enum
const (
	OperationComplexPolar OperationComplex = C.VIPS_OPERATION_COMPLEX_POLAR
	OperationComplexRect  OperationComplex = C.VIPS_OPERATION_COMPLEX_RECT
	OperationComplexConj  OperationComplex = C.VIPS_OPERATION_COMPLEX_CONJ
)

// OperationComplex2 represents VIPS_OPERATION_COMPLEX2 type
type OperationComplex2 int

// OperationComplex2 enum
const (
	OperationComplex2CrossPhase OperationComplex2 = C.VIPS_OPERATION_COMPLEX2_CROSS_PHASE
)

// OperationComplexGet represents VIPS_OPERATION_COMPLEXGET type
type OperationComplexGet int

// OperationComplexGet enum
const (
	OperationComplexReal OperationComplexGet = C.VIPS_OPERATION_COMPLEXGET_REAL
	OperationComplexImag OperationComplexGet = C.VIPS_OPERATION_COMPLEXGET_IMAG
)

// Extend represents VIPS_EXTEND type
type Extend int

// Extend enum
const (
	ExtendBlack      Extend = C.VIPS_EXTEND_BLACK
	ExtendCopy       Extend = C.VIPS_EXTEND_COPY
	ExtendRepeat     Extend = C.VIPS_EXTEND_REPEAT
	ExtendMirror     Extend = C.VIPS_EXTEND_MIRROR
	ExtendWhite      Extend = C.VIPS_EXTEND_WHITE
	ExtendBackground Extend = C.VIPS_EXTEND_BACKGROUND
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

// Coding represents VIPS_CODING type
type Coding int

// Coding enum
const (
	CodingError Coding = C.VIPS_CODING_ERROR
	CodingNone  Coding = C.VIPS_CODING_NONE
	CodingLABQ  Coding = C.VIPS_CODING_LABQ
	CodingRAD   Coding = C.VIPS_CODING_RAD
)

// Access represents VIPS_ACCESS
type Access int

// Access enum
const (
	AccessRandom               Access = C.VIPS_ACCESS_RANDOM
	AccessSequential           Access = C.VIPS_ACCESS_SEQUENTIAL
	AccessSequentialUnbuffered Access = C.VIPS_ACCESS_SEQUENTIAL_UNBUFFERED
)

// OperationMorphology represents VIPS_OPERATION_MORPHOLOGY
type OperationMorphology int

// OperationMorphology enum
const (
	MorphologyErode  OperationMorphology = C.VIPS_OPERATION_MORPHOLOGY_ERODE
	MorphologyDilate OperationMorphology = C.VIPS_OPERATION_MORPHOLOGY_DILATE
)

var ImageTypes = map[ImageType]string{
	ImageTypeGIF:    "gif",
	ImageTypeJPEG:   "jpeg",
	ImageTypeMagick: "magick",
	ImageTypePDF:    "pdf",
	ImageTypePNG:    "png",
	ImageTypeSVG:    "svg",
	ImageTypeTIFF:   "tiff",
	ImageTypeWEBP:   "webp",
}

type Composite struct {
	Image     *ImageRef
	BlendMode BlendMode
}

type BlendMode int

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

var (
	once                sync.Once
	typeLoaders         = make(map[string]ImageType)
	supportedImageTypes = make(map[ImageType]bool)
)

// DetermineImageType attempts to determine the image type of the given buffer
func DetermineImageType(buf []byte) ImageType {
	return vipsDetermineImageType(buf)
}

func IsTypeSupported(imageType ImageType) bool {
	return supportedImageTypes[imageType]
}

// InitTypes initializes caches and figures out which image types are supported
func initTypes() {
	once.Do(func() {
		cType := C.CString("VipsOperation")
		defer freeCString(cType)

		for k, v := range ImageTypes {
			name := strings.ToLower("VipsForeignLoad" + v)
			typeLoaders[name] = k
			typeLoaders[name+"buffer"] = k

			cFunc := C.CString(v + "load")
			defer freeCString(cFunc)

			ret := C.vips_type_find(
				cType,
				cFunc)
			log.Printf("Registered image type loader type=%s", v)
			supportedImageTypes[k] = int(ret) != 0
		}
	})
}
