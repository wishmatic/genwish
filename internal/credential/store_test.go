package credential

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/wishmatic/genwish/internal/provider"
)

func TestOpen(t *testing.T) {
	t.Run("rejects neo", func(t *testing.T) {
		if _, err := Open(provider.Neo); err == nil {
			t.Error("Open(neo) error = nil, want an error")
		}
	})

	t.Run("returns the auto store", func(t *testing.T) {
		store, err := Open(provider.Nano)
		if err != nil {
			t.Fatalf("Open(nano) error: %v", err)
		}

		if _, ok := store.(*autoStore); !ok {
			t.Errorf("Open(nano) = %T, want *autoStore", store)
		}
	})
}

func TestStorageWarning(t *testing.T) {
	tests := []struct {
		name       string
		keyringErr error
		path       string
		want       string
	}{
		{name: "silent on a healthy keyring", keyringErr: nil, path: "/tmp/key", want: ""},
		{
			name:       "silent when no key is stored",
			keyringErr: keyring.ErrNotFound,
			path:       "/tmp/key",
			want:       "",
		},
		{
			name:       "names the fallback path",
			keyringErr: errors.New("dbus unavailable"),
			path:       "/tmp/key",
			want:       "/tmp/key",
		},
		{
			name:       "falls back to a generic message",
			keyringErr: errors.New("dbus unavailable"),
			want:       "user-only file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := storageWarning(tt.keyringErr, tt.path)

			if tt.want == "" {
				if got != "" {
					t.Errorf("storageWarning() = %q, want silence", got)
				}

				return
			}

			if !strings.Contains(got, tt.want) {
				t.Errorf("storageWarning() = %q, want it to contain %q", got, tt.want)
			}
		})
	}
}

func TestFileStore(t *testing.T) {
	t.Run("round trips a key", func(t *testing.T) {
		store := NewFileStoreAt(filepath.Join(t.TempDir(), "api_key"))

		if err := store.Delete(); err != nil {
			t.Fatalf("Delete() before Save error: %v", err)
		}

		if _, err := store.Load(); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Load() before Save error = %v, want ErrNotFound", err)
		}

		if err := store.Save("  file-key  "); err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		key, err := store.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if key != "file-key" {
			t.Errorf("Load() = %q, want file-key", key)
		}

		if err := store.Delete(); err != nil {
			t.Fatalf("Delete() error: %v", err)
		}

		if _, err := store.Load(); !errors.Is(err, ErrNotFound) {
			t.Errorf("Load() after Delete error = %v, want ErrNotFound", err)
		}
	})

	t.Run("rejects an empty key", func(t *testing.T) {
		if err := NewFileStoreAt(filepath.Join(t.TempDir(), "api_key")).Save(" "); !errors.Is(
			err,
			ErrEmptyKey,
		) {
			t.Errorf("Save() error = %v, want ErrEmptyKey", err)
		}
	})

	t.Run("keeps providers apart", func(t *testing.T) {
		paths := make(map[string]string)

		for _, name := range []provider.Name{provider.Nano, provider.Novel, provider.Danbooru} {
			store, err := NewFileStore(name)
			if err != nil {
				t.Fatalf("NewFileStore(%s) error: %v", name, err)
			}

			if other, ok := paths[store.path]; ok {
				t.Errorf("%s and %s share the path %q", other, name, store.path)
			}

			paths[store.path] = name.String()
		}
	})

	t.Run("rejects neo", func(t *testing.T) {
		if _, err := NewFileStore(provider.Neo); err == nil {
			t.Error("NewFileStore(neo) error = nil, want an error")
		}
	})
}
