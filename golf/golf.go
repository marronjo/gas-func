package golf

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"
)

type Result struct {
	name     string
	selector string
}

func SearchFuncSelector(funcSignature string, numThreads int) (string, string, time.Duration) {
	start := time.Now()
	var sender sync.WaitGroup
	ch := make(chan Result, numThreads)
	goldenFound := make(chan bool, numThreads)
	tries := ^uint(0) / uint(numThreads)

	for thread := 0; thread < numThreads; thread++ {
		sender.Add(1)
		go runRoutinesFastest(&sender, ch, goldenFound, funcSignature, numThreads, uint(thread), tries)
	}

	sender.Wait()
	close(ch)

	result := aggregateFuncSelectors(ch, numThreads)

	return result.name, result.selector, time.Since(start)
}

func runRoutinesFastest(wg *sync.WaitGroup, ch chan<- Result, goldenFound chan bool, funcSignature string, numThreads int, thread uint, tries uint) {
	defer wg.Done()
	maxZeroes := 0
	min := Result{}
	startNum := thread * tries
	maxNum := (thread + 1) * tries

	for i := startNum; i < maxNum; i++ {
		select {
		case <-goldenFound:
			ch <- min
			return
		default:
		}

		newFuncSig := fmt.Sprintf(funcSignature, i)
		funcSelector := getFuncSelector(newFuncSig)
		numZeroes := countLeadingZeros(funcSelector)

		if numZeroes%2 == 0 && numZeroes > maxZeroes {
			maxZeroes = numZeroes
			min = Result{name: newFuncSig, selector: funcSelector}
			if funcSelector[0:6] == "000000" {
				for t := 0; t < numThreads; t++ {
					goldenFound <- true
				}
				ch <- min
				return
			}
		}
	}
	ch <- min
}

func aggregateFuncSelectors(ch chan Result, numThreads int) Result { //, receiver *sync.WaitGroup) {
	selectors := make([]Result, numThreads)
	for result := range ch {
		if result.selector != "" {
			selectors = append(selectors, result)
		}
	}
	return getLowestSelector(selectors)
}

func getLowestSelector(selectors []Result) Result {
	maxZeroes := 0
	lowestSelector := Result{}
	for i := 0; i < len(selectors); i++ {
		zeros := countLeadingZeros(selectors[i].selector)
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
