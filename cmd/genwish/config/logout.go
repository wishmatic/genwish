package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/provider"
)

func NewLogoutCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "logout <provider>",
		Short: "Delete a stored NanoGPT, NovelAI, or Danbooru credential",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			name, err := provider.ParseCredential(args[0])
			if err != nil {
				return err
			}

			return logout(deps, name)
		},
	}
}

func logout(deps cli.Deps, name provider.Name) error {
	store, err := cli.OpenCredentials(deps, name)
	if err != nil {
		return err
	}

	if err := store.Delete(); err != nil {
		return fmt.Errorf("delete API key: %w", err)
	}

	fmt.Fprintln(deps.Out, "logged out")

	return nil
}
