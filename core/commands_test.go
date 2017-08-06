package core

import (
	"fmt"
	"testing"
)

func TestProcessingInputByFieldNames(t *testing.T) {
	Field1 := "first_name"
	Field2 := "last_name"
	command := Command{
		Command:    "Hello",
		FieldNames: []string{Field1, Field2},
		Sep:        " ",
	}

	if err := ProcessingInputByFieldNames("Ivan Ivanov", &command); err != nil {
		t.Fatal(err)
	}

	if command.Args[Field1] != "Ivan" || command.Args[Field2] != "Ivanov" {
		t.FailNow()
	}

	if err := ProcessingInputByFieldNames("IvanIvanov", &command); err == nil {
		t.FailNow()
	}
}

func TestProcessingInput(t *testing.T) {
	command := "Hello"
	field1 := "first_name"
	fiend2 := "last_name"
	arg1 := "Ivan"
	arg2 := "Ivanov"

	c := ProcessingInput(fmt.Sprintf("%s %s=%s %s=%s", command, field1, arg1, fiend2, arg2), " ")
	if c.Command != command || c.Args[field1] != arg1 || c.Args[fiend2] != arg2 {
		t.FailNow()
	}
}

func TestCommand_String(t *testing.T) {
	var cmd Command
	cmd.Command = "hey!;"
	cmd.Args = map[string]string{
		"dfs _; 23y 98":       "dse 4s6&tt",
		"ffd; sf_ _fs; ; ;  ": "sf",
	}
	fmt.Println(cmd)
	strCmd := cmd.Encode()
	fmt.Println(strCmd)
	newCmd := UnescapedInput(strCmd)
	fmt.Println(newCmd)
}
