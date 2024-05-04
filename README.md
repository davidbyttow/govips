# <img src="https://raw.githubusercontent.com/davidbyttow/govips/master/assets/SVG/govips.svg" width="90" height="90"> <span style="font-size: 4em;">govips</span>

[![GoDoc](https://godoc.org/github.com/davidbyttow/govips?status.svg)](https://pkg.go.dev/mod/github.com/davidbyttow/govips/v2) [![Go Report Card](https://goreportcard.com/badge/github.com/davidbyttow/govips)](https://goreportcard.com/badge/github.com/davidbyttow/govips) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/davidbyttow/govips) ![License](https://img.shields.io/badge/license-MIT-blue.svg) [![Build Status](https://github.com/davidbyttow/govips/workflows/build/badge.svg)](https://github.com/davidbyttow/govips/actions) [![Coverage Status](https://img.shields.io/coveralls/github/davidbyttow/govips)](https://coveralls.io/github/davidbyttow/govips?branch=master)

## A lightning fast image processing and resizing library for Go

This package wraps the core functionality of [libvips](https://github.com/libvips/libvips) image processing library by exposing all image operations on first-class types in Go.

Libvips is generally 4-8x faster than other graphics processors such as GraphicsMagick and ImageMagick. Check the benchmark: [Speed and Memory Use](https://github.com/libvips/libvips/wiki/Speed-and-memory-use)

The intent for this is to enable developers to build extremely fast image processors in Go, which is suited well for concurrent requests.

## Requirements

-   [libvips](https://github.com/libvips/libvips) 8.10+
-   C compatible compiler such as gcc 4.6+ or clang 3.0+
-   Go 1.16+

## Dependencies

### MacOS

Use [homebrew](https://brew.sh/) to install vips and pkg-config:

```bash
brew install vips pkg-config
```

### Windows

The recommended approach on Windows is to use Govips via WSL and Ubuntu.

If you need to run Govips natively on Windows, it's not difficult but will require some effort. We don't have a recommended environment or setup at the moment. Windows is also not in our list of CI/CD targets so Govips is not regularly tested for compatibility. If you would be willing to setup and maintain a robust CI/CD Windows environment, please open a PR, we would be pleased to accept your contribution and support Windows as a platform.

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
	err = os.WriteFile("output.jpg", image1bytes, 0644)
	checkError(err)

}
```

See _examples/_ folder for more examples.

## Running tests

```bash
$ make test
```

## Memory usage note
### MALLOC_ARENA_MAX
`libvips` uses GLib for memory management, and it brings GLib memory fragmentation
issues to heavily multi-threaded programs. First thing you can try if you noticed
constantly growing RSS usage without Go's sys memory growth is set `MALLOC_ARENA_MAX`:

```
MALLOC_ARENA_MAX=2 application
```

This will reduce GLib memory appetites by reducing the number of malloc arenas
that it can create. By default GLib creates one are per thread, and this would
follow to memory fragmentation.

## Contributing

Feel free to file issues or create pull requests. See this [guide on contributing](https://github.com/davidbyttow/govips/blob/master/CONTRIBUTING.md) for more information.

## Credits

Thanks to:

-   [John Cupitt](https://github.com/jcupitt) for creating and maintaining libvips
-   [Toni Melisma](https://github.com/tonimelisma) for pushing to a 2.x release
-   [wix.com](https://wix.com/) for the govips logo and lots of great functionality
-   All of our fantastic [contributors](https://github.com/davidbyttow/govips/graphs/contributors)

## License

MIT
