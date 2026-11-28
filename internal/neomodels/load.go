package neomodels

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

type Registry map[string]Model

type file struct {
	Neo map[string]Model `toml:"neo"`
}

func Load(path string) (Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(data)
}

func Parse(data []byte) (Registry, error) {
	var decoded file

	meta, err := toml.Decode(string(data), &decoded)
	if err != nil {
		return nil, err
	}

	if unknown := meta.Undecoded(); len(unknown) > 0 {
		return nil, fmt.Errorf("unknown key %s", formatKeys(unknown))
	}

	registry := Registry(decoded.Neo)

	for _, id := range registry.IDs() {
		model, _ := registry.Get(id)
		if err := model.Validate(); err != nil {
			return nil, fmt.Errorf("neo.%s: %w", id, err)
		}
	}

	return registry, nil
}

func (r Registry) Get(id string) (Model, bool) {
	model, ok := r[id]

	return model, ok
}

// Match finds a model by its id, or else by its checkpoint filename.
func (r Registry) Match(idOrCheckpoint string) (Model, bool) {
	if model, ok := r[idOrCheckpoint]; ok {
		return model, true
	}

	for _, id := range r.IDs() {
		model := r[id]
		if model.Checkpoint == idOrCheckpoint {
			return model, true
		}
	}

	return Model{}, false
}

func (r Registry) IDs() []string {
	ids := make([]string, 0, len(r))

	for id := range r {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	return ids
}

func formatKeys(keys []toml.Key) string {
	names := make([]string, 0, len(keys))

	for _, key := range keys {
		names = append(names, fmt.Sprintf("%q", key.String()))
	}

	return strings.Join(names, ", ")
}
