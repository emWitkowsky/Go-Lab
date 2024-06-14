package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("liczby.txt")
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(data))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		digitsCount := len(strings.ReplaceAll(line, " ", "")) - 1
		fmt.Printf("%d %d\n", lineNumber, digitsCount)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
