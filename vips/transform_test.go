package vips

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestTransform_Resize_JPG(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_Resize_JPG_CMYK(t *testing.T) {
	goldenTest(t, resources+"jpg-32bit-cmyk-icc-swop.jpg", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_Resize_PNG_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_Resize_BMP(t *testing.T) {
	goldenTest(t, resources+"bmp.bmp", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_Resize_BMP_Alpha(t *testing.T) {
	goldenTest(t, resources+"bmp+alpha.bmp", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_Resize_HEIC(t *testing.T) {
	goldenTest(t, resources+"heic-24bit-exif.heic", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestTransform_StripMetadata(t *testing.T) {
	goldenTest(t, resources+"heic-24bit-exif.heic", func(tx *Transform) {
		tx.StripMetadata()
	})
}

func TestTransform_StripProfile(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-smpte.jpg", func(tx *Transform) {
		tx.StripProfile()
	})
}

func TestTransform_AutoRotate_0(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(tx *Transform) {
		tx.AutoRotate()
	})
}

func TestTransform_AutoRotate_1(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(tx *Transform) {
		tx.AutoRotate()
	})
}

func TestTransform_AutoRotate_5(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-5.jpg", func(tx *Transform) {
		tx.AutoRotate()
	})
}

func TestTransform_AutoRotate_6(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-6.jpg", func(tx *Transform) {
		tx.AutoRotate()
	})
}

func TestTransform_Scale_24bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		tx.Scale(0.2)
	})
}

func TestTransform_Scale_8bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(tx *Transform) {
		tx.Scale(0.25)
	})
}

func TestTransform_Sharpen_24bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		tx.Sharpen(sigma, x1, m2)
	})
}

func TestTransform_Sharpen_8bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(tx *Transform) {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		tx.Sharpen(sigma, x1, m2)
	})
}

func TestTransform_Modulate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(tx *Transform) {
		tx.Modulate(0.7, 0.5, 180)
	})
}

func TestTransform_Modulate_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		tx.Modulate(1.1, 1.2, 0)
	})
}

func TestTransform_ModulateHSV(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(tx *Transform) {
		tx.ModulateHSV(0.7, 0.5, 180)
	})
}

func TestTransform_ModulateHSV_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		tx.ModulateHSV(1.1, 1.2, 120)
	})
}

func TestTransform_BackgroundColor(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		tx.BackgroundColor(&Color{R: 255, G: 192, B: 203})
	})
}

func TestTransform_Resize_Crop(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(tx *Transform) {
		tx.Resize(100, 300).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_ScaleHeight(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-adobe-rgb.jpg", func(tx *Transform) {
		tx.ScaleHeight(0.5)
	})
}

func TestTransform_Resize_CenterCrop(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-adobe-rgb.jpg", func(tx *Transform) {
		tx.Resize(341, 256).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_BottomRightCrop(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-adobe-rgb.jpg", func(tx *Transform) {
		tx.Resize(341, 256).ResizeStrategy(ResizeStrategyCrop).Anchor(AnchorBottomRight)
	})
}

func TestTransform_OffsetCrop(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		tx.Resize(500, 720).CropOffsetX(120).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_OffsetCropBounds(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		tx.Resize(100, 100).CropOffsetX(120).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_RelativeOffsetCrop(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(tx *Transform) {
		tx.Resize(500, 720).CropRelativeOffsetX(0.1066).ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestTransform_Rotate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(tx *Transform) {
		tx.Rotate(Angle90)
	})
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

	t.Log("writing golden file: " + goldenFile)
	err := ioutil.WriteFile(goldenFile, buf, 0644)
	assert.NoError(t, err)
}
