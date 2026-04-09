package render

import (
	"strings"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/token"
)

func Join(tokens []token.Token) string {
	var builder strings.Builder

	for _, tok := range tokens {
		builder.WriteString(tok.Value)
	}

	return builder.String()
}
