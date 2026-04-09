package domain

import (
	"strconv"
	"strings"
	"unicode"
)

func TransformText(input string) string {
	tokens := Tokenize(input)
	processed := applyCommands(tokens)

	var builder strings.Builder
	for _, token := range processed {
		builder.WriteString(token.Value)
	}

	return builder.String()
}

func applyCommands(tokens []Token) []Token {
	result := make([]Token, 0, len(tokens))

	for _, token := range tokens {
		if !token.IsCommand() {
			result = append(result, token)
			continue
		}

		name, count, ok := parseCommand(token.Value)
		if !ok {
			result = append(result, token)
			continue
		}

		wordIndexes := findPreviousWordIndexes(result, count)
		applyWordTransformation(result, wordIndexes, name)
	}

	return result
}

func parseCommand(raw string) (string, int, bool) {
	if len(raw) < 3 || raw[0] != '(' || raw[len(raw)-1] != ')' {
		return "", 0, false
	}

	inner := raw[1 : len(raw)-1]
	parts := strings.Split(inner, ",")

	if len(parts) == 1 {
		switch parts[0] {
		case "up", "low", "cap", "hex", "bin":
			return parts[0], 1, true
		default:
			return "", 0, false
		}
	}

	if len(parts) != 2 {
		return "", 0, false
	}

	name := parts[0]
	if name != "up" && name != "low" && name != "cap" {
		return "", 0, false
	}

	countValue := strings.TrimSpace(parts[1])
	count, err := strconv.Atoi(countValue)
	if err != nil || count <= 0 {
		return "", 0, false
	}

	return name, count, true
}

func findPreviousWordIndexes(tokens []Token, count int) []int {
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
