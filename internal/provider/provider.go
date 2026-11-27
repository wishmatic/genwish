package provider

import (
	"fmt"
	"strings"
)

type Name string

const (
	Nano  Name = "nano"
	Neo   Name = "neo"
	Novel Name = "novel"
)

func (n Name) String() string {
	return string(n)
}

func Parse(value string) (Name, error) {
	switch normalize(value) {
	case "nano", "nanogpt":
		return Nano, nil
	case "neo":
		return Neo, nil
	case "novel", "novelai":
		return Novel, nil
	default:
		return "", fmt.Errorf("provider %q is not one of nano, neo, novel", value)
	}
}

func ParseBalance(value string) (Name, error) {
	name, err := Parse(value)
	if err != nil {
		return "", err
	}

	if name == Neo {
		return "", fmt.Errorf("provider %q has no balance", value)
	}

	return name, nil
}

func ParseCredential(value string) (Name, error) {
	name, err := Parse(value)
	if err != nil {
		return "", err
	}

	if name == Neo {
		return "", fmt.Errorf("provider %q does not use an API key", value)
	}

	return name, nil
}

func normalize(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	return normalized
}
