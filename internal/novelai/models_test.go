package novelai

import "testing"

func TestModels(t *testing.T) {
	t.Run("lists image models", func(t *testing.T) {
		got := Models()

		if len(got) == 0 {
			t.Fatal("Models() = none, want the known NovelAI image models")
		}

		if got[0].ID != "nai-diffusion-5-full" || got[0].Kind != "image" {
			t.Errorf("first model = %+v, want nai-diffusion-5-full", got[0])
		}
	})

	t.Run("returns a copy", func(t *testing.T) {
		first := Models()
		first[0].ID = "changed"

		second := Models()
		if second[0].ID != "nai-diffusion-5-full" {
			t.Errorf("Models() shares its backing array: %q", second[0].ID)
		}
	})
}
