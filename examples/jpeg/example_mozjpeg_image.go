// NOTE: Run from project root directory

package vips_test

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

// Mimick mozjpeg’s parameters for a JPEG export.
// Quality is 75/100. Progressive (interlaced), optimized, separate DCT scans,
// trellis optimization, overshoot deringing, quant table 3.
// These mozjpeg options are documented here:
// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-jpegsave
func ExampleMozJPEGEncode() {
	vips.Startup(nil)
	defer vips.Shutdown()

	inputImage, err := vips.NewImageFromFile("resources/jpg-24bit-icc-adobe-rgb.jpg")
	checkError(err)
	defer inputImage.Close()
	checkError(inputImage.OptimizeICCProfile())

	ep := vips.NewJpegExportParams()
	ep.StripMetadata = true
	ep.Quality = 75
	ep.Interlace = true
	ep.OptimizeCoding = true
	ep.SubsampleMode = vips.VipsForeignSubsampleAuto
	ep.TrellisQuant = true
	ep.OvershootDeringing = true
	ep.OptimizeScans = true
	ep.QuantTable = 3

	imageBytes, _, err := inputImage.ExportJpeg(ep)
	checkError(err)
	checkError(os.WriteFile("examples/jpeg/mozjpeg-output-govips.jpeg", imageBytes, 0644))
}
