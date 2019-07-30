package vips

// #cgo pkg-config: vips
// #include "image.h"
import "C"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"unsafe"
)

const (
	defaultQuality     = 90
	defaultCompression = 6
)

// ImageRef contains a libvips image and manages its lifecycle. You should
// close an image when done or it will leak until the next GC
type ImageRef struct {
	image  *C.VipsImage
	format ImageType

	// NOTE: We keep a reference to this so that the input buffer is
	// never garbage collected during processing. Some image loaders use random
	// access transcoding and therefore need the original buffer to be in memory.
	buf []byte
}

type ImageMetadata struct {
	Format ImageType
	Width  int
	Height int
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
	return vipsHasICCProfile(r.image)
}

// alias to HasProfile()
func (r *ImageRef) HasICCProfile() bool {
	return r.HasProfile()
}

// HasAlpha returns if the image has an alpha layer.
func (r *ImageRef) HasAlpha() bool {
	return vipsHasAlpha(r.image)
}

// Return the orientation number as appears in the EXIF, if present
func (r *ImageRef) GetOrientation() int {
	return vipsGetMetaOrientation(r.image)
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
		Format: format,
		Width:  r.Width(),
		Height: r.Height(),
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

func (r *ImageRef) Linear1(a, b float64) error {
	out, err := vipsLinear1(r.image, a, b)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
}

// Autorot executes the 'autorot' operation
func (r *ImageRef) AutoRotate() error {
	out, err := vipsAutoRotate(r.image)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
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
	vipsRemoveICCProfile(r.image)
	// this works in place on the header
	return nil
}

// won't remove the ICC profile
func (r *ImageRef) RemoveMetadata() error {
	vipsRemoveMetadata(r.image)
	// this works in place on the header
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
	out, err := vipsResize(r.image, scale, kernel)
	if err != nil {
		return err
	}
	r.setImage(out)
	return nil
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
	if r.image == image {
		return
	}

	un := r.image
	if un != nil {
		defer unrefImage(un)
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
