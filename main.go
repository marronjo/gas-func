package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"sync"
	"time"
)

func main() {
	funcSignature := "transfer%d(address,uint256)"
	tries := 100000000

	startSeq := time.Now()
	tryBruteForce(funcSignature, tries)
	fmt.Println("Sequential brute force done : %s", time.Since(startSeq))

	startThreaded := time.Now()
	var wg sync.WaitGroup
	tryBruteForceThreaded(&wg, funcSignature, tries)
	wg.Wait()
	fmt.Println("Threaded brute force done : %s", time.Since(startThreaded))
}

func tryBruteForceThreaded(wg *sync.WaitGroup, funcPattern string, tries int) {
	for i := 0; i < tries; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			checkZeros(funcPattern, index)
		}(i)
	}
}

func tryBruteForce(funcPattern string, tries int) {
	for i := 0; i < tries; i++ {
		checkZeros(funcPattern, i)
	}
}

func checkZeros(funcPattern string, index int) {
	newFuncSig := fmt.Sprintf(funcPattern, index)
	funcSelector := getFuncSelector(newFuncSig)
	if funcSelector[0:4] == "0000" {
		println(newFuncSig + " " + funcSelector)
	}
}

func getFuncSelector(funcSignature string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(funcSignature))
	funcSelBytes := hash.Sum(nil)
	return hex.EncodeToString(funcSelBytes[:4])
}
