package cli

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/resolve"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestImageDataURL(t *testing.T) {
	png := testutil.ValidPNG(t, 2, 2)

	t.Run("reads a local file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "photo.png")
		if err := os.WriteFile(path, png, 0o600); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		dataURL, err := ImageDataURL(context.Background(), resolve.New(), path)
		if err != nil {
			t.Fatalf("ImageDataURL() error: %v", err)
		}

		want := DataURL(resolve.Image{Data: png, MediaType: resolve.MediaPNG})
		if dataURL != want {
			t.Errorf("ImageDataURL() = %q, want %q", dataURL, want)
		}
	})

	t.Run("reads raw base64", func(t *testing.T) {
		encoded := base64.StdEncoding.EncodeToString(png)

		dataURL, err := ImageDataURL(context.Background(), resolve.New(), encoded)
		if err != nil {
			t.Fatalf("ImageDataURL() error: %v", err)
		}

		want := fmt.Sprintf("data:%s;base64,%s", resolve.MediaPNG, encoded)
		if dataURL != want {
			t.Errorf("ImageDataURL() = %q, want %q", dataURL, want)
		}
	})

	t.Run("rejects a local file that is not an image", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "notes.txt")
		if err := os.WriteFile(path, []byte("not an image"), 0o600); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		_, err := ImageDataURL(context.Background(), resolve.New(), path)
		if err == nil || !strings.Contains(err.Error(), "unsupported image media type") {
			t.Errorf("ImageDataURL() error = %v, want an unsupported media type error", err)
		}
	})
}
