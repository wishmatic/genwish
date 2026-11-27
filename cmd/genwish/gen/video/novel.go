package genvideo

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewNovelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "novel <prompt...>",
		Short: "Video generation is not supported by NovelAI yet",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return fmt.Errorf("video generation is not supported for provider novel")
		},
	}
}
