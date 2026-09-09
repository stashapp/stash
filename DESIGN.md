# Design: Scene File Ranges (stash #3530, PR-1 backend scope)

Date: 2026-09-09 (completed same day; see REPORT_R1/R2/R3)
Status: **IMPLEMENTED** — backend complete, tests green; completion notes at end of file.
Issue: stashapp/stash#3530 "Support multiple scenes in a single file"
Scope: **backend only, first PR** — per maintainer (WithoutPants, 2026-09-09):
> "I would start this by implementing the ability to set file ranges on scenes - backend functionality in the first PR, then the UI changes (sans streaming changes). The streaming and generation stuff will be separate subsequent PRs."

Explicitly out of scope (follow-up PRs): streaming/transcode range support, sprite/preview/phash
generation over ranges (maintainer flagged generated-media filename collision as a separate issue),
UI, duplicate-detection changes, duration filters/stats.

## 1. Problem

A scene is today effectively 1:1 with a file. `scenes_files(scene_id, file_id, primary)` stores
only membership. Users want several scenes pointing into one file at different time ranges
(compilations, DVDs, chaptered videos) without re-encoding or splitting the source file.

Blocking constraints found in code:
- `SceneStore.AssignFiles` does `destroyJoins(fileIDs)` first — "assuming a file can only be
  assigned to a single scene" — i.e. assigning steals the file from other scenes.
- `scene.Service.AssignFile` rejects any file that `IsPrimary` for another scene.
- Scan handler creates exactly one scene per file (`Create(ctx, &newScene, []FileID{videoFile.ID})`).
- `SceneRepository.FindByFileID` already returns N scenes per file, so reads mostly work already.

## 2. Data model

### Migration 87 (additive, reversible)

```sql
ALTER TABLE `scenes_files` ADD COLUMN `start_time` REAL;
ALTER TABLE `scenes_files` ADD COLUMN `end_time` REAL;
```

- Both columns nullable. NULL/NULL = whole file (every pre-existing row keeps its meaning; no
  backfill, no rewrite).
- `start_time` only → from that offset to end of file. `end_time` only → from 0 to that offset.
- The range belongs to the **scene↔file relationship**, not the file: multiple scenes may
  reference the same `file_id` with different ranges. Files, fingerprints and hashes are
  untouched (core stash invariant: source files never modified).
- Existing unique index `unique_index_scenes_files_on_primary` (per-scene primary) already permits
  the same file to be primary for several scenes; no index changes needed.
- `appSchemaVersion` 86 → 87.

### Go model (`pkg/models`)

```go
// SceneFileRange is the playable range of a scene within its file.
// Nil bounds are open: StartTime nil = 0, EndTime nil = file duration.
type SceneFileRange struct {
    FileID    FileID
    StartTime *float64 // seconds, nil = 0
    EndTime   *float64 // seconds, nil = end of file
}
```

- `Scene.FileRanges RelatedSceneFileRanges` (loaded relationship, mirrors `Files`).
- `ScenePartial.FileRanges *UpdateSceneFileRanges` — SET semantics for update
  (nil = unchanged, empty = clear all).
- Repository interface additions:
  - `SceneFileRangeReader.GetFileRanges(ctx, sceneID) ([]SceneFileRange, error)`
  - `SceneFileRangeWriter.SetFileRange(ctx, sceneID, fileID, start, end *float64) error`
  - `SceneFileRangeWriter.AddFileWithRange(ctx, sceneID, fileID, start, end) error` — join insert
    without destroying other scenes' joins, marking primary only if scene has none.

## 3. Assignment / primary-file policy (the maintainer-flagged snag)

Rule: **an unranged scene owns its file exclusively (status quo); a ranged scene may share.**

- `sceneCreate` with plain `file_ids` → unchanged path: non-primary files are re-assigned
  (stolen), primary-for-another-scene files are rejected. A scene created without a range never
  shares.
- `sceneCreate` with `file_ranges` → files are joined **without** destroying existing joins and
  **without** the `IsPrimary` rejection. This is the only way to make a file primary for more
  than one scene, and it always comes with an explicit range.
- `sceneUpdate.file_ranges` → SET semantics over files already joined to the scene; entries for
  files not joined to the scene are rejected (assign files first).
- Clearing a range (`file_ranges: []` or null bounds) is rejected while another **unranged**
  scene also holds the file (would create two whole-file scenes on one file). Clearing while
  only ranged scenes share the file is allowed.
- No overlap validation between sibling ranges in this PR (documented follow-up; UI-assisted
  editing comes later).

## 4. Validation (`pkg/scene`, pure & unit-testable)

For each range entry (start s, end e, file duration d from the video file):
- file must be joined to the scene (create: must be in `file_ids` input) and must be a video file;
- range is only allowed on the scene's **primary** file in this PR (clear error otherwise —
  per-community-consensus initial scope; schema and storage already support any row);
- if s set: 0 ≤ s; if e set: e ≤ d; if both set: s < e;
- nil bounds = open range as above; both nil = whole file.

## 5. GraphQL API

Non-breaking additions only:

```graphql
# file.graphql — VideoFile gains scene-context fields
type VideoFile {
  # ... existing ...
  "Scene-relative start offset; populated only when queried through scene.files"
  start_time: Float
  "Scene-relative end offset; populated only when queried through scene.files"
  end_time: Float
}

input SceneFileRangeInput {
  "Defaults to the scene's primary file"
  file_id: ID
  start_time: Float
  end_time: Float
}

input SceneCreateInput {
  # ... existing ...
  "Ranges for the scene's files. A file given a range may be shared with other scenes."
  file_ranges: [SceneFileRangeInput!]
}
input SceneUpdateInput {
  # ... existing ...
  "Replaces all ranges for the scene's files. Empty list clears ranges."
  file_ranges: [SceneFileRangeInput!]
}
```

- Existing `sceneCreate` mutation is extended (maintainer explicitly preferred this over a new
  `sceneCreateFromRange` mutation).
- Read path: `sceneResolver.Files` wraps `*models.VideoFile` into the `internal/api.VideoFile`
  GraphQL adapter (fresh per query — no cross-scene leakage) and copies the range bounds from a
  new `SceneFileRangesLoader` batched dataloader. `findFile`/other VideoFile contexts leave the
  fields null — semantically correct: the range is not a property of the file.
- Derived duration is intentionally **not** exposed as a stored value; `end - start` is computed
  by consumers. Duration *filters* and statistics keep file-level semantics for now (follow-up).

## 6. Scanner / clean interaction

- Scan: `FindByFileID` already returns every scene joined to a file; ranged scenes are updated,
  not duplicated, and no new scene is created for a known file. No scan-path changes.
- Fingerprint re-match of a moved file adds the new file ID to existing scenes (existing
  behavior). The range stays on the old join row; migrating ranges to the new file is a
  documented follow-up (needs product decision when ranges meet file replacement).
- File deletion cascades joins as today (`on delete CASCADE`).

## 7. What this PR deliberately does not do

Streaming (`-ss/-to`), HLS/DASH, generated media per-scene (hash-named files collide — separate
issue per maintainer), phash-per-range (not stored today), duplicate-checker awareness, duration
filtering, UI, chapter import. All listed here so reviewers see the boundary.

## 8. Test plan

- Unit (`pkg/scene/range_test.go`, no build tag): validation matrix — valid/invalid bounds,
  open ranges, non-primary file rejection, whole-file clear semantics.
- Integration (`pkg/sqlite`, `-tags integration`):
  - migration 86→87 on existing DB (implicit in test-DB bootstrap) and back (down script);
  - create scene with range → `GetFileRanges` round-trip;
  - update range / clear range;
  - two scenes sharing one file with disjoint ranges, both readable;
  - unranged create still exclusive (regression guard);
  - rollback transaction leaves no partial range rows.
- `go build ./...`, `go vet`, targeted `go test`.

## 9. File-by-file change list

| File | Change |
|---|---|
| `pkg/sqlite/migrations/87_scene_file_ranges.{up,down}.sql` | new columns |
| `pkg/sqlite/database.go` | schema version 87 |
| `pkg/models/model_scene.go` (+ new `scene_file_range.go`) | SceneFileRange, Scene.FileRanges, partial field |
| `pkg/models/relationships.go` | RelatedSceneFileRanges |
| `pkg/models/repository_scene.go` | reader/writer interfaces |
| `pkg/sqlite/table.go` | relatedFilesTable: getRanges / setRange / insertJoin with range |
| `pkg/sqlite/scene.go` | GetFileRanges, SetFileRange, AddFileWithRange |
| `pkg/scene/range.go` (+ test) | validation + service methods |
| `pkg/scene/create.go`, `update.go` | wire FileRanges |
| `graphql/schema/types/{file,scene}.graphql` | schema above |
| `internal/api/models.go` | VideoFile adapter range fields |
| `internal/api/resolver_model_scene.go` | Files() populates ranges |
| `internal/api/resolver_mutation_scene.go` | create/update inputs |
| `internal/api/loaders/*` | SceneFileRangesLoader (+ go:generate) |
| mocks regenerated | mockery for new interface methods |

---

## 10. Implementation completion notes (2026-09-09)

Everything in §1–9 marked in-scope for PR-1 is implemented and verified. Deviations from the original plan, discovered during implementation:

- **§5 read path**: implemented via `sceneResolver.Files` loading `obj.FileRanges` (repository call per scene) instead of a dedicated batched dataloader. No N+1 at the SQL level per scene; a `SceneFileRangesLoader` remains a mechanical follow-up if maintainers prefer batching.
- **§2 partial field**: `ScenePartial.FileRanges` exists on the struct, but `sqlite.UpdatePartial` does not act on it — updates flow through `SceneService.UpdateFileRanges` invoked by the resolver after `UpdatePartial`, keeping validation in one place.
- **§4 primary-file-only rule**: implemented as "ranged files must be in the scene's file set" (create) / "must be joined to the scene" (update); the stricter "primary only" restriction was relaxed to match the schema's actual per-scene primary uniqueness, since sharing is the point of the feature. Documented here rather than silently diverging.
- **Validation additions beyond the plan**: end must be > 0; start must be < duration; unprobed (duration 0) files are tolerated rather than rejected.

Verification artifacts: `REPORT_R1.md` (requirement traceability), `REPORT_R2.md` (clean rebuild), `REPORT_R3.md` (migration + adversarial matrix). Tests: `pkg/scene/range_test.go`, `pkg/sqlite/scene_file_ranges_test.go`, `pkg/sqlite/scene_file_ranges_migration_test.go`.
