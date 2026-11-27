package genvideo

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestNewCommand(t *testing.T) {
	t.Run("lists every provider", func(t *testing.T) {
		out := &bytes.Buffer{}
		cmd := NewCommand(cli.Deps{})
		cmd.SetOut(out)
		cmd.SetErr(&bytes.Buffer{})

		if err := testutil.Run(t, cmd); err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{"nano", "neo", "novel"} {
			if !strings.Contains(out.String(), want) {
				t.Errorf("stdout does not mention %q:\n%s", want, out)
			}
		}
	})
}

func TestCommandsShowHelpWithoutAPrompt(t *testing.T) {
	commands := map[string]func() *cobra.Command{
		"nano":  func() *cobra.Command { return NewNanoCommand(cli.Deps{}) },
		"neo":   NewNeoCommand,
		"novel": NewNovelCommand,
	}

	for name, newCommand := range commands {
		t.Run(name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := newCommand()
			cmd.SetOut(out)
			cmd.SetErr(&bytes.Buffer{})

			if err := testutil.Run(t, cmd); err != nil {
				t.Fatalf("Execute() error: %v", err)
			}

			if !strings.Contains(out.String(), "<prompt...>") {
				t.Errorf("stdout = %q, want the prompt usage", out)
			}
		})
	}
}
