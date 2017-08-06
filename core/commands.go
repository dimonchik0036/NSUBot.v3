package core

import (
	"errors"
	"strconv"
	"strings"
)

const (
	QuerySep = " "
)

type Command struct {
	Command    string            `json:"command"`
	MoreArgs   bool              `json:"more_args"`
	Args       map[string]string `json:"args"`
	ArgsStr    []string          `json:"args_str"`
	FieldNames []string          `json:"field_names"`
	Sep        string            `json:"sep"`
}

func (c *Command) GetArgInt64(key string) int64 {
	id, _ := strconv.ParseInt(c.Args[key], 10, 64)
	return id
}

func ProcessingInputByFieldNames(input string, command *Command) error {
	args := strings.Split(input, command.Sep)
	if command.MoreArgs {
		command.ArgsStr = args
		return nil
	}

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

func SearchCommand(input string, sep string) (command Command) {
	args := strings.Split(input, sep)
	command.Command = args[0]
	command.ArgsStr = append(command.ArgsStr, args[1:]...)
	command.Sep = sep
	return
}

func ProcessingInput(input string, sep string) (command Command) {
	args := strings.Split(input, sep)
	command.Command = args[0]

	command.Args = map[string]string{}

	for _, str := range args[1:] {
		queryRaw := strings.Split(str, "=")
		command.Args[queryRaw[0]] = queryRaw[1]
	}
	command.Sep = sep
	return
}
