package main

import "fmt"

func multiply_matrices(matrix1 [][]int, matrix2 [][]int) [][]int {
	var result [][]int

	if len(matrix1[0]) != len(matrix2) {
		panic("Matrices are not the same size")
	}

	for k := 0; k < len(matrix1); k++ {
		var temp []int
		for i := 0; i < len(matrix1[0]); i++ {
			var res int = 0
			for j := 0; j < len(matrix2); j++ {
				res += matrix1[k][j] * matrix2[j][i]
			}
			temp = append(temp, res)
		}
		result = append(result, temp)
	}

	return result
}

func main() {
	matrix1 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matrix2 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	result := multiply_matrices(matrix1, matrix2)

	fmt.Println("matrix1 = ", matrix1)
	fmt.Println("matrix2 = ", matrix2)
	fmt.Println("matrix1 * matrix2 = ", result)
}
