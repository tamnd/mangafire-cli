# mangafire

Browse the [MangaFire](https://mangafire.to) manga catalog from the command line.

`mangafire` is a single pure-Go binary. No API key required.

## Install

```bash
go install github.com/tamnd/mangafire-cli/cmd/mangafire@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/mangafire-cli/releases), or run the container image:

```bash
docker run --rm ghcr.io/tamnd/mangafire:latest --help
```

## Usage

```bash
# List top manga by most viewed
mangafire list

# List 60 most viewed manga (2 pages)
mangafire list -p 2 -o table

# List latest updated manhwa
mangafire list --sort latest_updated --type manhwa -n 20

# Sort options: most_viewed, latest_updated, new_release, title_az
mangafire list --sort new_release -n 10

# Output formats
mangafire list -o json
mangafire list -o csv -n 50
```

## Commands

| Command | Description |
|---------|-------------|
| `list` | List manga from the MangaFire catalog |
| `version` | Show version information |

## List flags

```
-s, --sort string    sort order: most_viewed|latest_updated|new_release|title_az (default "most_viewed")
-t, --type string    filter by type: manga|manhwa|manhua|novel|one-shot|doujinshi
-p, --pages int      number of pages to fetch (30 items per page) (default 1)
```

## Global flags

```
-o, --output string    output format: table|json|jsonl|csv|tsv|url|raw (default "auto")
-n, --limit int        limit number of records (0 = all on fetched pages)
    --fields strings   comma-separated columns to include
    --no-header        omit header row
    --template string  Go text/template per record
    --timeout duration per-request timeout (default 30s)
    --delay duration   minimum spacing between requests
    --retries int      retry attempts on 429/5xx (default 3)
```

## License

Apache-2.0. See [LICENSE](LICENSE).
