package main

import (
	"flag"
	"log"
	"os"
)

const (
	minargslen = 2
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < minargslen {
		log.Fatal("should be at least 2 arguments: dir to environment and execute command")
	}
	envdir := args[0]
	cmd := args[1:]
	env, err := ReadDir(envdir)
	if err != nil {
		log.Fatalf("can't read environment dir %s", envdir)
	}
	os.Exit(RunCmd(cmd, env))
}
