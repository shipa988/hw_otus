package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("good environment", func(t *testing.T) {
		//creation of my test files in code because the files in the "env" directory can be modified
		input, e := ioutil.TempDir("testdata", "myenv_")
		if e != nil {
			log.Fatalf("can't create temp directory with error: %s", e)
		}
		defer os.RemoveAll(input)

		BAR, e := os.Create(filepath.Join(input, `BAR`))
		defer func() {
			BAR.Close()
			os.Remove(BAR.Name())
		}()
		if e != nil {
			log.Fatal(e)
		}
		BAR.WriteString("bar\nPLEASE IGNORE SECOND LINE\n")

		FOO, e := os.Create(filepath.Join(input, `FOO`))
		defer func() {
			FOO.Close()
			os.Remove(FOO.Name())
		}()
		if e != nil {
			log.Fatal(e)
		}
		FOO.Write([]byte{0x20, 0x20, 0x20, 0x66, 0x6f, 0x6f, 0x00, 0x77, 0x69, 0x74, 0x68, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x6c, 0x69, 0x6e, 0x65}) //   foo\x00with new line

		HELLO, e := os.Create(filepath.Join(input, `HELLO`))
		defer func() {
			HELLO.Close()
			os.Remove(HELLO.Name())
		}()
		if e != nil {
			log.Fatal(e)
		}
		HELLO.WriteString(`"hello"`)

		UNSET, e := os.Create(filepath.Join(input, `UNSET`))
		defer func() {
			UNSET.Close()
			os.Remove(UNSET.Name())
		}()
		if e != nil {
			log.Fatal(e)
		}

		m := make(map[string]string)
		m[`BAR`] = `bar`
		m[`FOO`] = `   foo
with new line`
		m[`HELLO`] = `"hello"`
		m[`UNSET`] = ``

		expectedEnv := Environment(m)
		actualEnv, e := ReadDir(input)

		require.NoError(t, e)
		require.Equal(t, expectedEnv, actualEnv, "environment not equal")
	})
	t.Run("not exist environment dir", func(t *testing.T) {
		//env directory not exist
		input := path.Join("testdata", "myenv_123")
		actualEnv, e := ReadDir(input)

		require.Error(t, e)
		require.Equal(t, Environment(nil), actualEnv, "environment not equal")
	})
}
