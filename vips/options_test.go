package vips

import (
	"testing"
)

func TestOptionPrimitives(t *testing.T) {
	var b bool
	var i int
	var d float64
	var s string
	options := []*Option{
		InputBool("b", true),
		InputInt("i", 42),
		InputDouble("d", 42.2),
		InputString("s", "hi"),
		OutputBool("b", &b),
		OutputInt("i", &i),
		OutputDouble("d", &d),
		OutputString("s", &s),
	}

	// TODO: Write tests
	for i := 0; i < 8; i++ {
		opt := options[i]
		if !opt.Output() {

		}
	}
}
