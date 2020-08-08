package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Environment ...
type Environment map[string]string

func readFirstLine(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	line, _, err := reader.ReadLine()
	if err != nil && err != io.EOF {
		return "", err
	}
	result := string(bytes.Replace(line, []byte{0x00}, []byte{'\n'}, -1))

	return strings.TrimRight(result, "\t "), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return result, err
	}
	for _, f := range files {
		if f.IsDir() || !f.Mode().IsRegular() || strings.Contains(f.Name(), "=") {
			continue
		}
		fd, err := os.OpenFile(path.Join(dir, f.Name()), os.O_RDONLY, 0755)
		if err != nil {
			return result, err
		}
		fl, err := readFirstLine(fd)
		if err != nil {
			return result, err
		}
		result[f.Name()] = fl
	}

	return result, nil
}
