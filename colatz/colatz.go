package main

import (
	"flag"
	"fmt"
)

func collatz(num int) int {
	var wyk int = 0

	for num != 1 {
		wyk++
		fmt.Print(num, " ")
		if num%2 == 0 {
			num = num / 2
		} else {
			num = 3*num + 1
		}
	}

	fmt.Println(num)
	return wyk
}

func collatz_for_range(start_num int, end_num int) {
	var sum int = 0
	var sr int = 0

	for i := start_num; i <= end_num; i++ {
		sum += collatz(i)
		if i%1000 == 0 {
			sr = sum / i
			fmt.Println(i, sr)
		}
	}
}

func main() {

	var argN int = 10

	flag.IntVar(&argN, "N", argN, "Wartosc liczby max dla ktorej ma byc obliczona sekwencja Collatza")
	flag.Parse()

	// arr := [1000]int{}
	collatz_for_range(1, argN)
	// var N = strconv.ParseInt(argN, 32)

	// robKolacje(N)

}

// funkcja wypisuje wszystkie elementy Ciągu Collatza od 1 az do N

func robKolacje(N int64) {
	// for i := 1; i <= N; i++ {
	// 	fmt.Println(i, collatz(i))
	// }

	for N > 1 {
		if N%2 == 0 {
			N = N / 2
		} else {
			N = 3*N + 1
		}
		fmt.Printf("%d ", N)
	}
	fmt.Println("1")
}
