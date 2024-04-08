package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"time"
)

func main() {
	funcSignature := "placeOrder%d(uint24,address,uint)"
	tries := 10000000

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
	return hex.EncodeToString(hash.Sum(nil)[:4])
}
