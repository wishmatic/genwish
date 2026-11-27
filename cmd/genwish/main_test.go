package main

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestNewRoot(t *testing.T) {
	t.Run("lists the surviving commands", func(t *testing.T) {
		out, _, err := runCommand(t, cli.Deps{}, "help")
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
		out, _, err := runCommand(t, cli.Deps{}, "help", "login")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out, "genwish login") {
			t.Errorf("help output = %q, want the login usage", out)
		}
	})
}
