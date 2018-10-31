package main

import "testing"

func TestArgs(t *testing.T) {
	NewArgs("l,d*", []string{"-l", "-d", "bladsa"})
}
