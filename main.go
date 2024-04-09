package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"sync"
	"time"
)

func main() {
	funcSignature := "transfer%d(uint256,address)"
	//tries := 10000000
	//searchFuncSelectorBenchmark(funcSignature, tries)
	//searchFuncSelectorFast(funcSignature)
	searchFuncSelectorFastest(funcSignature)
}

// concurrent search, split search amount 10 goroutines and stop as soon as a minimum value function selector is found
func searchFuncSelectorFastest(funcSignature string) {
	start := time.Now()
	var sender sync.WaitGroup
	ch := make(chan string, 100)
	goldenFound := make(chan bool, 10)
	goldenFound <- false
	for i := 0; i < 10; i++ {
		sender.Add(1)
		go runRoutinesFastest(&sender, ch, goldenFound, funcSignature, i)
	}
	go receiveFuncSelectors(ch)
	sender.Wait()
	close(ch)
	fmt.Println(fmt.Sprintf("Complete fast run in %s", time.Since(start)))
}

func runRoutinesFastest(wg *sync.WaitGroup, ch chan<- string, goldenFound chan bool, funcSignature string, tries int) {
	defer wg.Done()
	maxZeroes := 0
	minFuncSelector := ""
	startNum := tries * 10000000
	maxNum := (tries + 1) * 10000000
	fmt.Println("Start -> Max : ", startNum, maxNum)
	for i := startNum; i < maxNum; i++ {
		select {
		case <-goldenFound:
			fmt.Println("Golden Set (Exiting)")
			return
		default:
			break
		}
		newFuncSig := fmt.Sprintf(funcSignature, i)
		funcSelector := getFuncSelector(newFuncSig)

		numZeroes := countLeadingZeros(funcSelector)
		if numZeroes%2 == 0 && numZeroes > maxZeroes {
			maxZeroes = numZeroes
			minFuncSelector = funcSelector
			ch <- minFuncSelector
			if funcSelector[0:6] == "000000" {
				fmt.Println("Golden Found : ", funcSelector)
				for a := 0; a < 10; a++ {
					goldenFound <- true
				}
				return
			}
		}
	}
}

// concurrent search, split search amount 10 goroutines print the minimum value function selectors from each
func searchFuncSelectorFast(funcSignature string) {
	start := time.Now()
	var sender sync.WaitGroup
	ch := make(chan string, 100)
	for i := 0; i < 10; i++ {
		sender.Add(1)
		go runRoutines(&sender, ch, funcSignature, i)
	}
	go receiveFuncSelectors(ch)
	sender.Wait()
	fmt.Println(fmt.Sprintf("Complete fast run in %s", time.Since(start)))
}

func receiveFuncSelectors(ch chan string) {
	open := true
	output := ""
	for open {
		output, open = <-ch
		fmt.Println("Received from channel : "+output+" with open status : ", open)
	}
}

func runRoutines(wg *sync.WaitGroup, ch chan<- string, funcSignature string, tries int) {
	defer wg.Done()
	maxZeroes := 0
	minFuncSelector := ""
	startNum := tries * 1000000
	maxNum := (tries + 1) * 1000000
	fmt.Println("Start -> Max : ", startNum, maxNum)
	for i := startNum; i < maxNum; i++ {
		newFuncSig := fmt.Sprintf(funcSignature, i)
		funcSelector := getFuncSelector(newFuncSig)

		numZeroes := countLeadingZeros(funcSelector)
		if numZeroes%2 == 0 && numZeroes > maxZeroes {
			maxZeroes = numZeroes
			minFuncSelector = funcSelector
			ch <- minFuncSelector
		}
	}
}

// sequential search, run x amount of times and print the minimum value function selector
func searchFuncSelectorBenchmark(funcSignature string, tries int) {
	maxZeroes := 0
	minFuncSelector := ""
	minFuncName := ""
	start := time.Now()
	for i := 0; i < tries; i++ {
		newFuncSig := fmt.Sprintf(funcSignature, i)
		funcSelector := getFuncSelector(newFuncSig)

		numZeroes := countLeadingZeros(funcSelector)
		if numZeroes%2 == 0 && numZeroes > maxZeroes {
			maxZeroes = numZeroes
			minFuncSelector = funcSelector
			minFuncName = newFuncSig
		}
	}
	fmt.Println("Min Selector\t: ", minFuncSelector)
	fmt.Println("Min Func Name \t: ", minFuncName)
	fmt.Println(fmt.Sprintf("Complete in %s with %d tries", time.Since(start), tries))
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
