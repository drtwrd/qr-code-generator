package main

import (
	"strings"
)

type EncodingMode int

const (
	NumericMode EncodingMode = iota
	AlphanumericMode
	ByteMode
	InvalidMode
)

type ErrorCorrectionLevel string

const (
	LevelL ErrorCorrectionLevel = "L"
	LevelM ErrorCorrectionLevel = "M"
	LevelQ ErrorCorrectionLevel = "Q"
	LevelH ErrorCorrectionLevel = "H"
)

const QrVersion = 1

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

	return ByteMode
}

func main() {
}
