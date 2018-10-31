package main

import "testing"

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

func TestArgsSetsDefaultValueForArgWithoutArgument(t *testing.T) {
	a, err := NewArgs("d*", []string{"-d"})
	if err != nil {
		t.Errorf(err.Error())
	}
	if !a.isValid() {
		t.Errorf("Not valid")
	}

	if "" != a.GetString('d') {
		t.Errorf("No default value set")
	}
}

func TestBla(t *testing.T) {
	a, err := NewArgs("d*", []string{"-d", "testing", "bla"})
	if err != nil {
		t.Errorf(err.Error())
	}
	if !a.isValid() {
		t.Errorf(a.ErrorMessage())
	}

	if "testing" != a.GetString('d') {
		t.Errorf("Invalid string arg")
	}
}
