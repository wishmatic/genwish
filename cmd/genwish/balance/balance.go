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
		Short: "Read the account balance for a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
