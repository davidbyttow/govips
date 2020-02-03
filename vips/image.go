package vips

// #cgo pkg-config: vips
// #include "image.h"
import "C"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"unsafe"
)

const (
	defaultQuality     = 80
	defaultCompression = 6
)

type PreMultiplicationState struct {
	bandFormat BandFormat
}

// ImageRef contains a libvips image and manages its lifecycle. You need to
// close an image when done or it will leak
type ImageRef struct {
	// NOTE: We keep a reference to this so that the input buffer is
	// never garbage collected during processing. Some image loaders use random
	// access transcoding and therefore need the original buffer to be in memory.
	buf               []byte
	image             *C.VipsImage
	format            ImageType
	lock              sync.Mutex
	preMultiplication *PreMultiplicationState
}

type ImageMetadata struct {
	Format      ImageType
	Width       int
	Height      int
	Colorspace  Interpretation
	Orientation int
}

// ExportParams are options when exporting an image to file or buffer
type ExportParams struct {
	Format      ImageType
	Quality     int
	Compression int
	Interlaced  bool
	Lossless    bool
}

// NewImageFromReader loads an ImageRef from the given reader
func NewImageFromReader(r io.Reader) (*ImageRef, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewImageFromBuffer(buf)
}

// NewImageFromFile loads an image from file and creates a new ImageRef
func NewImageFromFile(file string) (*ImageRef, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return NewImageFromBuffer(buf)
}

// NewImageFromBuffer loads an image buffer and creates a new Image
func NewImageFromBuffer(buf []byte) (*ImageRef, error) {
	startupIfNeeded()

	image, format, err := vipsLoadFromBuffer(buf)
	if err != nil {
		return nil, err
	}

	ref := newImageRef(image, format, buf)

	return ref, nil
}

func (r *ImageRef) Metadata() *ImageMetadata {
	return &ImageMetadata{
		Format: r.Format(),
		Width:  r.Width(),
		Height: r.Height(),
	}
}

// create a new ref
// deprecated
func (r *ImageRef) Copy() (*ImageRef, error) {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return nil, err
	}

	return newImageRef(out, r.format, r.buf), nil
}

func newImageRef(vipsImage *C.VipsImage, format ImageType, buf []byte) *ImageRef {
	image := &ImageRef{
		image:  vipsImage,
		format: format,
		buf:    buf,
	}

	return image
}

// Close closes an image and frees internal memory associated with it
func (r *ImageRef) Close() {
	r.lock.Lock()

	if r.image != nil {
		clearImage(r.image)
		r.image = nil
	}

	r.buf = nil

	r.lock.Unlock()
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
	return vipsHasICCProfile(r.image)
}

// alias to HasProfile()
func (r *ImageRef) HasICCProfile() bool {
	return r.HasProfile()
}

func (r *ImageRef) HasIPTC() bool {
	return vipsHasICPTC(r.image)
}

// HasAlpha returns if the image has an alpha layer.
func (r *ImageRef) HasAlpha() bool {
	return vipsHasAlpha(r.image)
}

// Return the orientation number as appears in the EXIF, if present
func (r *ImageRef) GetOrientation() int {
	return vipsGetMetaOrientation(r.image)
}

func (r *ImageRef) SetOrientation(orientation int) error {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return err
	}

	vipsSetMetaOrientation(out, orientation)

	r.setImage(out)
	return nil
}

func (r *ImageRef) RemoveOrientation() error {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return err
	}

	vipsRemoveMetaOrientation(out)

	r.setImage(out)
	return nil
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

// Alias to Interpretation()
func (r *ImageRef) ColorSpace() Interpretation {
	return r.Interpretation()
}

func (r *ImageRef) IsColorSpaceSupported() bool {
	return vipsIsColorSpaceSupported(r.image)
}

// Export exports the image
func (r *ImageRef) Export(params *ExportParams) ([]byte, *ImageMetadata, error) {
	p := params
	if p == nil {
		p = &ExportParams{}
	}

	if p.Format == ImageTypeUnknown {
		p.Format = r.format
	}

	// the exported buf is not necessarily in same format as the original buf, might default to JPEG as well.
	buf, format, err := r.exportBuffer(p)
	if err != nil {
		return nil, nil, err
	}

	metadata := &ImageMetadata{
		Format:      format,
		Width:       r.Width(),
		Height:      r.Height(),
		Colorspace:  r.ColorSpace(),
		Orientation: r.GetOrientation(),
	}

	return buf, metadata, nil
}

func (r *ImageRef) Composite(overlay *ImageRef, mode BlendMode, x, y int) error {
	out, err := vipsComposite2(r.image, overlay.image, mode, x, y)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// ExtractBand executes the 'extract_band' operation
func (r *ImageRef) ExtractBand(band int, num int) error {
	out, err := vipsExtractBand(r.image, band, num)
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

func (r *ImageRef) PremultiplyAlpha() error {
	if r.preMultiplication != nil || !vipsHasAlpha(r.image) {
		return nil
	}

	band := r.BandFormat()

	out, err := vipsPremultiplyAlpha(r.image)
	if err != nil {
		return err
	}
	r.preMultiplication = &PreMultiplicationState{
		bandFormat: band,
	}
	r.setImage(out)
	return nil
}

func (r *ImageRef) UnpremultiplyAlpha() error {
	if r.preMultiplication == nil {
		return nil
	}

	unpremultiplied, err := vipsUnpremultiplyAlpha(r.image)
	if err != nil {
		return err
	}
	defer clearImage(unpremultiplied)

	out, err := vipsCast(unpremultiplied, r.preMultiplication.bandFormat)
	if err != nil {
		return err
	}

	r.preMultiplication = nil
	r.setImage(out)
	return nil
}

func (r *ImageRef) Linear(a, b []float64) error {
	if len(a) != len(b) {
		return errors.New("a and b must be of same length")
	}

	out, err := vipsLinear(r.image, a, b, len(a))
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

func (r *ImageRef) Linear1(a, b float64) error {
	out, err := vipsLinear1(r.image, a, b)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

func getZeroedAngle(angle Angle) Angle {
	switch angle {
	case Angle0:
		return Angle0
	case Angle90:
		return Angle270
	case Angle180:
		return Angle180
	case Angle270:
		return Angle90
	}
	return Angle0
}

func GetRotationAngleFromExif(orientation int) (Angle, bool) {
	switch orientation {
	case 0, 1, 2:
		return Angle0, orientation == 2
	case 3, 4:
		return Angle180, orientation == 4
	case 5, 8:
		return Angle90, orientation == 5
	case 6, 7:
		return Angle270, orientation == 7
	}

	return Angle0, false
}

// Autorot do auto rotation
func (r *ImageRef) AutoRotate() error {
	// this is a full implementation of auto rotate as vips doesn't support auto rotating of mirrors exifs
	// https://jcupitt.github.io/libvips/API/current/libvips-conversion.html#vips-autorot
	angle, flipped := GetRotationAngleFromExif(r.GetOrientation())
	if flipped {
		err := r.Flip(DirectionHorizontal)
		if err != nil {
			return err
		}
	}

	zeroAngle := getZeroedAngle(angle)
	err := r.Rotate(zeroAngle)
	if err != nil {
		return err
	}

	return r.RemoveOrientation()
}

// ExtractArea executes the 'extract_area' operation
func (r *ImageRef) ExtractArea(left int, top int, width int, height int) error {
	out, err := vipsExtractArea(r.image, left, top, width, height)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

func (r *ImageRef) RemoveICCProfile() error {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return err
	}

	vipsRemoveICCProfile(out)

	r.setImage(out)
	return nil
}

// deprecated: use optimize
func (r *ImageRef) TransformICCProfile(isCmyk int) error {
	out, err := vipsOptimizeICCProfile(r.image, isCmyk)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

func (r *ImageRef) OptimizeICCProfile() error {
	isCMYK := 0
	if r.Interpretation() == InterpretationCMYK {
		isCMYK = 1
	}

	out, err := vipsOptimizeICCProfile(r.image, isCMYK)
	if err != nil {
		info(err.Error())
		return err
	}
	r.setImage(out)
	return nil
}

// won't remove the ICC profile
func (r *ImageRef) RemoveMetadata() error {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return err
	}

	vipsRemoveMetadata(out)

	r.setImage(out)
	return nil
}

func (r *ImageRef) ToColorSpace(interpretation Interpretation) error {
	out, err := vipsToColorSpace(r.image, interpretation)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Flatten executes the 'flatten' operation
func (r *ImageRef) Flatten(backgroundColor *Color) error {
	out, err := vipsFlatten(r.image, backgroundColor)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Gaussblur executes the 'gaussblur' operation
func (r *ImageRef) GaussianBlur(sigma float64) error {
	out, err := vipsGaussianBlur(r.image, sigma)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Sharpen executes the 'sharpen' operation
func (r *ImageRef) Sharpen(sigma float64, x1 float64, m2 float64) error {
	out, err := vipsSharpen(r.image, sigma, x1, m2)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Modulate the colors
func (r *ImageRef) Modulate(brightness, saturation float64, hue int) error {
	var err error
	var multiplications []float64
	var additions []float64

	colorspace := r.ColorSpace()
	if colorspace == InterpretationRGB {
		colorspace = InterpretationSRGB
	}

	if r.HasAlpha() {
		multiplications = []float64{brightness, saturation, 1, 1}
		additions = []float64{0, 0, float64(hue), 0}
	} else {
		multiplications = []float64{brightness, saturation, 1}
		additions = []float64{0, 0, float64(hue)}
	}

	err = r.ToColorSpace(InterpretationLCH)
	if err != nil {
		return err
	}

	err = r.Linear(multiplications, additions)
	if err != nil {
		return err
	}

	err = r.ToColorSpace(colorspace)
	if err != nil {
		return err
	}

	return nil
}

// Modulate the colors
func (r *ImageRef) ModulateHSV(brightness, saturation float64, hue int) error {
	var err error
	var multiplications []float64
	var additions []float64

	colorspace := r.ColorSpace()
	if colorspace == InterpretationRGB {
		colorspace = InterpretationSRGB
	}

	if r.HasAlpha() {
		multiplications = []float64{1, saturation, brightness, 1}
		additions = []float64{float64(hue), 0, 0, 0}
	} else {
		multiplications = []float64{1, saturation, brightness}
		additions = []float64{float64(hue), 0, 0}
	}

	err = r.ToColorSpace(InterpretationHSV)
	if err != nil {
		return err
	}

	err = r.Linear(multiplications, additions)
	if err != nil {
		return err
	}

	err = r.ToColorSpace(colorspace)
	if err != nil {
		return err
	}

	return nil
}

// Invert executes the 'invert' operation
func (r *ImageRef) Invert() error {
	out, err := vipsInvert(r.image)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Resize executes the 'resize' operation
func (r *ImageRef) Resize(scale float64, kernel Kernel) error {
	err := r.PremultiplyAlpha()
	if err != nil {
		return err
	}

	out, err := vipsResize(r.image, scale, kernel)
	if err != nil {
		return err
	}
	r.setImage(out)

	return r.UnpremultiplyAlpha()
}

// Resize executes the 'resize' operation
func (r *ImageRef) ResizeWithVScale(hScale, vScale float64, kernel Kernel) error {
	out, err := vipsResizeWithVScale(r.image, hScale, vScale, kernel)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Embed executes the 'embed' operation
func (r *ImageRef) Embed(left, top, width, height int, extend ExtendStrategy) error {
	out, err := vipsEmbed(r.image, left, top, width, height, extend)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Zoom executes the 'zoom' operation
func (r *ImageRef) Zoom(xFactor int, yFactor int) error {
	out, err := vipsZoom(r.image, xFactor, yFactor)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Flip executes the 'flip' operation
func (r *ImageRef) Flip(direction Direction) error {
	out, err := vipsFlip(r.image, direction)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Rotate executes the 'rot' operation
func (r *ImageRef) Rotate(angle Angle) error {
	out, err := vipsRotate(r.image, angle)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Label executes the 'label' operation
func (r *ImageRef) Label(labelParams *LabelParams) error {
	out, err := labelImage(r.image, labelParams)
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
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.image == image {
		return
	}

	if r.image != nil {
		clearImage(r.image)
	}

	r.image = image
}

func (r *ImageRef) exportBuffer(params *ExportParams) ([]byte, ImageType, error) {
	var buf []byte
	var err error

	format := params.Format
	if format != ImageTypeUnknown && !IsTypeSupported(format) {
		return nil, ImageTypeUnknown, fmt.Errorf("cannot save to %#v", ImageTypes[format])
	}

	if params.Quality == 0 {
		params.Quality = defaultQuality
	}

	if params.Compression == 0 {
		params.Compression = defaultCompression
	}

	switch format {
	case ImageTypeWEBP:
		buf, err = vipsSaveWebPToBuffer(r.image, false, params.Quality, params.Lossless)
	case ImageTypePNG:
		buf, err = vipsSavePNGToBuffer(r.image, false, params.Compression, params.Quality, params.Interlaced)
	case ImageTypeTIFF:
		buf, err = vipsSaveTIFFToBuffer(r.image)
	case ImageTypeHEIF:
		buf, err = vipsSaveHEIFToBuffer(r.image, params.Quality, params.Lossless)
	default:
		format = ImageTypeJPEG
		buf, err = vipsSaveJPEGToBuffer(r.image, params.Quality, false, params.Interlaced)
	}

	if err != nil {
		return nil, ImageTypeUnknown, err
	}

	return buf, format, nil
}

///////////////

func vipsHasAlpha(in *C.VipsImage) bool {
	return int(C.has_alpha_channel(in)) > 0
}

//////////////

func clearImage(ref *C.VipsImage) {
	C.clear_image(&ref)
}
