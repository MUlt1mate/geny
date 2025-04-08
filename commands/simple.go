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
	return &SimpleCommand{
		Type: CommandTypeSimple,
		Body: SimpleParts{Parts: strings.Split(input, " ")}}
}

func (s *SimpleCommand) String() string {
	return strings.Join(s.Body.Parts, " ")
}
