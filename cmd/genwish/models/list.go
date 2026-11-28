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
		Long: "List the models a provider serves.\n\n" +
			"Providers: nano, neo, novel.",
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

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
	cmd.Flags().IntVar(&limit, "limit", 25, "models to show, at most 100")

	return cmd
}
