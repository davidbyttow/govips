package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidbyttow/govips"
)

var (
	flagIn  = flag.String("in", "", "file to load")
	flagOut = flag.String("out", "", "file to write out")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "embed -in input_file -out output_file")
	}
	flag.Parse()

	vips.Startup(nil)
	embed(*flagIn, *flagOut)
	vips.Shutdown()

	vips.PrintObjectReport("embed")
}

func embed(inputFile, outputFile string) error {
	return vips.NewPipeline().
		LoadFile(inputFile).
		PadStrategy(vips.ExtendBlack).
		Resize(1200, 1200).
		OutputFile(outputFile)
}
