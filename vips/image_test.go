package vips

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	Shutdown()
	os.Exit(ret)
}

func TestImageRef_WebP(t *testing.T) {
	Startup(nil)

	srcBytes, err := os.ReadFile(resources + "webp+alpha.webp")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, _, err = img.ExportWebp(nil)
	assert.NoError(t, err)
}

func TestImageRef_WebP__ReducedEffort(t *testing.T) {
	Startup(nil)

	srcBytes, err := os.ReadFile(resources + "webp+alpha.webp")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)

	params := NewWebpExportParams()
	params.ReductionEffort = 2
	_, _, err = img.ExportWebp(params)
	assert.NoError(t, err)
}

func TestImageRef_WebP__NearLossless(t *testing.T) {
	Startup(nil)

	srcBytes, err := os.ReadFile(resources + "webp+alpha.webp")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)

	params := NewWebpExportParams()
	params.NearLossless = true
	_, _, err = img.ExportWebp(params)
	assert.NoError(t, err)
}

func TestImageRef_PNG(t *testing.T) {
	Startup(nil)

	srcBytes, err := os.ReadFile(resources + "png-24bit.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)

	// check random access by encoding twice
	_, _, err = img.ExportNative()
	assert.NoError(t, err)
	_, _, err = img.ExportNative()
	assert.NoError(t, err)
}

func TestImageRef_HEIF(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "heic-24bit-exif.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.ExportNative()
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, metadata.Format)
}

func TestImageRef_HEIF_MIF1(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "heic-24bit.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.ExportNative()
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, metadata.Format)
}

func TestImageRef_HEIF_ftypmsf1(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "heic-ftypmsf1.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.ExportNative()
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, metadata.Format)
}

func TestImageRef_BMP__ImplicitConversionToPNG(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "bmp.bmp")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	exported, metadata, err := img.ExportNative()
	assert.NoError(t, err)
	assert.Equal(t, ImageTypePNG, metadata.Format)
	assert.Equal(t, ImageTypeBMP, img.OriginalFormat())
	assert.NotNil(t, exported)
}

func TestImageRef_SVG(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "svg.svg")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, ImageTypeSVG, img.Metadata().Format)
}

func TestImageRef_SVG_1(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "svg_1.svg")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, ImageTypeSVG, img.Metadata().Format)
}

func TestImageRef_SVG_2(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "svg_2.svg")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, ImageTypeSVG, img.Metadata().Format)
}

func TestImageRef_OverSizedMetadata(t *testing.T) {
	Startup(nil)

	srcBytes, err := os.ReadFile(resources + "png-bad-metadata.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	assert.NoError(t, err)
	assert.NotNil(t, img)
}

func TestImageRef_Resize__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.Resize(-1, KernelLanczos3)
	require.Error(t, err)
}

func TestImageRef_ExtractArea__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.ExtractArea(1, 2, 10000, 4)
	require.Error(t, err)
}

func TestImageRef_HasAlpha__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit+alpha.png")
	require.NoError(t, err)

	assert.True(t, img.HasAlpha())
}

func TestImageRef_HasAlpha__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	assert.False(t, img.HasAlpha())
}

func TestImageRef_AddAlpha(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	require.NotNil(t, img)

	err = img.AddAlpha()
	assert.NoError(t, err)
	assert.True(t, img.HasAlpha(), "has alpha")

	_, _, err = img.ExportNative()
	assert.NoError(t, err)
}

func TestImageRef_AddAlpha__Idempotent(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit+alpha.png")
	require.NoError(t, err)
	require.NotNil(t, img)

	err = img.AddAlpha()
	assert.NoError(t, err)

	assert.True(t, img.HasAlpha(), "has alpha")
	_, _, err = img.ExportNative()
	assert.NoError(t, err)
}

func TestImageRef_HasProfile__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.True(t, img.HasProfile())
}

func TestImageRef_HasIPTC__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.True(t, img.HasIPTC())
}

func TestImageRef_HasIPTC__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.False(t, img.HasIPTC())
}

func TestImageRef_HasProfile__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	assert.False(t, img.HasProfile())
}

func TestImageRef_GetOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)

	assert.Equal(t, 6, image.Orientation())
}

func TestImageRef_GetOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	assert.Equal(t, 0, image.Orientation())
}

func TestImageRef_SetOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)

	err = image.SetOrientation(5)
	require.NoError(t, err)

	assert.Equal(t, 5, image.Orientation())
}

func TestImageRef_SetOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.SetOrientation(5)
	require.NoError(t, err)

	assert.Equal(t, 5, image.Orientation())
}

func TestImageRef_RemoveOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)

	err = image.RemoveOrientation()
	require.NoError(t, err)

	assert.Equal(t, 0, image.Orientation())
}

func TestImageRef_RemoveOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.RemoveOrientation()
	require.NoError(t, err)

	assert.Equal(t, 0, image.Orientation())
}

func TestImageRef_RemoveMetadata__RetainsProfile(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)

	require.True(t, image.HasIPTC())

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.False(t, image.HasIPTC())
	assert.True(t, image.HasICCProfile())
}

func TestImageRef_RemoveMetadata__RetainsOrientation(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-5.jpg")
	require.NoError(t, err)

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.Equal(t, 5, image.Orientation())
}

func TestImageRef_RemoveMetadata__RetainsNPages(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "gif-animated.gif")
	require.NoError(t, err)

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.Equal(t, 8, image.Pages())
}

func TestImageRef_RemoveMetadata__RetainsPageHeight(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "gif-animated.gif")
	require.NoError(t, err)

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.Equal(t, 128, image.PageHeight())
}

// Known issue: libvips does not write EXIF into WebP:
// https://github.com/libvips/libvips/pull/1745
func TestImageRef_RemoveMetadata__RetainsOrientation__WebP(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "webp-orientation-6.webp")
	require.NoError(t, err)

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.Equal(t, 6, image.Orientation())
}

func TestImageRef_RemoveICCProfile(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)

	require.True(t, image.HasIPTC())

	err = image.RemoveICCProfile()
	require.NoError(t, err)

	assert.False(t, image.HasICCProfile())
	assert.True(t, image.HasIPTC())
}

func TestImageRef_TransformICCProfile(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)

	require.True(t, image.HasIPTC())
	require.True(t, image.HasICCProfile())

	err = image.TransformICCProfile(SRGBIEC6196621ICCProfilePath)
	require.NoError(t, err)

	assert.True(t, image.HasIPTC())
	assert.True(t, image.HasICCProfile())
}

func TestImageRef_TransformICCProfileWithFallback(t *testing.T) {
	Startup(nil)

	t.Run("source with ICC", func(t *testing.T) {
		image, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
		require.NoError(t, err)

		require.True(t, image.HasIPTC())
		require.True(t, image.HasICCProfile())

		err = image.TransformICCProfileWithFallback(SRGBIEC6196621ICCProfilePath, SRGBV2MicroICCProfilePath)
		require.NoError(t, err)

		assert.True(t, image.HasIPTC())
		assert.True(t, image.HasICCProfile())
	})

	t.Run("source without ICC", func(t *testing.T) {
		image, err := NewImageFromFile(resources + "jpg-24bit.jpg")
		require.NoError(t, err)

		require.False(t, image.HasICCProfile())

		err = image.TransformICCProfileWithFallback(SRGBIEC6196621ICCProfilePath, SRGBV2MicroICCProfilePath)
		require.NoError(t, err)

		assert.True(t, image.HasICCProfile())
	})
}

func TestImageRef_Close(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)

	image.Close()
	assert.Nil(t, image.image)

	PrintObjectReport("Final")
}

func TestImageRef_Close__AlreadyClosed(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)

	go image.Close()
	go image.Close()
	go image.Close()
	go image.Close()
	defer image.Close()
	image.Close()

	assert.Nil(t, image.image)
	runtime.GC()
}

func TestImageRef_NotImage(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "txt.txt")
	require.Error(t, err)
	require.Nil(t, image)
}

func TestImageRef_Label(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	lp := &LabelParams{Text: "Text label"}

	err = image.Label(lp)
	require.NoError(t, err)
}

func TestImageRef_Composite(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	imageOverlay, err := NewImageFromFile(resources + "png-8bit+alpha.png")
	require.NoError(t, err)

	err = image.Composite(imageOverlay, BlendModeXOR, 10, 20)
	require.NoError(t, err)
}

func TestImageRef_Insert(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	imageOverlay, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.Insert(imageOverlay, 100, 200, false, nil)
	require.NoError(t, err)
}

func TestImageRef_Join(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	joinImage, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)
	width := image.Width() + joinImage.Width()
	height := joinImage.Height() // join appears to use the second image's height

	err = image.Join(joinImage, DirectionHorizontal)
	require.NoError(t, err)

	assert.True(t, width == image.Width(), "Join image width is incorrect: %d != %d", width, image.Width())
	assert.True(t, height == image.Height(), "Join image height is incorrect: %d != %d", height, image.Height())
}

func TestImageRef_ArrayJoin(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	joinImage1, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	joinImage2, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	joinImage3, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	joinImage4, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	images := []*ImageRef{image, joinImage1, joinImage2, joinImage3, joinImage4}
	width := image.Width() * 2 // arrayjoin appears to size based on the image's width and height
	height := image.Height() * 3

	err = image.ArrayJoin(images, 2)
	require.NoError(t, err)

	assert.True(t, width == image.Width(), "ArrayJoin image width is incorrect: %d != %d", width, image.Width())
	assert.True(t, height == image.Height(), "ArrayJoin image height is incorrect: %d != %d", height, image.Height())
}

func TestImageRef_Mapim(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	index, err := NewImageFromFile(resources + "png-8bit+alpha.png")
	require.NoError(t, err)

	_ = index.ExtractBand(0, 2)
	require.NoError(t, err)

	err = image.Mapim(index)
	require.NoError(t, err)
}

func TestImageRef_Mapim__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	index, err := NewImageFromFile(resources + "png-8bit+alpha.png")
	require.NoError(t, err)

	err = image.Mapim(index)
	assert.Error(t, err)
}

func TestImageRef_Maplut(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	lut, err := XYZ(1, 1)
	require.NoError(t, err)

	_ = image.ExtractBand(0, 2)
	require.NoError(t, err)

	err = image.Maplut(lut)
	require.NoError(t, err)
}

func TestImageRef_Maplut_Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	lut, err := XYZ(1, 1)
	require.NoError(t, err)

	err = image.Maplut(lut)
	assert.Error(t, err)
}

func TestImageRef_CompositeMulti(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	sources := []string{"png-8bit+alpha.png", "png-24bit+alpha.png"}
	images := make([]*ImageComposite, len(sources))
	for i, uri := range sources {
		image, err := NewImageFromFile(resources + uri)
		require.NoError(t, err)

		// add offset test
		images[i] = &ImageComposite{image, BlendModeOver, (i + 1) * 20, (i + 2) * 20}
	}

	err = image.CompositeMulti(images)
	require.NoError(t, err)

	_, _, err = image.ExportNative()
	require.NoError(t, err)
}

func TestImageRef_Recomb(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	matrix := [][]float64{
		{0.3588, 0.7044, 0.1368},
		{0.2990, 0.5870, 0.1140},
		{0.2392, 0.4696, 0.0912},
	}

	err = image.Recomb(matrix)
	require.NoError(t, err)
}

func TestImageRef_Recomb_Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	matrix := [][]float64{
		{0.3588, 0.7044, 0.1368, 0},
		{0.2990, 0.5870, 0.1140, 0},
		{0.2392, 0.4696, 0.0912, 0},
	}

	err = image.Recomb(matrix)
	require.Error(t, err)
}

func TestCopy(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	imageCopy, err := image.Copy()
	require.NoError(t, err)

	assert.Equal(t, image.buf, imageCopy.buf)
}

func BenchmarkExportImage(b *testing.B) {
	Startup(nil)

	fileBuf, err := os.ReadFile(resources + "heic-24bit.heic")
	require.NoError(b, err)

	img, err := NewImageFromBuffer(fileBuf)
	require.NoError(b, err)

	b.SetParallelism(100)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _, err = img.ExportJpeg(nil)
		require.NoError(b, err)
	}
	b.ReportAllocs()
}

func BenchmarkOpenBMPImage(b *testing.B) {
	Startup(nil)

	fileBuf, err := os.ReadFile(resources + "large.bmp")
	require.NoError(b, err)

	b.SetParallelism(100)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := NewImageFromBuffer(fileBuf)
		require.NoError(b, err)
	}
	b.ReportAllocs()
}

func TestMemstats(t *testing.T) {
	var stats MemoryStats
	ReadVipsMemStats(&stats)
	assert.NotNil(t, stats)
	assert.NotNil(t, stats.Allocs)
	assert.NotNil(t, stats.Files)
	assert.NotNil(t, stats.Mem)
	assert.NotNil(t, stats.MemHigh)
	govipsLog("govips", LogLevelInfo, fmt.Sprintf("MemoryStats: allocs: %d, files: %d, mem: %d, memhigh: %d", stats.Allocs, stats.Files, stats.Mem, stats.MemHigh))
}

func TestBands(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	bands := image.Bands()
	assert.Equal(t, bands, 3)
}

func TestCoding(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	coding := image.Coding()
	assert.Equal(t, coding, CodingNone)
}

func TestGetRotation(t *testing.T) {
	Startup(nil)

	rotation, flipped := GetRotationAngleFromExif(6)
	assert.Equal(t, rotation, Angle270)
	assert.Equal(t, flipped, false)

	rotation, flipped = GetRotationAngleFromExif(2)
	assert.Equal(t, rotation, Angle0)
	assert.Equal(t, flipped, true)

	rotation, flipped = GetRotationAngleFromExif(9)
	assert.Equal(t, rotation, Angle0)
	assert.Equal(t, flipped, false)

	rotation, flipped = GetRotationAngleFromExif(4)
	assert.Equal(t, rotation, Angle180)
	assert.Equal(t, flipped, true)

	rotation, flipped = GetRotationAngleFromExif(8)
	assert.Equal(t, rotation, Angle90)
	assert.Equal(t, flipped, false)
}

func TestResOffset(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	x := image.ResX()
	y := image.ResY()
	offsetX := image.OffsetX()
	offsetY := image.OffsetY()

	assert.Equal(t, x, float64(2.835))
	assert.Equal(t, y, float64(2.835))
	assert.Equal(t, offsetX, 0)
	assert.Equal(t, offsetY, 0)
}

func TestToBytes(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	buf1, err := image.ToBytes()
	assert.NoError(t, err)
	assert.Equal(t, 6220800, len(buf1))
}

func TestBandJoin(t *testing.T) {
	Startup(nil)

	image1, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	image2, err := NewImageFromFile(resources + "png-8bit.png")
	require.NoError(t, err)

	err = image1.BandJoin(image2)
	require.NoError(t, err)
}

func TestExtractBandToImage(t *testing.T) {
	Startup(nil)
	image1, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	v, err := image1.ExtractBandToImage(0, 2)
	require.NoError(t, err)
	require.Equal(t, v.Bands(), 2)

	_, err = v.ExtractBandToImage(0, 3)
	require.Error(t, err)
}

func TestBandSplit(t *testing.T) {
	Startup(nil)

	image1, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	bands, err := image1.BandSplit()
	require.NoError(t, err)

	require.Len(t, bands, 3)

	image2, err := NewImageFromFile(resources + "with_alpha.png")
	require.NoError(t, err)

	bands2, err := image2.BandSplit()
	require.NoError(t, err)

	require.Len(t, bands2, 4)
}

func TestIsColorSpaceSupport(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	supported := image.IsColorSpaceSupported()
	assert.True(t, supported)

	err = image.ToColorSpace(InterpretationError)
	assert.Error(t, err)
}

func TestPages_webp(t *testing.T) {
	Startup(nil)
	image, err := NewImageFromFile(resources + "webp-animated.webp")
	require.NoError(t, err)

	pages := image.Pages()
	assert.Equal(t, 8, pages)
}

func TestImageRef_Divide__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.ExtractBand(0, 2)
	require.NoError(t, err)

	denominator, err := NewImageFromFile(resources + "heic-24bit.heic")
	require.NoError(t, err)

	err = image.Divide(denominator)
	assert.Error(t, err)
}

func TestXYZ(t *testing.T) {
	Startup(nil)

	_, err := XYZ(100, 100)
	require.NoError(t, err)
}

func TestIdentity(t *testing.T) {
	Startup(nil)

	_, err := Identity(false)
	require.NoError(t, err)
	_, err = Identity(true)
	require.NoError(t, err)
}

func TestDeprecatedExportParams(t *testing.T) {
	Startup(nil)

	defaultExportParams := NewDefaultExportParams()
	assert.Equal(t, ImageTypeUnknown, defaultExportParams.Format)

	pngExportParams := NewPngExportParams()
	assert.Equal(t, 6, pngExportParams.Compression)
}

func TestNewImageFromReaderFail(t *testing.T) {
	r := strings.NewReader("")
	buf, err := NewImageFromReader(r)

	assert.Nil(t, buf)
	assert.Error(t, err)
}

func TestNewImageFromFileFail(t *testing.T) {
	buf, err := NewImageFromFile("/tmp/nonexistent-fasljdfalkjfadlafjladsfkjadfsljafdslk")

	assert.Nil(t, buf)
	assert.Error(t, err)
}

func TestImageRef_Cast(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)
	err = image.Cast(BandFormatUchar)
	assert.NoError(t, err)
	err = image.Cast(math.MaxInt8)
	assert.Error(t, err)
}

func TestImageRef_Average(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)
	average, err := image.Average()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, average)
}

func TestImageRef_FindTrim_White(t *testing.T) {
	image, err := NewImageFromFile(resources + "find_trim.png")
	assert.NoError(t, err)
	left, top, width, height, err := image.FindTrim(0, &Color{R: 255, G: 255, B: 255})
	assert.NoError(t, err)

	assert.Equal(t, 0, left)
	assert.Equal(t, 0, top)
	assert.Equal(t, 432, width)
	assert.Equal(t, 320, height)
}

func TestImageRef_FindTrim_Gray(t *testing.T) {
	image, err := NewImageFromFile(resources + "find_trim.png")
	assert.NoError(t, err)
	left, top, width, height, err := image.FindTrim(0, &Color{R: 238, G: 238, B: 238})
	assert.NoError(t, err)

	assert.Equal(t, 32, left)
	assert.Equal(t, 0, top)
	assert.Equal(t, 480, width)
	assert.Equal(t, 320, height)
}

func TestImageRef_FindTrim_Threshold(t *testing.T) {
	image, err := NewImageFromFile(resources + "find_trim.png")
	assert.NoError(t, err)
	left, top, width, height, err := image.FindTrim(17, &Color{R: 255, G: 255, B: 255})
	assert.NoError(t, err)

	assert.Equal(t, 80, left)
	assert.Equal(t, 32, top)
	assert.Equal(t, 352, width)
	assert.Equal(t, 256, height)
}

func TestImageRef_Height(t *testing.T) {
	image, err := NewImageFromFile(resources + "gif-animated-2.gif")
	assert.NoError(t, err)
	width := image.Height()
	assert.Equal(t, 90, width)
}

func TestImageRef_Linear_Fails(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)
	err = image.Linear([]float64{1, 2}, []float64{1, 2, 3})
	assert.Error(t, err)
}

func TestImageRef_AVIF(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "avif-8bit.avif")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.ExportNative()
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeAVIF, metadata.Format)
}

func TestImageRef_JP2K(t *testing.T) {
	if MajorVersion == 8 && MinorVersion < 11 {
		t.Skip("JPEG2000 is only supported in vips 8.11+")
	}
	Startup(nil)

	raw, err := os.ReadFile(resources + "jp2k-orientation-6.jp2")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.ExportJp2k(nil)
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeJP2K, metadata.Format)
	assert.Equal(t, 1, metadata.Pages)
}

func TestImageRef_CorruptedJPEG(t *testing.T) {
	Startup(nil)

	raw, err := os.ReadFile(resources + "jpg-corruption.jpg")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, _, err = img.ExportJpeg(nil)
	assert.Error(t, err, "VipsJpeg: Corrupt JPEG data: bad Huffman code")
}

func TestImageRef_Stats(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	bands := image.Bands()

	err = image.Stats()
	require.NoError(t, err)

	// May need updating if `vips_stats` adds more columns
	require.Equal(t, 10, image.Width())
	require.Equal(t, bands+1, image.Height())
}

func TestImageRef_HistogramFind(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.HistogramFind()
	require.NoError(t, err)
	require.Equal(t, 256, image.Width())
	require.Equal(t, 1, image.Height())
}

func TestImageRef_HistogramNormalize(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.HistogramFind()
	require.NoError(t, err)

	err = image.HistogramNormalise()
	require.NoError(t, err)
}

func TestImageRef_HistogramCumulative(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.HistogramFind()
	require.NoError(t, err)

	err = image.HistogramCumulative()
	require.NoError(t, err)
}

func TestImageRef_HistogramEntropy(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.HistogramFind()
	require.NoError(t, err)

	e, err := image.HistogramEntropy()
	require.NoError(t, err)
	require.True(t, e > 0)
}

func TestImageRef_SetPages(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "gif-animated.gif")
	require.NoError(t, err)
	require.Equal(t, 8, image.Pages())

	err = image.SetPages(3)
	require.NoError(t, err)
	require.Equal(t, 3, image.Pages())
}

func TestImageRef_SetGamma(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.Gamma(1.0 / 2.4)
	require.NoError(t, err)
}

// TODO unit tests to cover:
// NewImageFromReader failing test
// NewImageFromFile failing test
// Copy failing test
// SetOrientation failing test
// RemoveOrientation failing test
// ExportBuffer failing test
// Exporting a TIFF image
// Providing Linear() with different length a and b slices
// RemoveICCProfile failing test
// RemoveMetadata failing test
