# Stash - Organiser for adult media

Stash is a self-hosted web application built with Go (backend) and React/TypeScript (frontend), using GraphQL as the API layer.

## Architecture

- **Backend**: Go 1.25, entrypoint `cmd/stash/main.go`, second binary `cmd/phasher/main.go`
- **Frontend**: React + TypeScript + Vite, located in `ui/v2.5/`
- **API layer**: GraphQL via gqlgen — schema in `graphql/schema/*.graphql`, generated resolvers in `internal/api/`
- **Database**: SQLite with CGO (requires `CGO_ENABLED=1`), models in `internal/models/`
- **External integration**: stash-box client generated via gqlgenc (`.gqlgenc.yml`)
- **Code generation**: `make generate` runs both backend (`go generate`) and frontend (`npm run gqlgen`) codegen

## Key Directories

| Path | Purpose |
|------|---------|
| `cmd/stash/` | Main server binary |
| `cmd/phasher/` | Perceptual hash binary |
| `internal/api/` | GraphQL resolvers, API handlers, dataloaders |
| `internal/manager/` | Core session/manager, config, task management |
| `internal/models/` | Database models and repository interfaces |
| `pkg/` | Shared packages (session, sqlite, stashbox, logger, etc.) |
| `graphql/schema/` | GraphQL schema definition files |
| `ui/v2.5/` | React frontend (Vite + TypeScript) |
| `ui/v2.5/src/components/` | React UI components |
| `ui/v2.5/src/locales/` | i18n translation files |
| `docker/` | Docker build, CI, compiler, and production configs |
| `scripts/` | Build helper scripts |

## Build Commands

```bash
make pre-ui              # Install UI dependencies (pnpm install)
make generate            # Generate GraphQL backend + frontend code
make ui                  # Build the frontend
make stash               # Build the stash binary
make build-release       # Build release binaries (stripped, PIE)
make phasher             # Build the phasher binary
make lint                # Run golangci-lint
make test                # Run unit tests
make it                  # Run unit + integration tests
make fmt                 # Format Go code
make fmt-ui              # Format UI code (prettier)
make validate            # Run ALL PR checks (lint + tests + UI validation)
make validate-backend    # Backend-only PR checks
make validate-ui         # UI-only PR checks
make server-start        # Run dev server from .local/
make server-clean        # Remove .local/ dev directory
make ui-start            # Run UI dev server (port 3000)
```

## Development Workflow

1. `make pre-ui` — first time only, install dependencies
2. `make generate` — regenerate GraphQL code after schema changes
3. `make server-start` — start backend in one terminal
4. `make ui-start` — start frontend in another terminal
5. UI hot-reloads; backend changes require restart

After modifying GraphQL schema files in `graphql/schema/`, always run `make generate` to update generated Go resolvers and TypeScript types.

## PR Requirements

- `make validate` must pass (runs lint, tests, UI validation)
- Backend: `make validate-backend` (golangci-lint + integration tests)
- Frontend: `make validate-ui` (eslint, stylelint, tsc, prettier)

## Linting & Style

- Go: `golangci-lint` with config in `.golangci.yml` (gofmt with `simplify: false`)
- Frontend: ESLint + Stylelint + Prettier (config in `ui/v2.5/`)
- Mocks: generated via mockery (`.mockery.yml`)

## Go Build Tags

Default: `sqlite_stat4 sqlite_math_functions`
Static builds add: `sqlite_omit_load_extension osusergo netgo`

## Cross-Compilation

Uses Docker compiler container (`ghcr.io/stashapp/compiler:14`):
```bash
make start-compiler-container
docker exec -t build /bin/bash -c "make build-cc-<platform>"
```
Platforms: windows, macos, macos-intel, macos-arm, linux, linux-arm64v8, linux-arm32v7, linux-arm32v6, freebsd

## Notes

- CGO is required (`CGO_ENABLED=1`) — SQLite uses CGO bindings
- The frontend uses pnpm (not npm or yarn)
- Node.js >= 20 required
- GraphQL schema changes require `make generate` before building
- stash-box integration code is auto-generated from remote schema
