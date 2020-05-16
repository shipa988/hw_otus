package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const tempD = "test"

func TestGenCode(t *testing.T) {
	t.Run("not exist .go filename", func(t *testing.T) {
		e := GenValidate("fileNotExist.go")
		require.NotNil(t, e)
		require.NoFileExists(t, "fileNotExist_validation.go", "*_validation.go file should not be generate")
	})
	t.Run("simple test", func(t *testing.T) {
		newTempDir(tempD)
		defer os.RemoveAll(tempD)

		fname := filepath.Join(tempD, "test.go")
		fnameValitator := filepath.Join(tempD, "test_validation_generated.go")
		f, e := os.Create(fname)
		defer f.Close()
		if e != nil {
			log.Fatalf("can't create temp file with error: %s", e)
		}
		f.WriteString(`package test

		type Book struct {
			Title string    ` + "`" + `validate:"len:100"` + "`" + `
			PageCount int ` + "`" + `validate:"min:10|max:100"` + "`" + `
		}
		`)
		e = GenValidate(fname)
		require.Nil(t, e)
		require.FileExists(t, fnameValitator, "*_validation_generated.go file will not generate")

		f, _ = os.Open(fnameValitator)
		defer f.Close()
		b, _ := ioutil.ReadAll(f)
		require.True(t, strings.Contains(string(b), "CODE GENERATED AUTOMATICALLY"), "disclaimer is absent")
	})
	t.Run("bad tag values int len:one hundred", func(t *testing.T) {
		newTempDir(tempD)
		defer os.RemoveAll(tempD)

		fname := filepath.Join(tempD, "test.go")
		fnameValitator := filepath.Join(tempD, "test_validation_generated.go")
		f, e := os.Create(fname)
		defer f.Close()
		if e != nil {
			log.Fatalf("can't create temp file with error: %s", e)
		}
		f.WriteString(`package test

		type Book struct {
			Title string    ` + "`" + `validate:"len:one hundred"` + "`" + `
			PageCount int ` + "`" + `validate:"min:10|max:100"` + "`" + `
		}
		`)
		e = GenValidate(fname)
		require.NotNil(t, e)
		require.FileExists(t, fnameValitator, "*_validation_generated.go file will not generate")

	})
}

//newTempDir create tempDir returns it name,error.
func newTempDir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("unable to create temp dir fo test. Error: %v", err)
	}
}
