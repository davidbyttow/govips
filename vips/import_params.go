package vips

type ImportParams struct {
	// Fail on errors
	// Setting fail to TRUE makes the reader fail on any errors. This can be
	// useful for detecting truncated files, for example. Normally reading these
	// produces a warning, but no fatal error.
	//
	// Ignored for unsupported loaders.
	Fail *bool

	// Use exif Orientation tag to rotate the image during load.
	//
	// Ignored for unsupported loaders.
	AutoRotate *bool

	// Shrink by this much on load.
	//
	// Fail on unsupported loaders.
	Shrink *int
}

func NewImportParams() *ImportParams {
	return &ImportParams{}
}

func (p *ImportParams) WithFail(value bool) *ImportParams {
	p.Fail = &value
	return p
}

func (p *ImportParams) WithAutoRotate(value bool) *ImportParams {
	p.AutoRotate = &value
	return p
}

func (p *ImportParams) WithShrink(value int) *ImportParams {
	p.Shrink = &value
	return p
}

// MergeImportParams combines the given ImportParams instances into a single ImportParams in a last-one-wins
// fashion.
func MergeImportParams(opts ...*ImportParams) *ImportParams {
	d := NewImportParams()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Fail != nil {
			d.Fail = opt.Fail
		}
		if opt.AutoRotate != nil {
			d.Fail = opt.AutoRotate
		}
		if opt.Shrink != nil {
			d.Shrink = opt.Shrink
		}
	}
	return d
}
