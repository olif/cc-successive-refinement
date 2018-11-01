package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

var (
	valid         = true
	errorArgument = '0'
	errorCode     = ErrorCodeOk
)

const (
	// ErrorCodeOk - ok code
	ErrorCodeOk = 0
	// ErrorCodeMissingString  - missing string code
	ErrorCodeMissingString = 1
)

// Args -
type Args struct {
	schema              string
	args                []string
	nrOfArguments       int
	unexpectedArguments []rune
	argsFound           []rune
	booleanArgs         map[rune]*BooleanArgumentMarshaler
	stringArgs          map[rune]string
	currentArgument     int
}

type ArgumentMarshaler struct {
	boolVal bool
}

type BooleanArgumentMarshaler struct {
	ArgumentMarshaler
}

type StringArgumentMarshaler struct {
	ArgumentMarshaler
}

type IntegerArgumentMarshaler struct {
	ArgumentMarshaler
}

func (a *ArgumentMarshaler) setBool(value bool) {
	a.boolVal = value
}

func (a *ArgumentMarshaler) getBool() bool {
	return a.boolVal
}

// NewArgs - returns a new ArgParser
func NewArgs(schema string, args []string) (*Args, error) {
	var err error
	a := Args{
		schema:              schema,
		args:                args,
		nrOfArguments:       0,
		unexpectedArguments: make([]rune, 0),
		argsFound:           make([]rune, 0),
		booleanArgs:         make(map[rune]*BooleanArgumentMarshaler),
		stringArgs:          make(map[rune]string),
	}

	valid, err = a.parse()
	return &a, err
}

func (a *Args) isValid() bool {
	return valid
}

func (a *Args) parse() (bool, error) {
	if len(a.schema) == 0 && len(a.args) == 0 {
		return true, nil
	}
	if ok, err := a.parseSchema(); err != nil {
		return ok, err
	}
	if ok := a.parseArguments(); !ok {
		return ok, nil
	}
	return len(a.unexpectedArguments) == 0, nil
}

func (a *Args) parseSchema() (bool, error) {
	for _, element := range strings.Split(a.schema, ",") {
		if len(element) > 0 {
			trimmedElement := strings.TrimSpace(element)
			if err := a.parseSchemaElement(trimmedElement); err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (a *Args) parseSchemaElement(element string) error {
	elementID := rune(element[0])
	elementTail := element[1:]
	if err := a.validateSchemaElement(elementID); err != nil {
		return err
	}
	if isBooleanSchemaElement(elementTail) {
		a.parseBoolSchemaElement(elementID)
	} else if isStringSchemaElement(elementTail) {
		a.parseStringSchemaElement(elementID)
	}

	return nil
}

func (a *Args) validateSchemaElement(elementID rune) error {
	if !unicode.IsLetter(elementID) {
		return fmt.Errorf("Bad characted %s in Args format: %s", string(elementID), a.schema)
	}

	return nil
}

func (a *Args) parseStringSchemaElement(elementID rune) {
	a.stringArgs[elementID] = ""
}

func isStringSchemaElement(elementTail string) bool {
	return elementTail == "*"
}

func (a *Args) parseBoolSchemaElement(elementID rune) {
	a.booleanArgs[elementID] = &BooleanArgumentMarshaler{}
}

func isBooleanSchemaElement(elementTail string) bool {
	return len(elementTail) == 0
}

func (a *Args) parseArguments() bool {
	for a.currentArgument = 0; a.currentArgument < len(a.args); a.currentArgument++ {
		a.parseArgument(a.args[a.currentArgument])
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
	if a.setArgument(argChar) {
		a.argsFound = append(a.argsFound, argChar)
	} else {
		a.unexpectedArguments = append(a.unexpectedArguments, argChar)
		valid = false
	}
}

func (a *Args) setArgument(argChar rune) bool {
	set := true
	if a.isBool(argChar) {
		a.setBoolArg(argChar, true)
	} else if a.isString(argChar) {
		a.setStringArg(argChar, "")
	} else {
		set = false
	}

	return set
}

func (a *Args) setStringArg(argChar rune, s string) {
	a.currentArgument++
	if a.currentArgument >= len(a.args) {
		valid = false
		errorArgument = argChar
		errorCode = ErrorCodeMissingString
		return
	}
	arg := a.args[a.currentArgument]
	a.stringArgs[argChar] = arg
}

func (a *Args) isString(argChar rune) bool {
	if _, ok := a.stringArgs[argChar]; ok {
		return true
	}

	return false
}

func (a *Args) setBoolArg(argChar rune, value bool) {
	a.booleanArgs[argChar].setBool(value)
}

func (a *Args) isBool(argChar rune) bool {
	if _, ok := a.booleanArgs[argChar]; ok {
		return true
	}

	return false
}

// Cardinality - Returns the number of arguments
func (a *Args) Cardinality() int {
	return len(a.argsFound)
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
		switch errorCode {
		case ErrorCodeMissingString:
			return fmt.Sprintf("Could not find string parameter for -%s", string(errorArgument))
		case ErrorCodeOk:
			return fmt.Sprintf("TILT: Should not get here")
		}
	}

	return ""
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
	if argMarshaler, ok := a.booleanArgs[arg]; ok {
		return argMarshaler.getBool()
	}
	return false
}

func (a *Args) GetString(arg rune) string {
	return a.stringArgs[arg]
}

func (a *Args) Has(arg rune) bool {
	for _, v := range a.argsFound {
		if v == arg {
			return true
		}
	}
	return false
}
