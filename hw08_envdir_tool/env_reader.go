package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	newline   = 0x0a
	emptybyte = 0x00
)

type Environment map[string]string

func NewEnvironment() *Environment {
	env := Environment(make(map[string]string))
	return &env
}
func (e *Environment) Add(key, value string) {
	map[string]string(*e)[key] = value
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}
	env := NewEnvironment()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() || strings.ContainsRune(info.Name(), '=') {
			return nil
		}
		key := info.Name()
		value, e := readValue(path)
		if e != nil {
			return e
		}
		env.Add(key, value)
		return nil
	})
	return *env, err
}

//Read envirnment file and parse inner value.
func readValue(path string) (string, error) {
	file, e := os.Open(path)
	defer file.Close() //nolint:staticcheck
	if e != nil {
		return "", e
	}
	r := bufio.NewReader(file)
	value, e := r.ReadBytes(newline)
	if e != nil && e != io.EOF {
		return "", e
	}
	if e != io.EOF { //found new line
		value = value[0 : len(value)-1] //trim delimeter
	}
	value = bytes.ReplaceAll(value, []byte{emptybyte}, []byte{newline}) //replace 00
	value = bytes.TrimRight(value, ` \t`)                               //trim space tab
	return string(value), nil
}
