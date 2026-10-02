package vips

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveKeep(t *testing.T) {
	allButGainmap := int(ForeignKeepAll &^ ForeignKeepGainmap)

	tests := []struct {
		name         string
		keep         ForeignKeep
		major, minor int
		wantFlags    int
		wantSet      bool
		wantErr      error // nil, errKeepUnsupported, or errAny
	}{
		{"unset on 8.14", 0, 8, 14, 0, false, nil},
		{"unset on 8.18", 0, 8, 18, 0, false, nil},
		{"set on 8.14", ForeignKeepIcc, 8, 14, 0, false, errKeepUnsupported},
		{"set on 7.x", ForeignKeepIcc, 7, 99, 0, false, errKeepUnsupported},
		{"none maps to zero", ForeignKeepNone, 8, 15, 0, true, nil},
		{"none ored with icc", ForeignKeepNone | ForeignKeepIcc, 8, 15, int(ForeignKeepIcc), true, nil},
		{"icc", ForeignKeepIcc, 8, 15, int(ForeignKeepIcc), true, nil},
		{"unknown bit", ForeignKeep(1 << 6), 8, 18, 0, false, errAny},
		{"all on 8.17 drops gainmap", ForeignKeepAll, 8, 17, allButGainmap, true, nil},
		{"all on 8.18", ForeignKeepAll, 8, 18, int(ForeignKeepAll), true, nil},
		{"all on 9.0", ForeignKeepAll, 9, 0, int(ForeignKeepAll), true, nil},
		{"gainmap only on 8.17", ForeignKeepGainmap, 8, 17, 0, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, set, err := resolveKeep(tt.keep, tt.major, tt.minor)
			switch {
			case tt.wantErr == errAny:
				require.Error(t, err)
				assert.False(t, errors.Is(err, errKeepUnsupported))
			case tt.wantErr != nil:
				// testify v1.6.1 (go.mod) has no ErrorIs; repo tests use errors.Is.
				require.True(t, errors.Is(err, tt.wantErr), "want %v, got: %v", tt.wantErr, err)
			default:
				require.NoError(t, err)
				assert.Equal(t, tt.wantFlags, flags)
				assert.Equal(t, tt.wantSet, set)
			}
		})
	}
}

// errAny marks table rows that expect some error other than errKeepUnsupported.
var errAny = errors.New("any error")

const keepTestSource = "jpg-24bit-icc-iec.jpg" // has exif-data, xmp-data, icc-profile-data

type keepMeta struct {
	ICC, Exif, XMP bool
}

func keepSupported() bool {
	return MajorVersion > 8 || (MajorVersion == 8 && MinorVersion >= 15)
}

func loadKeepSource(t *testing.T) *ImageRef {
	t.Helper()
	img, err := NewImageFromFile(resources + keepTestSource)
	require.NoError(t, err)
	return img
}

func readKeepMeta(t *testing.T, buf []byte) keepMeta {
	t.Helper()
	img, err := NewImageFromBuffer(buf)
	require.NoError(t, err)
	defer img.Close()
	m := keepMeta{ICC: img.HasICCProfile(), Exif: img.HasExif()}
	for _, f := range img.ImageFields() {
		if f == "xmp-data" {
			m.XMP = true
		}
	}
	return m
}

type keepExporter struct {
	name   string
	export func(img *ImageRef, keep ForeignKeep, strip bool) ([]byte, error)
	// all and unset are the metadata the format actually round-trips
	// (measured on libvips 8.18.0, spec F11): TIFF drops EXIF, and unset
	// WebP drops ICC because govips passes profile="none" (spec F5).
	all, unset keepMeta
}

var everything = keepMeta{ICC: true, Exif: true, XMP: true}

var keepExporters = []keepExporter{
	{
		name: "jpeg",
		export: func(img *ImageRef, keep ForeignKeep, strip bool) ([]byte, error) {
			p := NewJpegExportParams()
			p.StripMetadata, p.Keep = strip, keep
			b, _, err := img.ExportJpeg(p)
			return b, err
		},
		all: everything, unset: everything,
	},
	{
		name: "png",
		export: func(img *ImageRef, keep ForeignKeep, strip bool) ([]byte, error) {
			p := NewPngExportParams()
			p.StripMetadata, p.Keep = strip, keep
			b, _, err := img.ExportPng(p)
			return b, err
		},
		all: everything, unset: everything,
	},
	{
		name: "tiff",
		export: func(img *ImageRef, keep ForeignKeep, strip bool) ([]byte, error) {
			p := NewTiffExportParams()
			p.StripMetadata, p.Keep = strip, keep
			b, _, err := img.ExportTiff(p)
			return b, err
		},
		all:   keepMeta{ICC: true, XMP: true},
		unset: keepMeta{ICC: true, XMP: true},
	},
}

func runKeepMetadataCases(t *testing.T, exporters []keepExporter) {
	require.NoError(t, Startup(nil))
	if !keepSupported() {
		t.Skipf("keep requires libvips 8.15+, found %s", Version)
	}
	before := OpenImageRefs()
	img := loadKeepSource(t)

	for _, e := range exporters {
		t.Run(e.name, func(t *testing.T) {
			cases := []struct {
				name  string
				keep  ForeignKeep
				strip bool
				want  keepMeta
			}{
				{"icc", ForeignKeepIcc, false, keepMeta{ICC: true}},
				{"none", ForeignKeepNone, false, keepMeta{}},
				{"all", ForeignKeepAll, false, e.all},
				{"icc overrides strip", ForeignKeepIcc, true, keepMeta{ICC: true}},
				{"unset", 0, false, e.unset},
				{"unset strip", 0, true, keepMeta{}},
			}
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					buf, err := e.export(img, c.keep, c.strip)
					require.NoError(t, err)
					assert.Equal(t, c.want, readKeepMeta(t, buf))
				})
			}
		})
	}

	img.Close()
	assertNoNewImageRefs(t, before)
}

func TestExportKeep_Metadata(t *testing.T) {
	runKeepMetadataCases(t, keepExporters)
}

// keepCall exports img with the given Keep through one public save path.
type keepCall struct {
	name string
	save func(img *ImageRef, keep ForeignKeep) error
}

// keepCalls lists every public save path that accepts Keep. Tasks 3 and 4
// append the WebP, AVIF, JP2K, JXL and Magick entries.
var keepCalls = []keepCall{
	{"ExportJpeg", func(img *ImageRef, k ForeignKeep) error {
		p := NewJpegExportParams()
		p.Keep = k
		_, _, err := img.ExportJpeg(p)
		return err
	}},
	{"ExportPng", func(img *ImageRef, k ForeignKeep) error {
		p := NewPngExportParams()
		p.Keep = k
		_, _, err := img.ExportPng(p)
		return err
	}},
	{"ExportTiff", func(img *ImageRef, k ForeignKeep) error {
		p := NewTiffExportParams()
		p.Keep = k
		_, _, err := img.ExportTiff(p)
		return err
	}},
	{"ExportGIF", func(img *ImageRef, k ForeignKeep) error {
		p := NewGifExportParams()
		p.Keep = k
		_, _, err := img.ExportGIF(p)
		return err
	}},
	{"ExportHeif", func(img *ImageRef, k ForeignKeep) error {
		p := NewHeifExportParams()
		p.Keep = k
		_, _, err := img.ExportHeif(p)
		return err
	}},
	{"SaveToWriterJpeg", func(img *ImageRef, k ForeignKeep) error {
		p := NewJpegExportParams()
		p.Keep = k
		return img.SaveToWriterJpeg(io.Discard, p)
	}},
	{"SaveToWriterPng", func(img *ImageRef, k ForeignKeep) error {
		p := NewPngExportParams()
		p.Keep = k
		return img.SaveToWriterPng(io.Discard, p)
	}},
	{"SaveToWriterTiff", func(img *ImageRef, k ForeignKeep) error {
		p := NewTiffExportParams()
		p.Keep = k
		return img.SaveToWriterTiff(io.Discard, p)
	}},
	{"SaveToWriterGif", func(img *ImageRef, k ForeignKeep) error {
		p := NewGifExportParams()
		p.Keep = k
		return img.SaveToWriterGif(io.Discard, p)
	}},
	{"SaveToWriterHeif", func(img *ImageRef, k ForeignKeep) error {
		p := NewHeifExportParams()
		p.Keep = k
		return img.SaveToWriterHeif(io.Discard, p)
	}},
}

// TestExportKeep_AllPathsHonorKeep proves every public save path reaches
// applyKeep: on libvips < 8.15 each must return the version error; on
// newer libvips each must save with Keep set. Formats this environment
// can't save at all (probed with Keep unset) are skipped.
func TestExportKeep_AllPathsHonorKeep(t *testing.T) {
	require.NoError(t, Startup(nil))
	before := OpenImageRefs()
	img := loadKeepSource(t)

	for _, c := range keepCalls {
		t.Run(c.name, func(t *testing.T) {
			if err := c.save(img, 0); err != nil {
				t.Skipf("%s unsupported in this environment: %v", c.name, err)
			}
			err := c.save(img, ForeignKeepIcc)
			if keepSupported() {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "8.15+")
			}
		})
	}

	img.Close()
	assertNoNewImageRefs(t, before)
}
