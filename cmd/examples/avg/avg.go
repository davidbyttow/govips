package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wix-playground/govips/vips"
)

var (
	flagFile        = flag.String("file", "", "file to compute average for")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "avg -file input_file")
	}
	flag.Parse()

	vips.Startup(nil)
	_ = avg(*flagFile)
	vips.ShutdownThread()
	vips.Shutdown()

	if *reportLeaksFlag {
		vips.PrintObjectReport("avg")
	}
}

func avg(file string) error {
	image, err := vips.NewImageFromFile(file)
	if err != nil {
		return err
	}
	defer image.Close()

	avg, err := image.Avg()
	if err != nil {
		return err
	}

	fmt.Printf("avg=%0.2f\n", avg)

	return nil
}
