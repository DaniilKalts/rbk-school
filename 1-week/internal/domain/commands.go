package domain

import (
	"strconv"
	"strings"
)

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
