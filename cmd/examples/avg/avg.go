package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidbyttow/govips/pkg/vips"
)

var (
	flagFile        = flag.String("file", "", "file to compute average for")
	reportLeaksFlag = flag.Bool("leaks", false, "Outputs vips memory")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "avg -file input_file")
	}
	flag.Parse()

	vips.Startup(nil)
	avg(*flagFile)
	vips.ShutdownThread()
	vips.Shutdown()

	if *reportLeaksFlag {
		vips.PrintObjectReport("avg")
	}
}

func avg(file string) error {
	image, err := vips.NewImageFromFile(file, nil)
	if err != nil {
		return err
	}
	defer image.Close()

	avg, _ := vips.Avg(image.Image())
	fmt.Printf("avg=%0.2f\n", avg)
	return nil
}
