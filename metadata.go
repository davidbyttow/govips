package govips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

// Represents the size of an image
type ImageSize struct {
	Width  int
	Height int
}

// Represents metadata for an image
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

// Loads metadata for a given image
func LoadMetadata(image *Image) (*ImageMetadata, error) {
	defer ShutdownThread()

	return &ImageMetadata{
		Size: ImageSize{
			Width:  image.Width(),
			Height: image.Height(),
		},
	}, nil
}
