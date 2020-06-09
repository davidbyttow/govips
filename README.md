# govips ![License](https://img.shields.io/badge/license-MIT-blue.svg)


## A libvips library for Go
This package wraps the core functionality of [libvips](https://github.com/libvips/libvips) image processing library by exposing all image operations on first-class types in Go.

How fast is libvips? See this: [Speed and Memory Use](https://github.com/libvips/libvips/wiki/Speed-and-memory-use)

The intent for this is to enable developers to build extremely fast image processors in Go, which is suited well for concurrent requests.

Libvips is generally 4-8x faster than other graphics processors such as GraphicsMagick and ImageMagick.

## Requirements
- [libvips](https://github.com/libvips/libvips) 8.10+
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.14+

## Installation
```bash
go get -u github.com/wix-playground/govips/vips
```

## Example usage
```go
package main

import (
	"github.com/wix-playground/govips/vips"
)

image, err := NewImageFromFile("image.jpg")
if err != nil {
	return nil, err
}
defer image.Close()

// Resize an image
return image.Resize(1200, 1200).Export(nil)
```

## Contributing
In short, feel free to file issues or send along pull requests. See this [guide on contributing](https://github.com/wix-playground/govips/blob/master/CONTRIBUTING.md) for more information.

## Credits
Thank you to [John Cupitt](https://github.com/jcupitt) for maintaining libvips and providing feedback on vips.

## License
MIT - David Byttow
