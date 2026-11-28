package util

import (
	"fmt"
	"strings"

	"github.com/wishmatic/genwish/internal/booru"
	"github.com/wishmatic/genwish/internal/catalog"
)

type tagsOutput struct {
	Search       string      `json:"search"`
	Exact        bool        `json:"exact"`
	Status       string      `json:"status"`
	SnapshotDate string      `json:"snapshot_date"`
	Offset       int         `json:"offset"`
	Limit        int         `json:"limit"`
	MinCount     int         `json:"min_count"`
	More         bool        `json:"more"`
	AliasOf      string      `json:"alias_of"`
	Withheld     int         `json:"withheld"`
	WithheldBest int         `json:"withheld_best_count"`
	Synonyms     []string    `json:"synonyms"`
	Tags         []tagOutput `json:"tags"`
}

type tagOutput struct {
	Name          string   `json:"name"`
	Category      string   `json:"category"`
	Count         int      `json:"count"`
	CountIsTarget bool     `json:"count_is_target"`
	AliasOf       string   `json:"alias_of"`
	Implications  []string `json:"implications"`
}

func newTagsOutput(result catalog.TagsResult, request catalog.TagsInput) tagsOutput {
	return tagsOutput{
		Search:       result.Search,
		Exact:        request.Exact,
		Status:       string(result.Status),
		SnapshotDate: result.SnapshotDate,
		Offset:       request.Offset,
		Limit:        request.Limit,
		MinCount:     request.MinCount,
		More:         result.More,
		AliasOf:      result.AliasOf,
		Withheld:     result.Withheld,
		WithheldBest: result.WithheldBest,
		Synonyms:     nonNilStrings(result.Synonyms),
		Tags:         toTagOutputs(result.Tags),
	}
}

func renderTags(out tagsOutput) string {
	if len(out.Tags) == 0 {
		return renderEmpty(out)
	}

	var b strings.Builder

	for _, tag := range out.Tags {
		fmt.Fprint(&b, renderTag(tag))

		if len(tag.Implications) > 0 {
			fmt.Fprintf(&b, " [implies %s]", strings.Join(tag.Implications, ", "))
		}

		b.WriteString("\n")
	}

	if len(out.Synonyms) > 0 {
		fmt.Fprintf(&b, "\nAlso known as: %s.\n", strings.Join(out.Synonyms, ", "))
	}

	if out.More {
		fmt.Fprintf(
			&b,
			"\nMore matches exist; continue with offset %d.\n",
			out.Offset+len(out.Tags),
		)
	}

	return strings.TrimRight(b.String(), "\n")
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}

	return values
}

func toTagOutputs(tags []booru.Tag) []tagOutput {
	out := make([]tagOutput, 0, len(tags))

	for _, tag := range tags {
		out = append(
			out,
			tagOutput{
				Name:          tag.Name,
				Category:      tag.Category.String(),
				Count:         tag.Count,
				CountIsTarget: tag.CountIsTarget,
				AliasOf:       tag.AliasOf,
				Implications:  nonNilStrings(tag.Implications),
			},
		)
	}

	return out
}

func renderTag(tag tagOutput) string {
	if tag.AliasOf == "" {
		return fmt.Sprintf("- %s (%s): %d works", tag.Name, tag.Category, tag.Count)
	}

	if tag.CountIsTarget {
		return fmt.Sprintf(
			"- %s (%s): alias of %s (%d works there)",
			tag.Name,
			tag.Category,
			tag.AliasOf,
			tag.Count,
		)
	}

	return fmt.Sprintf("- %s (%s): alias of %s", tag.Name, tag.Category, tag.AliasOf)
}

func renderEmpty(out tagsOutput) string {
	switch out.Status {
	case string(catalog.StatusNoSubstringMatch):
		return renderNoSubstringMatch(out)

	case string(catalog.StatusOffsetPastEnd):
		return fmt.Sprintf(
			"Offset %d is past the last result for %q; there is no further page. "+
				"Try a smaller offset.%s",
			out.Offset,
			out.Search,
			renderSynonyms(out),
		)

	case string(catalog.StatusExactNotFound):
		return fmt.Sprintf("No tag named %q exists.", out.Search)

	case string(catalog.StatusUnknown):
		return fmt.Sprintf(
			"Could not confirm whether %q exists because the tag index was unavailable; try again.",
			out.Search,
		)
	}

	return "No tags matched."
}

func renderSynonyms(out tagsOutput) string {
	if len(out.Synonyms) == 0 {
		return ""
	}

	return fmt.Sprintf(" Known aliases: %s.", strings.Join(out.Synonyms, ", "))
}

// renderNoSubstringMatch names the filter that emptied the page. A withheld canonical tag exists;
// saying only that the substring matched nothing would assert that it does not.
func renderNoSubstringMatch(out tagsOutput) string {
	if out.Withheld == 0 {
		return fmt.Sprintf("No tag name contains %q.%s", out.Search, renderSynonyms(out))
	}

	if out.Withheld == 1 {
		return fmt.Sprintf(
			"No tag name contains %q with at least %d works; "+
				"one canonical match with %d works was hidden by the floor. "+
				"Set min_count=0 to include it.%s",
			out.Search,
			out.MinCount,
			out.WithheldBest,
			renderSynonyms(out),
		)
	}

	return fmt.Sprintf(
		"No tag name contains %q with at least %d works; "+
			"%d canonical matches were hidden by the floor (the best has %d works). "+
			"Set min_count=0 to include them.%s",
		out.Search,
		out.MinCount,
		out.Withheld,
		out.WithheldBest,
		renderSynonyms(out),
	)
}
