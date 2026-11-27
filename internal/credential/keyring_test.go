package credential

import (
	"errors"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/wishmatic/genwish/internal/provider"
)

func TestKeyringStore(t *testing.T) {
	t.Run("round trips a nano key", func(t *testing.T) {
		keyring.MockInit()

		store := KeyringStore{provider: provider.Nano}

		if err := store.Delete(); err != nil {
			t.Fatalf("Delete() before Save error: %v", err)
		}

		if _, err := store.Load(); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Load() before Save error = %v, want ErrNotFound", err)
		}

		if err := store.Save("  keyring-key  "); err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		key, err := store.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if key != "keyring-key" {
			t.Errorf("Load() = %q, want keyring-key", key)
		}
	})

	t.Run("keeps providers apart", func(t *testing.T) {
		keyring.MockInit()

		nano := KeyringStore{provider: provider.Nano}
		novel := KeyringStore{provider: provider.Novel}

		if err := nano.Save("nano-key"); err != nil {
			t.Fatalf("nano Save() error: %v", err)
		}

		if err := novel.Save("novel-key"); err != nil {
			t.Fatalf("novel Save() error: %v", err)
		}

		nanoKey, err := nano.Load()
		if err != nil {
			t.Fatalf("nano Load() error: %v", err)
		}

		novelKey, err := novel.Load()
		if err != nil {
			t.Fatalf("novel Load() error: %v", err)
		}

		if nanoKey != "nano-key" || novelKey != "novel-key" {
			t.Errorf("keys = %q, %q, want them kept apart", nanoKey, novelKey)
		}
	})

	t.Run("rejects an empty key", func(t *testing.T) {
		keyring.MockInit()

		if err := (KeyringStore{provider: provider.Nano}).Save(" "); !errors.Is(err, ErrEmptyKey) {
			t.Errorf("Save() error = %v, want ErrEmptyKey", err)
		}
	})

	t.Run("reports an unavailable keyring", func(t *testing.T) {
		keyring.MockInitWithError(errors.New("keyring unavailable"))

		if _, err := (KeyringStore{provider: provider.Nano}).Load(); err == nil {
			t.Error("Load() error = nil, want the keyring error")
		}
	})

	t.Run("rejects neo", func(t *testing.T) {
		keyring.MockInit()

		if _, err := (KeyringStore{provider: provider.Neo}).Load(); err == nil {
			t.Error("Load() error = nil, want an error")
		}
	})
}
