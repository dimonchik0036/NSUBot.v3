package core

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"
)

const (
	QuerySep         = " "
	StrPreviousCmd   = "pc"
	StrJSPreviousCmd = "jpc"
	StrCmd           = "c"
)

type Command struct {
	Command    string            `json:"command"`
	MoreArgs   bool              `json:"more_args"`
	Args       map[string]string `json:"args"`
	ArgsStr    []string          `json:"args_str"`
	FieldNames []string          `json:"field_names"`
	Sep        string            `json:"sep"`
}

func (c *Command) SetArg(key string, value string) {
	if c.Args == nil {
		c.Args = map[string]string{}
	}

	c.Args[key] = value
}

func (c *Command) GetArg(key string) string {
	if c.Args == nil {
		return ""
	}

	return c.Args[key]
}

func (c *Command) GetArgInt64(key string) int64 {
	id, _ := strconv.ParseInt(c.GetArg(key), 10, 64)
	return id
}

func (c *Command) GetJSPreviousCommand() (Command, bool) {
	if c.Args == nil {
		return Command{}, false
	}

	cmd, ok := c.Args[StrJSPreviousCmd]
	if !ok {
		return Command{}, false
	}

	command := Command{}
	if err := json.Unmarshal([]byte(cmd), &command); err != nil {
		return Command{}, false
	}

	return command, true
}

func (c *Command) Encode() string {
	values := url.Values{}
	values.Set(StrCmd, c.Command)
	for key, val := range c.Args {
		if key == "cID" || key == "mID" {
			continue
		}

		values.Set(key, val)
	}

	return values.Encode()
}

func (c *Command) SetPreviousCommand() {
	c.SetArg(StrPreviousCmd, c.Encode())
}

func (c *Command) SetJSPreviousCommand() {
	cmd, err := json.Marshal(c)
	if err != nil {
		log.Print(err)
		return
	}

	if c.Args == nil {
		c.Args = map[string]string{}
	}

	c.Args[StrJSPreviousCmd] = string(cmd)
}

func GenerateCommandString(command string, args map[string]string) string {
	cmd := Command{
		Command: command,
		Args:    args,
	}
	return cmd.Encode()
}

func ProcessingInputByFieldNames(input string, command *Command) error {
	args := strings.Split(input, " ")
	command.ArgsStr = args

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
	command.SetArg("c_a", strings.Join(args[1:], " "))
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

func UnescapedInput(input string) (command Command) {
	values, err := url.ParseQuery(input)
	if err != nil {
		log.Print(err)
		return Command{}
	}

	command.Args = map[string]string{}
	for key, v := range values {
		if key == StrCmd {
			s := strings.Split(v[0], "*")
			command.Command = s[0]
			if len(s) > 1 {
				command.Args["c_a"] = s[1]
			}
		} else {
			command.Args[key] = v[0]
		}
	}

	return command
}
