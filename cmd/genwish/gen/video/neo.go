package genvideo

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewNeoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "neo <prompt...>",
		Short: "Video generation is not supported by Forge Neo yet",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return fmt.Errorf("video generation is not supported for provider neo")
		},
	}
}
