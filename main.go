package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"runtime"
	"sync"
	"time"
)

func main() {
	funcSignature := "balanceOf%d(address)"
	searchFuncSelector(funcSignature, runtime.GOMAXPROCS(runtime.NumCPU()))
}

func searchFuncSelector(funcSignature string, numThreads int) {
	start := time.Now()
	var sender sync.WaitGroup
	ch := make(chan string, 100)
	goldenFound := make(chan bool, numThreads)
	tries := ^uint(0) / uint(numThreads)

	for thread := 0; thread < numThreads; thread++ {
		sender.Add(1)
		go runRoutinesFastest(&sender, ch, goldenFound, funcSignature, numThreads, uint(thread), tries)
	}
	go receiveFuncSelectors(ch)
	sender.Wait()
	close(ch)
	fmt.Println(fmt.Sprintf("Complete fastest run in %s", time.Since(start)))
}

func runRoutinesFastest(wg *sync.WaitGroup, ch chan<- string, goldenFound chan bool, funcSignature string, numThreads int, thread uint, tries uint) {
	defer wg.Done()
	maxZeroes := 0
	startNum := thread * tries
	maxNum := (thread + 1) * tries

	for i := startNum; i < maxNum; i++ {
		select {
		case <-goldenFound:
			return
		default:
			break
		}

		newFuncSig := fmt.Sprintf(funcSignature, i)
		funcSelector := getFuncSelector(newFuncSig)
		numZeroes := countLeadingZeros(funcSelector)

		if numZeroes%2 == 0 && numZeroes > maxZeroes {
			maxZeroes = numZeroes
			ch <- funcSelector
			if funcSelector[0:6] == "000000" {
				fmt.Println(fmt.Sprintf("Found golden function selector '%s' with value '%s'", funcSelector, newFuncSig))
				for t := 0; t < numThreads; t++ {
					goldenFound <- true
				}
				return
			}
		}
	}
}

func receiveFuncSelectors(ch chan string) {
	open := true
	for open {
		_, open = <-ch
	}
}

func countLeadingZeros(funcSelector string) int {
	count := 0
	for _, a := range funcSelector {
		if a == '0' {
			count++
		} else {
			break
		}
	}
	return count
}

func getFuncSelector(funcSignature string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(funcSignature))
	funcSelBytes := hash.Sum(nil)
	return hex.EncodeToString(funcSelBytes[:4])
}
