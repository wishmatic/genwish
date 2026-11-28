package credential

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wishmatic/genwish/internal/provider"
)

type passRunner func(args []string, stdin io.Reader) (stdout, stderr []byte, err error)

type PassStore struct {
	entry string
	path  string
	run   passRunner
}

func newPassStore(p provider.Name) (Store, bool) {
	bin, err := exec.LookPath("pass")
	if err != nil {
		return nil, false
	}

	dir := passwordStoreDir()
	if dir == "" {
		return nil, false
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, false
	}

	entry, err := passEntry(p)
	if err != nil {
		return nil, false
	}

	return &PassStore{
		entry: entry,
		path:  filepath.Join(dir, fmt.Sprintf("%s.gpg", entry)),
		run: func(args []string, stdin io.Reader) ([]byte, []byte, error) {
			return runPass(bin, args, stdin)
		},
	}, true
}

func runPass(bin string, args []string, stdin io.Reader) ([]byte, []byte, error) {
	var stderr bytes.Buffer

	cmd := exec.Command(bin, args...)
	cmd.Stderr = &stderr

	if stdin != nil {
		cmd.Stdin = stdin
	}

	out, err := cmd.Output()

	return out, stderr.Bytes(), err
}

func (s *PassStore) Load() (string, error) {
	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		return "", ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("read pass entry: %w", err)
	}

	out, stderr, err := s.run([]string{"show", s.entry}, nil)
	if err != nil {
		return "", passError("read", stderr, err)
	}

	secret := strings.TrimSpace(string(out))
	if secret == "" {
		return "", ErrNotFound
	}

	return secret, nil
}

func (s *PassStore) Save(secret string) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ErrEmptyKey
	}

	stdin := strings.NewReader(fmt.Sprintf("%s\n", secret))

	_, stderr, err := s.run([]string{"insert", "-m", "-f", s.entry}, stdin)
	if err != nil {
		return passError("store", stderr, err)
	}

	return nil
}

func (s *PassStore) Delete() error {
	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("read pass entry: %w", err)
	}

	_, stderr, err := s.run([]string{"rm", "-f", s.entry}, nil)
	if err != nil {
		return passError("delete", stderr, err)
	}

	return nil
}

func passError(action string, stderr []byte, err error) error {
	if message := strings.TrimSpace(string(stderr)); message != "" {
		return fmt.Errorf("%s pass entry: %w: %s", action, err, message)
	}

	return fmt.Errorf("%s pass entry: %w", action, err)
}

func passwordStoreDir() string {
	if dir := strings.TrimSpace(os.Getenv("PASSWORD_STORE_DIR")); dir != "" {
		return dir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".password-store")
}

func passEntry(p provider.Name) (string, error) {
	switch p {
	case provider.Nano:
		return "genwish/nanogpt-api-key", nil
	case provider.Novel:
		return "genwish/novelai-api-key", nil
	case provider.Danbooru:
		return "genwish/danbooru", nil
	default:
		return "", fmt.Errorf("provider %q does not use an API key", p)
	}
}
