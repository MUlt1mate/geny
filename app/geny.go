package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/MUlt1mate/geny/commands"
)

const (
	goGeneratePrefix = "//go:generate "
)

type Geny struct {
}

func (g *Geny) ParseShellCommands(text string) (batch *CommandBatch, err error) {
	var command Command
	batch = &CommandBatch{}
	for line := range strings.Lines(text) {
		if command, err = g.parseLine(strings.TrimSpace(line)); err != nil {
			return nil, err
		}
		if command != nil {
			batch.Commands = append(batch.Commands, command)
		}
	}
	return batch, nil
}

func (g *Geny) ParseGoFile(text string) (batch *CommandBatch, err error) {
	var command Command
	batch = &CommandBatch{}
	for line := range strings.Lines(text) {
		if !strings.HasPrefix(line, goGeneratePrefix) {
			continue
		}
		if command, err = g.parseLine(strings.TrimSpace(strings.TrimPrefix(line, goGeneratePrefix))); err != nil {
			return nil, err
		}
		if command != nil {
			batch.Commands = append(batch.Commands, command)
		}
	}
	return batch, nil
}

func (g *Geny) parseLine(line string) (command Command, err error) {
	if line == "" || strings.HasPrefix(line, "//") {
		return nil, nil
	}
	switch {
	case strings.HasPrefix(line, "protoc"):
		if command, err = commands.ParseProtoc(line); err != nil {
			return nil, err
		}
		return command, nil
	default:
		return commands.ParseSimple(line), nil
	}
}

func (g *Geny) ParseYAML(data []byte) (batch *CommandBatch, err error) {
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

func (g *Geny) FormatGoFile(batch *CommandBatch, packageName string, showGenComment bool) (output string) {
	var (
		footer = []string{"", "package " + packageName, ""}
		lines  = make([]string, 0, len(batch.Commands)+len(footer)+1)
	)
	if showGenComment {
		lines = append(lines, "// Code generated by geny. DO NOT EDIT.")
	}

	for _, command := range batch.Commands {
		lines = append(lines, "//go:generate "+command.String())
	}
	lines = append(lines, footer...)
	return strings.Join(lines, "\n")
}

func (g *Geny) Exec(batch *CommandBatch) (err error) {
	for _, command := range batch.Commands {
		commandStr := command.String()
		log.Println(commandStr)
		parts := strings.Split(commandStr, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err = cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
