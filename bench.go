package benchmark

import (
	"sync"
	"time"
)

type UserData interface{}

type Benchmark struct {
	Samples          int
	Concurrency      int
	TestFunc         func(int, UserData) (int, string)
	UserDataFactory  func() UserData
	UserDataDisposer func(UserData)
}

type Result struct {
	Status int
	Start  int64
	End    int64
	Info   string
}

type Results []Result

func (b *Benchmark) Run() Results {
	if b.Concurrency <= 0 || b.Samples <= 0 || b.TestFunc == nil {
		panic("invalid configuration")
	}

	userData := make([]UserData, b.Concurrency)
	if b.UserDataFactory != nil {
		for i := 0; i < b.Concurrency; i++ {
			userData[i] = b.UserDataFactory()
		}
		if b.UserDataDisposer != nil {
			defer func() {
				for _, userData := range userData {
					b.UserDataDisposer(userData)
				}
			}()
		}
	}

	var wg sync.WaitGroup
	ch := make(chan int, b.Concurrency*2)
	results := make(Results, b.Samples)
	for i := 0; i < b.Concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			data := userData[i]
			for {
				if n, ok := <-ch; ok {
					start := time.Now()
					status, info := b.TestFunc(n, data)
					end := time.Now()
					results[n] = Result{Status: status, Start: start.UnixNano(), End: end.UnixNano(), Info: info}
				} else {
					return
				}
			}
		}(i)
	}

	for i := 0; i < b.Samples; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()

	return results
}
