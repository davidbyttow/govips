package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davidbyttow/govips"
)

func run(inputFile, outputFile string) error {
	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	in, err := vips.NewImageFromBuffer(buf,
		vips.IntInput("access", int(vips.AccessSequentialUnbuffered)))
	if err != nil {
		return err
	}

	out := in.Embed(10, 10, 1000, 1000,
		vips.IntInput("extend", int(vips.ExtendCopy)))

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
		fmt.Fprintf(os.Stderr, "embed -in input_file -out output_file")
	}
	flag.Parse()

	vips.Startup(nil)
	defer vips.Shutdown()

	err := run(*flagIn, *flagOut)
	if err != nil {
		os.Exit(1)
	}
}
