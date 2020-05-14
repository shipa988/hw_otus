package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("good cmd", func(t *testing.T) {
		input, e := ioutil.TempDir("testdata", "mycmd_")
		defer os.RemoveAll(input)
		if e != nil {
			log.Fatalf("can't create temp directory with error: %s", e)
		}

		testfile, e := ioutil.TempFile(input, "")
		testfile.Close()
		if e != nil {
			log.Fatalf("can't create temp file with error: %s", e)
		}

		var cmdName string
		switch runtime.GOOS {
		case "windows":
			cmd, e := os.Create(filepath.Join(input, "test.cmd"))
			cmdName = cmd.Name()
			if e != nil {
				log.Fatal(e)
			}
			cmd.WriteString("echo off\r\necho|set /p=\"HELLO is (%HELLO%) arguments are %*\">" + testfile.Name() + "\r\nexit 3")
			cmd.Close()
		case "linux":
			cmd, e := os.Create(path.Join(input, "test.sh"))
			cmdName = cmd.Name()
			if e != nil {
				log.Fatal(e)
			}
			cmd.WriteString("#!/usr/bin/env bash\necho -e " + `HELLO is (${HELLO}) arguments are $*">>` + filepath.Base(testfile.Name()) + "\nexit 3")
			cmd.Close()
		default:
			log.Fatal("unknown os")
		}

		m := make(map[string]string)
		m[`HELLO`] = `"hello"`
		env := Environment(m)
		expectedString := `HELLO is ("hello") arguments are 1 test`

		code := RunCmd([]string{cmdName, "1", "test"}, env)

		require.Equal(t, 3, code, "not expected exit code")

		//read test file after cmd execute
		f, e := os.Open(testfile.Name())
		defer f.Close()
		if os.IsNotExist(e) {
			t.Fatal("cmd not execute")
		}
		b, e := ioutil.ReadAll(f)
		if e != nil {
			log.Fatalf("can't open testfile %s; error: %v", testfile.Name(), e)
		}

		require.Equal(t, expectedString, string(b), "not expected environment variable and arguments")
	})
}
