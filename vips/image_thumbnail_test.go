package vips

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThumbnail_NoCrop(t *testing.T) {
	goldenCreateTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(path string) (*ImageRef, error) {
			return NewThumbnailFromFile(path, 36, 36, InterestingNone)
		},
		func(buf []byte) (*ImageRef, error) {
			return NewThumbnailFromBuffer(buf, 36, 36, InterestingNone)
		},
		nil,
		func(result *ImageRef) {
			assert.Equal(t, 36, result.Width())
			assert.Equal(t, 24, result.Height())
			assert.Equal(t, ImageTypeJPEG, result.Format())
		}, nil)
}

func TestThumbnail_NoUpscale(t *testing.T) {
	goldenCreateTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(path string) (*ImageRef, error) {
			return NewThumbnailWithSizeFromFile(path, 9999, 9999, InterestingNone, SizeDown)
		},
		func(buf []byte) (*ImageRef, error) {
			return NewThumbnailWithSizeFromBuffer(buf, 9999, 9999, InterestingNone, SizeDown)
		},
		nil,
		func(result *ImageRef) {
			assert.Equal(t, 715, result.Width())
			assert.Equal(t, 483, result.Height())
			assert.Equal(t, ImageTypeJPEG, result.Format())
		}, nil)
}

func TestThumbnail_CropCentered(t *testing.T) {
	goldenCreateTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(path string) (*ImageRef, error) {
			return NewThumbnailFromFile(path, 25, 25, InterestingCentre)
		},
		func(buf []byte) (*ImageRef, error) {
			return NewThumbnailFromBuffer(buf, 25, 25, InterestingCentre)
		},
		nil,
		func(result *ImageRef) {
			assert.Equal(t, 25, result.Width())
			assert.Equal(t, 25, result.Height())
			assert.Equal(t, ImageTypeJPEG, result.Format())
		}, nil)
}

func TestThumbnail_PNG_CropCentered(t *testing.T) {
	goldenCreateTest(t, resources+"png-24bit.png",
		func(path string) (*ImageRef, error) {
			return NewThumbnailFromFile(path, 25, 25, InterestingCentre)
		},
		func(buf []byte) (*ImageRef, error) {
			return NewThumbnailFromBuffer(buf, 25, 25, InterestingCentre)
		},
		nil,
		func(result *ImageRef) {
			assert.Equal(t, 25, result.Width())
			assert.Equal(t, 25, result.Height())
			assert.Equal(t, ImageTypePNG, result.Format())
		}, nil)
}

func TestThumbnail_Decode_BMP(t *testing.T) {
	goldenCreateTest(t, resources+"bmp.bmp",
		func(path string) (*ImageRef, error) {
			return NewThumbnailWithSizeFromFile(path, 9999, 9999, InterestingNone, SizeDown)
		},
		func(buf []byte) (*ImageRef, error) {
			return NewThumbnailWithSizeFromBuffer(buf, 9999, 9999, InterestingNone, SizeDown)
		},
		nil,
		func(img *ImageRef) {
			assert.Equal(t, 164, img.Width())
			assert.Equal(t, 211, img.Height())
		}, nil)
}

func TestThumbnail_GIF(t *testing.T) {
	goldenCreateTest(t, resources+"gif-animated.gif",
		func(path string) (*ImageRef, error) {
			return NewThumbnailFromFile(path, 50, 70, InterestingNone)
		},
		func(buf []byte) (*ImageRef, error) {
			return NewThumbnailFromBuffer(buf, 50, 70, InterestingNone)
		},
		func(img *ImageRef) error {
			pages := img.Pages()
			assert.Equal(t, 8, pages)
			assert.Equal(t, img.Format(), ImageTypeGIF)
			return nil
		}, nil, exportGif(NewGifExportParams()))
}

func TestThumbnail_GIF_Animated(t *testing.T) {
	importParams := NewImportParams()
	importParams.NumPages.Set(-1)

	goldenCreateTest(t, resources+"gif-animated.gif",
		func(path string) (*ImageRef, error) {
			return LoadThumbnailFromFile(path, 50, 70, InterestingNone, SizeBoth, importParams)
		},
		func(buf []byte) (*ImageRef, error) {
			return LoadThumbnailFromBuffer(buf, 50, 70, InterestingNone, SizeBoth, importParams)
		},
		nil, nil, exportGif(NewGifExportParams()))
}

func TestThumbnail_GIF_Animated_Force(t *testing.T) {
	importParams := NewImportParams()
	importParams.NumPages.Set(-1)

	goldenCreateTest(t, resources+"gif-animated.gif",
		func(path string) (*ImageRef, error) {
			return LoadThumbnailFromFile(path, 50, 100, InterestingNone, SizeForce, importParams)
		},
		func(buf []byte) (*ImageRef, error) {
			return LoadThumbnailFromBuffer(buf, 50, 100, InterestingNone, SizeForce, importParams)
		},
		nil, nil, exportGif(NewGifExportParams()))
}

func TestThumbnail_GIF_ExportNative(t *testing.T) {
	importParams := NewImportParams()
	importParams.NumPages.Set(3)

	goldenCreateTest(t, resources+"gif-animated.gif",
		func(path string) (*ImageRef, error) {
			return LoadThumbnailFromFile(path, 50, 70, InterestingNone, SizeBoth, importParams)
		},
		func(buf []byte) (*ImageRef, error) {
			return LoadThumbnailFromBuffer(buf, 50, 70, InterestingNone, SizeBoth, importParams)
		},
		nil, nil, nil)
}

func TestThumbnail_GIF_ExportWebP(t *testing.T) {
	importParams := NewImportParams()
	importParams.NumPages.Set(3)

	goldenCreateTest(t, resources+"gif-animated.gif",
		func(path string) (*ImageRef, error) {
			return LoadThumbnailFromFile(path, 50, 70, InterestingNone, SizeBoth, importParams)
		},
		func(buf []byte) (*ImageRef, error) {
			return LoadThumbnailFromBuffer(buf, 50, 70, InterestingNone, SizeBoth, importParams)
		},
		nil, nil, exportWebp(NewWebpExportParams()))
}
