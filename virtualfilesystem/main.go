package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func main() {
	fsys := os.DirFS("/")
	path := "bin/ls"
	if err := run(fsys, path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var ErrDifferent = errors.New("different!")

func run(fsys fs.FS, path string) error {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return err
	}
	if string(data) != "foo" {
		return ErrDifferent
	}
	return nil
}
