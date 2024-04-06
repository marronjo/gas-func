package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
)

func main() {
	funcSignature := "transfer%d(address,uint256)"
	tryBruteForce(funcSignature, 100000)
}

func tryBruteForce(funcPattern string, tries int) {
	for i := 0; i < tries; i++ {
		newFuncSig := fmt.Sprintf(funcPattern, i)
		funcSelector := getFuncSelector(newFuncSig)
		if funcSelector[0:4] == "0000" {
			println(newFuncSig + " " + funcSelector)
		}
	}
}

func getFuncSelector(funcSignature string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(funcSignature))
	funcSelBytes := hash.Sum(nil)
	return hex.EncodeToString(funcSelBytes[:4])
}
