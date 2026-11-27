package models

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestModelsGetNovel(t *testing.T) {
	t.Run("prints the static model", func(t *testing.T) {
		out, err := runGet(t, cli.Deps{}, "novel", "nai-diffusion-5-full")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if out != "id: nai-diffusion-5-full\nkind: image\n" {
			t.Errorf("stdout = %q, want the image model", out)
		}
	})

	t.Run("reports an unserved model", func(t *testing.T) {
		_, err := runGet(t, cli.Deps{}, "novel", "absent")
		if err == nil || !strings.Contains(err.Error(), "is not served by novel") {
			t.Errorf("Execute() error = %v, want the unserved model", err)
		}
	})

	t.Run("rejects video", func(t *testing.T) {
		_, err := runGet(t, cli.Deps{}, "novel", "nai-diffusion-5-full", "video")
		if err == nil || !strings.Contains(err.Error(), "does not serve video") {
			t.Errorf("Execute() error = %v, want a video rejection", err)
		}
	})
}
