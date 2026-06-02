# Changelog

All notable changes to the **stash-llm** fork are documented here. This tracks *fork-specific* changes
only — the embedded LLM assistant and related infrastructure. For changes to upstream Stash itself,
see [stashapp/stash releases](https://github.com/stashapp/stash/releases) and the in-app changelog under
`ui/v2.5/src/docs/en/Changelog/`.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this fork's
versioning is independent of upstream — it is anchored to the upstream base release it sits on
(currently **stash v0.31.1**).

## [Unreleased]

### Added
- **`docs/llm/RUNBOOK.md`** — single operations runbook consolidating every external/off-queue tool's
  copy-paste command (worker tasks, the identify/parse/enrichment scripts, maintenance tools, db_backup,
  facerec), with the real NAS/`--media-prefix`/`--ffmpeg` values, dry-run-first defaults, the pipeline run
  order, the don't-run-together rules, and the NAS-disk concurrency cap. Linked from the README. Audited
  every tool's flags against the code first (no drift found). Replaces hopping between five per-tool READMEs.
- **Seven more speed/convenience tools** (built in parallel; all default to dry-run/report):
  - **`worker` image-phash task** (`worker/internal/image.go`) — `--tasks image-phash` enumerates
    `findImages`, decodes each ImageFile, computes `goimagehash.PerceptionHash`, writes via
    `fileSetFingerprints`. Images don't support `is_missing:"phash"` server-side, so it enumerates all
    and skips phashed files client-side (stable pagination, like previews). Unlocks image dedup.
  - **`tools/maintenance/scene_dedup.py`** — clusters near-duplicate scenes by phash Hamming distance
    (`--max-distance`, default 8), reports keeper vs dup-candidates + reclaimable space, `--apply` tags
    non-keepers `_dupe-candidate`.
  - **`tools/maintenance/db_backup.py`** — consistent SQLite **online-backup** of the metadata DB +
    `config.yml` → gzip tar, retention (`--keep`), optional off-NAS `--remote` rsync. Meant as a NAS cron.
  - **`tools/maintenance/ingest_pipeline.py`** — orchestrates the off-queue chain (identify → filename-parse,
    in that safe order) over a growing library, with before/after gap snapshots. Single cron entry point.
  - **`tag_consolidate.py --emit-plan`** — emits a structured JSON merge plan of the subjective tag
    families (POV/hair/age/1-edit) for the in-app assistant or a human to execute via `tagsMerge`.
  - **`tools/facerec/`** (InsightFace) — performer face tagger for the scenes that resist both fingerprint
    and filename matching: build per-performer embeddings from stash images, match scene frames, link
    **existing** performers only (never creates), dry-run default. Runs on the RTX box (GPU deps + model
    download required — see `tools/facerec/README.md`).
  - **`external_identify_performers.py` documented** (`tools/identify/README.md`) — modes, invocations,
    NAS cron for periodic performer metadata enrichment.
- **Assistant `run_command` (shell escape hatch) + action-first prompt** (`internal/llm/tools_exec.go`).
  The assistant can now run shell commands inside the Stash container — the escape hatch for anything
  GraphQL/`define_tool` can't express: running the bundled helper scripts with custom flags (e.g.
  `external_identify.py --phashed-only`), editing config, filesystem inspection, package installs. It's
  `Writes:true` so every command is **confirm-gated** (shown to the user before it runs), and its
  registration is gated behind **`assistant_dev_loop_enabled`** (`STASH_ASSISTANT_DEV_LOOP_ENABLED=true`)
  so the whole capability has an explicit on/off switch. The system prompt now **defaults to action** —
  when asked to change a built-in's behavior it wraps it via `define_tool` or runs the underlying helper
  via `run_command` instead of refusing (it cannot recompile the Go binary, the one true limit).
  Also fixed `/llm/confirm` to return a well-formed body for plain-text tool output (it had wrapped
  output as `json.RawMessage`, breaking on non-JSON results like `run_command`'s).
- **Library maintenance + transcode tooling** (four off-queue/external helpers):
  - **`worker` transcode task** — new `--tasks transcode` pre-generates a
    browser-friendly h264/aac MP4 into `generated/transcodes/<hash>.mp4`, which stash
    serves directly (`GetStreamPath`, `pkg/models/paths/paths_scenes.go:28-35`) so the
    weak NAS never live-transcodes. Smart: if the video is already h264 it **remuxes**
    (`-c:v copy`, lossless, just fixes audio/container); otherwise a full **NVENC h264**
    re-encode at full resolution (CPU decode for universal codec support). Targets
    everything stash can't stream directly **plus HEVC** (which stash serves direct but
    browsers often can't play). `worker/internal/transcode.go` + `NeedsTranscode` +
    `GeneratedPaths.Transcode`; scene projection extended with
    `video_codec/audio_codec/format/width/height`. ~25-30% of the library qualifies.
  - **`tools/maintenance/tag_consolidate.py`** — clusters duplicate tags; auto-merges
    the SAFE set (case/punctuation/plural/`blond`↔`blonde`) via `tagsMerge`, preserving
    merged names as aliases, and reports subjective families (POV/hair/age/1-edit) for
    review. Dry-run by default.
  - **`tools/maintenance/generated_sweeper.py`** — audits `generated/` for orphaned,
    zero-byte, and stale-tmp artifacts against live file hashes; `--prune` quarantines
    orphans (reversible). Doubles as a coverage report. Runs on the NAS.
  - **`tools/identify/external_filename_parse.py`** — links unidentified scenes to
    **existing** studios/performers parsed from their path (multi-word name/alias
    matching only — never creates records, so it can't introduce junk). Fills empty
    studio/date (+ optional title), merges performers. Dry-run by default.
- **Assistant autonomy — generic GraphQL access + self-defined tools** (`internal/llm/tools_graphql.go`,
  `tools_dynamic.go`). The assistant is no longer limited to hand-coded Go tools; new capabilities no
  longer require a rebuild/redeploy:
  - **In-process GraphQL executor** wired in `internal/api/server.go`: runs arbitrary operations against
    stash's own schema through the full resolver chain **and the dataloader middleware** (an httptest
    round-trip through the dataloader-wrapped gql handler). Passed into `llm.Deps` alongside the parsed
    `*ast.Schema` and a persistent `ToolsDir`.
  - `graphql_schema` (introspection, scoped by `type`/`section`/`search`), `graphql_query` (read-only —
    rejects mutation/subscription docs), and `graphql_mutate` (`Writes:true` → confirm-gated).
    Operations are validated against the live schema before execution, so the model gets precise
    correction errors. "Combine similar tags" is now just `tagsMerge` via `graphql_mutate` — no new Go code.
  - **Persisted dynamic tools:** `define_tool` lets the assistant create a named, reusable tool (a GraphQL
    op + a JSON-Schema for its args; the call's args pass through as GraphQL variables). Definitions are
    saved as JSON under `<config>/llm_tools/` (a mounted volume, **not** the image), hot-registered for
    immediate use, and reloaded on boot — so self-authored capabilities survive restarts with no rebuild.
    `list_dynamic_tools` / `delete_dynamic_tool` manage them. A dynamic tool whose GraphQL mutates is
    auto-marked write-gated.
  - **Safety:** every library mutation (`graphql_mutate` + mutating dynamic tools) flows through the
    existing write-policy confirm gate (`write_policy=ask`); reads stay un-gated; built-in tool names
    can't be shadowed. The system prompt now teaches introspect → read → mutate and self-extension.
- **External GPU worker** (`worker/`) — separate Go module producing a single Windows `.exe` that
  offloads stash's heaviest media-generation work to a Windows + NVIDIA box, with outputs landing
  directly in stash's `generated/` share over SMB. Runs **outside stash's job queue** (parallel
  with whatever stash is doing). **Live-validated 2026-06-01 on RTX 5080** against the production NAS.
  Tasks shipped:
  - **previews** — 12-segment NVENC encode + concat demuxer matching stash's `transcoder.Splice`
    output. CRF→`-cq 21`, scale=640:-2, yuv420p high@4.2, 128k AAC, faststart. Two-pass slow-seek
    retry on first-pass failure (matches `task_generate_preview.go:67-72`). VFR `-vsync 2` branch
    when ffprobe reports ≤0.01 fps. ~7s/scene; 98% success rate observed on a 50-scene batch.
  - **covers** — single-frame extract at 20% of duration (mirrors `screenshotDurationProportion`
    in `pkg/scene/generate/screenshot.go:17`). Base64 JPEG → `sceneUpdate(cover_image:…)` mutation
    (no fork additions needed — `SceneUpdateInput.cover_image` already exists in upstream). ~3s/scene.
  - **sprites** — per-frame `-ss` seek with NVDEC, tiled in Go via `image/draw` instead of ffmpeg's
    `tile` filter. WebVTT cells derived from the actual decoded frame dimensions. Initial
    implementation used the `fps` filter (full-pass decode) at 73s/scene; later optimized to
    parallel-bounded seek + Go tiling at ~30s/scene (commit `45206dc8`).
  - **phash** (Phase C) — a **bit-for-bit** replica of stash's `pkg/hash/videophash/phash.go`:
    25 lossless BMP frames (5×5) at `offset=0.05·dur, step=0.9·dur/25`, scaled `scale=160:-2`,
    montaged with `disintegration/imaging` v1.6.2, hashed with `corona10/goimagehash` v1.1.0
    (the exact libs + versions stash uses). Frames are **CPU-decoded** (no NVDEC — hardware decode
    yields different pixels and breaks the hash). Uses the **stash-stored duration** from the API,
    not a fresh ffprobe, so sampled timestamps match stash's. Writes per-FILE via the existing
    `fileSetFingerprints` mutation (value = `FormatUint(hash, 16)`; stash round-trips it through
    `ParseUint(_,16,64)`). No fork additions. `is_missing:"phash"` server-side filter (shrinks as
    phashes are applied → same re-fetch-page-1 pagination as covers). ~15-25s/scene.
  - **`--verify-phash N` gate** — recomputes `N` files that already have a native stash phash and
    compares; exits non-zero on any mismatch. Bit-exactness depends on the ffmpeg build's
    decoder/scaler, so this proves the worker matches the NAS **before** writing hashes in bulk.
    First gate run 2026-06-01: 25/25 bit-identical against native phashes.
  - CLI: `--tasks previews,covers,sprites,phash` dispatcher; `--limit N` caps **encoded items**
    (skips don't count); `--max-failures` safety brake; `--watch` keep-running mode; `--dry-run`
    prints commands without writing. `--generated-prefix` is now required only for previews/sprites
    (covers and phash write via the API, not the filesystem).
  - Worker-side path translation (`internal/paths.go`): `normalizeUNC` recovers from
    bash/PowerShell/MSYS collapsing `\\server\share` to `\server\share`, so SMB UNC paths "just work"
    regardless of shell quoting layer.
  - SMB atomicity: `.partial` files live in `<generated>/tmp/` (stash's own scratch dir), atomic
    rename to destination so stash never sees half-written outputs.
  - `--tasks` filter strategy per task: `is_missing:"cover"` filter for covers (stash tracks
    missingness in SQLite); filesystem walk for previews and sprites (stash doesn't track those —
    `pkg/sqlite/scene_filter.go:396-436`).
  - Two concurrent `.exe` instances are supported and useful — run a `--tasks covers` worker in
    parallel with a `--tasks previews,covers` worker so the cover sweep finishes much sooner than
    waiting for previews to complete. NAS disk is the real bottleneck at that point.
- **`phashed_only` filter** on `identify_scenes_fast` / `external_identify.py` — when set, only scenes
  with a phash are considered. Lets you target the **high-yield** scenes first while stash is still
  generating phashes for the rest (work in parallel, no waiting on the queue). The script now also
  prints a diagnostic line (`N candidate scene(s) [X with phash+oshash, Y with oshash only]`) and the
  assistant tool parses this into the result so the model can report fingerprint coverage of each run.
- **`identify_scenes_fast` assistant tool** — bundles `external_identify.py` into the runtime image
  (`/usr/local/bin/identify_external.py`) and exposes it as an assistant tool that shells out via
  `os/exec`. Runs **outside stash's job queue** (parallel with Generate→Phash etc.) and uses **batched
  fingerprint lookups** (~40× fewer round-trips than native Identify). Default 200 unorganized
  scenes/run, `--set-organized`, multi-match skip; 100s timeout (under the LLM client's 120s) so a
  run fits in one chat turn. Parses the summary into compact JSON (matched/applied/skipped/no_match).
  Native `identify_scenes` remains for full-fidelity, no-urgency runs.
- **`tools/identify/external_identify.py`** — standalone, dependency-free (stdlib) fingerprint
  identifier that runs **outside stash's job queue** (so it works in parallel with a long
  Generate→Phash, or on another machine). Reads scene oshashes from the stash API, queries
  StashDB/TPDB via `findScenesBySceneFingerprints`, and applies MERGE-like updates via `sceneUpdate`
  (find-or-create studio/performers by name + stash_id stamp). Defaults to `--dry-run`. Lower fidelity
  than native Identify; documented in `tools/identify/README.md`.
  - **Tags now applied by default (MERGE).** The stash-box query now requests `tags { id name }`,
    and matched tags are added on top of the scene's existing tags (never replacing them), creating
    any tag that doesn't exist yet (stamped with its stash-box id). Idempotent — only writes `tag_ids`
    when there's something new to add. `--no-tags` opts out. This closes the previous "no tags" gap;
    the in-app `identify_scenes_fast` inherits it (it shells out to the same script).
- **Bulk identify / mass tagging** — the assistant can now drive stash-box identification:
  - `identify_scenes` (write): wraps `metadataIdentify` / `manager.CreateIdentifyJob` — runs the
    Identify task over all or selected scenes against the configured stash-box endpoints (StashDB,
    ThePornDB, …) in priority order, as a background job. Supports `set_organized`,
    `skip_multiple_matches`, and `flag_multiple_as_tag` (tags ambiguous skips for manual review).
  - `fingerprint_coverage` (read): reports phash coverage (with/missing) so you know whether to run
    Generate → Phash before identifying (mass matching is fingerprint-based).
- **Scraper management** — the assistant can now manage metadata scrapers itself:
  - `internal/llm/tools_scrapers.go`: `list_scrapers` (read) and `reload_scrapers` (write), wrapping
    `manager.GetInstance().ScraperCache.ListScrapers` and `RefreshScraperCache` in-process.
  - The runtime image now ships **Python 3** + the common CommunityScrapers deps (requests, lxml,
    beautifulsoup4, dateutil via apk; cloudscraper + stashapp-tools via pip) so python-backed scrapers
    actually run (`docker/build/x86_64/Dockerfile`). Previously the alpine image had no Python.
- **Phase 1 — library assistant** (see `docs/llm/DESIGN.md`):
  - `internal/llm/`: a dependency-free **OpenAI-compatible** chat-completions client (function
    calling), a tool registry, the agent loop with an in-memory conversation store, and the Phase 1
    library tools — `library_stats`, `find_scenes` / `find_performers` / `find_studios` / `find_tags`
    (read) and `create_tag`, `add_tags_to_scenes`, `set_scenes_organized` (write). Tools call stash's
    repository in-process via `txn.WithReadTxn` / `txn.WithTxn`.
  - **Multi-provider via a gateway:** stash talks to one OpenAI-compatible endpoint (a LiteLLM
    gateway) that fronts the real providers and owns their auth/OAuth refresh — `docker/llm/litellm/`
    config exposes `minimax` (MiniMax-M2.7 via API key) and `grok` (grok.com OAuth subscription via a
    `cli-proxy-api` bridge container, not metered API billing).
  - Write policy: `readonly` omits write tools, `auto` executes them, `ask` (default) emits a
    `confirm_required` event and defers to `POST /llm/confirm` after the user approves in the UI.
  - `internal/api/routes_llm.go`: `GET /llm/status`, `POST /llm/chat` (Server-Sent Events),
    `POST /llm/confirm` — behind stash's global auth middleware.
  - Config keys `assistant_base_url` / `assistant_api_key` / `assistant_model` +
    `assistant_{enabled,write_policy,dev_loop_enabled}`, env overrides (`STASH_ASSISTANT_BASE_URL`, …).
  - `ui/v2.5/src/components/Assistant/`: a floating chat widget (SSE-over-fetch) mounted on all
    authenticated views, with inline tool activity and write-confirmation prompts.
  - Unit tests for the registry, client request/response, API-error surfacing, and tool-result
    marshalling. Verified end-to-end with a full Docker build + boot (UI 200, all 8 tools listed).
- Fork scaffolding and documentation:
  - `README.md` fork banner (purpose, added code paths, quick build), preserving upstream's README below it.
  - `docs/llm/DESIGN.md` — architecture, phased plan, Phase 1 tool surface, config/secrets, and safety model, grounded in the real v0.31.1 integration points (`internal/api/server.go` routing, the `Resolver` services, `internal/manager/config` keys, UI layout).
  - `docs/llm/DEPLOY-NAS.md` — NAS Docker deployment guide (stash + litellm, reuse Jellyfin media folders read-only, app state under `/volume1/docker/stash`, ghcr vs save-load transfer, verification).
  - `docker/llm/`: `docker-compose.nas.yml` (stash + litellm + cli-proxy-api services), `build-and-push.sh`, `litellm/config.yaml` (gateway: minimax + grok), `cliproxy/config.example.yaml` (Grok-OAuth bridge), and `.env.example` (stash) + `litellm.env.example` (provider secrets).
- `.gitignore` rules blocking secret env files (`stash.env`, `*.env`) and saved image tars (`stash-llm-*.tar`).
- `.gitattributes` rule forcing `*.sh` to LF line endings so the build script's shebang works in bash.

### Changed
- Fork is based on the stable upstream tag **v0.31.1** rather than the `develop` tip, which currently
  fails to build via the official Dockerfile (Makefile regression in commit `92349790`, 2026-05-05, that
  makes the node-only frontend stage invoke Go). Rationale captured in `docs/llm/DESIGN.md`.

### Fixed
- Assistant tools no longer crash the chat stream when the library DB isn't ready (e.g. before stash
  Setup has been run): `runToolSafely` recovers any tool-handler panic and returns it as a graceful
  tool error the model can report. Found during the first NAS deploy — chatting before Setup hit a
  nil-pointer in `library_stats` and the SSE stream died mid-turn.

### Notes
- **First NAS deploy** 2026-05-31: stash + litellm + cli-proxy-api up on overwatch-stash via
  `docker compose`; `/llm/status` → `configured:true, model:minimax`. Live chat pending stash Setup
  (DB creation + library scan) and, for grok, the one-time OAuth login.
- **Phase 0 (baseline) validated** on 2026-05-31: `docker/build/x86_64/Dockerfile` builds UI + Go from
  source; the container boots on `:9999`, returns HTTP 200, and reports `stash v0.31.1`.

### Planned
- **Phase 1 polish:** token-level streaming (currently streams per loop turn), a Settings UI to set
  the API key/model via GraphQL (currently env/config.yml + the `/llm/status` notice), and persisting
  conversations across restarts (currently in-memory).
- **Phase 2 — agentic dev loop:** flag-gated, path-jailed, diff-for-review.

[Unreleased]: https://github.com/Ryokushen/stash/compare/v0.31.1...llm-interface
