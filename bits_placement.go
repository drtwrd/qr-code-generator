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
			if cell == 1 {
				fmt.Print("██")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
}

func addFinderPatterns(matrix [][]int) {
	size := len(matrix)

	// Top left square
	for i := range 7 {
		for j := range 7 {
			if i == 0 || i == 6 || j == 0 || j == 6 {
				matrix[i][j] = 1
			}

			if (i >= 2 && i <= 4) && (j >= 2 && j <= 4) {
				matrix[i][j] = 1
			}
		}
	}

	// Top right square
	for i := range 7 {
		for j := size - 7; j < size; j++ {
			if i == 0 || i == 6 || j == size-7 || j == size-1 {
				matrix[i][j] = 1
			}

			if (i >= 2 && i <= 4) && (j >= size-5 && j <= size-3) {
				matrix[i][j] = 1
			}
		}
	}

	// Bottom left square
	for i := size - 7; i < size; i++ {
		for j := range 7 {
			if i == size-7 || i == size-1 || j == 0 || j == 6 {
				matrix[i][j] = 1
			}

			if (i >= size-5 && i <= size-3) && (j >= 2 && j <= 4) {
				matrix[i][j] = 1
			}
		}
	}
}

// For version 1 (21x21), alignment patterns are not required.
// TODO: implement for versions 2 and higher
func addAlignmentPatterns(matrix [][]int) {}

func addTimingPatterns(matrix [][]int) {
	size := len(matrix)

	// Horizontal timing pattern
	for i := 8; i < size-8; i++ {
		if i%2 == 0 {
			matrix[6][i] = 1
		} else {
			matrix[6][i] = 0
		}
	}

	// Vertical timing pattern
	for j := 8; j < size-8; j++ {
		if j%2 == 0 {
			matrix[j][6] = 1
		} else {
			matrix[j][6] = 0
		}
	}
}

func main() {
	matrix := createEmptyMatrix(21)

	addFinderPatterns(matrix)
	addTimingPatterns(matrix)

	printMatrix(matrix)
}
