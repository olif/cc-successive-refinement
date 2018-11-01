package main

// Create a list for all errors encountered

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	// ErrorCodeOk - ok code
	ErrorCodeOk = 0
	// ErrorCodeUnexpectedArgument -
	ErrorCodeUnexpectedArgument = 1
	// ErrorCodeMissingString  - missing string code
	ErrorCodeMissingString = 2
	// ErrorCodeMissingInteger  - missing integer code
	ErrorCodeMissingInteger = 3
	// ErrorCodeInvalidInteger  - invalid integer code
	ErrorCodeInvalidInteger = 4
	// ErrorCodeMissingFloat  - missing float code
	ErrorCodeMissingFloat = 5
	// ErrorCodeInvalidFloat  - invalid float code
	ErrorCodeInvalidFloat = 6
	// ErrorCodeInvalidArgumentName  - invalid name code
	ErrorCodeInvalidArgumentName = 7
)

type argsErrorCode int

type argsError struct {
	errorArgumentID rune
	errorParameter  string
	errorCode       argsErrorCode
}

func (e *argsError) Error() string {
	switch e.errorCode {
	case ErrorCodeOk:
		return fmt.Sprintf("TILT: Should not get here")
	case ErrorCodeUnexpectedArgument:
		return fmt.Sprintf("Argument -%s unexpected", string(e.errorArgumentID))
	case ErrorCodeMissingString:
		return fmt.Sprintf("Could not find string parameter for -%s", string(e.errorArgumentID))
	case ErrorCodeInvalidInteger:
		return fmt.Sprintf("Argument -%s expects an integer but was '%s'", string(e.errorArgumentID), e.errorParameter)
	case ErrorCodeInvalidFloat:
		return fmt.Sprintf("Argument -%s expects a float but was '%s'", string(e.errorArgumentID), e.errorParameter)
	case ErrorCodeMissingFloat:
		return fmt.Sprintf("Could not find double parameter for -%s", string(e.errorArgumentID))
	case ErrorCodeInvalidArgumentName:
		return fmt.Sprintf("'%s' is not a valid argument name", string(e.errorArgumentID))
	}

	return ""
}

// Args -
type Args struct {
	schema              string
	args                []string
	unexpectedArguments []rune
	argsFound           []rune
	marhalers           map[rune]ArgumentMarshaler
	currentArgument     *iterator
}

// NewArgs - returns a new ArgParser
func NewArgs(schema string, args []string) (*Args, error) {
	a := Args{
		schema:              schema,
		args:                args,
		unexpectedArguments: make([]rune, 0),
		argsFound:           make([]rune, 0),
		marhalers:           make(map[rune]ArgumentMarshaler),
		currentArgument:     nil,
	}

	return &a, a.parse()
}

func (a *Args) parse() error {
	if len(a.schema) == 0 && len(a.args) == 0 {
		return nil
	}
	if err := a.parseSchema(); err != nil {
		return err
	}
	if err := a.parseArguments(); err != nil {
		return err
	}
	//return len(a.unexpectedArguments) == 0, nil
	return nil
}

func (a *Args) parseSchema() error {
	for _, element := range strings.Split(a.schema, ",") {
		if len(element) > 0 {
			trimmedElement := strings.TrimSpace(element)
			if err := a.parseSchemaElement(trimmedElement); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Args) parseSchemaElement(element string) error {
	elementID := rune(element[0])
	elementTail := element[1:]
	if err := a.validateSchemaElement(elementID); err != nil {
		return err
	}

	if len(elementTail) == 0 {
		a.marhalers[elementID] = &booleanArgumentMarshaler{}
	} else if elementTail == "*" {
		a.marhalers[elementID] = &stringArgumentMarshaler{}
	} else if elementTail == "#" {
		a.marhalers[elementID] = &integerArgumentMarshaler{}
	} else if elementTail == "##" {
		a.marhalers[elementID] = &floatArgumentMarshaler{}
	} else {
		return fmt.Errorf("%s has invalid format: %s", string(elementID), elementTail)
	}

	return nil
}

func (a *Args) validateSchemaElement(elementID rune) error {
	if !unicode.IsLetter(elementID) {
		return fmt.Errorf("Bad characted %s in Args format: %s", string(elementID), a.schema)
	}
	return nil
}

func (a *Args) parseArguments() error {
	for a.currentArgument = newIterator(a.args); a.currentArgument.hasNext(); {
		arg := a.currentArgument.next()
		if err := a.parseArgument(arg); err != nil {
			return err
		}
	}

	return nil
}

func (a *Args) parseArgument(arg string) error {
	if string(arg[0]) == "-" {
		return a.parseElements(arg)
	}
	return nil
}

func (a *Args) parseElements(arg string) error {
	for i := 1; i < len(arg); i++ {
		if err := a.parseElement(rune(arg[i])); err != nil {
			return err
		}
	}
	return nil
}

func (a *Args) parseElement(argChar rune) error {
	if err := a.setArgument(argChar); err == nil {
		a.argsFound = append(a.argsFound, argChar)
	} else {
		a.unexpectedArguments = append(a.unexpectedArguments, argChar)
		argsErr := err.(*argsError)
		argsErr.errorArgumentID = argChar
		return argsErr
	}

	return nil
}

func (a *Args) setArgument(argChar rune) error {
	m, ok := a.marhalers[argChar]
	if !ok {
		return fmt.Errorf("Not a valid arg type")
	}
	return m.set(a.currentArgument)
}

// Cardinality - Returns the number of arguments
func (a *Args) Cardinality() int {
	return len(a.argsFound)
}

// Usage - Returns a string describing the usage
func (a *Args) Usage() string {
	if len(a.schema) > 0 {
		return fmt.Sprintf("[%s]", a.schema)
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

// GetString - Returns a string given its argument key
func (a *Args) GetString(arg rune) string {
	if stringMarshaler, ok := a.marhalers[arg]; ok {
		return stringMarshaler.get().(string)
	}
	return ""
}

// GetInteger - Returns an integer given its argument key
func (a *Args) GetInteger(arg rune) int {
	if intMarshaler, ok := a.marhalers[arg]; ok {
		return intMarshaler.get().(int)
	}

	return 0
}

// GetDouble - Returns an integer given its argument key
func (a *Args) GetDouble(arg rune) float64 {
	if doubleMarshaler, ok := a.marhalers[arg]; ok {
		return doubleMarshaler.get().(float64)
	}

	return 0
}

// Has - Returns true if the argument key exists, otherwise false
func (a *Args) Has(arg rune) bool {
	for _, v := range a.argsFound {
		if v == arg {
			return true
		}
	}
	return false
}

// ArgumentMarshaler - Parses cmd line arguments and makes them retrievable by key
type ArgumentMarshaler interface {
	get() interface{}
	set(i *iterator) error
}

type booleanArgumentMarshaler struct {
	boolVal bool
}

func (b *booleanArgumentMarshaler) set(i *iterator) error {
	b.boolVal = true
	return nil
}

func (b *booleanArgumentMarshaler) get() interface{} {
	return b.boolVal
}

type stringArgumentMarshaler struct {
	stringVal string
}

func (s *stringArgumentMarshaler) set(i *iterator) error {
	if !i.hasNext() {
		return &argsError{errorCode: ErrorCodeMissingString}
	}
	arg := i.next()
	s.stringVal = arg
	return nil
}

func (s *stringArgumentMarshaler) get() interface{} {
	return s.stringVal
}

type integerArgumentMarshaler struct {
	intVal int
}

func (m *integerArgumentMarshaler) set(i *iterator) error {
	if !i.hasNext() {
		return &argsError{errorCode: ErrorCodeMissingInteger}
	}

	param := i.next()
	intVal, err := strconv.Atoi(param)
	if err != nil {
		return &argsError{
			errorCode:      ErrorCodeInvalidInteger,
			errorParameter: param,
		}
	}
	m.intVal = intVal
	return nil
}

func (m *integerArgumentMarshaler) get() interface{} {
	return m.intVal
}

type floatArgumentMarshaler struct {
	floatVal float64
}

func (m *floatArgumentMarshaler) set(i *iterator) error {
	if !i.hasNext() {
		return &argsError{errorCode: ErrorCodeMissingFloat}
	}

	param := i.next()
	floatVal, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return &argsError{
			errorCode:      ErrorCodeInvalidFloat,
			errorParameter: param,
		}
	}
	m.floatVal = floatVal
	return nil
}

func (m *floatArgumentMarshaler) get() interface{} {
	return m.floatVal
}

type iterator struct {
	pos      int
	elements []string
}

func newIterator(elements []string) *iterator {
	return &iterator{
		pos:      -1,
		elements: elements,
	}
}

func (i *iterator) next() string {
	i.pos++
	return i.elements[i.pos]
}

func (i *iterator) hasNext() bool {
	return i.pos < len(i.elements)-1
}
