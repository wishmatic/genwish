package genimage

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/credential"
)

func TestGenImageNovel(t *testing.T) {
	t.Run("reports a missing key", func(t *testing.T) {
		deps := cli.Deps{
			Credentials: credential.NewFileStoreAt(filepath.Join(t.TempDir(), "key")),
		}

		_, _, err := runNovel(t, deps, "a cat")
		if err == nil || !strings.Contains(err.Error(), "genwish login novel") {
			t.Errorf("Execute() error = %v, want a login hint", err)
		}
	})

	t.Run("rejects an empty model", func(t *testing.T) {
		_, _, err := runNovel(t, cli.Deps{}, "a cat", "--model", " ")
		if err == nil || !strings.Contains(err.Error(), "--model must not be empty") {
			t.Errorf("Execute() error = %v, want the empty-model rejection", err)
		}
	})
}
