package vips

import (
	"bytes"
	"image"
	"testing"

	"golang.org/x/image/bmp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: The JPEG spec requires some minimal exif data including exif-ifd0-Orientation.
// libvips always adds these fields back but they should not be a privacy concern.
// HEIC images require the same fields and behave the same way in libvips.
func TestImage_RemoveMetadata_Removes_Exif(t *testing.T) {
	skipIfHeifSaveUnsupported(t)
	var initialEXIFCount int
	goldenTest(t, resources+"heic-24bit-exif.heic",
		func(img *ImageRef) error {
			exifData := img.GetExif()
			initialEXIFCount = len(exifData)
			assert.Greater(t, initialEXIFCount, 0)
			return img.RemoveMetadata()
		},
		func(img *ImageRef) {
			exifData := img.GetExif()
			finalEXIFCount := len(exifData)
			assert.Less(t, finalEXIFCount, initialEXIFCount)
		}, nil)
}

func TestImage_SetExifField(t *testing.T) {
	skipIfHeifSaveUnsupported(t)
	var originalExifValue string
	goldenTest(t, resources+"heic-24bit-exif.heic",
		func(img *ImageRef) error {
			originalExifValue = img.GetString("exif-ifd0-Model")
			assert.NotEqual(t, originalExifValue, "iPhone (iPhone, ASCII, 7 components, 7 bytes)")
			img.SetString("exif-ifd0-Model", "iPhone (iPhone, ASCII, 7 components, 7 bytes)")
			updatedExifValue := img.GetString("exif-ifd0-Model")
			assert.Equal(t, updatedExifValue, "iPhone (iPhone, ASCII, 7 components, 7 bytes)")
			return nil
		},
		func(img *ImageRef) {
			updatedExifValue := img.GetString("exif-ifd0-Model")
			assert.Equal(t, updatedExifValue, "iPhone (iPhone, ASCII, 7 components, 7 bytes)")
		}, nil)
}

func TestImage_AutoRotate_6__heic_to_jpg(t *testing.T) {
	goldenTest(t, resources+"heic-orientation-6.heic",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.Orientation())
		}, exportJpeg(nil),
	)
}

func TestImage_Export_AVIF_8_Bit(t *testing.T) {
	skipIfHeifSaveUnsupported(t)
	avifExportParams := NewAvifExportParams()
	goldenTest(t, resources+"avif-8bit.avif",
		func(img *ImageRef) error {
			return nil
		},
		func(result *ImageRef) {
		}, exportAvif(avifExportParams),
	)
}

func TestImage_TIF_16_Bit_To_AVIF_12_Bit(t *testing.T) {
	skipIfHeifSaveUnsupported(t)
	avifExportParams := NewAvifExportParams()
	avifExportParams.Bitdepth = 12
	goldenTest(t, resources+"tif-16bit.tif",
		func(img *ImageRef) error {
			// TIFF images don't use regular exif fields -- they iptc and/or xmp instead.
			fields := img.GetFields()
			assert.Greater(t, len(fields), 0)
			xmpData := img.GetBlob("xmp-data")
			assert.Greater(t, len(xmpData), 0)
			return nil
		},
		func(result *ImageRef) {
		}, exportAvif(avifExportParams),
	)
}

func TestImage_PSD_To_PNG(t *testing.T) {
	goldenTest(t, resources+"psd.example.psd",
		func(img *ImageRef) error {
			fields := img.GetFields()
			assert.Greater(t, len(fields), 0)
			xmpData := img.GetBlob("xmp-data")
			assert.Greater(t, len(xmpData), 0)
			return nil
		},
		nil,
		exportPng(NewPngExportParams()),
	)
}

func TestImage_PSD_To_JPEG(t *testing.T) {
	goldenTest(t, resources+"psd.example.psd",
		func(img *ImageRef) error {
			fields := img.GetFields()
			assert.Greater(t, len(fields), 0)
			xmpData := img.GetBlob("xmp-data")
			assert.Greater(t, len(xmpData), 0)
			return nil
		},
		nil,
		exportJpeg(NewJpegExportParams()),
	)
}

func TestImage_Decode_BMP(t *testing.T) {
	goldenTest(t, resources+"bmp.bmp", func(img *ImageRef) error {
		goImg, err := img.ToImage(nil)
		assert.NoError(t, err)

		buf := new(bytes.Buffer)
		err = bmp.Encode(buf, goImg)
		assert.Nil(t, err)

		config, format, err := image.DecodeConfig(buf)
		assert.Nil(t, err)
		assert.Equal(t, "bmp", format)
		assert.NotNil(t, config)
		assert.True(t, config.Height > 0)
		assert.True(t, config.Width > 0)
		return nil
	}, nil, nil)
}

func TestImage_Tiff(t *testing.T) {
	goldenTest(t, resources+"tif.tif", func(img *ImageRef) error {
		return img.OptimizeICCProfile()
	}, nil, nil)
}

func TestImage_Black(t *testing.T) {
	require.NoError(t, Startup(nil))
	i, err := Black(10, 20)
	require.NoError(t, err)
	buf, metadata, err := i.ExportNative()
	require.NoError(t, err)

	assertGoldenMatch(t, resources+"jpg-24bit.jpg", buf, metadata.Format)
}

func TestNewTransparentCanvas(t *testing.T) {
	require.NoError(t, Startup(nil))

	ref, err := NewTransparentCanvas(200, 100)
	require.NoError(t, err)
	defer ref.Close()

	assert.Equal(t, 200, ref.Width())
	assert.Equal(t, 100, ref.Height())
	assert.Equal(t, 4, ref.Bands())
	assert.True(t, ref.HasAlpha())
	assert.Equal(t, InterpretationSRGB, ref.Interpretation())

	pixel, err := ref.GetPoint(0, 0)
	require.NoError(t, err)
	assert.Equal(t, []float64{0, 0, 0, 0}, pixel)

	pixel, err = ref.GetPoint(199, 99)
	require.NoError(t, err)
	assert.Equal(t, []float64{0, 0, 0, 0}, pixel)
}

func TestNewTransparentCanvas_Composite(t *testing.T) {
	require.NoError(t, Startup(nil))

	canvas, err := NewTransparentCanvas(200, 200)
	require.NoError(t, err)
	defer canvas.Close()

	overlay, err := Black(50, 50)
	require.NoError(t, err)
	defer overlay.Close()
	require.NoError(t, overlay.ToColorSpace(InterpretationSRGB))

	err = canvas.Composite(overlay, BlendModeOver, 10, 10)
	require.NoError(t, err)

	assert.Equal(t, 200, canvas.Width())
	assert.Equal(t, 200, canvas.Height())
}

func TestImage_Grey(t *testing.T) {
	require.NoError(t, Startup(nil))

	ref, err := Grey(256, 1, true)
	require.NoError(t, err)
	defer ref.Close()

	assert.Equal(t, 256, ref.Width())
	assert.Equal(t, 1, ref.Height())
	assert.Equal(t, 1, ref.Bands())

	// First pixel should be near 0, last pixel should be near 255
	p0, err := ref.GetPoint(0, 0)
	require.NoError(t, err)
	assert.InDelta(t, 0, p0[0], 2)

	p255, err := ref.GetPoint(255, 0)
	require.NoError(t, err)
	assert.InDelta(t, 255, p255[0], 2)
}

func TestImage_Grey_Gradient_Composite(t *testing.T) {
	// Demonstrates the gradient overlay use case from issue #287:
	// Create a vertical gradient (black->transparent) and composite over an image
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			// Create a gradient ramp matching image dimensions
			gradient, err := Grey(img.Width(), img.Height(), true)
			if err != nil {
				return err
			}
			defer gradient.Close()

			// Rotate 90 degrees to make it vertical (top=black, bottom=white)
			if err := gradient.Rotate(Angle90); err != nil {
				return err
			}

			// Use the gradient as an alpha channel on a black overlay:
			// black image + gradient alpha = black-to-transparent overlay
			overlay, err := Black(img.Width(), img.Height())
			if err != nil {
				return err
			}
			defer overlay.Close()

			if err := overlay.ToColorSpace(InterpretationSRGB); err != nil {
				return err
			}

			// BandJoin the gradient as the alpha channel
			if err := overlay.BandJoin(gradient); err != nil {
				return err
			}

			// Composite the gradient overlay on top
			return img.Composite(overlay, BlendModeOver, 0, 0)
		},
		nil,
		exportPng(NewPngExportParams()),
	)
}

func TestPDF_WithOffsetStart(t *testing.T) {
	goldenTest(t, resources+"PDF-2.0-with-offset-start.pdf",
		nil, func(img *ImageRef) {
			assert.Equal(t, 612, img.Width())
			assert.Equal(t, 396, img.Height())
		}, nil)
}
