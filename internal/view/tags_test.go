package view

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/booru"
	"github.com/wishmatic/genwish/internal/catalog"
)

func TestTagsEchoesTheRequest(t *testing.T) {
	result := catalog.TagsResult{
		Search:       "blue_hair",
		Status:       catalog.StatusOK,
		SnapshotDate: "2026-09-24",
		More:         true,
		AliasOf:      "futanari",
		Withheld:     2,
		WithheldBest: 3,
		Synonyms:     []string{"pee"},
		Tags: []booru.Tag{{
			Name:          "blue_hair",
			Category:      booru.CategoryGeneral,
			Count:         482913,
			CountIsTarget: true,
			AliasOf:       "futanari",
			Implications:  []string{"aqua_hair"},
		}},
	}

	out := Tags(
		result,
		catalog.TagsInput{Search: "blue hair", Offset: 5, Limit: 10, MinCount: 2, Exact: true},
	)

	if out.Search != "blue_hair" || out.Status != "ok" || out.SnapshotDate != "2026-09-24" {
		t.Errorf("output = %+v, want the result's search, status, and snapshot date", out)
	}

	if !out.Exact || out.Offset != 5 || out.Limit != 10 || out.MinCount != 2 {
		t.Errorf("output = %+v, want the echoed request", out)
	}

	if !out.More || out.AliasOf != "futanari" || out.Withheld != 2 || out.WithheldBest != 3 {
		t.Errorf("output = %+v, want the result's page and alias data", out)
	}

	if len(out.Synonyms) != 1 || out.Synonyms[0] != "pee" {
		t.Errorf("synonyms = %v, want [pee]", out.Synonyms)
	}

	if len(out.Tags) != 1 {
		t.Fatalf("tags = %+v, want one", out.Tags)
	}

	tag := out.Tags[0]
	if tag.Name != "blue_hair" || tag.Category != "general" || tag.Count != 482913 {
		t.Errorf("tag = %+v, want the name, category, and count", tag)
	}

	if !tag.CountIsTarget || tag.AliasOf != "futanari" || len(tag.Implications) != 1 {
		t.Errorf("tag = %+v, want the alias and implication data", tag)
	}
}

func TestTagsKeepsNilSlicesEmpty(t *testing.T) {
	out := Tags(catalog.TagsResult{}, catalog.TagsInput{})

	if out.Synonyms == nil {
		t.Error("synonyms = nil, want an empty slice")
	}

	if out.Tags == nil {
		t.Error("tags = nil, want an empty slice")
	}
}

func TestTagsEmptyImplicationsAreAnEmptySlice(t *testing.T) {
	result := catalog.TagsResult{Tags: []booru.Tag{{Name: "blue_hair"}}}

	out := Tags(result, catalog.TagsInput{})

	if out.Tags[0].Implications == nil {
		t.Error("implications = nil, want an empty slice")
	}
}

func TestRenderPage(t *testing.T) {
	out := TagsOutput{
		Offset:   0,
		Synonyms: []string{"pee", "urine"},
		More:     true,
		Tags: []TagOutput{
			{
				Name:         "blue_hair",
				Category:     "general",
				Count:        482913,
				Implications: []string{"aqua_hair"},
			},
			{
				Name:          "dickgirl",
				Category:      "general",
				Count:         52264,
				AliasOf:       "futanari",
				CountIsTarget: true,
			},
		},
	}

	got := Render(out)
	want := "- blue_hair (general): 482913 works [implies aqua_hair]\n" +
		"- dickgirl (general): alias of futanari (52264 works there)\n" +
		"\nAlso known as: pee, urine.\n" +
		"\nMore matches exist; continue with offset 2."

	if got != want {
		t.Errorf("Render() =\n%s\nwant\n%s", got, want)
	}
}

func TestRenderAliasWithoutTheTargetCount(t *testing.T) {
	out := TagsOutput{
		Tags: []TagOutput{{Name: "dickgirl", Category: "general", AliasOf: "futanari"}},
	}

	if got, want := Render(out), "- dickgirl (general): alias of futanari"; got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRenderHonestEmpties(t *testing.T) {
	tests := map[string]struct {
		out  TagsOutput
		want string
	}{
		"substring miss": {
			out:  TagsOutput{Search: "qzx", Status: string(catalog.StatusNoSubstringMatch)},
			want: `No tag name contains "qzx".`,
		},
		"one withheld": {
			out: TagsOutput{
				Search:       "ass_grabeye_contact",
				Status:       string(catalog.StatusNoSubstringMatch),
				MinCount:     1,
				Withheld:     1,
				WithheldBest: 0,
			},
			want: `No tag name contains "ass_grabeye_contact" with at least 1 works; ` +
				`one canonical match with 0 works was hidden by the floor. ` +
				`Set min_count=0 to include it.`,
		},
		"many withheld": {
			out: TagsOutput{
				Search:       "blue",
				Status:       string(catalog.StatusNoSubstringMatch),
				MinCount:     2,
				Withheld:     3,
				WithheldBest: 1,
			},
			want: `No tag name contains "blue" with at least 2 works; ` +
				`3 canonical matches were hidden by the floor (the best has 1 works). ` +
				`Set min_count=0 to include them.`,
		},
		"offset past end": {
			out: TagsOutput{
				Search:   "blue",
				Status:   string(catalog.StatusOffsetPastEnd),
				Offset:   50,
				Synonyms: []string{"pee"},
			},
			want: `Offset 50 is past the last result for "blue"; there is no further page. ` +
				`Try a smaller offset. Known aliases: pee.`,
		},
		"exact miss": {
			out:  TagsOutput{Search: "no_such_tag", Status: string(catalog.StatusExactNotFound)},
			want: `No tag named "no_such_tag" exists.`,
		},
		"unknown": {
			out: TagsOutput{Search: "maybe_tag", Status: string(catalog.StatusUnknown)},
			want: `Could not confirm whether "maybe_tag" exists because the tag index was ` +
				`unavailable; try again.`,
		},
		"unset status": {
			out:  TagsOutput{},
			want: "No tags matched.",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := Render(tt.out); got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTagsJSONShape(t *testing.T) {
	out := Tags(catalog.TagsResult{
		Search: "blue_hair",
		Status: catalog.StatusOK,
		Tags:   []booru.Tag{{Name: "blue_hair", Category: booru.CategoryGeneral, Count: 1}},
	}, catalog.TagsInput{Limit: 25, MinCount: 1})

	raw, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}

	for _, key := range []string{
		`"search"`, `"exact"`, `"status"`, `"snapshot_date"`, `"offset"`,
		`"limit"`, `"min_count"`, `"more"`, `"alias_of"`, `"withheld"`,
		`"withheld_best_count"`, `"synonyms"`, `"tags"`, `"name"`,
		`"category"`, `"count"`, `"count_is_target"`, `"implications"`,
		`"count_is_target"`, `"implications"`,
	} {
		if !strings.Contains(string(raw), key) {
			t.Errorf("JSON %s does not contain %s", raw, key)
		}
	}

	if !strings.Contains(string(raw), `"synonyms":[]`) ||
		!strings.Contains(string(raw), `"implications":[]`) {
		t.Errorf("JSON %s renders empty slices as null", raw)
	}
}
