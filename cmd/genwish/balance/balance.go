package balance

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/provider"
)

func NewCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "balance <provider>",
		Short: "Read account credit balance",
		Long:  "Read the account credit balance for a provider (nano or novel).",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			name, err := provider.ParseBalance(args[0])
			if err != nil {
				return err
			}

			switch name {
			case provider.Nano:
				return RunNano(cmd.Context(), deps)
			case provider.Novel:
				return RunNovel(cmd.Context(), deps)
			default:
				return fmt.Errorf("balance %s is not supported yet", name)
			}
		},
	}
}
