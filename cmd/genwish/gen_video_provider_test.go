package main

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestGenVideoProviders(t *testing.T) {
	t.Run("lists every provider", func(t *testing.T) {
		out, _, err := runCommand(t, cli.Deps{}, "gen", "video")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{"nano", "neo", "novel"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not mention %q:\n%s", want, out)
			}
		}
	})

	for _, name := range []string{"neo", "novel"} {
		t.Run("rejects "+name, func(t *testing.T) {
			_, _, err := runCommand(t, cli.Deps{}, "gen", "video", name, "a wave")
			if err == nil || !strings.Contains(err.Error(), "not supported for provider "+name) {
				t.Errorf("Execute() error = %v, want the unsupported-provider error", err)
			}
		})
	}
}
