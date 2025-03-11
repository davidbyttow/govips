package vips_test

import (
	"bytes"
	"fmt"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidbyttow/govips/v2/vips"
)

func TestBMPLoading(t *testing.T) {
	// Define the minimum ratio of pixels that should have significant variation
	const minSignificantPixelRatio = 0.3 // 30%

	// Check if BMP is supported via the magick loader
	// BMP files are loaded using the magick loader in govips
	if !vips.IsTypeSupported(vips.ImageTypeMagick) {
		t.Skip("Magick loader is not available, skipping BMP tests")
	}

	// List of koala BMP files to test
	koalaFiles := []string{
		"koala.bmp",
		"koala.bmp2",
		"koala.bmp3",
		"koala-rle.bmp",
		"koala-rle.bmp2",
		"koala-rle.bmp3",
	}

	// Get the resources directory
	resourcesDir := filepath.Join("..", "resources")

	// Test each file
	for _, filename := range koalaFiles {
		t.Run(filename, func(t *testing.T) {
			// Construct full path to the file
			filePath := filepath.Join(resourcesDir, filename)

			// Read the file into a buffer
			buf, err := os.ReadFile(filePath)
			require.NoError(t, err, "Failed to read file: %s", filename)

			// Check if the file is a BMP
			imageType := vips.DetermineImageType(buf)
			require.Equal(t, vips.ImageTypeBMP, imageType, "File should be detected as BMP: %s", filename)

			// Load the BMP file using the magick loader
			// The vipsLoadFromBuffer function will automatically use the magick loader for BMP files
			importParams := vips.NewImportParams()
			img, err := vips.LoadImageFromBuffer(buf, importParams)
			require.NoError(t, err, "Failed to load BMP file: %s", filename)
			require.NotNil(t, img, "Image should not be nil")

			// Make sure to close the image when done
			defer img.Close()

			// Check that the image has valid dimensions
			width := img.Width()
			height := img.Height()
			assert.Greater(t, width, 0, "Image width should be greater than 0")
			assert.Greater(t, height, 0, "Image height should be greater than 0")

			// Export as PNG
			pngParams := vips.NewPngExportParams()
			pngBytes, metadata, err := img.ExportPng(pngParams)
			require.NoError(t, err, "Failed to export as PNG: %s", filename)
			require.NotNil(t, pngBytes, "PNG bytes should not be nil")
			require.NotNil(t, metadata, "Metadata should not be nil")

			// Verify the exported PNG has the same dimensions
			assert.Equal(t, width, metadata.Width, "Exported PNG width should match original")
			assert.Equal(t, height, metadata.Height, "Exported PNG height should match original")

			// Load exported PNG bytes directly into standard library image
			pngReader := bytes.NewReader(pngBytes)
			stdImg, err := png.Decode(pngReader)
			require.NoError(t, err, "Failed to decode exported PNG")

			bounds := stdImg.Bounds()
			totalPixels := bounds.Dx() * bounds.Dy()
			pixelValues := make([]float64, totalPixels)

			// Calculate average grayscale value and collect all values
			var sum float64
			idx := 0
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := stdImg.At(x, y).RGBA()
					// Convert from uint32 to uint8 range and average RGB
					gray := float64(uint8(r>>8)+uint8(g>>8)+uint8(b>>8)) / 3.0
					pixelValues[idx] = gray
					sum += gray
					idx++
				}
			}
			mean := sum / float64(totalPixels)

			// Calculate standard deviation
			var variance float64
			for _, v := range pixelValues {
				diff := v - mean
				variance += diff * diff
			}
			variance /= float64(totalPixels)
			stdDev := math.Sqrt(variance)

			// Count pixels more than 1 stddev from mean
			significantPixels := 0
			for _, v := range pixelValues {
				if math.Abs(v-mean) > stdDev {
					significantPixels++
				}
			}

			significantRatio := float64(significantPixels) / float64(totalPixels)
			t.Logf("Image stats: ratio=%.1f%%, avg=%.1f, stddev=%.1f, total_pixels=%d, significant_pixels=%d",
				significantRatio*100, mean, stdDev, totalPixels, significantPixels)
			assert.Greater(t, significantRatio, minSignificantPixelRatio, "Image should have at least 30%% pixels with significant variation (got %.1f%%)", significantRatio*100)

			fmt.Printf("Successfully processed %s: %dx%d\n", filename, width, height)
		})
	}
}
