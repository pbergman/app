package app

import (
	"fmt"
	"io"
	"os"
)

type Error struct {
	Err     error
	Code	int8
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) Exit(w io.Writer) {
	fmt.Fprintln(w, e.Err.Error())
	os.Exit(int(e.Code))
}

