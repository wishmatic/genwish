package main

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func newGenVideoCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "video",
		Short: "Generate a video through a provider",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newGenVideoNanoCommand(deps), newGenVideoNeoCommand(), newGenVideoNovelCommand())

	return cmd
}
