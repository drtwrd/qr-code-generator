package main

import "fmt"

type EncodingMode int

const (
	NumericMode EncodingMode = iota
	AlphanumericMode
	ByteMode
	KanjiMode
	InvalidMode
)

func isNumeric(text string) bool {
	for _, char := range text {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

func determineEncodingMode(textToEncode string) EncodingMode {
	if len(textToEncode) == 0 {
		return InvalidMode
	}

	if isNumeric(textToEncode) {
		return NumericMode
	}

	return InvalidMode
}

func main() {
	test1 := determineEncodingMode("")
	test2 := determineEncodingMode("1029")
	test3 := determineEncodingMode("wasd")

	fmt.Println("Test 1: ", test1)
	fmt.Println("Test 2: ", test2)
	fmt.Println("Test 3: ", test3)
}
