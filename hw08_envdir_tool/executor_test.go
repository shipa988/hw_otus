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
	t.Run("bad cmd", func(t *testing.T) {
		code := RunCmd([]string{}, Environment{})
		require.Equalf(t, 1, code, "should be error code 1")
	})
	t.Run("good cmd", func(t *testing.T) {
		tempdir, e := ioutil.TempDir("testdata", "mycmd_")
		defer os.RemoveAll(tempdir)

		testfile := newTempFile(tempdir, "echofile")

		var cmdName string
		switch runtime.GOOS {
		case "windows":
			cmd, e := os.Create(filepath.Join(tempdir, "test.cmd"))
			cmdName = cmd.Name()
			if e != nil {
				log.Fatal(e)
			}
			cmd.WriteString("echo off\r\necho|set /p=\"HELLO is %HELLO% arguments are %*\">" + testfile + "\r\nexit 3")
			cmd.Close()
		case "linux", "darwin":
			cmd, e := os.Create(path.Join(tempdir, "test.sh"))
			cmdName = cmd.Name()
			if e != nil {
				log.Fatal(e)
			}
			cmd.WriteString("#!/usr/bin/env bash\necho -e -n \"HELLO is ${HELLO} arguments are $*\">>" + testfile + " && exit 3")
			if e = cmd.Chmod(0777); e != nil {
				log.Fatal(e)
				return
			}
			cmd.Close()
		default:
			log.Fatal("unknown os")
		}

		m := make(map[string]string)
		m[`HELLO`] = `"hello"`
		env := Environment(m)
		expectedString := `HELLO is "hello" arguments are 1 test`

		code := RunCmd([]string{cmdName, "1", "test"}, env)

		require.Equal(t, 3, code, "not expected exit code")

		//read test file after cmd execute
		f, e := os.Open(testfile)
		defer f.Close()
		if os.IsNotExist(e) {
			t.Fatal("cmd not execute")
		}
		b, e := ioutil.ReadAll(f)
		if e != nil {
			log.Fatalf("can't open testfile %s; error: %v", testfile, e)
		}

		require.Equal(t, expectedString, string(b), "not expected environment variable and arguments")
	})
}

//Create tempfile in "dir" with "templ" template and return filename of tempfile
func newTempFile(dir, templ string) string {
	tempFile, err := ioutil.TempFile(dir, "*_"+templ)
	if err != nil {
		log.Fatalf("unable to create temp file fo test. Error: %v", err)
	}
	defer tempFile.Close()
	return tempFile.Name()
}
