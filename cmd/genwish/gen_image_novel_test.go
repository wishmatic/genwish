package main

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestGenImageNovel(t *testing.T) {
	t.Run("reports a missing key", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{Credentials: &memoryStore{}}, "gen", "image", "novel", "a cat")
		if err == nil || !strings.Contains(err.Error(), "genwish login novel") {
			t.Errorf("Execute() error = %v, want a login hint", err)
		}
	})

	t.Run("rejects an empty model", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "gen", "image", "novel", "a cat", "--model", " ")
		if err == nil || !strings.Contains(err.Error(), "--model must not be empty") {
			t.Errorf("Execute() error = %v, want the empty-model rejection", err)
		}
	})
}
