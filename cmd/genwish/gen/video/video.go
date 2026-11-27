package genvideo

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func NewCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "video",
		Short: "Generate a video through a provider",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		NewNanoCommand(deps),
		NewNeoCommand(),
		NewNovelCommand(),
	)

	return cmd
}
