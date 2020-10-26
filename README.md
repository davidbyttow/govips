# govips [![GoDoc](https://godoc.org/github.com/davidbyttow/govips?status.svg)](https://godoc.org/github.com/davidbyttow/govips) [![Go Report Card](http://goreportcard.com/badge/davidbyttow/govips)](http://goreportcard.com/report/davidbyttow/govips) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/govips/govips) ![License](https://img.shields.io/badge/license-MIT-blue.svg) [![Build Status](https://travis-ci.com/govips/govips.svg?branch=master)](https://travis-ci.com/govips/govips) [![Coverage Status](https://coveralls.io/repos/github/govips/govips/badge.svg?branch=master)](https://coveralls.io/github/govips/govips?branch=master)

## A fast image processing library for Go
This package wraps the core functionality of [libvips](https://github.com/libvips/libvips) image processing library by exposing all image operations on first-class types in Go.

Libvips is generally 4-8x faster than other graphics processors such as GraphicsMagick and ImageMagick. Check the benchmark: [Speed and Memory Use](https://github.com/libvips/libvips/wiki/Speed-and-memory-use)

The intent for this is to enable developers to build extremely fast image processors in Go, which is suited well for concurrent requests.

## Requirements
- [libvips](https://github.com/libvips/libvips) 8.10+
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.14+

## Installation
```bash
go get -u github.com/davidbyttow/govips/v2/vips
```

### Dependencies on macOS

Use [homebrew](https://brew.sh/) to install vips and pkg-config

```bash
brew install vips pkg-config
```

### Dependencies on Ubuntu

You need at least libvips 8.10.2 to work with govips. Groovy (20.10) repositories have the latest version. However on Bionic (18.04) and Focal (20.04), you need to install libvips and dependencies from a backports repository:

```bash
sudo add-apt-repository ppa:tonimelisma/ppa
```

Then:

```bash
sudo apt -y install libvips-dev
```

## Example usage
```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func main() {
	vips.Startup(nil)

	image1, err := vips.NewImageFromFile("input.jpg")
	checkError(err)
	defer image1.Close()

	// Rotate the picture upright and reset EXIF orientation tag
	err = image1.AutoRotate()
	checkError(err)

	image1bytes, _, err := image1.Export(&vips.ExportParams{Format: vips.ImageTypeJPEG})
	err = ioutil.WriteFile("output.jpg", image1bytes, 0644)
	checkError(err)

	vips.Shutdown()
}
```

## Contributing
In short, feel free to file issues or send along pull requests. See this [guide on contributing](https://github.com/davidbyttow/govips/blob/master/CONTRIBUTING.md) for more information.

## Credits
Thank you to [John Cupitt](https://github.com/jcupitt) for maintaining libvips and providing feedback on vips.

## License
MIT - David Byttow
