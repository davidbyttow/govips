# <img src="https://raw.githubusercontent.com/davidbyttow/govips/master/assets/SVG/govips.svg" width="90" height="90"> <span style="font-size: 4em;">govips</span>

[![GoDoc](https://godoc.org/github.com/davidbyttow/govips?status.svg)](https://pkg.go.dev/mod/github.com/davidbyttow/govips/v2) [![Go Report Card](https://goreportcard.com/badge/github.com/davidbyttow/govips)](https://goreportcard.com/badge/github.com/davidbyttow/govips) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/davidbyttow/govips) ![License](https://img.shields.io/badge/license-MIT-blue.svg) [![Build Status](https://github.com/davidbyttow/govips/workflows/build/badge.svg)](https://github.com/davidbyttow/govips/actions) [![Coverage Status](https://img.shields.io/coveralls/github/davidbyttow/govips)](https://coveralls.io/github/davidbyttow/govips?branch=master)

## A lightning fast image processing and resizing library for Go

This package wraps the core functionality of [libvips](https://github.com/libvips/libvips) image processing library by exposing all image operations on first-class types in Go.

Libvips is generally 4-8x faster than other graphics processors such as GraphicsMagick and ImageMagick. Check the benchmark: [Speed and Memory Use](https://github.com/libvips/libvips/wiki/Speed-and-memory-use)

The intent for this is to enable developers to build extremely fast image processors in Go, which is suited well for concurrent requests.

## Project Status

govips now includes a built-in code generator (`cmd/vipsgen/`) that uses libvips GObject introspection to auto-generate type-safe Go bindings. This covers 193+ operations across 9 categories (arithmetic, colour, conversion, convolution, create, freqfilt, histogram, morphology, resample), while complex operations like image I/O remain as hand-written bindings for full control.

## Requirements

-   [libvips](https://github.com/libvips/libvips) 8.14+
-   C compatible compiler such as gcc 4.6+ or clang 3.0+
-   Go 1.23+

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

## Examples

Every example below assumes this setup:

```go
package main

import (
	"fmt"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()

	// ... example code goes here
}
```

### 1. Load an image and export as JPEG

The basics: load from a file, auto-rotate based on EXIF data, and write it back out.

```go
image, err := vips.NewImageFromFile("input.jpg")
if err != nil {
	log.Fatal(err)
}

// Fix orientation from EXIF metadata
if err := image.AutoRotate(); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(vips.NewJpegExportParams())
if err != nil {
	log.Fatal(err)
}
os.WriteFile("output.jpg", buf, 0644)
```

### 2. Resize an image

Scale an image down by 50% using the Lanczos3 kernel (the sharpest option).

```go
image, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

// Scale to 50%
if err := image.Resize(0.5, vips.KernelLanczos3); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(&vips.JpegExportParams{Quality: 85})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("resized.jpg", buf, 0644)
```

### 3. Create a thumbnail with smart crop

`NewThumbnailFromFile` is the fastest way to generate thumbnails. It decodes only the pixels it needs.

```go
// Load and shrink to fit within 200x200, cropping to the most interesting region
image, err := vips.NewThumbnailFromFile("photo.jpg", 200, 200, vips.InterestingAttention)
if err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(&vips.JpegExportParams{Quality: 80})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("thumb.jpg", buf, 0644)
```

### 4. Convert between formats

Load a JPEG and export it as WebP and PNG.

```go
image, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

// Export as WebP (lossy)
webpBuf, _, err := image.ExportWebp(&vips.WebpExportParams{Quality: 75})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("photo.webp", webpBuf, 0644)

// Export as PNG
pngBuf, _, err := image.ExportPng(&vips.PngExportParams{Compression: 6})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("photo.png", pngBuf, 0644)
```

### 5. Crop and extract a region

Pull out a specific rectangle from an image.

```go
image, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

// Extract a 300x300 region starting at (50, 100)
if err := image.ExtractArea(50, 100, 300, 300); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportPng(vips.NewPngExportParams())
if err != nil {
	log.Fatal(err)
}
os.WriteFile("cropped.png", buf, 0644)
```

### 6. Blur and sharpen

Apply a Gaussian blur or sharpen an image.

```go
image, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

// Gaussian blur with sigma=3.0
if err := image.GaussianBlur(3.0); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(vips.NewJpegExportParams())
if err != nil {
	log.Fatal(err)
}
os.WriteFile("blurred.jpg", buf, 0644)

// Or sharpen instead: sigma=1.0, x1=2.0, m2=3.0
image2, _ := vips.NewImageFromFile("photo.jpg")
if err := image2.Sharpen(1.0, 2.0, 3.0); err != nil {
	log.Fatal(err)
}

buf2, _, _ := image2.ExportJpeg(vips.NewJpegExportParams())
os.WriteFile("sharpened.jpg", buf2, 0644)
```

### 7. Rotate and flip

Rotate by fixed angles and flip along an axis.

```go
image, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

// Rotate 90 degrees clockwise
if err := image.Rotate(vips.Angle90); err != nil {
	log.Fatal(err)
}

// Flip horizontally
if err := image.Flip(vips.DirectionHorizontal); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(vips.NewJpegExportParams())
if err != nil {
	log.Fatal(err)
}
os.WriteFile("rotated.jpg", buf, 0644)
```

### 8. Adjust brightness, saturation, and hue

`Modulate` works in the LCH color space. Brightness and saturation are multipliers (1.0 = no change), hue is an angle shift in degrees.

```go
image, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

// Bump brightness by 20%, desaturate by 30%, shift hue by 45 degrees
if err := image.Modulate(1.2, 0.7, 45); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(&vips.JpegExportParams{Quality: 90})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("adjusted.jpg", buf, 0644)
```

### 9. Composite two images (watermark overlay)

Layer one image on top of another using Porter-Duff blending.

```go
base, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

overlay, err := vips.NewImageFromFile("watermark.png")
if err != nil {
	log.Fatal(err)
}

// Place the watermark at position (20, 20) using "over" blending
if err := base.Composite(overlay, vips.BlendModeOver, 20, 20); err != nil {
	log.Fatal(err)
}

buf, _, err := base.ExportJpeg(&vips.JpegExportParams{Quality: 90})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("watermarked.jpg", buf, 0644)
```

### 10. Add a text label

Overlay text directly onto an image.

```go
image, err := vips.NewImageFromFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

err = image.Label(&vips.LabelParams{
	Text:      "govips",
	Font:      "sans bold 16",
	OffsetX:   vips.ValueOf(20),
	OffsetY:   vips.ValueOf(20),
	Opacity:   0.8,
	Color:     vips.Color{R: 255, G: 255, B: 255},
	Alignment: vips.AlignLow,
	Width:     vips.ValueOf(200),
	Height:    vips.ValueOf(40),
})
if err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportPng(vips.NewPngExportParams())
if err != nil {
	log.Fatal(err)
}
os.WriteFile("labeled.png", buf, 0644)
```

### 11. Flatten transparency onto a background color

Remove the alpha channel by compositing onto a solid color.

```go
image, err := vips.NewImageFromFile("logo.png")
if err != nil {
	log.Fatal(err)
}

// Flatten alpha onto white
if err := image.Flatten(&vips.Color{R: 255, G: 255, B: 255}); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(&vips.JpegExportParams{Quality: 90})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("flattened.jpg", buf, 0644)
```

### 12. Embed an image in a larger canvas

Center an image inside a larger canvas with a colored background.

```go
image, err := vips.NewImageFromFile("icon.png")
if err != nil {
	log.Fatal(err)
}

w, h := image.Width(), image.Height()

// Center the image inside a 800x600 canvas with a dark background
if err := image.EmbedBackgroundRGBA(
	(800-w)/2, (600-h)/2, 800, 600,
	&vips.ColorRGBA{R: 30, G: 30, B: 30, A: 255},
); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportPng(vips.NewPngExportParams())
if err != nil {
	log.Fatal(err)
}
os.WriteFile("embedded.png", buf, 0644)
```

### 13. Load from a byte buffer and strip metadata

Useful when you're reading images from HTTP requests or databases.

```go
inputBytes, err := os.ReadFile("photo.jpg")
if err != nil {
	log.Fatal(err)
}

image, err := vips.NewImageFromBuffer(inputBytes)
if err != nil {
	log.Fatal(err)
}

// Strip all EXIF/metadata for privacy
if err := image.RemoveMetadata(); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportJpeg(&vips.JpegExportParams{
	Quality:       80,
	StripMetadata: true,
})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("clean.jpg", buf, 0644)
```

### 14. Build a pipeline: thumbnail, sharpen, and export as AVIF

Chain multiple operations together. Each method mutates the image in place, so you can pipeline them naturally.

```go
image, err := vips.NewThumbnailFromFile("photo.jpg", 800, 600, vips.InterestingCentre)
if err != nil {
	log.Fatal(err)
}

if err := image.Sharpen(0.7, 1.0, 2.0); err != nil {
	log.Fatal(err)
}

buf, _, err := image.ExportAvif(&vips.AvifExportParams{
	Quality: 50,
	Effort:  4,
})
if err != nil {
	log.Fatal(err)
}
os.WriteFile("output.avif", buf, 0644)
```

See the _examples/_ folder for more.

## Running tests

```bash
$ make test
```

## Code Generation

The built-in generator lives in `cmd/vipsgen/` and uses libvips GObject introspection to discover operations and their arguments at build time. To regenerate the bindings:

```bash
go generate ./vips/
```

This produces the generated files in `vips/gen_*.{c,h,go}`. You should not need to regenerate unless you are adding support for new libvips operations.

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

### Jemalloc
If the arena option doesn't help, you can try replacing the standard allocator with `jemalloc`,
which emphasizes fragmentation avoidance and scalable concurrency support.

To do this, you need to install the `libjemalloc-dev` package. 
And pass the following flags for build command:

```
CGO_CFLAGS="-fno-builtin-malloc -fno-builtin-calloc -fno-builtin-realloc -fno-builtin-free" CGO_LDFLAGS="-ljemalloc" go build
```


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
