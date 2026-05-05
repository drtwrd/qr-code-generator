package main

import (
	"fmt"
)

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
				fmt.Print("\033[48;5;16m  \033[0m")
			} else {
				fmt.Print("\033[48;5;231m  \033[0m")
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

func canBeDrawn(x, y, size int) bool {
	// Handle top left finder + separator
	if y <= 7 && x <= 7 {
		return false
	}

	// Handle top right finder + separator
	if y <= 7 && x >= size-8 {
		return false
	}

	// Handle bottom left finder + separator
	if y >= size-8 && x <= 7 {
		return false
	}

	return true
}

func addAlignmentPatterns(matrix [][]int) {
	size := len(matrix)
	version := getVersion(size)

	// Alignment patterns are not required for version 1
	if version == 1 {
		return
	}

	squaresToDraw := [][]int{}
	squaresCenters := alignmentCenters[version]
	for _, row := range squaresCenters {
		for _, col := range squaresCenters {
			canDraw := canBeDrawn(row, col, size)
			if canDraw {
				squaresToDraw = append(squaresToDraw, []int{row, col})
			}
		}
	}

	for _, center := range squaresToDraw {
		row, col := center[0], center[1]

		matrix[row][col] = 1

		for i := row - 2; i <= row+2; i++ {
			for j := col - 2; j <= col+2; j++ {
				if i == row-2 || i == row+2 || j == col-2 || j == col+2 {
					matrix[i][j] = 1
				}
			}
		}
	}

}

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

func getVersion(size int) int {
	version := ((size - 21) / 4) + 1

	return version
}

func addDarkSquare(matrix [][]int) {
	size := len(matrix)
	version := getVersion(size)

	matrix[(4*version)+9][8] = 1
}

// TODO: Reserve the Version Information Area
// QR codes versions 7 and larger must contain two areas where version information bits are placed.
// The areas are a 6x3 block above the bottom-left finder pattern and a 3x6 block to the left of the top-right finder pattern.
func isFinderZone(row, col, size int) bool {
	// Top-left
	if row <= 7 && col <= 7 {
		return true
	}

	// Top-right
	if row <= 7 && col >= size-8 {
		return true
	}

	// Bottom-left
	if row >= size-8 && col <= 7 {
		return true
	}

	return false
}

func putBits(matrix [][]int, row int, col int, bits []rune, bitIndex *int) {
	if !isFinderZone(row, col, len(matrix)) && matrix[row][col] == 0 && *bitIndex < len(bits) {
		matrix[row][col] = int(bits[*bitIndex] - '0')
		*bitIndex++
	}

	if !isFinderZone(row, col, len(matrix)) && col-1 >= 0 && matrix[row][col-1] == 0 && *bitIndex < len(bits) {
		matrix[row][col-1] = int(bits[*bitIndex] - '0')
		*bitIndex++
	}
}

func placeDataBits(text string, matrix [][]int) {
	bitString, err := encodeData(text)
	if err != nil {
		panic(err)
	}
	bits := []rune(bitString)
	bitIndex := 0

	size := len(matrix)
	goUp := true

	for col := size - 1; col >= 0; col -= 2 {
		if goUp {
			for row := size - 1; row >= 0; row-- {
				putBits(matrix, row, col, bits, &bitIndex)
			}
		} else {
			for row := range size {
				putBits(matrix, row, col, bits, &bitIndex)
			}
		}
		goUp = !goUp
	}
}

func main() {
	matrix := createEmptyMatrix(21)

	addFinderPatterns(matrix)
	addAlignmentPatterns(matrix)
	addTimingPatterns(matrix)
	addDarkSquare(matrix)

	placeDataBits("IT WORKS", matrix)

	printMatrix(matrix)
}
