package main

import (
	"bytes"
	"fmt"
	"strconv"
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
	// ErrorCodeMissingInteger  - missing integer code
	ErrorCodeMissingInteger = 2
	// ErrorCodeInvalidInteger  - invalid integer code
	ErrorCodeInvalidInteger = 3
)

// Args -
type Args struct {
	schema              string
	args                []string
	nrOfArguments       int
	unexpectedArguments []rune
	argsFound           []rune
	marhalers           map[rune]ArgumentMarshaler
	currentArgument     int
}

type ArgumentMarshaler interface {
	get() interface{}
	set(s string) error
}

type BooleanArgumentMarshaler struct {
	ArgumentMarshaler
	boolVal bool
}

func (b *BooleanArgumentMarshaler) set(val string) error {
	b.boolVal = true
	return nil
}

func (b *BooleanArgumentMarshaler) get() interface{} {
	return b.boolVal
}

type StringArgumentMarshaler struct {
	ArgumentMarshaler
	stringVal string
}

func (s *StringArgumentMarshaler) set(val string) error {
	s.stringVal = val
	return nil
}

func (s *StringArgumentMarshaler) get() interface{} {
	return s.stringVal
}

type IntegerArgumentMarshaler struct {
	ArgumentMarshaler
	intVal int
}

func (m *IntegerArgumentMarshaler) set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	m.intVal = val
	return nil
}

func (s *IntegerArgumentMarshaler) get() interface{} {
	return s.intVal
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
		marhalers:           make(map[rune]ArgumentMarshaler),
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
		a.marhalers[elementID] = &BooleanArgumentMarshaler{}
	} else if isStringSchemaElement(elementTail) {
		a.marhalers[elementID] = &StringArgumentMarshaler{}
	} else if isIntegerSchemaElement(elementTail) {
		a.marhalers[elementID] = &IntegerArgumentMarshaler{}
	}

	return nil
}

func (a *Args) validateSchemaElement(elementID rune) error {
	if !unicode.IsLetter(elementID) {
		return fmt.Errorf("Bad characted %s in Args format: %s", string(elementID), a.schema)
	}
	return nil
}

func isIntegerSchemaElement(elementTail string) bool {
	return elementTail == "#"
}

func isStringSchemaElement(elementTail string) bool {
	return elementTail == "*"
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
	if err := a.setArgument(argChar); err == nil {
		a.argsFound = append(a.argsFound, argChar)
	} else {
		a.unexpectedArguments = append(a.unexpectedArguments, argChar)
		valid = false
	}
}

func (a *Args) setArgument(argChar rune) error {
	m := a.marhalers[argChar]
	switch m.(type) {
	case *BooleanArgumentMarshaler:
		return a.setBoolArg(m)
	case *StringArgumentMarshaler:
		return a.setStringArg(m)
	case *IntegerArgumentMarshaler:
		return a.setIntArg(m)
	}

	return fmt.Errorf("Not a valid arg type")
}

func (a *Args) setStringArg(m ArgumentMarshaler) error {
	a.currentArgument++
	if a.currentArgument >= len(a.args) {
		valid = false
		errorCode = ErrorCodeMissingString
		return fmt.Errorf("Missing string argument")
	}
	arg := a.args[a.currentArgument]
	return m.set(arg)
}

func (a *Args) setBoolArg(m ArgumentMarshaler) error {
	return m.set("true")
}

func (a *Args) setIntArg(m ArgumentMarshaler) error {
	a.currentArgument++
	if a.currentArgument > len(a.args) {
		valid = false
		errorCode = ErrorCodeMissingInteger
		return fmt.Errorf("Missing integer argument")
	}

	arg := a.args[a.currentArgument]
	return m.set(arg)
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
	if argMarshaler, ok := a.marhalers[arg]; ok {
		return argMarshaler.get().(bool)
	}
	return false
}

func (a *Args) GetString(arg rune) string {
	if stringMarshaler, ok := a.marhalers[arg]; ok {
		return stringMarshaler.get().(string)
	}
	return ""
}

func (a *Args) GetInteger(arg rune) int {
	if intMarshaler, ok := a.marhalers[arg]; ok {
		return intMarshaler.get().(int)
	}

	return 0
}

func (a *Args) Has(arg rune) bool {
	for _, v := range a.argsFound {
		if v == arg {
			return true
		}
	}
	return false
}
