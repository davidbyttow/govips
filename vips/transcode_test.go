package vips

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- fixtures and instrumentation ---

// bigSide is chosen so the decoded frame (bigSide² bytes, single band)
// dwarfs every memory budget asserted below: 144 MB decoded from a
// compressed input of well under 1 MB — a mild pixel bomb.
const bigSide = 12000

var (
	bigJPEGOnce sync.Once
	bigJPEGBuf  []byte
	bigJPEGErr  error
)

// bigJPEG returns (building once) a tiny compressed JPEG that decodes to
// a bigSide x bigSide single-band frame. Generation itself streams
// through libvips, so building the fixture does not need the frame in
// RAM either.
func bigJPEG(t *testing.T) []byte {
	t.Helper()
	require.NoError(t, Startup(nil))

	bigJPEGOnce.Do(func() {
		img, err := Black(bigSide, bigSide)
		if err != nil {
			bigJPEGErr = err
			return
		}
		defer img.Close()
		// Baseline (non-interlaced) JPEG: progressive JPEG cannot be
		// decoded incrementally — the codec needs every scan before it
		// can emit complete rows, which would defeat the pacing
		// assertions below (and is called out in the streaming docs).
		bigJPEGBuf, _, bigJPEGErr = img.ExportJpeg(&JpegExportParams{Quality: 80, Interlace: false})
	})
	require.NoError(t, bigJPEGErr)
	return bigJPEGBuf
}

// countingReader is a non-seekable reader that counts delivered bytes.
// The counter is atomic because libvips reads from worker threads while
// the test goroutine asserts on it.
type countingReader struct {
	r io.Reader
	n atomic.Int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n.Add(int64(n))
	return n, err
}

// instrumentedWriter records, per chunk: the count, how much input had
// been consumed when the first chunk arrived, and the peak libvips
// tracked memory observed across chunk deliveries (a proxy for peak RSS
// of pixel buffers). Optionally it stalls once at chunk stallAt,
// sampling memory growth during the stall. Writes are serialized by the
// target entry's per-instance mutex, so plain fields are safe.
type instrumentedWriter struct {
	buf              bytes.Buffer
	chunks           int
	reader           *countingReader
	readAtFirstChunk int64
	peakVipsMem      int64

	stallAt     int
	stallFor    time.Duration
	stallGrowth int64
}

func (w *instrumentedWriter) sampleMem() int64 {
	var ms MemoryStats
	ReadVipsMemStats(&ms)
	if ms.Mem > w.peakVipsMem {
		w.peakVipsMem = ms.Mem
	}
	return ms.Mem
}

func (w *instrumentedWriter) Write(p []byte) (int, error) {
	w.chunks++
	if w.chunks == 1 && w.reader != nil {
		w.readAtFirstChunk = w.reader.n.Load()
	}
	base := w.sampleMem()

	if w.stallAt > 0 && w.chunks == w.stallAt {
		// Simulate a network/multipart flush stalling mid-encode and
		// watch whether libvips buffers grow while we block.
		deadline := time.Now().Add(w.stallFor)
		for time.Now().Before(deadline) {
			time.Sleep(50 * time.Millisecond)
			if m := w.sampleMem(); m-base > w.stallGrowth {
				w.stallGrowth = m - base
			}
		}
	}

	return w.buf.Write(p)
}

// baselineJPEG produces non-progressive JPEG output. govips' default
// export params use Interlace (progressive) encoding, which buffers the
// whole image inside libjpeg before writing the first byte — legal, but
// it defeats chunked-output and pacing assertions. True streaming needs
// baseline JPEG on both ends (see the README streaming caveats).
var baselineJPEG = &ExportParams{Quality: 80, Interlaced: false}

func resetStreamKnobs() {
	SetStreamDiscThreshold(-1)
	SetStreamScratchDir("")
}

// --- byte identity (requirement 4) ---

func TestTranscodeStream_SequentialByteIdentity(t *testing.T) {
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
			err = TranscodeStream(&nonSeekable{r: bytes.NewReader(srcBuf)}, &w,
				&TranscodeOptions{Format: format})
			require.NoError(t, err)

			assert.True(t, bytes.Equal(expected, w.Bytes()),
				"sequential transcode must be byte-identical to Export* (got %d bytes, want %d)",
				w.Len(), len(expected))
		})
	}

	sources, targets := streamRegistrySizes()
	assert.Zero(t, sources)
	assert.Zero(t, targets)
}

func TestTranscodeStream_FormatUnknownKeepsInputFormat(t *testing.T) {
	require.NoError(t, Startup(nil))

	srcBuf, err := os.ReadFile(resources + "png-24bit.png")
	require.NoError(t, err)
	img, err := LoadImageFromBuffer(srcBuf, nil)
	require.NoError(t, err)
	defer img.Close()
	expected, _, err := img.ExportPng(nil)
	require.NoError(t, err)

	var w bytes.Buffer
	require.NoError(t, TranscodeStream(bytes.NewReader(srcBuf), &w, nil))
	assert.True(t, bytes.Equal(expected, w.Bytes()))
}

func TestDiscBackedLoad_ByteIdentityAndScratchCleanup(t *testing.T) {
	require.NoError(t, Startup(nil))

	scratch := t.TempDir()
	SetStreamScratchDir(scratch)
	SetStreamDiscThreshold(1) // force every materialization to disc
	defer resetStreamKnobs()

	srcBuf, err := os.ReadFile(resources + "png-24bit.png")
	require.NoError(t, err)

	expectedImg, err := LoadImageFromBuffer(srcBuf, nil)
	require.NoError(t, err)
	defer expectedImg.Close()
	expected, _, err := expectedImg.ExportJpeg(nil)
	require.NoError(t, err)

	img, err := LoadImageFromReader(bytes.NewReader(srcBuf), nil)
	require.NoError(t, err)
	defer img.Close()

	// The scratch file must already be unlinked while in use.
	entries, err := os.ReadDir(scratch)
	require.NoError(t, err)
	assert.Empty(t, entries, "scratch files must be unlinked immediately after open")

	var w bytes.Buffer
	require.NoError(t, img.SaveToWriter(&w, ImageTypeJPEG, nil))
	assert.True(t, bytes.Equal(expected, w.Bytes()),
		"disc-backed image must encode byte-identically to the buffer path")
}

// --- memory budget / pixel bomb / incremental I/O (requirements 5a-5d) ---

// memBudget is the hard ceiling for libvips tracked memory while
// processing the bigSide² frame: a quarter of the decoded size. Staying
// under it proves the frame is not resident in RAM.
const memBudget = int64(bigSide) * bigSide / 4

func TestTranscodeStream_SequentialMemoryBounded_PixelBomb(t *testing.T) {
	input := bigJPEG(t)
	t.Logf("pixel bomb: %d compressed bytes -> %d decoded bytes", len(input), int64(bigSide)*bigSide)

	reader := &countingReader{r: bytes.NewReader(input)}
	w := &instrumentedWriter{reader: reader}

	require.NoError(t, TranscodeStream(reader, w,
		&TranscodeOptions{Format: ImageTypeJPEG, ExportParams: baselineJPEG}))

	// 5a: the decoded frame never lived in RAM.
	assert.Less(t, w.peakVipsMem, memBudget,
		"peak libvips memory %d must stay under %d (decoded frame is %d)",
		w.peakVipsMem, memBudget, int64(bigSide)*bigSide)

	// 5c: output arrived in chunks during encoding.
	assert.GreaterOrEqual(t, w.chunks, 4, "output must arrive in chunks during encoding")

	// 5b: input was pulled gradually — encoding had begun long before
	// the input was fully consumed.
	assert.Less(t, w.readAtFirstChunk, int64(len(input))/2,
		"first output chunk arrived after consuming %d of %d input bytes — input was slurped",
		w.readAtFirstChunk, len(input))

	// The result must still be a valid full-size image.
	out, err := LoadImageFromBuffer(w.buf.Bytes(), nil)
	require.NoError(t, err)
	defer out.Close()
	assert.Equal(t, bigSide, out.Width())
	assert.Equal(t, bigSide, out.Height())

	sources, targets := streamRegistrySizes()
	assert.Zero(t, sources)
	assert.Zero(t, targets)
}

func TestDiscBackedLoad_MemoryBounded(t *testing.T) {
	input := bigJPEG(t)

	SetStreamScratchDir(t.TempDir())
	SetStreamDiscThreshold(8 << 20) // 144 MB decoded frame goes to disc
	defer resetStreamKnobs()

	// Default (random-access) load: the frame materializes to scratch
	// disc in one streaming pass; RAM stays bounded throughout.
	img, err := LoadImageFromReader(bytes.NewReader(input), nil)
	require.NoError(t, err)
	defer img.Close()

	var ms MemoryStats
	ReadVipsMemStats(&ms)
	assert.Less(t, ms.Mem, memBudget,
		"libvips memory %d after disc-backed load must stay under %d", ms.Mem, memBudget)

	// Random access now works (the whole point of this mode)...
	require.NoError(t, img.Flip(DirectionVertical))

	// ...and encoding from disc stays within budget too.
	w := &instrumentedWriter{}
	require.NoError(t, img.SaveToWriter(w, ImageTypeJPEG, baselineJPEG))
	assert.Less(t, w.peakVipsMem, memBudget)
	assert.GreaterOrEqual(t, w.chunks, 4)
}

func TestLoadImageFromReader_SequentialIsLazy(t *testing.T) {
	input := bigJPEG(t)

	reader := &countingReader{r: bytes.NewReader(input)}
	params := &ImportParams{}
	params.Access.Set(AccessSequential)

	img, err := LoadImageFromReader(&nonSeekable{r: reader}, params)
	require.NoError(t, err)

	// Only the header was sniffed; the bulk of the input is unread.
	assert.Less(t, reader.n.Load(), int64(len(input))/5,
		"sequential load must read only the header, consumed %d of %d bytes",
		reader.n.Load(), len(input))
	assert.Equal(t, bigSide, img.Width())

	// The source stays registered until Close, then everything is
	// released even though no pixel was ever decoded.
	sources, _ := streamRegistrySizes()
	assert.Equal(t, 1, sources, "sequential image must keep its source registered")
	img.Close()
	sources, _ = streamRegistrySizes()
	assert.Zero(t, sources, "Close must release the streaming source")
}

// --- truncated input (requirement 5e) ---

func TestTranscodeStream_TruncatedSequential(t *testing.T) {
	input := bigJPEG(t)
	truncated := input[:len(input)/3]

	var w bytes.Buffer
	err := TranscodeStream(&nonSeekable{r: bytes.NewReader(truncated)}, &w,
		&TranscodeOptions{Format: ImageTypeJPEG, ExportParams: baselineJPEG})
	// On the sequential path the error surfaces during the encode pass
	// (the first full pass over the pixels), wrapped by TranscodeStream.
	require.Error(t, err, "truncated sequential stream must fail, not emit a corrupt image")

	sources, targets := streamRegistrySizes()
	assert.Zero(t, sources, "failed transcode must release the source")
	assert.Zero(t, targets, "failed transcode must release the target")
}

func TestTranscodeStream_SequentialReaderErrorTIFF(t *testing.T) {
	input := bigJPEG(t)

	readerErr := errors.New("connection reset by peer")
	r := &errAfterReader{data: input[:64*1024], err: readerErr}

	// TIFF takes the buffer-encode fallback inside SaveToWriter; the
	// sequential reader's error must survive that path too.
	var w bytes.Buffer
	err := TranscodeStream(&nonSeekable{r: r}, &w, &TranscodeOptions{Format: ImageTypeTIFF})
	require.Error(t, err)
	assert.True(t, errors.Is(err, readerErr),
		"reader error must be wrapped through the TIFF fallback, got: %v", err)
}

func TestDiscBackedLoad_Truncated(t *testing.T) {
	input := bigJPEG(t)
	truncated := input[:len(input)/3]

	SetStreamDiscThreshold(8 << 20)
	defer resetStreamKnobs()

	// On the materialized path the error surfaces at load time, during
	// the render-to-scratch pass.
	img, err := LoadImageFromReader(bytes.NewReader(truncated), nil)
	if img != nil {
		img.Close()
	}
	require.Error(t, err, "truncated input must fail during disc materialization")

	sources, _ := streamRegistrySizes()
	assert.Zero(t, sources)
}

// --- slow/blocking writer (requirement 5f) ---

func TestTranscodeStream_SlowWriterBackpressure(t *testing.T) {
	input := bigJPEG(t)

	// Reference run with a fast writer.
	var fast bytes.Buffer
	require.NoError(t, TranscodeStream(bytes.NewReader(input), &fast,
		&TranscodeOptions{Format: ImageTypeJPEG, ExportParams: baselineJPEG}))

	// Stalled run: block for 2s mid-encode, sampling memory growth.
	slow := &instrumentedWriter{stallAt: 3, stallFor: 2 * time.Second}
	start := time.Now()
	require.NoError(t, TranscodeStream(bytes.NewReader(input), slow,
		&TranscodeOptions{Format: ImageTypeJPEG, ExportParams: baselineJPEG}))
	require.True(t, time.Since(start) >= 2*time.Second, "stall must actually have happened")

	// No unbounded buffering while the sink was blocked: encoding
	// backpressures instead of piling up pixels.
	assert.Less(t, slow.stallGrowth, int64(32<<20),
		"libvips memory grew by %d during a 2s sink stall", slow.stallGrowth)

	// No corruption: byte-identical to the fast run.
	assert.True(t, bytes.Equal(fast.Bytes(), slow.buf.Bytes()),
		"stalled run output differs from fast run")
}

// --- automatic path selection (requirement 1+2) ---

func TestTranscodeStream_AutoRotateMaterializes(t *testing.T) {
	require.NoError(t, Startup(nil))

	srcBuf, err := os.ReadFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)

	// Buffer-path reference: load, autorotate, export.
	ref, err := LoadImageFromBuffer(srcBuf, nil)
	require.NoError(t, err)
	defer ref.Close()
	origW, origH := ref.Width(), ref.Height()
	require.NoError(t, ref.AutoRotate())
	expected, _, err := ref.ExportJpeg(nil)
	require.NoError(t, err)

	// Orientation 6 needs a 90° rotation → forces materialization.
	var w bytes.Buffer
	require.NoError(t, TranscodeStream(bytes.NewReader(srcBuf), &w,
		&TranscodeOptions{Format: ImageTypeJPEG, AutoRotate: true}))

	assert.True(t, bytes.Equal(expected, w.Bytes()),
		"autorotated transcode must be byte-identical to the buffer path")

	out, err := LoadImageFromBuffer(w.Bytes(), nil)
	require.NoError(t, err)
	defer out.Close()
	assert.Equal(t, origH, out.Width(), "90° rotation must swap dimensions")
	assert.Equal(t, origW, out.Height())

	sources, targets := streamRegistrySizes()
	assert.Zero(t, sources, "materialization must release the source")
	assert.Zero(t, targets)
}
