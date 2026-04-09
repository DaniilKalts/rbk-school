package punctuation

import (
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/grammar"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/stages"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/token"
)

func Apply(tokens []token.Token) []token.Token {
	result := make([]token.Token, 0, len(tokens))
	quoteOpen := false

	for i, tok := range tokens {
		if tok.IsSpace() {
			handleSpaceToken(&result, tokens, i, quoteOpen)
			continue
		}

		handled, nextQuoteOpen := handleNonSpaceToken(&result, tokens, i, quoteOpen)
		quoteOpen = nextQuoteOpen
		if handled {
			continue
		}

		result = append(result, tok)
	}

	return result
}

func handleNonSpaceToken(result *[]token.Token, tokens []token.Token, index int, quoteOpen bool) (bool, bool) {
	tok := tokens[index]

	if shouldInsertSpaceAfterTightPunctuation(*result, tok, quoteOpen) {
		appendSpaceToken(result)
	}

	return preprocessPunctuationToken(result, tokens, index, quoteOpen)
}

func handleSpaceToken(result *[]token.Token, tokens []token.Token, index int, quoteOpen bool) {
	nextIndex, ok := nextNonSpaceTokenIndex(tokens, index+1)
	if !ok || len(*result) == 0 {
		return
	}

	nextToken := tokens[nextIndex]
	prevToken, ok := lastNonSpaceToken(*result)
	if !ok {
		return
	}

	if shouldSkipSpace(prevToken, nextToken, quoteOpen) {
		return
	}

	appendSpaceToken(result)
}

func preprocessPunctuationToken(result *[]token.Token, tokens []token.Token, index int, quoteOpen bool) (bool, bool) {
	tok := tokens[index]
	if !tok.IsPunctuation() {
		return false, quoteOpen
	}

	if tok.Value != "'" {
		if isTightPunctuationToken(tok) {
			trimTrailingSpaces(result)
		}

		return false, quoteOpen
	}

	if isApostropheInWord(tokens, index) {
		return false, quoteOpen
	}

	if quoteOpen {
		trimTrailingSpaces(result)
		*result = append(*result, tok)
		return true, false
	}

	*result = append(*result, tok)
	return true, true
}

func isTightPunctuationToken(tok token.Token) bool {
	if !tok.IsPunctuation() || len(tok.Value) != 1 {
		return false
	}

	return grammar.IsTightPunctuationRune(rune(tok.Value[0]))
}

func shouldInsertSpaceAfterTightPunctuation(result []token.Token, current token.Token, quoteOpen bool) bool {
	prevToken, ok := lastNonSpaceToken(result)
	if !ok || !isTightPunctuationToken(prevToken) {
		return false
	}

	if isNoSpaceSuccessor(current, quoteOpen) {
		return false
	}

	return true
}

func shouldSkipSpace(prev token.Token, next token.Token, quoteOpen bool) bool {
	if isNoSpaceSuccessor(next, quoteOpen) {
		return true
	}

	if prev.IsPunctuation() && prev.Value == "'" && quoteOpen {
		return true
	}

	return false
}

func isNoSpaceSuccessor(tok token.Token, quoteOpen bool) bool {
	if !tok.IsPunctuation() {
		return false
	}

	if isTightPunctuationToken(tok) {
		return true
	}

	return tok.Value == "'" && quoteOpen
}

func nextNonSpaceTokenIndex(tokens []token.Token, start int) (int, bool) {
	return stages.FindNextIndex(tokens, start, func(tok token.Token) bool {
		return !tok.IsSpace()
	})
}

func appendSpaceToken(tokens *[]token.Token) {
	if len(*tokens) == 0 {
		return
	}

	if (*tokens)[len(*tokens)-1].IsSpace() {
		(*tokens)[len(*tokens)-1].Value = " "
		return
	}

	*tokens = append(*tokens, token.Token{Kind: token.KindSpace, Value: " "})
}

func trimTrailingSpaces(tokens *[]token.Token) {
	for len(*tokens) > 0 {
		lastIndex := len(*tokens) - 1
		if !(*tokens)[lastIndex].IsSpace() {
			break
		}

		*tokens = (*tokens)[:lastIndex]
	}
}

func lastNonSpaceToken(tokens []token.Token) (token.Token, bool) {
	index, ok := stages.FindPrevIndex(tokens, len(tokens)-1, func(tok token.Token) bool {
		return !tok.IsSpace()
	})
	if !ok {
		return token.Token{}, false
	}

	return tokens[index], true
}

func isApostropheInWord(tokens []token.Token, index int) bool {
	if index <= 0 || index >= len(tokens)-1 {
		return false
	}

	prev := tokens[index-1]
	next := tokens[index+1]

	return prev.IsWord() && next.IsWord()
}
