package commands

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/grammar"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/stages"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/token"
)

func Apply(tokens []token.Token) []token.Token {
	result := make([]token.Token, 0, len(tokens))

	for _, tok := range tokens {
		if !tok.IsCommand() {
			result = append(result, tok)
			continue
		}

		spec, ok := grammar.ParseCommandSpec(tok.Value)
		if !ok {
			result = append(result, tok)
			continue
		}

		wordIndexes := findPreviousWordIndexes(result, spec.Count)
		applyWordTransformation(result, wordIndexes, spec.Name)
	}

	return result
}

func findPreviousWordIndexes(tokens []token.Token, count int) []int {
	indexes := make([]int, 0, count)
	searchStart := len(tokens) - 1

	for len(indexes) < count {
		index, ok := stages.FindPrevIndex(tokens, searchStart, func(tok token.Token) bool {
			return tok.IsWord()
		})
		if !ok {
			break
		}

		indexes = append(indexes, index)
		searchStart = index - 1
	}

	for left, right := 0, len(indexes)-1; left < right; left, right = left+1, right-1 {
		indexes[left], indexes[right] = indexes[right], indexes[left]
	}

	return indexes
}

func applyWordTransformation(tokens []token.Token, wordIndexes []int, commandName string) {
	for _, index := range wordIndexes {
		tokens[index].Value = transformWord(tokens[index].Value, commandName)
	}
}

func transformWord(word, commandName string) string {
	switch commandName {
	case grammar.CommandUp:
		return strings.ToUpper(word)
	case grammar.CommandLow:
		return strings.ToLower(word)
	case grammar.CommandCap:
		return capitalizeWord(word)
	case grammar.CommandHex:
		return convertBaseToDecimal(word, 16)
	case grammar.CommandBin:
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
