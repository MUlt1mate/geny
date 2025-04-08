package main

import (
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	CommandTypeSimple = "simple"
)

type (
	Command interface {
		String() string
	}
	CommandBatch struct {
		Commands []Command
	}
	RawBatch struct {
		Commands []RawCommand
	}
	RawCommand struct {
		Type string
		Body yaml.Node
	}
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
