package golf

import (
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/crypto/sha3"
)

const (
	MAX_UINT    = ^uint(0)
	PLACEHOLDER = "@"
)

type Result struct {
	Name      string
	Selector  string
	TimeTaken time.Duration
}

func SearchFuncSelector(funcSignature string) (Result, error) {

	inputError := validateInput(funcSignature)
	if inputError != nil {
		return Result{}, inputError
	}

	start := time.Now()
	numThreads := runtime.NumCPU()

	var wg sync.WaitGroup
	ch := make(chan Result, numThreads)
	goldenFound := make(chan bool, numThreads)
	tries := MAX_UINT / uint(numThreads)

	for thread := 0; thread < numThreads; thread++ {
		wg.Add(1)
		go runSearcher(&wg, ch, goldenFound, funcSignature, numThreads, uint(thread), tries)
	}

	wg.Wait()
	close(ch)

	result := aggregateFuncSelectors(ch)
	result.TimeTaken = time.Since(start)

	return result, nil
}

func runSearcher(wg *sync.WaitGroup, ch chan<- Result, goldenFound chan bool, funcSignature string, numThreads int, thread uint, tries uint) {
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
			min = Result{Name: newFuncSig, Selector: funcSelector}
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

func aggregateFuncSelectors(ch chan Result) Result {
	maxZeroes := 0
	lowestSelector := Result{}
	for result := range ch {
		sel := result.Selector
		if sel != "" {
			zeros := countLeadingZeros(sel)
			if zeros > maxZeroes {
				lowestSelector = result
				maxZeroes = zeros
			}
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

func validateInput(funcSelector string) error {
	// check how many place holder characters, revert if more than 1
	phError, _ := checkPlaceholderCharacters(funcSelector)
	if phError != nil {
		return phError
	}

	// ensure all characters in function name are valid e.g. no symbols
	funcNameError := checkFunctionName(funcSelector)
	if funcNameError != nil {
		return funcNameError
	}

	// ensure all specified input types are valid e.g. uint not valid ... uint256 valid
	funcStructureError := checkFunctionStructure(funcSelector)
	if funcStructureError != nil {
		return funcStructureError
	}
	return nil
}

func checkPlaceholderCharacters(funcSelector string) (error, int) {
	phCount := strings.Count(funcSelector, PLACEHOLDER)
	if phCount > 1 {
		return errors.New("too many placeholders found in input"), 0
	}
	return nil, phCount
}

func checkFunctionName(funcSelector string) error {
	funcName := strings.Split(funcSelector, "(")
	if len(funcName) == 1 {
		return errors.New("no opening bracket found in function name")
	}

	for _, c := range funcName[0] {
		if !unicode.IsLetter(c) || !unicode.IsDigit(c) {
			return errors.New(fmt.Sprintf("character %c not allowed in function name", c))
		}
	}
	return nil
}

func checkFunctionStructure(funcSelector string) error {
	if strings.Count(funcSelector, "(") != 1 || strings.Count(funcSelector, ")") != 1 {
		return errors.New("invalid number of brackets")
	}
	//check valid input types
	return nil
}
