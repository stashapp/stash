# tools/maintenance — off-queue library maintenance scripts

Dependency-free (Python stdlib) maintenance tools that talk to the stash GraphQL
API and run **outside** stash's job queue, so they work in parallel with whatever
stash is doing and can run on another machine or as a NAS cron. Both default to a
**dry-run / report** mode and only write when told to.

## tag_consolidate.py — merge duplicate / near-duplicate tags

Clusters tags and auto-merges the trivially-safe ones, while reporting the
judgement-call families for manual review.

- **SAFE (auto-merged with `--apply`):** case/punctuation/whitespace variants,
  en/em-dash differences, singular↔plural, `blond`↔`blonde`. Destination = the
  member with the most scenes; each merged name is preserved as an **alias** on the
  survivor (so future scraper matches still hit). Uses `tagsMerge`.
- **SUBJECTIVE (report-only, never auto-merged):** POV / parenthetical families
  (`Cowgirl`, `Cowgirl (DP)`, `Cowgirl - POV`…), hair-colour schemes, age brackets,
  and 1-edit near-duplicates (e.g. `Doggy Style`/`Doggystyle`, typos like
  `Ass Smacking`/`Ass Stacking`). Merge these yourself or via the assistant.

```bash
python3 tag_consolidate.py --stash-url http://overwatch-stash:9999            # dry-run report
python3 tag_consolidate.py --stash-url http://overwatch-stash:9999 --apply     # merge SAFE clusters
```

## generated_sweeper.py — orphan + integrity audit of generated/

Builds the set of live file hashes (scene + image fingerprints) from the API, walks
stash's `generated/` tree, and reports:
- **orphans** — artifacts whose hash is no longer in the library (deleted/re-hashed scenes),
- **zero-byte / corrupt** artifacts,
- **stale** `tmp/` / `download_stage/` leftovers,
- a **coverage** summary (how many previews / sprites / transcodes exist).

`--prune` MOVES orphans into `generated/_quarantine/` (reversible — never hard-deletes).
Run it on the NAS for a fast local walk:

```bash
python3 generated_sweeper.py --stash-url http://localhost:9999 \
    --generated-dir /volume1/docker/stash/generated            # report only
python3 generated_sweeper.py --stash-url http://localhost:9999 \
    --generated-dir /volume1/docker/stash/generated --prune     # quarantine orphans
```

Good as a periodic cron — its yield is low on a tidy library but it doubles as a
coverage dashboard and catches corruption early.

## Notes

- Both send a browser `User-Agent` (stash's stack rejects the default `Python-urllib`
  UA with an empty body) — see the same pattern in `../identify/`.
- Don't run `tag_consolidate.py --apply` or heavy maintenance at the same time as a
  native stash job over the same data — last-write-wins / SQLite lock contention.
