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

	f := func(_ int, _ UserData) (Status, OptionalInfo) {
		time.Sleep(time.Duration(rand.Int63n(100 * 1000 * 1000)))
		return Status(rand.Intn(2)), nil
	}
	bench := Benchmark{Samples: 500, Concurrency: 50, TestFunc: f}
	r := bench.Run()

	fmt.Print(r.FormatByStatus())
}
