package vips

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nonSeekable hides the io.Seeker implementation of the wrapped reader,
// forcing libvips into sequential (pipe) mode.
type nonSeekable struct {
	r io.Reader
}

func (n *nonSeekable) Read(p []byte) (int, error) {
	return n.r.Read(p)
}

// errAfterReader returns errAt-prefixed bytes, then fails.
type errAfterReader struct {
	data []byte
	pos  int
	err  error
}

func (e *errAfterReader) Read(p []byte) (int, error) {
	if e.pos >= len(e.data) {
		return 0, e.err
	}
	n := copy(p, e.data[e.pos:])
	e.pos += n
	return n, nil
}

// errWriter fails on the first write.
type errWriter struct {
	err error
}

func (e *errWriter) Write(p []byte) (int, error) {
	return 0, e.err
}

func streamRegistrySizes() (int, int) {
	streamCallbacks.Lock()
	defer streamCallbacks.Unlock()
	return len(streamCallbacks.sources), len(streamCallbacks.targets)
}

// assertNoNewImageRefs verifies a test did not leak ImageRefs. It checks
// the delta rather than calling AssertNoLeaks because other tests in the
// suite legitimately rely on GC finalizers for cleanup, making the global
// counter nonzero at arbitrary times.
func assertNoNewImageRefs(t *testing.T, before int64) {
	t.Helper()
	// LessOrEqual rather than Equal: a GC run during this test may
	// finalize images leaked by earlier tests, lowering the counter.
	assert.LessOrEqual(t, OpenImageRefs(), before, "test leaked ImageRef(s)")
}

var streamTestFiles = map[ImageType]string{
	ImageTypeJPEG: "jpg-24bit.jpg",
	ImageTypePNG:  "png-24bit.png",
	ImageTypeWEBP: "webp+alpha.webp",
	ImageTypeHEIF: "heic-24bit.heic",
	ImageTypeTIFF: "tif.tif",
	ImageTypeGIF:  "gif-animated.gif",
}

// --- Foundational: callback registry (T008) ---

func TestStreamRegistry_ConcurrentRegisterDeregister(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sh, _ := registerSource(bytes.NewReader([]byte("x")))
				th, _ := registerTarget(&bytes.Buffer{})
				assert.NotNil(t, lookupSource(sh))
				assert.NotNil(t, lookupTarget(th))
				deregisterSource(sh)
				deregisterTarget(th)
			}
		}()
	}
	wg.Wait()

	sources, targets := streamRegistrySizes()
	assert.Zero(t, sources)
	assert.Zero(t, targets)
}

func TestStreamRegistry_HandlesAreUnique(t *testing.T) {
	seen := make(map[int]bool)
	var handles []int
	for i := 0; i < 100; i++ {
		h, _ := registerSource(bytes.NewReader(nil))
		require.False(t, seen[h], "handle %d reused", h)
		seen[h] = true
		handles = append(handles, h)
	}
	for i := 1; i < len(handles); i++ {
		assert.Greater(t, handles[i], handles[i-1], "handles must be monotonically increasing")
	}
	for _, h := range handles {
		deregisterSource(h)
	}
}

func TestStreamRegistry_StaleHandleReturnsError(t *testing.T) {
	sh, _ := registerSource(bytes.NewReader([]byte("data")))
	th, _ := registerTarget(&bytes.Buffer{})
	deregisterSource(sh)
	deregisterTarget(th)

	buf := make([]byte, 16)
	assert.EqualValues(t, -1, sourceRead(sh, buf))
	assert.EqualValues(t, -1, sourceSeek(sh, 0, io.SeekStart))
	assert.EqualValues(t, -1, targetWrite(th, buf))
	assert.EqualValues(t, -1, targetEnd(th))
}

func TestStreamRegistry_HandleWraparound(t *testing.T) {
	streamCallbacks.Lock()
	saved := streamCallbacks.nextHandle
	streamCallbacks.nextHandle = math.MaxInt32 - 1
	streamCallbacks.Unlock()
	defer func() {
		streamCallbacks.Lock()
		streamCallbacks.nextHandle = saved
		streamCallbacks.Unlock()
	}()

	h1, _ := registerSource(bytes.NewReader(nil))
	h2, _ := registerSource(bytes.NewReader(nil))
	defer deregisterSource(h1)
	defer deregisterSource(h2)

	assert.Equal(t, math.MaxInt32, h1, "handle must stay within int32 range")
	assert.Less(t, h2, 10, "next handle must wrap to the start of the range")
	assert.NotEqual(t, h1, h2)
	assert.NotNil(t, lookupSource(h1))
	assert.NotNil(t, lookupSource(h2))
}

// strictSeeker fails any attempt to position beyond the end, like some
// custom io.Seeker implementations do. The seek callback must never ask
// it for a past-EOF position.
type strictSeeker struct {
	r *bytes.Reader
}

func (s *strictSeeker) Read(p []byte) (int, error) { return s.r.Read(p) }
func (s *strictSeeker) Seek(offset int64, whence int) (int64, error) {
	pos, err := s.r.Seek(offset, whence)
	if err == nil && pos > s.r.Size() {
		return 0, errors.New("strictSeeker: seek beyond EOF")
	}
	return pos, err
}

func TestStreamRegistry_SeekPastEOFIsPosix(t *testing.T) {
	data := []byte("0123456789")
	h, entry := registerSource(bytes.NewReader(data))
	defer deregisterSource(h)

	// Codecs probe past EOF (libheif box probing, see issue #3). POSIX
	// allows it: the callback reports the requested position — libvips
	// core range-checks it and signals the codec glue — while reads at
	// that position return 0 bytes.
	assert.EqualValues(t, 18, sourceSeek(h, 18, io.SeekStart), "past-EOF SEEK_SET must report the POSIX position")
	buf := make([]byte, 4)
	assert.EqualValues(t, 0, sourceRead(h, buf), "read past EOF must return 0 bytes")

	assert.EqualValues(t, 18, sourceSeek(h, 8, io.SeekEnd), "past-EOF SEEK_END must report the POSIX position")

	require.EqualValues(t, 4, sourceSeek(h, 4, io.SeekStart))
	assert.EqualValues(t, 104, sourceSeek(h, 100, io.SeekCurrent), "past-EOF SEEK_CUR must report the POSIX position")

	// Negative targets are rejected (POSIX EINVAL) without poisoning
	// the stream error state.
	assert.EqualValues(t, -1, sourceSeek(h, -1, io.SeekStart))
	assert.NoError(t, entry.takeErr(), "a rejected probe seek must not become a stream error")

	// Normal seeking still works afterwards.
	assert.EqualValues(t, 2, sourceSeek(h, 2, io.SeekStart))
	assert.EqualValues(t, 4, sourceRead(h, buf))
}

func TestStreamRegistry_SeekNeverAsksSeekerPastEOF(t *testing.T) {
	data := []byte("0123456789")
	h, entry := registerSource(&strictSeeker{r: bytes.NewReader(data)})
	defer deregisterSource(h)

	assert.EqualValues(t, 25, sourceSeek(h, 25, io.SeekStart),
		"the POSIX position is reported without invoking the seeker out of range")
	assert.NoError(t, entry.takeErr())
}

func TestLoadImageFromReader_HEICSeekProbing(t *testing.T) {
	require.NoError(t, Startup(nil))

	// libheif probes a few bytes past the final box of this fixture
	// (issue #3): the load must succeed with no "bad seek" in sight.
	buf, err := os.ReadFile(resources + "heic-24bit.heic")
	require.NoError(t, err)

	img, err := LoadImageFromReader(bytes.NewReader(buf), nil)
	require.NoError(t, err)
	defer img.Close()
	assert.Equal(t, ImageTypeHEIF, img.Format())
}

func TestStreamRegistry_SeekWithoutSeeker(t *testing.T) {
	h, _ := registerSource(&nonSeekable{r: bytes.NewReader([]byte("data"))})
	defer deregisterSource(h)

	assert.EqualValues(t, -1, sourceSeek(h, 0, io.SeekStart))
}

// --- User Story 1: LoadImageFromReader (T012, T013, T014) ---

func TestLoadImageFromReader_RoundtripAllFormats(t *testing.T) {
	require.NoError(t, Startup(nil))

	for format, file := range streamTestFiles {
		t.Run(ImageTypes[format]+"/"+file, func(t *testing.T) {
			path := resources + file

			f, err := os.Open(path)
			require.NoError(t, err)
			defer f.Close()

			streamed, err := LoadImageFromReader(f, nil)
			require.NoError(t, err)
			defer streamed.Close()

			buf, err := os.ReadFile(path)
			require.NoError(t, err)
			buffered, err := LoadImageFromBuffer(buf, nil)
			require.NoError(t, err)
			defer buffered.Close()

			assert.Equal(t, buffered.Width(), streamed.Width())
			assert.Equal(t, buffered.Height(), streamed.Height())
			assert.Equal(t, format, streamed.Format())
			assert.Nil(t, streamed.buf, "stream-loaded image must have no buffer backing")

			// The image must be fully usable after the source is released.
			require.NoError(t, streamed.Resize(0.5, KernelLanczos3))
			assert.Equal(t, (buffered.Width()+1)/2, streamed.Width())
		})
	}

	sources, _ := streamRegistrySizes()
	assert.Zero(t, sources, "all sources must be deregistered after load")
}

func TestLoadImageFromReader_SeekableAndSequential(t *testing.T) {
	require.NoError(t, Startup(nil))

	for _, file := range []string{"jpg-24bit.jpg", "heic-24bit.heic"} {
		t.Run("seekable/"+file, func(t *testing.T) {
			f, err := os.Open(resources + file)
			require.NoError(t, err)
			defer f.Close()

			img, err := LoadImageFromReader(f, nil)
			require.NoError(t, err)
			defer img.Close()
			assert.Greater(t, img.Width(), 0)
		})

		t.Run("sequential/"+file, func(t *testing.T) {
			buf, err := os.ReadFile(resources + file)
			require.NoError(t, err)

			img, err := LoadImageFromReader(&nonSeekable{r: bytes.NewReader(buf)}, nil)
			require.NoError(t, err)
			defer img.Close()
			assert.Greater(t, img.Width(), 0)
		})
	}
}

func TestLoadImageFromReader_ReaderErrorPropagates(t *testing.T) {
	require.NoError(t, Startup(nil))
	before := OpenImageRefs()

	buf, err := os.ReadFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	readerErr := errors.New("connection reset by peer")
	r := &errAfterReader{data: buf[:512], err: readerErr}

	img, err := LoadImageFromReader(&nonSeekable{r: r}, nil)
	if img != nil {
		img.Close()
	}
	require.Error(t, err)
	assert.True(t, errors.Is(err, readerErr), "original reader error must be wrapped, got: %v", err)

	sources, _ := streamRegistrySizes()
	assert.Zero(t, sources, "source must be deregistered after a failed load")
	assertNoNewImageRefs(t, before)
}

func TestLoadImageFromReader_EmptyStream(t *testing.T) {
	require.NoError(t, Startup(nil))
	before := OpenImageRefs()

	img, err := LoadImageFromReader(bytes.NewReader(nil), nil)
	if img != nil {
		img.Close()
	}
	require.Error(t, err)

	sources, _ := streamRegistrySizes()
	assert.Zero(t, sources)
	assertNoNewImageRefs(t, before)
}

func TestLoadImageFromReader_TruncatedInput(t *testing.T) {
	require.NoError(t, Startup(nil))
	before := OpenImageRefs()

	buf, err := os.ReadFile(resources + "png-24bit.png")
	require.NoError(t, err)

	img, err := LoadImageFromReader(bytes.NewReader(buf[:len(buf)/3]), nil)
	if img != nil {
		img.Close()
	}
	require.Error(t, err, "truncated input must fail, not produce a corrupt image")

	sources, _ := streamRegistrySizes()
	assert.Zero(t, sources)
	assertNoNewImageRefs(t, before)
}

func TestLoadImageFromReader_NilReader(t *testing.T) {
	img, err := LoadImageFromReader(nil, nil)
	require.Error(t, err)
	assert.Nil(t, img)
}

// --- User Story 2: SaveToWriter (T019, T020) ---

func TestSaveToWriter_ByteIdenticalToExport(t *testing.T) {
	require.NoError(t, Startup(nil))

	srcBuf, err := os.ReadFile(resources + "png-24bit.png")
	require.NoError(t, err)
	img, err := LoadImageFromBuffer(srcBuf, nil)
	require.NoError(t, err)
	defer img.Close()

	exports := map[ImageType]func() ([]byte, error){
		ImageTypeJPEG: func() ([]byte, error) { b, _, err := img.ExportJpeg(nil); return b, err },
		ImageTypePNG:  func() ([]byte, error) { b, _, err := img.ExportPng(nil); return b, err },
		ImageTypeWEBP: func() ([]byte, error) { b, _, err := img.ExportWebp(nil); return b, err },
		ImageTypeHEIF: func() ([]byte, error) { b, _, err := img.ExportHeif(nil); return b, err },
		ImageTypeTIFF: func() ([]byte, error) { b, _, err := img.ExportTiff(nil); return b, err },
		ImageTypeGIF:  func() ([]byte, error) { b, _, err := img.ExportGIF(nil); return b, err },
	}

	for format, export := range exports {
		t.Run(ImageTypes[format], func(t *testing.T) {
			if format == ImageTypeHEIF {
				skipIfHeifSaveUnsupported(t)
			}

			expected, err := export()
			require.NoError(t, err)

			var w bytes.Buffer
			require.NoError(t, img.SaveToWriter(&w, format, nil))

			assert.True(t, bytes.Equal(expected, w.Bytes()),
				"streaming output must be byte-identical to Export* (got %d bytes, want %d)",
				w.Len(), len(expected))

			// The output must re-load as a valid image of the saved format.
			reloaded, err := LoadImageFromBuffer(w.Bytes(), nil)
			require.NoError(t, err)
			defer reloaded.Close()
			assert.Equal(t, format, reloaded.Format())
		})
	}

	_, targets := streamRegistrySizes()
	assert.Zero(t, targets, "all targets must be deregistered after save")
}

// TestSaveToWriterTyped_ByteIdenticalToExport exercises format-specific
// options that the generic SaveToWriter/ExportParams path cannot express,
// verifying the typed savers stay byte-identical to their Export*
// counterparts.
func TestSaveToWriterTyped_ByteIdenticalToExport(t *testing.T) {
	require.NoError(t, Startup(nil))

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	t.Run("png palette", func(t *testing.T) {
		params := NewPngExportParams()
		params.Palette = true
		params.Bitdepth = 8
		expected, _, err := img.ExportPng(params)
		require.NoError(t, err)

		var w bytes.Buffer
		require.NoError(t, img.SaveToWriterPng(&w, params))
		assert.True(t, bytes.Equal(expected, w.Bytes()),
			"typed PNG streaming output must match ExportPng (got %d bytes, want %d)", w.Len(), len(expected))
	})

	t.Run("webp near-lossless", func(t *testing.T) {
		params := NewWebpExportParams()
		params.NearLossless = true
		expected, _, err := img.ExportWebp(params)
		require.NoError(t, err)

		var w bytes.Buffer
		require.NoError(t, img.SaveToWriterWebp(&w, params))
		assert.True(t, bytes.Equal(expected, w.Bytes()),
			"typed WebP streaming output must match ExportWebp (got %d bytes, want %d)", w.Len(), len(expected))
	})

	t.Run("tiff deflate", func(t *testing.T) {
		params := NewTiffExportParams()
		params.Compression = TiffCompressionDeflate
		expected, _, err := img.ExportTiff(params)
		require.NoError(t, err)

		var w bytes.Buffer
		require.NoError(t, img.SaveToWriterTiff(&w, params))
		assert.True(t, bytes.Equal(expected, w.Bytes()),
			"typed TIFF streaming output must match ExportTiff (got %d bytes, want %d)", w.Len(), len(expected))
	})

	_, targets := streamRegistrySizes()
	assert.Zero(t, targets, "all targets must be deregistered after save")
}

func TestSaveToWriter_GenericParamsMatchExport(t *testing.T) {
	require.NoError(t, Startup(nil))

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	params := &ExportParams{Quality: 85, StripMetadata: true, Interlaced: true}

	expected, _, err := img.Export(&ExportParams{
		Format: ImageTypeJPEG, Quality: 85, StripMetadata: true, Interlaced: true,
	})
	require.NoError(t, err)

	var w bytes.Buffer
	require.NoError(t, img.SaveToWriter(&w, ImageTypeJPEG, params))
	assert.True(t, bytes.Equal(expected, w.Bytes()),
		"generic ExportParams mapping must match (*ImageRef).Export")
}

func TestSaveToWriter_WriterErrorPropagates(t *testing.T) {
	require.NoError(t, Startup(nil))
	before := OpenImageRefs()

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)

	writerErr := errors.New("disk full")
	err = img.SaveToWriter(&errWriter{err: writerErr}, ImageTypeJPEG, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, writerErr), "original writer error must be wrapped, got: %v", err)

	_, targets := streamRegistrySizes()
	assert.Zero(t, targets, "target must be deregistered after a failed save")

	img.Close()
	assertNoNewImageRefs(t, before)
}

func TestSaveToWriter_ClosedImageRef(t *testing.T) {
	require.NoError(t, Startup(nil))

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	img.Close()

	var w bytes.Buffer
	err = img.SaveToWriter(&w, ImageTypeJPEG, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "closed")
}

func TestSaveToWriter_UnsupportedFormat(t *testing.T) {
	require.NoError(t, Startup(nil))

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	var w bytes.Buffer
	err = img.SaveToWriter(&w, ImageTypeBMP, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not support")
	assert.Zero(t, w.Len())
}

// --- User Story 3: end-to-end pipeline (T021-T024) ---

func TestStreamingPipeline_HEICToJPEG(t *testing.T) {
	require.NoError(t, Startup(nil))

	f, err := os.Open(resources + "heic-24bit.heic")
	require.NoError(t, err)
	defer f.Close()

	img, err := LoadImageFromReader(f, nil)
	require.NoError(t, err)
	defer img.Close()

	origWidth := img.Width()
	require.NoError(t, img.AutoRotate())
	require.NoError(t, img.Resize(0.5, KernelLanczos3))

	outPath := filepath.Join(t.TempDir(), "out.jpg")
	out, err := os.Create(outPath)
	require.NoError(t, err)
	require.NoError(t, img.SaveToWriter(out, ImageTypeJPEG, nil))
	require.NoError(t, out.Close())

	result, err := NewImageFromFile(outPath)
	require.NoError(t, err)
	defer result.Close()

	assert.Equal(t, ImageTypeJPEG, result.Format())
	assert.Equal(t, (origWidth+1)/2, result.Width())
}

func TestStreamingPipeline_ByteIdenticalToBufferPipeline(t *testing.T) {
	require.NoError(t, Startup(nil))

	path := resources + "heic-24bit.heic"

	// Buffer-based pipeline.
	buf, err := os.ReadFile(path)
	require.NoError(t, err)
	bufImg, err := LoadImageFromBuffer(buf, nil)
	require.NoError(t, err)
	defer bufImg.Close()
	require.NoError(t, bufImg.Resize(0.5, KernelLanczos3))
	expected, _, err := bufImg.ExportJpeg(nil)
	require.NoError(t, err)

	// Streaming pipeline with the same operations.
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()
	streamImg, err := LoadImageFromReader(f, nil)
	require.NoError(t, err)
	defer streamImg.Close()
	require.NoError(t, streamImg.Resize(0.5, KernelLanczos3))
	var w bytes.Buffer
	require.NoError(t, streamImg.SaveToWriter(&w, ImageTypeJPEG, nil))

	assert.True(t, bytes.Equal(expected, w.Bytes()),
		"streaming pipeline output must be byte-identical to the buffer pipeline")
}

func TestStreaming_GoldenOutput(t *testing.T) {
	goldenTest(t, resources+"jpg-24bit.jpg",
		func(img *ImageRef) error {
			return img.Resize(0.25, KernelLanczos3)
		},
		nil,
		func(img *ImageRef) ([]byte, *ImageMetadata, error) {
			var w bytes.Buffer
			if err := img.SaveToWriter(&w, ImageTypeJPEG, nil); err != nil {
				return nil, nil, err
			}
			return w.Bytes(), img.newMetadata(ImageTypeJPEG), nil
		},
	)
}

func TestStreaming_ConcurrentPipelines(t *testing.T) {
	require.NoError(t, Startup(nil))

	srcBuf, err := os.ReadFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	const goroutines = 8
	const iterations = 5

	var wg sync.WaitGroup
	errCh := make(chan error, goroutines*iterations)

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				err := func() error {
					img, err := LoadImageFromReader(&nonSeekable{r: bytes.NewReader(srcBuf)}, nil)
					if err != nil {
						return fmt.Errorf("load: %w", err)
					}
					defer img.Close()

					if err := img.Resize(0.5, KernelLanczos3); err != nil {
						return fmt.Errorf("resize: %w", err)
					}

					var w bytes.Buffer
					if err := img.SaveToWriter(&w, ImageTypePNG, nil); err != nil {
						return fmt.Errorf("save: %w", err)
					}
					if w.Len() == 0 {
						return errors.New("empty output")
					}
					return nil
				}()
				if err != nil {
					errCh <- err
				}
			}
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}

	sources, targets := streamRegistrySizes()
	assert.Zero(t, sources)
	assert.Zero(t, targets)
}

func TestStreaming_NoLeaks(t *testing.T) {
	require.NoError(t, Startup(nil))

	before := OpenImageRefs()

	// Success path.
	f, err := os.Open(resources + "jpg-24bit.jpg")
	require.NoError(t, err)
	img, err := LoadImageFromReader(f, nil)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	var w bytes.Buffer
	require.NoError(t, img.SaveToWriter(&w, ImageTypePNG, nil))
	img.Close()

	// Load error path.
	_, err = LoadImageFromReader(bytes.NewReader([]byte("not an image")), nil)
	require.Error(t, err)

	// Save error path.
	img2, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	require.Error(t, img2.SaveToWriter(&errWriter{err: errors.New("boom")}, ImageTypeJPEG, nil))
	img2.Close()

	assertNoNewImageRefs(t, before)
	sources, targets := streamRegistrySizes()
	assert.Zero(t, sources, "callback registry must hold no sources")
	assert.Zero(t, targets, "callback registry must hold no targets")
}

// --- Memory benchmark (T025, SC-002) ---

// TestStreamingLoadAllocsBounded proves the streaming path does not hold
// the compressed input in Go memory: its Go-side allocations must total
// less than 10% of the input file size, while the buffer path allocates
// at least 100% by definition.
func TestStreamingLoadAllocsBounded(t *testing.T) {
	require.NoError(t, Startup(nil))

	path := resources + "tif.tif" // largest test resource (~2.9 MB)
	info, err := os.Stat(path)
	require.NoError(t, err)
	fileSize := info.Size()

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)

	img, err := LoadImageFromReader(f, nil)
	require.NoError(t, err)

	runtime.ReadMemStats(&after)
	img.Close()

	allocated := int64(after.TotalAlloc - before.TotalAlloc)
	limit := fileSize / 10
	assert.Less(t, allocated, limit,
		"streaming load allocated %d Go bytes for a %d byte input (limit %d)",
		allocated, fileSize, limit)
	t.Logf("streaming load: %d Go bytes allocated for %d byte input (%.1f%%)",
		allocated, fileSize, 100*float64(allocated)/float64(fileSize))
}

// Note when comparing the two load benchmarks: streaming decodes the
// full image during load (the source is released immediately after), so
// it pays the decode cost upfront. Buffer load is lazy and defers decode
// to first pixel access. Compare allocated bytes, not wall time; for a
// full load+process+export pipeline the total work is equivalent (see
// TestStreamingPipeline_ByteIdenticalToBufferPipeline).
func BenchmarkStreamLoad(b *testing.B) {
	if err := Startup(nil); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		f, err := os.Open(resources + "tif.tif")
		if err != nil {
			b.Fatal(err)
		}
		img, err := LoadImageFromReader(f, nil)
		if err != nil {
			b.Fatal(err)
		}
		img.Close()
		f.Close()
	}
}

func BenchmarkBufferLoad(b *testing.B) {
	if err := Startup(nil); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf, err := os.ReadFile(resources + "tif.tif")
		if err != nil {
			b.Fatal(err)
		}
		img, err := LoadImageFromBuffer(buf, nil)
		if err != nil {
			b.Fatal(err)
		}
		img.Close()
	}
}

// --- SetPipeReadLimit (T027) ---

func TestSetPipeReadLimit(t *testing.T) {
	require.NoError(t, Startup(nil))

	// Restore the libvips default (~1 GB) after the test so other tests
	// are unaffected.
	defer SetPipeReadLimit(1024 * 1024 * 1024)

	SetPipeReadLimit(64 * 1024 * 1024)

	buf, err := os.ReadFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)

	img, err := LoadImageFromReader(&nonSeekable{r: bytes.NewReader(buf)}, nil)
	require.NoError(t, err, "small image must load under a lowered pipe limit")
	img.Close()
}
