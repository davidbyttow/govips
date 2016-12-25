package gimage

import (
	"math"

	"github.com/simplethingsllc/gimage/vips"
)

func Process(buf []byte, options *Options) ([]byte, error) {
	defer vips.FinalizeRequest()

	image, err := vips.NewImage(buf)
	if err != nil {
		return nil, err
	}

	image, err = process(image, options)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func process(image *vips.Image, options *Options) (*vips.Image, error) {

	// TODO(d): Rotation

	imageWidth := float64(image.Width())
	imageHeight := float64(image.Height())

	xFactor := 1.0
	yFactor := 1.0
	desiredWidth := float64(options.Width)
	desiredHeight := float64(options.Height)
	if desiredWidth > 0 && desiredHeight > 0 {
		xFactor = imageWidth / desiredWidth
		yFactor = imageHeight / desiredHeight
		switch options.CanvasStrategy {
		case CanvasStrategyCrop:
			fallthrough
		case CanvasStrategyEmbed:
			crop := options.CanvasStrategy == CanvasStrategyCrop
			if (crop && xFactor < yFactor) || (!crop && xFactor > yFactor) {
				desiredHeight = Round(imageHeight / xFactor)
				yFactor = xFactor
			} else {
				desiredWidth = Round(imageWidth / yFactor)
				xFactor = yFactor
			}
		case CanvasStrategyMax:
			fallthrough
		case CanvasStrategyMin:
			max := options.CanvasStrategy == CanvasStrategyMax
			if (max && xFactor > yFactor) || (!max && xFactor < yFactor) {
				desiredHeight = Round(imageHeight / xFactor)
				options.Height = int(desiredHeight)
				yFactor = xFactor
			} else {
				desiredWidth = Round(imageWidth / yFactor)
				options.Width = int(desiredWidth)
				xFactor = yFactor
			}
		case CanvasStrategyIgnoreAspect:
			// Nothing to do unless there's a rotation
		}
	} else if desiredWidth > 0 {
		xFactor = imageWidth / desiredWidth
		if options.CanvasStrategy == CanvasStrategyIgnoreAspect {
			desiredHeight = imageHeight
			options.Height = int(desiredHeight)
		} else {
			yFactor = xFactor
			desiredHeight = Round(imageHeight / yFactor)
			options.Height = int(desiredHeight)
		}
	} else if desiredHeight > 0 {
		yFactor = imageHeight / desiredHeight
		if options.CanvasStrategy == CanvasStrategyIgnoreAspect {
			desiredWidth = imageWidth
			options.Width = int(desiredWidth)
		} else {
			xFactor = yFactor
			desiredWidth = Round(imageWidth / xFactor)
			options.Width = int(desiredWidth)
		}
	} else {
		options.Width = int(imageWidth)
		options.Height = int(imageHeight)
	}

	xShrink := int(math.Max(1.0, math.Floor(xFactor)))
	yShrink := int(math.Max(1.0, math.Floor(yFactor)))

	// xResidual := float64(xShrink) / xFactor
	// yResidual := float64(yShrink) / yFactor

	// Optionally prevent enlargement

	imageType := image.Type()

	hasGammaAdjustment := options.Gamma > 0
	canShrinkOnLoad := (imageType == vips.ImageTypeJpeg || imageType == vips.ImageTypeWebp) &&
		!hasGammaAdjustment

	shrinkFactor := 1
	if canShrinkOnLoad && xShrink == yShrink && xShrink >= 2 {
		if xShrink >= 8 {
			shrinkFactor = 8
		} else if xShrink >= 4 {
			shrinkFactor = 4
		} else if xShrink >= 2 {
			shrinkFactor = 2
		}
		if shrinkFactor > 1 {
			xFactor /= float64(shrinkFactor)
			yFactor = xFactor
		}
	}

	// Reload the image with a shrink factor on load
	if shrinkFactor > 1 {
		buf := CopyBuffer(image.SourceBytes())
		var err error
		if imageType == vips.ImageTypeJpeg {
			image, err = vips.NewJpegImage(buf, shrinkFactor)
		} else {
			image, err = vips.NewWebpImage(buf, shrinkFactor)
		}
		if err != nil {
			return image, err
		}
		shrunkWidth := image.Width()
		shrunkHeight := image.Height()
		xFactor = float64(shrunkWidth) / desiredWidth
		yFactor = float64(shrunkHeight) / desiredHeight
		xShrink = int(math.Max(1.0, math.Floor(xFactor)))
		yShrink = int(math.Max(1.0, math.Floor(yFactor)))
		// xResidual = float64(xShrink) / xFactor
		// yResidual = float64(yShrink) / yFactor
	}

	// TODO(d): Remove alpha channel?

	// TODO(d): Negate image if needed

	// TODO(d): Gamma darkening

	// TODO(d): Greyscale

	// TODO(d): Overlay setup

	// shrink := xShrink > 1 || yShrink > 1
	// reduce := xResidual != 1.0 || yResidual != 1.0
	// blur := options.GaussianBlur.Sigma != 0.0
	// sharpen := options.Sharpen.Sigma != 0.0

	// TODO(d): Premultiply alpha if needed

	// if shrink {
	// 	if yShrink > 1 {
	// 		image = image.ShrinkV(yShrink)
	// 	}
	// 	if xShrink > 1 {
	// 		image = image.ShrinkH(xShrink)
	// 	}
	// 	shrunkWidth := image.Width()
	// 	shrunkHeight := image.Height
	// }

	return image, nil
}
