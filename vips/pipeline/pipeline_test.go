package pipeline_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/davidbyttow/govips/v2/vips/pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const resources = "../../resources/"

func TestCompose_SingleOperation(t *testing.T) {
	tr := pipeline.Compose(pipeline.AutoRotate())
	// Verify it's composable by applying to a real image
	require.NoError(t, vips.Startup(nil))
	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()
	require.NoError(t, tr.Apply(img))
}

func TestCompose_FlattenTransforms(t *testing.T) {
	step1 := pipeline.Compose(pipeline.Resize(0.5, vips.KernelLanczos3), pipeline.AutoRotate())
	step2 := pipeline.Compose(pipeline.Invert())

	// Compose transforms together
	combined := pipeline.Compose(step1, step2)

	require.NoError(t, vips.Startup(nil))
	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	require.NoError(t, combined.Apply(img))
}

func TestCompose_MixedOperationsAndTransforms(t *testing.T) {
	step1 := pipeline.Compose(pipeline.Resize(0.5, vips.KernelLanczos3))
	combined := pipeline.Compose(step1, pipeline.AutoRotate(), pipeline.Invert())

	require.NoError(t, vips.Startup(nil))
	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	require.NoError(t, combined.Apply(img))
}

func TestCompose_Empty(t *testing.T) {
	require.NoError(t, vips.Startup(nil))
	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	tr := pipeline.Compose()
	require.NoError(t, tr.Apply(img))
}

func TestApply_MutatesInPlace(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	origWidth := img.Width()

	tr := pipeline.Compose(pipeline.Resize(0.5, vips.KernelLanczos3))
	require.NoError(t, tr.Apply(img))

	assert.InDelta(t, origWidth/2, img.Width(), 1)
}

func TestApplyNew_PreservesOriginal(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	origWidth := img.Width()

	tr := pipeline.Compose(pipeline.Resize(0.5, vips.KernelLanczos3))
	result, err := tr.ApplyNew(img)
	require.NoError(t, err)
	defer result.Close()

	assert.Equal(t, origWidth, img.Width())
	assert.InDelta(t, origWidth/2, result.Width(), 1)
}

func TestApply_ErrorContext(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	failing := pipeline.NewOperation("BadOp", func(img *vips.ImageRef) error {
		return errors.New("something broke")
	})

	tr := pipeline.Compose(pipeline.AutoRotate(), failing, pipeline.Invert())
	err = tr.Apply(img)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "step 2/3 BadOp")
	assert.Contains(t, err.Error(), "something broke")
}

func TestApply_ErrorUnwrap(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	sentinel := errors.New("sentinel error")
	failing := pipeline.NewOperation("Fail", func(img *vips.ImageRef) error {
		return sentinel
	})

	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	err = pipeline.Compose(failing).Apply(img)
	require.Error(t, err)
	assert.True(t, errors.Is(err, sentinel))
}

func TestApply_StopsOnFirstError(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	var called []string

	op1 := pipeline.NewOperation("Op1", func(img *vips.ImageRef) error {
		called = append(called, "op1")
		return nil
	})
	op2 := pipeline.NewOperation("Op2", func(img *vips.ImageRef) error {
		called = append(called, "op2")
		return errors.New("fail")
	})
	op3 := pipeline.NewOperation("Op3", func(img *vips.ImageRef) error {
		called = append(called, "op3")
		return nil
	})

	err = pipeline.Compose(op1, op2, op3).Apply(img)
	require.Error(t, err)
	assert.Equal(t, []string{"op1", "op2"}, called)
}

func TestApplyNew_ErrorClosesImage(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	failing := pipeline.NewOperation("Fail", func(img *vips.ImageRef) error {
		return errors.New("boom")
	})

	result, err := pipeline.Compose(failing).ApplyNew(img)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestProcess(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	img, err := pipeline.Process(resources+"png-24bit.png",
		pipeline.Resize(0.5, vips.KernelLanczos3),
		pipeline.AutoRotate(),
	)
	require.NoError(t, err)
	defer img.Close()

	assert.Greater(t, img.Width(), 0)
}

func TestProcess_BadPath(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	_, err := pipeline.Process("/nonexistent/file.png", pipeline.AutoRotate())
	require.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "load "))
}

func TestProcess_WithComposedTransforms(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	resize := pipeline.Compose(pipeline.Resize(0.5, vips.KernelLanczos3))
	cleanup := pipeline.Compose(pipeline.AutoRotate())

	img, err := pipeline.Process(resources+"png-24bit.png", resize, cleanup)
	require.NoError(t, err)
	defer img.Close()

	assert.Greater(t, img.Width(), 0)
}

func TestProcessBytes(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	src, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	buf, _, err := src.ExportPng(vips.NewPngExportParams())
	src.Close()
	require.NoError(t, err)

	img, err := pipeline.ProcessBytes(buf, pipeline.Resize(0.5, vips.KernelLanczos3))
	require.NoError(t, err)
	defer img.Close()

	assert.Greater(t, img.Width(), 0)
}

func TestMultiStepPipeline(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	// Define reusable presets
	webOptimize := pipeline.Compose(
		pipeline.Resize(0.5, vips.KernelLanczos3),
		pipeline.AutoRotate(),
		pipeline.Sharpen(1.0, 2.0, 10.0),
	)

	colorCorrect := pipeline.Compose(
		pipeline.Gamma(1.2),
	)

	// Combine presets into a full pipeline
	fullPipeline := pipeline.Compose(webOptimize, colorCorrect)

	img, err := vips.NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)
	defer img.Close()

	origWidth := img.Width()
	require.NoError(t, fullPipeline.Apply(img))

	// Verify resize happened
	assert.InDelta(t, origWidth/2, img.Width(), 1)

	// Verify the image is still exportable (all operations succeeded)
	buf, meta, err := img.ExportJpeg(vips.NewJpegExportParams())
	require.NoError(t, err)
	assert.NotEmpty(t, buf)
	assert.Equal(t, vips.ImageTypeJPEG, meta.Format)
}

func TestMultiStepPipeline_WithApplyNew(t *testing.T) {
	require.NoError(t, vips.Startup(nil))

	img, err := vips.NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	origWidth := img.Width()
	origHeight := img.Height()

	// Create two different outputs from the same source
	thumbTransform := pipeline.Compose(
		pipeline.Resize(0.25, vips.KernelLanczos3),
		pipeline.Sharpen(0.5, 1.0, 2.0),
	)

	previewTransform := pipeline.Compose(
		pipeline.Resize(0.5, vips.KernelLanczos3),
	)

	thumb, err := thumbTransform.ApplyNew(img)
	require.NoError(t, err)
	defer thumb.Close()

	preview, err := previewTransform.ApplyNew(img)
	require.NoError(t, err)
	defer preview.Close()

	// Original unchanged
	assert.Equal(t, origWidth, img.Width())
	assert.Equal(t, origHeight, img.Height())

	// Thumb is 25%
	assert.InDelta(t, origWidth/4, thumb.Width(), 1)

	// Preview is 50%
	assert.InDelta(t, origWidth/2, preview.Width(), 1)
}
