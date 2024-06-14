package main

import (
	"fmt"
)

func peselIsValid(pesel []int) {
	if len(pesel) == 11 {
		var sum int = 0
		var weights = [10]int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}
		for i := 0; i < 10; i++ {
			sum += weights[i] * pesel[i]
		}
		controlDigit := 10 - (sum % 10)
		if controlDigit == 10 {
			controlDigit = 0
		}
		if controlDigit == pesel[10] {
			fmt.Println("PESEL is valid")
		} else {
			fmt.Println("PESEL is not valid")
		}
	}
}

func main() {

	var pesel1 = []int{1, 9, 1, 0, 2, 4, 1, 2, 3, 2, 0}
	var pesel2 = []int{8, 8, 6, 6, 2, 4, 5, 3, 2, 7, 2}
	var pesel3 = []int{9, 7, 7, 2, 1, 0, 0, 3, 1, 2, 1}

	peselIsValid(pesel1)
	peselIsValid(pesel2)
	peselIsValid(pesel3)
}
