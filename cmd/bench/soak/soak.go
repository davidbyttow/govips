package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jhford/govips/pkg/vips"
)

var (
	batchFlag        = flag.Int("batch", 5, "Transforms per batch")
	varianceFlag     = flag.Int("variance", 10, "Target height and width variance")
	delayFlag        = flag.Int("delay", 10, "Delay in milliseconds between batches")
	limitFlag        = flag.Int("limit", 100, "Total number of images to process. 0 = infinite")
	widthFlag        = flag.Int("width", 300, "Target width of each image")
	heightFlag       = flag.Int("height", 300, "Target height of each image")
	cacheFlag        = flag.Bool("cache", true, "Cache remote images")
	concurrencyFlag  = flag.Int("concurrency", 1, "Concurrency level")
	maxCacheMemFlag  = flag.Int("maxCacheMem", 0, "Max cache memory")
	maxCacheSizeFlag = flag.Int("maxCacheSize", 0, "Max cache size")
	imageFlag        = flag.String("image", "cmd/bench/soak/dwi.jpg", "Image file to load")
	cpuProfileFlag   = flag.String("cpuprofile", "", "write cpu profile `file`")
	memProfileFlag   = flag.String("memprofile", "", "write memory profile to `file`")
)

func main() {
	flag.Parse()

	if *cpuProfileFlag != "" {
		f, err := os.Create(*cpuProfileFlag)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	cfg := &vips.Config{
		ConcurrencyLevel: *concurrencyFlag,
		MaxCacheMem:      *maxCacheMemFlag,
		MaxCacheSize:     *maxCacheSizeFlag,
		CollectStats:     true,
	}
	vips.Startup(cfg)

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1)

	outputStats()

	soak()

	vips.ShutdownThread()
	vips.Shutdown()
	runtime.GC()
	vips.PrintObjectReport("Soak test")
	outputStats()

	if *memProfileFlag != "" {
		f, err := os.Create(*memProfileFlag)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

func soak() {
	batch := *batchFlag
	delay := time.Millisecond * time.Duration(*delayFlag)
	width := *widthFlag
	height := *heightFlag
	variance := *varianceFlag
	limit := int64(*limitFlag)
	var total int64

	fmt.Printf("Running soak with caching=%t batch=%d delay=%d width=%d height=%d variance=%d limit=%d\n",
		*cacheFlag, batch, delay, width, height, variance, limit)

	for limit == 0 || total < limit {
		var wg sync.WaitGroup
		for i := 0; i < batch; i++ {
			wg.Add(1)
			go func(i int) {
				defer func() {
					if err := recover(); err != nil {
						log.Print(err)
					}
				}()
				defer wg.Done()
				count := int(atomic.AddInt64(&total, 1))
				w, h := width, height
				if variance > 0 {
					w += rand.Intn(variance)
					h += rand.Intn(variance)
				}
				buf := loadFile(*imageFlag)
				_, _, err := vips.NewTransform().
					LoadBuffer(buf).
					Resize(w, h).
					OutputBytes().
					Apply()

				if err != nil {
					log.Print(err)
				}

				if count%batch == 0 {
					log.Printf("Processed %d...\n", count)
				}
			}(i)
			time.Sleep(delay)
		}
		wg.Wait()
	}
}

var cachedImages = map[string][]byte{}

func loadFile(file string) []byte {
	if *cacheFlag {
		if buf, ok := cachedImages[file]; ok {
			return buf
		}
	}
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if *cacheFlag {
		cachedImages[file] = buf
	}
	return buf
}

func outputStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	var vipsMem vips.VipsMemoryStats
	vips.ReadVipsMemStats(&vipsMem)

	var stats vips.RuntimeStats
	vips.ReadRuntimeStats(&stats)

	fmt.Println(stats)

	lines := []string{
		fmt.Sprintf("alloc=%d\ntotal=%d\nsys=%d\ninuse=%d", mem.Alloc, mem.TotalAlloc, mem.Sys, mem.HeapInuse),
		fmt.Sprintf("vips mem=%d high=%d allocs=%d", vipsMem.Mem, vipsMem.MemHigh, vipsMem.Allocs),
	}
	fmt.Println(strings.Join(lines, "\n"))
}
