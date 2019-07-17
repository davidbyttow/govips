package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"os"
)

// InputParams are options when importing an image from file or buffer
type InputParams struct {
	Reader io.Reader
	Image  *ImageRef
}

// TransformParams are parameters for the transformation
type TransformParams struct {
	PadStrategy             Extend
	ResizeStrategy          ResizeStrategy
	CropAnchor              Anchor
	ReductionSampler        Kernel
	EnlargementInterpolator Interpolator
	ZoomX                   int
	ZoomY                   int
	Invert                  bool
	Rotate                  Angle
	BlurSigma               float64
	Flip                    FlipDirection
	Width                   Scalar
	Height                  Scalar
	CropOffsetX             Scalar
	CropOffsetY             Scalar
	MaxScale                float64
	Label                   *LabelParams
}

// Transform handles single image transformations
type Transform struct {
	input        *InputParams
	tx           *TransformParams
	export       *ExportParams
	targetWidth  int
	targetHeight int
	cropOffsetX  int
	cropOffsetY  int
	source       []byte
}

// NewTransform constructs a new transform for execution
func NewTransform() *Transform {
	return &Transform{
		input: &InputParams{},
		tx: &TransformParams{
			ResizeStrategy:          ResizeStrategyAuto,
			CropAnchor:              AnchorAuto,
			ReductionSampler:        KernelLanczos3,
			EnlargementInterpolator: InterpolateBicubic,
		},
		export: &ExportParams{
			Format:         ImageTypeUnknown,
			Quality:        90,
			Interpretation: InterpretationSRGB,
		},
	}
}

// Image sets the image to operate on
func (t *Transform) Image(image *ImageRef) *Transform {
	t.input.Image = image
	return t
}

// LoadFile loads a file into the transform
func (t *Transform) LoadFile(file string) *Transform {
	t.input.Reader = LazyOpen(file)
	return t
}

// LoadBuffer loads a buffer into the transform
func (t *Transform) LoadBuffer(buf []byte) *Transform {
	t.input.Reader = bytes.NewBuffer(buf)
	return t
}

// Load loads a buffer into the transform
func (t *Transform) Load(reader io.Reader) *Transform {
	t.input.Reader = reader
	return t
}

// Output outputs the transform to a buffer and closes it
func (t *Transform) Output(writer io.Writer) *Transform {
	t.export.Writer = writer
	return t
}

// OutputBytes outputs the transform to a buffer and closes it
func (t *Transform) OutputBytes() *Transform {
	t.export.Writer = nil
	return t
}

// OutputFile outputs the transform to a file and closes it
func (t *Transform) OutputFile(file string) *Transform {
	t.export.Writer = LazyCreate(file)
	return t
}

// Zoom an image by repeating pixels. This is fast nearest-neighbour zoom.
func (t *Transform) Zoom(x, y int) *Transform {
	t.tx.ZoomX = x
	t.tx.ZoomY = y
	return t
}

// Anchor sets the anchor for cropping
func (t *Transform) Anchor(anchor Anchor) *Transform {
	t.tx.CropAnchor = anchor
	return t
}

// CropOffsetX sets the target offset from the crop position
func (t *Transform) CropOffsetX(x int) *Transform {
	t.tx.CropOffsetX.SetInt(x)
	return t
}

// CropOffsetY sets the target offset from the crop position
func (t *Transform) CropOffsetY(y int) *Transform {
	t.tx.CropOffsetY.SetInt(y)
	return t
}

// CropRelativeOffsetX sets the target offset from the crop position
func (t *Transform) CropRelativeOffsetX(x float64) *Transform {
	t.tx.CropOffsetX.SetScale(x)
	return t
}

// CropRelativeOffsetY sets the target offset from the crop position
func (t *Transform) CropRelativeOffsetY(y float64) *Transform {
	t.tx.CropOffsetY.SetScale(y)
	return t
}

// Kernel sets the sampling kernel for the transform when down-scaling. Defaults to lancosz3
func (t *Transform) Kernel(kernel Kernel) *Transform {
	t.tx.ReductionSampler = kernel
	return t
}

// Interpolator sets the resampling interpolator when upscaling, defaults to bicubic
func (t *Transform) Interpolator(interp Interpolator) *Transform {
	t.tx.EnlargementInterpolator = interp
	return t
}

// ResizeStrategy sets the strategy when resizing an image
func (t *Transform) ResizeStrategy(strategy ResizeStrategy) *Transform {
	t.tx.ResizeStrategy = strategy
	return t
}

// PadStrategy sets the strategy when the image must be padded to maintain aspect ratoi
func (t *Transform) PadStrategy(strategy Extend) *Transform {
	t.tx.PadStrategy = strategy
	return t
}

// Invert inverts the image color
func (t *Transform) Invert() *Transform {
	t.tx.Invert = true
	return t
}

// Flip flips the image horizontally or vertically
func (t *Transform) Flip(flip FlipDirection) *Transform {
	t.tx.Flip = flip
	return t
}

// GaussBlur applies a gaussian blur to the image
func (t *Transform) GaussBlur(sigma float64) *Transform {
	t.tx.BlurSigma = sigma
	return t
}

// Rotate rotates image by a multiple of 90 degrees
func (t *Transform) Rotate(angle Angle) *Transform {
	t.tx.Rotate = angle
	return t
}

// Embed this image appropriately if resized according to a new aspect ratio
func (t *Transform) Embed(extend Extend) *Transform {
	t.tx.ResizeStrategy = ResizeStrategyEmbed
	t.tx.PadStrategy = extend
	return t
}

// Crop an image, width and height must be equal to or less than image size
func (t *Transform) Crop(anchor Anchor) *Transform {
	t.tx.ResizeStrategy = ResizeStrategyCrop
	return t
}

// Stretch an image without maintaining aspect ratio
func (t *Transform) Stretch() *Transform {
	t.tx.ResizeStrategy = ResizeStrategyCrop
	return t
}

// ScaleWidth scales the image by its width proportionally
func (t *Transform) ScaleWidth(scale float64) *Transform {
	t.tx.Width.SetScale(scale)
	return t
}

// ScaleHeight scales the height of the image proportionally
func (t *Transform) ScaleHeight(scale float64) *Transform {
	t.tx.Height.SetScale(scale)
	return t
}

// Scale the image
func (t *Transform) Scale(scale float64) *Transform {
	t.tx.Width.SetScale(scale)
	t.tx.Height.SetScale(scale)
	return t
}

// MaxScale sets the max scale factor that this image can be enlarged or reduced by
func (t *Transform) MaxScale(max float64) *Transform {
	t.tx.MaxScale = max
	return t
}

// ResizeWidth resizes the image to the given width, maintaining aspect ratio
func (t *Transform) ResizeWidth(width int) *Transform {
	t.tx.Width.SetInt(width)
	return t
}

// ResizeHeight resizes the image to the given height, maintaining aspect ratio
func (t *Transform) ResizeHeight(height int) *Transform {
	t.tx.Height.SetInt(height)
	return t
}

// Resize resizes the image to the given width and height
func (t *Transform) Resize(width, height int) *Transform {
	t.tx.Width.SetInt(width)
	t.tx.Height.SetInt(height)
	return t
}

func (t *Transform) Label(lp *LabelParams) *Transform {
	if lp.Text == "" {
		t.tx.Label = nil
		return t
	}

	label := *lp

	// Defaults
	if label.Width.IsZero() {
		label.Width.SetScale(1)
	}
	if label.Height.IsZero() {
		label.Height.SetScale(1)
	}
	if label.Font == "" {
		label.Font = DefaultFont
	}
	if label.Opacity == 0 {
		label.Opacity = 1
	}
	t.tx.Label = &label
	return t
}

// Format sets the image format of the input image when exporting. Defaults to JPEG
func (t *Transform) Format(format ImageType) *Transform {
	t.export.Format = format
	return t
}

// Quality sets the quality value for image formats that support it
func (t *Transform) Quality(quality int) *Transform {
	t.export.Quality = quality
	return t
}

// Compression sets the compression value for image formats that support it
func (t *Transform) Compression(compression int) *Transform {
	t.export.Compression = compression
	return t
}

// Lossless uses lossless compression for image formats that support both lossy and lossless e.g. webp
func (t *Transform) Lossless() *Transform {
	t.export.Lossless = true
	return t
}

// StripMetadata strips metadata from the image
func (t *Transform) StripMetadata() *Transform {
	t.export.StripMetadata = true
	return t
}

// StripProfile strips ICC profile from the image
func (t *Transform) StripProfile() *Transform {
	t.export.StripProfile = true
	return t
}

// BackgroundColor sets the background color of the image when a transparent
// image is flattened
func (t *Transform) BackgroundColor(color Color) *Transform {
	t.export.BackgroundColor = &color
	return t
}

// Interpretation sets interpretation for image
func (t *Transform) Interpretation(interpretation Interpretation) *Transform {
	t.export.Interpretation = interpretation
	return t
}

// Interlaced uses interlaced for image that support it
func (t *Transform) Interlaced() *Transform {
	t.export.Interlaced = true
	return t
}

// Apply loads the image, applies the transform, and exports it according
// to the parameters specified
func (t *Transform) Apply() ([]byte, ImageType, error) {
	defer ShutdownThread()
	defer func() {
		t.source = nil
	}()
	startupIfNeeded()

	input, imageType, err := t.importImage()
	if err != nil {
		return nil, ImageTypeUnknown, err
	}
	if input == nil {
		return nil, ImageTypeUnknown, errors.New("vips: image not found")
	}

	transformed, err := t.transform(input, imageType)
	if err != nil {
		return nil, ImageTypeUnknown, err
	}
	defer unrefImage(transformed)

	return t.exportImage(transformed, imageType)
}

func (t *Transform) importImage() (*C.VipsImage, ImageType, error) {
	if t.input.Image != nil {
		cpy, err := vipsCopyImage(t.input.Image.image)
		return cpy, t.input.Image.Format(), err
	}

	if t.input.Reader == nil {
		panic("no input source specified")
	}

	var err error
	t.source, err = ioutil.ReadAll(t.input.Reader)
	if err != nil {
		return nil, ImageTypeUnknown, nil
	}

	return vipsLoadFromBuffer(t.source)
}

func (t *Transform) exportImage(image *C.VipsImage, imageType ImageType) ([]byte, ImageType, error) {
	if t.export.Format == ImageTypeUnknown {
		t.export.Format = imageType
	}

	buf, format, err := vipsExportBuffer(image, t.export)
	if err != nil {
		return nil, ImageTypeUnknown, err
	}

	if t.export.Writer != nil {
		_, err = t.export.Writer.Write(buf)
		if err != nil {
			return buf, format, err
		}
	}

	return buf, format, err
}

// Blackboard is an object that tracks transient data during a transformation
type Blackboard struct {
	*TransformParams
	image        *C.VipsImage
	imageType    ImageType
	aspectRatio  float64
	targetWidth  int
	targetHeight int
	targetScale  float64
	cropOffsetX  int
	cropOffsetY  int
}

// newBlackboard creates a new blackboard object meant for transformation data
func newBlackboard(image *C.VipsImage, imageType ImageType, p *TransformParams) *Blackboard {
	bb := &Blackboard{
		TransformParams: p,
		image:           image,
		imageType:       imageType,
	}
	imageWidth := int(image.Xsize)
	imageHeight := int(image.Ysize)
	bb.aspectRatio = ratio(imageWidth, imageHeight)
	bb.cropOffsetX = p.CropOffsetX.GetRounded(imageWidth)
	bb.cropOffsetY = p.CropOffsetY.GetRounded(imageHeight)

	if p.Width.Value == 0 && p.Height.Value == 0 {
		return bb
	}

	bb.targetWidth = p.Width.GetRounded(imageWidth)
	bb.targetHeight = p.Height.GetRounded(imageHeight)

	if bb.MaxScale > 0 {
		if bb.targetWidth > 0 && ratio(bb.targetWidth, imageWidth) > bb.MaxScale {
			bb.targetWidth = int(float64(imageWidth) * bb.MaxScale)
		}
		if bb.targetHeight > 0 && ratio(bb.targetHeight, imageHeight) > bb.MaxScale {
			bb.targetHeight = int(float64(imageHeight) * bb.MaxScale)
		}
	}

	switch {
	case bb.targetWidth > 0 && bb.targetHeight > 0:
		// Nothing to do
	case bb.targetWidth > 0:
		bb.targetHeight = roundFloat(ratio(bb.targetWidth, imageWidth) * float64(imageHeight))
	case bb.targetHeight > 0:
		bb.targetWidth = roundFloat(ratio(bb.targetHeight, imageHeight) * float64(imageWidth))
	}

	if p.Width.Relative && p.Height.Relative {
		sx, sy := p.Width.Value, p.Height.Value
		if sx == 0 {
			sx = sy
		} else if sy == 0 {
			sy = sx
		}
		if sx == sy {
			bb.targetScale = sx
		}
	}

	if bb.MaxScale != 0 && bb.targetScale > bb.MaxScale {
		bb.targetScale = bb.MaxScale
	}

	return bb
}

// Width returns the width of the in-flight image
func (bb *Blackboard) Width() int {
	return int(bb.image.Xsize)
}

// Height returns the height of the in-flight image
func (bb *Blackboard) Height() int {
	return int(bb.image.Ysize)
}

func (t *Transform) transform(image *C.VipsImage, imageType ImageType) (*C.VipsImage, error) {
	bb := newBlackboard(image, imageType, t.tx)
	if err := resize(bb); err != nil {
		return image, err
	}

	if err := postProcess(bb); err != nil {
		return image, err
	}

	return bb.image, nil
}

func resize(bb *Blackboard) error {
	var err error
	kernel := bb.ReductionSampler

	// Check for the simple scale down cases
	if bb.targetScale != 0 {
		bb.image, err = vipsResize(bb.image, bb.targetScale, bb.targetScale, kernel)
		if err != nil {
			return err
		}
	}

	if bb.targetHeight == 0 && bb.targetWidth == 0 {
		return nil
	}

	shrinkX := ratio(bb.Width(), bb.targetWidth)
	shrinkY := ratio(bb.Height(), bb.targetHeight)

	cropMode := bb.ResizeStrategy == ResizeStrategyCrop
	stretchMode := bb.ResizeStrategy == ResizeStrategyStretch

	if !stretchMode {
		if shrinkX > 0 && shrinkY > 0 {
			if cropMode {
				shrinkX = math.Min(shrinkX, shrinkY)
			} else {
				shrinkX = math.Max(shrinkX, shrinkY)
			}
		} else {
			if cropMode {
				shrinkX = math.Min(shrinkX, shrinkY)
			} else {
				shrinkX = math.Max(shrinkX, shrinkY)
			}
		}
		shrinkY = shrinkX
	}

	if shrinkX != 1 || shrinkY != 1 {
		bb.image, err = vipsResize(bb.image, 1.0/shrinkX, 1.0/shrinkY, kernel)
		if err != nil {
			return err
		}

		// If stretching then we're done.
		if stretchMode {
			return nil
		}
	}

	// Crop if necessary
	if cropMode {
		if err := maybeCrop(bb); err != nil {
			return err
		}
	}

	if err := maybeEmbed(bb); err != nil {
		return err
	}

	return nil
}

func maybeCrop(bb *Blackboard) error {
	var err error
	imageW, imageH := bb.Width(), bb.Height()

	if bb.targetWidth >= imageW && bb.targetHeight >= imageH {
		return nil
	}

	width := minInt(bb.targetWidth, imageW)
	height := minInt(bb.targetHeight, imageH)
	left, top := 0, 0
	middleX := (imageW - bb.targetWidth + 1) >> 1
	middleY := (imageH - bb.targetHeight + 1) >> 1
	if bb.cropOffsetX != 0 || bb.cropOffsetY != 0 {
		if bb.cropOffsetX >= 0 {
			left = middleX + minInt(bb.cropOffsetX, middleX)
		} else {
			left = middleX - maxInt(bb.cropOffsetX, middleX)
		}
		if bb.cropOffsetY >= 0 {
			top = middleY + minInt(bb.cropOffsetY, middleY)
		} else {
			top = middleY - maxInt(bb.cropOffsetY, middleY)
		}
	} else {
		switch bb.CropAnchor {
		case AnchorTop:
			left = middleX
		case AnchorBottom:
			left = middleX
			top = imageH - bb.targetHeight
		case AnchorRight:
			left = imageW - bb.targetWidth
			top = middleY
		case AnchorLeft:
			top = middleY
		case AnchorTopRight:
			left = imageW - bb.targetWidth
		case AnchorTopLeft:
		case AnchorBottomRight:
			left = imageW - bb.targetWidth
			top = imageH - bb.targetHeight
		case AnchorBottomLeft:
			top = imageH - bb.targetHeight
		default:
			left = middleX
			top = middleY
		}
	}
	left = maxInt(left, 0)
	top = maxInt(top, 0)
	if left+width > imageW {
		width = imageW - left
		bb.targetWidth = width
	}
	if top+height > imageH {
		height = imageH - top
		bb.targetHeight = height
	}
	bb.image, err = vipsExtractArea(bb.image, left, top, width, height)
	return err
}

func maybeEmbed(bb *Blackboard) error {
	var err error
	imageW, imageH := bb.Width(), bb.Height()

	// Now we might need to embed to match the target dimensions
	if bb.targetWidth > imageW || bb.targetHeight > imageH {
		var left, top int
		width, height := imageW, imageH
		if bb.targetWidth > imageW {
			width = bb.targetWidth
			left = (bb.targetWidth - imageW) >> 1
		}
		if bb.targetHeight > imageH {
			height = bb.targetHeight
			top = (bb.targetHeight - imageH) >> 1
		}
		bb.image, err = vipsEmbed(bb.image, left, top, width, height, bb.PadStrategy)
		if err != nil {
			return err
		}
	}

	return nil
}

func postProcess(bb *Blackboard) error {
	var err error
	if bb.ZoomX > 0 || bb.ZoomY > 0 {
		bb.image, err = vipsZoom(bb.image, bb.ZoomX, bb.ZoomY)
		if err != nil {
			return err
		}
	}

	if bb.Flip != FlipNone {
		var err error
		switch bb.Flip {
		case FlipHorizontal:
			bb.image, err = vipsFlip(bb.image, DirectionHorizontal)
		case FlipVertical:
			bb.image, err = vipsFlip(bb.image, DirectionVertical)
		case FlipBoth:
			bb.image, err = vipsFlip(bb.image, DirectionHorizontal)
			if err == nil {
				bb.image, err = vipsFlip(bb.image, DirectionVertical)
			}
		}
		if err != nil {
			return err
		}
	}

	if bb.Invert {
		bb.image, err = vipsInvert(bb.image)
		if err != nil {
			return err
		}
	}

	if bb.BlurSigma > 0 {
		bb.image, err = vipsGaussianBlur(bb.image, bb.BlurSigma)
		if err != nil {
			return err
		}
	}

	if bb.Rotate > 0 {
		bb.image, err = vipsRotate(bb.image, bb.Rotate)
		if err != nil {
			return err
		}
	}

	if bb.Label != nil {
		bb.image, err = vipsLabel(bb.image, *bb.Label)
		if err != nil {
			return err
		}
	}

	return nil
}

func minInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func maxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func ratio(x, y int) float64 {
	if x == y {
		return 1
	}
	return float64(x) / float64(y)
}

func roundFloat(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

// LazyFile is a lazy reader or writer
// TODO(d): Move this to AF
type LazyFile struct {
	name string
	file *os.File
}

func LazyOpen(file string) io.Reader {
	return &LazyFile{name: file}
}

func LazyCreate(file string) io.Writer {
	return &LazyFile{name: file}
}

func (r *LazyFile) Read(p []byte) (n int, err error) {
	if r.file == nil {
		f, err := os.Open(r.name)
		if err != nil {
			return 0, err
		}
		r.file = f
	}
	return r.file.Read(p)
}

func (r *LazyFile) Close() error {
	if r.file != nil {
		_ = r.file.Close()
		r.file = nil
	}
	return nil
}

func (r *LazyFile) Write(p []byte) (n int, err error) {
	if r.file == nil {
		f, err := os.Create(r.name)
		if err != nil {
			return 0, err
		}
		r.file = f
	}
	return r.file.Write(p)
}

type Scalar struct {
	Value    float64
	Relative bool
}

func ValueOf(value float64) Scalar {
	return Scalar{value, false}
}

func ScaleOf(value float64) Scalar {
	return Scalar{value, true}
}

func (s *Scalar) IsZero() bool {
	return s.Value == 0 && !s.Relative
}

func (s *Scalar) SetInt(value int) {
	s.Set(float64(value))
}

func (s *Scalar) Set(value float64) {
	s.Value = value
	s.Relative = false
}

func (s *Scalar) SetScale(f float64) {
	s.Value = f
	s.Relative = true
}

func (s *Scalar) Get(base int) float64 {
	if s.Relative {
		return s.Value * float64(base)
	}
	return s.Value
}

func (s *Scalar) GetRounded(base int) int {
	return roundFloat(s.Get(base))
}
