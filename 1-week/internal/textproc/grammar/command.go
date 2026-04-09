package grammar

import (
	"strconv"
	"strings"
)

type CommandSpec struct {
	Name  string
	Count int
}

const (
	CommandUp  = "up"
	CommandLow = "low"
	CommandCap = "cap"
	CommandHex = "hex"
	CommandBin = "bin"
)

func ParseCommandSpec(raw string) (CommandSpec, bool) {
	if len(raw) < 3 || raw[0] != '(' || raw[len(raw)-1] != ')' {
		return CommandSpec{}, false
	}

	inner := raw[1 : len(raw)-1]
	if inner == "" || inner != strings.TrimSpace(inner) {
		return CommandSpec{}, false
	}

	parts := strings.Split(inner, ",")
	if len(parts) == 1 {
		switch parts[0] {
		case CommandUp, CommandLow, CommandCap, CommandHex, CommandBin:
			return CommandSpec{Name: parts[0], Count: 1}, true
		default:
			return CommandSpec{}, false
		}
	}

	if len(parts) != 2 {
		return CommandSpec{}, false
	}

	name := parts[0]
	if !isCountCommandName(name) {
		return CommandSpec{}, false
	}

	countValue := strings.TrimSpace(parts[1])
	count, err := strconv.Atoi(countValue)
	if err != nil || count <= 0 {
		return CommandSpec{}, false
	}

	return CommandSpec{Name: name, Count: count}, true
}

func isCountCommandName(command string) bool {
	switch command {
	case CommandUp, CommandLow, CommandCap:
		return true
	default:
		return false
	}
}
