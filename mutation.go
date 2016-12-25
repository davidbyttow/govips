package gimage

import "github.com/simplethingsllc/gimage/vips"

type Mutation struct {
	image   vips.Image
	options *Options
}

func NewMutation(buf []byte) (Mutation, error) {
	image, err := vips.NewImage(buf)
	if err != nil {
		return Mutation{}, err
	}
	return Mutation{image, &Options{}}, nil
}

type ResizeOptions struct {
	Kernel         Kernel
	Interpolator   Interpolator
	CenterSampling bool
}

func (o Mutation) ResizeH(width int, options *ResizeOptions) Mutation {
	o.Resize(width, 0, options)
	return o
}

func (o Mutation) ResizeV(height int, options *ResizeOptions) Mutation {
	o.Resize(0, height, options)
	return o
}

func (o Mutation) Resize(width, height int, options *ResizeOptions) Mutation {
	o.options.Width = width
	o.options.Height = height
	if options != nil {
		o.options.Kernel = options.Kernel
		o.options.Interpolator = options.Interpolator
		o.options.CenterSampling = options.CenterSampling
	}
	return o
}

func (o Mutation) Apply() ([]byte, error) {
	var err error
	o.image, err = process(o.image, o.options)
	if err != nil {
		return nil, err
	}
	return nil, nil
	//return o.image.Bytes(), nil
}
