package main

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestGenImageNeo(t *testing.T) {
	t.Run("requires a model", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "gen", "image", "neo", "a cat")
		if err == nil || !strings.Contains(err.Error(), "--model is required") {
			t.Errorf("Execute() error = %v, want the missing-model rejection", err)
		}
	})

	t.Run("requires a forge preset", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "gen", "image", "neo", "a cat", "--model", "x")
		if err == nil || !strings.Contains(err.Error(), "--forge-preset is required") {
			t.Errorf("Execute() error = %v, want the missing-preset rejection", err)
		}
	})
}

func TestEditConvertRequiresAFormat(t *testing.T) {
	_, _, err := runCommand(t, cli.Deps{}, "edit", "convert", "image.png")
	if err == nil || !strings.Contains(err.Error(), "--format is required") {
		t.Errorf("Execute() error = %v, want the missing-format rejection", err)
	}
}
