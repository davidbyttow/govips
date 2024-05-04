// NOTE: Run from project root directory

package package_test

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

	inputImage, err := vips.NewImageFromFile("examples/tiff/input.jpg")
	checkError(err)
	defer inputImage.Close()

	fmt.Println(inputImage.Height(), inputImage.Width())
	err = inputImage.Resize(0.5, vips.KernelLanczos3)
	checkError(err)
	fmt.Println(inputImage.Height(), inputImage.Width())

	exportParams := vips.NewTiffExportParams()
	exportParams.Quality = 100
	imageBytes, _, err := inputImage.ExportTiff(exportParams)
	err = os.WriteFile("examples/tiff/output-govips.tiff", imageBytes, 0644)
	checkError(err)
}
