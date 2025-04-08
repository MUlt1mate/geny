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
		Body struct {
			Parts []string
		}
	}
)

func ParseSimple(input string) (command *SimpleCommand) {
	return &SimpleCommand{
		Type: CommandTypeSimple,
		Body: struct{ Parts []string }{Parts: strings.Split(input, " ")},
	}
}

func (s *SimpleCommand) String() string {
	return strings.Join(s.Body.Parts, " ")
}
