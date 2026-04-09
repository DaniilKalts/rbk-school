package punctuation

import (
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/model"
)

func Apply(tokens []model.Token) []model.Token {
	result := make([]model.Token, 0, len(tokens))
	quoteOpen := false

	for i, token := range tokens {
		if token.IsSpace() {
			nextIndex := nextNonSpaceTokenIndex(tokens, i+1)
			if nextIndex == -1 || len(result) == 0 {
				continue
			}

			nextToken := tokens[nextIndex]
			prevToken, ok := lastNonSpaceToken(result)
			if !ok {
				continue
			}

			if shouldSkipSpace(prevToken, nextToken, quoteOpen) {
				continue
			}

			if shouldForceSingleSpace(prevToken, nextToken, quoteOpen) {
				appendSpaceToken(&result)
				continue
			}

			appendSpaceToken(&result)
			continue
		}

		if needsSpaceAfterTightPunctuation(result, token, quoteOpen) {
			appendSpaceToken(&result)
		}

		if token.IsPunctuation() {
			if token.Value == "'" {
				if isApostropheInWord(tokens, i) {
					result = append(result, token)
					continue
				}

				if quoteOpen {
					trimTrailingSpaces(&result)
					result = append(result, token)
					quoteOpen = false
					continue
				}

				result = append(result, token)
				quoteOpen = true
				continue
			}

			if isTightPunctuationToken(token) {
				trimTrailingSpaces(&result)
			}
		}

		result = append(result, token)
	}

	return result
}

func isTightPunctuationToken(token model.Token) bool {
	if !token.IsPunctuation() || len(token.Value) != 1 {
		return false
	}

	switch token.Value[0] {
	case '.', ',', '!', '?', ':', ';':
		return true
	default:
		return false
	}
}

func needsSpaceAfterTightPunctuation(result []model.Token, current model.Token, quoteOpen bool) bool {
	prevToken, ok := lastNonSpaceToken(result)
	if !ok || !isTightPunctuationToken(prevToken) {
		return false
	}

	if current.IsPunctuation() {
		if isTightPunctuationToken(current) {
			return false
		}

		if current.Value == "'" && quoteOpen {
			return false
		}
	}

	return true
}

func shouldSkipSpace(prev model.Token, next model.Token, quoteOpen bool) bool {
	if next.IsPunctuation() {
		if isTightPunctuationToken(next) {
			return true
		}

		if next.Value == "'" && quoteOpen {
			return true
		}
	}

	if prev.IsPunctuation() && prev.Value == "'" && quoteOpen {
		return true
	}

	return false
}

func shouldForceSingleSpace(prev model.Token, next model.Token, quoteOpen bool) bool {
	if !isTightPunctuationToken(prev) {
		return false
	}

	if next.IsPunctuation() {
		if isTightPunctuationToken(next) {
			return false
		}

		if next.Value == "'" && quoteOpen {
			return false
		}
	}

	return true
}

func nextNonSpaceTokenIndex(tokens []model.Token, start int) int {
	for i := start; i < len(tokens); i++ {
		if !tokens[i].IsSpace() {
			return i
		}
	}

	return -1
}

func appendSpaceToken(tokens *[]model.Token) {
	if len(*tokens) == 0 {
		return
	}

	if (*tokens)[len(*tokens)-1].IsSpace() {
		(*tokens)[len(*tokens)-1].Value = " "
		return
	}

	*tokens = append(*tokens, model.Token{Kind: model.KindSpace, Value: " "})
}

func trimTrailingSpaces(tokens *[]model.Token) {
	for len(*tokens) > 0 {
		lastIndex := len(*tokens) - 1
		if !(*tokens)[lastIndex].IsSpace() {
			break
		}

		*tokens = (*tokens)[:lastIndex]
	}
}

func lastNonSpaceToken(tokens []model.Token) (model.Token, bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		if !tokens[i].IsSpace() {
			return tokens[i], true
		}
	}

	return model.Token{}, false
}

func isApostropheInWord(tokens []model.Token, index int) bool {
	if index <= 0 || index >= len(tokens)-1 {
		return false
	}

	prev := tokens[index-1]
	next := tokens[index+1]

	return prev.IsWord() && next.IsWord()
}
