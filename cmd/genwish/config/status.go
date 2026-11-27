package config

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/credential"
	"github.com/wishmatic/genwish/internal/provider"
)

var credentialProviders = []provider.Name{provider.Nano, provider.Novel, provider.Danbooru}

type storeOpener func(provider.Name) (credential.Store, error)

func NewStatusCommand(deps cli.Deps) *cobra.Command {
	return newStatusCommand(deps, func(name provider.Name) (credential.Store, error) {
		return cli.OpenCredentials(deps, name)
	})
}

func newStatusCommand(deps cli.Deps, open storeOpener) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show which providers have a stored login",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runStatus(deps, open)
		},
	}
}

func runStatus(deps cli.Deps, open storeOpener) error {
	writer := tabwriter.NewWriter(deps.Out, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(writer, "PROVIDER\tSTATUS"); err != nil {
		return err
	}

	for _, name := range credentialProviders {
		status, err := providerStatus(open, name)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(writer, "%s\t%s\n", name, status); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func providerStatus(open storeOpener, name provider.Name) (string, error) {
	store, err := open(name)
	if err != nil {
		return "", err
	}

	if name == provider.Danbooru {
		return danbooruStatus(store)
	}

	if _, err := store.Load(); err != nil {
		if errors.Is(err, credential.ErrNotFound) {
			return "not logged in", nil
		}

		return "", err
	}

	return "logged in", nil
}

func danbooruStatus(store credential.Store) (string, error) {
	creds, err := credential.LoadDanbooru(store)
	if err != nil {
		if errors.Is(err, credential.ErrNotFound) {
			return "not logged in", nil
		}

		return "", err
	}

	return "logged in as " + creds.Login, nil
}
