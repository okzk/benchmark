package benchmark

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBenchRun(_ *testing.T) {
	if testing.Short() {
		return
	}

	f := func(_ int, _ UserData) (int, string) {
		time.Sleep(time.Duration(rand.Int63n(100 * 1000 * 1000)))
		return rand.Intn(2), "sleep test"
	}
	bench := Benchmark{Samples: 500, Concurrency: 50, TestFunc: f}
	r := bench.Run()

	fmt.Print(r.FormatByStatus())
}
