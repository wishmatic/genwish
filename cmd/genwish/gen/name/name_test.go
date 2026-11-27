package genname

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func runName(t *testing.T, deps cli.Deps, args ...string) (string, error) {
	t.Helper()

	out := &bytes.Buffer{}
	deps.Out = out
	if deps.Err == nil {
		deps.Err = &bytes.Buffer{}
	}

	err := testutil.Run(t, NewCommand(deps), args...)

	return out.String(), err
}

func TestNewCommand(t *testing.T) {
	t.Run("generates one name", func(t *testing.T) {
		out, err := runName(t, cli.Deps{})
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != 1 || strings.TrimSpace(lines[0]) == "" {
			t.Errorf("stdout = %q, want one name", out)
		}
	})

	t.Run("generates the requested count", func(t *testing.T) {
		out, err := runName(t, cli.Deps{}, "--count", "3")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != 3 {
			t.Errorf("stdout = %q, want three names", out)
		}
	})

	t.Run("pins the gender", func(t *testing.T) {
		out, err := runName(t, cli.Deps{}, "--gender", "female", "--count", "2", "--json")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		var names []generatedName
		if err := json.Unmarshal([]byte(out), &names); err != nil {
			t.Fatalf("decode stdout: %v\n%s", err, out)
		}

		if len(names) != 2 {
			t.Fatalf("JSON = %s, want two names", out)
		}

		for _, name := range names {
			if name.Gender != "female" {
				t.Errorf("gender = %q, want female", name.Gender)
			}

			if name.GivenName == "" || name.Surname == "" {
				t.Errorf("name = %+v, want a given name and surname", name)
			}
		}
	})

	t.Run("accepts a culture", func(t *testing.T) {
		out, err := runName(t, cli.Deps{}, "--culture", "japanese")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if strings.TrimSpace(out) == "" {
			t.Error("stdout is empty, want a name")
		}
	})

	t.Run("prints JSON with every field", func(t *testing.T) {
		out, err := runName(t, cli.Deps{}, "--json")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		var names []map[string]any
		if err := json.Unmarshal([]byte(out), &names); err != nil {
			t.Fatalf("decode stdout: %v\n%s", err, out)
		}

		if len(names) != 1 {
			t.Fatalf("JSON = %s, want one name", out)
		}

		for _, key := range []string{"given_name", "surname", "gender"} {
			if _, ok := names[0][key]; !ok {
				t.Errorf("JSON %s does not contain %q", out, key)
			}
		}
	})

	t.Run("rejects an unknown gender", func(t *testing.T) {
		_, err := runName(t, cli.Deps{}, "--gender", "elf")
		if err == nil || !strings.Contains(err.Error(), "unknown gender") {
			t.Errorf("Execute() error = %v, want the unknown gender", err)
		}
	})

	t.Run("rejects an unknown culture", func(t *testing.T) {
		_, err := runName(t, cli.Deps{}, "--culture", "atlantean")
		if err == nil || !strings.Contains(err.Error(), "unknown culture") {
			t.Errorf("Execute() error = %v, want the unknown culture", err)
		}
	})

	t.Run("rejects a count outside the cap", func(t *testing.T) {
		for _, count := range []string{"0", "101"} {
			_, err := runName(t, cli.Deps{}, "--count", count)
			if err == nil || !strings.Contains(err.Error(), "between 1 and 100") {
				t.Errorf("Execute() error for count %s = %v, want the accepted range", count, err)
			}
		}
	})
}
