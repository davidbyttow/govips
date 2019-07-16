package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wix-playground/govips/vips"
)

var (
	flagIn          = flag.String("in", "", "file to load")
	flagMessage     = flag.String("message", "", "message to write")
	flagOut         = flag.String("out", "", "file to write out")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "text -in input_file -message message -out output_file")
	}
	flag.Parse()

	vips.Startup(nil)
	if err := text(*flagIn, *flagMessage, *flagOut); err != nil {
		panic(err)
	}
	vips.Shutdown()

	if *reportLeaksFlag {
		vips.PrintObjectReport("text")
	}
}

func text(inputFile, message, outputFile string) error {
	_, _, err := vips.NewTransform().
		Label(&vips.LabelParams{
			Text:      message,
			Opacity:   1.0,
			Width:     vips.ScaleOf(0.9),
			Height:    vips.ScaleOf(1.0),
			Alignment: vips.AlignCenter,
			Color:     vips.Color{R: 255, G: 255, B: 255},
		}).
		LoadFile(inputFile).
		OutputFile(outputFile).
		Apply()
	return err
}
