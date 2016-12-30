package avg

import (
	"flag"
	"fmt"
	"os"

	"github.com/davidbyttow/gimage"
	"github.com/davidbyttow/gomore/io"
)

const usage = "avg [input]"

func loadAverage(file string) error {
	buf, err := io.ReadFile(file)
	if err != nil {
		return err
	}

	image, err := gimage.NewImageFromBuffer(buf)
	if err != nil {
		return err
	}

	avg := image.Avg(nil)
	fmt.Printf("avg=%0.2f\n", avg)

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	flag.Parse()

	file := os.Args[1]

	defer gimage.Shutdown()
	err := loadAverage(file)
	if err != nil {
		os.Exit(1)
	}
}
