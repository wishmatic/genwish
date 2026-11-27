package genimage

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func runImage(
	t *testing.T, newCommand func(cli.Deps) *cobra.Command, deps cli.Deps, args ...string,
) (string, string, error) {
	t.Helper()

	testutil.SetHome(t)

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	deps.Out = out
	deps.Err = errOut

	err := testutil.Run(t, newCommand(deps), args...)

	return out.String(), errOut.String(), err
}

func runNano(t *testing.T, deps cli.Deps, args ...string) (string, string, error) {
	t.Helper()

	return runImage(t, NewNanoCommand, deps, args...)
}

func runNeo(t *testing.T, deps cli.Deps, args ...string) (string, string, error) {
	t.Helper()

	return runImage(t, NewNeoCommand, deps, args...)
}

func runNovel(t *testing.T, deps cli.Deps, args ...string) (string, string, error) {
	t.Helper()

	return runImage(t, NewNovelCommand, deps, args...)
}
