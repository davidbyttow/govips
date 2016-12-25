package gimage

type CanvasStrategy int

const (
	CanvasStrategyCrop CanvasStrategy = iota
	CanvasStrategyEmbed
	CanvasStrategyMax
	CanvasStrategyMin
	CanvasStrategyIgnoreAspect
)

type Gravity int

const (
	GravityCenter Gravity = iota
	GravityN
	GravityE
	GravityS
	GravityW
	GravithNE
	GravitySE
	GravitySW
	GravityNW
)

type Kernel int

const (
	KernelAuto Kernel = iota
	KernelNearest
	KernelLinear
	KernelCubic
	KernelLanczos2
	KernelLanczos3
)

type CropStrategy int

const (
	CropStrategyEntropy CropStrategy = iota
	CropStrategyAttention
)

type Interpolator int

const (
	InterpolatorAuto Interpolator = iota
	InterpolatorNearest
	InterpolatorBicubic
	InterpolatorNohalo
	InterpolatorLocallyBoundedBicubic
	InterpolatorVertexSplitQuadraticBSpline
)

type GaussianBlur struct {
	Sigma float64
}

type Sharpen struct {
	Sigma float64
}

type Options struct {
	CanvasStrategy CanvasStrategy
	CenterSampling bool
	Gamma          float64
	GaussianBlur   GaussianBlur
	Gravity        Gravity
	Height         int
	Interpolator   Interpolator
	Kernel         Kernel
	Sharpen        Sharpen
	Width          int
}
