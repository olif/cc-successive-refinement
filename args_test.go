package main

import (
	"testing"
)

func TestCanParseBoolArg(t *testing.T) {
	a, err := NewArgs("l", []string{"-l"})
	if err != nil {
		t.Fail()
	}

	if ok := a.GetBoolean('l'); !ok {
		t.Fail()
	}
}

func TestCanParseStringArg(t *testing.T) {
	a, err := NewArgs("d*", []string{"-d", "testing"})
	if err != nil {
		t.Errorf("Could not parse args")
	}
	if "testing" != a.GetString('d') {
		t.Errorf("Incorrect string value")
	}
}

func TestCanParseInteger(t *testing.T) {
	_, err := NewArgs("d#", []string{"-d", "1"})
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCanParseBothStringAndBoolArg(t *testing.T) {
	a, err := NewArgs("l,d*", []string{"-l", "-d", "testing"})
	if err != nil {
		t.Errorf("Could not parse args")
	}

	if ok := a.GetBoolean('l'); !ok {
		t.Errorf("Invalid boolean arg")
	}

	if "testing" != a.GetString('d') {
		t.Errorf("Invalid string arg")
	}
}

func TestDoublePresent(t *testing.T) {
	a, err := NewArgs("x##", []string{"-x", "42.3"})
	if err != nil {
		t.Errorf(err.Error())
	}

	if !a.Has('x') {
		t.Errorf("Element not found")
	}

	if a.Cardinality() != 1 {
		t.Errorf("Invalid number of arguments found")
	}

	if 42.3 != a.GetDouble('x') {
		t.Errorf("Value is invalid")
	}
}

func TestMissingDouble(t *testing.T) {
	_, err := NewArgs("x##", []string{"-x"})
	if err != nil {
		if err.Error() != "Could not find double parameter for -x" {
			t.Errorf("Got: %s", err.Error())
		}
	} else {
		t.Fail()
	}
}

func TestBla(t *testing.T) {
	a, err := NewArgs("d*", []string{"-d", "testing", "bla"})
	if err != nil {
		t.Errorf(err.Error())
	}

	if "testing" != a.GetString('d') {
		t.Errorf("Invalid string arg")
	}
}
