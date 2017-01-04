package examples

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

	image, err := govips.NewImageFromBuffer(buf, nil)
	if err != nil {
		return err
	}

	avg := image.Avg(nil)
	fmt.Printf("avg=%0.2f\n", avg)

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "avg [input]")
	}
	flag.Parse()

	file := os.Args[1]

	defer govips.Shutdown()
	err := loadAverage(file)
	if err != nil {
		os.Exit(1)
	}
}
