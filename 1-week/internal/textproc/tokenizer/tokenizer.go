package tokenizer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/model"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/stages/commands"
)

func Tokenize(input string) []model.Token {
	tokens := make([]model.Token, 0)

	for len(input) > 0 {
		r, size := utf8.DecodeRuneInString(input)

		switch {
		case unicode.IsSpace(r):
			value := readSpace(input)
			tokens = append(tokens, model.Token{Kind: model.KindSpace, Value: value})
			input = input[len(value):]
		case isPunctuation(r):
			tokens = append(tokens, model.Token{Kind: model.KindPunctuation, Value: input[:size]})
			input = input[size:]
		case r == '(':
			if value, ok := readCommand(input); ok {
				tokens = append(tokens, model.Token{Kind: model.KindCommand, Value: value})
				input = input[len(value):]
				continue
			}

			value := readTextToken(input)
			tokens = append(tokens, model.Token{Kind: model.KindWord, Value: value})
			input = input[len(value):]
		default:
			value := readTextToken(input)
			tokens = append(tokens, model.Token{Kind: model.KindWord, Value: value})
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
		if unicode.IsSpace(r) || isPunctuation(r) {
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
	if _, ok := commands.ParseSpec(raw); !ok {
		return "", false
	}

	return raw, true
}

func isPunctuation(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';', '\'':
		return true
	default:
		return false
	}
}
