package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

var valid = true

// Args -
type Args struct {
	schema              string
	args                []string
	nrOfArguments       int
	unexpectedArguments []rune
	booleanArgs         map[rune]bool
}

// NewArgs - returns a new ArgParser
func NewArgs(schema string, args []string) *Args {
	a := Args{
		schema:              schema,
		args:                args,
		nrOfArguments:       0,
		unexpectedArguments: make([]rune, 0),
		booleanArgs:         make(map[rune]bool),
	}

	valid = a.parse()
	return &a
}

func (a *Args) isValid() bool {
	return valid
}

func (a *Args) parse() bool {
	if len(a.schema) == 0 && len(a.args) == 0 {
		return true
	}
	a.parseSchema()
	a.parseArguments()
	return len(a.unexpectedArguments) == 0
}

func (a *Args) parseSchema() (bool, error) {
	for _, element := range strings.Split(a.schema, ",") {
		if len(element) > 0 {
			trimmedElement := strings.TrimSpace(element)
			a.parseSchemaElement(trimmedElement)
		}
	}
	return true, nil
}

func (a *Args) parseSchemaElement(element string) {
	if len(element) == 1 {
		a.parseBooleanSchemaElement(element)
	}
}

func (a *Args) parseBooleanSchemaElement(element string) {
	c := rune(element[0])
	if unicode.IsLetter(c) {
		a.booleanArgs[c] = false
	}
}

func (a *Args) parseArguments() bool {
	for _, arg := range a.args {
		a.parseArgument(arg)
	}

	return true
}

func (a *Args) parseArgument(arg string) {
	if string(arg[0]) == "-" {
		a.parseElements(arg)
	}
}

func (a *Args) parseElements(arg string) {
	for i := 1; i < len(arg); i++ {
		a.parseElement(rune(arg[i]))
	}
}

func (a *Args) parseElement(argChar rune) {
	if a.isBool(argChar) {
		a.setBoolArg(argChar, true)
	} else {
		a.unexpectedArguments = append(a.unexpectedArguments, argChar)
	}
}

func (a *Args) setBoolArg(argChar rune, value bool) {
	a.booleanArgs[argChar] = value
}

func (a *Args) isBool(argChar rune) bool {
	if _, ok := a.booleanArgs[argChar]; ok {
		return true
	}

	return false
}

// Cardinality - Returns the number of arguments
func (a *Args) Cardinality() int {
	return a.nrOfArguments
}

// Usage - Returns a string describing the usage
func (a *Args) Usage() string {
	if len(a.schema) > 0 {
		return fmt.Sprintf("[%s]", a.schema)
	} else {
		return ""
	}
}

// ErrorMessage - Returns a string containing an error message for all unexpected arguments
func (a *Args) ErrorMessage() string {
	if len(a.unexpectedArguments) > 0 {
		return a.unexpectedArgumentMessage()
	} else {
		return ""
	}
}

func (a *Args) unexpectedArgumentMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString("Argument(s) -")
	for _, v := range a.unexpectedArguments {
		buffer.WriteRune(v)
	}
	buffer.WriteString(" unexpected.")

	return buffer.String()
}

// GetBoolean - Returns the value of the boolean arg
func (a *Args) GetBoolean(arg rune) bool {
	return a.booleanArgs[arg]
}
