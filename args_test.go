package args

import (
	"testing"
)

func TestCreateWithNoSchemaOrArguments(t *testing.T) {
	a, _ := NewArgs("", []string{})
	assertEquals(t, 0, a.Cardinality(), "Exptected 0 cardinality")
}

func TestWithNoSchemaButWithOneArgument(t *testing.T) {
	_, err := NewArgs("", []string{"-x"})
	assertEquals(t, ErrUnexpectedArgument, err.(*ArgsError).Code(), "")
	assertEquals(t, 'x', err.(*ArgsError).ArgumentID(), "Unexpected argument")
}

func TestWithNoSchemaButWithMultipleArguments(t *testing.T) {
	_, err := NewArgs("", []string{"-x", "-y"})
	assertEquals(t, ErrUnexpectedArgument, err.(*ArgsError).Code(), "Unexpected argument")
	assertEquals(t, 'x', err.(*ArgsError).ArgumentID(), "Unexpected argument")
}

func TestNonLetterSchema(t *testing.T) {
	_, err := NewArgs("*", []string{})
	assertEquals(t, ErrInvalidArgumentName, err.(*ArgsError).Code(), "Invalid argument name")
	assertEquals(t, '*', err.(*ArgsError).ArgumentID(), "Invalid argument name")
}

func TestInvalidArgumentFormat(t *testing.T) {
	_, err := NewArgs("f-", []string{})
	assertEquals(t, ErrInvalidFormat, err.(*ArgsError).Code(), "Invalid format")
	assertEquals(t, 'f', err.(*ArgsError).ArgumentID(), "Invalid format")
}

func TestSimpleBooleanPresent(t *testing.T) {
	a, err := NewArgs("x", []string{"-x"})
	assertNil(t, err, "Expected nil")
	assertEquals(t, 1, a.Cardinality(), "Expected 1 argument")
	assertEquals(t, true, a.Boolean('x'), "Expected true")
}

func TestSimpleStringPresent(t *testing.T) {
	a, err := NewArgs("x*", []string{"-x", "test"})
	assertNil(t, err, "Expected nil")
	assertEquals(t, 1, a.Cardinality(), "Expected 1 argument")
	assertEquals(t, "test", a.String('x'), "Expected true")
}

func TestMissingStringParameter(t *testing.T) {
	_, err := NewArgs("x*", []string{"-x"})
	assertEquals(t, ErrMissingString, err.(*ArgsError).Code(), "Expected missing string")
	assertEquals(t, 'x', err.(*ArgsError).ArgumentID(), "Expected 'x' as missing")
}

func TestSpacesInFormat(t *testing.T) {
	a, err := NewArgs("x, y", []string{"-xy"})
	assertNil(t, err, "Expecting no errors")
	assertEquals(t, 2, a.Cardinality(), "Expected 2 in cardinality")
	assertEquals(t, true, a.Has('x'), "Expected x in args")
	assertEquals(t, true, a.Has('y'), "Expected y in args")
}

func TestSimpleIntPresent(t *testing.T) {
	a, err := NewArgs("x#", []string{"-x", "1"})
	assertNil(t, err, "Expected no errors")
	assertEquals(t, 1, a.Cardinality(), "Expected 1 argument")
	assertEquals(t, 1, a.Integer('x'), "Expected 1")
}

func TestInvalidIntParameter(t *testing.T) {
	_, err := NewArgs("x#", []string{"-x", "Forty two"})
	assertEquals(t, ErrInvalidInteger, err.(*ArgsError).Code(), "Expected invalid integer")
	assertEquals(t, 'x', err.(*ArgsError).ArgumentID(), "Expected 'x' as invalid")
}

func TestMissingIntParameter(t *testing.T) {
	_, err := NewArgs("x#", []string{"-x"})
	assertEquals(t, ErrMissingInteger, err.(*ArgsError).Code(), "Expected missing integer")
	assertEquals(t, 'x', err.(*ArgsError).ArgumentID(), "Expected 'x' as missing")
}

func TestSimpleFloatPresent(t *testing.T) {
	a, err := NewArgs("x##", []string{"-x", "42.3"})
	assertNil(t, err, "Expected no errors")
	assertEquals(t, 1, a.Cardinality(), "Expected 1 argument")
	assertEquals(t, 42.3, a.Float('x'), "Expected 42.3")
}

func TestInvalidFloatParameter(t *testing.T) {
	_, err := NewArgs("x##", []string{"-x", "Forty two"})
	assertEquals(t, ErrInvalidFloat, err.(*ArgsError).Code(), "Expected invalid double")
	assertEquals(t, 'x', err.(*ArgsError).ArgumentID(), "Expected 'x' as invalid")
}

func TestMissingFloatParameter(t *testing.T) {
	_, err := NewArgs("x##", []string{"-x"})
	assertEquals(t, ErrMissingFloat, err.(*ArgsError).Code(), "Expected missing float")
	assertEquals(t, 'x', err.(*ArgsError).ArgumentID(), "Expected 'x' as missing")
}

func assertEquals(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		t.Fatalf(message)
	}
}

func assertNil(t *testing.T, a interface{}, message string) {
	if a != nil {
		t.Fatalf(message)
	}
}
