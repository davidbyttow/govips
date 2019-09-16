package vips

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransform_Resize(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_BMP__Alpha(t *testing.T) {
	goldenTest(t, resources+"with-alpha.bmp", func(tx *Transform) {
		tx.AutoRotate()
	})
}

func TestTransform_HEIC__Resize(t *testing.T) {
	goldenTest(t, resources+"citron.heic", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_Flatten(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.BackgroundColor(&Color{R: 255, G: 192, B: 203})
	})
}

func TestTransform_Resize_RetainICC(t *testing.T) {
	goldenTest(t, resources+"icc.jpg", func(tx *Transform) {
		tx.ResizeWidth(300).StripMetadata()
	})
}

func TestTransform_Resize_StripICC(t *testing.T) {
	goldenTest(t, resources+"icc.jpg", func(tx *Transform) {
		tx.ResizeWidth(300).StripProfile()
	})
}

func TestTransform_AdobeRGB_sRGB_StripICC(t *testing.T) {
	goldenTest(t, resources+"adobe-rgb.jpg", func(tx *Transform) {
		tx.StripProfile()
	})
}

func TestTransform_AdobeRGB_sRGB_StripMetadata(t *testing.T) {
	// this strips the ICC profile as well
	goldenTest(t, resources+"adobe-rgb.jpg", func(tx *Transform) {
		tx.StripMetadata()
	})
}

func TestTransform_AdobeRGB_sRGB_Resize_RetainMetadata(t *testing.T) {
	// this strips the ICC profile as well
	goldenTest(t, resources+"adobe-rgb.jpg", func(tx *Transform) {
		tx.Resize(1000, 1000)
	})
}

func TestTransform_Resize_Crop(t *testing.T) {
	goldenTest(t, resources+"colors.png", func(tx *Transform) {
		tx.Resize(100, 300).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_ResizeShapes(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(341, 256)
	})
}

func TestTransform_RelativeResizeShapes(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.ScaleHeight(0.5)
	})
}

func TestTransform_CenterCrop(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(341, 256).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_BottomRightCrop(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(341, 256).ResizeStrategy(ResizeStrategyCrop).Anchor(AnchorBottomRight)
	})
}

func TestTransform_OffsetCrop(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Resize(500, 720).CropOffsetX(120).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_OffsetCropBounds(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Resize(100, 100).CropOffsetX(120).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_RelativeOffsetCrop(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Resize(500, 720).CropRelativeOffsetX(0.1066).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_Rotate(t *testing.T) {
	goldenTest(t, resources+"canyon.jpg", func(tx *Transform) {
		tx.Rotate(Angle90)
	})
}

func TestTransform_AutoRotate(t *testing.T) {
	goldenTest(t, resources+"canyon.jpg", func(tx *Transform) {
		tx.AutoRotate()
	})
}

func TestTransform_Scale3x(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Scale(3.0)
	})
}

func TestTransform_MaxScale(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.MaxScale(1.0).ResizeWidth(100000)
	})
}

func TestTransform_Overlay(t *testing.T) {
	if testing.Short() {
		return
	}
	var tomatoesData, cloverData []byte
	t.Run("tomatoes", func(t *testing.T) {
		tomatoesData = goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
			tx.ResizeWidth(320)
		})
	})
	t.Run("clover", func(t *testing.T) {
		cloverData = goldenTest(t, resources+"clover.png", func(tx *Transform) {
			tx.ResizeWidth(64)
		})
	})
	tomatoes, err := NewImageFromBuffer(tomatoesData)
	require.NoError(t, err)
	defer tomatoes.Close()

	clover, err := NewImageFromBuffer(cloverData)
	require.NoError(t, err)
	defer clover.Close()

	err = tomatoes.Composite(clover, BlendModeOver, 0, 0)
	require.NoError(t, err)

	buf, _, err := tomatoes.Export(nil)
	require.NoError(t, err)
	assertGoldenMatch(t, resources+"tomatoes.png", buf)
}

func TestTransform_BandJoin(t *testing.T) {
	image1, err := NewImageFromFile(resources + "tomatoes.png")
	require.NoError(t, err)
	defer image1.Close()

	image2, err := NewImageFromFile(resources + "clover.png")
	require.NoError(t, err)
	defer image2.Close()

	err = image1.BandJoin(image2)
	require.NoError(t, err)

	buf, _, err := image1.Export(nil)
	require.NoError(t, err)
	assertGoldenMatch(t, resources+"tomatoes.png", buf)
}

func goldenTest(t *testing.T, file string, fn func(t *Transform)) []byte {
	if testing.Short() {
		return nil
	}

	Startup(nil)

	i, err := NewImageFromFile(file)
	require.NoError(t, err)
	defer i.Close()

	tx := NewTransform()

	fn(tx)

	buf, _, err := tx.ApplyAndExport(i)
	require.NoError(t, err)
	assertGoldenMatch(t, file, buf)

	return buf
}

func assertGoldenMatch(t *testing.T, file string, buf []byte) {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		panic("bad filename")
	}
	name := strings.Replace(t.Name(), "/", "_", -1)
	prefix := file[:i] + "." + name
	ext := file[i:]
	goldenFile := prefix + ".golden" + ext
	golden, _ := ioutil.ReadFile(goldenFile)
	if golden != nil {
		if !assert.Equal(t, golden, buf) {
			failed := prefix + ".failed" + ext
			err := ioutil.WriteFile(failed, buf, 0666)
			if err != nil {
				panic(err)
			}
		}
		return
	}
	t.Log("Writing golden file: " + goldenFile)
	err := ioutil.WriteFile(goldenFile, buf, 0644)
	assert.NoError(t, err)
}
