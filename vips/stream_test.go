package vips

import (
	"bytes"
	"fmt"
	"io"
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
		name         string
		files        []string
		reader       func(string) io.Reader
		expSourceErr error
		expLoadErr   bool
	}{
		{
			name:  "nil reader",
			files: []string{"empty"},
			reader: func(name string) io.Reader {
				return nil
			},
			expSourceErr: io.EOF,
		},
		{
			name:  "bad data",
			files: []string{"empty"},
			reader: func(name string) io.Reader {
				return bytes.NewReader([]byte("not an image"))
			},
			expLoadErr: true,
		},
		{
			name:  "good data",
			files: allImageTypeFiles,
			reader: func(name string) io.Reader {
				file, err := os.Open(resources + name)
				require.NoError(t, err)
				return file
			},
		},
		{
			name:  "incomplete data",
			files: allImageTypeFiles,
			reader: func(name string) io.Reader {
				file, _ := os.Open(resources + name)
				reader := io.LimitReader(file, 128)
				return reader
			},
			expLoadErr: true,
		},
	}

	for _, test := range tests {
		for _, file := range test.files {
			t.Run(fmt.Sprintf("[%s] %s", file, test.name), func(t *testing.T) {

				reader := test.reader(file)

				source, err := NewSourceFromReader(reader)
				assert.Equal(t, test.expSourceErr, err)

				if test.expSourceErr != nil {
					require.Nil(t, source, nil)
					return
				}

				require.NotNil(t, source)

				img, err := NewImageFromSource(source)
				if test.expLoadErr {
					assert.NotNil(t, err)
					require.Nil(t, img)
					return
				}

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
			name:       "file writer target",
			imageTypes: allSupportedTypes,
			target: func() *Target {
				file, err := os.CreateTemp("", "test")
				require.NoError(t, err)

				target, err := NewTargetToWriter(file)
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
		{
			name:       "buffer target - no tiff",
			imageTypes: allTypesExceptTiff,
			target: func() *Target {
				target, err := NewTargetToWriter(&bytes.Buffer{})
				require.NoError(t, err)
				return target
			},
		},
		{
			name:       "buffer target - tiff",
			imageTypes: []ImageType{ImageTypeTIFF},
			target: func() *Target {
				target, err := NewTargetToWriter(&bytes.Buffer{})
				require.NoError(t, err)
				return target
			},
			expExportError: true,
		},
	}

	for _, test := range tests {
		for _, format := range test.imageTypes {

			t.Run(fmt.Sprintf("[%s] %s", format, test.name), func(t *testing.T) {
				img, err := NewImageFromFile(resources + "png-24bit.png")
				assert.NoError(t, err)
				assert.NotNil(t, img)

				target := test.target()
				if file, ok := target.writer.(*os.File); ok {
					defer os.Remove(file.Name())
				}

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
