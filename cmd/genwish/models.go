package main

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func newModelsCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "List the models a provider serves",
	}

	cmd.AddCommand(newModelsListCommand(deps))

	return cmd
}
