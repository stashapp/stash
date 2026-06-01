# External performer identifier — design

Companion to `external_identify.py` (which identifies **scenes**). This tool enriches **performers**
with full metadata from StashDB / ThePornDB, running outside stash's job queue. Shipped as
`tools/identify/external_identify_performers.py` (commit `fac8b029`) — `refresh`/`search`/`both`
modes per §6. This doc remains the design spec it was built to.

Same architectural pattern as the scene identifier: stdlib Python, talks to stash GraphQL + the
configured stash-box endpoints, applies MERGE-like updates. The bundled fragment query at
`graphql/stash-box/query.graphql` (`PerformerFragment`) already returns every field we need.

> **Critical reference: `pkg/stashbox/performer.go`** — stash itself contains the canonical
> stash-box → stash translation logic. **The Python tool MUST mirror this file byte-for-byte** for
> measurements formatting, body-mod string formatting, enum mapping, null handling, alias dedup,
> and merged-performer redirects. Don't re-derive — port. Quoted line refs throughout this doc.

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
| **`search`** | Unlinked (no `stash_id`) | `searchPerformer(term: <name>)` — **fuzzy, returns up to ~30 results** | Yes — see below |

`refresh` mode is the high-yield, zero-risk first pass — every match is unambiguous because we have
the canonical id. `search` mode handles the long tail (older performers, ones imported without
stash-box links) at higher ambiguity risk.

**`search` mode ambiguity handling (important):** stash-box's `searchPerformer` is a fuzzy name
match — searching "Mia" returns many results, and the API does not document ordering as "best match
first." Stash itself filters with `strings.EqualFold(performer.Name, name)` post-search
(`pkg/stashbox/performer.go:311-325`). **The tool must do the same** — apply an exact case-
insensitive name filter before any apply step. `--allow-multiple` only kicks in for genuine
post-filter ambiguity.

**Filtering candidate performers from stash** (B9 from the review):
`PerformerFilterType` has `is_missing: String` but it only accepts **one value per call**
(`pkg/sqlite/performer_filter.go:331-365`). There is no compound "all these fields are empty"
filter. So the tool does **client-side filtering**: fetch all performers in pages (including all
fields we might fill), then filter in memory. At "few hundred bare performers" scale this is fine.

---

## 3. Architecture

```
stash findPerformers (paged, with files/aliases/stash_ids/all fields)
   ↓
client-side filter (refresh = stash_ids set & some field empty; search = no stash_ids)
   ↓
per-performer:  stash-box findPerformer(id) OR searchPerformer(term)→EqualFold(name) filter
   ↓
build SceneUpdate from PerformerFragment using pkg/stashbox/performer.go translation
   ↓
**merge with existing alias_list + stash_ids** (Set mode REPLACES wholesale; must union explicitly)
   ↓
stash performerUpdate
```

- One Python script (stdlib only) — proposed name `external_identify_performers.py`.
- Reuses `gql()` helper + `Stash`/`StashBox` classes + `fetch_stash_boxes` from `external_identify.py`.
- `--dry-run` default, `--apply` to write, `--mode refresh|search|both`.
- Sources in stash-box priority order (StashDB first, TPDB fallback) — same as scene identifier.

**Don't run alongside native stash performer scraping.** `performerUpdate` is full-record partial
replace under one txn; concurrent updates clobber each other (last write wins). Same caveat as
`external_identify.py` for scenes — document it in the script's `--help`.

---

## 4. Field mapping (stash-box `PerformerFragment` → stash `PerformerUpdateInput`)

**Port `pkg/stashbox/performer.go` directly.** That file already does this translation for stash's
internal scraper — every formatting decision and null-handling rule below comes from there:

| Stash-box field | Type | → | Stash field | Type | Translation |
|---|---|---|---|---|---|
| `name` | String | | `name` | String | direct (refresh mode: never overwrite existing name) |
| `disambiguation` | String | | `disambiguation` | String | direct |
| `aliases` | [String!] | | `alias_list` | [String!] | filter aliases case-equal to name (`performer.go:273-280`), case-fold dedupe with existing, then UNION with existing (Set mode replaces — see §5) |
| `gender` | GenderEnum | | `gender` | GenderEnum | direct — enum names match exactly (MALE/FEMALE/TRANSGENDER_MALE/TRANSGENDER_FEMALE/INTERSEX/NON_BINARY) |
| `urls[].url` | [URL] | | `urls` | [String!] | flatten to URL strings; merge with existing |
| `images[0].url` | [Image] | | `image` | String | URL passthrough (stash downloads it). **Pre-HEAD the URL** — skip image on non-200 to avoid full-update rollback on CDN hiccup. No `is_primary` field exists on Image; convention is index 0 (`performer.go:225-227`) |
| `birth_date` | String | | `birthdate` | String | direct (just renamed) — `YYYY-MM-DD` |
| `death_date` | String | | `death_date` | String | direct |
| `ethnicity` | enum/String | | `ethnicity` | String | stash-box returns an enum; map to lowercase string (`performer.go` converts via `String()`) |
| `country` | String | | `country` | String | direct |
| `eye_color` | enum | | `eye_color` | String | stash-box enum → lowercase string |
| `hair_color` | enum | | `hair_color` | String | stash-box enum → lowercase string |
| `height` | Int (cm) | | `height_cm` | Int | **skip when `0`** (`performer.go:229`) |
| `measurements{band_size,cup_size,waist,hip}` | struct | | `measurements` | String | format as `"<band><cup>-<waist>-<hip>"` — **only when ALL FOUR sub-fields are non-null** (`performer.go:128-135`). Skip otherwise. Practically: female-only data. |
| `career_start_year` | Int | | `career_start` | String | stringify (`"2018"`) |
| `career_end_year` | Int | | `career_end` | String | stringify |
| `tattoos[]{location,description}` | [BodyMod] | | `tattoos` | String | port `formatBodyModifications` from `performer.go` (exact format used by stash itself) |
| `piercings[]{location,description}` | [BodyMod] | | `piercings` | String | same source function |
| `breast_type` | enum | | `fake_tits` | String | port the exact mapping in `performer.go` — verify before guessing (the doc draft's `NATURAL→"No"` was unverified) |
| `id` (stash-box id) | ID | | `stash_ids` | [StashIDInput!] | UNION with existing (Set mode replaces — see §5) |

**Fields with no source — acknowledge the gap:**
- stash has `weight: Int`, `penis_length: Float`, `details: String`, `tag_ids: [ID!]`, `favorite`,
  `rating100`, `circumcised`, `fake_tits` (existing user values) → stash-box's PerformerFragment has
  no analog. The tool **never touches these** (MERGE rule: missing source = don't write).
- stash-box has `merged_ids`, `breast_type` (handled), `merged_into_id`, `deleted` → see merge
  redirect handling in §5.

**MERGE semantics**: only write fields where stash's current value is empty/null *and* stash-box's
value is non-empty. Never overwrite existing user edits. For multi-value fields (`urls`,
`alias_list`, `stash_ids`), always pass the union with the existing stash values — see §5.

---

## 5. Multi-value fields & merge semantics — the trap

The earlier draft of this doc claimed `alias_list` would "union with existing automatically." **That
is wrong.** Verified in `internal/api/resolver_mutation_performer.go:401`: even in `Set` mode the
resolver calls `updatedPerformer.Aliases.Apply(p.Aliases.List())` which **replaces wholesale**. The
"dupes auto-filtered server-side" schema comment is intra-list dedupe, not union-with-existing.

So for **`alias_list`, `stash_ids`, and `urls`**, the tool's apply step must:
1. Fetch the performer's current values (already done — the candidate enumeration query selects all
   fields needed).
2. Build the union: `existing ∪ new_from_stashbox`.
3. Pass the unioned list in `performerUpdate`. The server's case-fold dedupe trims dupes but won't
   restore values you didn't include.

For `stash_ids` specifically: keep all existing entries from any endpoint; add `{endpoint, stash_id}`
only if not already present. Never drop entries from other stash-box endpoints.

**Merge redirects** (`merged_into_id`, `deleted` on PerformerFragment): if a fetched performer has
`deleted == true` and `merged_into_id != nil`, refetch that target id **once** (single hop, no
recursion — the doc draft said "follow the redirect" without bounds; that's an infinite-loop risk if
stash-box data is corrupt). If the target is also deleted, log and skip the performer.

**Other gotchas:**
- **Images**: pre-HEAD the URL before passing it in `performerUpdate`. A 404/timeout makes stash
  fail the full mutation, rolling back every field. Better to skip just the image.
- **`weight`, `penis_length`, `details`**: in stash, not in stash-box. Don't touch.
- **`circumcised`, `favorite`, `rating100`, `tag_ids`**: same — never overwrite user values.
- **Birthdate / death_date format**: stash-box returns ISO `YYYY-MM-DD`; stash accepts that exactly.

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
| ~~`measurements` format~~ | — | **Closed.** Port `formatMeasurements` from `pkg/stashbox/performer.go` — matches whatever stash itself produces. No format flag needed. |
| ~~`breast_type` mapping~~ | — | **Closed.** Port the exact map from `performer.go`. Stop guessing. |
| Multi-match in `search` mode | Wrong metadata applied to wrong performer | EqualFold(name) post-filter is mandatory (§2). `--allow-multiple` only after that filter. Default skip. |
| ~~Image URL changes / 404s~~ | — | **Closed.** Pre-HEAD before update; skip image on non-200. Documented in §4 + §5. |
| Stash-box rate limits, especially with `--mode both` | Cumulative cost doubles | Inherit `max_requests_per_minute` from stash config; emit progress meter so big runs are observable. |
| Performer merges (`merged_into_id`) | We update against a stale record | 1-hop redirect only, no recursion. If target also deleted, skip with warning. |
| Race with native stash performer scrape | Last-write-wins under one txn | Document in `--help`; recommend disabling native auto-scrape during bulk runs. |
| `weight`/`penis_length`/`details` are stash-only | Could be enriched from elsewhere later | Acknowledged in §4 — never touched by this tool. Future scope. |

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
