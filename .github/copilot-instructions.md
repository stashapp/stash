## Quick context for an AI coding agent

This repository is a Go web application with a separate UI (React/Vite) in `ui/v2.5` and GraphQL used as the API surface. Key entry points are `./cmd/stash` (main webapp) and `./cmd/phasher` (helper binary). Generated GraphQL code and UI GraphQL types are present; do not modify generated files directly—regenerate instead (see below).

Key locations:
- Backend entry points: `cmd/stash/main.go`, `cmd/phasher/main.go`
- GraphQL schema: `graphql/schema/schema.graphql`, `gqlgen.yml`
- Backend resolver/translator code: `internal/resolver_*.go`, `internal/resolver_model_*.go`, `internal/api` (loaders, auth)
- UI: `ui/v2.5` (Yarn-based; run `yarn build` / `yarn start` from that dir)
- Makefile: top-level `Makefile` drives build, test, codegen and common developer workflows

Big-picture architecture and patterns
- Backend: monolithic Go server exposing a GraphQL API. Business logic and models live under `internal/` and `pkg/`.
- Frontend: single-page app in `ui/v2.5` that talks to the backend API (default platform URL http://localhost:9999). The UI build artifacts are copied into the Go app at build time (look for `ui/v2.5/build`).
- Code generation: GraphQL code is generated for both backend (Go) and frontend (TypeScript). The Makefile orchestrates generation (`make generate`, `make generate-ui`, `make generate-backend`).
- Tests: unit tests are `go test ./...` and integration tests use the `integration` build tag (run with `make it`).

Developer workflows (use `make` targets rather than ad-hoc commands)
- Setup UI deps: `make pre-ui` (runs `yarn install` in `ui/v2.5`).
- Generate code: `make generate` (runs both backend and UI generation). If you change GraphQL schema, run `make generate`.
- Build UI: `make ui` or `make ui-only` (uses Yarn/Vite). For local dev run `make ui-start` to run UI in dev mode.
- Run local server: `make server-start` (creates `.local` and runs the backend using `config.yml`). Use `make server-clean` to clear `.local`.
- Run tests: `make test` (unit), `make it` (including integration tests). `make validate` runs lint and tests required for PRs.
- Build binaries: `make build` (debug), `make build-release` (release flags). Several cross-compile targets exist (`build-cc-*`) that are intended for the compiler docker image.

Important Makefile details to follow exactly
- Use `flags-release`, `flags-pie`, `flags-static` variants via the Makefile when creating release or static builds (examples in `Makefile` header).
- Integration tests: `make it` appends `integration` to `GO_BUILD_TAGS` while running `go test -tags "$(GO_BUILD_TAGS)" ./...`.

GraphQL & generation rules
- Schema lives in `graphql/schema/schema.graphql`. To change GraphQL types, update schema and run:
  - `make pre-ui` (once after clone or when UI deps change)
  - `make generate` (or `make generate-ui` + `make generate-backend`)
- Backend generation uses `go generate ./cmd/stash`. UI generation uses `cd ui/v2.5 && yarn run gqlgen`.
- Do not edit generated files under `internal/resolver_model_*.go` or files in `ui/v2.5/build` — repeat codegen instead.

Conventions & project-specific patterns
- Package layout: core server code lives under `internal/` (private API surface). Reusable libraries and utilities appear in `pkg/`.
- Resolver naming: GraphQL resolver-model files follow `resolver_model_<type>.go` naming in `internal/` — search there for examples when adding fields.
- Config: the dev server uses `config.yml` and writes runtime state under `.local` when using `make server-start`.
- Plugins/scrapers: extensible plugin system; look at `pkg/plugin` and `internal/plugin_map.go` for integration patterns.

Debugging & profiling
- Start server with: `make server-start` then attach debugger or use logs. For CPU profiling, run the binary with `--cpuprofile <file>` (see README: Profiling section) and open with `go tool pprof`.

Testing & PR checks
- Pre-PR checks: `make validate` runs `golangci-lint` and tests. Frontend checks are `make validate-ui`.
- Quick-check changed UI code: `make validate-ui-quick` and `make fmt-ui-quick` exist for PRs to speed checks.

Examples (explicit)
- Full local dev: `make pre-ui && make generate && make server-start` (in one terminal), `make ui-start` (in another). Open `http://localhost:3000` for the UI dev server.
- Run all tests and linters: `make validate`
- Regenerate after schema change: `make generate`

If something is missing or ambiguous
- Point me at the file or behavior you'd like clarified (for example, a new GraphQL field, or where to add an integration test). I can update this guide with explicit snippets and commands.

References (examples):
- `Makefile` — drive builds, codegen, and test targets
- `graphql/schema/schema.graphql` and `gqlgen.yml` — GraphQL schema and generation config
- `internal/resolver_model_*.go` — generated resolver models (do not edit directly)
- `cmd/stash/main.go` — backend entrypoint

Please review and tell me any missing workflows or files you'd like called out; I'll iterate the guide accordingly.
