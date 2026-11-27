package main

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func newGenImageCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Generate an image through a provider",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newGenImageNanoCommand(deps), newGenImageNeoCommand(deps), newGenImageNovelCommand(deps))

	return cmd
}
