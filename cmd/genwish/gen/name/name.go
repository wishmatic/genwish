package genname

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wishmatic/namegen/gender"
	"github.com/wishmatic/namegen/name"

	"github.com/wishmatic/genwish/internal/cli"
)

const (
	defaultCount = 1
	maxCount     = 100
)

var genderNames = []string{string(gender.Male), string(gender.Female), string(gender.NonBinary)}

type generatedName struct {
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
	Gender    string `json:"gender"`
}

func NewCommand(deps cli.Deps) *cobra.Command {
	var (
		genderValue  string
		cultureValue string
		count        int
		asJSON       bool
	)

	cmd := &cobra.Command{
		Use:   "name",
		Short: "Generate one or more names",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runGenerate(deps, genderValue, cultureValue, count, asJSON)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(
		&genderValue,
		"gender",
		"",
		"male, female, or non-binary; unset picks one at random",
	)
	flags.StringVar(
		&cultureValue,
		"culture",
		"",
		"naming culture; unset picks one at random, and an unknown value lists them",
	)
	flags.IntVar(&count, "count", defaultCount, "how many names to generate")
	flags.BoolVar(&asJSON, "json", false, "print a JSON array instead of one name per line")

	return cmd
}

func runGenerate(deps cli.Deps, genderValue, cultureValue string, count int, asJSON bool) error {
	if count < 1 || count > maxCount {
		return fmt.Errorf("count must be between 1 and %d, got %d", maxCount, count)
	}

	resolvedGender, err := parseGender(genderValue)
	if err != nil {
		return err
	}

	culture, err := parseCulture(cultureValue)
	if err != nil {
		return err
	}

	names := make([]generatedName, 0, count)

	for range count {
		generated := name.Generate(resolvedGender, culture)

		names = append(
			names,
			generatedName{
				GivenName: generated.GivenName,
				Surname:   generated.Surname,
				Gender:    string(generated.Gender),
			},
		)
	}

	if asJSON {
		return cli.WriteJSON(deps.Out, names)
	}

	for _, generated := range names {
		line := fmt.Sprintf("%s %s", generated.GivenName, generated.Surname)

		if _, err := fmt.Fprintln(deps.Out, line); err != nil {
			return fmt.Errorf("write name: %w", err)
		}
	}

	return nil
}

func parseGender(value string) (gender.Gender, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}

	for _, candidate := range genderNames {
		if strings.EqualFold(candidate, trimmed) {
			return gender.Gender(candidate), nil
		}
	}

	return "", fmt.Errorf(
		"unknown gender %q; acceptable values: %s",
		value,
		strings.Join(genderNames, ", "),
	)
}

func parseCulture(value string) (*name.Culture, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	for _, candidate := range name.SupportedCultures {
		if strings.EqualFold(string(candidate), trimmed) {
			culture := candidate

			return &culture, nil
		}
	}

	return nil, fmt.Errorf(
		"unknown culture %q; acceptable values: %s",
		value,
		strings.Join(cultureNames(), ", "),
	)
}

func cultureNames() []string {
	names := make([]string, len(name.SupportedCultures))

	for i, culture := range name.SupportedCultures {
		names[i] = string(culture)
	}

	return names
}
