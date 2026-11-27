package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func runRoot(t *testing.T, deps cli.Deps, args ...string) (string, string, error) {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	deps.Out = out
	deps.Err = errOut

	root := newRoot(deps)
	root.SetArgs(args)

	err := root.Execute()

	return out.String(), errOut.String(), err
}

func TestNewRoot(t *testing.T) {
	t.Run("lists the surviving commands", func(t *testing.T) {
		out, _, err := runRoot(t, cli.Deps{}, "help")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, command := range []string{"gen", "edit", "models", "balance", "login", "logout", "util"} {
			if !strings.Contains(out, command) {
				t.Errorf("help output does not mention %q:\n%s", command, out)
			}
		}

		for _, removed := range []string{"  serve ", "  web ", "  mcp "} {
			if strings.Contains(out, removed) {
				t.Errorf("help output still mentions removed command %q:\n%s", removed, out)
			}
		}
	})

	t.Run("shows a subcommand usage", func(t *testing.T) {
		out, _, err := runRoot(t, cli.Deps{}, "help", "login")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out, "genwish login") {
			t.Errorf("help output = %q, want the login usage", out)
		}
	})
}
