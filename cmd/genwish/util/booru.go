package util

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/booru"
	"github.com/wishmatic/genwish/internal/catalog"
	"github.com/wishmatic/genwish/internal/cli"
)

const defaultMinCount = 1

type booruOptions struct {
	offset   int
	limit    int
	minCount int
	exact    bool
	category []string
	json     bool
}

func NewBooruCommand(deps cli.Deps) *cobra.Command {
	var opts booruOptions

	cmd := &cobra.Command{
		Use:   "booru <search...>",
		Short: "Search Danbooru tag names",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return runBooruSearch(cmd.Context(), deps, opts, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.offset, "offset", 0, "how many matching tags to skip")
	flags.IntVar(&opts.limit, "limit", catalog.DefaultTagLimit, "how many matching tags to return")
	flags.IntVar(
		&opts.minCount,
		"min-count",
		defaultMinCount,
		"lowest work count a substring match may have",
	)
	flags.BoolVar(
		&opts.exact,
		"exact",
		false,
		"search is the whole tag name and aliases resolve to their target",
	)
	flags.StringArrayVar(
		&opts.category,
		"category",
		nil,
		fmt.Sprintf("keep only this category; repeatable (%s)", booru.CategoryList()),
	)
	flags.BoolVar(&opts.json, "json", false, "print JSON instead of text")

	return cmd
}

func runBooruSearch(ctx context.Context, deps cli.Deps, opts booruOptions, args []string) error {
	search := strings.TrimSpace(strings.Join(args, " "))
	if search == "" {
		return errors.New("search is empty")
	}

	categories, err := parseCategories(opts.category)
	if err != nil {
		return err
	}

	service, err := cli.CatalogService(deps)
	if err != nil {
		return err
	}

	request := catalog.TagsInput{
		Search:     search,
		Offset:     opts.offset,
		Limit:      opts.limit,
		MinCount:   opts.minCount,
		Exact:      opts.exact,
		Categories: categories,
	}

	result, err := service.Tags(ctx, request)
	if err != nil {
		return fmt.Errorf("tags: %w", err)
	}

	out := newTagsOutput(result, request)

	if opts.json {
		return cli.WriteJSON(deps.Out, out)
	}

	_, err = fmt.Fprintln(deps.Out, renderTags(out))

	return err
}

func parseCategories(names []string) ([]booru.TagCategory, error) {
	categories := make([]booru.TagCategory, 0, len(names))

	for _, name := range names {
		category, ok := booru.ParseTagCategory(name)
		if !ok {
			return nil, fmt.Errorf("category %q is not one of %s", name, booru.CategoryList())
		}

		categories = append(categories, category)
	}

	return categories, nil
}
