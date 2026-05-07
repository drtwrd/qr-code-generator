package main

import (
	"errors"
	"fmt"
	"strings"
)

func createEmptyMatrix(size int) [][]int {
	matrix := make([][]int, size)

	for i := range matrix {
		matrix[i] = make([]int, size)
	}

	return matrix
}

func printZeroesAndOnes(matrix [][]int) {
	for _, row := range matrix {
		for _, cell := range row {
			fmt.Printf("%d ", cell)
		}
		fmt.Println()
	}
}

func printMatrix(matrix [][]int) {
	quietZone := 4

	for range quietZone {
		fmt.Printf("%s\n", strings.Repeat("\033[48;5;231m  \033[0m", quietZone*2+len(matrix)))
	}

	for _, row := range matrix {
		fmt.Print(strings.Repeat("\033[48;5;231m  \033[0m", quietZone))
		for _, cell := range row {
			if cell == 1 {
				fmt.Print("\033[48;5;16m  \033[0m")
			} else {
				fmt.Print("\033[48;5;231m  \033[0m")
			}
		}
		fmt.Printf("%s\n", strings.Repeat("\033[48;5;231m  \033[0m", quietZone))
	}

	for range quietZone {
		fmt.Printf("%s\n", strings.Repeat("\033[48;5;231m  \033[0m", quietZone*2+len(matrix)))
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

func addDarkModule(matrix [][]int) {
	size := len(matrix)
	version := getVersion(size)

	matrix[(4*version)+9][8] = 1
}

// TODO: Reserve the Version Information Area
// QR codes versions 7 and larger must contain two areas where version information bits are placed.
// The areas are a 6x3 block above the bottom-left finder pattern and a 3x6 block to the left of the top-right finder pattern.
func isReservedArea(row, col, version, size int) bool {
	// Finder patterns + separators (всегда)
	if (row <= 8 && col <= 8) ||
		(row <= 8 && col >= size-8) ||
		(row >= size-8 && col <= 8) {
		return true
	}

	if version >= 7 {
		// Top-right area (6x3)
		if row <= 5 && col >= size-11 && col <= size-9 {
			return true
		}
		// Bottom-left area (3x6)
		if col <= 5 && row >= size-11 && row <= size-9 {
			return true
		}
	}

	return false
}

func putBits(matrix [][]int, row int, col int, bits []rune, bitIndex *int) {
	size := len(matrix)
	version := getVersion(size)

	if *bitIndex >= len(bits) {
		return
	}

	skipped := 0
	if isReservedArea(row, col, version, size) {
		skipped++
	}
	if col-1 >= 0 && isReservedArea(row, col-1, version, size) {
		skipped++
	}
	if skipped == 2 {
		fmt.Printf("WARNING: both cells in pair (%d,%d) and (%d,%d) are reserved\n", row, col, row, col-1)
	}

	if !isReservedArea(row, col, version, size) && *bitIndex < len(bits) {
		matrix[row][col] = int(bits[*bitIndex] - '0')
		*bitIndex++
	}

	if *bitIndex >= len(bits) {
		return
	}
	if col-1 >= 0 && !isReservedArea(row, col-1, version, size) && *bitIndex < len(bits) {
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

func getMatrixSize(version int) int {
	return ((version - 1) * 4) + 21
}

func GenerateQRCode(text string) ([][]int, error) {
	mode := determineEncodingMode(text)
	if mode == InvalidMode {
		return nil, errors.New("Invalid encoding mode")
	}

	ecLevel := LevelM

	version, err := determineSmallestVersion(mode, ecLevel, len(text))
	if err != nil {
		return nil, err
	}

	size := getMatrixSize(version)
	matrix := createEmptyMatrix(size)

	addFinderPatterns(matrix)
	addTimingPatterns(matrix)
	if version >= 2 {
		addAlignmentPatterns(matrix)
	}
	addDarkModule(matrix)

	placeDataBits(text, matrix)

	return matrix, nil
}
func compareWithExpected(matrix [][]int) {
	expected := [][]int{
		{1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	for i := range 21 {
		for j := range 21 {
			if matrix[i][j] != expected[i][j] {
				fmt.Printf("First mismatch at (%d,%d): got %d, expected %d\n", i, j, matrix[i][j], expected[i][j])
				return
			}
		}
	}
	fmt.Println("Matrix matches expected!")
}
func main() {
	matrix, err := GenerateQRCode("HELLO WORLD")
	if err != nil {
		panic(err)
	}

	printMatrix(matrix)
	// printZeroesAndOnes(matrix)
	compareWithExpected(matrix)
}
