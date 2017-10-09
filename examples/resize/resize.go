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
		fmt.Fprintf(os.Stderr, "resize -in input_file -out output_file")
	}
	flag.Parse()

	vips.Startup(nil)
	resize(*flagIn, *flagOut)
	vips.Shutdown()

	vips.PrintObjectReport("resize")
}

func resize(inputFile, outputFile string) error {
	_, err := vips.NewTransform().
		LoadFile(inputFile).
		Scale(0.2).
		OutputFile(outputFile).
		Apply()
	return err
}
