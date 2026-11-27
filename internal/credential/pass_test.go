package credential

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/provider"
)

type passCall struct {
	args  []string
	stdin string
}

type fakePass struct {
	response func() ([]byte, []byte, error)
	calls    []passCall
}

func (f *fakePass) run(args []string, stdin io.Reader) ([]byte, []byte, error) {
	call := passCall{args: slices.Clone(args)}

	if stdin != nil {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return nil, nil, err
		}

		call.stdin = string(data)
	}

	f.calls = append(f.calls, call)

	if f.response == nil {
		return nil, nil, nil
	}

	return f.response()
}

func testPassStore(t *testing.T, runner *fakePass) (*PassStore, string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "genwish", "nanogpt-api-key.gpg")

	return &PassStore{
		entry: "genwish/nanogpt-api-key",
		path:  path,
		run:   runner.run,
	}, path
}

func writePassEntry(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}

	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}
}

func TestPassStore(t *testing.T) {
	t.Run("round trips a key", func(t *testing.T) {
		runner := &fakePass{
			response: func() ([]byte, []byte, error) {
				return []byte("  pass-key  \n"), nil, nil
			},
		}
		store, path := testPassStore(t, runner)

		if err := store.Save("pass-key"); err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		if !slices.Equal(runner.calls[0].args,
			[]string{"insert", "-m", "-f", "genwish/nanogpt-api-key"},
		) {
			t.Errorf("Save() args = %v, want the insert flags", runner.calls[0].args)
		}

		if runner.calls[0].stdin != "pass-key\n" {
			t.Errorf("Save() stdin = %q, want the key", runner.calls[0].stdin)
		}

		writePassEntry(t, path)

		key, err := store.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if key != "pass-key" {
			t.Errorf("Load() = %q, want pass-key", key)
		}

		if !slices.Equal(runner.calls[1].args,
			[]string{"show", "genwish/nanogpt-api-key"},
		) {
			t.Errorf("Load() args = %v, want the show command", runner.calls[1].args)
		}
	})

	t.Run("returns not found when the entry is missing", func(t *testing.T) {
		runner := &fakePass{}
		store, _ := testPassStore(t, runner)

		if _, err := store.Load(); !errors.Is(err, ErrNotFound) {
			t.Errorf("Load() error = %v, want ErrNotFound", err)
		}

		if len(runner.calls) != 0 {
			t.Errorf("calls = %v, want none for a missing entry", runner.calls)
		}
	})

	t.Run("reports a pass failure", func(t *testing.T) {
		runner := &fakePass{
			response: func() ([]byte, []byte, error) {
				return nil, []byte("gpg: decryption failed"), errors.New("exit status 2")
			},
		}
		store, path := testPassStore(t, runner)
		writePassEntry(t, path)

		_, err := store.Load()
		if err == nil || !strings.Contains(err.Error(), "read pass entry") ||
			!strings.Contains(err.Error(), "gpg: decryption failed") {
			t.Errorf("Load() error = %v, want the pass failure", err)
		}
	})

	t.Run("deletes an existing entry", func(t *testing.T) {
		runner := &fakePass{}
		store, path := testPassStore(t, runner)
		writePassEntry(t, path)

		if err := store.Delete(); err != nil {
			t.Fatalf("Delete() error: %v", err)
		}

		want := []string{"rm", "-f", "genwish/nanogpt-api-key"}
		if !slices.Equal(runner.calls[0].args, want) {
			t.Errorf("Delete() args = %v, want %v", runner.calls[0].args, want)
		}
	})

	t.Run("keeps delete quiet when the entry is missing", func(t *testing.T) {
		runner := &fakePass{}
		store, _ := testPassStore(t, runner)

		if err := store.Delete(); err != nil {
			t.Errorf("Delete() error = %v, want nil", err)
		}

		if len(runner.calls) != 0 {
			t.Errorf("calls = %v, want none for a missing entry", runner.calls)
		}
	})

	t.Run("rejects an empty key", func(t *testing.T) {
		runner := &fakePass{}
		store, _ := testPassStore(t, runner)

		if err := store.Save(" "); !errors.Is(err, ErrEmptyKey) {
			t.Errorf("Save() error = %v, want ErrEmptyKey", err)
		}

		if len(runner.calls) != 0 {
			t.Errorf("calls = %v, want none for an empty key", runner.calls)
		}
	})
}

func TestNewPassStore(t *testing.T) {
	t.Run("detects an installed pass with a store", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("pass is a Unix tool")
		}

		binDir := t.TempDir()
		bin := filepath.Join(binDir, "pass")
		if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
		t.Setenv("PASSWORD_STORE_DIR", t.TempDir())

		store, ok := newPassStore(provider.Nano)
		if !ok {
			t.Fatal("newPassStore() = false, want the pass backend")
		}

		if _, err := store.Load(); !errors.Is(err, ErrNotFound) {
			t.Errorf("Load() error = %v, want ErrNotFound", err)
		}
	})

	t.Run("requires an existing store", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("pass is a Unix tool")
		}

		binDir := t.TempDir()
		bin := filepath.Join(binDir, "pass")
		if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
		t.Setenv("PASSWORD_STORE_DIR", filepath.Join(t.TempDir(), "missing"))

		if _, ok := newPassStore(provider.Nano); ok {
			t.Error("newPassStore() = true, want a missing store rejected")
		}
	})
}

func TestPassEntry(t *testing.T) {
	t.Run("keeps providers apart", func(t *testing.T) {
		seen := make(map[string]provider.Name)

		for _, name := range []provider.Name{provider.Nano, provider.Novel, provider.Danbooru} {
			entry, err := passEntry(name)
			if err != nil {
				t.Fatalf("passEntry(%s) error: %v", name, err)
			}

			if other, ok := seen[entry]; ok {
				t.Errorf("%s and %s share the entry %q", other, name, entry)
			}

			seen[entry] = name
		}
	})

	t.Run("rejects neo", func(t *testing.T) {
		if _, err := passEntry(provider.Neo); err == nil {
			t.Error("passEntry(neo) error = nil, want an error")
		}
	})
}

func TestPasswordStoreDir(t *testing.T) {
	t.Run("uses PASSWORD_STORE_DIR", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("PASSWORD_STORE_DIR", dir)

		if got := passwordStoreDir(); got != dir {
			t.Errorf("passwordStoreDir() = %q, want %q", got, dir)
		}
	})

	t.Run("falls back to the home directory", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("PASSWORD_STORE_DIR", "")
		t.Setenv("HOME", home)
		t.Setenv("USERPROFILE", home)

		want := filepath.Join(home, ".password-store")
		if got := passwordStoreDir(); got != want {
			t.Errorf("passwordStoreDir() = %q, want %q", got, want)
		}
	})
}
