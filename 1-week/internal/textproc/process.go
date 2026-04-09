package textproc

import (
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/render"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/stages/articles"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/stages/commands"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/stages/punctuation"
	"github.com/DaniilKalts/rbk-school/1-week/internal/textproc/tokenizer"
)

const defaultWrapWidth = 96

func Process(input string) string {
	tokens := tokenizer.Tokenize(input)
	tokens = commands.Apply(tokens)
	tokens = articles.Apply(tokens)
	tokens = punctuation.Apply(tokens)

	return render.Wrap(render.Join(tokens), defaultWrapWidth)
}
