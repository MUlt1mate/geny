package main

import (
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type geny struct {
}

func (g *geny) ParseText(text string) (batch *CommandBatch, err error) {
	batch = &CommandBatch{}
	for line := range strings.Lines(text) {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch {
		default:
			batch.Commands = append(batch.Commands, ParseSimple(line))
		}
	}
	return batch, nil
}

func (g *geny) ParseYAML(data []byte) (batch *CommandBatch, err error) {
	var listRaw RawBatch
	if err = yaml.Unmarshal(data, &listRaw); err != nil {
		return nil, err
	}
	batch = &CommandBatch{}
	batch.Commands = make([]Command, len(listRaw.Commands))
	for i, command := range listRaw.Commands {
		switch command.Type {
		case CommandTypeSimple:
			newCommand := &SimpleCommand{}
			if err = command.Body.Decode(&newCommand.Body); err != nil {
				return nil, err
			}
			batch.Commands[i] = newCommand

		}
	}
	return batch, nil
}

func (g *geny) FormatGoFile(batch *CommandBatch) (output string) {
	var (
		footer = []string{"", "package main", ""}
		lines  = make([]string, 0, len(batch.Commands)+len(footer))
	)
	for _, command := range batch.Commands {
		lines = append(lines, "//go:generate "+command.String())
	}
	lines = append(lines, footer...)
	return strings.Join(lines, "\n")
}

func (g *geny) Exec(batch *CommandBatch) (err error) {
	for _, command := range batch.Commands {
		cmd := exec.Command(command.String())
		if err = cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
