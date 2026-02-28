package vips

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func skipIfHeifSaveUnsupported(t *testing.T) {
	t.Helper()
	require.NoError(t, Startup(nil))
	img, err := Black(1, 1)
	require.NoError(t, err)
	_, _, err = img.ExportHeif(NewHeifExportParams())
	if err != nil {
		t.Skip("HEIF save is not supported in this environment")
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

	require.NoError(t, Startup(nil))

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

	require.NoError(t, Startup(nil))

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

	require.NoError(t, Startup(nil))

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
