package main

import (
	"fmt"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

/* govips example: how to perform streaming operations
 *
 * Opens a file reader, prints out header info, autorotates and writes out jpeg
 * without buffering the image in memory.
 */

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func main() {

	vips.Startup(nil)
	defer vips.Shutdown()

	inputImage, err := os.Open("input.jpg")
	checkError(err)

	defer inputImage.Close()

	img, err := vips.NewImageFromReader(inputImage)
	checkError(err)

	fmt.Println("Image width: ", img.Metadata().Width)
	fmt.Println("Image height: ", img.Metadata().Height)
	fmt.Println("Image orientation: ", img.Metadata().Orientation)

	err = img.AutoRotate()
	checkError(err)

	outImage, err := os.Create("streaming-example-reorient.jpg")
	checkError(err)

	defer outImage.Close()

	params := vips.NewDefaultJPEGExportParams()
	_, err = img.ExportWriter(outImage, params)
	checkError(err)

}
