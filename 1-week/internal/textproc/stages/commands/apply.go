package commands

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/model"
)

func Apply(tokens []model.Token) []model.Token {
	result := make([]model.Token, 0, len(tokens))

	for _, token := range tokens {
		if !token.IsCommand() {
			result = append(result, token)
			continue
		}

		spec, ok := ParseSpec(token.Value)
		if !ok {
			result = append(result, token)
			continue
		}

		wordIndexes := findPreviousWordIndexes(result, spec.Count)
		applyWordTransformation(result, wordIndexes, spec.Name)
	}

	return result
}

func findPreviousWordIndexes(tokens []model.Token, count int) []int {
	indexes := make([]int, 0, count)

	for i := len(tokens) - 1; i >= 0 && len(indexes) < count; i-- {
		if tokens[i].IsWord() {
			indexes = append(indexes, i)
		}
	}

	for left, right := 0, len(indexes)-1; left < right; left, right = left+1, right-1 {
		indexes[left], indexes[right] = indexes[right], indexes[left]
	}

	return indexes
}

func applyWordTransformation(tokens []model.Token, wordIndexes []int, commandName string) {
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
