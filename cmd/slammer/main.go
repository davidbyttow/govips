package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"sync"
	"sync/atomic"

	"github.com/wix-playground/govips/pkg/vips"
)

var (
	dirFlag = flag.String("dir", "", "directory to scan for images")
)

func main() {
	flag.Parse()

	vips.Startup(&vips.Config{
		ConcurrencyLevel: 250,
		MaxCacheMem:      2 * 1024 * 1024,
		MaxCacheSize:     1024,
		MaxCacheFiles:    1000,
	})
	defer vips.Shutdown()
	var totalBytes int64

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := run(*dirFlag)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddInt64(&totalBytes, int64(n))
		}()
	}
	wg.Wait()

	vips.PrintObjectReport("slammer")
	fmt.Printf("Total bytes processed: %d\n", totalBytes)
	fmt.Println("Press enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func run(dir string) (int, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, fi := range files {
		x, y := rand.Float64(), rand.Float64()
		f := path.Join(dir, fi.Name())
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return 0, err
		}
		imageRef, err := vips.NewImageFromBuffer(data)
		if err != nil {
			return 0, err
		}
		tx := vips.NewTransform().
			Image(imageRef).
			ResizeStrategy(vips.ResizeStrategyCrop).
			ScaleWidth(x).
			ScaleHeight(y).
			Format(vips.ImageTypeJPEG).
			OutputBytes()
		b, _, err := tx.Apply()
		if err != nil {
			return 0, err
		}
		total += len(b)
		fmt.Printf("Processed file=%s w=%v h=%v bytes=%d\n", f, x, y, len(b))
		imageRef.Close()
	}

	return total, err
}
