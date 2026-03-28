package pipeline

import "github.com/davidbyttow/govips/v2/vips"

// Geometric transforms

// Resize returns an operation that resizes the image by scale using the given kernel.
func Resize(scale float64, kernel vips.Kernel) Operation {
	return Operation{name: "Resize", fn: func(img *vips.ImageRef) error {
		return img.Resize(scale, kernel)
	}}
}

// ResizeWithVScale returns an operation that resizes with separate horizontal and vertical scaling.
func ResizeWithVScale(hScale, vScale float64, kernel vips.Kernel) Operation {
	return Operation{name: "ResizeWithVScale", fn: func(img *vips.ImageRef) error {
		return img.ResizeWithVScale(hScale, vScale, kernel)
	}}
}

// Thumbnail returns an operation that resizes to exact dimensions with cropping.
func Thumbnail(width, height int, crop vips.Interesting) Operation {
	return Operation{name: "Thumbnail", fn: func(img *vips.ImageRef) error {
		return img.Thumbnail(width, height, crop)
	}}
}

// ThumbnailWithSize returns an operation that resizes with a size strategy.
func ThumbnailWithSize(width, height int, crop vips.Interesting, size vips.Size) Operation {
	return Operation{name: "ThumbnailWithSize", fn: func(img *vips.ImageRef) error {
		return img.ThumbnailWithSize(width, height, crop, size)
	}}
}

// Rotate returns an operation that rotates the image by multiples of 90 degrees.
func Rotate(angle vips.Angle) Operation {
	return Operation{name: "Rotate", fn: func(img *vips.ImageRef) error {
		return img.Rotate(angle)
	}}
}

// AutoRotate returns an operation that rotates based on EXIF orientation.
func AutoRotate() Operation {
	return Operation{name: "AutoRotate", fn: func(img *vips.ImageRef) error {
		return img.AutoRotate()
	}}
}

// Flip returns an operation that flips the image horizontally or vertically.
func Flip(direction vips.Direction) Operation {
	return Operation{name: "Flip", fn: func(img *vips.ImageRef) error {
		return img.Flip(direction)
	}}
}

// Crop returns an operation that crops the image.
func Crop(left, top, width, height int) Operation {
	return Operation{name: "Crop", fn: func(img *vips.ImageRef) error {
		return img.Crop(left, top, width, height)
	}}
}

// SmartCrop returns an operation that crops based on interesting content.
func SmartCrop(width, height int, interesting vips.Interesting) Operation {
	return Operation{name: "SmartCrop", fn: func(img *vips.ImageRef) error {
		return img.SmartCrop(width, height, interesting)
	}}
}

// ExtractArea returns an operation that extracts a region from the image.
func ExtractArea(left, top, width, height int) Operation {
	return Operation{name: "ExtractArea", fn: func(img *vips.ImageRef) error {
		return img.ExtractArea(left, top, width, height)
	}}
}

// Embed returns an operation that embeds the image in a larger canvas.
func Embed(left, top, width, height int, extend vips.ExtendStrategy) Operation {
	return Operation{name: "Embed", fn: func(img *vips.ImageRef) error {
		return img.Embed(left, top, width, height, extend)
	}}
}

// EmbedBackground returns an operation that embeds with a background color.
func EmbedBackground(left, top, width, height int, bg *vips.Color) Operation {
	return Operation{name: "EmbedBackground", fn: func(img *vips.ImageRef) error {
		return img.EmbedBackground(left, top, width, height, bg)
	}}
}

// EmbedBackgroundRGBA returns an operation that embeds with an RGBA background color.
func EmbedBackgroundRGBA(left, top, width, height int, bg *vips.ColorRGBA) Operation {
	return Operation{name: "EmbedBackgroundRGBA", fn: func(img *vips.ImageRef) error {
		return img.EmbedBackgroundRGBA(left, top, width, height, bg)
	}}
}

// Zoom returns an operation that zooms by pixel repetition.
func Zoom(xFactor, yFactor int) Operation {
	return Operation{name: "Zoom", fn: func(img *vips.ImageRef) error {
		return img.Zoom(xFactor, yFactor)
	}}
}

// Similarity returns an operation that scales, rotates, and offsets in a single step.
func Similarity(scale, angle float64, bg *vips.ColorRGBA, idx, idy, odx, ody float64) Operation {
	return Operation{name: "Similarity", fn: func(img *vips.ImageRef) error {
		return img.Similarity(scale, angle, bg, idx, idy, odx, ody)
	}}
}

// Gravity returns an operation that positions the image using gravity.
func Gravity(gravity vips.Gravity, width, height int) Operation {
	return Operation{name: "Gravity", fn: func(img *vips.ImageRef) error {
		return img.Gravity(gravity, width, height)
	}}
}

// Replicate returns an operation that repeats the image across and down.
func Replicate(across, down int) Operation {
	return Operation{name: "Replicate", fn: func(img *vips.ImageRef) error {
		return img.Replicate(across, down)
	}}
}

// Filters and effects

// GaussianBlur returns an operation that applies a Gaussian blur.
func GaussianBlur(sigma float64) Operation {
	return Operation{name: "GaussianBlur", fn: func(img *vips.ImageRef) error {
		return img.GaussianBlur(sigma)
	}}
}

// Sharpen returns an operation that sharpens the image.
func Sharpen(sigma, x1, m2 float64) Operation {
	return Operation{name: "Sharpen", fn: func(img *vips.ImageRef) error {
		return img.Sharpen(sigma, x1, m2)
	}}
}

// Sobel returns an operation that applies the Sobel edge detector.
func Sobel() Operation {
	return Operation{name: "Sobel", fn: func(img *vips.ImageRef) error {
		return img.Sobel()
	}}
}

// Flatten returns an operation that removes alpha and replaces with a background color.
func Flatten(bg *vips.Color) Operation {
	return Operation{name: "Flatten", fn: func(img *vips.ImageRef) error {
		return img.Flatten(bg)
	}}
}

// Color operations

// ToColorSpace returns an operation that converts to a different color space.
func ToColorSpace(interpretation vips.Interpretation) Operation {
	return Operation{name: "ToColorSpace", fn: func(img *vips.ImageRef) error {
		return img.ToColorSpace(interpretation)
	}}
}

// Modulate returns an operation that modulates brightness, saturation, and hue via LCH.
func Modulate(brightness, saturation, hue float64) Operation {
	return Operation{name: "Modulate", fn: func(img *vips.ImageRef) error {
		return img.Modulate(brightness, saturation, hue)
	}}
}

// ModulateHSV returns an operation that modulates via HSV.
func ModulateHSV(brightness, saturation float64, hue int) Operation {
	return Operation{name: "ModulateHSV", fn: func(img *vips.ImageRef) error {
		return img.ModulateHSV(brightness, saturation, hue)
	}}
}

// Invert returns an operation that inverts the image.
func Invert() Operation {
	return Operation{name: "Invert", fn: func(img *vips.ImageRef) error {
		return img.Invert()
	}}
}

// Gamma returns an operation that adjusts gamma.
func Gamma(gamma float64) Operation {
	return Operation{name: "Gamma", fn: func(img *vips.ImageRef) error {
		return img.Gamma(gamma)
	}}
}

// Recomb returns an operation that recombines bands using a matrix.
func Recomb(matrix [][]float64) Operation {
	return Operation{name: "Recomb", fn: func(img *vips.ImageRef) error {
		return img.Recomb(matrix)
	}}
}

// Arithmetic

// Linear returns an operation that applies a linear transformation (output = input * a + b).
func Linear(a, b []float64) Operation {
	return Operation{name: "Linear", fn: func(img *vips.ImageRef) error {
		return img.Linear(a, b)
	}}
}

// Linear1 returns an operation that applies a single-constant linear transformation.
func Linear1(a, b float64) Operation {
	return Operation{name: "Linear1", fn: func(img *vips.ImageRef) error {
		return img.Linear1(a, b)
	}}
}

// Band operations

// AddAlpha returns an operation that adds an alpha channel.
func AddAlpha() Operation {
	return Operation{name: "AddAlpha", fn: func(img *vips.ImageRef) error {
		return img.AddAlpha()
	}}
}

// ExtractBand returns an operation that extracts bands in place.
func ExtractBand(band, num int) Operation {
	return Operation{name: "ExtractBand", fn: func(img *vips.ImageRef) error {
		return img.ExtractBand(band, num)
	}}
}

// BandJoinConst returns an operation that appends constant bands.
func BandJoinConst(constants []float64) Operation {
	return Operation{name: "BandJoinConst", fn: func(img *vips.ImageRef) error {
		return img.BandJoinConst(constants)
	}}
}

// Cast returns an operation that converts to a target band format.
func Cast(format vips.BandFormat) Operation {
	return Operation{name: "Cast", fn: func(img *vips.ImageRef) error {
		return img.Cast(format)
	}}
}

// ICC profile management

// RemoveICCProfile returns an operation that removes the ICC profile.
func RemoveICCProfile() Operation {
	return Operation{name: "RemoveICCProfile", fn: func(img *vips.ImageRef) error {
		return img.RemoveICCProfile()
	}}
}

// OptimizeICCProfile returns an operation that optimizes the ICC profile.
func OptimizeICCProfile() Operation {
	return Operation{name: "OptimizeICCProfile", fn: func(img *vips.ImageRef) error {
		return img.OptimizeICCProfile()
	}}
}

// TransformICCProfile returns an operation that transforms to a target ICC profile.
func TransformICCProfile(outputProfilePath string) Operation {
	return Operation{name: "TransformICCProfile", fn: func(img *vips.ImageRef) error {
		return img.TransformICCProfile(outputProfilePath)
	}}
}

// Metadata

// RemoveMetadata returns an operation that removes EXIF metadata.
func RemoveMetadata(keep ...string) Operation {
	return Operation{name: "RemoveMetadata", fn: func(img *vips.ImageRef) error {
		return img.RemoveMetadata(keep...)
	}}
}

// RemoveOrientation returns an operation that removes the EXIF orientation tag.
func RemoveOrientation() Operation {
	return Operation{name: "RemoveOrientation", fn: func(img *vips.ImageRef) error {
		return img.RemoveOrientation()
	}}
}

// SetOrientation returns an operation that sets the EXIF orientation.
func SetOrientation(orientation int) Operation {
	return Operation{name: "SetOrientation", fn: func(img *vips.ImageRef) error {
		return img.SetOrientation(orientation)
	}}
}

// Label returns an operation that overlays label text on the image.
func Label(params *vips.LabelParams) Operation {
	return Operation{name: "Label", fn: func(img *vips.ImageRef) error {
		return img.Label(params)
	}}
}
