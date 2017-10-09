package vips

import (
	"bytes"
	"io"
	"math"
	"os"
)

// InputParams are options when importing an image from file or buffer
type InputParams struct {
	InputFile string
	Reader    io.Reader
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
	BlurSigma               float64
	Flip                    FlipDirection
	Width                   int
	Height                  int
	ScaleX                  float64
	ScaleY                  float64
}

func (t *TransformParams) SetTargets(width, height int, scaleX, scaleY float64) {
	t.Width = width
	t.Height = height
	t.ScaleX = scaleX
	t.ScaleY = scaleY
}

// Transform handles single image transformations
type Transform struct {
	input  *InputParams
	tx     *TransformParams
	export *ExportParams
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
	return t.Scale(scale, 0)
}

// ScaleHeight scales the height of the image proportionally
func (t *Transform) ScaleHeight(scale float64) *Transform {
	return t.Scale(0, scale)
}

// Scale the image
func (t *Transform) Scale(scaleX, scaleY float64) *Transform {
	t.tx.SetTargets(0, 0, scaleX, scaleY)
	return t
}

// Reduce the image proportionally
func (t *Transform) Reduce(scale float64) *Transform {
	if scale >= 1 {
		panic("scale must be less than 1")
	}
	return t.Scale(scale, scale)
}

// Enlarge the image proportionally
func (t *Transform) Enlarge(scale float64) *Transform {
	if scale <= 1 {
		panic("scale must be greater than 1")
	}
	return t.Scale(scale, scale)
}

// ResizeWidth resizes the image to the given width, maintaining aspect ratio
func (t *Transform) ResizeWidth(width int) *Transform {
	return t.Resize(width, 0)
}

// ResizeHeight resizes the image to the given height, maintaining aspect ratio
func (t *Transform) ResizeHeight(height int) *Transform {
	return t.Resize(0, height)
}

// Resize resizes the image to the given width and height
func (t *Transform) Resize(width, height int) *Transform {
	t.tx.SetTargets(width, height, 0, 0)
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

// StripMetadata strips ICC profile and metadata from the image
func (t *Transform) StripMetadata() *Transform {
	t.export.StripProfile = true
	t.export.StripMetadata = true
	return t
}

// Apply loads the image, applies the transform, and exports it according
// to the parameters specified
func (t *Transform) Apply() ([]byte, error) {
	defer ShutdownThread()
	startupIfNeeded()

	image, err := t.importImage()
	if err != nil {
		return nil, err
	}

	defer image.Close()

	err = t.transform(image)
	if err != nil {
		return nil, err
	}

	return t.exportImage(image)
}

func (t *Transform) importImage() (*ImageRef, error) {
	if t.input.Reader == nil {
		panic("no input source specified")
	}
	return LoadImage(t.input.Reader)
}

func (t *Transform) exportImage(image *ImageRef) ([]byte, error) {
	buf, err := vipsExportBuffer(image.Image(), t.export)
	if err != nil {
		return nil, err
	}

	if t.export.Writer != nil {
		_, err = t.export.Writer.Write(buf)
		if err != nil {
			return buf, err
		}
	}

	return buf, err
}

func (t *Transform) transform(image *ImageRef) error {
	if err := resize(image, t.tx); err != nil {
		return err
	}

	if err := postProcess(image, t.tx); err != nil {
		return err
	}

	// TODO(d): Flatten image

	return nil
}

func resize(image *ImageRef, p *TransformParams) error {
	kernel := p.ReductionSampler

	// Check for the simple scale down cases
	if (p.ScaleX > 0 && p.ScaleX <= 1) || (p.ScaleY > 0 && p.ScaleY < 1) {
		scaleX := p.ScaleX
		scaleY := p.ScaleY
		if scaleX == 0 {
			scaleX = scaleY
		} else if scaleY == 0 {
			scaleY = scaleX
		}
		if scaleX == scaleY {
			return image.Resize(scaleX, InputInt("kernel", int(kernel)))
		}
	}

	if p.Width == 0 {
		p.Width = roundFloat(p.ScaleX * float64(image.Width()))
	}
	if p.Height == 0 {
		p.Height = roundFloat(p.ScaleY * float64(image.Height()))
	}

	if p.Width == 0 || p.Height == 0 {
		return nil
	}

	shrinkX := scale(image.Width(), p.Width)
	shrinkY := scale(image.Height(), p.Height)

	cropMode := p.ResizeStrategy == ResizeStrategyCrop

	if cropMode {
		if shrinkX > 0 && shrinkY > 0 {
			shrinkX = math.Min(shrinkX, shrinkY)
		} else {
			shrinkX = math.Max(shrinkX, shrinkY)
		}
		shrinkY = shrinkX
	}

	if shrinkX != 1 || shrinkY != 1 {
		if err := image.Resize(
			1.0/shrinkX,
			InputDouble("vscale", 1.0/shrinkY),
			InputInt("kernel", int(kernel)),
		); err != nil {
			return err
		}

		// If stretching then we're done.
		if p.ResizeStrategy == ResizeStrategyStretch {
			return nil
		}
	}

	// Crop if necessary
	if cropMode {
		if err := maybeCrop(image, p); err != nil {
			return err
		}
	}

	// Now we might need to embed to match the target dimensions
	if p.Width > image.Width() || p.Height > image.Height() {
		var left, top int
		width, height := image.Width(), image.Height()
		if p.Width > image.Width() {
			width = p.Width
			left = (p.Width - image.Width()) >> 1
		}
		if p.Height > image.Height() {
			height = p.Height
			top = (p.Height - image.Height()) >> 1
		}
		if err := image.Embed(left, top, width, height, InputInt("extend", int(p.PadStrategy))); err != nil {
			return err
		}
	}

	return nil
}

func maybeCrop(image *ImageRef, p *TransformParams) error {
	if p.Width >= image.Width() && p.Height >= image.Height() {
		return nil
	}
	imageW, imageH := image.Width(), image.Height()
	width := minInt(p.Width, imageW)
	height := minInt(p.Height, imageH)
	left, top := 0, 0
	middleX := (imageW - p.Width + 1) >> 1
	middleY := (imageH - p.Height + 1) >> 1
	switch p.CropAnchor {
	case AnchorTop:
		left = middleX
	case AnchorBottom:
		left = middleX
		top = imageH - p.Height
	case AnchorRight:
		left = imageW - p.Width
		top = middleY
	case AnchorLeft:
		top = middleY
	case AnchorTopRight:
		left = imageW - p.Width
	case AnchorTopLeft:
	case AnchorBottomRight:
		left = imageW - p.Width
		top = imageH - p.Height
	case AnchorBottomLeft:
		top = imageH - p.Height
	default:
		left = middleX
		top = middleY
	}
	left = maxInt(left, 0)
	top = maxInt(top, 0)
	return image.ExtractArea(left, top, width, height)
}

func postProcess(image *ImageRef, p *TransformParams) error {
	if p.ZoomX > 0 || p.ZoomY > 0 {
		if err := image.Zoom(p.ZoomX, p.ZoomY); err != nil {
			return err
		}
	}

	if p.Flip != FlipNone {
		var err error
		switch p.Flip {
		case FlipHorizontal:
			err = image.Flip(DirectionHorizontal)
		case FlipVertical:
			err = image.Flip(DirectionVertical)
		case FlipBoth:
			err = image.Flip(DirectionHorizontal)
			if err == nil {
				err = image.Flip(DirectionVertical)
			}
		}
		if err != nil {
			return err
		}
	}

	if p.Invert {
		if err := image.Invert(); err != nil {
			return err
		}
	}

	if p.BlurSigma > 0 {
		if err := image.Gaussblur(p.BlurSigma); err != nil {
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

func scale(x, y int) float64 {
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
		r.file.Close()
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
