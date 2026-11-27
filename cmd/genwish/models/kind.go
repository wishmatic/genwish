package models

import (
	"fmt"
	"strings"
)

func noVideoKind(name string, args []string) error {
	if len(args) == 0 {
		return nil
	}

	kind := strings.ToLower(strings.TrimSpace(args[0]))
	if kind != "image" {
		return fmt.Errorf("%s does not serve %s models", name, kind)
	}

	return nil
}
