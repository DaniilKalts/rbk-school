package domain

import "strings"

func joinTokens(tokens []Token) string {
	var builder strings.Builder

	for _, token := range tokens {
		builder.WriteString(token.Value)
	}

	return builder.String()
}

func nextNonSpaceTokenIndex(tokens []Token, start int) int {
	for i := start; i < len(tokens); i++ {
		if !tokens[i].IsSpace() {
			return i
		}
	}

	return -1
}

func appendSpaceToken(tokens *[]Token) {
	if len(*tokens) == 0 {
		return
	}

	if (*tokens)[len(*tokens)-1].IsSpace() {
		(*tokens)[len(*tokens)-1].Value = " "
		return
	}

	*tokens = append(*tokens, Token{Kind: KindSpace, Value: " "})
}

func trimTrailingSpaces(tokens *[]Token) {
	for len(*tokens) > 0 {
		lastIndex := len(*tokens) - 1
		if !(*tokens)[lastIndex].IsSpace() {
			break
		}

		*tokens = (*tokens)[:lastIndex]
	}
}

func lastNonSpaceToken(tokens []Token) (Token, bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		if !tokens[i].IsSpace() {
			return tokens[i], true
		}
	}

	return Token{}, false
}
