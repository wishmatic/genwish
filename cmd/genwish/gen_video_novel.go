package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGenVideoNovelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "novel <prompt...>",
		Short: "Video generation is not supported by NovelAI yet",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(*cobra.Command, []string) error {
			return fmt.Errorf("video generation is not supported for provider novel")
		},
	}
}
