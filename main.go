package main

import (
	"log"
	"os"
)

func main() {
	a := NewArgs("l", os.Args)
	if a.isValid() {
		boolArg := a.GetBoolean('l')
		if boolArg {
			log.Println("Bool on")
		} else {
			log.Println("Bool off")
		}
	} else {
		log.Println("Could not parse args")
		log.Println(a.ErrorMessage())
	}
}
