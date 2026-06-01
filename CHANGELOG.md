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
- **`tools/identify/external_identify.py`** — standalone, dependency-free (stdlib) fingerprint
  identifier that runs **outside stash's job queue** (so it works in parallel with a long
  Generate→Phash, or on another machine). Reads scene oshashes from the stash API, queries
  StashDB/TPDB via `findScenesBySceneFingerprints`, and applies MERGE-like updates via `sceneUpdate`
  (find-or-create studio/performers by name + stash_id stamp). Defaults to `--dry-run`. Lower fidelity
  than native Identify; documented in `tools/identify/README.md`.
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
