package vips_test

import (
	"testing"

	"github.com/wix-playground/govips/pkg/vips"
)

func TestOptionPrimitives(t *testing.T) {
	var b bool
	var i int
	var d float64
	var s string
	options := []*vips.Option{
		vips.InputBool("b", true),
		vips.InputInt("i", 42),
		vips.InputDouble("d", 42.2),
		vips.InputString("s", "hi"),
		vips.OutputBool("b", &b),
		vips.OutputInt("i", &i),
		vips.OutputDouble("d", &d),
		vips.OutputString("s", &s),
	}

	// TODO(d): Write tests
	for i := 0; i < 8; i++ {
		opt := options[i]
		if !opt.Output() {

		}
	}
}
