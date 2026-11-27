package genimage

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestGenImageNeo(t *testing.T) {
	t.Run("requires a model", func(t *testing.T) {
		_, _, err := runNeo(t, cli.Deps{}, "a cat")
		if err == nil || !strings.Contains(err.Error(), "--model is required") {
			t.Errorf("Execute() error = %v, want the missing-model rejection", err)
		}
	})

	t.Run("requires a forge preset", func(t *testing.T) {
		_, _, err := runNeo(t, cli.Deps{}, "a cat", "--model", "x")
		if err == nil || !strings.Contains(err.Error(), "--forge-preset is required") {
			t.Errorf("Execute() error = %v, want the missing-preset rejection", err)
		}
	})
}
