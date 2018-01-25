package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidbyttow/govips"
)

var (
	flagIn          = flag.String("in", "", "file to load")
	flagOut         = flag.String("out", "", "file to write out")
	page         		= flag.Int("page", 1, "page number in PDF")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "pdf -in input_file -out output_file -page 2")
	}
	flag.Parse()

	vips.Startup(nil)
	resize(*flagIn, *flagOut)
	vips.Shutdown()

	if *reportLeaksFlag {
		vips.PrintObjectReport("resize")
	}
}

func resize(inputFile, outputFile string) error {
	_, err := vips.NewTransform().
		LoadFile(inputFile).
		Page(0).
		LoadScale(2).
		OutputFile(outputFile).
		Apply()
	return err
}
