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

	inputImage, err := vips.NewImageFromFile("./jpg-24bit-icc-adobe-rgb.jpg")
	checkError(err)
	defer inputImage.Close()
	checkError(inputImage.OptimizeICCProfile())

	params := vips.NewJpegExportParams()
	params.RestartInterval = 1

	imageBytes, _, err := inputImage.ExportJpeg(params)
	checkError(err)
	checkError(ioutil.WriteFile("./result/rst-mrk-output-govips.jpeg", imageBytes, 0644))
}
