package articles

import (
	"strings"
	"unicode/utf8"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/model"
)

func Apply(tokens []model.Token) []model.Token {
	for i := range tokens {
		if !tokens[i].IsWord() || !isArticleA(tokens[i].Value) {
			continue
		}

		nextWord, ok := nextWordToken(tokens, i+1)
		if !ok {
			continue
		}

		if startsWithVowelOrH(nextWord.Value) {
			tokens[i].Value = toAn(tokens[i].Value)
		}
	}

	return tokens
}

func isArticleA(word string) bool {
	return word == "a" || word == "A"
}

func nextWordToken(tokens []model.Token, start int) (model.Token, bool) {
	for i := start; i < len(tokens); i++ {
		if tokens[i].IsWord() {
			return tokens[i], true
		}
	}

	return model.Token{}, false
}

func startsWithVowelOrH(word string) bool {
	if word == "" {
		return false
	}

	r, _ := utf8.DecodeRuneInString(strings.ToLower(word))

	switch r {
	case 'a', 'e', 'i', 'o', 'u', 'h':
		return true
	default:
		return false
	}
}

func toAn(article string) string {
	if article == "A" {
		return "An"
	}

	return "an"
}
