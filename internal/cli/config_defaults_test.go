package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/config"
)

func TestApplyConfigDefaults(t *testing.T) {
	t.Run("shows the configured value", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("format", "", "output format: png, jpeg, jxl, or webp")
		SetConfigDefault(cmd, "format", "output.format")

		ApplyConfigDefaults(cmd, config.Config{
			Output: config.OutputConfig{Format: "png"},
		})

		usage := cmd.Flags().Lookup("format").Usage
		if usage != "output format: png, jpeg, jxl, or webp (default: png)" {
			t.Errorf("usage = %q, want the configured format", usage)
		}
	})

	t.Run("replaces an existing default", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("out", "", "base directory (default: output.path)")
		SetConfigDefault(cmd, "out", "output.path")

		ApplyConfigDefaults(cmd, config.Config{
			Output: config.OutputConfig{Path: "/tmp/out"},
		})

		usage := cmd.Flags().Lookup("out").Usage
		if usage != "base directory (default: /tmp/out)" {
			t.Errorf("usage = %q, want the configured path", usage)
		}
	})

	t.Run("updates on a second pass", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("format", "", "output format")
		SetConfigDefault(cmd, "format", "output.format")

		ApplyConfigDefaults(cmd, config.Config{Output: config.OutputConfig{Format: "webp"}})
		ApplyConfigDefaults(cmd, config.Config{Output: config.OutputConfig{Format: "png"}})

		usage := cmd.Flags().Lookup("format").Usage
		if strings.Count(usage, "default:") != 1 || !strings.Contains(usage, "(default: png)") {
			t.Errorf("usage = %q, want a single updated default", usage)
		}
	})
}

func TestWithConfigDefaults(t *testing.T) {
	t.Run("fills help from the config", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: "test",
			RunE: func(*cobra.Command, []string) error {
				return nil
			},
		}
		cmd.Flags().String("format", "", "output format")
		SetConfigDefault(cmd, "format", "output.format")

		WithConfigDefaults(cmd, Deps{
			Config: config.Config{Output: config.OutputConfig{Format: "jxl"}},
		})

		out := &bytes.Buffer{}
		cmd.SetOut(out)
		cmd.SetErr(&bytes.Buffer{})

		if err := cmd.Help(); err != nil {
			t.Fatalf("Help() error: %v", err)
		}

		if !strings.Contains(out.String(), "(default: jxl)") {
			t.Errorf("help = %q, want the configured default", out)
		}
	})
}
