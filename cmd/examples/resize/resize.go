package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jhford/govips/pkg/vips"
)

var (
	flagIn          = flag.String("in", "", "file to load")
	flagOut         = flag.String("out", "", "file to write out")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "resize -in input_file -out output_file")
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
	_, _, err := vips.NewTransform().
		LoadFile(inputFile).
		Scale(0.2).
		OutputFile(outputFile).
		Apply()
	return err
}
