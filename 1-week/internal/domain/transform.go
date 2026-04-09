package domain

func TransformText(input string) string {
	tokens := Tokenize(input)
	processed := applyCommands(tokens)
	processed = applyArticles(processed)
	normalized := formatPunctuation(processed)

	return wrapText(normalized, 96)
}
