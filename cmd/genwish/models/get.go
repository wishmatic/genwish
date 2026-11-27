package models

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/provider"
)

func NewGetCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "get <provider> <id> [image|video]",
		Short: "Show everything known about one model",
		Long: "Show everything known about one model.\n\n" +
			"Providers: nano, neo, novel.",
		Args: cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return cmd.Help()
			}

			name, err := provider.Parse(args[0])
			if err != nil {
				return err
			}

			id := args[1]

			switch name {
			case provider.Nano:
				return getNano(cmd.Context(), deps, id, args[2:])
			case provider.Neo:
				return getNeo(deps, id, args[2:])
			case provider.Novel:
				return getNovel(deps, id, args[2:])
			default:
				return fmt.Errorf("models get %s is not supported yet", name)
			}
		},
	}
}
