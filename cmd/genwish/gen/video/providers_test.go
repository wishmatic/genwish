package genvideo

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/testutil"
)

func TestNewNeoCommand(t *testing.T) {
	t.Run("rejects video", func(t *testing.T) {
		err := testutil.Run(t, NewNeoCommand(), "a wave")
		if err == nil || !strings.Contains(err.Error(), "not supported for provider neo") {
			t.Errorf("Execute() error = %v, want the unsupported-provider error", err)
		}
	})
}

func TestNewNovelCommand(t *testing.T) {
	t.Run("rejects video", func(t *testing.T) {
		err := testutil.Run(t, NewNovelCommand(), "a wave")
		if err == nil || !strings.Contains(err.Error(), "not supported for provider novel") {
			t.Errorf("Execute() error = %v, want the unsupported-provider error", err)
		}
	})
}
