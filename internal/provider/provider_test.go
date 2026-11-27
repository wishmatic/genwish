package provider

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  Name
	}{
		{name: "nano", value: "nano", want: Nano},
		{name: "nano-gpt", value: "nano-gpt", want: Nano},
		{name: "nanogpt", value: "nanogpt", want: Nano},
		{name: "nano gpt", value: "nano gpt", want: Nano},
		{name: "nano gpt upper", value: "NANO GPT", want: Nano},
		{name: "neo", value: "neo", want: Neo},
		{name: "novel", value: "novel", want: Novel},
		{name: "novel-ai", value: "novel-ai", want: Novel},
		{name: "novel ai", value: "novel ai", want: Novel},
		{name: "novelai", value: "novelai", want: Novel},
		{name: "novel ai mixed", value: "Novel AI", want: Novel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.value)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.value, err)
			}

			if got != tt.want {
				t.Errorf("Parse(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}

	for _, value := range []string{"", "  ", "atlantean"} {
		t.Run("rejects "+value, func(t *testing.T) {
			if _, err := Parse(value); err == nil {
				t.Errorf("Parse(%q) error = nil, want an error", value)
			}
		})
	}
}

func TestParseBalance(t *testing.T) {
	for _, value := range []string{"nano", "nano-gpt", "nanogpt", "novel", "novel-ai", "novel ai", "novelai"} {
		t.Run("accepts "+value, func(t *testing.T) {
			if _, err := ParseBalance(value); err != nil {
				t.Errorf("ParseBalance(%q) error: %v", value, err)
			}
		})
	}

	t.Run("rejects neo", func(t *testing.T) {
		if _, err := ParseBalance("neo"); err == nil {
			t.Error("ParseBalance(neo) error = nil, want an error")
		}
	})
}

func TestParseCredential(t *testing.T) {
	for _, value := range []string{"nano", "nanogpt", "novel", "novelai", "danbooru", "Danbooru"} {
		t.Run("accepts "+value, func(t *testing.T) {
			if _, err := ParseCredential(value); err != nil {
				t.Errorf("ParseCredential(%q) error: %v", value, err)
			}
		})
	}

	t.Run("rejects neo", func(t *testing.T) {
		if _, err := ParseCredential("neo"); err == nil {
			t.Error("ParseCredential(neo) error = nil, want an error")
		}
	})
}

func TestParseRejectsDanbooru(t *testing.T) {
	if _, err := Parse("danbooru"); err == nil {
		t.Error("Parse(danbooru) error = nil, want the generation providers only")
	}
}
