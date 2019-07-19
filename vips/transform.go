package vips

import (
	"math"
)

// TransformParams are parameters for the transformation
type TransformParams struct {
	PadStrategy             ExtendStrategy
	ResizeStrategy          ResizeStrategy
	CropAnchor              Anchor
	ReductionSampler        Kernel
	EnlargementInterpolator Interpolator
	ZoomX                   int
	ZoomY                   int
	Invert                  bool
	Rotate                  Angle
	AutoRotate              bool
	BlurSigma               float64
	Flip                    FlipDirection
	Width                   Scalar
	Height                  Scalar
	CropOffsetX             Scalar
	CropOffsetY             Scalar
	MaxScale                float64
	Label                   *LabelParams
	SharpSigma              float64
	SharpX1                 float64
	SharpM2                 float64
}

// Transform handles single image transformations
type Transform struct {
	transformParams *TransformParams
	exportParams    *ExportParams
}

// NewTransform constructs a new transform for execution
func NewTransform() *Transform {
	return &Transform{
		transformParams: &TransformParams{
			ResizeStrategy:          ResizeStrategyAuto,
			CropAnchor:              AnchorAuto,
			ReductionSampler:        KernelLanczos3,
			EnlargementInterpolator: InterpolateBicubic,
		},
		exportParams: &ExportParams{
			Format:         ImageTypeUnknown,
			Quality:        90,
			Interpretation: InterpretationSRGB,
		},
	}
}

// Zoom an image by repeating pixels. This is fast nearest-neighbour zoom.
func (t *Transform) Zoom(x, y int) *Transform {
	t.transformParams.ZoomX = x
	t.transformParams.ZoomY = y
	return t
}

// Anchor sets the anchor for cropping
func (t *Transform) Anchor(anchor Anchor) *Transform {
	t.transformParams.CropAnchor = anchor
	return t
}

// CropOffsetX sets the target offset from the crop position
func (t *Transform) CropOffsetX(x int) *Transform {
	t.transformParams.CropOffsetX.SetInt(x)
	return t
}

// CropOffsetY sets the target offset from the crop position
func (t *Transform) CropOffsetY(y int) *Transform {
	t.transformParams.CropOffsetY.SetInt(y)
	return t
}

// CropRelativeOffsetX sets the target offset from the crop position
func (t *Transform) CropRelativeOffsetX(x float64) *Transform {
	t.transformParams.CropOffsetX.SetScale(x)
	return t
}

// CropRelativeOffsetY sets the target offset from the crop position
func (t *Transform) CropRelativeOffsetY(y float64) *Transform {
	t.transformParams.CropOffsetY.SetScale(y)
	return t
}

// Kernel sets the sampling kernel for the transform when down-scaling. Defaults to lancosz3
func (t *Transform) Kernel(kernel Kernel) *Transform {
	t.transformParams.ReductionSampler = kernel
	return t
}

// Interpolator sets the resampling interpolator when upscaling, defaults to bicubic
func (t *Transform) Interpolator(interp Interpolator) *Transform {
	t.transformParams.EnlargementInterpolator = interp
	return t
}

// ResizeStrategy sets the strategy when resizing an image
func (t *Transform) ResizeStrategy(strategy ResizeStrategy) *Transform {
	t.transformParams.ResizeStrategy = strategy
	return t
}

// PadStrategy sets the strategy when the image must be padded to maintain aspect ratoi
func (t *Transform) PadStrategy(strategy ExtendStrategy) *Transform {
	t.transformParams.PadStrategy = strategy
	return t
}

// Invert inverts the image color
func (t *Transform) Invert() *Transform {
	t.transformParams.Invert = true
	return t
}

// Flip flips the image horizontally or vertically
func (t *Transform) Flip(flip FlipDirection) *Transform {
	t.transformParams.Flip = flip
	return t
}

// GaussianBlur applies a gaussian blur to the image
func (t *Transform) GaussianBlur(sigma float64) *Transform {
	t.transformParams.BlurSigma = sigma
	return t
}

// Sharpen applies a sharpen to the image
func (t *Transform) Sharpen(sigma float64, x1 float64, m2 float64) *Transform {
	t.transformParams.SharpSigma = sigma
	t.transformParams.SharpX1 = x1
	t.transformParams.SharpM2 = m2
	return t
}

// AutoRotate rotates image by a the embedded metadata (EXIF Orientation, etc.)
func (t *Transform) AutoRotate() *Transform {
	t.transformParams.AutoRotate = true
	return t
}

// Rotate rotates image by a multiple of 90 degrees
func (t *Transform) Rotate(angle Angle) *Transform {
	t.transformParams.Rotate = angle
	return t
}

// Embed this image appropriately if resized according to a new aspect ratio
func (t *Transform) Embed(extend ExtendStrategy) *Transform {
	t.transformParams.ResizeStrategy = ResizeStrategyEmbed
	t.transformParams.PadStrategy = extend
	return t
}

// Crop an image, width and height must be equal to or less than image size
func (t *Transform) Crop(anchor Anchor) *Transform {
	t.transformParams.ResizeStrategy = ResizeStrategyCrop
	return t
}

// Stretch an image without maintaining aspect ratio
func (t *Transform) Stretch() *Transform {
	t.transformParams.ResizeStrategy = ResizeStrategyCrop
	return t
}

// ScaleWidth scales the image by its width proportionally
func (t *Transform) ScaleWidth(scale float64) *Transform {
	t.transformParams.Width.SetScale(scale)
	return t
}

// ScaleHeight scales the height of the image proportionally
func (t *Transform) ScaleHeight(scale float64) *Transform {
	t.transformParams.Height.SetScale(scale)
	return t
}

// Scale the image
func (t *Transform) Scale(scale float64) *Transform {
	t.transformParams.Width.SetScale(scale)
	t.transformParams.Height.SetScale(scale)
	return t
}

// MaxScale sets the max scale factor that this image can be enlarged or reduced by
func (t *Transform) MaxScale(max float64) *Transform {
	t.transformParams.MaxScale = max
	return t
}

// ResizeWidth resizes the image to the given width, maintaining aspect ratio
func (t *Transform) ResizeWidth(width int) *Transform {
	t.transformParams.Width.SetInt(width)
	return t
}

// ResizeHeight resizes the image to the given height, maintaining aspect ratio
func (t *Transform) ResizeHeight(height int) *Transform {
	t.transformParams.Height.SetInt(height)
	return t
}

// Resize resizes the image to the given width and height
func (t *Transform) Resize(width, height int) *Transform {
	t.transformParams.Width.SetInt(width)
	t.transformParams.Height.SetInt(height)
	return t
}

func (t *Transform) Label(lp *LabelParams) *Transform {
	if lp.Text == "" {
		t.transformParams.Label = nil
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
	t.transformParams.Label = &label
	return t
}

// Format sets the image format of the input image when exporting. Defaults to JPEG
func (t *Transform) Format(format ImageType) *Transform {
	t.exportParams.Format = format
	return t
}

// Quality sets the quality value for image formats that support it
func (t *Transform) Quality(quality int) *Transform {
	t.exportParams.Quality = quality
	return t
}

// Compression sets the compression value for image formats that support it
func (t *Transform) Compression(compression int) *Transform {
	t.exportParams.Compression = compression
	return t
}

// Lossless uses lossless compression for image formats that support both lossy and lossless e.g. webp
func (t *Transform) Lossless() *Transform {
	t.exportParams.Lossless = true
	return t
}

// StripMetadata strips metadata from the image
func (t *Transform) StripMetadata() *Transform {
	t.exportParams.StripMetadata = true
	return t
}

// StripProfile strips ICC profile from the image
func (t *Transform) StripProfile() *Transform {
	t.exportParams.StripProfile = true
	return t
}

// BackgroundColor sets the background color of the image when a transparent
// image is flattened
func (t *Transform) BackgroundColor(color *Color) *Transform {
	t.exportParams.BackgroundColor = color
	return t
}

// Interpretation sets interpretation for image
func (t *Transform) Interpretation(interpretation Interpretation) *Transform {
	t.exportParams.Interpretation = interpretation
	return t
}

// Interlaced uses interlaced for image that support it
func (t *Transform) Interlaced() *Transform {
	t.exportParams.Interlaced = true
	return t
}

// Apply the transform, returns the modified image
func (t *Transform) Apply(image *ImageRef) (*ImageRef, error) {
	startupIfNeeded()

	defer ShutdownThread()

	return newBlackboard(image, t.transformParams).execute()
}

// Return the formatted buffer of the transformed image, and its metadata
func (t *Transform) ApplyAndExport(image *ImageRef) ([]byte, *ImageMetadata, error) {
	i, err := t.Apply(image)
	if err != nil {
		return nil, nil, err
	}

	return i.Export(t.exportParams)
}

// blackboard is an object that tracks transient data during a transformation
type blackboard struct {
	*TransformParams
	image        *ImageRef
	aspectRatio  float64
	targetWidth  int
	targetHeight int
	targetScale  float64
	cropOffsetX  int
	cropOffsetY  int
}

// newBlackboard creates a new blackboard object meant for transformation data
func newBlackboard(imageRef *ImageRef, transformParams *TransformParams) *blackboard {
	bb := &blackboard{
		TransformParams: transformParams,
		image:           imageRef,
	}
	imageWidth := imageRef.Width()
	imageHeight := imageRef.Height()
	bb.aspectRatio = ratio(imageWidth, imageHeight)
	bb.cropOffsetX = transformParams.CropOffsetX.GetRounded(imageWidth)
	bb.cropOffsetY = transformParams.CropOffsetY.GetRounded(imageHeight)

	if transformParams.Width.Value == 0 && transformParams.Height.Value == 0 {
		return bb
	}

	bb.targetWidth = transformParams.Width.GetRounded(imageWidth)
	bb.targetHeight = transformParams.Height.GetRounded(imageHeight)

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

	if transformParams.Width.Relative && transformParams.Height.Relative {
		sx, sy := transformParams.Width.Value, transformParams.Height.Value
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

func (b *blackboard) execute() (*ImageRef, error) {
	if err := b.resize(); err != nil {
		return nil, err
	}

	if err := b.postProcess(); err != nil {
		return nil, err
	}

	return b.image, nil
}

func (b *blackboard) resize() error {
	var err error
	kernel := b.ReductionSampler

	// Check for the simple scale down cases
	if b.targetScale != 0 {
		err := b.image.Resize(b.targetScale, b.targetScale, kernel)
		if err != nil {
			return err
		}
	}

	if b.targetHeight == 0 && b.targetWidth == 0 {
		return nil
	}

	shrinkX := ratio(b.width(), b.targetWidth)
	shrinkY := ratio(b.height(), b.targetHeight)

	cropMode := b.ResizeStrategy == ResizeStrategyCrop
	stretchMode := b.ResizeStrategy == ResizeStrategyStretch

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
		err = b.image.Resize(1.0/shrinkX, 1.0/shrinkY, kernel)
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
		if err := b.maybeCrop(); err != nil {
			return err
		}
	}

	if err := b.maybeEmbed(); err != nil {
		return err
	}

	return nil
}

func (b *blackboard) maybeCrop() error {
	var err error
	imageW, imageH := b.width(), b.height()

	if b.targetWidth >= imageW && b.targetHeight >= imageH {
		return nil
	}

	width := minInt(b.targetWidth, imageW)
	height := minInt(b.targetHeight, imageH)
	left, top := 0, 0
	middleX := (imageW - b.targetWidth + 1) >> 1
	middleY := (imageH - b.targetHeight + 1) >> 1
	if b.cropOffsetX != 0 || b.cropOffsetY != 0 {
		if b.cropOffsetX >= 0 {
			left = middleX + minInt(b.cropOffsetX, middleX)
		} else {
			left = middleX - maxInt(b.cropOffsetX, middleX)
		}
		if b.cropOffsetY >= 0 {
			top = middleY + minInt(b.cropOffsetY, middleY)
		} else {
			top = middleY - maxInt(b.cropOffsetY, middleY)
		}
	} else {
		switch b.CropAnchor {
		case AnchorTop:
			left = middleX
		case AnchorBottom:
			left = middleX
			top = imageH - b.targetHeight
		case AnchorRight:
			left = imageW - b.targetWidth
			top = middleY
		case AnchorLeft:
			top = middleY
		case AnchorTopRight:
			left = imageW - b.targetWidth
		case AnchorTopLeft:
		case AnchorBottomRight:
			left = imageW - b.targetWidth
			top = imageH - b.targetHeight
		case AnchorBottomLeft:
			top = imageH - b.targetHeight
		default:
			left = middleX
			top = middleY
		}
	}
	left = maxInt(left, 0)
	top = maxInt(top, 0)
	if left+width > imageW {
		width = imageW - left
		b.targetWidth = width
	}
	if top+height > imageH {
		height = imageH - top
		b.targetHeight = height
	}
	err = b.image.ExtractArea(left, top, width, height)

	return err
}

func (b *blackboard) maybeEmbed() error {
	var err error
	imageW, imageH := b.width(), b.height()

	// Now we might need to embed to match the target dimensions
	if b.targetWidth > imageW || b.targetHeight > imageH {
		var left, top int
		width, height := imageW, imageH
		if b.targetWidth > imageW {
			width = b.targetWidth
			left = (b.targetWidth - imageW) >> 1
		}
		if b.targetHeight > imageH {
			height = b.targetHeight
			top = (b.targetHeight - imageH) >> 1
		}
		err = b.image.Embed(left, top, width, height, b.PadStrategy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *blackboard) postProcess() error {
	var err error
	if b.ZoomX > 0 || b.ZoomY > 0 {
		err = b.image.Zoom(b.ZoomX, b.ZoomY)
		if err != nil {
			return err
		}
	}

	if b.Flip != FlipNone {
		var err error
		switch b.Flip {
		case FlipHorizontal:
			err = b.image.Flip(DirectionHorizontal)
		case FlipVertical:
			err = b.image.Flip(DirectionVertical)
		case FlipBoth:
			err = b.image.Flip(DirectionHorizontal)
			if err == nil {
				err = b.image.Flip(DirectionVertical)
			}
		}
		if err != nil {
			return err
		}
	}

	if b.Invert {
		err = b.image.Invert()
		if err != nil {
			return err
		}
	}

	if b.BlurSigma > 0 {
		err = b.image.GaussianBlur(b.BlurSigma)
		if err != nil {
			return err
		}
	}

	if b.SharpSigma > 0 {
		err = b.image.Sharpen(b.SharpSigma, b.SharpX1, b.SharpM2)
		if err != nil {
			return err
		}
	}

	if b.AutoRotate {
		err = b.image.AutoRotate()
		if err != nil {
			return err
		}
	}

	if b.Rotate > 0 {
		err = b.image.Rotate(b.Rotate)
		if err != nil {
			return err
		}
	}

	if b.Label != nil {
		err = b.image.Label(b.Label)
		if err != nil {
			return err
		}
	}

	return nil
}

// width returns the width of the in-flight image
func (b *blackboard) width() int {
	return b.image.Width()
}

// height returns the height of the in-flight image
func (b *blackboard) height() int {
	return b.image.Height()
}
