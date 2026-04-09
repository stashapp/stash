# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Is

Stash is a self-hosted Go + React web application for organizing and serving media content (videos, images). It uses GraphQL as the API layer, SQLite for storage, and FFmpeg for media processing.

## Build & Development Commands

### Prerequisites
- Go 1.24.3, Node.js 18+, pnpm, FFmpeg, golangci-lint

### Initial Setup
```bash
make pre-ui      # Install UI dependencies (pnpm install)
make generate    # Generate GraphQL code (backend + frontend types)
```

### Running for Development (two terminals)
```bash
make server-start   # Backend on :9999, stores data in .local/
make ui-start       # Vite dev server on :3000, proxies API to :9999
```

### Building
```bash
make stash          # Backend binary only
make ui             # Frontend only
make build-release  # Optimized release binary
```

### Linting
```bash
make lint           # golangci-lint on backend
make validate-ui    # ESLint + TypeScript check + Prettier on UI
make validate       # Both
```

### Testing
```bash
make test           # Go unit tests: go test ./...
make it             # Unit + integration tests (requires SQLite): go test -tags "... integration" ./...

# Single test:
go test -v ./pkg/models/... -run TestScenePartial
go test -v -tags integration ./pkg/sqlite/... -run TestSceneQuery
```

### Code Generation
After modifying `.graphql` schema files or adding new resolvers:
```bash
go generate ./cmd/stash             # Regenerates internal/api/generated_*.go
cd ui/v2.5 && npm run gqlgen        # Regenerates UI TypeScript types
```

After modifying interfaces that have mocks:
```bash
go generate ./pkg/models/...        # Regenerates mocks via mockery
```

## Architecture Overview

### Request Flow
Browser → GraphQL (`:9999/graphql`) → `internal/api` resolvers → service layer (`pkg/scene`, `pkg/image`, `pkg/gallery`) → `pkg/sqlite` → SQLite DB

### Key Layers

**`internal/api/`** — GraphQL resolvers only. ~95 resolver files split by domain (scene, image, gallery, performer, studio, tag, etc). No business logic here; resolvers delegate to services or the manager. Generated code: `generated_exec.go`, `generated_models.go`.

**`internal/manager/`** — Central singleton (`Manager`) wiring together config, DB, services, tasks, and plugins. `task/` contains all async operations (scan, generate, autotag, import/export) run via the job queue in `pkg/job`.

**`pkg/models/`** — Domain model interfaces and types. Repositories (e.g., `SceneReader`, `SceneWriter`) are interfaces here; implementations live in `pkg/sqlite/`.

**`pkg/sqlite/`** — SQLite implementation of all repository interfaces. Uses `sqlx`. Database migrations are in `pkg/sqlite/migrations/`. Integration tests live here and require the `integration` build tag.

**`pkg/scene/`, `pkg/image/`, `pkg/gallery/`** — Service layer with business logic for each media type (generate, scan, identify, etc).

**`pkg/scraper/`** — Scraper system supporting both XPath/CSS scrapers (via config YAML) and CDP (Chrome DevTools Protocol) for JavaScript-heavy sites.

**`pkg/plugin/`** — Plugin system supporting JavaScript (via Goja) and external Python plugins with hook system.

**`ui/v2.5/src/`** — React 17 frontend using Apollo Client for GraphQL, Bootstrap + SCSS for styling. GraphQL operations defined in `ui/v2.5/graphql/` and auto-generated as TypeScript hooks.

### GraphQL Schema
Schema source files: `graphql/schema/*.graphql`  
gqlgen config: `gqlgen.yml`  
Resolver stubs: `internal/api/resolver*.go` — new schema fields need corresponding resolver methods added here.

### File Identity & Deduplication
Files are identified by fingerprints: MD5, oshash (OpenSubtitles hash), and phash (perceptual hash for visual similarity). These drive duplicate detection and external metadata matching.

### Task/Job System
Long-running operations (scan, generate, identify) are submitted to a job queue (`pkg/job`). Tasks in `internal/manager/task/` implement the `Task` interface with `Start(ctx)`. Progress is streamed via GraphQL subscriptions.

### DataLoaders
To avoid N+1 queries in GraphQL, `internal/api/dataloader/` contains batch loaders for related entities (e.g., loading all performers for a list of scenes in one query).

## Important Conventions

- **Build tags**: Integration tests use `//go:build integration`. Always check if a test file has this tag before running without `-tags integration`.
- **Model mutations**: `pkg/models` types use `UpdateInput` partial structs with `OptionalString`/`OptionalInt` etc. for nullable field updates — don't use raw pointers for optional update fields.
- **Resolver pattern**: Mutation resolvers follow `mutationResolver` struct in `internal/api/resolver_mutation_*.go`; query resolvers use `queryResolver` in `resolver_query_*.go`.
- **Database queries**: Filter types in `pkg/models/` map to SQL via `pkg/sqlite/` query builders — adding a new filter field requires changes in both places.
