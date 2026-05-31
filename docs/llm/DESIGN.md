# stash-llm — Design

Embedding an LLM-powered assistant inside Stash (model-agnostic, via a gateway). This document is the source of truth for the
architecture and the phased plan. It references the real upstream code paths so implementation can
follow it directly.

---

## 1. Goals & non-goals

**Goals**
- A **library assistant** (Phase 1): chat with an LLM to search, summarize, and curate the media
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

As built (mapped from the v0.31.1 source). We deliberately kept the surface minimal — a REST + SSE
route, **no GraphQL schema changes** — so upstream merges stay clean:

| Concern | Upstream location | How we hooked in |
|---|---|---|
| HTTP routing (chi) | `internal/api/server.go` (`r.Mount("/scene", …)` block, ~L214) | `r.Mount("/llm", server.getLLMRoutes())` next to the others ✓ |
| Assistant endpoints | new `internal/api/routes_llm.go` | `GET /llm/status`, `POST /llm/chat` (SSE), `POST /llm/confirm`. Behind the global `authenticateHandler` middleware. **No GraphQL schema/resolver was added.** |
| Data access | `models.Repository` + the `*Service` types in `internal/manager` | LLM **tools call these directly in-process** via `txn.WithReadTxn`/`txn.WithTxn` — no HTTP round-trip back to `/graphql` |
| Config / secrets | `internal/manager/config/config.go` + the env allowlist in `init.go` | `assistant_base_url` / `assistant_api_key` / `assistant_model` (+ `enabled`/`write_policy`/`dev_loop_enabled`) keys & getters. **Each key must also be added to the `envBinds` map in `init.go`** or its `STASH_*` env var is silently dropped. |
| UI | `ui/v2.5/src/components/` + `App.tsx` | A floating `components/Assistant/AssistantWidget.tsx` mounted once in `App.tsx` (gated off setup views). No navbar entry or Settings section — config is via env/`config.yml`. |

New Go code: **`internal/llm/`** (OpenAI-compatible client, tool registry, agent loop). It depends on
`models.Repository` and the repository's reader/writers — injected from `server.go` (`getLLMRoutes`).

---

## 3. Architecture

```
┌──────────────────────────── browser (React, ui/v2.5) ────────────────────────────┐
│  components/Assistant/AssistantWidget.tsx — floating chat widget mounted in App.tsx │
│   parses the SSE stream (text / tool_call / tool_result / confirm_required / done)  │
└───────────────────────────────────┬───────────────────────────────────────────────┘
                                     │  POST /llm/chat  (SSE event stream)
                                     ▼
┌──────────────────── Go backend (internal/api  +  internal/llm) ────────────────────┐
│  internal/api/routes_llm.go   — chi route (/status, /chat SSE, /confirm), auth      │
│  internal/llm/                                                                      │
│   • client.go        — OpenAI-compatible chat-completions client (function calling) │
│   • service.go       — agent loop + in-memory conversations + settings/status        │
│   • tools.go         — tool registry; each tool = JSON schema + Go handler           │
│   • tools_library.go — Phase 1 handlers calling repository / *Service               │
└───────────┬─────────────────────────────────────────────────┬───────────────────────┘
            │ in-process calls (no self-HTTP)                   │ OpenAI /v1/chat/completions
            ▼                                                   ▼
  models.Repository · Scene/Performer/             ┌──────── LiteLLM gateway ────────┐
  Studio/Tag reader-writers                        │  minimax → MiniMax API (key)     │
            │                                       │  claude  → Claude-OAuth bridge   │
            ▼                                       └──────────────────────────────────┘
   SQLite · ffmpeg · generated content
```

### Transport choice — REST + SSE for the chat, in-process for tools
- **Chat stream:** `POST /llm/chat` returning **Server-Sent Events**. SSE is the simplest fit for
  emitting the assistant's text plus structured events (`text`, `tool_call`, `tool_result`,
  `confirm_required`, `error`, `done`) per loop turn. (A GraphQL subscription over Stash's existing
  websocket transport was the alternative; SSE keeps the assistant self-contained with no schema churn.)
- **Tools:** the LLM service calls the repository reader/writers **directly in Go** — it does *not* call
  back into `/graphql` over HTTP. Same process, so calls are faster, transactional, and need no auth
  re-plumbing. No GraphQL schema or resolver was added (see §2).

### Request/response shape (`POST /llm/chat`)
```jsonc
// request
{ "conversationId": "uuid|null", "message": "tag every unorganized scene from studio X as 'review'" }

// SSE events (text/event-stream) — one event per line-pair, terminated by a blank line
event: text             data: {"text":"I found "}
event: tool_call        data: {"id":"call_1","name":"find_scenes","input":{"studio":"X","organized":false}}
event: tool_result      data: {"id":"call_1","name":"find_scenes","is_error":false,"summary":"{\"count\":12,...}"}
event: confirm_required data: {"id":"call_2","name":"add_tags_to_scenes","input":{"scene_ids":[...],"tag_ids":[7]}}
event: text             data: {"text":"That will tag 12 scenes — approve?"}
event: done             data: {"conversationId":"uuid","usage":{"prompt_tokens":1234,"completion_tokens":567,"total_tokens":1801}}
```
`confirm_required` appears only under the `ask` write policy; the UI then calls `POST /llm/confirm`
`{name, input}` to execute the approved write.

---

## 4. Tool surface (Phase 1)

Each tool is a JSON-schema definition (sent to the model as an OpenAI function) plus a Go handler that
wraps the data layer. **Read tools** are always on; **write tools** are gated by the write policy
(see §6). The **8 tools shipped in Phase 1** (`internal/llm/tools_library.go`):

| Tool | Wraps | Kind |
|---|---|---|
| `library_stats` | counts via `repo.{Scene,Performer,Studio,Tag}.Count` | read |
| `find_scenes` | scene query (filters: free-text, organized, tag names, studio name, performer names; paginated) | read |
| `find_performers` | performer query (free-text) | read |
| `find_studios` | studio query (free-text) | read |
| `find_tags` | tag query (free-text) | read |
| `create_tag` | `Tag.Create` | write |
| `add_tags_to_scenes` | `Scene.UpdatePartial` with `UpdateIDs{Mode: Add}` over many scene ids | write |
| `set_scenes_organized` | `Scene.UpdatePartial` setting the organized flag over many scene ids | write |

Tools return compact, model-friendly JSON (ids + a few display fields), never raw rows, to keep token
use down. `find_scenes` pagination defaults to 25 with a hard cap of 100.

Natural follow-ups (not yet built): `get_scene` detail, `find_duplicates` (phash), `remove_tags`,
`set_scene_studio`/`set_scene_performers`, `create_studio`, and `trigger_scan`/`trigger_generate` jobs.

---

## 5. Providers, config & secrets

The assistant does **not** call providers directly. It speaks one OpenAI-compatible chat-completions
API to a **gateway** (LiteLLM) that fronts the real providers and owns their auth + OAuth refresh. This
keeps four very different auth schemes (MiniMax API key, Claude/Codex/xAI OAuth) out of stash entirely,
and lets the model set change without touching stash code.

- **Gateway:** `docker/llm/litellm/config.yaml` exposes named models to stash:
  - `minimax` → `minimax/MiniMax-M2.7` via `MINIMAX_API_KEY` (OpenAI-compatible; headless default).
  - `claude` → routed to an external **Claude-OAuth bridge** (`CLAUDE_BRIDGE_URL`) that holds and
    refreshes the setup-token. LiteLLM's own raw-OAuth passthrough is unreliable
    ([BerriAI/litellm#19618](https://github.com/BerriAI/litellm/issues/19618)) — the bridge is what
    makes Claude-OAuth work headless.
- **Config keys** in `internal/manager/config/config.go` (mirroring `GetHandyKey()` / `GetAPIKey()`):
  - `assistant_base_url` (string) — gateway OpenAI base, e.g. `http://litellm:4000/v1`. Env `STASH_ASSISTANT_BASE_URL`.
  - `assistant_api_key` (string) — the gateway (LiteLLM) master key. Env `STASH_ASSISTANT_API_KEY`.
  - `assistant_model` (string, default `minimax`) — gateway model name. Env `STASH_ASSISTANT_MODEL`.
  - `assistant_enabled` (bool, default true)
  - `assistant_write_policy` (enum `ask` | `auto` | `readonly`, default `ask`)
  - `assistant_dev_loop_enabled` (bool, default **false** — Phase 2 gate)
- **Secrets stay in the gateway.** The stash container holds only the gateway URL/key/model; provider
  keys live in `litellm.env` (git-ignored; `.env.example` files hold names only). On the NAS both env
  files are injected via compose `env_file`, and litellm is not published to the host — only stash
  reaches it over the private Docker network.

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
- **Phase 1 — library assistant** ✅ — shipped:
  1. ✅ `internal/llm/` client + loop + tool registry (8 read/write tools).
  2. ✅ `routes_llm.go` SSE endpoint + `/status` + `/confirm`, behind global auth; config keys.
  3. ✅ `components/Assistant/` floating chat widget mounted in `App.tsx`.
  4. ✅ Write tools behind the `ask` policy + confirm round-trip via `/llm/confirm`.
  5. ✅ Unit tests; full Docker build + boot smoke. Live chat smoke (real key + library) is the
     operator's step. Follow-ups: token streaming, a Settings GraphQL UI for the key, persistent
     conversations (see CHANGELOG "Planned").
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
