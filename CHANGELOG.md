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
- **Phase 1 — library assistant:** `internal/llm/` (Anthropic client + agent loop + tool registry),
  `/llm/chat` SSE route, `ui/v2.5/src/components/Assistant/` chat panel, read tools then write tools
  behind an `ask | auto | readonly` policy with confirmation cards.
- **Phase 2 — agentic dev loop:** flag-gated, path-jailed, diff-for-review.

[Unreleased]: https://github.com/Ryokushen/stash/compare/v0.31.1...llm-interface
