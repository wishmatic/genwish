package models

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestNewCommand(t *testing.T) {
	t.Run("list requires a provider", func(t *testing.T) {
		err := testutil.Run(t, NewCommand(cli.Deps{}), "list")
		if err == nil || !strings.Contains(err.Error(), "between 1 and 2") {
			t.Errorf("Execute() error = %v, want the missing provider reported", err)
		}
	})

	t.Run("list rejects an unknown provider", func(t *testing.T) {
		err := testutil.Run(t, NewCommand(cli.Deps{}), "list", "atlantean")
		if err == nil || !strings.Contains(err.Error(), "is not one of") {
			t.Errorf("Execute() error = %v, want the unknown provider reported", err)
		}
	})

	t.Run("get requires a provider and an id", func(t *testing.T) {
		err := testutil.Run(t, NewCommand(cli.Deps{}), "get")
		if err == nil || !strings.Contains(err.Error(), "between 2 and 3") {
			t.Errorf("Execute() error = %v, want the missing arguments reported", err)
		}
	})

	t.Run("get rejects an unknown provider", func(t *testing.T) {
		err := testutil.Run(t, NewCommand(cli.Deps{}), "get", "atlantean", "model")
		if err == nil || !strings.Contains(err.Error(), "is not one of") {
			t.Errorf("Execute() error = %v, want the unknown provider reported", err)
		}
	})
}
