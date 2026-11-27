package models

import (
	"bytes"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func runList(t *testing.T, deps cli.Deps, args ...string) (string, error) {
	t.Helper()

	testutil.SetHome(t)

	out := &bytes.Buffer{}
	deps.Out = out
	if deps.Err == nil {
		deps.Err = &bytes.Buffer{}
	}

	err := testutil.Run(t, NewListCommand(deps), args...)

	return out.String(), err
}
