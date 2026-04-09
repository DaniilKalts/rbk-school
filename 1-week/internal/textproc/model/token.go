package model

type Kind int

const (
	KindWord Kind = iota
	KindCommand
	KindSpace
	KindPunctuation
)

type Token struct {
	Kind  Kind
	Value string
}

func (t Token) IsWord() bool {
	return t.Kind == KindWord
}

func (t Token) IsCommand() bool {
	return t.Kind == KindCommand
}

func (t Token) IsSpace() bool {
	return t.Kind == KindSpace
}

func (t Token) IsPunctuation() bool {
	return t.Kind == KindPunctuation
}
