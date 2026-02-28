package vips

import (
	"bytes"
	"image"
	jpeg2 "image/jpeg"
	"testing"

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

func TestImage_Gravity(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.Gravity(GravityNorthWest, 500, 500)
		}, nil, nil)
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

func TestImage_Modulate(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(img *ImageRef) error {
		return img.Modulate(0.7, 0.5, 180)
	}, nil, nil)
}

func TestImage_ModulateHSV(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit-icc-iec.jpg", func(img *ImageRef) error {
		return img.ModulateHSV(0.7, 0.5, 180)
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
