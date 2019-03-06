package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/davidbyttow/govips/pkg/vips"
)

var (
	flagIn          = flag.String("in", "", "file to load")
	flagWidth       = flag.Int("width", 352, "thumbnail width")
	flagHeight      = flag.Int("height", -1, "thumbnail height")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "thumbnail -in input_file -width width -height height\n")
	}
	flag.Parse()

	vips.Startup(nil)
	thumbnail(*flagIn, *flagWidth, *flagHeight)
	vips.Shutdown()

	if *reportLeaksFlag {
		vips.PrintObjectReport("thumbnail")
	}
}

func thumbnail(inputFile string, width, height int) error {
	dir, file := path.Split(inputFile)
	outFile := path.Join(dir, fmt.Sprintf("tn_%s", file))
	outImage, err := vips.ThumbnailWithSize(inputFile, width, height)
	if err != nil {
		return err
	}
	return vips.Jpegsave(outImage, outFile)
}
