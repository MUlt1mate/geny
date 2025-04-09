package app

import (
	"gopkg.in/yaml.v3"
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
)
