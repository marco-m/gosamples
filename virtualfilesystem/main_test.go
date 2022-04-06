package main

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestRun(t *testing.T) {
	fsys := fstest.MapFS{
		"hello": {
			Data: []byte("hello, world"),
		},
		"ciao": {
			Data: []byte("foo"),
		},
	}

	{
		want := &fs.PathError{}
		err := run(fsys, "pippo")
		if !errors.As(err, &want) {
			t.Errorf("have: %s; want: PathError", err)
		}
	}

	{
		want := ErrDifferent
		err := run(fsys, "hello")
		if !errors.Is(err, want) {
			t.Errorf("have: %s; want: %s", err, want)
		}
	}

	if err := run(fsys, "ciao"); err != nil {
		t.Errorf("have: %s; want: <no error>", err)
	}
}
