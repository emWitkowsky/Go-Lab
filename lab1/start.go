package main

import (
	"fmt"
	"math"
)

func isTriangle(a, b, c int) bool {
	if a+b > c && b+c > a && a+c > b {
		return true
	}
	return false
}

func check() string {
	if isTriangle(1, 2, 1) {
		return "Yes"
	} else {
		return "No"
	}
}

// program który liczy pierwiastki trójmianu kwadratowego

func delta(a, b, c float64) float64 {
	return b*b - 4*a*c
}

func particle(a, b, c float64) (float64, float64) {
	var x1, x2 float64

	del := delta(a, b, c)

	x1 = (-b + math.Sqrt(del)) / (2 * a)
	x2 = (-b - math.Sqrt(del)) / (2 * a)

	return x1, x2
}

func main() {
	fmt.Println("Hello, World!")

	var x int = 1

	fmt.Println(x)

	fmt.Println(check())

	given1, given2 := particle(1, -3, 2)
	fmt.Printf("x1 = %f, x2 = %f\n", given1, given2)

}
