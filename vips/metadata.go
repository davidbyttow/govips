package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

type ImageSize struct {
	Width  int
	Height int
}

type ImageMetadata struct {
	Alpha       bool
	Channels    int
	Colourspace string
	Orientation int
	Profile     bool
	Size        ImageSize
	Space       string
	Type        string
}

func LoadMetadata(image *Image) (*ImageMetadata, error) {
	defer ShutdownThread()

	return &ImageMetadata{
		Size: ImageSize{
			Width:  image.Width(),
			Height: image.Height(),
		},
	}, nil
}
