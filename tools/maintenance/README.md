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

## scene_dedup.py — find near-duplicate scenes (phash)

Clusters scenes whose perceptual hashes are within a Hamming distance (`--max-distance`,
default 8), picks a keeper (highest resolution → largest file → lowest id), and reports the
rest as dup-candidates with reclaimable space. `--apply` tags non-keepers `_dupe-candidate`
(merge, never replace) for review/deletion. Dry-run default.

```bash
python3 scene_dedup.py --stash-url http://overwatch-stash:9999                 # report
python3 scene_dedup.py --stash-url http://overwatch-stash:9999 --apply          # tag dup-candidates
```

## db_backup.py — consistent metadata-DB + config backup (run on the NAS)

Online-backup (sqlite3 `.backup` API — safe while stash runs) of `stash-go.sqlite` + `config.yml`
→ timestamped gzip tar in `--out`, keeps the most recent `--keep` (default 14), optional off-NAS
`--remote user@host:/path` rsync. **Dry-run by default; pass `--apply` to write.** Recommended NAS
cron (DSM Task Scheduler, nightly): `db_backup.py --apply`.

## ingest_pipeline.py — off-queue metadata orchestrator

Runs the chain in the only safe order — `external_identify.py` (fingerprint) **then**
`external_filename_parse.py` (link-only) for the leftovers (they must never run together) — with
before/after gap snapshots. The single cron entry point to keep a growing library's metadata current.

```bash
python3 ingest_pipeline.py --stash-url http://localhost:9999            # dry-run both steps
python3 ingest_pipeline.py --stash-url http://localhost:9999 --apply     # apply
```

## tag_consolidate.py --emit-plan

In addition to its SAFE auto-merge, `tag_consolidate.py --emit-plan` outputs a structured JSON merge
plan of the **subjective** families (POV/hair/age/1-edit) — destination + sources + confidence — for the
in-app assistant or a human to execute via `tagsMerge`. Read-only.

## Notes

- Both send a browser `User-Agent` (stash's stack rejects the default `Python-urllib`
  UA with an empty body) — see the same pattern in `../identify/`.
- Don't run `tag_consolidate.py --apply` or heavy maintenance at the same time as a
  native stash job over the same data — last-write-wins / SQLite lock contention.
