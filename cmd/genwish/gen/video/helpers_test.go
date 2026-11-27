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

	cmd := NewNanoCommand(deps)
	cmd.SetOut(out)
	cmd.SetErr(errOut)

	err := testutil.Run(t, cmd, args...)

	return out.String(), errOut.String(), err
}
