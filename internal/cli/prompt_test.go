package cli

import (
	"strings"
	"testing"
)

func TestPromptFromArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr string
	}{
		{name: "joins arguments", args: []string{"a", "cat"}, want: "a cat"},
		{name: "trims surrounding space", args: []string{"  a cat  "}, want: "a cat"},
		{name: "rejects whitespace", args: []string{"  ", "\t"}, wantErr: "prompt is empty"},
		{name: "rejects no arguments", wantErr: "prompt is empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PromptFromArgs(tt.args)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("PromptFromArgs() error: %v", err)
				}

				if got != tt.want {
					t.Errorf("PromptFromArgs() = %q, want %q", got, tt.want)
				}

				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("PromptFromArgs() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}
