package edit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestCommandsShowHelpWithoutAnImage(t *testing.T) {
	commands := map[string]func(cli.Deps) *cobra.Command{
		"trim":    NewTrimCommand,
		"square":  NewSquareCommand,
		"circle":  NewCircleCommand,
		"convert": NewConvertCommand,
		"bgkill":  NewBgkillCommand,
	}

	for name, newCommand := range commands {
		t.Run(name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := newCommand(cli.Deps{})
			cmd.SetOut(out)
			cmd.SetErr(&bytes.Buffer{})

			if err := testutil.Run(t, cmd); err != nil {
				t.Fatalf("Execute() error: %v", err)
			}

			if !strings.Contains(out.String(), "<image>") {
				t.Errorf("stdout = %q, want the image usage", out)
			}
		})
	}
}
