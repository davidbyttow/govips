package vips

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImage_Animated_GIF(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		nil,
		nil,
		exportGif(NewGifExportParams()))
}

func TestImage_Animated_GIF_ExportNative(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		nil,
		nil,
		nil)
}

func TestImage_Animated_GIF_to_WebP(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		nil,
		nil,
		exportWebp(NewWebpExportParams()))
}

func TestImage_Animated_GIF_Resize(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		func(img *ImageRef) error {
			return img.Resize(2, KernelCubic)
		},
		nil,
		nil)
}

func goldenAnimatedTest(
	t *testing.T,
	path string,
	pages int,
	exec func(img *ImageRef) error,
	validate func(img *ImageRef),
	export func(img *ImageRef) ([]byte, *ImageMetadata, error),
) []byte {
	if exec == nil {
		exec = func(*ImageRef) error { return nil }
	}

	if validate == nil {
		validate = func(*ImageRef) {}
	}

	if export == nil {
		export = func(img *ImageRef) ([]byte, *ImageMetadata, error) { return img.ExportNative() }
	}

	Startup(nil)

	importParams := NewImportParams()
	importParams.NumPages.Set(pages)

	img, err := LoadImageFromFile(path, importParams)
	require.NoError(t, err)

	err = exec(img)
	require.NoError(t, err)

	buf, metadata, err := export(img)
	require.NoError(t, err)

	result, err := NewImageFromBuffer(buf)
	require.NoError(t, err)

	validate(result)

	assertGoldenMatch(t, path, buf, metadata.Format)

	return buf
}
