package main

import (
	"errors"

	"golang.org/x/mod/modfile"
)

type lintError struct {
	error
	code int
}

func (e lintError) Code() int {
	return e.code
}

func validateNoReplace(modFile *modfile.File) error {
	if len(modFile.Replace) == 0 {
		return nil
	}
	return lintError{
		error: errors.New("contains `replace` directive(s)"),
		code:  2,
	}
}
