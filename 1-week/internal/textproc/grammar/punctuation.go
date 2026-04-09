package grammar

func IsPunctuationRune(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';', '\'':
		return true
	default:
		return false
	}
}

func IsTightPunctuationRune(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';':
		return true
	default:
		return false
	}
}
