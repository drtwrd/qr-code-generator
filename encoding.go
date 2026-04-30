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
	InvalidMode
)

func (m EncodingMode) Indicator() string {
	switch m {
	case NumericMode:
		return "0001"
	case AlphanumericMode:
		return "0010"
	case ByteMode:
		return "0100"
	default:
		return "0000"
	}

}

type ErrorCorrectionLevel string

const (
	LevelL ErrorCorrectionLevel = "L"
	LevelM ErrorCorrectionLevel = "M"
	LevelQ ErrorCorrectionLevel = "Q"
	LevelH ErrorCorrectionLevel = "H"
)

const QrVersion = 1

func getModeIndicator(mode EncodingMode) string {
	return mode.Indicator()
}

func getCharCountIndicator(mode EncodingMode, version int, textLength int) string {
	// For alphanumeric mode and verison 1 needs 9 bits long
	return fmt.Sprintf("%09b", textLength)
}

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
	test := "HELLO WORLD"
	fmt.Println(getModeIndicator(AlphanumericMode))
	fmt.Println(getCharCountIndicator(AlphanumericMode, QrVersion, len(test)))
}
