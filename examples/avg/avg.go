package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davidbyttow/govips"
)

func loadAverage(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	image, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		return err
	}

	avg := image.Avg()
	fmt.Printf("avg=%0.2f\n", avg)

	return nil
}

var (
	flagFile = flag.String("file", "", "file to compute average for")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "avg -file input_file")
	}
	flag.Parse()

	vips.Startup(nil)
	defer vips.Shutdown()

	err := loadAverage(*flagFile)
	if err != nil {
		os.Exit(1)
	}
}
