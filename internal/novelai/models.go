package novelai

type Model struct {
	ID   string
	Kind string
}

var models = []Model{
	{ID: "nai-diffusion-5-full", Kind: "image"},
	{ID: "nai-diffusion-5-full-medium", Kind: "image"},
	{ID: "nai-diffusion-5-full-inpainting", Kind: "image"},
	{ID: "nai-diffusion-5-full-medium-inpainting", Kind: "image"},
	{ID: "nai-diffusion-5-curated", Kind: "image"},
	{ID: "nai-diffusion-4-5-full", Kind: "image"},
	{ID: "nai-diffusion-4-5-full-inpainting", Kind: "image"},
	{ID: "nai-diffusion-4-5-curated", Kind: "image"},
	{ID: "nai-diffusion-4-5-curated-inpainting", Kind: "image"},
	{ID: "nai-diffusion-4-full", Kind: "image"},
	{ID: "nai-diffusion-4-full-inpainting", Kind: "image"},
	{ID: "nai-diffusion-4-curated-preview", Kind: "image"},
	{ID: "nai-diffusion-4-curated-inpainting", Kind: "image"},
	{ID: "nai-diffusion-3", Kind: "image"},
	{ID: "nai-diffusion-3-inpainting", Kind: "image"},
	{ID: "nai-diffusion-furry-3", Kind: "image"},
	{ID: "nai-diffusion-furry-3-inpainting", Kind: "image"},
}

func Models() []Model {
	out := make([]Model, len(models))
	copy(out, models)

	return out
}
