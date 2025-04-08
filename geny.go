package main

import (
	"fmt"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/MUlt1mate/geny/commands"
)

type geny struct {
}

func (g *geny) ParseText(text string) (batch *CommandBatch, err error) {
	var command Command
	batch = &CommandBatch{}
	for line := range strings.Lines(text) {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		switch {
		case strings.HasPrefix(line, "protoc"):
			if command, err = commands.ParseProtoc(line); err != nil {
				return nil, err
			}
			batch.Commands = append(batch.Commands, command)
		default:
			batch.Commands = append(batch.Commands, commands.ParseSimple(line))
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
		case commands.CommandTypeSimple:
			newCommand := &commands.SimpleCommand{}
			if err = command.Body.Decode(&newCommand.Body); err != nil {
				return nil, err
			}
			batch.Commands[i] = newCommand
		case commands.CommandTypeProtoc:
			newCommand := &commands.ProtocCommand{}
			if err = command.Body.Decode(&newCommand.Body); err != nil {
				return nil, err
			}
			batch.Commands[i] = newCommand
		default:
			return nil, fmt.Errorf("geny: unknown command type: %s", command.Type)
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
