package credential

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/provider"
)

func TestOpen(t *testing.T) {
	t.Run("rejects neo", func(t *testing.T) {
		if _, err := Open(provider.Neo); err == nil {
			t.Error("Open(neo) error = nil, want an error")
		}
	})

	t.Run("rejects an unknown backend", func(t *testing.T) {
		t.Setenv(envBackend, "bogus")

		_, err := Open(provider.Nano)
		if err == nil || !strings.Contains(err.Error(), envBackend) {
			t.Errorf("Open(nano) error = %v, want it to name %s", err, envBackend)
		}
	})

	t.Run("selects the file backend", func(t *testing.T) {
		t.Setenv(envBackend, BackendFile)

		store, err := Open(provider.Nano)
		if err != nil {
			t.Fatalf("Open(nano) error: %v", err)
		}

		if _, ok := store.(*FileStore); !ok {
			t.Errorf("Open(nano) = %T, want *FileStore", store)
		}
	})
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
		if err := NewFileStoreAt(filepath.Join(t.TempDir(), "api_key")).Save(" "); !errors.Is(err, ErrEmptyKey) {
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
