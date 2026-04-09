package domain

import (
	"strconv"
	"strings"
	"unicode"
)

func applyWordTransformation(tokens []Token, wordIndexes []int, commandName string) {
	for _, index := range wordIndexes {
		tokens[index].Value = transformWord(tokens[index].Value, commandName)
	}
}

func transformWord(word, commandName string) string {
	switch commandName {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		return capitalizeWord(word)
	case "hex":
		return convertBaseToDecimal(word, 16)
	case "bin":
		return convertBaseToDecimal(word, 2)
	default:
		return word
	}
}

func capitalizeWord(word string) string {
	if word == "" {
		return word
	}

	runes := []rune(strings.ToLower(word))
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func convertBaseToDecimal(word string, base int) string {
	value, err := strconv.ParseInt(word, base, 64)
	if err != nil {
		return word
	}

	return strconv.FormatInt(value, 10)
}
