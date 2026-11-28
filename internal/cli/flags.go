package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func OptionalBool(cmd *cobra.Command, name string) *bool {
	if !cmd.Flags().Changed(name) {
		return nil
	}

	value, _ := cmd.Flags().GetBool(name)

	return &value
}

func OptionalInt(cmd *cobra.Command, name string) *int {
	if !cmd.Flags().Changed(name) {
		return nil
	}

	value, _ := cmd.Flags().GetInt(name)

	return &value
}

func OptionalFloat(cmd *cobra.Command, name string) *float64 {
	if !cmd.Flags().Changed(name) {
		return nil
	}

	value, _ := cmd.Flags().GetFloat64(name)

	return &value
}

func NormalizeFlagDefaults(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if !flagShowsDefault(flag) {
			return
		}

		if !strings.Contains(flag.Usage, "(default:") {
			flag.Usage += " (default: " + flag.DefValue + ")"
		}

		flag.DefValue = zeroDefault(flag)
	})

	for _, child := range cmd.Commands() {
		NormalizeFlagDefaults(child)
	}
}

func zeroDefault(flag *pflag.Flag) string {
	switch flag.Value.Type() {
	case "string":
		return ""
	case "bool":
		return "false"
	case "duration":
		return "0s"
	case "intSlice", "stringSlice", "stringArray":
		return "[]"
	case "ip", "ipMask", "ipNet":
		return "<nil>"
	default:
		return "0"
	}
}

func flagShowsDefault(flag *pflag.Flag) bool {
	if flag.DefValue == "" {
		return false
	}

	switch flag.Value.Type() {
	case "bool":
		return flag.DefValue != "false"
	case "duration":
		return flag.DefValue != "0" && flag.DefValue != "0s"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return flag.DefValue != "0"
	case "intSlice", "stringSlice", "stringArray":
		return flag.DefValue != "[]"
	case "ip", "ipMask", "ipNet":
		return flag.DefValue != "<nil>"
	default:
		return flag.DefValue != "false" && flag.DefValue != "0" &&
			flag.DefValue != "[]" && flag.DefValue != "<nil>"
	}
}
