package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func silnia(nickBytes []byte) (*big.Int, int64) {
	factorial := big.NewInt(1)
	i := int64(1)

	for {
		factorialStr := factorial.String()
		allPresent := true
		for _, b := range nickBytes {
			if !strings.Contains(factorialStr, strconv.Itoa(int(b))) {
				allPresent = false
				break
			}
		}

		if allPresent {
			return factorial, i
		}

		i++
		bigI := big.NewInt(i)
		factorial.Mul(factorial, bigI)
	}
}

var callCount = make(map[int]int)
var result = make(map[int]int)

func fibRecursive(n int) int {
	callCount[n]++

	if n <= 1 {
		return callCount[n]
	}
	return fibRecursive(n-1) + fibRecursive(n-2)
}

func turnPositive(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	nick := "MicWit"
	nickSmall := strings.ToLower(nick)
	nickBytes := []byte(nickSmall)

	//silna liczba
	_, strongestNum := silnia(nickBytes)
	fmt.Println("Twoją silną liczbą jest:", strongestNum)

	//słaba liczba
	fibRecursive(30)
	for i := 1; i <= 30; i++ {
		result[i] = callCount[i]
		fmt.Printf("Number: %d, Calls: %d\n", i, callCount[i])
	}

	weakestNum := 0
	minDiff := math.MaxInt64
	for num, calls := range result {
		diff := turnPositive(strongestNum - int64(calls))
		if diff < int64(minDiff) {
			minDiff = int(diff)
			weakestNum = num
		}
	}

	fmt.Printf("Nasza słaba liczba to: %d\n", weakestNum)
}
