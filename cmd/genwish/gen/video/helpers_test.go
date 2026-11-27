package genvideo

import (
	"bytes"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func runNano(t *testing.T, deps cli.Deps, args ...string) (string, string, error) {
	t.Helper()

	testutil.SetHome(t)

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	deps.Out = out
	deps.Err = errOut

	err := testutil.Run(t, NewNanoCommand(deps), args...)

	return out.String(), errOut.String(), err
}
