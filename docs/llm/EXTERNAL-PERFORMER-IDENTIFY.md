# External performer identifier — design

Companion to `external_identify.py` (which identifies **scenes**). This tool enriches **performers**
with full metadata from StashDB / ThePornDB, running outside stash's job queue. No code yet — this is
the spec.

Same architectural pattern as the scene identifier: stdlib Python, talks to stash GraphQL + the
configured stash-box endpoints, applies MERGE-like updates. The bundled fragment query at
`graphql/stash-box/query.graphql` (`PerformerFragment`) already returns every field we need.

---

## 1. Why this is needed

When `external_identify.py` matches a scene with new performers, it creates them via
`performerCreate` with **only `name` + `stash_id`** — no birthdate, no images, no measurements, no
aliases. The result is hundreds of bare performer pages.

Native stash's `scrapeSinglePerformer` + manual UI tagger works but:
- Sequential (one performer at a time, behind whatever's in the job queue)
- No batching against stash-box endpoints
- Per-scene "Identify" task doesn't populate performer fields

This tool fixes that: a one-shot batch enrich for all bare performers, running in parallel with
anything stash is doing.

---

## 2. Two operating modes

| Mode | Target performers | stash-box query | Ambiguity |
|---|---|---|---|
| **`refresh`** | Already linked (`stash_ids` not empty) but fields empty | `findPerformer(id: <stash_id>)` — direct, by id | None |
| **`search`** | Unlinked (no `stash_id`) | `searchPerformer(term: <name>)` | Multi-match possible |

`refresh` mode is the high-yield, zero-risk first pass — every match is unambiguous because we have
the canonical id. `search` mode handles the long tail (older performers, ones imported without
stash-box links) at higher ambiguity risk.

---

## 3. Architecture

```
stash GraphQL                  stash-box GraphQL                stash GraphQL
   findPerformers   ─→ for each ─→  findPerformer(id)   ─→ MERGE diff ─→  performerUpdate
                                    searchPerformer(term)
```

- One Python script (stdlib only) — proposed name `external_identify_performers.py`.
- Reuses `gql()` helper, `Stash` / `StashBox` classes from `external_identify.py`.
- `--dry-run` default, `--apply` to write, `--mode refresh|search|both`.
- Sources in stash-box priority order (StashDB first, TPDB fallback) — same as scene identifier.

---

## 4. Field mapping (stash-box `PerformerFragment` → stash `PerformerUpdateInput`)

Verbatim from `graphql/stash-box/query.graphql` (already in repo) and
`graphql/schema/types/performer.graphql:107`. Type translations called out where needed:

| Stash-box field | Type | → | Stash field | Type | Translation |
|---|---|---|---|---|---|
| `name` | String | | `name` | String | direct |
| `disambiguation` | String | | `disambiguation` | String | direct |
| `aliases` | [String!] | | `alias_list` | [String!] | direct (dupes auto-filtered server-side) |
| `gender` | GenderEnum | | `gender` | GenderEnum | direct (enums align: MALE/FEMALE/…) |
| `urls[].url` | [URL] | | `urls` | [String!] | flatten to URL strings |
| `images[0].url` | [Image] | | `image` | String | pass URL — stash field accepts URL OR base64 |
| `birth_date` | String | | `birthdate` | String | direct (rename only) |
| `death_date` | String | | `death_date` | String | direct |
| `ethnicity` | String | | `ethnicity` | String | direct |
| `country` | String | | `country` | String | direct |
| `eye_color` | String | | `eye_color` | String | direct |
| `hair_color` | String | | `hair_color` | String | direct |
| `height` | Int (cm) | | `height_cm` | Int | direct (rename only) |
| `measurements{band_size,cup_size,waist,hip}` | struct | | `measurements` | String | **format as `"<band><cup>-<waist>-<hip>"`** (e.g. `32C-26-36`) |
| `career_start_year` | Int | | `career_start` | String | stringify (`"2018"`) |
| `career_end_year` | Int | | `career_end` | String | stringify |
| `tattoos[]{location,description}` | [BodyMod] | | `tattoos` | String | **join as `"<location>: <description>"` newline-separated** |
| `piercings[]{location,description}` | [BodyMod] | | `piercings` | String | same as tattoos |
| `breast_type` | enum | | `fake_tits` | String | map: `NATURAL→"No"`, `FAKE→"Yes"`, `NA→""` (best-effort; verify against stash UI strings) |
| `id` (stash-box id) | ID | | `stash_ids` | [StashIDInput!] | append `{endpoint, stash_id: id}` if not already present |

**MERGE semantics**: only write fields where stash's current value is empty/null. Always append the
stash-box `stash_id` if missing. Never overwrite existing user edits.

**Open question on measurements format** — stash UI accepts arbitrary strings, but multiple
conventions exist (`32C-26-36`, `32-26-36`, etc.). The format proposed matches the most common
StashDB convention. Worth verifying against a few real performers before committing.

---

## 5. Field-by-field decisions worth flagging

- **Aliases**: union with existing stash aliases (stash server dedupes case-insensitively per the
  schema comment at line 128). Don't drop user-added aliases.
- **Images**: stash-box performer `images` is an array; take `images[0]` (StashDB returns the
  primary first). Passing the URL means stash downloads it server-side — saves bandwidth on our end.
- **Empty stash-box fields**: many are nullable. Skip the update for any field stash-box returned
  null/empty.
- **Birthdate / death_date format**: stash-box returns ISO `YYYY-MM-DD`; stash accepts the same.

---

## 6. Phase plan

**Phase 1 — `refresh` mode.** Lowest risk, highest yield. Performers with stash_ids but missing
fields. No ambiguity (direct id lookup). Validates the whole MERGE + apply pipeline.

**Phase 2 — `search` mode.** Unlinked performers. Per-name search with multi-match handling.
Defaults to skipping ambiguous; `--allow-multiple` takes first match (parallels the scene tool's
flag). Tag ambiguous performers (`Needs Review`) for manual triage.

**Phase 3 — assistant tool**. Mirror `identify_scenes_fast`: a new `identify_performers_fast`
assistant tool that shells out to the Python script. Same UI integration; same approve/decline
write-policy semantics; same `100s` exec timeout (performer enrichment is much cheaper than scene
identify — only a few hundred bare performers max in any library).

---

## 7. Verification before bulk run

Always do a `--dry-run` over the first ~20 performers and skim the proposed diffs. Things to spot-
check:
- Measurements formatting matches what you've manually typed elsewhere in stash
- Tattoos/piercings string formatting renders reasonably in the UI
- Image URLs from StashDB are reachable from the NAS (some are CDN'd)
- Gender enum maps cleanly

---

## 8. Open questions / risks

| Question | Risk | Mitigation |
|---|---|---|
| `measurements` string format conventions vary | Cosmetic only — values still searchable | Verify against existing entries; offer a `--measurements-format` flag if needed |
| `breast_type` mapping to stash's `fake_tits` string | Misleading display ("No" vs "Natural") | Confirm against stash UI text; allow override via config |
| Multi-match in `search` mode could pick wrong person | Wrong metadata applied to wrong performer | Skip by default; offer `flag_multiple_as_tag` like scene identifier |
| Image URL changes / 404s | Performer ends up without image | Pre-HEAD the URL before passing to stash; fall back to skipping image |
| Stash-box rate limits | Throttle/ban on rapid enrichment of large libraries | Inherit `max_requests_per_minute` from stash config (same as scene identifier) |
| Stash-box performer merges (`merged_into_id`) | We update against a stale/merged performer record | If `merged_into_id` is set on a fetch, follow the redirect to the current id |

---

## 9. Build order relative to other work

Suggested sequence:
1. Finish the **EXTERNAL-WORKERS.md** revisions (blockers from the design review).
2. Build **Phase A** of the GPU worker (previews) — highest impact for NAS speed.
3. Then **Phase 1** of this doc (`refresh` mode) — independent of GPU work, mostly clerical.
4. Phase B of workers (sprites), then Phase 2/3 of this doc, in either order.

This sequencing front-loads the biggest pain (slow NAS generation) and treats performer enrichment
as a polish pass once the bulk is identified. None of these block each other — they could run in
any order, or in parallel by different people.
