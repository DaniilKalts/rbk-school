package token

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/grammar"
)

func Tokenize(input string) []Token {
	tokens := make([]Token, 0)

	for len(input) > 0 {
		r, size := utf8.DecodeRuneInString(input)

		switch {
		case unicode.IsSpace(r):
			value := readSpace(input)
			tokens = append(tokens, Token{Kind: KindSpace, Value: value})
			input = input[len(value):]
		case grammar.IsPunctuationRune(r):
			tokens = append(tokens, Token{Kind: KindPunctuation, Value: input[:size]})
			input = input[size:]
		case r == '(':
			if value, ok := readCommand(input); ok {
				tokens = append(tokens, Token{Kind: KindCommand, Value: value})
				input = input[len(value):]
				continue
			}

			value := readTextToken(input)
			tokens = append(tokens, Token{Kind: KindWord, Value: value})
			input = input[len(value):]
		default:
			value := readTextToken(input)
			tokens = append(tokens, Token{Kind: KindWord, Value: value})
			input = input[len(value):]
		}
	}

	return tokens
}

func readSpace(input string) string {
	end := 0

	for i, r := range input {
		if !unicode.IsSpace(r) {
			break
		}

		end = i + utf8.RuneLen(r)
	}

	return input[:end]
}

func readTextToken(input string) string {
	end := 0

	for i, r := range input {
		if unicode.IsSpace(r) || grammar.IsPunctuationRune(r) {
			break
		}

		if r == '(' {
			if _, ok := readCommand(input[i:]); ok {
				break
			}
		}

		end = i + utf8.RuneLen(r)
	}

	if end == 0 {
		_, size := utf8.DecodeRuneInString(input)
		return input[:size]
	}

	return input[:end]
}

func readCommand(input string) (string, bool) {
	end := strings.IndexByte(input, ')')
	if end == -1 {
		return "", false
	}

	raw := input[:end+1]
	if _, ok := grammar.ParseCommandSpec(raw); !ok {
		return "", false
	}

	return raw, true
}
