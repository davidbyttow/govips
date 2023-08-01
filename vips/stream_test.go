package vips

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSource(t *testing.T) {
	Startup(nil)

	allImageTypeFiles := []string{
		"avif-8bit.avif",
		"jpg-24bit.jpg",
		"bmp.bmp",
		"large.bmp",
		"gif-animated.gif",
		"png-24bit.png",
		"svg.svg",
		"svg_1.svg",
		"svg_2.svg",
		"tif.tif",
		"heic-24bit.heic",
		"webp-animated.webp",
		"jp2k-orientation-6.jp2",
	}
	tests := []struct {
		name       string
		files      []string
		source     func(string) *Source
		expLoadErr bool
	}{
		{
			name:  "pipe",
			files: allImageTypeFiles,
			source: func(path string) *Source {
				file, err := os.Open(resources + path)
				require.NoError(t, err)
				source, err := NewSourceFromPipe(file)
				require.NoError(t, err)
				return source
			},
		},
		{
			name:  "file path",
			files: allImageTypeFiles,
			source: func(path string) *Source {
				source, err := NewSourceFromFile(resources + path)
				require.NoError(t, err)
				return source
			},
		},
		{
			name:  "incomplete data",
			files: allImageTypeFiles,
			source: func(path string) *Source {
				file, err := os.Open(resources + path)
				require.NoError(t, err)
				fi, err := file.Stat()
				require.NoError(t, err)
				file.Seek(0, int(fi.Size()/2))
				source, err := NewSourceFromPipe(file)
				require.NoError(t, err)
				return source
			},
			expLoadErr: true,
		},
	}

	for _, test := range tests {
		for _, file := range test.files {
			t.Run(fmt.Sprintf("[%s] %s", file, test.name), func(t *testing.T) {

				source := test.source(file)

				require.NotNil(t, source)

				img, err := NewImageFromSource(source)

				require.NoError(t, err)
				require.NotNil(t, img)
				_, _, err = img.ExportNative()
				require.NoError(t, err)

			})
		}
	}
}

func TestTarget(t *testing.T) {
	Startup(nil)

	var allSupportedTypes []ImageType
	var allTypesExceptTiff []ImageType // tiff cannot be written to a writer that is not ReadSeeker (i.e. not a file)

	for t, ok := range supportedImageTypes {
		if !ok {
			continue
		}
		allSupportedTypes = append(allSupportedTypes, t)
		if t != ImageTypeTIFF {
			allTypesExceptTiff = append(allTypesExceptTiff, t)
		}
	}

	tests := []struct {
		name           string
		imageTypes     []ImageType
		target         func() *Target
		expErr         error
		expExportError bool
	}{
		{
			name:       "pipe target",
			imageTypes: allSupportedTypes,
			target: func() *Target {
				file, err := os.CreateTemp("", "test")
				require.NoError(t, err)

				target, err := NewTargetToPipe(file)
				require.NoError(t, err)
				return target
			},
		},
		{
			name:       "file path target",
			imageTypes: allSupportedTypes,
			target: func() *Target {
				file, err := os.CreateTemp("", "test")
				require.NoError(t, err)

				target, err := NewTargetToFile(file.Name())
				require.NoError(t, err)
				return target
			},
		},
	}

	for _, test := range tests {
		for _, format := range test.imageTypes {

			t.Run(fmt.Sprintf("[%s] %s", format, test.name), func(t *testing.T) {
				file, err := os.Open(resources + "png-24bit.png")
				require.NoError(t, err)
				buf, err := ioutil.ReadAll(file)
				require.NoError(t, err)

				img, err := NewImageFromBuffer(buf)
				assert.NoError(t, err)
				assert.NotNil(t, img)

				target := test.target()

				assert.NotNil(t, target)

				params := NewDefaultExportParams()
				params.Format = format

				_, err = img.ExportTarget(target, params)

				if test.expExportError {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			})
		}
	}
}
