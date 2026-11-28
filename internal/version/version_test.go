package version

import "testing"

func TestString(t *testing.T) {
	t.Run("uses the build-time version", func(t *testing.T) {
		previous := Version
		Version = "v1.2.3"
		t.Cleanup(func() { Version = previous })

		if got := String(); got != "v1.2.3" {
			t.Errorf("String() = %q, want v1.2.3", got)
		}
	})

	t.Run("falls back to dev", func(t *testing.T) {
		previous := Version
		Version = ""
		t.Cleanup(func() { Version = previous })

		if got := String(); got != "dev" {
			t.Errorf("String() = %q, want dev", got)
		}
	})
}
