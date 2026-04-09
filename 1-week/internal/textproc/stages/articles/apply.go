package articles

import (
	"strings"
	"unicode/utf8"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/stages"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/token"
)

func Apply(tokens []token.Token) []token.Token {
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

func nextWordToken(tokens []token.Token, start int) (token.Token, bool) {
	index, ok := stages.FindNextIndex(tokens, start, func(tok token.Token) bool {
		return tok.IsWord()
	})
	if !ok {
		return token.Token{}, false
	}

	return tokens[index], true
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
