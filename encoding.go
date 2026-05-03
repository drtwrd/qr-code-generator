package main

import (
	"errors"
	"fmt"
	"strconv"
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

func getModeIndicator(mode EncodingMode) string {
	return mode.Indicator()
}

func getCharCountIndicator(mode EncodingMode, version int, textLength int) string {
	bits := 0

	if version >= 1 && version <= 9 {
		switch mode {
		case NumericMode:
			bits = 10
		case AlphanumericMode:
			bits = 9
		case ByteMode:
			bits = 8
		}
	} else if version >= 10 && version <= 26 {
		switch mode {
		case NumericMode:
			bits = 12
		case AlphanumericMode:
			bits = 11
		case ByteMode:
			bits = 16
		}
	} else if version >= 27 && version <= 40 {
		switch mode {
		case NumericMode:
			bits = 14
		case AlphanumericMode:
			bits = 13
		case ByteMode:
			bits = 16
		}
	} else {
		return ""
	}

	format := "%0" + strconv.Itoa(bits) + "b"
	return fmt.Sprintf(format, textLength)
}

func encodeInNumericMode(text string) string {
	var encodedText strings.Builder

	for i := 0; i < len(text); i += 3 {
		group := text[i:min(i+3, len(text))]
		length := len(group)

		num, _ := strconv.Atoi(group)

		switch length {
		case 3:
			fmt.Fprintf(&encodedText, "%010b", num)
		case 2:
			fmt.Fprintf(&encodedText, "%07b", num)
		case 1:
			fmt.Fprintf(&encodedText, "%04b", num)
		}
	}

	return encodedText.String()
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

func encodeInByteMode(text string) string {
	var encodedText strings.Builder

	for _, char := range text {
		charToBinary := fmt.Sprintf("%08b", char)
		fmt.Fprintf(&encodedText, charToBinary)
	}

	return encodedText.String()
}

func addPadBytes(binaryString string, totalBits int) string {
	var outputBinary strings.Builder
	outputBinary.WriteString(binaryString)

	// Add a Terminator of 0s if Necessary
	remaining := totalBits - len(binaryString)
	if remaining > 0 {
		terminatorBits := min(remaining, 4)
		outputBinary.WriteString(strings.Repeat("0", terminatorBits))
	}

	// Add More 0s to Make the Length a Multiple of 8
	remainder := len(outputBinary.String()) % 8
	if remainder != 0 {
		zeroesToAdd := 8 - remainder
		outputBinary.WriteString(strings.Repeat("0", zeroesToAdd))
	}

	// Add Pad Bytes if the String is Still too Short
	for len(outputBinary.String()) < totalBits {
		outputBinary.WriteString("11101100") // 236
		if len(outputBinary.String()) < totalBits {
			outputBinary.WriteString("00010001") // 17
		}
	}

	return outputBinary.String()
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

func determineSmallestVersion(mode EncodingMode, ecLevel ErrorCorrectionLevel, textLength int) (int, error) {
	for version := 1; version <= 40; version++ {
		maxCapacity := QrCapacityByVersionErrLevelMode[version][ecLevel][mode]
		if textLength <= maxCapacity {
			return version, nil
		}
	}
	return 0, errors.New("data too long for QR code version 40")
}
