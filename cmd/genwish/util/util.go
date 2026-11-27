package util

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func NewCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "util",
		Short: "Utility commands",
	}

	cmd.AddCommand(NewBooruCommand(deps))

	return cmd
}
