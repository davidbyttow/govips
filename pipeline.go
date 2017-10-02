package vips

import (
	"io/ioutil"
	"math"
)

// Pipeline is an interface for chaining together image operations
type Pipeline struct {
	bb    *blackboard
	image *ImageRef
	err   error
}

// NewPipeline constructs a new pipeline for execution
func NewPipeline() *Pipeline {
	startupIfNeeded()
	return &Pipeline{newBlackboard(), nil, nil}
}

// LoadFile loads a file into the pipeline
func (pipe *Pipeline) LoadFile(file string) *Pipeline {
	pipe.reset()
	pipe.image, pipe.err = NewImageFromFile(file)
	return pipe
}

// LoadBuffer loads a buffer into the pipeline
func (pipe *Pipeline) LoadBuffer(buf []byte) *Pipeline {
	pipe.reset()
	pipe.image, pipe.err = NewImageFromBuffer(buf)
	return pipe
}

// Output exports the pipeline to a buffer and closes it
func (pipe *Pipeline) Output() ([]byte, error) {
	pipe.checkImage()
	return pipe.export()
}

// OutputFile outputs the pipeline to a file and closes it
func (pipe *Pipeline) OutputFile(file string) error {
	pipe.checkImage()
	buf, err := pipe.export()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, buf, 0644)
}

func (pipe *Pipeline) export() ([]byte, error) {
	defer ShutdownThread()
	defer pipe.reset()

	if pipe.err != nil {
		return nil, pipe.err
	}
	return vipsExportBuffer(pipe.image.Image(), &pipe.bb.ExportParams)
}

// Error returns the error encountered during the pipeline
func (pipe *Pipeline) Error() error {
	return pipe.err
}

// Zoom an image by repeating pixels. This is fast nearest-neighbour zoom.
func (pipe *Pipeline) Zoom(x, y int) *Pipeline {
	pipe.checkImage()
	if err := pipe.image.Zoom(x, y); err != nil {
		pipe.err = err
	}
	return pipe
}

// Anchor sets the anchor for cropping
func (pipe *Pipeline) Anchor(anchor Anchor) *Pipeline {
	pipe.bb.CropAnchor = anchor
	return pipe
}

// Kernel sets the sampling kernel for the pipeline when down-scaling. Defaults to lancosz3
func (pipe *Pipeline) Kernel(kernel Kernel) *Pipeline {
	pipe.bb.ReductionSampler = kernel
	return pipe
}

// Interpolator sets the resampling interpolator when upscaling, defaults to bicubic
func (pipe *Pipeline) Interpolator(interp Interpolator) *Pipeline {
	pipe.bb.EnlargementInterpolator = interp
	return pipe
}

// ResizeStrategy sets the strategy when resizing an image
func (pipe *Pipeline) ResizeStrategy(strategy ResizeStrategy) *Pipeline {
	pipe.bb.ResizeStrategy = strategy
	return pipe
}

// PadStrategy sets the strategy when the image must be padded to maintain aspect ratoi
func (pipe *Pipeline) PadStrategy(strategy Extend) *Pipeline {
	pipe.bb.PadStrategy = strategy
	return pipe
}

// Format sets the image format of the input image when exporting. Defaults to JPEG
func (pipe *Pipeline) Format(format ImageType) *Pipeline {
	pipe.bb.Format = format
	return pipe
}

// Quality sets the quality value for image formats that support it
func (pipe *Pipeline) Quality(quality int) *Pipeline {
	pipe.bb.Quality = quality
	return pipe
}

// StripMetadata strips ICC profile and metadata from the image
func (pipe *Pipeline) StripMetadata() *Pipeline {
	pipe.bb.StripProfile = true
	pipe.bb.StripMetadata = true
	return pipe
}

// Invert inverts the image color
func (pipe *Pipeline) Invert() *Pipeline {
	pipe.checkImage()
	if err := pipe.image.Invert(); err != nil {
		pipe.err = err
	}
	return pipe
}

// Flip flips the image horizontally or vertically
func (pipe *Pipeline) Flip(direction Direction) *Pipeline {
	pipe.checkImage()
	if err := pipe.image.Flip(direction); err != nil {
		pipe.err = err
	}
	return pipe
}

// GaussBlur applies a gaussian blur to the image
func (pipe *Pipeline) GaussBlur(sigma float64) *Pipeline {
	if err := pipe.image.Gaussblur(sigma); err != nil {
		pipe.err = err
	}
	return pipe
}

// Crop an image, width and height must be equal to or less than image size
func (pipe *Pipeline) Crop(width, height int) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.ResizeStrategy = ResizeStrategyCrop
	pipe.bb.targetWidth = width
	pipe.bb.targetHeight = height
	return pipe.resize()
}

// Stretch an image without maintaining aspect ratio
func (pipe *Pipeline) Stretch(width, height int) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.ResizeStrategy = ResizeStrategyStretch
	pipe.bb.targetWidth = width
	pipe.bb.targetHeight = height
	return pipe.resize()
}

// ScaleWidth scales the image by its width proportionally
func (pipe *Pipeline) ScaleWidth(scale float64) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.targetScaleX = scale
	return pipe.resize()
}

// ScaleHeight scales the height of the image proportionally
func (pipe *Pipeline) ScaleHeight(scale float64) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.targetScaleY = scale
	return pipe.resize()
}

// Scale the image
func (pipe *Pipeline) Scale(scaleX, scaleY float64) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.targetScaleX = scaleX
	pipe.bb.targetScaleY = scaleY
	return pipe.resize()
}

// Reduce the image proportionally
func (pipe *Pipeline) Reduce(scale float64) *Pipeline {
	pipe.checkImage()
	if scale >= 1 {
		panic("scale must be less than 1")
	}
	pipe.bb.clear()
	pipe.bb.targetScaleX = scale
	pipe.bb.targetScaleY = scale
	return pipe.resize()
}

// ResizeWidth resizes the image to the given width, maintaining aspect ratio
func (pipe *Pipeline) ResizeWidth(width int) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.targetScaleX = scale(width, pipe.image.Width())
	return pipe.resize()
}

// ResizeHeight resizes the image to the given height, maintaining aspect ratio
func (pipe *Pipeline) ResizeHeight(height int) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.targetScaleY = scale(height, pipe.image.Height())
	return pipe.resize()
}

// Resize resizes the image to the given width and height
func (pipe *Pipeline) Resize(width, height int) *Pipeline {
	pipe.checkImage()
	pipe.bb.clear()
	pipe.bb.targetWidth = width
	pipe.bb.targetHeight = height
	return pipe.resize()
}

func (pipe *Pipeline) resize() *Pipeline {
	if err := resize(pipe.image, pipe.bb); err != nil {
		pipe.err = err
	}
	return pipe
}

func resize(image *ImageRef, bb *blackboard) error {
	kernel := bb.ReductionSampler

	// Check for the simple scale down cases
	if (bb.targetScaleX > 0 && bb.targetScaleX <= 1) || (bb.targetScaleY > 0 && bb.targetScaleY < 1) {
		scaleX := bb.targetScaleX
		scaleY := bb.targetScaleY
		if scaleX == 0 {
			scaleX = scaleY
		} else if scaleY == 0 {
			scaleY = scaleX
		}
		if scaleX == scaleY {
			return image.Resize(scaleX, InputInt("kernel", int(kernel)))
		}
	}

	if bb.targetWidth == 0 {
		bb.targetWidth = roundFloat(bb.targetScaleX * float64(image.Width()))
	}
	if bb.targetHeight == 0 {
		bb.targetHeight = roundFloat(bb.targetScaleY * float64(image.Height()))
	}

	if bb.targetWidth == 0 || bb.targetHeight == 0 {
		return nil
	}

	shrinkX := scale(image.Width(), bb.targetWidth)
	shrinkY := scale(image.Height(), bb.targetHeight)

	crop := bb.ResizeStrategy == ResizeStrategyCrop

	direction := DirectionHorizontal

	if crop && shrinkX >= shrinkY {
		direction = DirectionVertical
	} else if !crop && shrinkX < shrinkY {
		direction = DirectionVertical
	}

	if direction == DirectionHorizontal {
		shrinkY = shrinkX
	} else {
		shrinkX = shrinkY
	}

	if err := image.Resize(
		1.0/shrinkX,
		InputDouble("vscale", 1.0/shrinkY),
		InputInt("kernel", int(kernel)),
	); err != nil {
		return err
	}

	// If stretching then we're done.
	if bb.ResizeStrategy == ResizeStrategyStretch {
		return nil
	}

	// Crop if necessary
	if crop {
		if bb.targetWidth < image.Width() || bb.targetHeight < image.Height() {
			var left, top int
			width, height := image.Width(), image.Height()
			if bb.targetWidth < image.Width() {
				width = bb.targetWidth
				left = (image.Width() - bb.targetWidth) >> 1
			}
			if bb.targetHeight < image.Height() {
				height = bb.targetHeight
				top = (image.Height() - bb.targetHeight) >> 1
			}
			if err := image.ExtractArea(left, top, width, height); err != nil {
				return err
			}
		}
		return nil
	}

	// Now we might need to embed to match the target dimensions
	if bb.targetWidth > image.Width() || bb.targetHeight > image.Height() {
		var left, top int
		width, height := image.Width(), image.Height()
		if bb.targetWidth > image.Width() {
			width = bb.targetWidth
			left = (bb.targetWidth - image.Width()) >> 1
		}
		if bb.targetHeight > image.Height() {
			height = bb.targetHeight
			top = (bb.targetHeight - image.Height()) >> 1
		}
		if err := image.Embed(left, top, width, height, InputInt("extend", int(bb.PadStrategy))); err != nil {
			return err
		}
	}

	return nil
}

func (pipe *Pipeline) reset() {
	if pipe.image != nil {
		pipe.image.Close()
		pipe.image = nil
	}
	pipe.err = nil
	pipe.bb.clear()
}

func (pipe *Pipeline) checkImage() {
	if pipe.image == nil {
		panic("no image loaded in pipeline")
	}
}

type blackboard struct {
	ResizeParams
	ExportParams
	targetWidth  int
	targetHeight int
	targetScaleX float64
	targetScaleY float64
}

func newBlackboard() *blackboard {
	bb := &blackboard{}
	bb.ResizeStrategy = ResizeStrategyAuto
	bb.CropAnchor = AnchorAuto
	bb.ReductionSampler = KernelLanczos3
	bb.EnlargementInterpolator = InterpolateBicubic
	bb.Quality = 80
	bb.Interpretation = InterpretationSRGB
	bb.clear()
	return bb
}

func (bb *blackboard) clear() {
	bb.targetWidth = 0
	bb.targetHeight = 0
	bb.targetScaleX = 0
	bb.targetScaleY = 0
}

func scale(x, y int) float64 {
	return float64(x) / float64(y)
}

func roundFloat(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}
