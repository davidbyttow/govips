package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidbyttow/govips/pkg/vips"
)

var (
	flagIn          = flag.String("in", "", "file to load")
	flagOut         = flag.String("out", "", "file to write out")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "embed -in input_file -out output_file")
	}
	flag.Parse()

	vips.Startup(nil)
	embed(*flagIn, *flagOut)
	vips.Shutdown()

	if *reportLeaksFlag {
		vips.PrintObjectReport("invert")
	}

	vips.PrintObjectReport("embed")
}

func embed(inputFile, outputFile string) error {
	_, err := vips.NewTransform().
		LoadFile(inputFile).
		Resize(1200, 1200).
		Embed(vips.ExtendBlack).
		OutputFile(outputFile).
		Apply()
	return err
}
