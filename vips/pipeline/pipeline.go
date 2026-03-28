// Package pipeline provides a composable, functional API for image transforms.
//
// Operations are values that describe a transformation without performing it.
// Compose them into reusable Transform chains, then apply with Apply (mutate)
// or ApplyNew (copy first).
//
//	t := pipeline.Compose(
//	    pipeline.Resize(0.5, vips.KernelLanczos3),
//	    pipeline.AutoRotate(),
//	    pipeline.Sharpen(1.0, 2.0, 3.0),
//	)
//	err := t.Apply(img)
//
// Export is intentionally outside the pipeline. Apply your transform, then
// call the appropriate Export method on the resulting image.
package pipeline

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

// Composable is implemented by types that can be used in Compose.
// Both Operation and Transform implement this interface.
type Composable interface {
	operations() []Operation
}

// Operation is a single image transformation step.
// Create operations using the provided wrapper functions (e.g. Resize, AutoRotate).
// For custom operations, use NewOperation.
type Operation struct {
	name string
	fn   func(*vips.ImageRef) error
}

func (o Operation) operations() []Operation {
	return []Operation{o}
}

// NewOperation creates a custom Operation with the given name and function.
// Use this to wrap custom transformations not covered by the built-in operations.
func NewOperation(name string, fn func(*vips.ImageRef) error) Operation {
	return Operation{name: name, fn: fn}
}

// Transform is a composable sequence of operations.
// Transforms are created with Compose and applied with Apply or ApplyNew.
type Transform struct {
	ops []Operation
}

func (t Transform) operations() []Operation {
	return t.ops
}

// Compose flattens any combination of Operations and Transforms into a single Transform.
// Transforms nest: Compose(transformA, transformB) produces a flat sequence of all operations.
func Compose(items ...Composable) Transform {
	var ops []Operation
	for _, item := range items {
		ops = append(ops, item.operations()...)
	}
	return Transform{ops: ops}
}

// Apply executes all operations on img in place.
// If any operation fails, it stops and returns an error with step context
// (e.g. "step 2/5 Sharpen: invalid sigma"). Errors are wrapped with %w
// so errors.Is and errors.As work on the underlying error.
func (t Transform) Apply(img *vips.ImageRef) error {
	total := len(t.ops)
	for i, op := range t.ops {
		if err := op.fn(img); err != nil {
			return fmt.Errorf("step %d/%d %s: %w", i+1, total, op.name, err)
		}
	}
	return nil
}

// ApplyNew copies the image first, then applies all operations to the copy.
// The original image is left untouched. The caller is responsible for closing
// the returned image when done.
func (t Transform) ApplyNew(img *vips.ImageRef) (*vips.ImageRef, error) {
	cp, err := img.Copy()
	if err != nil {
		return nil, fmt.Errorf("copy before transform: %w", err)
	}
	if err := t.Apply(cp); err != nil {
		cp.Close()
		return nil, err
	}
	return cp, nil
}

// Process loads an image from file, applies the given operations, and returns the result.
// The caller is responsible for closing the returned image when done.
// Export is intentionally not part of the pipeline — call Export on the returned image.
func Process(path string, items ...Composable) (*vips.ImageRef, error) {
	img, err := vips.NewImageFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", path, err)
	}
	t := Compose(items...)
	if err := t.Apply(img); err != nil {
		img.Close()
		return nil, err
	}
	return img, nil
}

// ProcessBytes loads an image from a byte buffer, applies the given operations,
// and returns the result. The caller is responsible for closing the returned image when done.
func ProcessBytes(buf []byte, items ...Composable) (*vips.ImageRef, error) {
	img, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		return nil, fmt.Errorf("load from buffer: %w", err)
	}
	t := Compose(items...)
	if err := t.Apply(img); err != nil {
		img.Close()
		return nil, err
	}
	return img, nil
}
