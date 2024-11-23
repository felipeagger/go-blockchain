package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	counterOne int
	atmCount   atomic.Int32
	mu         sync.Mutex
)

func incrementWithLockContention(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		mu.Lock()
		counterOne++
		//time.Sleep(5 * time.Millisecond)
		mu.Unlock()
	}
}

func incrementWithoutLockContention(ch chan int, w *sync.WaitGroup) {
	defer w.Done()
	for i := 0; i < 100; i++ {
		ch <- 1
	}
}

func incrementWithAtomic(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		atmCount.Add(1)
	}
}

// run bench: go test -bench=. -cpuprofile cpu.prof -memprofile mem.prof -blockprofile block.out -mutexprofile mutex.out -race -benchmem -count=3 -v
/*func BenchmarkLockContention(b *testing.B) {
	var wg sync.WaitGroup
	now := time.Now()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go incrementWithLockContention(&wg)
	}

	wg.Wait()
	dur := time.Now().Sub(now)
	fmt.Printf("\nFinal counterOne value: %v; duration: %d ns\n", counterOne, dur.Nanoseconds())
}

/*func BenchmarkLockContentionAtomic(b *testing.B) {
	var wg sync.WaitGroup
	now := time.Now()

	atmCount.Store(0)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go incrementWithAtomic(&wg)
	}

	wg.Wait()
	dur := time.Now().Sub(now)
	fmt.Printf("\nFinal atmCount value: %v; duration: %d ns\n", atmCount.Load(), dur.Nanoseconds())
}*/

// run bench: go test -bench=. -cpuprofile cpu.prof -memprofile mem.prof -blockprofile block.out -mutexprofile mutex.out -race -benchmem -count=1 -benchtime=1000 -v
func BenchmarkWithoutLockContention(b *testing.B) {
	var wgTwo sync.WaitGroup
	counterTwo := 0
	ch := make(chan int, 1000000)

	now := time.Now()

	for i := 0; i < b.N; i++ {
		wgTwo.Add(1)
		go incrementWithoutLockContention(ch, &wgTwo)
	}

	go func() {
		wgTwo.Wait()
		close(ch)
	}()

	for val := range ch {
		counterTwo += val
	}

	dur := time.Now().Sub(now)
	fmt.Printf("\nFinal counterTwo value: %v; duration: %d ns\n", counterTwo, dur.Nanoseconds())
}
