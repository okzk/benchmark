package main

import (
	"flag"
	"fmt"
	"github.com/okzk/benchmark"
	"net/http"
	"os"
)

var (
	samples     int
	concurrency int
)

func init() {
	flag.IntVar(&samples, "n", 1, "Number of samples")
	flag.IntVar(&concurrency, "c", 1, "Concurrency")
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "URL required.")
		os.Exit(1)
	}
	url := flag.Arg(0)

	f := func(_ int, _ benchmark.UserData) (int, string) {
		res, err := http.Get(url)
		if err != nil {
			return 0, url
		}
		res.Body.Close()
		return res.StatusCode, url
	}

	b := &benchmark.Benchmark{Samples: samples, Concurrency: concurrency, TestFunc: f}

	r := b.Run()
	fmt.Print(r.FormatByStatus())
}
