package main

import "fmt"

func createEmptyMatrix(size int) [][]int {
	matrix := make([][]int, size)

	for i := range matrix {
		matrix[i] = make([]int, size)
	}

	return matrix
}

func printMatrix(matrix [][]int) {
	for _, row := range matrix {
		for _, cell := range row {
			fmt.Printf("%2d", cell)
		}
		fmt.Println()
	}
}
