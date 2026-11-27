package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/credential"
	"github.com/wishmatic/genwish/internal/provider"
)

const maxAPIKeyBytes = 4096

func newLoginCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "login <provider>",
		Short: "Store a NanoGPT or NovelAI API key",
		Long: "Reads an API key from STDIN and stores it in the OS keyring, with a " +
			"user-only file fallback when no keyring is available",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name, err := provider.ParseCredential(args[0])
			if err != nil {
				return err
			}

			return login(deps, name)
		},
	}
}

func login(deps cli.Deps, name provider.Name) error {
	key, err := readAPIKey(deps, name)
	if err != nil {
		return err
	}

	store, err := cli.OpenCredentials(deps, name)
	if err != nil {
		return err
	}

	if err := store.Save(key); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	fmt.Fprintln(deps.Out, "logged in")

	return nil
}

func readAPIKey(deps cli.Deps, name provider.Name) (string, error) {
	if file, ok := deps.In.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		return readAPIKeyFromTerminal(deps, file, name)
	}

	return readAPIKeyFromReader(deps.In)
}

func readAPIKeyFromTerminal(deps cli.Deps, file *os.File, name provider.Name) (string, error) {
	fmt.Fprint(deps.Err, keyPrompt(name))

	key, err := term.ReadPassword(int(file.Fd()))
	fmt.Fprintln(deps.Err)

	if err != nil {
		return "", fmt.Errorf("read API key: %w", err)
	}

	return normalizeAPIKey(key)
}

func readAPIKeyFromReader(r io.Reader) (string, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxAPIKeyBytes+1))
	if err != nil {
		return "", fmt.Errorf("read API key: %w", err)
	}

	if len(data) > maxAPIKeyBytes {
		return "", fmt.Errorf("API key is longer than %d bytes", maxAPIKeyBytes)
	}

	return normalizeAPIKey(data)
}

func normalizeAPIKey(data []byte) (string, error) {
	key := strings.TrimSpace(string(data))
	if key == "" {
		return "", credential.ErrEmptyKey
	}

	if strings.ContainsAny(key, "\r\n") {
		return "", fmt.Errorf("API key must be a single line")
	}

	return key, nil
}

func keyPrompt(name provider.Name) string {
	switch name {
	case provider.Nano:
		return "NanoGPT API key: "
	case provider.Novel:
		return "NovelAI API key: "
	default:
		return "API key: "
	}
}
