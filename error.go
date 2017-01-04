package govips

import "errors"

var (
	// ErrUnsupportedImageFormat when image type is unsupported
	ErrUnsupportedImageFormat = errors.New("UnsupportedImageFormat")

	// ErrInvalidInterpolator when interpolator is invalid
	ErrInvalidInterpolator = errors.New("Invalid interpolator")
)
