package render

import "strings"

func Wrap(text string, width int) string {
	if width <= 0 || text == "" {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	lines := make([]string, 0)
	current := words[0]

	for _, word := range words[1:] {
		if len(current)+1+len(word) <= width {
			current += " " + word
			continue
		}

		lines = append(lines, current)
		current = word
	}

	lines = append(lines, current)

	return strings.Join(lines, "\n") + "\n"
}
