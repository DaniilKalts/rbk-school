package stages

import "github.com/DaniilKalts/rbk-school/1-week/internal/textproc/token"

func FindNextIndex(tokens []token.Token, start int, pred func(token.Token) bool) (int, bool) {
	for i := start; i < len(tokens); i++ {
		if pred(tokens[i]) {
			return i, true
		}
	}

	return -1, false
}

func FindPrevIndex(tokens []token.Token, start int, pred func(token.Token) bool) (int, bool) {
	for i := start; i >= 0; i-- {
		if pred(tokens[i]) {
			return i, true
		}
	}

	return -1, false
}
