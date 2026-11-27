package edit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestNewConvertCommand(t *testing.T) {
	t.Run("requires a format", func(t *testing.T) {
		deps := cli.Deps{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}}

		err := testutil.Run(t, NewConvertCommand(deps), "image.png")
		if err == nil || !strings.Contains(err.Error(), "--format is required") {
			t.Errorf("Execute() error = %v, want the missing-format rejection", err)
		}
	})
}
