package core

import (
	"errors"
	"strings"
)

const (
	QuerySep = " "
)

type Command struct {
	Command    string            `json:"command"`
	Args       map[string]string `json:"args"`
	FieldNames []string          `json:"field_names"`
	Sep        string            `json:"sep"`
}

func ProcessingInputByFieldNames(input string, command *Command) error {
	args := strings.Split(input, command.Sep)
	if len(args) < len(command.FieldNames) {
		return errors.New("Args len smaller fields len")
	}

	if command.Args == nil {
		command.Args = map[string]string{}
	}

	for i, field := range command.FieldNames {
		command.Args[field] = args[i]
	}

	return nil
}

func ProcessingInput(input string, sep string) (command Command) {
	args := strings.Split(input, sep)
	command.Command = args[0]

	command.Args = map[string]string{}

	for _, str := range args[1:] {
		queryRaw := strings.Split(str, "=")
		command.Args[queryRaw[0]] = queryRaw[1]
	}

	return
}
