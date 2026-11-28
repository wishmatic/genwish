package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/booru"
	"github.com/wishmatic/genwish/internal/catalog"
	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

type stubSource struct {
	page  booru.TagPage
	err   error
	calls int

	tag     booru.Tag
	found   bool
	findErr error

	aliasTarget string
	aliased     bool
	aliasErr    error

	implications []string
	implErr      error

	query booru.TagQuery
}

func (s *stubSource) SearchTags(_ context.Context, query booru.TagQuery) (booru.TagPage, error) {
	s.calls++
	s.query = query

	return s.page, s.err
}

func (s *stubSource) FindTag(_ context.Context, name string) (booru.Tag, bool, error) {
	s.calls++
	s.query = booru.TagQuery{Search: name}

	return s.tag, s.found, s.findErr
}

func (s *stubSource) AliasTarget(_ context.Context, name string) (string, bool, error) {
	s.calls++
	s.query = booru.TagQuery{Search: name}

	return s.aliasTarget, s.aliased, s.aliasErr
}

func (s *stubSource) Implications(_ context.Context, name string) ([]string, error) {
	s.calls++
	s.query = booru.TagQuery{Search: name}

	return s.implications, s.implErr
}

func runBooru(t *testing.T, source catalog.TagSource, args ...string) (string, error) {
	t.Helper()

	out := &bytes.Buffer{}
	deps := cli.Deps{
		Out: out,
		Err: &bytes.Buffer{},
		Catalog: catalog.New(
			source,
			catalog.Options{MaxLimit: 100, MaxOffset: 1000},
		),
	}

	err := testutil.Run(t, NewBooruCommand(deps), args...)

	return out.String(), err
}

func TestUtilBooru(t *testing.T) {
	t.Run("prints the tool text", func(t *testing.T) {
		source := &stubSource{page: booru.TagPage{
			Tags: []booru.Tag{{Name: "blue_hair", Category: booru.CategoryGeneral, Count: 482913}},
			More: true,
		}}

		out, err := runBooru(t, source, "blue", "hair", "--limit", "5")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{
			"blue_hair (general): 482913 works", "More matches exist; continue with offset 1",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}

		if source.query.Search != "blue_hair" || source.query.Offset != 0 ||
			source.query.Limit != 5 {
			t.Errorf("query = %+v, want the joined search and the requested window", source.query)
		}

		if source.query.MinCount != defaultMinCount {
			t.Errorf("min count = %d, want %d", source.query.MinCount, defaultMinCount)
		}
	})

	t.Run("sends every flag", func(t *testing.T) {
		source := &stubSource{}

		_, err := runBooru(
			t,
			source,
			"blue",
			"hair",
			"--offset",
			"2",
			"--limit",
			"7",
			"--min-count",
			"0",
			"--category",
			"artist",
			"--category",
			"meta",
		)
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if source.query.Offset != 2 || source.query.Limit != 7 || source.query.MinCount != 0 {
			t.Errorf("window = %+v, want offset 2, limit 7, min count 0", source.query)
		}

		want := []booru.TagCategory{booru.CategoryArtist, booru.CategoryMeta}
		if len(source.query.Categories) != len(want) {
			t.Fatalf("categories = %v, want %v", source.query.Categories, want)
		}

		for i, category := range want {
			if source.query.Categories[i] != category {
				t.Errorf("category %d = %q, want %q", i, source.query.Categories[i], category)
			}
		}
	})

	t.Run("sends exact", func(t *testing.T) {
		source := &stubSource{
			tag:   booru.Tag{Name: "1girl", Category: booru.CategoryGeneral, Count: 8455477},
			found: true,
		}

		out, err := runBooru(t, source, "1girl", "--exact")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out, "- 1girl (general): 8455477 works") {
			t.Errorf("stdout = %q, want the exact tag", out)
		}
	})

	t.Run("prints JSON", func(t *testing.T) {
		source := &stubSource{page: booru.TagPage{
			Tags: []booru.Tag{{Name: "blue_hair", Category: booru.CategoryGeneral, Count: 482913}},
		}}

		out, err := runBooru(t, source, "blue", "hair", "--json")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		var decoded tagsOutput
		if err := json.Unmarshal([]byte(out), &decoded); err != nil {
			t.Fatalf("decode stdout: %v\n%s", err, out)
		}

		if decoded.Search != "blue_hair" || decoded.Offset != 0 ||
			decoded.Limit != catalog.DefaultTagLimit {
			t.Errorf("output = %+v, want the echoed request", decoded)
		}

		if len(decoded.Tags) != 1 || decoded.Tags[0].Name != "blue_hair" ||
			decoded.Tags[0].Count != 482913 {
			t.Errorf("tags = %+v, want the matching tag", decoded.Tags)
		}

		if decoded.Synonyms == nil || decoded.Tags[0].Implications == nil {
			t.Errorf("output = %+v, want empty arrays rather than null", decoded)
		}
	})

	badInputs := map[string]struct {
		args      []string
		wantNamed string
	}{
		"unknown category": {
			args:      []string{"blue", "--category", "nonsense"},
			wantNamed: "general",
		},
		"negative offset": {args: []string{"blue", "--offset=-1"}, wantNamed: "offset"},
		"huge limit":      {args: []string{"blue", "--limit=101"}, wantNamed: "limit"},
		"negative floor":  {args: []string{"blue", "--min-count=-1"}, wantNamed: "min_count"},
		"empty search":    {args: []string{"   "}, wantNamed: "search"},
	}

	for name, tt := range badInputs {
		t.Run("rejects "+name, func(t *testing.T) {
			source := &stubSource{}

			_, err := runBooru(t, source, tt.args...)
			if err == nil {
				t.Fatal("Execute() error = nil, want an error")
			}

			if !strings.Contains(err.Error(), tt.wantNamed) {
				t.Errorf("error = %q, want it to name %s", err, tt.wantNamed)
			}

			if source.calls != 0 {
				t.Errorf("provider calls = %d, want none", source.calls)
			}
		})
	}

	t.Run("reports an upstream failure", func(t *testing.T) {
		source := &stubSource{err: errors.New("danbooru: upstream is down")}

		out, err := runBooru(t, source, "blue")
		if err == nil || !strings.Contains(err.Error(), "tags:") ||
			!strings.Contains(err.Error(), "upstream is down") {
			t.Errorf("Execute() error = %v, want the upstream failure", err)
		}

		if out != "" {
			t.Errorf("stdout = %q, want no page printed", out)
		}
	})
}

func TestBooruShowsHelpWithoutASearch(t *testing.T) {
	out := &bytes.Buffer{}
	cmd := NewBooruCommand(cli.Deps{})
	cmd.SetOut(out)
	cmd.SetErr(&bytes.Buffer{})

	if err := testutil.Run(t, cmd); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(out.String(), "booru <search...>") {
		t.Errorf("stdout = %q, want the search usage", out)
	}
}
