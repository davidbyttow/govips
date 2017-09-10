package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidbyttow/govips"
)

func run(inputFile, outputFile string) error {
	in, err := govips.NewImageFromFile(inputFile,
		govips.IntInput("access", int(govips.AccessSequentialUnbuffered)))
	if err != nil {
		return err
	}

	out := in.Embed(10, 10, 1000, 1000,
		govips.IntInput("extend", int(govips.ExtendCopy)))

	out.WriteToFile(outputFile)

	return nil
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

	govips.Startup(nil)
	defer govips.Shutdown()

	err := run(*flagIn, *flagOut)
	if err != nil {
		os.Exit(1)
	}
}
