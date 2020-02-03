# govips [![Build Status](https://travis-ci.org/jhford/govips.svg)](https://travis-ci.org/jhford/govips) [![GoDoc](https://godoc.org/github.com/jhford/govips?status.svg)](https://godoc.org/github.com/jhford/govips) [![Go Report Card](http://goreportcard.com/badge/jhford/govips)](http://goreportcard.com/report/jhford/govips) ![License](https://img.shields.io/badge/license-MIT-blue.svg)

# A libvips library for Go
This package wraps the core functionality of [libvips](https://github.com/libvips/libvips) image processing library by exposing all image operations on first-class types in Go. Additionally, it exposes raw access to call operations directly, for forward compatibility.

How fast is libvips? See this: [Speed and Memory Use](https://github.com/libvips/libvips/wiki/Speed-and-memory-use)

This library was inspired primarily based on the C++ wrapper in libvips.

The intent for this is to enable developers to build extremely fast image processors in Go, which is suited well for concurrent requests.

Libvips is generally 4-8x faster than other graphics processors such as GraphicsMagick and ImageMagick.

# Supported image operations
This library supports all known operations available to libvips found here:
- [VIPS function list](http://libvips.github.io/libvips/API/current/VipsImage.html)
- [VipsImage](http://libvips.github.io/libvips/API/current/VipsImage.html)
- [VipsOperation](http://libvips.github.io/libvips/API/current/VipsOperation.html)

# Requirements
- [libvips](https://github.com/libvips/libvips) 8+ (8.5.8+ recommended for GIF, PDF, SVG support)
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.9+

# Installation
```bash
go get -u github.com/jhford/govips/pkg/vips
```

# Example usage
```go
// Resize an image with padding
return vips.NewPipeline().
	LoadFile(inputFile).
	PadStrategy(vips.ExtendBlack).
	Resize(1200, 1200).
	OutputFile(outputFile)
```

# Contributing
In short, feel free to file issues or send along pull requests. See this [guide on contributing](https://github.com/jhford/govips/blob/master/CONTRIBUTING.md) for more information.

# Credits
Thank you to [John Cupitt](https://github.com/jcupitt) for maintaining libvips and providing feedback on vips.

# License
MIT - David Byttow
