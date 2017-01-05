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

	image, err := govips.NewImageFromBuffer(buf, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d x %d pixel image from %s\n",
		image.Width(), image.Height(), inputFile)

	buf, err = image.WriteToBuffer(".png", nil)
	if err != nil {
		return err
	}

	fmt.Printf("Written to memory %p in png format, %d bytes\n", buf, len(buf))

	image, err = govips.NewImageFromBuffer(buf, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded from memory, %d x %d pixel image\n", image.Width(), image.Height())

	image.WriteToFile(outputFile, nil)
	fmt.Printf("Written back to %s\n", outputFile)

	return nil
}

var (
	flagIn  = flag.String("in", "", "file to load")
	flagOut = flag.String("out", "", "file to write out")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "buffer -in input_file -out output_file")
	}
	flag.Parse()

	defer govips.Shutdown()
	err := run(*flagIn, *flagOut)
	if err != nil {
		os.Exit(1)
	}
}
