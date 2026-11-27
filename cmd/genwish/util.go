package main

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func newUtilCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "util",
		Short: "Utility commands",
	}

	cmd.AddCommand(newUtilBooruCommand(deps))

	return cmd
}
