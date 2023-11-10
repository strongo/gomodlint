package main

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
)

// Run lint rules on go.mod file content
func Run(data []byte, rules []string) error {

	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	var errs []error
	for _, rule := range rules {
		switch rule {
		case "no-replace":
			if err := validateNoReplace(modFile); err != nil {
				errs = append(errs, err)
			}
		}
	}
	switch len(errs) {
	case 0:
		return nil
	case 1:
		err = errs[0]
	default:
		err = errors.Join(errs...)
	}
	return err
}
