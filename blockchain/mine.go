package blockchain

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type MinedResult struct {
	hash      string
	pow       int
	genHashes int
}

func (b *Block) SequentialMine(difficulty int) (int, string) {
	started := time.Now()
	qtdHashes := 0
	var hash string
	pow := b.Pow
	for !strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
		pow++
		hash = CalculateHash(b.Timestamp, pow, b.PreviousHash, b.Data)
		qtdHashes = qtdHashes + 1
	}
	PrintHashRate(qtdHashes, time.Now().Sub(started))
	return pow, hash
}

func generatePow(c context.Context, ch chan<- int, w *sync.WaitGroup, currPow int) {
	defer w.Done()
	pow := currPow
	for {
		select {
		case <-c.Done():
			fmt.Println("\nshutting down pow worker")
			return
		default:
			pow += 1
			ch <- pow
		}
	}
}

func (b *Block) Mine(difficulty int) (int, string) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	ch := make(chan MinedResult, 1)
	nCores := runtime.NumCPU()
	powCh := make(chan int, nCores*10)

	started := time.Now()
	wg.Add(1)
	go generatePow(ctx, powCh, &wg, b.Pow)

	for i := 0; i < nCores; i++ {
		wg.Add(1)
		go func(c context.Context, powChan <-chan int, block Block, startedTime time.Time) {
			var qtdHashes int
			defer func() {
				PrintHashRate(qtdHashes, time.Now().Sub(startedTime))
				wg.Done()
			}()

			var hash string
			prefix := strings.Repeat("0", difficulty)

			for {
				select {
				case <-c.Done():
					return

				case pow, ok := <-powChan:
					if !ok {
						return
					}

					//mining
					qtdHashes++
					hash = CalculateHash(block.Timestamp, pow, block.PreviousHash, block.Data)

					if strings.HasPrefix(hash, prefix) {
						ch <- MinedResult{hash: hash, pow: pow, genHashes: qtdHashes}
						return
					}
				}

			}

		}(ctx, powCh, *b, started)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(powCh)
	}()

	// Process hashes
	var mined MinedResult
	for minedRes := range ch {
		if strings.HasPrefix(minedRes.hash, strings.Repeat("0", difficulty)) {
			fmt.Println("Hash Gerado", minedRes.hash)
			mined = minedRes
			break
		}
	}

	cancel()
	PrintHashRate(mined.genHashes, time.Now().Sub(started))
	return mined.pow, mined.hash
}
