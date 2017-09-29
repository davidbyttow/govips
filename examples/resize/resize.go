package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davidbyttow/govips"
)

func run(inputFile, outputFile string) error {
	buf, _ := ioutil.ReadFile(inputFile)
	in, err := vips.NewImageFromBuffer(buf,
		vips.IntInput("access", int(vips.AccessSequential)))
	if err != nil {
		return err
	}

	interp, err := vips.NewInterpolator("nohalo")
	if err != nil {
		return err
	}

	out := in.Resize(0.2, vips.InterpolatorInput("interpolate", interp))

	buf, err = out.Export(vips.ExportOptions{})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFile, buf, 0644)
}

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
	defer vips.Shutdown()

	err := run(*flagIn, *flagOut)
	if err != nil {
		os.Exit(1)
	}
}
