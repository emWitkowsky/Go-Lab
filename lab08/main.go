package main

import (
	"fmt"
	"math"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Point[T Number] struct {
	X T
	Y T
}

func main() {
	// Simple example with int
	point_int := Point[int]{X: 3, Y: 4}
	fmt.Println(distance(point_int))

	// Simple example with float64
	point_float := Point[float64]{X: 3.0, Y: 4.0}
	fmt.Println(distance(point_float))

	// Simple example with int
	point_area_int := Point[int]{X: 3, Y: 4}
	fmt.Println(area(point_area_int))

	// Simple example with float64
	point_area_float := Point[float64]{X: 3.5, Y: 4.0}
	fmt.Println(area(point_area_float))

}

// NAPISZ ELEGANCKa STRUKTURĘ OBSŁUGUJąCą PUNKT 2D
// tak aby działał na wszystkich typach numerycznych

// Napisz fajną funkcję Distance(), która oblicza odległość punktu od punktu (0,0)
// wynikiem powinien być ten sam typ co punkt (np. dla int to int, dla float to float)

func distance[T Number](ourPoint Point[T]) T {
	pointZero := Point[T]{0, 0}

	var xDist = ourPoint.X - pointZero.X
	var yDist = ourPoint.Y - pointZero.Y

	if xDist < 0 {
		xDist = -xDist
	}

	if yDist < 0 {
		yDist = -yDist
	}

	// Niestety Sqrt kowertuje wszystko do float64
	finalDist := math.Sqrt(float64(xDist*xDist + yDist*yDist))

	return T(finalDist)
}

func area[T Number](ourPoint Point[T]) T {
	pointZero := Point[T]{0, 0}

	var xDist = ourPoint.X - pointZero.X
	var yDist = ourPoint.Y - pointZero.Y

	if xDist < 0 {
		xDist = -xDist
	}

	if yDist < 0 {
		yDist = -yDist
	}

	finalArea := T(xDist * yDist)

	return T(finalArea)
}
