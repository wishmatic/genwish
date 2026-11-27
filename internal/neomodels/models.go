package neomodels

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Model struct {
	Checkpoint   string         `toml:"checkpoint"`
	Encoders     []string       `toml:"encoders"`
	Preset       string         `toml:"preset"`
	Samplers     Choice         `toml:"samplers"`
	Schedulers   Choice         `toml:"schedulers"`
	Steps        Range[int]     `toml:"steps"`
	CFG          Range[float64] `toml:"cfg"`
	Hires        Hires          `toml:"hires"`
	Sizes        []string       `toml:"sizes"`
	Instructions Instructions   `toml:"instructions"`
}

type Choice struct {
	Default string   `toml:"default"`
	Options []string `toml:"options"`
}

type Range[T int | float64] struct {
	Default *T `toml:"default"`
	Min     *T `toml:"min"`
	Max     *T `toml:"max"`
}

type Hires struct {
	Enabled bool `toml:"enabled"`
}

type Instructions struct {
	Usage string `toml:"usage"`
	Style string `toml:"style"`
}

func (m Model) Validate() error {
	if err := m.Samplers.validate("samplers"); err != nil {
		return err
	}

	if err := m.Schedulers.validate("schedulers"); err != nil {
		return err
	}

	if err := m.Steps.validate("steps"); err != nil {
		return err
	}

	if err := m.CFG.validate("cfg"); err != nil {
		return err
	}

	for _, size := range m.Sizes {
		if err := validateSize(size); err != nil {
			return err
		}
	}

	return nil
}

func (c Choice) validate(name string) error {
	if len(c.Options) == 0 {
		return nil
	}

	if c.Default == "" {
		return fmt.Errorf("%s.default is required when options are set", name)
	}

	if !slices.Contains(c.Options, c.Default) {
		return fmt.Errorf("%s.default %q is not one of the options", name, c.Default)
	}

	return nil
}

func (r Range[T]) validate(name string) error {
	if r.Min == nil && r.Max == nil {
		return nil
	}

	if r.Default == nil {
		return fmt.Errorf("%s.default is required when min or max is set", name)
	}

	if r.Min == nil || r.Max == nil {
		return fmt.Errorf("%s needs both min and max", name)
	}

	if *r.Min > *r.Max {
		return fmt.Errorf("%s.min must not exceed %s.max", name, name)
	}

	if *r.Default < *r.Min || *r.Default > *r.Max {
		return fmt.Errorf("%s.default is outside min and max", name)
	}

	return nil
}

func validateSize(size string) error {
	width, height, ok := strings.Cut(size, "x")
	if !ok {
		return fmt.Errorf("size %q is not WxH", size)
	}

	parsedWidth, err := strconv.Atoi(width)
	if err != nil || parsedWidth <= 0 {
		return fmt.Errorf("size %q is not WxH", size)
	}

	parsedHeight, err := strconv.Atoi(height)
	if err != nil || parsedHeight <= 0 {
		return fmt.Errorf("size %q is not WxH", size)
	}

	return nil
}
