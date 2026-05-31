# stash-llm — Design

Embedding a Claude-powered assistant inside Stash. This document is the source of truth for the
architecture and the phased plan. It references the real upstream code paths so implementation can
follow it directly.

---

## 1. Goals & non-goals

**Goals**
- A **library assistant** (Phase 1): chat with Claude to search, summarize, and curate the media
  library in natural language — find scenes, bulk-tag, clean up metadata, answer "what do I have"
  questions — by calling Stash's existing data layer through a defined tool surface.
- An **agentic dev loop** (Phase 2): an in-app agent that can read/write this fork's own source,
  scaffold features/plugins, run a build, and surface a diff for review — sandboxed and gated.
- Stay **mergeable with upstream**: all new code lives under namespaced paths; we touch upstream files
  in as few, as small, and as obvious places as possible (route registration, config keys, nav entry).

**Non-goals**
- Not reimplementing Stash. We build on its media pipeline, DB, and GraphQL layer.
- Not exposing the assistant unauthenticated. It rides Stash's existing auth/session.
- Phase 1 does **not** modify source or run shell commands. That capability is Phase 2 only, off by
  default, behind an explicit feature flag.

---

## 2. Where this plugs into Stash

Mapped from the v0.31.1 source:

| Concern | Upstream location | How we hook in |
|---|---|---|
| HTTP routing (chi) | `internal/api/server.go` (`r.Mount("/scene", …)` block, ~L214) | Add `r.Mount("/llm", server.getLLMRoutes())` next to the others |
| GraphQL schema | `graphql/schema/schema.graphql` + `graphql/schema/types/*.graphql` | Add `graphql/schema/types/assistant.graphql`; `make generate-backend` regenerates |
| GraphQL resolvers | `internal/api/resolver.go` (`Resolver` struct holds `repository`, `sceneService`, `imageService`, `galleryService`, `groupService`) | New `resolver_*_assistant.go`; reuse the same services/repository the resolver already has |
| Data access | `models.Repository` + the `*Service` types in `internal/manager` | LLM **tools call these directly in-process** — no HTTP round-trip back to `/graphql` |
| Config / secrets | `internal/manager/config/config.go` (keys like `ApiKey = "api_key"`, getters like `GetAPIKey()`, viper-backed with env override) | Add `AnthropicAPIKey = "anthropic_api_key"` + `GetAnthropicAPIKey()`, env override `STASH_ANTHROPIC_API_KEY` |
| UI components | `ui/v2.5/src/components/` (domain folders), `MainNavbar.tsx`, `Settings/` | Add `components/Assistant/`, a navbar entry, and a Settings section for the API key/model |

New Go code: **`internal/llm/`** (provider client, tool registry, chat session/loop). It depends on
`models.Repository` and the manager services — injected the same way the resolver gets them.

---

## 3. Architecture

```
┌──────────────────────────── browser (React, ui/v2.5) ────────────────────────────┐
│  components/Assistant/                                                            │
│   • ChatPanel.tsx  — message list, composer, streaming tokens                     │
│   • useAssistant.ts — opens SSE to POST /llm/chat, renders tool-call cards         │
│  components/Settings/  → Anthropic key + model + feature flags                     │
└───────────────────────────────────┬───────────────────────────────────────────────┘
                                     │  POST /llm/chat  (SSE stream of tokens + events)
                                     ▼
┌──────────────────── Go backend (internal/api  +  internal/llm) ────────────────────┐
│  internal/api/routes_llm.go        — chi route, auth, request/SSE plumbing          │
│  internal/llm/                                                                      │
│   • client.go      — Anthropic Messages API client (streaming, tool use)            │
│   • session.go     — the agent loop: send → tool_use → run tool → tool_result → …   │
│   • tools.go       — tool registry; each tool = JSON schema + Go handler            │
│   • tools_library.go — Phase 1 handlers calling repository / *Service               │
│   • config.go      — reads key/model/flags from manager/config                      │
└───────────────────────────────────┬───────────────────────────────────────────────┘
                                     │  in-process calls (no self-HTTP)
                                     ▼
        models.Repository · SceneService · ImageService · GalleryService · GroupService
                                     │
                                     ▼
                          SQLite · ffmpeg · generated content
```

### Transport choice — REST + SSE for the chat, in-process for tools
- **Chat stream:** `POST /llm/chat` returning **Server-Sent Events**. SSE is the simplest fit for
  streaming Claude tokens + structured events (`tool_call`, `tool_result`, `error`, `done`). A
  GraphQL subscription (Stash already runs a websocket transport, `server.go` ~L176) is the
  alternative; we pick SSE to keep the assistant self-contained and avoid schema churn on the hot path.
- **Tools:** the LLM service calls `repository`/`*Service` **directly in Go**. We do *not* have Claude
  call back into `/graphql` over HTTP — same process, so direct calls are faster, transactional, and
  avoid auth re-plumbing. (The GraphQL schema addition in §2 is only for *non-streaming* helpers like
  listing available tools / conversation history, not the token stream.)

### Request/response shape (`POST /llm/chat`)
```jsonc
// request
{ "conversationId": "uuid|null", "message": "tag every untagged scene from studio X as 'review'" }

// SSE events (text/event-stream)
event: token        data: {"text":"I found "}
event: tool_call    data: {"id":"t1","name":"find_scenes","input":{"studio":"X","untagged":true}}
event: tool_result  data: {"id":"t1","summary":"12 scenes"}
event: token        data: {"text":"…tagging them now."}
event: done         data: {"conversationId":"uuid","usage":{"input":1234,"output":567}}
```

---

## 4. Tool surface (Phase 1)

Each tool is a JSON-schema definition (sent to Claude as a `tool`) plus a Go handler that wraps the
existing data layer. **Read tools** are always on; **write tools** are gated by a confirmation policy
(see §6). Initial set:

| Tool | Wraps | Kind |
|---|---|---|
| `find_scenes` | scene query (filters: text, studio, performer, tags, rating, organized, date, duration) | read |
| `find_performers` / `find_studios` / `find_tags` | corresponding queries | read |
| `get_scene` | scene by id incl. files, tags, performers, markers | read |
| `library_stats` | the stats resolver (`internal/api`, `Stats.tsx` backs this in UI) | read |
| `find_duplicates` | phash duplicate query (Stash already computes phash) | read |
| `add_tags_to_scenes` / `remove_tags_from_scenes` | bulk scene update via `SceneService` | write |
| `set_scene_studio` / `set_scene_performers` | scene update | write |
| `set_organized` | bulk organized flag | write |
| `create_tag` / `create_studio` | tag/studio create | write |
| `trigger_scan` / `trigger_generate` | the metadata task mutations (`resolver_mutation_metadata.go`) | write (job) |

Tools return compact, model-friendly JSON (ids + a few display fields), never raw rows, to keep token
use down. Pagination is enforced (default 25, hard cap) and surfaced to the model.

---

## 5. Config & secrets

- New config keys in `internal/manager/config/config.go`, mirroring the existing `GetHandyKey()` /
  `GetAPIKey()` pattern:
  - `anthropic_api_key` (string) — env override `STASH_ANTHROPIC_API_KEY`
  - `assistant_model` (string, default `claude-opus-4-8` / a configured default)
  - `assistant_enabled` (bool, default true)
  - `assistant_dev_loop_enabled` (bool, default **false** — Phase 2 gate)
  - `assistant_write_policy` (enum: `ask` | `auto` | `readonly`, default `ask`)
- Key is settable in the **Settings UI** (write-only field; never echoed back) or via env for the NAS
  container. **Never committed** — `docker/llm/.env.example` holds key *names* only, matching the
  fleet's redaction discipline.
- On the NAS the key is injected via the container environment (compose `env_file`), consistent with
  how other fleet containers handle secrets.

---

## 6. Safety model

**Phase 1 (library writes)**
- `assistant_write_policy=ask` (default): write tools don't execute immediately. The model proposes the
  action; the UI renders a confirmation card with the exact change + affected count; the user approves.
- `auto`: writes execute without confirmation (power users).
- `readonly`: write tools are not registered at all.
- All write tools log to Stash's logger with a clear `[assistant]` prefix and record an audit line
  (who/what/when/affected ids) so changes are traceable.
- Pagination/row caps prevent a single prompt from mutating the whole library by accident.

**Phase 2 (dev loop) — designed, not built yet**
- Off by default (`assistant_dev_loop_enabled=false`).
- Filesystem + build tools are **scoped to the repo working tree only** (path-jailed), run in a
  dedicated worktree/branch, and produce a **diff + build result for human review** — they never push
  or restart services autonomously.
- Intended to run only in a dev/staging container, never the production NAS deployment, until proven.
- Reuses the same provider client and loop; adds `read_file`, `write_file`, `run_build`, `run_tests`
  tools behind the flag, each audited.

---

## 7. Phased plan

- **Phase 0 — baseline** ✅ — fork builds & runs (v0.31.1); image `stash-llm:dev`; validated on :9999.
- **Phase 1 — library assistant**
  1. `internal/llm/` client + loop + tool registry (read tools first).
  2. `routes_llm.go` SSE endpoint + auth + config keys.
  3. `components/Assistant/` chat panel + Settings section.
  4. Write tools behind the `ask` policy + confirmation cards.
  5. Regression tests under `internal/llm/` and `internal/api/`; live smoke against a seeded library.
- **Phase 2 — agentic dev loop** — flag-gated, path-jailed, diff-for-review (per §6).

---

## 8. Build & deploy

`docker/build/x86_64/Dockerfile` builds UI + Go from source in one multi-stage build — our added code
compiles in with no extra steps. See [`DEPLOY-NAS.md`](DEPLOY-NAS.md) and `docker/llm/`. Local dev can
also run `make server-start` (backend) + `make ui-start` (Vite dev server) once a Go/Node toolchain is
present, but the Docker path is the supported one here.

---

## Appendix: why v0.31.1, not develop

The `develop` tip (as of 2026-05-31, commit `5646b793`) fails to build via the official Dockerfile.
Commit `92349790` ("[ci] add explicit flow for makefile", 2026-05-05) changed the Makefile so
`ui-only` depends on `generate` → `generate-backend` → `go generate`. The Dockerfile's **frontend**
stage (`node:24-alpine`, no Go) calls `make ui-only`, so it dies with `make: go: No such file or
directory`. **v0.31.1** (2026-04-13) predates this; its `ui-only` target is just `ui-env` +
`npm run build`, and it builds cleanly with the identical Dockerfile. We base on v0.31.1 and track
upstream via the `upstream` remote; we'll rebase forward once upstream fixes the regression.
