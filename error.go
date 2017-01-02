package govips

import "errors"

var (
	ErrUnsupportedImageFormat = errors.New("UnsupportedImageFormat")
	ErrInvalidInterpolator    = errors.New("Invalid interpolator")
)
