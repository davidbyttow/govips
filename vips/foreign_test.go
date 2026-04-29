package vips

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DetermineImageType__JPEG(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "jpg-24bit-icc-iec.jpg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJPEG, imageType)
}

func Test_DetermineImageType__HEIF_HEIC(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "heic-24bit-exif.heic")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeHEIF, imageType)
}

func Test_DetermineImageType__PSD(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "psd.example.psd")

	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePSD, imageType)
}

func Test_DetermineImageType__HEIF_MIF1(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "heic-24bit.heic")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeHEIF, imageType)
}

func Test_DetermineImageType__PNG(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "png-24bit+alpha.png")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePNG, imageType)
}

func Test_DetermineImageType__TIFF(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "tif.tif")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeTIFF, imageType)
}

func Test_DetermineImageType__BigTIFF(t *testing.T) {
	for name, buf := range map[string][]byte{
		"little_endian": {0x49, 0x49, 0x2B, 0x00, 0x08, 0x00, 0x00, 0x00, 0, 0, 0, 0},
		"big_endian":    {0x4D, 0x4D, 0x00, 0x2B, 0x00, 0x08, 0x00, 0x00, 0, 0, 0, 0},
	} {
		t.Run(name, func(t *testing.T) {
			imageType := DetermineImageType(buf)
			assert.Equal(t, ImageTypeTIFF, imageType)
		})
	}
}

func Test_DetermineImageType__WEBP(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "webp+alpha.webp")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeWEBP, imageType)
}

func Test_DetermineImageType__SVG(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "svg.svg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeSVG, imageType)
}

func Test_DetermineImageType__SVG_1(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "svg_1.svg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeSVG, imageType)
}

func Test_DetermineImageType__PDF(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "pdf.pdf")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePDF, imageType)
}

func Test_DetermineImageType__PDF_1(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "PDF-2.0-with-offset-start.pdf")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePDF, imageType)
}

func Test_DetermineImageType__BMP(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "bmp.bmp")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeBMP, imageType)
}

func Test_DetermineImageType__AVIF(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "avif-8bit.avif")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeAVIF, imageType)
}

func Test_DetermineImageType__AVIF_Animated(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	// Construct a minimal ftyp box with "avis" brand (Animated AVIF)
	buf := make([]byte, 16)
	copy(buf[4:8], "ftyp")
	copy(buf[8:12], "avis")

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeAVIF, imageType)
}

func Test_DetermineImageType__JP2K(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "jp2k-orientation-6.jp2")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJP2K, imageType)
}

func Test_DetermineImageType__JXL(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "jxl-8bit-grey-icc-dot-gain.jxl")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJXL, imageType)
}

func Test_DetermineImageType__JXL_ISOBMFF(t *testing.T) {
	require.NoError(t, Startup(&Config{}))

	buf, err := os.ReadFile(resources + "jxl-isobmff.jxl")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJXL, imageType)
}

func TestIsPDF(t *testing.T) {
	t.Run("real_pdf", func(t *testing.T) {
		buf, err := os.ReadFile(resources + "PDF-2.0-with-offset-start.pdf")
		require.NoError(t, err)
		assert.True(t, isPDF(buf))
	})

	t.Run("non_pdf", func(t *testing.T) {
		buf, err := os.ReadFile(resources + "jpg-24bit.jpg")
		require.NoError(t, err)
		assert.False(t, isPDF(buf))
	})

	t.Run("small_buffer_with_sig", func(t *testing.T) {
		// Buffer <= 1024 bytes with PDF signature
		buf := make([]byte, 1024)
		copy(buf[1020:], []byte("%PDF"))
		assert.True(t, isPDF(buf))
	})

	t.Run("large_buffer_sig_within_window", func(t *testing.T) {
		// Buffer > 1024 bytes, signature within first 1024 bytes
		buf := make([]byte, 2048)
		copy(buf[100:], []byte("%PDF"))
		assert.True(t, isPDF(buf))
	})

	t.Run("large_buffer_sig_outside_window", func(t *testing.T) {
		// Buffer > 1024 bytes, signature at position 1024 (outside window)
		buf := make([]byte, 2048)
		copy(buf[1024:], []byte("%PDF"))
		assert.False(t, isPDF(buf))
	})
}
