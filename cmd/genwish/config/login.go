package config

import (
	"errors"
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

func NewLoginCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "login <provider>",
		Short: "Store a NanoGPT, NovelAI, or Danbooru credential",
		Long: "Reads credentials from STDIN and stores them in the OS keyring, with a " +
			"user-only file fallback when no keyring is available. Danbooru reads a login " +
			"line followed by an API key line.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			name, err := provider.ParseCredential(args[0])
			if err != nil {
				return err
			}

			if name == provider.Danbooru {
				return loginDanbooru(deps)
			}

			return login(deps, name)
		},
	}
}

func login(deps cli.Deps, name provider.Name) error {
	warnAboutStorage(deps, name)

	store, err := cli.OpenCredentials(deps, name)
	if err != nil {
		return err
	}

	warnAboutStored(deps, name, store)

	key, err := readAPIKey(deps, name)
	if err != nil {
		return err
	}

	if err := store.Save(key); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	fmt.Fprintln(deps.Out, "logged in")

	return nil
}

func loginDanbooru(deps cli.Deps) error {
	warnAboutStorage(deps, provider.Danbooru)

	store, err := cli.OpenCredentials(deps, provider.Danbooru)
	if err != nil {
		return err
	}

	warnAboutStored(deps, provider.Danbooru, store)

	creds, err := readDanbooruCredentials(deps)
	if err != nil {
		return err
	}

	if err := credential.SaveDanbooru(store, creds); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	fmt.Fprintln(deps.Out, "logged in")

	return nil
}

func warnAboutStored(deps cli.Deps, name provider.Name, store credential.Store) {
	if _, err := store.Load(); err != nil {
		return
	}

	fmt.Fprintf(
		deps.Err,
		"warning: %s is already logged in; the stored credential will be replaced\n",
		name,
	)
}

func warnAboutStorage(deps cli.Deps, name provider.Name) {
	if deps.Credentials != nil {
		return
	}

	if warning := credential.StorageWarning(name); warning != "" {
		fmt.Fprintln(deps.Err, warning)
	}
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

func readDanbooruCredentials(deps cli.Deps) (credential.DanbooruCredentials, error) {
	if file, ok := deps.In.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		return readDanbooruFromTerminal(deps, file)
	}

	return readDanbooruFromReader(deps.In)
}

func readDanbooruFromTerminal(
	deps cli.Deps,
	file *os.File,
) (credential.DanbooruCredentials, error) {
	fmt.Fprint(deps.Err, "Danbooru login: ")

	login, err := readTerminalLine(deps.In)
	if err != nil {
		return credential.DanbooruCredentials{}, fmt.Errorf("read login: %w", err)
	}

	login = strings.TrimSpace(login)
	if login == "" {
		return credential.DanbooruCredentials{}, fmt.Errorf("Danbooru login is empty")
	}

	fmt.Fprint(deps.Err, "Danbooru API key: ")

	key, err := term.ReadPassword(int(file.Fd()))
	fmt.Fprintln(deps.Err)

	if err != nil {
		return credential.DanbooruCredentials{}, fmt.Errorf("read API key: %w", err)
	}

	apiKey, err := normalizeAPIKey(key)
	if err != nil {
		return credential.DanbooruCredentials{}, err
	}

	return credential.DanbooruCredentials{Login: login, APIKey: apiKey}, nil
}

func readDanbooruFromReader(r io.Reader) (credential.DanbooruCredentials, error) {
	data, err := io.ReadAll(io.LimitReader(r, 2*maxAPIKeyBytes+2))
	if err != nil {
		return credential.DanbooruCredentials{}, fmt.Errorf("read credentials: %w", err)
	}

	lines := make([]string, 0, 2)

	for _, line := range strings.Split(string(data), "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}

	if len(lines) != 2 {
		return credential.DanbooruCredentials{}, fmt.Errorf(
			"expected a Danbooru login line then an API key line",
		)
	}

	apiKey, err := normalizeAPIKey([]byte(lines[1]))
	if err != nil {
		return credential.DanbooruCredentials{}, err
	}

	return credential.DanbooruCredentials{Login: lines[0], APIKey: apiKey}, nil
}

func readTerminalLine(r io.Reader) (string, error) {
	var b strings.Builder

	buf := make([]byte, 1)

	for {
		n, err := r.Read(buf)
		if n > 0 {
			switch buf[0] {
			case '\n':
				return b.String(), nil
			case '\r':
			default:
				b.WriteByte(buf[0])
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				return b.String(), nil
			}

			return "", err
		}
	}
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
