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
