package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGenVideoNeoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "neo <prompt...>",
		Short: "Video generation is not supported by Forge Neo yet",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(*cobra.Command, []string) error {
			return fmt.Errorf("video generation is not supported for provider neo")
		},
	}
}
