package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Println("cmd args is empty")
		return 1
	}

	name := cmd[0]
	command := exec.Command(name, cmd[1:]...)

	for s, t := range map[string]string(env) {
		os.Unsetenv(s)
		if t != "" { //else unset envar and not append
			os.Setenv(s, t) // append envar
		}
	}

	stdin, err := command.StdinPipe()
	if err != nil {
		log.Println(err)
	}
	go func() {
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err = command.Start(); err != nil {
		log.Println(err)
	}
	if err = command.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		} else {
			log.Println(err)
		}
	}

	return returnCode
}
