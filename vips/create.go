package vips

// #include "create.h"
import "C"
import "unsafe"

type TextWrap int

type TextParams struct {
	Text      string
	Font      string
	Width     int
	Height    int
	Alignment Align
	DPI       int
	RGBA      bool
	Justify   bool
	Spacing   int
	Wrap      TextWrap
}

type vipsTextOptions struct {
	Text      *C.char
	Font      *C.char
	Width     C.int
	Height    C.int
	DPI       C.int
	RGBA      C.gboolean
	Justify   C.gboolean
	Spacing   C.int
	Alignment C.VipsAlign
	Wrap      C.VipsTextWrap
}

// TextWrap enum
const (
	TextWrapWord     TextWrap = C.VIPS_TEXT_WRAP_WORD
	TextWrapChar     TextWrap = C.VIPS_TEXT_WRAP_CHAR
	TextWrapWordChar TextWrap = C.VIPS_TEXT_WRAP_WORD_CHAR
	TextWrapNone     TextWrap = C.VIPS_TEXT_WRAP_NONE
)

// https://libvips.github.io/libvips/API/current/libvips-create.html#vips-text
func vipsText(params *TextParams) (*C.VipsImage, error) {
	var out *C.VipsImage

	text := C.CString(params.Text)
	defer freeCString(text)

	font := C.CString(params.Font)
	defer freeCString(font)

	opts := vipsTextOptions{
		Text:      text,
		Font:      font,
		Width:     C.int(params.Width),
		Height:    C.int(params.Height),
		DPI:       C.int(params.DPI),
		Alignment: C.VipsAlign(params.Alignment),
		Spacing:   C.int(params.Spacing),
		Wrap:      C.VipsTextWrap(params.Wrap),
	}

	if params.RGBA {
		opts.RGBA = C.TRUE
	}

	if params.Justify {
		opts.Justify = C.TRUE
	}

	err := C.text(&out, (*C.TextOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}
