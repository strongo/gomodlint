package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/mod/modfile"
)

const validGoMod = `module example.com/test

go 1.22
`

const replaceGoMod = `module example.com/test

go 1.22

replace example.com/foo => ../foo
`

func TestRun(t *testing.T) {
	t.Run("valid_no_replace", func(t *testing.T) {
		err := Run([]byte(validGoMod), []string{"no-replace"})
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("unhandled_rule", func(t *testing.T) {
		err := Run([]byte(validGoMod), []string{"unknown-rule"})
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("replace_detected_single", func(t *testing.T) {
		err := Run([]byte(replaceGoMod), []string{"no-replace"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "contains `replace` directive(s)") {
			t.Errorf("unexpected error message: %v", err)
		}
		if coder, ok := err.(interface{ Code() int }); !ok || coder.Code() != 2 {
			t.Errorf("expected exit code 2, got %v (ok=%v)", coder, ok)
		}
	})

	t.Run("replace_detected_multiple_errs", func(t *testing.T) {
		err := Run([]byte(replaceGoMod), []string{"no-replace", "no-replace"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid_go_mod_syntax", func(t *testing.T) {
		err := Run([]byte("this is not valid go.mod syntax :::"), []string{"no-replace"})
		if err == nil {
			t.Fatal("expected parse error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to parse") {
			t.Errorf("expected 'failed to parse', got %v", err)
		}
	})
}

func TestValidateNoReplace(t *testing.T) {
	t.Run("no_replaces", func(t *testing.T) {
		f, err := modfile.Parse("go.mod", []byte(validGoMod), nil)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if err := validateNoReplace(f); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("with_replaces", func(t *testing.T) {
		f, err := modfile.Parse("go.mod", []byte(replaceGoMod), nil)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		err = validateNoReplace(f)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		coder, ok := err.(interface{ Code() int })
		if !ok {
			t.Fatal("expected error to implement Code() int")
		}
		if coder.Code() != 2 {
			t.Errorf("expected code 2, got %d", coder.Code())
		}
	})
}

func TestMainFunc(t *testing.T) {
	origArgs := os.Args
	origExit := osExit
	defer func() {
		os.Args = origArgs
		osExit = origExit
	}()

	var lastExitCode int
	osExit = func(code int) {
		lastExitCode = code
	}

	t.Run("valid_mod_specified_dir", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(validGoMod), 0600); err != nil {
			t.Fatal(err)
		}

		lastExitCode = -1
		os.Args = []string{"gomodlint", tmpDir}
		main()
		if lastExitCode != -1 {
			t.Errorf("expected no exit called, got %d", lastExitCode)
		}
	})

	t.Run("valid_mod_default_dir", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(validGoMod), 0600); err != nil {
			t.Fatal(err)
		}

		origWd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.Chdir(origWd) }()

		lastExitCode = -1
		os.Args = []string{"gomodlint"}
		main()
		if lastExitCode != -1 {
			t.Errorf("expected no exit called, got %d", lastExitCode)
		}
	})

	t.Run("invalid_mod_with_replace", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(replaceGoMod), 0600); err != nil {
			t.Fatal(err)
		}

		lastExitCode = -1
		os.Args = []string{"gomodlint", tmpDir}
		main()
		if lastExitCode != 2 {
			t.Errorf("expected exit code 2, got %d", lastExitCode)
		}
	})

	t.Run("file_read_error", func(t *testing.T) {
		lastExitCode = -1
		os.Args = []string{"gomodlint", filepath.Join(t.TempDir(), "nonexistent")}
		main()
		if lastExitCode != 1 {
			t.Errorf("expected exit code 1, got %d", lastExitCode)
		}
	})
}
