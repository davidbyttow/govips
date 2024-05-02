package vips

import (
	"bytes"
	"image"
	jpeg2 "image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/image/bmp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImage_PNG_64bit_OptimizeICCProfile(t *testing.T) {
	goldenTest(t, resources+"png-alpha-64bit.png",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		nil,
		exportPng(NewPngExportParams()))
}

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

func TestImage_Embed_ExtendWhite_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Embed(0, 0, 1000, 500, ExtendWhite)
	}, func(img *ImageRef) {
		point, err := img.GetPoint(999, 0)
		assert.NoError(t, err)
		assert.Equal(t, point, []float64{255, 255, 255, 255})
	}, nil)
}

func TestImage_EmbedBackgroundRGBA_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.EmbedBackgroundRGBA(0, 0, 1000, 500, &ColorRGBA{R: 238, G: 238, B: 238, A: 50})
	}, func(img *ImageRef) {
		point, err := img.GetPoint(999, 0)
		assert.NoError(t, err)
		assert.Equal(t, point, []float64{238, 238, 238, 50})
	}, nil)
}

func TestImage_EmbedBackground_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.EmbedBackground(0, 0, 1000, 500, &Color{R: 238, G: 238, B: 238})
	}, func(img *ImageRef) {
		point, err := img.GetPoint(999, 0)
		assert.NoError(t, err)
		assert.Equal(t, point, []float64{238, 238, 238, 255})
	}, nil)
}

func TestImage_EmbedBackground_NoAlpha(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.EmbedBackground(0, 0, 500, 300, &Color{R: 238, G: 238, B: 238})
		},
		func(result *ImageRef) {
			point, err := result.GetPoint(499, 0)
			assert.NoError(t, err)
			assert.Equal(t, point, []float64{238, 238, 238})
		}, nil)
}

func TestImage_TransformICCProfile_RGB_No_Profile(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.TransformICCProfile(SRGBIEC6196621ICCProfilePath)
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			iccProfileData := result.GetICCProfile()
			assert.Greater(t, len(iccProfileData), 0)
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_TransformICCProfile_RGB_Embedded(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-adobe-rgb.jpg",
		func(img *ImageRef) error {
			return img.TransformICCProfile(SRGBIEC6196621ICCProfilePath)
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile())
			assert.Equal(t, InterpretationSRGB, result.Interpretation())
		}, nil)
}

func TestImage_TransformICCProfileWithFallback(t *testing.T) {
	t.Run("RGB source without ICC", func(t *testing.T) {
		goldenTest(t, resources+"jpg-24bit-rgb-no-icc.jpg",
			func(img *ImageRef) error {
				return img.TransformICCProfileWithFallback(SRGBIEC6196621ICCProfilePath, resources+"adobe-rgb.icc")
			},
			func(result *ImageRef) {
				assert.True(t, result.HasICCProfile())
				iccProfileData := result.GetICCProfile()
				assert.Greater(t, len(iccProfileData), 0)
				assert.Equal(t, InterpretationSRGB, result.Interpretation())
			}, nil)
	})
	t.Run("RGB source with ICC", func(t *testing.T) {
		goldenTest(t, resources+"jpg-24bit-icc-adobe-rgb.jpg",
			func(img *ImageRef) error {
				return img.TransformICCProfileWithFallback(SRGBIEC6196621ICCProfilePath, SRGBV2MicroICCProfilePath)
			},
			func(result *ImageRef) {
				assert.True(t, result.HasICCProfile())
				iccProfileData := result.GetICCProfile()
				assert.Greater(t, len(iccProfileData), 0)
				assert.Equal(t, InterpretationSRGB, result.Interpretation())
			}, nil)
	})
	t.Run("CMYK source without ICC", func(t *testing.T) {
		goldenTest(t, resources+"jpg-32bit-cmyk-no-icc.jpg",
			func(img *ImageRef) error {
				return img.TransformICCProfileWithFallback(SRGBIEC6196621ICCProfilePath, "cmyk")
			},
			func(result *ImageRef) {
				assert.True(t, result.HasICCProfile())
				iccProfileData := result.GetICCProfile()
				assert.Greater(t, len(iccProfileData), 0)
				assert.Equal(t, InterpretationSRGB, result.Interpretation())
			}, nil)
	})
	t.Run("CMYK source with ICC", func(t *testing.T) {
		goldenTest(t, resources+"jpg-32bit-cmyk-icc-swop.jpg",
			func(img *ImageRef) error {
				return img.TransformICCProfileWithFallback(SRGBIEC6196621ICCProfilePath, "cmyk")
			},
			func(result *ImageRef) {
				assert.True(t, result.HasICCProfile())
				iccProfileData := result.GetICCProfile()
				assert.Greater(t, len(iccProfileData), 0)
				assert.Equal(t, InterpretationSRGB, result.Interpretation())
			}, nil)
	})
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

func TestImageRef_PngToWebp_OptimizeICCProfile_Lossless(t *testing.T) {
	exportParams := NewWebpExportParams()
	exportParams.Quality = 90
	exportParams.Lossless = true

	testWebpOptimizeIccProfile(t, exportParams)
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
			assert.True(t, img.HasICCProfile())
			return img.RemoveICCProfile()
		},
		func(result *ImageRef) {
			assert.False(t, result.HasICCProfile())
		}, nil)
}

// NOTE: The JPEG spec requires some minimal exif data including exif-ifd0-Orientation.
// libvips always adds these fields back but they should not be a privacy concern.
// HEIC images require the same fields and behave the same way in libvips.
func TestImage_RemoveMetadata_Removes_Exif(t *testing.T) {
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

func TestImageRef_RemoveMetadata_Leave_Orientation(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-5.jpg",
		func(img *ImageRef) error {
			return img.RemoveMetadata()
		},
		func(result *ImageRef) {
			assert.Equal(t, 5, result.Orientation())
		}, nil)
}

// https://iptc.org/std/photometadata/specification/IPTC-PhotoMetadata#creator
// https://iptc.org/std/photometadata/specification/IPTC-PhotoMetadata#credit-line
// https://iptc.org/std/photometadata/specification/IPTC-PhotoMetadata#copyright-notice
func TestImageRef_RemoveMetadata_Leave_Copyright(t *testing.T) {
	goldenTest(t, resources+"copyright.jpeg",
		func(img *ImageRef) error {
			return img.RemoveMetadata("exif-ifd0-Copyright", "exif-ifd0-Artist")
		},
		func(result *ImageRef) {
			assert.Contains(t, result.ImageFields(), "exif-ifd0-Copyright")
		}, nil)
}

func TestImageRef_Orientation_Issue(t *testing.T) {
	goldenTest(t, resources+"orientation-issue-1.jpg",
		func(img *ImageRef) error {
			return img.Resize(0.9, KernelLanczos3)
		},
		func(result *ImageRef) {
			assert.Equal(t, 6, result.Orientation())
		},
		exportWebp(nil),
	)
}

func TestImageRef_RemoveMetadata_Leave_Profile(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(img *ImageRef) error {
			return img.RemoveMetadata()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile(), "should have an ICC profile")
		}, nil)
}

func TestImage_AutoRotate_0(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 0, result.Orientation())
		}, nil)
}

func TestImage_AutoRotate_1(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.Orientation())
		}, nil)
}

func TestImage_AutoRotate_5(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-5.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.Orientation())
		}, nil)
}

func TestImage_AutoRotate_6(t *testing.T) {
	goldenTest(t, resources+"jpg-orientation-6.jpg",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 1, result.Orientation())
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
			// assert.Equal(t, 0, result.Orientation())
		},
		exportWebp(nil),
	)
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

func TestImage_Sharpen_24bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		// usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2)
	}, nil, nil)
}

func TestImage_Sharpen_8bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		// usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2)
	}, nil, nil)
}

func TestImage_Sobel(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Sobel()
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

func TestImage_Thumbnail_NoUpscale(t *testing.T) {
	goldenTest(t, resources+"jpg-8bit-grey-icc-dot-gain.jpg",
		func(img *ImageRef) error {
			return img.ThumbnailWithSize(9999, 9999, InterestingNone, SizeDown)
		},
		func(result *ImageRef) {
			assert.Equal(t, 715, result.Width())
			assert.Equal(t, 483, result.Height())
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
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.GaussianBlur(10.5, 0.2)
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

func TestImage_Crop(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.Crop(10, 10, 60, 80)
	}, func(result *ImageRef) {
		assert.Equal(t, 60, result.Width())
		assert.Equal(t, 80, result.Height())
	}, nil)
}

func TestImage_Replicate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg", func(img *ImageRef) error {
		return img.Replicate(3, 2)
	}, func(result *ImageRef) {
		assert.Equal(t, 300, result.Width())
		assert.Equal(t, 200, result.Height())
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

func TestImage_Rank(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		return img.Rank(15, 15, 224)
	}, nil, nil)
}

func TestImage_GetPointWhite(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		point, err := img.GetPoint(10, 10)

		assert.Equal(t, 3, len(point))
		assert.Equal(t, 255.0, point[0])
		assert.Equal(t, 255.0, point[1])
		assert.Equal(t, 255.0, point[2])

		return err
	}, nil, nil)
}

func TestImage_GetPointYellow(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		point, err := img.GetPoint(400, 10)

		assert.Equal(t, 3, len(point))
		assert.Equal(t, 255.0, point[0])
		assert.Equal(t, 255.0, point[1])
		assert.Equal(t, 0.0, point[2])

		return err
	}, nil, nil)
}

func TestImage_GetPointWhiteR(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		point, err := img.GetPoint(10, 10)

		assert.Equal(t, 3, len(point))
		assert.Equal(t, 255.0, point[0])
		assert.Equal(t, 255.0, point[1])
		assert.Equal(t, 255.0, point[2])

		return err
	}, nil, nil)
}

func TestImage_GetPoint_WithAlpha(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		point, err := img.GetPoint(10, 10)

		assert.Equal(t, 4, len(point))
		assert.Equal(t, 0.0, point[0])
		assert.Equal(t, 0.0, point[1])
		assert.Equal(t, 0.0, point[2])
		assert.Equal(t, 255.0, point[3])

		return err
	}, nil, nil)
}

func TestImage_GetPoint_WithAlpha2(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		point, err := img.GetPoint(0, 0)

		assert.Equal(t, 4, len(point))
		assert.Equal(t, 0.0, point[0])
		assert.Equal(t, 0.0, point[1])
		assert.Equal(t, 0.0, point[2])
		assert.Equal(t, 0.0, point[3])

		return err
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
		assert.NoError(t, err)

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

func TestImage_Decode_PNG(t *testing.T) {
	goldenTest(t, resources+"png-8bit.png", func(img *ImageRef) error {
		goImg, err := img.ToImage(nil)
		assert.NoError(t, err)

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
	buf, metadata, err := i.ExportNative()
	require.NoError(t, err)

	assertGoldenMatch(t, resources+"jpg-24bit.jpg", buf, metadata.Format)
}

// vips jpegsave resources/jpg-24bit-icc-iec.jpg test.jpg --Q=75 --profile=none --strip --subsample-mode=auto --interlace --optimize-coding
func TestImage_OptimizeCoding(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		nil,
		nil,
		exportJpeg(&JpegExportParams{
			SubsampleMode:  VipsForeignSubsampleAuto,
			StripMetadata:  true,
			Quality:        75,
			Interlace:      true,
			OptimizeCoding: true,
		}),
	)
}

// vips jpegsave resources/jpg-24bit-icc-iec.jpg test.jpg --Q=75 --profile=none --strip --subsample-mode=on
func TestImage_SubsampleMode(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		nil,
		nil,
		exportJpeg(&JpegExportParams{
			SubsampleMode: VipsForeignSubsampleOn,
			StripMetadata: true,
			Quality:       75,
		}),
	)
}

// vips jpegsave resources/jpg-24bit-icc-iec.jpg test.jpg --Q=75 --profile=none --strip --trellis-quant
func TestImage_TrellisQuant(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		nil,
		nil,
		exportJpeg(&JpegExportParams{
			SubsampleMode: VipsForeignSubsampleAuto,
			StripMetadata: true,
			Quality:       75,
			TrellisQuant:  true,
		}),
	)
}

// vips jpegsave resources/jpg-24bit-icc-iec.jpg test.jpg --Q=75 --profile=none --strip --overshoot-deringing
func TestImage_OvershootDeringing(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		nil,
		nil,
		exportJpeg(&JpegExportParams{
			SubsampleMode:      VipsForeignSubsampleAuto,
			StripMetadata:      true,
			Quality:            75,
			OvershootDeringing: true,
		}),
	)
}

// vips jpegsave resources/jpg-24bit-icc-iec.jpg test.jpg --Q=75 --profile=none --strip --interlace --optimize-scans
func TestImage_OptimizeScans(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		nil,
		nil,
		exportJpeg(&JpegExportParams{
			SubsampleMode: VipsForeignSubsampleAuto,
			StripMetadata: true,
			Quality:       75,
			Interlace:     true,
			OptimizeScans: true,
		}),
	)
}

// vips jpegsave resources/jpg-24bit-icc-iec.jpg test.jpg --Q=75 --profile=none --strip --quant-table=3
func TestImage_QuantTable(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		nil,
		nil,
		exportJpeg(&JpegExportParams{
			SubsampleMode: VipsForeignSubsampleAuto,
			StripMetadata: true,
			Quality:       75,
			QuantTable:    3,
		}),
	)
}

func TestImage_Pixelate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg",
		func(img *ImageRef) error {
			return Pixelate(img, 24)
		},
		nil, nil)
}

func testWebpOptimizeIccProfile(t *testing.T, exportParams *WebpExportParams) []byte {
	return goldenTest(t, resources+"has-icc-profile.png",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		func(result *ImageRef) {
			assert.True(t, result.HasICCProfile(), "should have an ICC profile")
		},
		exportWebp(exportParams),
	)
}

func exportWebp(exportParams *WebpExportParams) func(img *ImageRef) ([]byte, *ImageMetadata, error) {
	return func(img *ImageRef) ([]byte, *ImageMetadata, error) {
		return img.ExportWebp(exportParams)
	}
}

func exportJpeg(exportParams *JpegExportParams) func(img *ImageRef) ([]byte, *ImageMetadata, error) {
	return func(img *ImageRef) ([]byte, *ImageMetadata, error) {
		return img.ExportJpeg(exportParams)
	}
}

func exportAvif(exportParams *AvifExportParams) func(img *ImageRef) ([]byte, *ImageMetadata, error) {
	return func(img *ImageRef) ([]byte, *ImageMetadata, error) {
		return img.ExportAvif(exportParams)
	}
}

func exportPng(exportParams *PngExportParams) func(img *ImageRef) ([]byte, *ImageMetadata, error) {
	return func(img *ImageRef) ([]byte, *ImageMetadata, error) {
		return img.ExportPng(exportParams)
	}
}

func exportGif(exportParams *GifExportParams) func(img *ImageRef) ([]byte, *ImageMetadata, error) {
	return func(img *ImageRef) ([]byte, *ImageMetadata, error) {
		return img.ExportGIF(exportParams)
	}
}

func goldenTest(
	t *testing.T,
	path string,
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

	img, err := NewImageFromFile(path)
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

func goldenCreateTest(
	t *testing.T,
	path string,
	createFromFile func(path string) (*ImageRef, error),
	createFromBuffer func(buf []byte) (*ImageRef, error),
	exec func(img *ImageRef) error,
	validate func(img *ImageRef),
	export func(img *ImageRef) ([]byte, *ImageMetadata, error),
) []byte {
	if createFromFile == nil {
		createFromFile = NewImageFromFile
	}
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

	img, err := createFromFile(path)
	require.NoError(t, err)

	err = exec(img)
	require.NoError(t, err)

	buf, metadata, err := export(img)
	require.NoError(t, err)

	result, err := NewImageFromBuffer(buf)
	require.NoError(t, err)

	validate(result)

	assertGoldenMatch(t, path, buf, metadata.Format)

	buf2, err := os.ReadFile(path)
	require.NoError(t, err)

	img2, err := createFromBuffer(buf2)
	require.NoError(t, err)

	err = exec(img2)
	require.NoError(t, err)

	buf2, metadata2, err := export(img2)
	require.NoError(t, err)

	result2, err := NewImageFromBuffer(buf2)
	require.NoError(t, err)

	validate(result2)

	assertGoldenMatch(t, path, buf2, metadata2.Format)

	return buf
}

func getEnvironment() string {
	sanitizedVersion := strings.ReplaceAll(Version, ":", "-")
	switch runtime.GOOS {
	case "windows":
		// Missing Windows version detection. Windows is not a supported CI target right now
		return "windows_" + runtime.GOARCH + "_libvips-" + sanitizedVersion
	case "darwin":
		out, err := exec.Command("sw_vers", "-productVersion").Output()
		if err != nil {
			return "macos-unknown_" + runtime.GOARCH + "_libvips-" + sanitizedVersion
		}
		majorVersion := strings.Split(strings.TrimSpace(string(out)), ".")[0]
		return "macos-" + majorVersion + "_" + runtime.GOARCH + "_libvips-" + sanitizedVersion
	case "linux":
		out, err := exec.Command("lsb_release", "-cs").Output()
		if err != nil {
			return "linux-unknown_" + runtime.GOARCH
		}
		strout := strings.TrimSuffix(string(out), "\n")
		return "linux-" + strout + "_" + runtime.GOARCH + "_libvips-" + sanitizedVersion
	}
	// default to unknown assets otherwise
	return "unknown_" + runtime.GOARCH + "_libvips-" + sanitizedVersion
}

func assertGoldenMatch(t *testing.T, file string, buf []byte, format ImageType) {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		panic("bad filename")
	}

	name := strings.Replace(t.Name(), "/", "_", -1)
	name = strings.Replace(name, "TestImage_", "", -1)
	prefix := file[:i] + "." + name
	ext := format.FileExt()
	goldenFile := prefix + "-" + getEnvironment() + ".golden" + ext

	golden, _ := os.ReadFile(goldenFile)
	if golden != nil {
		sameAsGolden := assert.True(t, bytes.Equal(buf, golden), "Actual image (size=%d) didn't match expected golden file=%s (size=%d)", len(buf), goldenFile, len(golden))
		if !sameAsGolden {
			failed := prefix + "-" + getEnvironment() + ".failed" + ext
			err := os.WriteFile(failed, buf, 0666)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	t.Log("writing golden file: " + goldenFile)
	err := os.WriteFile(goldenFile, buf, 0644)
	assert.NoError(t, err)
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
