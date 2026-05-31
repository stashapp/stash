# Changelog

All notable changes to the **stash-llm** fork are documented here. This tracks *fork-specific* changes
only — the embedded Claude assistant and related infrastructure. For changes to upstream Stash itself,
see [stashapp/stash releases](https://github.com/stashapp/stash/releases) and the in-app changelog under
`ui/v2.5/src/docs/en/Changelog/`.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this fork's
versioning is independent of upstream — it is anchored to the upstream base release it sits on
(currently **stash v0.31.1**).

## [Unreleased]

### Added
- **Phase 1 — library assistant** (see `docs/llm/DESIGN.md`):
  - `internal/llm/`: a dependency-free **OpenAI-compatible** chat-completions client (function
    calling), a tool registry, the agent loop with an in-memory conversation store, and the Phase 1
    library tools — `library_stats`, `find_scenes` / `find_performers` / `find_studios` / `find_tags`
    (read) and `create_tag`, `add_tags_to_scenes`, `set_scenes_organized` (write). Tools call stash's
    repository in-process via `txn.WithReadTxn` / `txn.WithTxn`.
  - **Multi-provider via a gateway:** stash talks to one OpenAI-compatible endpoint (a LiteLLM
    gateway) that fronts the real providers and owns their auth/OAuth refresh — `docker/llm/litellm/`
    config exposes `minimax` (MiniMax-M2.7, API key) and `claude` (via an external OAuth bridge).
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
  - `docs/llm/DEPLOY-NAS.md` — NAS Docker deployment guide (reuse Jellyfin media folders read-only, app state under `/volume1/docker/stash`, ghcr vs save-load transfer, verification).
  - `docker/llm/docker-compose.nas.yml`, `docker/llm/build-and-push.sh`, `docker/llm/.env.example`.
- `.gitignore` rules blocking secret env files (`stash.env`, `*.env`) and saved image tars (`stash-llm-*.tar`).
- `.gitattributes` rule forcing `*.sh` to LF line endings so the build script's shebang works in bash.

### Changed
- Fork is based on the stable upstream tag **v0.31.1** rather than the `develop` tip, which currently
  fails to build via the official Dockerfile (Makefile regression in commit `92349790`, 2026-05-05, that
  makes the node-only frontend stage invoke Go). Rationale captured in `docs/llm/DESIGN.md`.

### Notes
- **Phase 0 (baseline) validated** on 2026-05-31: `docker/build/x86_64/Dockerfile` builds UI + Go from
  source; the container boots on `:9999`, returns HTTP 200, and reports `stash v0.31.1`.

### Planned
- **Phase 1 polish:** token-level streaming (currently streams per loop turn), a Settings UI to set
  the API key/model via GraphQL (currently env/config.yml + the `/llm/status` notice), and persisting
  conversations across restarts (currently in-memory).
- **Phase 2 — agentic dev loop:** flag-gated, path-jailed, diff-for-review.

[Unreleased]: https://github.com/Ryokushen/stash/compare/v0.31.1...llm-interface
