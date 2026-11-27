# Booru

`genwish util booru` searches Danbooru tag names and prints text or JSON.

```sh
genwish util booru blue_hair
genwish util booru blue hair --limit 5 --category general
genwish util booru blue_hair --json
```

Flags:

| Flag          | Default | Purpose                                                               |
| ------------- | ------- | --------------------------------------------------------------------- |
| `--offset`    | `0`     | How many matching tags to skip.                                       |
| `--limit`     | `25`    | How many matching tags to return.                                     |
| `--min-count` | `1`     | Lowest work count a match may have; `0` includes zero-work tags.      |
| `--exact`     | `false` | The search is the whole tag name and aliases resolve to their target. |
| `--category`  | all     | Keep only this category; repeatable.                                  |
| `--json`      | `false` | Print JSON instead of text.                                           |

The search is a literal substring of the tag name, so `blue hair` and `blue_hair` are the same
search.

## Results

Results are ordered by work count, name ascending on ties, and form one page of the match list.
Walk the list with `--offset` and stop when `more` is false.

Substring results exclude tags below `--min-count` works by default because Danbooru's tag table
also holds concatenated and punctuation-mangled zero-work names. `withheld` and
`withheld_best_count` report what the floor removed. `--min-count 0` returns them.
Counts are read live, and `snapshot_date` records the date they were read.

An empty page reports a status:

- `ok`: matches were returned.
- `no_substring_match`: the substring genuinely matched nothing.
- `exact_not_found`: the exact name is absent.
- `offset_past_end`: the offset ran past the last match.
- `unknown`: an unanswerable case was reached.

Use `--exact` to ask whether a tag exists; a substring result describes matches only.

## Authentication

Danbooru answers anonymous requests. Logging in raises the rate limit:

```sh
genwish login danbooru
```

The command prompts for the Danbooru login first, then the API key. `genwish logout danbooru`
deletes the stored pair.

## Configuration

See [configuration.md](configuration.md). The rate limit is configured with the `booru.rate_limit`
table (`rps`, `burst`, `timeout`). Errors always include a short upstream body excerpt. The tag
page size is fixed at 100 and the largest offset at 1000. The user agent is hardcoded.
