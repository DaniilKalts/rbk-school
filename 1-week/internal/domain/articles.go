package domain

import "strings"

func applyArticles(tokens []Token) []Token {
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

func nextWordToken(tokens []Token, start int) (Token, bool) {
	for i := start; i < len(tokens); i++ {
		if tokens[i].IsWord() {
			return tokens[i], true
		}
	}

	return Token{}, false
}

func startsWithVowelOrH(word string) bool {
	if word == "" {
		return false
	}

	firstRune := []rune(strings.ToLower(word))[0]

	switch firstRune {
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
