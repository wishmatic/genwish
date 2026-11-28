package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNormalizeFlagDefaults(t *testing.T) {
	t.Run("moves the default into the usage", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("model", "wan-3.0", "video model id")

		NormalizeFlagDefaults(cmd)

		flag := cmd.Flags().Lookup("model")
		if flag.Usage != "video model id (default: wan-3.0)" {
			t.Errorf("usage = %q, want the default moved in", flag.Usage)
		}

		if flag.DefValue != "" {
			t.Errorf("DefValue = %q, want it cleared", flag.DefValue)
		}
	})

	t.Run("leaves zero defaults alone", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().Int("count", 0, "how many")

		NormalizeFlagDefaults(cmd)

		flag := cmd.Flags().Lookup("count")
		if flag.Usage != "how many" || flag.DefValue != "0" {
			t.Errorf("flag = %+v, want the zero default untouched", flag)
		}
	})

	t.Run("keeps a manual default", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("out", "", "base directory (default: output.path)")

		NormalizeFlagDefaults(cmd)

		flag := cmd.Flags().Lookup("out")
		if strings.Count(flag.Usage, "default:") != 1 {
			t.Errorf("usage = %q, want one default", flag.Usage)
		}
	})

	t.Run("normalizes children", func(t *testing.T) {
		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{Use: "child"}
		child.Flags().Int("steps", 20, "sampling steps")
		parent.AddCommand(child)

		NormalizeFlagDefaults(parent)

		if usage := child.Flags().Lookup("steps").Usage; usage != "sampling steps (default: 20)" {
			t.Errorf("usage = %q, want the child's default", usage)
		}
	})

	t.Run("renders the default with a colon", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().Int("steps", 20, "sampling steps")

		NormalizeFlagDefaults(cmd)

		usage := cmd.Flags().FlagUsages()
		if !strings.Contains(usage, "(default: 20)") || strings.Contains(usage, "(default )") {
			t.Errorf("usages = %q, want only the colon form", usage)
		}
	})
}
