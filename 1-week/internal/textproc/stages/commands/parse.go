package commands

import (
	"strconv"
	"strings"
)

type Spec struct {
	Name  string
	Count int
}

func ParseSpec(raw string) (Spec, bool) {
	if len(raw) < 3 || raw[0] != '(' || raw[len(raw)-1] != ')' {
		return Spec{}, false
	}

	inner := raw[1 : len(raw)-1]
	if inner == "" || inner != strings.TrimSpace(inner) {
		return Spec{}, false
	}

	parts := strings.Split(inner, ",")
	if len(parts) == 1 {
		switch parts[0] {
		case "up", "low", "cap", "hex", "bin":
			return Spec{Name: parts[0], Count: 1}, true
		default:
			return Spec{}, false
		}
	}

	if len(parts) != 2 {
		return Spec{}, false
	}

	name := parts[0]
	if !isCountCommandName(name) {
		return Spec{}, false
	}

	countValue := strings.TrimSpace(parts[1])
	count, err := strconv.Atoi(countValue)
	if err != nil || count <= 0 {
		return Spec{}, false
	}

	return Spec{Name: name, Count: count}, true
}

func isCountCommandName(command string) bool {
	switch command {
	case "up", "low", "cap":
		return true
	default:
		return false
	}
}
