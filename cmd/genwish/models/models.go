package models

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func NewCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Look up the models a provider serves",
	}

	cmd.AddCommand(NewListCommand(deps), NewGetCommand(deps))

	return cmd
}
