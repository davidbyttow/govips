package vips

import (
	"io/ioutil"
	"math"
)

type Pipeline struct {
	bb    *blackboard
	image *ImageRef
	err   error
}

func NewPipeline() *Pipeline {
	startupIfNeeded()
	return &Pipeline{newBlackboard(), nil, nil}
}

func (pipe *Pipeline) LoadFile(file string) *Pipeline {
	pipe.reset()
	pipe.image, pipe.err = NewImageFromFile(file)
	return pipe
}

func (pipe *Pipeline) LoadBuffer(buf []byte) *Pipeline {
	pipe.reset()
	pipe.image, pipe.err = NewImageFromBuffer(buf)
	return pipe
}

func (pipe *Pipeline) Output() ([]byte, error) {
	pipe.checkImage()
	return pipe.Export(nil)
}

func (pipe *Pipeline) OutputFile(file string) error {
	pipe.checkImage()
	buf, err := pipe.Export(nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, buf, 0644)
}

func (pipe *Pipeline) Export(params *ExportParams) ([]byte, error) {
	defer ShutdownThread()
	defer pipe.reset()

	if pipe.err != nil {
		return nil, pipe.err
	}
	if params == nil {
		params = &ExportParams{}
	}
	return vipsExportBuffer(pipe.image.Image(), params)
}

func (pipe *Pipeline) Error() error {
	return pipe.err
}

func (pipe *Pipeline) Zoom(x, y int) *Pipeline {
	pipe.checkImage()
	if err := pipe.image.Zoom(x, y); err != nil {
		pipe.err = err
	}
	return pipe
}

func (pipe *Pipeline) Invert() *Pipeline {
	pipe.checkImage()
	if err := pipe.image.Invert(); err != nil {
		pipe.err = err
	}
	return pipe
}

func (pipe *Pipeline) Flip(direction Direction) *Pipeline {
	pipe.checkImage()
	if err := pipe.image.Flip(direction); err != nil {
		pipe.err = err
	}
	return pipe
}

func (pipe *Pipeline) GaussBlur(sigma float64) *Pipeline {
	if err := pipe.image.Gaussblur(sigma); err != nil {
		pipe.err = err
	}
	return pipe
}

func (pipe *Pipeline) Crop(width, height int, anchor Anchor) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.ResizeStrategy = ResizeStrategyCrop
	pipe.bb.CropAnchor = anchor
	pipe.bb.targetWidth = width
	pipe.bb.targetHeight = height
	return pipe.resize()
}

func (pipe *Pipeline) Embed(width, height int, strategy Extend) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.ResizeStrategy = ResizeStrategyEmbed
	pipe.bb.EmbedStrategy = strategy
	pipe.bb.targetWidth = width
	pipe.bb.targetHeight = height
	return pipe.resize()
}

func (pipe *Pipeline) Stretch(width, height int) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.ResizeStrategy = ResizeStrategyStretch
	pipe.bb.targetWidth = width
	pipe.bb.targetHeight = height
	return pipe.resize()
}

func (pipe *Pipeline) ScaleWidth(scale float64) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.targetScaleX = scale
	return pipe.resize()
}

func (pipe *Pipeline) ScaleHeight(scale float64) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.targetScaleY = scale
	return pipe.resize()
}

func (pipe *Pipeline) Scale(scaleX, scaleY float64) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.targetScaleX = scaleX
	pipe.bb.targetScaleY = scaleY
	return pipe.resize()
}

func (pipe *Pipeline) Reduce(scale float64) *Pipeline {
	pipe.checkImage()
	if scale >= 1 {
		panic("scale must be less than 1")
	}
	pipe.bb.Clear()
	pipe.bb.targetScaleX = scale
	pipe.bb.targetScaleY = scale
	return pipe.resize()
}

func (pipe *Pipeline) ResizeWidth(width int) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.targetScaleX = scale(width, pipe.image.Width())
	return pipe.resize()
}

func (pipe *Pipeline) ResizeHeight(height int) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
	pipe.bb.targetScaleY = scale(height, pipe.image.Height())
	return pipe.resize()
}

type ResizeParams struct {
	EmbedStrategy           Extend
	ResizeStrategy          ResizeStrategy
	CropAnchor              Anchor
	ReductionSampler        Kernel
	EnlargementInterpolator Interpolator
}

func (pipe *Pipeline) Resize(width, height int, params ResizeParams) *Pipeline {
	pipe.checkImage()
	pipe.bb.Clear()
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
	if kernel == 0 {
		kernel = KernelLanczos3
	}

	interpolator := bb.EnlargementInterpolator
	if interpolator == "" {
		interpolator = InterpolateBicubic
	}

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
		if err := image.Embed(left, top, width, height, InputInt("extend", int(bb.EmbedStrategy))); err != nil {
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
	pipe.bb.Clear()
}

func (pipe *Pipeline) checkImage() {
	if pipe.image == nil {
		panic("no image loaded in pipeline")
	}
}

type blackboard struct {
	ResizeParams
	targetWidth  int
	targetHeight int
	targetScaleX float64
	targetScaleY float64
}

func newBlackboard() *blackboard {
	bb := &blackboard{}
	bb.Clear()
	return bb
}

func (bb *blackboard) Clear() {
	bb.ResizeStrategy = ResizeStrategyAuto
	bb.CropAnchor = AnchorAuto
	bb.ReductionSampler = KernelLanczos3
	bb.EnlargementInterpolator = InterpolateBicubic
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
