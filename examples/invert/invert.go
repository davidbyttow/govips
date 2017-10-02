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
		fmt.Fprintf(os.Stderr, "invert -in input_file -out output_file")
	}
	flag.Parse()

	vips.Startup(nil)
	invert(*flagIn, *flagOut)
	vips.Shutdown()

	vips.PrintObjectReport("invert")
}

func invert(inputFile, outputFile string) error {
	return vips.NewPipeline().
		LoadFile(inputFile).
		Invert().
		OutputFile(outputFile)
}
