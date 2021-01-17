package vips

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

	srcBytes, err := ioutil.ReadFile(resources + "webp+alpha.webp")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, _, err = img.Export(nil)
	assert.NoError(t, err)
}

func TestImageRef_WebP__ReducedEffort(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "webp+alpha.webp")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)

	params := NewDefaultWEBPExportParams()
	params.Effort = 2
	_, _, err = img.Export(params)
	assert.NoError(t, err)
}

func TestImageRef_PNG(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "png-24bit.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)

	// check random access by encoding twice
	_, _, err = img.Export(nil)
	assert.NoError(t, err)
	_, _, err = img.Export(nil)
	assert.NoError(t, err)
}

func TestImageRef_HEIF(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "heic-24bit-exif.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.Export(nil)
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, metadata.Format)
}

func TestImageRef_HEIF_MIF1(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "heic-24bit.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.Export(nil)
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, metadata.Format)
}

func TestImageRef_HEIF_ftypmsf1(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "heic-ftypmsf1.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	_, metadata, err := img.Export(nil)
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, metadata.Format)
}

func TestImageRef_BMP(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "bmp.bmp")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	exported, metadata, err := img.Export(nil)
	assert.NoError(t, err)
	assert.Equal(t, ImageTypePNG, metadata.Format)
	assert.NotNil(t, exported)
}

func TestImageRef_SVG(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "svg.svg")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, ImageTypeSVG, img.Metadata().Format)
}

func TestImageRef_SVG_1(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "svg_1.svg")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, ImageTypeSVG, img.Metadata().Format)
}

func TestImageRef_SVG_2(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "svg_2.svg")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, ImageTypeSVG, img.Metadata().Format)
}

func TestImageRef_OverSizedMetadata(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "png-bad-metadata.png")
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

	_, _, err = img.Export(nil)
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
	_, _, err = img.Export(nil)
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

	assert.Equal(t, 6, image.GetOrientation())
}

func TestImageRef_GetOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	assert.Equal(t, 0, image.GetOrientation())
}

func TestImageRef_SetOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)

	err = image.SetOrientation(5)
	require.NoError(t, err)

	assert.Equal(t, 5, image.GetOrientation())
}

func TestImageRef_SetOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.SetOrientation(5)
	require.NoError(t, err)

	assert.Equal(t, 5, image.GetOrientation())
}

func TestImageRef_RemoveOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)

	err = image.RemoveOrientation()
	require.NoError(t, err)

	assert.Equal(t, 0, image.GetOrientation())
}

func TestImageRef_RemoveOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	err = image.RemoveOrientation()
	require.NoError(t, err)

	assert.Equal(t, 0, image.GetOrientation())
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

	assert.Equal(t, 5, image.GetOrientation())
}

// Known issue: libvips does not write EXIF into WebP:
// https://github.com/libvips/libvips/pull/1745
func TestImageRef_RemoveMetadata__RetainsOrientation__WebP(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "webp-orientation-6.webp")
	require.NoError(t, err)

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.Equal(t, 6, image.GetOrientation())
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

func TestImageRef_Close(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)

	image.Close()
	assert.NotNil(t, image.image)

	image.close()
	assert.Nil(t, image.image)

	PrintObjectReport("Final")
}

func TestImageRef_Close__AlreadyClosed(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)

	go image.close()
	go image.close()
	go image.close()
	go image.close()
	defer image.close()
	image.close()

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

	images := []*ImageComposite{}
	for i, uri := range []string{"png-8bit+alpha.png", "png-24bit+alpha.png"} {
		image, err := NewImageFromFile(resources + uri)
		require.NoError(t, err)

		//add offset test
		images = append(images, &ImageComposite{image, BlendModeOver, (i + 1) * 20, (i + 2) * 20})
	}

	err = image.CompositeMulti(images)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
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

	fileBuf, err := ioutil.ReadFile(resources + "heic-24bit.heic")
	require.NoError(b, err)

	img, err := NewImageFromBuffer(fileBuf)
	require.NoError(b, err)

	b.SetParallelism(100)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _, err = img.Export(NewDefaultJPEGExportParams())
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
	offx := image.OffsetX()
	offy := image.OffsetY()

	assert.Equal(t, x, float64(2.835))
	assert.Equal(t, y, float64(2.835))
	assert.Equal(t, offx, 0)
	assert.Equal(t, offy, 0)
}

func TestToBytes(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	buf1, err := image.ToBytes()
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

func TestIsColorSpaceSupport(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	supported := image.IsColorSpaceSupported()
	assert.True(t, supported)

	err = image.ToColorSpace(InterpretationError)
	assert.Error(t, err)
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

func TestImageRef_Linear_Fails(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)
	err = image.Linear([]float64{1,2}, []float64{1,2,3})
	assert.Error(t, err)
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
