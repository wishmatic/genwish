package models

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestNewCommand(t *testing.T) {
	t.Run("list shows help without a provider", func(t *testing.T) {
		out := &bytes.Buffer{}
		cmd := NewCommand(cli.Deps{Out: out})
		cmd.SetOut(out)
		cmd.SetErr(&bytes.Buffer{})

		if err := testutil.Run(t, cmd, "list"); err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out.String(), "Providers: nano, neo, novel") {
			t.Errorf("stdout = %q, want the provider options", out)
		}
	})

	t.Run("list rejects an unknown provider", func(t *testing.T) {
		err := testutil.Run(t, NewCommand(cli.Deps{}), "list", "atlantean")
		if err == nil || !strings.Contains(err.Error(), "is not one of") {
			t.Errorf("Execute() error = %v, want the unknown provider reported", err)
		}
	})

	t.Run("get shows help without an id", func(t *testing.T) {
		out := &bytes.Buffer{}
		cmd := NewCommand(cli.Deps{Out: out})
		cmd.SetOut(out)
		cmd.SetErr(&bytes.Buffer{})

		if err := testutil.Run(t, cmd, "get", "nano"); err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out.String(), "Providers: nano, neo, novel") {
			t.Errorf("stdout = %q, want the provider options", out)
		}
	})

	t.Run("get rejects an unknown provider", func(t *testing.T) {
		err := testutil.Run(t, NewCommand(cli.Deps{}), "get", "atlantean", "model")
		if err == nil || !strings.Contains(err.Error(), "is not one of") {
			t.Errorf("Execute() error = %v, want the unknown provider reported", err)
		}
	})
}
