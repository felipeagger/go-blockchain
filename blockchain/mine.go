package blockchain

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MinedResult struct {
	hash      string
	nonce     string
	genHashes int
}

func (b *Block) SequentialMine(difficulty int) (int, string) {
	started := time.Now()
	qtdHashes := 0
	var hash string
	nonce := b.Nonce
	date := b.Timestamp.Format("2006-01-02 15:04:05")
	for !strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
		nonce++
		hash = CalculateHash(date, strconv.Itoa(nonce), b.PreviousHash, b.GetDataInBytes())
		qtdHashes = qtdHashes + 1
	}
	PrintHashRate(qtdHashes, time.Now().Sub(started))
	return nonce, hash
}

func generateNonce(c context.Context, ch chan<- string, w *sync.WaitGroup, currNonce int) {
	defer w.Done()
	nonce := currNonce
	for {
		select {
		case <-c.Done():
			//fmt.Println("\nshutting down nonce worker")
			return
		default:
			nonce += 1
			ch <- strconv.Itoa(nonce)
		}
	}
}

func (b *Block) Mine(difficulty int) (int, string) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	ch := make(chan MinedResult, 1)
	nCores := runtime.NumCPU()
	nonceCh := make(chan string, nCores*10)

	started := time.Now()
	wg.Add(1)
	go generateNonce(ctx, nonceCh, &wg, b.Nonce)

	for i := 0; i < nCores; i++ {
		wg.Add(1)
		go func(c context.Context, nonceChan <-chan string, block Block, startedTime time.Time) {
			var qtdHashes int
			defer func() {
				//PrintHashRate(qtdHashes, time.Now().Sub(startedTime))
				wg.Done()
			}()

			var hash string
			prefix := strings.Repeat("0", difficulty)
			date := block.Timestamp.Format("2006-01-02 15:04:05")

			for {
				select {
				case <-c.Done():
					return

				case nonce, ok := <-nonceChan:
					if !ok {
						return
					}

					//mining
					qtdHashes++
					hash = CalculateHash(date, nonce, block.PreviousHash, block.GetDataInBytes())

					if strings.HasPrefix(hash, prefix) {
						ch <- MinedResult{hash: hash, nonce: nonce, genHashes: qtdHashes}
						return
					}
				}

			}

		}(ctx, nonceCh, *b, started)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(nonceCh)
	}()

	// Process hashes
	var mined MinedResult
	for minedRes := range ch {
		if strings.HasPrefix(minedRes.hash, strings.Repeat("0", difficulty)) {
			//fmt.Println("Hash Gerado", minedRes.hash)
			mined = minedRes
			break
		}
	}

	cancel()
	PrintHashRate(mined.genHashes, time.Now().Sub(started))
	nonce, _ := strconv.Atoi(mined.nonce)
	return nonce, mined.hash
}
