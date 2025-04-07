package main

import "strings"

type (
	Command interface {
		Parse(input string) error
		String() string
	}
	CommandBatch struct {
		Simple []*SimpleCommand
	}
	SimpleCommand struct {
		Parts []string
	}
)

func (s *SimpleCommand) Parse(input string) (err error) {
	s.Parts = strings.Split(input, " ")
	return nil
}

func (s *SimpleCommand) String() string {
	return strings.Join(s.Parts, " ")
}
