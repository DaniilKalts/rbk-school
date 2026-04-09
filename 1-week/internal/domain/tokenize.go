package domain

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
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
		case isPunctuation(r):
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
	if !isValidCommand(raw) {
		return "", false
	}

	return raw, true
}

func isValidCommand(raw string) bool {
	n := len(raw)
	if n < 3 || raw[0] != '(' || raw[n-1] != ')' {
		return false
	}

	inner := raw[1 : n-1]
	if inner == "" || inner != strings.TrimSpace(inner) {
		return false
	}

	switch inner {
	case "up", "low", "cap", "hex", "bin":
		return true
	}

	parts := strings.Split(inner, ",")
	if len(parts) != 2 {
		return false
	}

	command := parts[0]
	if !isCountCommandName(command) {
		return false
	}

	count := strings.TrimSpace(parts[1])
	if count == "" {
		return false
	}

	n, err := strconv.Atoi(count)
	if err != nil || n <= 0 {
		return false
	}

	return true
}

func isCountCommandName(command string) bool {
	switch command {
	case "up", "low", "cap":
		return true
	}

	return false
}

func isPunctuation(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';', '\'':
		return true
	default:
		return false
	}
}
