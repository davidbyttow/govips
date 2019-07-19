package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wix-playground/govips/vips"
)

var (
	flagIn          = flag.String("in", "", "file to load")
	flagOut         = flag.String("out", "", "file to write out")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "resize -in input_file -out output_file")
	}
	flag.Parse()

	vips.Startup(nil)
	_ = resize(*flagIn, *flagOut)
	vips.Shutdown()

	if *reportLeaksFlag {
		vips.PrintObjectReport("resize")
	}
}

func resize(inputFile, outputFile string) error {
	img, err := vips.NewImageFromFile(inputFile)
	if err != nil {
		return err
	}
	defer img.Close()

	b, _, err := vips.NewTransform().Scale(0.2).ApplyAndExport(img)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputFile, b, os.ModeAppend)
}
