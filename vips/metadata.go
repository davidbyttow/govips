package vips

// #cgo pkg-config: vips
// #include "bridge.h"
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

func Metadata(image Image) (*ImageMetadata, error) {
	defer ShutdownThread()
	return &ImageMetadata{}, nil
}
