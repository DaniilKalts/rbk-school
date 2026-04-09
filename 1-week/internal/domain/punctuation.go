package domain

func formatPunctuation(tokens []Token) string {
	result := make([]Token, 0, len(tokens))
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

			result = append(result, token)
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

	return joinTokens(result)
}

func isTightPunctuationToken(token Token) bool {
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

func needsSpaceAfterTightPunctuation(result []Token, current Token, quoteOpen bool) bool {
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

func shouldSkipSpace(prev Token, next Token, quoteOpen bool) bool {
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

func shouldForceSingleSpace(prev Token, next Token, quoteOpen bool) bool {
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

func isApostropheInWord(tokens []Token, index int) bool {
	if index <= 0 || index >= len(tokens)-1 {
		return false
	}

	prev := tokens[index-1]
	next := tokens[index+1]

	return prev.IsWord() && next.IsWord()
}
