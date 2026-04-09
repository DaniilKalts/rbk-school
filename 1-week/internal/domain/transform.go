package domain

func TransformText(input string) string {
	tokens := Tokenize(input)
	processed := applyCommands(tokens)

	return formatPunctuation(processed)
}
