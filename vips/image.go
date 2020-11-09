package vips

// #cgo pkg-config: vips
// #include "image.h"
// #include "icc_profiles.h"
import "C"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"sync"
	"unsafe"
)

// PreMultiplicationState stores the premultiplication band format of the image
type PreMultiplicationState struct {
	bandFormat BandFormat
}

// ImageRef contains a libvips image and manages its lifecycle. You need to
// close an image when done or it will leak
type ImageRef struct {
	// NOTE: We keep a reference to this so that the input buffer is
	// never garbage collected during processing. Some image loaders use random
	// access transcoding and therefore need the original buffer to be in memory.
	buf                 []byte
	image               *C.VipsImage
	format              ImageType
	lock                sync.Mutex
	preMultiplication   *PreMultiplicationState
	optimizedIccProfile string
}

// ImageMetadata is a data structure holding the width, height, orientation and other metadata of the picture.
type ImageMetadata struct {
	Format      ImageType
	Width       int
	Height      int
	Colorspace  Interpretation
	Orientation int
}

// ExportParams are options when exporting an image to file or buffer.
type ExportParams struct {
	Format        ImageType
	Quality       int
	Compression   int
	Interlaced    bool
	Lossless      bool
	Effort        int
	StripMetadata bool
}

// NewDefaultExportParams creates default values for an export when image type is not JPEG, PNG or WEBP.
// By default, govips creates interlaced, lossy images with a quality of 80/100 and compression of 6/10.
// As these are default values for a wide variety of image formats, their application varies.
// Some formats use the quality parameters, some compression, etc.
func NewDefaultExportParams() *ExportParams {
	return &ExportParams{
		Format:      ImageTypeUnknown, // defaults to the starting encoder
		Quality:     80,
		Compression: 6,
		Interlaced:  true,
		Lossless:    false,
		Effort:      4,
	}
}

// NewDefaultJPEGExportParams creates default values for an export of a JPEG image.
// By default, govips creates interlaced JPEGs with a quality of 80/100.
func NewDefaultJPEGExportParams() *ExportParams {
	return &ExportParams{
		Format:     ImageTypeJPEG,
		Quality:    80,
		Interlaced: true,
	}
}

// NewDefaultPNGExportParams creates default values for an export of a PNG image.
// By default, govips creates non-interlaced PNGs with a compression of 6/10.
func NewDefaultPNGExportParams() *ExportParams {
	return &ExportParams{
		Format:      ImageTypePNG,
		Compression: 6,
		Interlaced:  false,
	}
}

// NewDefaultWEBPExportParams creates default values for an export of a WEBP image.
// By default, govips creates lossy images with a quality of 75/100.
func NewDefaultWEBPExportParams() *ExportParams {
	return &ExportParams{
		Format:   ImageTypeWEBP,
		Quality:  75,
		Lossless: false,
		Effort:   4,
	}
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

// Metadata returns the metadata (ImageMetadata struct) of the associated ImageRef
func (r *ImageRef) Metadata() *ImageMetadata {
	return &ImageMetadata{
		Format: r.Format(),
		Width:  r.Width(),
		Height: r.Height(),
	}
}

// Copy creates a new copy of the given image.
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
	runtime.SetFinalizer(image, finalizeImage)
	return image
}

func finalizeImage(ref *ImageRef) {
	ref.Close()
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

// Format returns the initial format of the vips image when loaded.
func (r *ImageRef) Format() ImageType {
	return r.format
}

// Width returns the width of this image.
func (r *ImageRef) Width() int {
	return int(r.image.Xsize)
}

// Height returns the height of this image.
func (r *ImageRef) Height() int {
	return int(r.image.Ysize)
}

// Bands returns the number of bands for this image.
func (r *ImageRef) Bands() int {
	return int(r.image.Bands)
}

// HasProfile returns if the image has an ICC profile embedded.
func (r *ImageRef) HasProfile() bool {
	return vipsHasICCProfile(r.image)
}

// HasICCProfile checks whether the image has an ICC profile embedded. Alias to HasProfile
func (r *ImageRef) HasICCProfile() bool {
	return r.HasProfile()
}

// HasIPTC returns a boolean whether the image in question has IPTC data associated with it.
func (r *ImageRef) HasIPTC() bool {
	return vipsHasIPTC(r.image)
}

// HasAlpha returns if the image has an alpha layer.
func (r *ImageRef) HasAlpha() bool {
	return vipsHasAlpha(r.image)
}

// GetOrientation returns the orientation number as it appears in the EXIF, if present
func (r *ImageRef) GetOrientation() int {
	return vipsGetMetaOrientation(r.image)
}

// SetOrientation sets the orientation in the EXIF header of the associated image.
func (r *ImageRef) SetOrientation(orientation int) error {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return err
	}

	vipsSetMetaOrientation(out, orientation)

	r.setImage(out)
	return nil
}

// RemoveOrientation removes the EXIF orientation information of the image.
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

// Interpretation returns the current interpretation of the color space of the image.
func (r *ImageRef) Interpretation() Interpretation {
	return Interpretation(int(r.image.Type))
}

// ColorSpace returns the interpreptation of the current color space. Alias to Interpretation().
func (r *ImageRef) ColorSpace() Interpretation {
	return r.Interpretation()
}

// IsColorSpaceSupported returns a boolean whether the image's color space is supported by libvips.
func (r *ImageRef) IsColorSpaceSupported() bool {
	return vipsIsColorSpaceSupported(r.image)
}

// Export creates a byte array of the image for use.
// The function returns a byte array that can be written to a file e.g. via ioutil.WriteFile().
// N.B. govips does not currently have built-in support for directly exporting to a file.
// The function also returns a copy of the image metadata as well as an error.
func (r *ImageRef) Export(params *ExportParams) ([]byte, *ImageMetadata, error) {
	params = r.resolveExportParams(params)

	// the exported buf is not necessarily in same format as the original buf, might default to JPEG as well.
	buf, format, err := r.exportBuffer(params)
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

// Composite composites the given overlay image on top of the associated image with provided blending mode.
func (r *ImageRef) resolveExportParams(params *ExportParams) *ExportParams {
	if params == nil {
		switch r.format {
		case ImageTypeJPEG:
			params = NewDefaultJPEGExportParams()
		case ImageTypePNG:
			params = NewDefaultPNGExportParams()
		case ImageTypeWEBP:
			params = NewDefaultWEBPExportParams()
		default:
			params = NewDefaultExportParams()
		}
	}

	if params.Format == ImageTypeUnknown {
		params.Format = r.format
	}

	return params
}

func (r *ImageRef) Composite(overlay *ImageRef, mode BlendMode, x, y int) error {
	out, err := vipsComposite2(r.image, overlay.image, mode, x, y)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// ExtractBand extracts one or more bands out of the image (replacing the associated ImageRef)
func (r *ImageRef) ExtractBand(band int, num int) error {
	out, err := vipsExtractBand(r.image, band, num)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// BandJoin joins a set of images together, bandwise.
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

// AddAlpha adds an alpha channel to the associated image.
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

// PremultiplyAlpha premultiplies the alpha channel.
// See https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-premultiply
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

// UnpremultiplyAlpha unpremultiplies any alpha channel.
// See https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-unpremultiply
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

// Linear passes an image through a linear transformation (ie. output = input * a + b).
// See https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-linear
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

// Linear1 runs Linear() with a single constant.
// See https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-linear1
func (r *ImageRef) Linear1(a, b float64) error {
	out, err := vipsLinear1(r.image, a, b)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// GetRotationAngleFromExif returns the angle which the image is currently rotated in.
// This is based on the EXIF orientation tag standard.
// If no proper orientation number is provided, the picture will be assumed to be upright.
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

// AutoRotate rotates the image upright based on the EXIF Orientation tag.
// It also resets the orientation information in the EXIF tag to be 1 (i.e. upright).
// N.B. libvips does not flip images currently (i.e. no support for orientations 2, 4, 5 and 7).
// N.B. due to the HEIF image standard, HEIF images are always autorotated by default on load.
// Thus, calling AutoRotate for HEIF images is not needed.
func (r *ImageRef) AutoRotate() error {
	out, err := vipsAutoRotate(r.image)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// ExtractArea crops the image to a specified area
func (r *ImageRef) ExtractArea(left, top, width, height int) error {
	out, err := vipsExtractArea(r.image, left, top, width, height)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// RemoveICCProfile removes the ICC Profile information from the image.
// Typically browsers and other software assume images without profile to be in the sRGB color space.
func (r *ImageRef) RemoveICCProfile() error {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return err
	}

	vipsRemoveICCProfile(out)

	r.setImage(out)
	return nil
}

// OptimizeICCProfile optimizes the ICC color profile of the image.
// For two color channel images, it sets a grayscale profile.
// For color images, it sets a CMYK or non-CMYK profile based on the image metadata.
func (r *ImageRef) OptimizeICCProfile() error {
	inputProfile := r.determineInputICCProfile()
	if !r.HasICCProfile() && (inputProfile == "") {
		//No embedded ICC profile in the input image and no input profile determined, nothing to do.
		return nil
	}

	r.optimizedIccProfile = C.GoString(C.SRGB_V2_MICRO_ICC_PATH)
	if r.Bands() <= 2 {
		r.optimizedIccProfile = C.GoString(C.SGRAY_V2_MICRO_ICC_PATH)
	}

	embedded := r.HasICCProfile() && (inputProfile == "")

	out, err := vipsICCTransform(r.image, r.optimizedIccProfile, inputProfile, IntentPerceptual, 0, embedded)
	if err != nil {
		govipsLog("govips", LogLevelError, err.Error())
		return err
	}

	r.setImage(out)
	return nil
}

// RemoveMetadata removes the EXIF metadata from the image.
// N.B. this function won't remove the ICC profile and orientation because
// govips needs it to correctly display the image.
func (r *ImageRef) RemoveMetadata() error {
	out, err := vipsCopyImage(r.image)
	if err != nil {
		return err
	}

	vipsRemoveMetadata(out)

	r.setImage(out)
	return nil
}

// ToColorSpace changes the color space of the image to the interpreptation supplied as the parameter.
func (r *ImageRef) ToColorSpace(interpretation Interpretation) error {
	out, err := vipsToColorSpace(r.image, interpretation)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Flatten removes the alpha channel from the image and replaces it with the background color
func (r *ImageRef) Flatten(backgroundColor *Color) error {
	out, err := vipsFlatten(r.image, backgroundColor)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// GaussianBlur blurs the image
func (r *ImageRef) GaussianBlur(sigma float64) error {
	out, err := vipsGaussianBlur(r.image, sigma)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Sharpen sharpens the image
// sigma: sigma of the gaussian
// x1: flat/jaggy threshold
// m2: slope for jaggy areas
func (r *ImageRef) Sharpen(sigma float64, x1 float64, m2 float64) error {
	out, err := vipsSharpen(r.image, sigma, x1, m2)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Modulate the colors
func (r *ImageRef) Modulate(brightness, saturation, hue float64) error {
	var err error
	var multiplications []float64
	var additions []float64

	colorspace := r.ColorSpace()
	if colorspace == InterpretationRGB {
		colorspace = InterpretationSRGB
	}

	multiplications = []float64{brightness, saturation, 1}
	additions = []float64{0, 0, hue}

	if r.HasAlpha() {
		multiplications = append(multiplications, 1)
		additions = append(additions, 0)
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

// ModulateHSV modulates the image HSV values based on the supplier parameters.
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

// Invert inverts the image
func (r *ImageRef) Invert() error {
	out, err := vipsInvert(r.image)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Resize resizes the image based on the scale, maintaining aspect ratio
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

// ResizeWithVScale resizes the image with both horizontal as well as vertical scaling.
// The parameters are the scaling factors.
func (r *ImageRef) ResizeWithVScale(hScale, vScale float64, kernel Kernel) error {
	out, err := vipsResizeWithVScale(r.image, hScale, vScale, kernel)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Embed embeds the given picture in a new one, i.e. the opposite of ExtractArea
func (r *ImageRef) Embed(left, top, width, height int, extend ExtendStrategy) error {
	out, err := vipsEmbed(r.image, left, top, width, height, extend)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Zoom zooms the image by repeating pixels (fast nearest-neighbour)
func (r *ImageRef) Zoom(xFactor int, yFactor int) error {
	out, err := vipsZoom(r.image, xFactor, yFactor)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Flip flips the image either horizontally or vertically based on the parameter
func (r *ImageRef) Flip(direction Direction) error {
	out, err := vipsFlip(r.image, direction)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Rotate rotates the image based on the given angle
func (r *ImageRef) Rotate(angle Angle) error {
	out, err := vipsRotate(r.image, angle)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Label overlays a label on top of the image
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

func (r *ImageRef) determineInputICCProfile() (inputProfile string) {
	if r.Interpretation() == InterpretationCMYK {
		inputProfile = "cmyk"
	}
	return
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

	switch format {
	case ImageTypeWEBP:
		buf, err = vipsSaveWebPToBuffer(r.image, false, params.Quality, params.Lossless, params.Effort,
			r.optimizedIccProfile)
	case ImageTypePNG:
		buf, err = vipsSavePNGToBuffer(r.image, params.StripMetadata, params.Compression, params.Interlaced)
	case ImageTypeTIFF:
		buf, err = vipsSaveTIFFToBuffer(r.image)
	case ImageTypeHEIF:
		buf, err = vipsSaveHEIFToBuffer(r.image, params.Quality, params.Lossless)
	default:
		format = ImageTypeJPEG
		buf, err = vipsSaveJPEGToBuffer(r.image, params.Quality, params.StripMetadata, params.Interlaced)
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

// Coding represents VIPS_CODING type
type Coding int

// Coding enum
const (
	CodingError Coding = C.VIPS_CODING_ERROR
	CodingNone  Coding = C.VIPS_CODING_NONE
	CodingLABQ  Coding = C.VIPS_CODING_LABQ
	CodingRAD   Coding = C.VIPS_CODING_RAD
)
