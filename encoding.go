package main

import (
	"fmt"
	"strings"
)

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

func isAlphanumeric(text string) bool {
	const specialsChars = "$%*+-./: "

	for _, char := range text {
		if !(char >= '0' && char <= '9') &&
			!(char >= 'A' && char <= 'Z') &&
			!strings.ContainsRune(specialsChars, char) {
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

	if isAlphanumeric(textToEncode) {
		return AlphanumericMode
	}

	return InvalidMode
}

func main() {
	fmt.Println(isAlphanumeric("HELLO123"))
	fmt.Println(isAlphanumeric("Hello123"))
	fmt.Println(isAlphanumeric("HELLO$%*"))
	fmt.Println(isAlphanumeric("HELLO WORLD"))
	fmt.Println(isAlphanumeric("HELLO@WORLD"))
}
