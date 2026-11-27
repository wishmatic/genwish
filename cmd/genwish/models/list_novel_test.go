package models

import (
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestModelsListNovel(t *testing.T) {
	t.Run("lists the static models", func(t *testing.T) {
		out, err := runList(t, cli.Deps{}, "novel")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{"NOVEL MODELS", "nai-diffusion-5-full", "17 of 17 models shown"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}
	})

	t.Run("filters and limits", func(t *testing.T) {
		out, err := runList(t, cli.Deps{}, "novel", "--query", "4-5", "--limit", "2")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out, "nai-diffusion-4-5-full") || !strings.Contains(out, "2 of 17 models shown (4 matched)") {
			t.Errorf("stdout = %q, want the filtered window", out)
		}
	})

	t.Run("rejects video", func(t *testing.T) {
		_, err := runList(t, cli.Deps{}, "novel", "video")
		if err == nil || !strings.Contains(err.Error(), "does not serve video") {
			t.Errorf("Execute() error = %v, want a video rejection", err)
		}
	})
}
