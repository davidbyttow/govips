package vips

import (
	"bytes"
	"fmt"
	"image"
	jpeg2 "image/jpeg"
	"image/png"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/image/bmp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestImage_Embed_ExtendBackground_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Embed(0, 0, 1000, 500, ExtendBackground)
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

//This test is disabled until this issue is resolved: https://github.com/libvips/libvips/pull/1745
func _TestImageRef_Orientation_Issue(t *testing.T) {
	goldenTest(t, resources+"orientation-issue-1.jpg",
		func(img *ImageRef) error {
			return img.Resize(0.9, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 6, result.GetOrientation())
		},
		NewDefaultWEBPExportParams())
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
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		//usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2, SharpenModeRGB)
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

func TestImage_Zoom(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.Zoom(2, 3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 200, result.Width())
			assert.Equal(t, 300, result.Height())
		}, nil)
}

func TestImage_Thumbnail_NoCrop(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(img *ImageRef) error {
			return img.Thumbnail(36, 36, InterestingNone)
		},
		func(result *ImageRef) {
			assert.Equal(t, 36, result.Width())
			assert.Equal(t, 24, result.Height())
		}, nil)
}

func TestImage_Thumbnail_CropCentered(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(img *ImageRef) error {
			return img.Thumbnail(25, 25, InterestingCentre)
		},
		func(result *ImageRef) {
			assert.Equal(t, 25, result.Width())
			assert.Equal(t, 25, result.Height())
		}, nil)
}

func TestImage_ResizeWithVScale(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.ResizeWithVScale(1.1, 1.2, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 110, result.Width())
			assert.Equal(t, 120, result.Height())
		}, nil)
}

func TestImage_Add(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			addend, err := NewImageFromFile(resources + "heic-24bit.heic")
			require.NoError(t, err)

			return img.Add(addend)
		},
		func(result *ImageRef) {
			assert.Equal(t, 1440, result.Width())
			assert.Equal(t, 960, result.Height())
		}, nil)
}

func TestImage_Multiply(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			multiplier, err := NewImageFromFile(resources + "heic-24bit.heic")
			require.NoError(t, err)

			return img.Multiply(multiplier)
		},
		func(result *ImageRef) {
			assert.Equal(t, 1440, result.Width())
			assert.Equal(t, 960, result.Height())
		}, nil)
}

func TestImageRef_Divide(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			denominator, err := NewImageFromFile(resources + "heic-24bit.heic")
			require.NoError(t, err)

			return img.Divide(denominator)
		},
		func(result *ImageRef) {
			assert.Equal(t, 1440, result.Width())
			assert.Equal(t, 960, result.Height())
		}, nil)
}

func TestImage_GaussianBlur(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.GaussianBlur(10.5)
	}, nil, nil)
}

func TestImage_BandJoinConst(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.BandJoinConst([]float64{255})
	}, nil, nil)
}

func TestImage_SmartCrop(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.SmartCrop(60, 80, InterestingCentre)
	}, func(result *ImageRef) {
		assert.Equal(t, 60, result.Width())
		assert.Equal(t, 80, result.Height())
	}, nil)
}

func TestImage_DrawRect(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.DrawRect(ColorRGBA{
			R: 255,
			G: 255,
			B: 0,
			A: 0,
		}, 20, 20, img.Width()-40, img.Height()-40, true)
	}, nil, nil)
}

func TestImage_DrawRectRGBA(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		err := img.AddAlpha()
		assert.Nil(t, err)
		return img.DrawRect(ColorRGBA{
			R: 255,
			G: 255,
			B: 0,
			A: 255,
		}, 20, 20, img.Width()-40, img.Height()-40, false)
	}, nil, nil)
}

func TestImage_SimilarityRGB(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.Similarity(0.5, 5, &ColorRGBA{R: 127, G: 127, B: 127, A: 127},
			10, 10, 20, 20)
	}, func(result *ImageRef) {
		assert.Equal(t, 3, result.Bands())
	}, nil)
}

func TestImage_SimilarityRGBA(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		err := img.AddAlpha()
		assert.Nil(t, err)
		err = img.Similarity(0.5, 5, &ColorRGBA{R: 127, G: 127, B: 127, A: 127},
			10, 10, 20, 20)
		assert.Nil(t, err)
		assert.Equal(t, 4, img.Bands())
		return nil
	}, nil, nil)
}

func TestImage_Decode_JPG(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		goImg, err := img.ToImage(nil)

		buf := new(bytes.Buffer)
		err = jpeg2.Encode(buf, goImg, nil)
		assert.Nil(t, err)

		config, format, err := image.DecodeConfig(buf)
		assert.Nil(t, err)
		assert.Equal(t, "jpeg", format)
		assert.NotNil(t, config)
		assert.True(t, config.Height > 0)
		assert.True(t, config.Width > 0)
		return nil
	}, nil, nil)
}

func TestImage_Decode_BMP(t *testing.T) {
	goldenTest(t, resources+"bmp.bmp", func(img *ImageRef) error {
		goImg, err := img.ToImage(nil)

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

func TestImage_Decode_PNG(t *testing.T) {
	goldenTest(t, resources+"png-8bit.png", func(img *ImageRef) error {
		goImg, err := img.ToImage(nil)

		buf := new(bytes.Buffer)
		err = png.Encode(buf, goImg)
		assert.Nil(t, err)

		config, format, err := image.DecodeConfig(buf)
		assert.Nil(t, err)
		assert.Equal(t, "png", format)
		assert.NotNil(t, config)
		assert.Equal(t, 150, config.Height)
		assert.Equal(t, 200, config.Width)
		return nil
	}, nil, nil)
}

func TestImage_Invert(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.Invert()
	}, nil, nil)
}

func TestImage_Flip(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.Flip(DirectionVertical)
	}, nil, nil)
}

func TestImage_Embed(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.Embed(10, 10, 20, 20, ExtendBlack)
	}, nil, nil)
}

func TestImage_ExtractBand(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		return img.ExtractBand(2, 1)
	}, nil, nil)
}

func TestImage_Flatten(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		return img.Flatten(&Color{R: 32, G: 64, B: 128})
	}, nil, nil)
}

func TestImage_Tiff(t *testing.T) {
	goldenTest(t, resources+"tif.tif", func(img *ImageRef) error {
		return img.OptimizeICCProfile()
	}, nil, nil)
}

func TestImage_Black(t *testing.T) {
	Startup(nil)
	i, err := Black(10, 20)
	require.NoError(t, err)
	buf, metadata, err := i.Export(nil)
	require.NoError(t, err)

	assertGoldenMatch(t, resources+"jpg-24bit.jpg", buf, metadata.Format)
}

func goldenTest(t *testing.T, file string, exec func(img *ImageRef) error, validate func(img *ImageRef), params *ExportParams) []byte {
	Startup(nil)

	i, err := NewImageFromFile(file)
	require.NoError(t, err)

	err = exec(i)
	require.NoError(t, err)

	buf, metadata, err := i.Export(params)
	require.NoError(t, err)

	if validate != nil {
		result, err := NewImageFromBuffer(buf)
		require.NoError(t, err)

		validate(result)
	}

	assertGoldenMatch(t, file, buf, metadata.Format)

	return buf
}

func getEnvironment() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "darwin":
		out, err := exec.Command("sw_vers", "-productVersion").Output()
		if err != nil {
			return "macos-unknown"
		}
		majorVersion := strings.Split(strings.TrimSpace(string(out)), ".")[0]
		return "macos-" + majorVersion
	case "linux":
		out, err := exec.Command("lsb_release", "-cs").Output()
		if err != nil {
			return "linux"
		}
		strout := strings.TrimSuffix(string(out), "\n")
		return "linux-" + strout
	}
	// default to linux assets otherwise
	return "linux"
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
			fmt.Printf("To diff and prompt to replace golden file: ./diff-replace-golden.sh %v %v\n", goldenPath, failedPath)
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
