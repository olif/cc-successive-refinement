package main

import (
	"log"
	"os"
)

func main() {
	log.Printf("%v", os.Args)
	a, err := NewArgs("l,d*", os.Args)
	if err == nil {
		boolArg := a.GetBoolean('l')
		stringArg := a.GetString('d')
		log.Printf("String arg: %s", stringArg)
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
