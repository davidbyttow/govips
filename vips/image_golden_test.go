package vips

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestImage_Resize_Downscale(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.Resize(0.9, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 90, result.Width())
			assert.Equal(t, 90, result.Height())
		}, nil)
}

func TestImage_Resize_Upscale(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.Resize(1.1, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 110, result.Width())
			assert.Equal(t, 110, result.Height())
		}, nil)
}

func TestImage_Resize_Downscale_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Resize(0.9, KernelLanczos3)
	}, nil, nil)
}

func TestImage_Resize_Upscale_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Resize(1.1, KernelLanczos3)
	}, nil, nil)
}

func TestImage_OptimizeICCProfile_CMYK(t *testing.T) {
	goldenTest(t, resources+"jpg-32bit-cmyk-icc-swop.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_OptimizeICCProfile_RGB_No_Profile(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.False(t, result.HasICCProfile())
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_OptimizeICCProfile_RGB_Embedded(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-adobe-rgb.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_OptimizeICCProfile_Grey(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-gray-scale-with-icc-profile.jpg",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			assert.Equal(t, InterpretationBW, result.Interpretation())
		}, nil)
}

func TestImage_RemoveICCProfile(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-smpte.jpg",
		func(img *ImageRef) error {
			return img.RemoveICCProfile()
		},
		func(result *ImageRef) {
			assert.False(t, result.HasICCProfile())
		}, nil)
}

func TestImage_RemoveMetadata(t *testing.T) {
	goldenTest(t, resources+"heic-24bit-exif.heic", func(img *ImageRef) error {
		return img.RemoveMetadata()
	}, nil, nil)
}

func TestImageRef_RemoveMetadata_Leave_Orientation(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-5.jpg",
		func(img *ImageRef) error {
			return img.RemoveMetadata()
		},
		func(result *ImageRef) {
			assert.Equal(t, 5, result.GetOrientation())
		}, nil)
}

func TestImageRef_Orientation_Issue(t *testing.T) {
	goldenTest(t, resources+"orientation-issue-1.jpg",
		func(img *ImageRef) error {
			return img.Resize(0.9, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 6, result.GetOrientation())
		},
		NewDefaultWEBPExportParams())
}

func TestImageRef_PngToWebp_OptimizeICCProfile_NearLossless_HasProfile(t *testing.T) {
	exportParams := NewDefaultWEBPExportParams()
	exportParams.Quality = 90
	exportParams.NearLossless = true

	goldenTest(t, resources+"has-icc-profile.png",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
		}, exportParams)
}

func TestImageRef_PngToWebp_OptimizeICCProfile_LosslessAndNearLossless(t *testing.T) {
	exportParams := NewDefaultWEBPExportParams()
	exportParams.Quality = 90
	exportParams.Lossless = true
	exportParams.NearLossless = true

	goldenTest(t, resources+"has-icc-profile.png",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
		}, exportParams)
}

func TestImageRef_PngToWebp_OptimizeICCProfile_Lossless(t *testing.T) {
	exportParams := NewDefaultWEBPExportParams()
	exportParams.Quality = 90
	exportParams.Lossless = true

	goldenTest(t, resources+"has-icc-profile.png",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
		}, exportParams)
}

func TestImageRef_RemoveMetadata_Leave_Profile(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(img *ImageRef) error {
			return img.RemoveMetadata()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
		}, nil)
}

func TestImage_AutoRotate_0(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 0, result.GetOrientation())
		}, nil)
}

func TestImage_AutoRotate_1(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.GetOrientation())
		}, nil)
}

func TestImage_AutoRotate_5(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-5.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.GetOrientation())
		}, nil)
}

func TestImage_AutoRotate_6(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-6.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.GetOrientation())
		}, nil)
}

func TestImage_AutoRotate_6__jpeg_to_webp(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-6.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			// expected should be 1
			// Known issue: libvips does not write EXIF into WebP:
			// https://github.com/libvips/libvips/pull/1745
			//assert.Equal(t, 0, result.GetOrientation())
		}, NewDefaultWEBPExportParams())
}

func TestImage_AutoRotate_6__heic_to_jpg(t *testing.T) {
	goldenTest(t, resources+"heic-orientation-6.heic",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.GetOrientation())
		}, NewDefaultJPEGExportParams())
}

func TestImage_Sharpen_Luminescence_24bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2, SharpenModeLuminescence)
	}, nil, nil)
}

func TestImage_Sharpen_RGB_24bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2, SharpenModeRGB)
	}, nil, nil)
}

func TestImage_Sharpen_8bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2, SharpenModeLuminescence)
	}, nil, nil)
}

func TestImage_Modulate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(img *ImageRef) error {
		return img.Modulate(0.7, 0.5, 180)
	}, nil, nil)
}

func TestImage_Modulate_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.Modulate(1.1, 1.2, 0)
	}, nil, nil)
}

func TestImage_ModulateHSV(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(img *ImageRef) error {
		return img.ModulateHSV(0.7, 0.5, 180)
	}, nil, nil)
}

func TestImage_ModulateHSV_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.ModulateHSV(1.1, 1.2, 120)
	}, nil, nil)
}

func TestImage_ExtractArea(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return img.ExtractArea(10, 20, 200, 100)
		},
		func(result *ImageRef) {
			assert.Equal(t, 200, result.Width())
			assert.Equal(t, 100, result.Height())
		}, nil)
}

func TestImage_Rotate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return img.Rotate(Angle90)
		},
		func(result *ImageRef) {
			assert.Equal(t, 1600, result.Width())
			assert.Equal(t, 2560, result.Height())
		}, nil)
}

func TestImageRef_Linear1(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		return img.Linear1(3, 4)
	}, nil, nil)
}

func TestImageRef_Linear_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.Linear([]float64{1.1, 1.2, 1.3, 1.4}, []float64{1, 2, 3, 4})
	}, nil, nil)
}

func goldenTest(t *testing.T, goldenFileName string, exec func(img *ImageRef) error, validate func(img *ImageRef),
	exportParams *ExportParams) []byte {
	Startup(nil)

	goldenImage, err := NewImageFromFile(goldenFileName)
	require.NoError(t, err)
	defer goldenImage.Close()

	err = exec(goldenImage)
	require.NoError(t, err)

	buf, metadata, err := goldenImage.Export(exportParams)
	require.NoError(t, err)

	if validate != nil {
		result, err := NewImageFromBuffer(buf)
		require.NoError(t, err)
		defer result.Close()

		validate(result)
	}

	assertGoldenMatch(t, goldenFileName, buf, metadata.Format)

	return buf
}

func assertGoldenMatch(t *testing.T, imagePath string, imageData []byte, format ImageType) {
	imagePath, err := filepath.Abs(imagePath)
	panicOnError(err)

	dotIndex := strings.LastIndex(imagePath, ".")
	if dotIndex < 0 {
		panic("bad filename")
	}

	testName := strings.Replace(t.Name(), "/", "_", -1)
	testName = strings.Replace(testName, "TestImage_", "", -1)
	prefix := imagePath[:dotIndex] + "." + testName
	ext := format.FileExt()
	goldenPath, err := filepath.Abs(prefix + ".golden" + ext)
	panicOnError(err)

	expectedData, _ := ioutil.ReadFile(goldenPath)
	if expectedData != nil {
		if !bytes.Equal(expectedData, imageData) {
			failedPath := prefix + ".failed" + ext
			assert.Fail(t, "Output not equal to golden result",
				"expected\t%v (%v bytes)\n"+
					"but got\t\t%v (%v bytes)",
				goldenPath, len(expectedData), failedPath, len(imageData))
			fmt.Printf("To diff and prompt to replace golden file: ./diff-replace-golden.sh %v %v", goldenPath, failedPath)
			err := ioutil.WriteFile(failedPath, imageData, 0666)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	t.Log("writing golden file: " + goldenPath)
	err = ioutil.WriteFile(goldenPath, imageData, 0644)
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
