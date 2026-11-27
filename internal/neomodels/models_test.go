package neomodels

import (
	"strings"
	"testing"
)

func ptr[T any](value T) *T {
	return &value
}

func TestModelValidate(t *testing.T) {
	tests := map[string]struct {
		model Model
		want  string
	}{
		"accepts an empty model": {
			model: Model{},
		},
		"accepts a default with no options": {
			model: Model{Samplers: Choice{Default: "Euler"}},
		},
		"requires a default when options are set": {
			model: Model{Samplers: Choice{Options: []string{"Euler"}}},
			want:  "samplers.default is required",
		},
		"requires the default among the options": {
			model: Model{Schedulers: Choice{Default: "Karras", Options: []string{"Automatic"}}},
			want:  "schedulers.default \"Karras\" is not one of the options",
		},
		"accepts a complete range": {
			model: Model{Steps: Range[int]{Default: ptr(20), Min: ptr(10), Max: ptr(40)}},
		},
		"requires a default with a bound": {
			model: Model{CFG: Range[float64]{Min: ptr(1.0), Max: ptr(2.0)}},
			want:  "cfg.default is required",
		},
		"requires both bounds": {
			model: Model{Steps: Range[int]{Default: ptr(20), Min: ptr(10)}},
			want:  "steps needs both min and max",
		},
		"rejects a reversed range": {
			model: Model{Steps: Range[int]{Default: ptr(20), Min: ptr(40), Max: ptr(10)}},
			want:  "steps.min must not exceed steps.max",
		},
		"rejects a default outside the range": {
			model: Model{Steps: Range[int]{Default: ptr(5), Min: ptr(10), Max: ptr(40)}},
			want:  "steps.default is outside min and max",
		},
		"accepts sizes": {
			model: Model{Sizes: []string{"1024x1024", "832x1216"}},
		},
		"rejects a size without a height": {
			model: Model{Sizes: []string{"1024"}},
			want:  "is not WxH",
		},
		"rejects a zero dimension": {
			model: Model{Sizes: []string{"0x512"}},
			want:  "is not WxH",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.model.Validate()

			if tt.want == "" {
				if err != nil {
					t.Fatalf("Validate() error: %v", err)
				}

				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Validate() error = %v, want it to name %q", err, tt.want)
			}
		})
	}
}
