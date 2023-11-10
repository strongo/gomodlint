package main

import (
	"errors"
	"golang.org/x/mod/modfile"
)

func validateNoReplace(modFile *modfile.File) error {
	if len(modFile.Replace) == 0 {
		return nil
	}
	return errors.New("contains `replace` directive(s)")
}
