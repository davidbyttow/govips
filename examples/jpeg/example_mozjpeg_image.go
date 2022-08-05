// NOTE: Run from project root directory

package main

import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io/ioutil"
	"os"
	"runtime"

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
func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	createCoverImage("1.pdf", "1-1.jpeg")
	return

	vips.Startup(nil)
	defer vips.Shutdown()

	inputImage, err := vips.NewImageFromFile("1.pdf")
	checkError(err)
	defer inputImage.Close()
	//checkError(inputImage.OptimizeICCProfile())

	ep := vips.NewJpegExportParams()
	//ep.StripMetadata = true
	ep.Quality = 100
	//ep.Interlace = true
	//ep.OptimizeCoding = true
	//ep.SubsampleMode = vips.VipsForeignSubsampleAuto
	//ep.TrellisQuant = true
	//ep.OvershootDeringing = true
	//ep.OptimizeScans = true
	//ep.QuantTable = 3

	imageBytes, _, err := inputImage.ExportJpeg(ep)
	checkError(err)
	checkError(ioutil.WriteFile("govips.jpeg", imageBytes, 0644))
}

func clearImagickWand(mw *imagick.MagickWand) {
	mw.Clear()
	mw.Destroy()
	runtime.SetFinalizer(mw, nil)
	mw = nil
}

func createCoverImage(pathNoExtension string, coverPathName string) bool {
	//sourceImagePath := getSourceImageForCover(filepath.Dir(pathNoExtension))
	mw := imagick.NewMagickWand()
	mw.SetResolution(300, 300)

	err := mw.ReadImage(pathNoExtension)
	if err != nil {
		clearImagickWand(mw)
		return false
	}
	pix := imagick.NewPixelWand()
	pix.SetColor("white")
	mw.SetBackgroundColor(pix)
	mw.SetImageBackgroundColor(pix)
	mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE)
	width, height, err := mw.GetSize()
	if err != nil {
		clearImagickWand(mw)
		return false
	}
	fmt.Println("width", width)
	fmt.Println("height", height)

	// create thumb if we are requesting cover
	//if width > 320 {
	//	scaleRatio := 320 / width
	//	width = width * scaleRatio
	//	height = height * scaleRatio
	//
	//	err = mw.ResizeImage(width, height, imagick.FILTER_LANCZOS, -0.1)
	//	if err != nil {
	//		clearImagickWand(mw)
	//		return false
	//	}
	//}

	mw.SetImageFormat("png")
	err = mw.WriteImage(coverPathName)
	if err != nil {
		clearImagickWand(mw)
		return false
	}

	clearImagickWand(mw)

	return true
}
