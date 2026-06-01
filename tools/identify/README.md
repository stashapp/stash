# external_identify.py — parallel / off-box scene identification

Identify stash scenes against your stash-box endpoints (StashDB, ThePornDB, …)
**without using stash's internal job queue**, so it runs in parallel with a long
Generate→Phash (or on a different machine entirely).

## Why
stash runs Tasks sequentially — a native **Identify** queues behind a running
**Generate→Phash** (which can take hours). This script talks directly to the stash
GraphQL API and the stash-box GraphQL API (`findScenesBySceneFingerprints`), so it's
independent of stash's job queue. It matches by **fingerprint** (oshash + phash);
since oshash is computed at scan time, it works with no phash and no waiting.

## Requirements
- Python 3 (standard library only — no pip installs).
- Network access to stash and to the stash-box endpoints.
- If stash has auth enabled, a stash **API key** (Settings → Security → API Key).
  Run on the stash host (`localhost`) and you typically need none.

## Usage
```bash
# safe preview (default): shows what would be matched/changed, writes nothing
python3 external_identify.py --stash-url http://localhost:9999

# from another machine over Tailscale (with an API key if auth is on)
python3 external_identify.py --stash-url http://overwatch-stash:9999 --stash-api-key <KEY>

# actually apply, and mark matched scenes organized
python3 external_identify.py --stash-url http://localhost:9999 --apply --set-organized
```
It reads your configured stash-box endpoints **and their API keys** straight from
stash, and queries them in the order stash has them (first match wins).

### Options
| flag | effect |
|---|---|
| (default) | **dry-run** — preview only, no writes |
| `--apply` | write changes via `sceneUpdate` |
| `--all-scenes` | consider all scenes (default: only `organized: false`) |
| `--set-organized` | mark matched scenes organized |
| `--allow-multiple` | apply the first match when a scene has several (default: skip ambiguous) |
| `--batch N` | scenes per stash-box query (default 40) |
| `--limit N` | cap scenes processed (good for a first test) |

## What it writes (MERGE-like, non-destructive)
For each single fingerprint match it fills only what's **missing**: title / date /
details (only if empty), studio and performers (only if the scene has none —
find-or-create by name, stamping the stash-box `stash_id`), always adds the scene's
stash-box `stash_id`, and optionally sets `organized`. It does **not** overwrite
existing single-value fields. **Tags are added by default** with a MERGE strategy —
the match's tags are added on top of the scene's existing tags (never replacing them),
creating any missing tag (stamped with its stash-box id). Pass `--no-tags` to disable.
Images are still not handled.

## Important caveats
- **Don't run this at the same time as a native Identify** over the same library —
  you'd double the stash-box queries and race on writes.
- **Slightly lower fidelity than stash's native Identify**: name-based performer/studio/
  tag linking (vs stash's richer stash-id + alias matching), no images. Tags ARE applied
  by default (merge + create-missing). Prefer native Identify when you can run it; use
  this for parallel/off-box runs.
- **Rate limits**: it paces requests from each endpoint's `maxRequestsPerMinute`
  (defaults to a gentle ~4/s when unset). ThePornDB is stricter than StashDB.
- Always do a `--dry-run` pass first and skim the matches.

## Companion: `external_filename_parse.py` (no stash-box needed)

For the long tail of scenes that **no stash-box can match**, this companion links them
to metadata parsed from their **path** — but **LINK-ONLY**: it matches the path against
your *existing* studios/performers (multi-word names/aliases, whole-word) and never
creates new records, so it can't pollute the library. MERGE semantics: fills an empty
studio, merges performers on top of existing, fills an empty date, and (with
`--set-title`) a cleaned title. Single-token names are deliberately ignored (a first-name
alias would mislink). Lower fidelity than fingerprint identify — `--dry-run` is the
default; review before `--apply`.

```bash
python3 external_filename_parse.py --stash-url http://overwatch-stash:9999            # dry-run
python3 external_filename_parse.py --stash-url http://overwatch-stash:9999 --apply
```

## Companion: `external_identify_performers.py` — enrich **performers** (not scenes)

The scene tools above link *scenes*. This one enriches **performers** with canonical
metadata (birthdate, ethnicity, eye/hair colour, measurements, country, aliases, URLs,
image, …) pulled straight from your stash-box endpoints (StashDB, ThePornDB) — again
**outside** stash's job queue. It already sends the browser `User-Agent` that ThePornDB's
Cloudflare requires (the default `Python-urllib` UA gets a 403 / error 1010).

It reads your configured stash-box endpoints and their API keys from stash, and is a
faithful port of stash's own `pkg/stashbox/performer.go` translation logic, so the values
it writes match what a native performer scrape would produce.

### Modes
| mode | what it does |
|---|---|
| `refresh` (default) | for **linked** performers (those that already have a `stash_id`) that are **missing some fields**, fetch the canonical record via `findPerformer(id)` and **fill only the empty fields** |
| `search` | for **unlinked** performers, run `searchPerformer(term=name)`, post-filter to a case-insensitive **exact** name match, and link/merge the single survivor |
| `both` | refresh, then search |

### Usage
```bash
# safe preview (default is dry-run): refresh linked performers' empty fields
python3 external_identify_performers.py --stash-url http://overwatch-stash:9999

# find linkage for performers that have no stash_id yet
python3 external_identify_performers.py --stash-url http://overwatch-stash:9999 --mode search

# refresh THEN search, small first pass to eyeball the matches
python3 external_identify_performers.py --stash-url http://overwatch-stash:9999 --mode both --limit 20

# actually write (after reviewing a dry-run)
python3 external_identify_performers.py --stash-url http://overwatch-stash:9999 --mode both --apply
```

### Options
| flag | effect |
|---|---|
| (default) | **dry-run** — preview only, no writes |
| `--apply` | write changes via `performerUpdate` |
| `--mode {refresh,search,both}` | which candidate set to process (default `refresh`) |
| `--limit N` | cap performers processed (good for a first test; `0` = no cap) |
| `--allow-multiple` | search mode: take the first remaining result after the exact-name filter instead of skipping ambiguous |
| `--per-page N` | stash pagination page size (default 100) |
| `--stash-api-key KEY` | stash API key, if auth is on (omit when running on the stash host) |

### What it writes (MERGE + union, non-destructive)
Single-value fields are filled **only when stash's value is empty** — it never overwrites
an existing value (e.g. `height_cm` only when stash has none, `height > 0` skip). Multi-value
fields (`alias_list`, `urls`, `stash_ids`) are **unioned** with the existing values before
sending, because stash's update replaces these wholesale — so existing aliases/URLs and
stash-ids from *other* endpoints are preserved, never dropped. Aliases case-equal to the
performer's own name are filtered out. Images are pre-`HEAD`ed and skipped if the CDN URL is
dead, so a bad image can't roll back the whole update. A `[noop]` line means the stash-box
record had nothing new to merge.

### Caveats
- **Don't run it concurrently with a native stash performer scrape** — both issue
  `performerUpdate` in a single txn and it's last-write-wins.
- It paces stash-box requests off each endpoint's `maxRequestsPerMinute` (gentle ~4/s when
  unset); ThePornDB is stricter than StashDB.
- Always do the default **dry-run** first and skim the `[match] … would set [...]` lines
  before adding `--apply`.

### NAS cron (periodic enrichment)
The canonical NAS copy lives next to the scene tool at
`/volume1/docker/stash/identify/` (or `~shadowshark/stash-tools/`). The NAS host has
**Python 3.8** and reaches stash at `http://localhost:9999` with no API key. Push the
script the same way as the others:

```bash
# from this repo (tar-over-ssh; the NAS has no scp/sftp):
tar cf - -C tools/identify external_identify_performers.py \
  | ssh -p 3239 shadowshark@overwatch-stash 'tar xf - -C /volume1/docker/stash/identify'
```

Recommended cron line — a **weekly** dry-run-then-apply refresh+search pass, logged so you
can eyeball what it touched (Synology has no per-user `crontab`; add this via **DSM →
Control Panel → Task Scheduler** as a *Scheduled Task → User-defined script*, run as
`shadowshark`, or append it to root's `/etc/crontab` — note `/etc/crontab` needs the
**user column**):

```cron
# m h dom mon dow  user        command
  30 4 *   *   1    shadowshark /usr/bin/python3 /volume1/docker/stash/identify/external_identify_performers.py --stash-url http://localhost:9999 --mode both --apply >> /volume1/docker/stash/identify/enrich_performers.log 2>&1
```

Drop `--apply` for a report-only schedule. New performers only appear after you've scanned
new media, so weekly (or right after a big import) is plenty — this is low-yield on an
already-enriched library but cheap to run.
