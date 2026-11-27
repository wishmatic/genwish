package cli

import (
	"fmt"
	"strings"
)

func PromptFromArgs(args []string) (string, error) {
	prompt := strings.TrimSpace(strings.Join(args, " "))
	if prompt == "" {
		return "", fmt.Errorf("prompt is empty")
	}

	return prompt, nil
}
