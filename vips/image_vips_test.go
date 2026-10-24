package vips

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireVipsCLI(t *testing.T) {
	t.Helper()
	for _, tool := range []string{"vips", "vipsheader"} {
		if _, err := exec.LookPath(tool); err != nil {
			if os.Getenv("CI") != "" {
				t.Fatalf("native VIPS round-trip tests require %s: %v", tool, err)
			}
			t.Skipf("native VIPS round-trip tests require %s (install libvips-tools): %v", tool, err)
		}
	}
}

func vipsCLI(t *testing.T, tool string, args ...string) string {
	t.Helper()
	out, err := exec.Command(tool, args...).CombinedOutput()
	require.NoError(t, err, "%s %v: %s", tool, args, out)
	return strings.TrimSpace(string(out))
}

func vipsSnapshotImage(t *testing.T) *ImageRef {
	t.Helper()
	require.NoError(t, Startup(nil))
	before := OpenImageRefs()
	img, err := NewImageFromMemory([]byte{101, 51, 25, 128, 10, 200, 255, 64}, 2, 1, 4, BandFormatUchar, InterpretationSRGB)
	require.NoError(t, err)
	t.Cleanup(func() {
		img.Close()
		assertNoNewImageRefs(t, before)
	})
	return img
}

func assertVipsSnapshot(t *testing.T, path string, img *ImageRef) {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	// Native VIPS files can use either byte order.
	assert.True(t, bytes.HasPrefix(data, []byte{0xb6, 0xa6, 0xf2, 0x08}) ||
		bytes.HasPrefix(data, []byte{0x08, 0xf2, 0xa6, 0xb6}), "expected native VIPS signature")
	for field, want := range map[string]string{
		"width": strconv.Itoa(img.Width()), "height": strconv.Itoa(img.Height()),
		"bands": strconv.Itoa(img.Bands()),
	} {
		assert.Equal(t, want, vipsCLI(t, "vipsheader", "-f", field, path), field)
	}
	for x := 0; x < img.Width(); x++ {
		want, err := img.GetPoint(x, 0)
		require.NoError(t, err)
		values := strings.Fields(vipsCLI(t, "vips", "getpoint", path, strconv.Itoa(x), "0"))
		require.Len(t, values, len(want))
		for band, value := range values {
			got, err := strconv.ParseFloat(value, 64)
			require.NoError(t, err)
			assert.InDelta(t, want[band], got, 0.001, "x=%d band=%d", x, band)
		}
	}
}

func TestImageRef_SaveVips_Pixels(t *testing.T) {
	requireVipsCLI(t)
	for _, kind := range []string{"rgba", "premultiplied", "extra-band"} {
		t.Run(kind, func(t *testing.T) {
			img := vipsSnapshotImage(t)
			format := "VIPS_FORMAT_UCHAR"
			switch kind {
			case "premultiplied":
				require.NoError(t, img.PremultiplyAlpha())
				format = "VIPS_FORMAT_FLOAT"
			case "extra-band":
				require.NoError(t, img.BandJoinConst([]float64{37}))
			}
			bandFormat := img.BandFormat()
			imageFormat, originalFormat := img.Format(), img.OriginalFormat()
			path := filepath.Join(t.TempDir(), "snapshot.v")
			require.NoError(t, img.SaveVips(path))
			assert.Contains(t, vipsCLI(t, "vipsheader", "-f", "format", path), format)
			assert.Contains(t, vipsCLI(t, "vipsheader", "-f", "interpretation", path), "VIPS_INTERPRETATION_sRGB")
			assertVipsSnapshot(t, path, img)
			assert.Equal(t, bandFormat, img.BandFormat())
			assert.Equal(t, imageFormat, img.Format())
			assert.Equal(t, originalFormat, img.OriginalFormat())
			if kind == "premultiplied" {
				// Saving must retain govips's state for later unpremultiplication.
				require.NoError(t, img.UnpremultiplyAlpha())
				assert.Equal(t, BandFormatUchar, img.BandFormat())
				point, err := img.GetPoint(0, 0)
				require.NoError(t, err)
				assert.InDeltaSlice(t, []float64{101, 51, 25, 128}, point, 1)
			}
		})
	}
}

func TestImageRef_SaveVips_Metadata(t *testing.T) {
	requireVipsCLI(t)
	img := vipsSnapshotImage(t)
	require.NoError(t, img.SetOrientation(6))
	img.SetString("checkpoint-note", "native snapshot")
	blob := []byte{0, 1, 2, 127, 128, 255}
	img.SetBlob("checkpoint-blob", blob)
	path := filepath.Join(t.TempDir(), "metadata.vips")
	require.NoError(t, img.SaveVips(path))
	assert.Equal(t, "6", vipsCLI(t, "vipsheader", "-f", "orientation", path))
	assert.Equal(t, "native snapshot", vipsCLI(t, "vipsheader", "-f", "checkpoint-note", path))
	got, err := base64.StdEncoding.DecodeString(vipsCLI(t, "vipsheader", "-f", "checkpoint-blob", path))
	require.NoError(t, err)
	assert.Equal(t, blob, got)
	assert.Equal(t, blob, img.GetBlob("checkpoint-blob"))
}

func TestImageRef_SaveVips_Paths(t *testing.T) {
	requireVipsCLI(t)
	img := vipsSnapshotImage(t)
	for _, name := range []string{"snapshot.v", "snapshot.vips", "snapshot.V", "snapshot.ViPs", "快照 with spaces.v"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			// A longer existing file checks that overwrite also truncates it.
			require.NoError(t, os.WriteFile(path, bytes.Repeat([]byte{255}, 65536), 0600))
			require.NoError(t, img.SaveVips(path))
			assertVipsSnapshot(t, path, img)
			info, err := os.Stat(path)
			require.NoError(t, err)
			assert.Less(t, info.Size(), int64(65536))
		})
	}
}

func TestImageRef_SaveVips_InvalidPaths(t *testing.T) {
	img := vipsSnapshotImage(t)
	dir := t.TempDir()
	for _, path := range []string{"", filepath.Join(dir, "x.v") + "\x00.v", filepath.Join(dir, "x.png"), filepath.Join(dir, "x"), filepath.Join(dir, "x.v[strip]")} {
		require.Error(t, img.SaveVips(path), "path=%q", path)
	}
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Empty(t, entries, "invalid paths must not create files")

	missing := filepath.Join(dir, "missing", "x.v")
	require.Error(t, img.SaveVips(missing))
	_, err = os.Stat(filepath.Dir(missing))
	assert.True(t, os.IsNotExist(err), "missing parent must not be created")
	directory := filepath.Join(dir, "directory.v")
	require.NoError(t, os.Mkdir(directory, 0700))
	require.Error(t, img.SaveVips(directory))
	// A C-side save error must leave the input usable.
	_, err = img.GetPoint(0, 0)
	require.NoError(t, err)

	img.Close()
	path := filepath.Join(dir, "closed.v")
	require.Error(t, img.SaveVips(path))
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}

// snapshotReader can fail after the lazy loader has read the image header.
type snapshotReader struct {
	*bytes.Reader
	err error
}

func (r *snapshotReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.Reader.Read(p)
}

func TestImageRef_SaveVips_Sequential(t *testing.T) {
	// An uncompressed PNG is larger than the loader's header buffer, ensuring
	// that saving has to read from the source after the lazy load returns.
	pixels := image.NewNRGBA(image.Rect(0, 0, 256, 256))
	for i := range pixels.Pix {
		pixels.Pix[i] = byte(i*37 + i/1024)
	}
	var encoded bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.NoCompression}
	require.NoError(t, encoder.Encode(&encoded, pixels))
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprintf("reader-error=%t", fail), func(t *testing.T) {
			if !fail {
				requireVipsCLI(t)
			}
			before := OpenImageRefs()
			sourcesBefore, targetsBefore := streamRegistrySizes()
			reader := &snapshotReader{Reader: bytes.NewReader(encoded.Bytes())}
			params := NewImportParams()
			params.Access.Set(AccessSequential)
			img, err := LoadImageFromReader(&nonSeekable{r: reader}, params)
			require.NoError(t, err)
			t.Cleanup(func() {
				img.Close()
				assertNoNewImageRefs(t, before)
				sources, targets := streamRegistrySizes()
				assert.Equal(t, sourcesBefore, sources)
				assert.Equal(t, targetsBefore, targets)
			})
			readerErr := errors.New("snapshot source failed")
			if fail {
				reader.err = readerErr
			}
			path := filepath.Join(t.TempDir(), "sequential.v")
			err = img.SaveVips(path)
			if fail {
				require.Error(t, err)
				assert.True(t, errors.Is(err, readerErr), "source error must be wrapped: %v", err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "256", vipsCLI(t, "vipsheader", "-f", "height", path))
				assert.Equal(t, "0 37 74 111", vipsCLI(t, "vips", "getpoint", path, "0", "0"))
			}
		})
	}
}
