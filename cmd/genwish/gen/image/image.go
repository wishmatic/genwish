package genimage

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func NewCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Generate an image through a provider",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		NewNanoCommand(deps),
		NewNeoCommand(deps),
		NewNovelCommand(deps),
	)

	return cmd
}
