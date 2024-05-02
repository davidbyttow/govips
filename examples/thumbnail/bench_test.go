package thumbnail

import (
	"fmt"
	"os"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
)

var file = "../../resources/jpg-24bit-icc-iec.jpg"

func init() {
	vips.LoggingSettings(func(domain string, level vips.LogLevel, msg string) {
		fmt.Println(domain, level, msg)
	}, vips.LogLevelError)
	vips.Startup(nil)
}

func BenchmarkNewImageFromFile(b *testing.B) {
	resizeToTest := func(size int) {
		image, err := vips.NewImageFromFile(file)
		if err != nil {
			panic(err)
		}
		if err := image.Thumbnail(size*100, size*100, vips.InterestingCentre); err != nil {
			panic(err)
		}
		if err := image.Flip(vips.DirectionVertical); err != nil {
			panic(err)
		}
		if _, _, err = image.ExportNative(); err != nil {
			panic(err)
		}
	}
	for i := 0; i < b.N; i++ {
		resizeToTest(1)
		resizeToTest(2)
		resizeToTest(3)
		resizeToTest(4)
		resizeToTest(5)
	}
}

func BenchmarkNewImageFromBuffer(b *testing.B) {
	resizeToTest := func(size int) {
		buf, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		image, err := vips.NewImageFromBuffer(buf)
		if err != nil {
			panic(err)
		}
		if err := image.Thumbnail(size*100, size*100, vips.InterestingCentre); err != nil {
			panic(err)
		}
		if err := image.Flip(vips.DirectionVertical); err != nil {
			panic(err)
		}
		if _, _, err = image.ExportNative(); err != nil {
			panic(err)
		}
	}
	for i := 0; i < b.N; i++ {
		resizeToTest(1)
		resizeToTest(2)
		resizeToTest(3)
		resizeToTest(4)
		resizeToTest(5)
	}
}

func BenchmarkNewThumbnailFromFile(b *testing.B) {
	resizeToTest := func(size int) {
		image, err := vips.NewThumbnailFromFile(file, size*100, size*100, vips.InterestingCentre)
		if err != nil {
			panic(err)
		}
		if err := image.Flip(vips.DirectionVertical); err != nil {
			panic(err)
		}
		if _, _, err = image.ExportNative(); err != nil {
			panic(err)
		}
	}
	for i := 0; i < b.N; i++ {
		resizeToTest(1)
		resizeToTest(2)
		resizeToTest(3)
		resizeToTest(4)
		resizeToTest(5)
	}
}

func BenchmarkNewThumbnailFromBuffer(b *testing.B) {
	resizeToTest := func(size int) {
		buf, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		image, err := vips.NewThumbnailFromBuffer(buf, size*100, size*100, vips.InterestingCentre)
		if err != nil {
			panic(err)
		}
		if err := image.Flip(vips.DirectionVertical); err != nil {
			panic(err)
		}
		if _, _, err = image.ExportNative(); err != nil {
			panic(err)
		}
	}
	for i := 0; i < b.N; i++ {
		resizeToTest(1)
		resizeToTest(2)
		resizeToTest(3)
		resizeToTest(4)
		resizeToTest(5)
	}
}
