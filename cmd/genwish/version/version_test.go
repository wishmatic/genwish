package version

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
	buildversion "github.com/wishmatic/genwish/internal/version"
)

func TestNewCommand(t *testing.T) {
	t.Run("prints the version", func(t *testing.T) {
		out := &bytes.Buffer{}
		cmd := NewCommand(cli.Deps{Out: out})

		if err := testutil.Run(t, cmd); err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if got := strings.TrimSpace(out.String()); got != buildversion.String() {
			t.Errorf("stdout = %q, want %q", got, buildversion.String())
		}
	})
}
