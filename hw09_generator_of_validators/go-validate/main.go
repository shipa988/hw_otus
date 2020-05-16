package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Please enter the .go filename for go-validate")
	}
	fname := args[0]
	err := GenValidate(fname)
	if err != nil {
		log.Printf("Error occurred:%v", err)
	}
}
