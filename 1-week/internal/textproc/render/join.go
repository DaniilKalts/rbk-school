package render

import (
	"strings"

	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/model"
)

func Join(tokens []model.Token) string {
	var builder strings.Builder

	for _, token := range tokens {
		builder.WriteString(token.Value)
	}

	return builder.String()
}
