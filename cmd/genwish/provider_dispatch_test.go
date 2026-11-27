package main

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestProviderDispatch(t *testing.T) {
	t.Run("models list requires a provider", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "models", "list")
		if err == nil || !strings.Contains(err.Error(), "between 1 and 2") {
			t.Errorf("Execute() error = %v, want the missing provider reported", err)
		}
	})

	t.Run("models list rejects an unknown provider", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "models", "list", "atlantean")
		if err == nil || !strings.Contains(err.Error(), "is not one of") {
			t.Errorf("Execute() error = %v, want the unknown provider reported", err)
		}
	})

	t.Run("login rejects neo", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "login", "neo")
		if err == nil || !strings.Contains(err.Error(), "does not use an API key") {
			t.Errorf("Execute() error = %v, want neo rejected", err)
		}
	})

	t.Run("logout rejects neo", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "logout", "neo")
		if err == nil || !strings.Contains(err.Error(), "does not use an API key") {
			t.Errorf("Execute() error = %v, want neo rejected", err)
		}
	})

	t.Run("balance rejects an unknown provider", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "balance", "atlantean")
		if err == nil || !strings.Contains(err.Error(), "is not one of") {
			t.Errorf("Execute() error = %v, want the unknown provider reported", err)
		}
	})

	t.Run("gen image lists every provider", func(t *testing.T) {
		out, _, err := runCommand(t, cli.Deps{}, "gen", "image")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{"nano", "neo", "novel"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not mention %q:\n%s", want, out)
			}
		}
	})
}
