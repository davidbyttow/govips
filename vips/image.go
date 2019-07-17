package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"errors"
	"io"
	"io/ioutil"
	"runtime"
	"unsafe"
)

// ImageRef contains a libvips image and manages its lifecycle. You should
// close an image when done or it will leak until the next GC
type ImageRef struct {
	image  *C.VipsImage
	format ImageType

	// NOTE(d): We keep a reference to this so that the input buffer is
	// never garbage collected during processing. Some image loaders use random
	// access transcoding and therefore need the original buffer to be in memory.
	buf []byte
}

type LoadOption func(o *vipsLoadOptions)

func WithAccessMode(a Access) LoadOption {
	return func(o *vipsLoadOptions) {
		switch a {
		case AccessRandom:
			o.cOpts.access = C.VIPS_ACCESS_RANDOM
		case AccessSequential:
			o.cOpts.access = C.VIPS_ACCESS_SEQUENTIAL
		case AccessSequentialUnbuffered:
			o.cOpts.access = C.VIPS_ACCESS_SEQUENTIAL_UNBUFFERED
		default:
			o.cOpts.access = C.VIPS_ACCESS_RANDOM
		}
	}
}

// LoadImage loads an ImageRef from the given reader
func LoadImage(r io.Reader, opts ...LoadOption) (*ImageRef, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewImageFromBuffer(buf, opts...)
}

// NewImageFromFile loads an image from file and creates a new ImageRef
func NewImageFromFile(file string, opts ...LoadOption) (*ImageRef, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return NewImageFromBuffer(buf, opts...)
}

// NewImageFromBuffer loads an image buffer and creates a new Image
func NewImageFromBuffer(buf []byte, opts ...LoadOption) (*ImageRef, error) {
	startupIfNeeded()

	image, format, err := vipsLoadFromBuffer(buf, opts...)
	if err != nil {
		return nil, err
	}

	ref := newImageRef(image, format, buf)

	return ref, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-copy
// create a new ref
func (r *ImageRef) Copy(options ...*Option) (*ImageRef, error) {
	out, err := Copy(r.image, options...)
	if err != nil {
		return nil, err
	}

	var buf []byte
	if r.buf != nil {
		buf = make([]byte, len(r.buf))
		copy(buf, r.buf)
	}

	return newImageRef(out, r.format, buf), nil
}

func newImageRef(vipsImage *C.VipsImage, format ImageType, buf []byte) *ImageRef {
	image := &ImageRef{
		image:  vipsImage,
		format: format,
		buf:    buf,
	}
	runtime.SetFinalizer(image, finalizeImage)

	return image
}

func finalizeImage(ref *ImageRef) {
	ref.Close()
}

// Close closes an image and frees internal memory associated with it
func (r *ImageRef) Close() {
	r.setImage(nil)
	r.buf = nil
}

// Format returns the initial format of the vips image when loaded
func (r *ImageRef) Format() ImageType {
	return r.format
}

// Width returns the width of this image
func (r *ImageRef) Width() int {
	return int(r.image.Xsize)
}

// Height returns the height of this iamge
func (r *ImageRef) Height() int {
	return int(r.image.Ysize)
}

// Bands returns the number of bands for this image
func (r *ImageRef) Bands() int {
	return int(r.image.Bands)
}

// HasProfile returns if the image has an ICC profile embedded.
func (r *ImageRef) HasProfile() bool {
	return vipsHasProfile(r.image)
}

// HasAlpha returns if the image has an alpha layer.
func (r *ImageRef) HasAlpha() bool {
	return vipsHasAlpha(r.image)
}

// ResX returns the X resolution
func (r *ImageRef) ResX() float64 {
	return float64(r.image.Xres)
}

// ResY returns the Y resolution
func (r *ImageRef) ResY() float64 {
	return float64(r.image.Yres)
}

// OffsetX returns the X offset
func (r *ImageRef) OffsetX() int {
	return int(r.image.Xoffset)
}

// OffsetY returns the Y offset
func (r *ImageRef) OffsetY() int {
	return int(r.image.Yoffset)
}

// BandFormat returns the current band format
func (r *ImageRef) BandFormat() BandFormat {
	return BandFormat(int(r.image.BandFmt))
}

// Coding returns the image coding
func (r *ImageRef) Coding() Coding {
	return Coding(int(r.image.Coding))
}

// Interpretation returns the current interpretation
func (r *ImageRef) Interpretation() Interpretation {
	return Interpretation(int(r.image.Type))
}

// Avg executes the 'avg' operation
func (r *ImageRef) Avg(options ...*Option) (float64, error) {
	return Avg(r.image, options...)
}

// Export exports the image
func (r *ImageRef) Export(params ExportParams) ([]byte, ImageType, error) {
	if params.Format == ImageTypeUnknown {
		params.Format = r.format
	}
	return vipsExportBuffer(r.image, &params)
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-composite
func (r *ImageRef) Composite(overlay *ImageRef, mode BlendMode) error {
	out, err := vipsComposite([]*C.VipsImage{r.image, overlay.image}, mode)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// ExtractBand executes the 'extract_band' operation
// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-extract-band
func (r *ImageRef) ExtractBand(band int, options ...*Option) error {
	out, err := ExtractBand(r.image, band, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

func (r *ImageRef) BandJoin(images ...*ImageRef) error {
	vipsImages := []*C.VipsImage{r.image}
	for _, image := range images {
		vipsImages = append(vipsImages, image.image)
	}

	out, err := vipsBandJoin(vipsImages)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

func (r *ImageRef) AddAlpha() error {
	if vipsHasAlpha(r.image) {
		return nil
	}

	out, err := vipsAddAlpha(r.image)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

//  https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-linear1
func (r *ImageRef) Linear1(a, b float64) error {
	out, err := vipsLinear1(r.image, a, b)
	if err != nil {
		return err
	}

	r.setImage(out)
	return nil
}

// Abs executes the 'abs' operation
func (r *ImageRef) Abs(options ...*Option) error {
	out, err := Abs(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Autorot executes the 'autorot' operation
func (r *ImageRef) Autorot(options ...*Option) error {
	out, err := Autorot(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Bandbool executes the 'bandbool' operation
func (r *ImageRef) Bandbool(boolean OperationBoolean, options ...*Option) error {
	out, err := Bandbool(r.image, boolean, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Bandfold executes the 'bandfold' operation
func (r *ImageRef) Bandfold(options ...*Option) error {
	out, err := Bandfold(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Bandunfold executes the 'bandunfold' operation
func (r *ImageRef) Bandunfold(options ...*Option) error {
	out, err := Bandunfold(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Buildlut executes the 'buildlut' operation
func (r *ImageRef) Buildlut(options ...*Option) error {
	out, err := Buildlut(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Byteswap executes the 'byteswap' operation
func (r *ImageRef) Byteswap(options ...*Option) error {
	out, err := Byteswap(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Cache executes the 'cache' operation
func (r *ImageRef) Cache(options ...*Option) error {
	out, err := Cache(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Cast executes the 'cast' operation
func (r *ImageRef) Cast(format BandFormat, options ...*Option) error {
	out, err := Cast(r.image, format, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Cmc2Lch executes the 'CMC2LCh' operation
func (r *ImageRef) Cmc2Lch(options ...*Option) error {
	out, err := Cmc2Lch(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Colourspace executes the 'colourspace' operation
func (r *ImageRef) Colourspace(space Interpretation, options ...*Option) error {
	out, err := Colourspace(r.image, space, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Complex executes the 'complex' operation
func (r *ImageRef) Complex(cmplx OperationComplex, options ...*Option) error {
	out, err := Complex(r.image, cmplx, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Complexget executes the 'complexget' operation
func (r *ImageRef) Complexget(get OperationComplexGet, options ...*Option) error {
	out, err := Complexget(r.image, get, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Embed executes the 'embed' operation
func (r *ImageRef) Embed(x int, y int, width int, height int, options ...*Option) error {
	out, err := Embed(r.image, x, y, width, height, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// ExtractArea executes the 'extract_area' operation
func (r *ImageRef) ExtractArea(left int, top int, width int, height int, options ...*Option) error {
	out, err := ExtractArea(r.image, left, top, width, height, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Falsecolour executes the 'falsecolour' operation
func (r *ImageRef) Falsecolour(options ...*Option) error {
	out, err := Falsecolour(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// FillNearest executes the 'fill_nearest' operation
func (r *ImageRef) FillNearest(options ...*Option) error {
	out, err := FillNearest(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Flatten executes the 'flatten' operation
func (r *ImageRef) Flatten(options ...*Option) error {
	out, err := Flatten(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Flip executes the 'flip' operation
func (r *ImageRef) Flip(direction Direction, options ...*Option) error {
	out, err := Flip(r.image, direction, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Float2Rad executes the 'float2rad' operation
func (r *ImageRef) Float2Rad(options ...*Option) error {
	out, err := Float2Rad(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Fwfft executes the 'fwfft' operation
func (r *ImageRef) Fwfft(options ...*Option) error {
	out, err := Fwfft(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Gamma executes the 'gamma' operation
func (r *ImageRef) Gamma(options ...*Option) error {
	out, err := Gamma(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Gaussblur executes the 'gaussblur' operation
func (r *ImageRef) Gaussblur(sigma float64, options ...*Option) error {
	out, err := Gaussblur(r.image, sigma, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Globalbalance executes the 'globalbalance' operation
func (r *ImageRef) Globalbalance(options ...*Option) error {
	out, err := Globalbalance(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Grid executes the 'grid' operation
func (r *ImageRef) Grid(tileHeight int, across int, down int, options ...*Option) error {
	out, err := Grid(r.image, tileHeight, across, down, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HistCum executes the 'hist_cum' operation
func (r *ImageRef) HistCum(options ...*Option) error {
	out, err := HistCum(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HistEqual executes the 'hist_equal' operation
func (r *ImageRef) HistEqual(options ...*Option) error {
	out, err := HistEqual(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HistFind executes the 'hist_find' operation
func (r *ImageRef) HistFind(options ...*Option) error {
	out, err := HistFind(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HistFindNdim executes the 'hist_find_ndim' operation
func (r *ImageRef) HistFindNdim(options ...*Option) error {
	out, err := HistFindNdim(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HistLocal executes the 'hist_local' operation
func (r *ImageRef) HistLocal(width int, height int, options ...*Option) error {
	out, err := HistLocal(r.image, width, height, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HistNorm executes the 'hist_norm' operation
func (r *ImageRef) HistNorm(options ...*Option) error {
	out, err := HistNorm(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HistPlot executes the 'hist_plot' operation
func (r *ImageRef) HistPlot(options ...*Option) error {
	out, err := HistPlot(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HoughCircle executes the 'hough_circle' operation
func (r *ImageRef) HoughCircle(options ...*Option) error {
	out, err := HoughCircle(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// HoughLine executes the 'hough_line' operation
func (r *ImageRef) HoughLine(options ...*Option) error {
	out, err := HoughLine(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Hsv2Srgb executes the 'HSV2sRGB' operation
func (r *ImageRef) Hsv2Srgb(options ...*Option) error {
	out, err := Hsv2Srgb(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// IccExport executes the 'icc_export' operation
func (r *ImageRef) IccExport(options ...*Option) error {
	out, err := IccExport(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// IccImport executes the 'icc_import' operation
func (r *ImageRef) IccImport(options ...*Option) error {
	out, err := IccImport(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// IccTransform executes the 'icc_transform' operation
func (r *ImageRef) IccTransform(outputProfile string, options ...*Option) error {
	out, err := IccTransform(r.image, outputProfile, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Invert executes the 'invert' operation
func (r *ImageRef) Invert(options ...*Option) error {
	out, err := Invert(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Invertlut executes the 'invertlut' operation
func (r *ImageRef) Invertlut(options ...*Option) error {
	out, err := Invertlut(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Invfft executes the 'invfft' operation
func (r *ImageRef) Invfft(options ...*Option) error {
	out, err := Invfft(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Lab2Labq executes the 'Lab2LabQ' operation
func (r *ImageRef) Lab2Labq(options ...*Option) error {
	out, err := Lab2Labq(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Lab2Labs executes the 'Lab2LabS' operation
func (r *ImageRef) Lab2Labs(options ...*Option) error {
	out, err := Lab2Labs(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Lab2Lch executes the 'Lab2LCh' operation
func (r *ImageRef) Lab2Lch(options ...*Option) error {
	out, err := Lab2Lch(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Lab2Xyz executes the 'Lab2XYZ' operation
func (r *ImageRef) Lab2Xyz(options ...*Option) error {
	out, err := Lab2Xyz(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Labelregions executes the 'labelregions' operation
func (r *ImageRef) Labelregions(options ...*Option) error {
	out, err := Labelregions(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Labq2Lab executes the 'LabQ2Lab' operation
func (r *ImageRef) Labq2Lab(options ...*Option) error {
	out, err := Labq2Lab(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Labq2Labs executes the 'LabQ2LabS' operation
func (r *ImageRef) Labq2Labs(options ...*Option) error {
	out, err := Labq2Labs(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Labq2Srgb executes the 'LabQ2sRGB' operation
func (r *ImageRef) Labq2Srgb(options ...*Option) error {
	out, err := Labq2Srgb(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Labs2Lab executes the 'LabS2Lab' operation
func (r *ImageRef) Labs2Lab(options ...*Option) error {
	out, err := Labs2Lab(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Labs2Labq executes the 'LabS2LabQ' operation
func (r *ImageRef) Labs2Labq(options ...*Option) error {
	out, err := Labs2Labq(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Linecache executes the 'linecache' operation
func (r *ImageRef) Linecache(options ...*Option) error {
	out, err := Linecache(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Math executes the 'math' operation
func (r *ImageRef) Math(math OperationMath, options ...*Option) error {
	out, err := Math(r.image, math, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Measure executes the 'measure' operation
func (r *ImageRef) Measure(h int, v int, options ...*Option) error {
	out, err := Measure(r.image, h, v, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Msb executes the 'msb' operation
func (r *ImageRef) Msb(options ...*Option) error {
	out, err := Msb(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Premultiply executes the 'premultiply' operation
func (r *ImageRef) Premultiply(options ...*Option) error {
	out, err := Premultiply(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Rad2Float executes the 'rad2float' operation
func (r *ImageRef) Rad2Float(options ...*Option) error {
	out, err := Rad2Float(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Rank executes the 'rank' operation
func (r *ImageRef) Rank(width int, height int, index int, options ...*Option) error {
	out, err := Rank(r.image, width, height, index, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Reduce executes the 'reduce' operation
func (r *ImageRef) Reduce(hshrink float64, vshrink float64, options ...*Option) error {
	out, err := Reduce(r.image, hshrink, vshrink, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Reduceh executes the 'reduceh' operation
func (r *ImageRef) Reduceh(hshrink float64, options ...*Option) error {
	out, err := Reduceh(r.image, hshrink, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Reducev executes the 'reducev' operation
func (r *ImageRef) Reducev(vshrink float64, options ...*Option) error {
	out, err := Reducev(r.image, vshrink, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Replicate executes the 'replicate' operation
func (r *ImageRef) Replicate(across int, down int, options ...*Option) error {
	out, err := Replicate(r.image, across, down, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Resize executes the 'resize' operation
func (r *ImageRef) Resize(scale float64, options ...*Option) error {
	out, err := Resize(r.image, scale, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Rot executes the 'rot' operation
func (r *ImageRef) Rot(angle Angle, options ...*Option) error {
	out, err := Rot(r.image, angle, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Rot45 executes the 'rot45' operation
func (r *ImageRef) Rot45(options ...*Option) error {
	out, err := Rot45(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Round executes the 'round' operation
func (r *ImageRef) Round(round OperationRound, options ...*Option) error {
	out, err := Round(r.image, round, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Scale executes the 'scale' operation
func (r *ImageRef) Scale(options ...*Option) error {
	out, err := Scale(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Scrgb2Bw executes the 'scRGB2BW' operation
func (r *ImageRef) Scrgb2Bw(options ...*Option) error {
	out, err := Scrgb2Bw(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Scrgb2Srgb executes the 'scRGB2sRGB' operation
func (r *ImageRef) Scrgb2Srgb(options ...*Option) error {
	out, err := Scrgb2Srgb(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Scrgb2Xyz executes the 'scRGB2XYZ' operation
func (r *ImageRef) Scrgb2Xyz(options ...*Option) error {
	out, err := Scrgb2Xyz(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Sequential executes the 'sequential' operation
func (r *ImageRef) Sequential(options ...*Option) error {
	out, err := Sequential(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Sharpen executes the 'sharpen' operation
func (r *ImageRef) Sharpen(options ...*Option) error {
	out, err := Sharpen(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Shrink executes the 'shrink' operation
func (r *ImageRef) Shrink(hshrink float64, vshrink float64, options ...*Option) error {
	out, err := Shrink(r.image, hshrink, vshrink, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Shrinkh executes the 'shrinkh' operation
func (r *ImageRef) Shrinkh(hshrink int, options ...*Option) error {
	out, err := Shrinkh(r.image, hshrink, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Shrinkv executes the 'shrinkv' operation
func (r *ImageRef) Shrinkv(vshrink int, options ...*Option) error {
	out, err := Shrinkv(r.image, vshrink, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Sign executes the 'sign' operation
func (r *ImageRef) Sign(options ...*Option) error {
	out, err := Sign(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Similarity executes the 'similarity' operation
func (r *ImageRef) Similarity(options ...*Option) error {
	out, err := Similarity(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Smartcrop executes the 'smartcrop' operation
func (r *ImageRef) Smartcrop(width int, height int, options ...*Option) error {
	out, err := Smartcrop(r.image, width, height, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Spectrum executes the 'spectrum' operation
func (r *ImageRef) Spectrum(options ...*Option) error {
	out, err := Spectrum(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Srgb2Hsv executes the 'sRGB2HSV' operation
func (r *ImageRef) Srgb2Hsv(options ...*Option) error {
	out, err := Srgb2Hsv(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Srgb2Scrgb executes the 'sRGB2scRGB' operation
func (r *ImageRef) Srgb2Scrgb(options ...*Option) error {
	out, err := Srgb2Scrgb(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Stats executes the 'stats' operation
func (r *ImageRef) Stats(options ...*Option) error {
	out, err := Stats(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Stdif executes the 'stdif' operation
func (r *ImageRef) Stdif(width int, height int, options ...*Option) error {
	out, err := Stdif(r.image, width, height, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Subsample executes the 'subsample' operation
func (r *ImageRef) Subsample(xfac int, yfac int, options ...*Option) error {
	out, err := Subsample(r.image, xfac, yfac, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// ThumbnailImage executes the 'thumbnail_image' operation
func (r *ImageRef) ThumbnailImage(width int, options ...*Option) error {
	out, err := ThumbnailImage(r.image, width, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Tilecache executes the 'tilecache' operation
func (r *ImageRef) Tilecache(options ...*Option) error {
	out, err := Tilecache(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Unpremultiply executes the 'unpremultiply' operation
func (r *ImageRef) Unpremultiply(options ...*Option) error {
	out, err := Unpremultiply(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Wrap executes the 'wrap' operation
func (r *ImageRef) Wrap(options ...*Option) error {
	out, err := Wrap(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Xyz2Lab executes the 'XYZ2Lab' operation
func (r *ImageRef) Xyz2Lab(options ...*Option) error {
	out, err := Xyz2Lab(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Xyz2Scrgb executes the 'XYZ2scRGB' operation
func (r *ImageRef) Xyz2Scrgb(options ...*Option) error {
	out, err := Xyz2Scrgb(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Xyz2Yxy executes the 'XYZ2Yxy' operation
func (r *ImageRef) Xyz2Yxy(options ...*Option) error {
	out, err := Xyz2Yxy(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Yxy2Xyz executes the 'Yxy2XYZ' operation
func (r *ImageRef) Yxy2Xyz(options ...*Option) error {
	out, err := Yxy2Xyz(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Zoom executes the 'zoom' operation
func (r *ImageRef) Zoom(xfac int, yfac int, options ...*Option) error {
	out, err := Zoom(r.image, xfac, yfac, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Bandmean executes the 'bandmean' operation
func (r *ImageRef) Bandmean(options ...*Option) error {
	out, err := Bandmean(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Lch2Cmc executes the 'LCh2CMC' operation
func (r *ImageRef) Lch2Cmc(options ...*Option) error {
	out, err := Lch2Cmc(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Lch2Lab executes the 'LCh2Lab' operation
func (r *ImageRef) Lch2Lab(options ...*Option) error {
	out, err := Lch2Lab(r.image, options...)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// ToBytes writes the image to memory in VIPs format and returns the raw bytes, useful for storage.
func (r *ImageRef) ToBytes() ([]byte, error) {
	var cSize C.size_t
	cData := C.vips_image_write_to_memory(r.image, &cSize)
	if cData == nil {
		return nil, errors.New("failed to write image to memory")
	}
	defer C.free(cData)

	bytes := C.GoBytes(unsafe.Pointer(cData), C.int(cSize))
	return bytes, nil
}

// setImage resets the image for this image and frees the previous one
func (r *ImageRef) setImage(image *C.VipsImage) {
	if r.image != nil {
		defer unrefImage(r.image)
	}

	r.image = image
}
