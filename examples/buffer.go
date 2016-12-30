package buffer

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidbyttow/gimage"
	"github.com/davidbyttow/gomore/io"
)

const usage = "avg [input]"

func run(file string) error {
	buf, err := io.ReadFile(file)
	if err != nil {
		return err
	}

	image, err := gimage.NewImageFromBuffer(buf)
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

	defer gimage.Shutdown()
	err := run(file)
	if err != nil {
		os.Exit(1)
	}
}
