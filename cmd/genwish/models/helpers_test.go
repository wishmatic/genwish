package models

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func runList(t *testing.T, deps cli.Deps, args ...string) (string, error) {
	t.Helper()

	return runModels(t, NewListCommand, deps, args...)
}

func runGet(t *testing.T, deps cli.Deps, args ...string) (string, error) {
	t.Helper()

	return runModels(t, NewGetCommand, deps, args...)
}

func runModels(
	t *testing.T, newCommand func(cli.Deps) *cobra.Command, deps cli.Deps, args ...string,
) (string, error) {
	t.Helper()

	testutil.SetHome(t)

	out := &bytes.Buffer{}
	deps.Out = out
	if deps.Err == nil {
		deps.Err = &bytes.Buffer{}
	}

	err := testutil.Run(t, newCommand(deps), args...)

	return out.String(), err
}
