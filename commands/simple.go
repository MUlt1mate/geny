package commands

import (
	"strings"
)

const (
	CommandTypeSimple = "simple"
)

type (
	SimpleCommand struct {
		Type string
		Body SimpleParts
	}
	SimpleParts struct {
		Parts []string
	}
)

func ParseSimple(input string) (command *SimpleCommand) {
	parts := strings.Split(input, " ")
	s := &SimpleCommand{
		Type: CommandTypeSimple,
		Body: SimpleParts{Parts: make([]string, 0, len(parts))},
	}
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		s.Body.Parts = append(s.Body.Parts, strings.TrimSpace(part))
	}
	return s
}

func (s *SimpleCommand) String() string {
	return strings.Join(s.Body.Parts, " ")
}
