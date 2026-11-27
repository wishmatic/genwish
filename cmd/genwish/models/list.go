package models

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/provider"
)

func NewListCommand(deps cli.Deps) *cobra.Command {
	var (
		query string
		limit int
	)

	cmd := &cobra.Command{
		Use:   "list <provider> [image|video]",
		Short: "List the models a provider serves",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := provider.Parse(args[0])
			if err != nil {
				return err
			}

			switch name {
			case provider.Nano:
				return runNano(cmd.Context(), deps, args[1:], query, limit)
			case provider.Neo:
				return runNeo(deps, args[1:], query, limit)
			case provider.Novel:
				return runNovel(deps, args[1:], query, limit)
			default:
				return fmt.Errorf("models list %s is not supported yet", name)
			}
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "filter by id, name, description, owner, or tags")
	cmd.Flags().IntVar(&limit, "limit", 0, "models to show, 25 default and 100 at most")

	return cmd
}
