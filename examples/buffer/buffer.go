package examples

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davidbyttow/govips"
)

const usage = "avg [input]"

func run(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	image, err := govips.NewImageFromBuffer(buf, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d x %d pixel image from %s\n",
		image.Width(), image.Height(), file)

	// TODO(d): Resave
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	flag.Parse()

	file := os.Args[1]

	defer govips.Shutdown()
	err := run(file)
	if err != nil {
		os.Exit(1)
	}
}
