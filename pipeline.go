package vips

// type Pipeline struct {
// 	image *ImageRef
// 	err   error
// }
//
// func New(image *ImageRef) *Pipeline {
// 	return &Pipeline{image, nil}
// }
//
// func (p *Pipeline) ReduceWidth(width int) *Pipeline {
// 	scale := scale(width, p.image.Width())
// 	if err := p.image.Reduceh(scale); err != nil {
// 		p.err = err
// 	}
// 	return p
// }
//
// func (p *Pipeline) ReduceHeight(height int) *Pipeline {
// 	scale := scale(height, p.image.Height())
// 	if err := p.image.Reducev(scale); err != nil {
// 		p.err = err
// 	}
// 	return p
// }
//
// func (p *Pipeline) Resize(width, height int) *Pipeline {
// 	scaleX := scale(width, p.image.Width())
// 	scaleY := scale(height, p.image.Height())
// 	if err := p.image.Reduce(scaleX, scaleY); err != nil {
// 		p.err = err
// 	}
// 	return p
// }
//
// func scale(x, y int) float64 {
// 	return float64(x) / float64(y)
// }
