# govips [![GoDoc](https://godoc.org/github.com/davidbyttow/govips?status.svg)](https://pkg.go.dev/mod/github.com/davidbyttow/govips/v2) [![Go Report Card](http://goreportcard.com/badge/davidbyttow/govips)](http://goreportcard.com/report/davidbyttow/govips) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/davidbyttow/govips) ![License](https://img.shields.io/badge/license-MIT-blue.svg) [![Build Status](https://github.com/davidbyttow/govips/workflows/build/badge.svg)](https://github.com/davidbyttow/govips/actions) [![Coverage Status](https://img.shields.io/coveralls/github/davidbyttow/govips)](https://coveralls.io/github/davidbyttow/govips?branch=master)

## A lightning fast image processing and resizing library for Go

This package wraps the core functionality of [libvips](https://github.com/libvips/libvips) image processing library by exposing all image operations on first-class types in Go.

Libvips is generally 4-8x faster than other graphics processors such as GraphicsMagick and ImageMagick. Check the benchmark: [Speed and Memory Use](https://github.com/libvips/libvips/wiki/Speed-and-memory-use)

The intent for this is to enable developers to build extremely fast image processors in Go, which is suited well for concurrent requests.

## Requirements

-   [libvips](https://github.com/libvips/libvips) 8.10+
-   C compatible compiler such as gcc 4.6+ or clang 3.0+
-   Go 1.14+

## Dependencies

### MacOS

Use [homebrew](https://brew.sh/) to install vips and pkg-config:

```bash
brew install vips pkg-config
```

### Ubuntu

You need at least libvips 8.10.2 to work with govips. Groovy (20.10) repositories have the latest version. However on Bionic (18.04) and Focal (20.04), you need to install libvips and dependencies from a backports repository:

```bash
sudo add-apt-repository ppa:tonimelisma/ppa
```

Then:

```bash
sudo apt -y install libvips-dev
```

## Installation

```bash
go get -u github.com/davidbyttow/govips/v2/vips
```

### MacOS note

On MacOS, govips may not compile without first setting an environment variable:

```bash
export CGO_CFLAGS_ALLOW="-Xpreprocessor"
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
	defer vips.Shutdown()

	image1, err := vips.NewImageFromFile("input.jpg")
	checkError(err)

	// Rotate the picture upright and reset EXIF orientation tag
	err = image1.AutoRotate()
	checkError(err)

	ep := vips.NewDefaultJPEGExportParams()
	image1bytes, _, err := image1.Export(ep)
	err = ioutil.WriteFile("output.jpg", image1bytes, 0644)
	checkError(err)

}
```

See *examples/* folder for more examples.

## Running tests

```bash
$ make test
```

## Contributing

Feel free to file issues or create pull requests. See this [guide on contributing](https://github.com/davidbyttow/govips/blob/master/CONTRIBUTING.md) for more information.

## Credits

Thanks to:

-   [John Cupitt](https://github.com/jcupitt) for creating and maintaining libvips
-   [Toni Melisma](https://github.com/tonimelisma) for pushing to a 2.x release
-   All of our fantastic [contributors](https://github.com/davidbyttow/govips/graphs/contributors)

## License

MIT - David Byttow
