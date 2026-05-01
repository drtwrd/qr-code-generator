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

func encodeInAlphanumericMode(text string) string {
	alphanumericTable := map[rune]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'A': 10,
		'B': 11,
		'C': 12,
		'D': 13,
		'E': 14,
		'F': 15,
		'G': 16,
		'H': 17,
		'I': 18,
		'J': 19,
		'K': 20,
		'L': 21,
		'M': 22,
		'N': 23,
		'O': 24,
		'P': 25,
		'Q': 26,
		'R': 27,
		'S': 28,
		'T': 29,
		'U': 30,
		'V': 31,
		'W': 32,
		'X': 33,
		'Y': 34,
		'Z': 35,
		' ': 36,
		'$': 37,
		'%': 38,
		'*': 39,
		'+': 40,
		'-': 41,
		'.': 42,
		'/': 43,
		':': 44,
	}

	var encodedText strings.Builder

	runes := []rune(text)

	for i := 0; i < len(runes)-1; i += 2 {
		firstChar := alphanumericTable[runes[i]]
		secondChar := alphanumericTable[runes[i+1]]

		numberToEncode := (45 * firstChar) + secondChar

		fmt.Fprintf(&encodedText, "%011b", numberToEncode)
	}

	if len(runes)%2 != 0 {
		lastChar := alphanumericTable[runes[len(runes)-1]]
		fmt.Fprintf(&encodedText, "%06b", lastChar)
	}

	return encodedText.String()
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
	modeIndicator := getModeIndicator(AlphanumericMode)
	charCountIndicator := getCharCountIndicator(AlphanumericMode, QrVersion, 11)
	textEncodedInAlphanumericMode := encodeInAlphanumericMode(test)

	encocedText := modeIndicator + charCountIndicator + textEncodedInAlphanumericMode

	fmt.Println(encocedText)
}
