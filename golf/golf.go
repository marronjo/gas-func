package golf

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"
)

func SearchFuncSelector(funcSignature string, numThreads int) (string, time.Duration) {
	start := time.Now()
	var sender sync.WaitGroup
	ch := make(chan string, numThreads)
	goldenFound := make(chan bool, numThreads)
	tries := ^uint(0) / uint(numThreads)

	for thread := 0; thread < numThreads; thread++ {
		sender.Add(1)
		go runRoutinesFastest(&sender, ch, goldenFound, funcSignature, numThreads, uint(thread), tries)
	}

	sender.Wait()
	close(ch)

	bestFuncSelector := aggregateFuncSelectors(ch, numThreads)

	fmt.Printf("Best Selector : %s\n", bestFuncSelector)
	fmt.Printf("Complete fastest run in %s", time.Since(start))
	return bestFuncSelector, time.Since(start)
}

func runRoutinesFastest(wg *sync.WaitGroup, ch chan<- string, goldenFound chan bool, funcSignature string, numThreads int, thread uint, tries uint) {
	defer wg.Done()
	maxZeroes := 0
	minFuncSelector := ""
	startNum := thread * tries
	maxNum := (thread + 1) * tries

	for i := startNum; i < maxNum; i++ {
		select {
		case <-goldenFound:
			ch <- minFuncSelector
			return
		default:
		}

		newFuncSig := fmt.Sprintf(funcSignature, i)
		funcSelector := getFuncSelector(newFuncSig)
		numZeroes := countLeadingZeros(funcSelector)

		if numZeroes%2 == 0 && numZeroes > maxZeroes {
			maxZeroes = numZeroes
			minFuncSelector = funcSelector
			if funcSelector[0:6] == "000000" {
				fmt.Printf("Found golden function selector '%s' with value '%s'\n", funcSelector, newFuncSig)
				for t := 0; t < numThreads; t++ {
					goldenFound <- true
				}
				ch <- funcSelector
				return
			}
		}
	}
	ch <- minFuncSelector
}

func aggregateFuncSelectors(ch chan string, numThreads int) string { //, receiver *sync.WaitGroup) {
	selectors := make([]string, numThreads)
	for val := range ch {
		if val != "" {
			selectors = append(selectors, val)
		}
	}
	return getLowestSelector(selectors)
}

func getLowestSelector(selectors []string) string {
	maxZeroes := 0
	lowestSelector := ""
	for i := 0; i < len(selectors); i++ {
		zeros := countLeadingZeros(selectors[i])
		if zeros > maxZeroes {
			lowestSelector = selectors[i]
			maxZeroes = zeros
		}
	}
	return lowestSelector
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
