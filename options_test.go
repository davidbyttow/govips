package gimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	options := NewOptions()
	assert.NotNil(t, options)
}
