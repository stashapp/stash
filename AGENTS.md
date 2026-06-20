# Repository Guidelines

## Project Structure & Module Organization

Stash is a Go application with a React/TypeScript UI. CLI entry points live in `cmd/stash` and `cmd/phasher`. Backend packages are split between reusable code in `pkg/` and application internals in `internal/`. The main UI is in `ui/v2.5`, with source under `ui/v2.5/src`, public assets under `ui/v2.5/public`, and locale files under `ui/v2.5/src/locales`. Docker and release assets live in `docker/` and `scripts/`.

## Build, Test, and Development Commands

- `make pre-ui`: install UI dependencies with pnpm; run after cloning or dependency changes.
- `make generate`: generate Go and UI GraphQL files.
- `make ui`: build the frontend.
- `make stash`: build the main `stash` binary.
- `make build`: build the default binaries for local development.
- `make server-start`: run a local development server using `.local`.
- `make ui-start`: start the Vite UI on port `3000` or the next available port.
- `make validate`: run the full PR validation suite.
- `make lint`, `make it`, `make validate-ui`: run backend linting, backend tests, or UI checks separately.

## Coding Style & Naming Conventions

Format Go code with `make fmt`; the project uses `gofmt` without simplification and `golangci-lint` via `.golangci.yml`. Place Go tests beside implementation files as `*_test.go`. UI code uses Biome and Stylelint; run `make fmt-ui` or, inside `ui/v2.5`, `pnpm run format`. UI formatting uses spaces, double quotes, and ES5 trailing commas. Keep package, component, and file names consistent with nearby code.

## Testing Guidelines

Backend tests use Go's standard test framework plus project helpers; run all unit and integration tests with `make it` or targeted tests with `go test ./pkg/scene -run TestName`. UI validation is `make validate-ui`, which runs linting, TypeScript checks, and format checks. Add focused tests for behavioral changes. Produce generated-code changes through `make generate`; do not edit them manually.

## Commit & Pull Request Guidelines

Recent history uses concise imperative subjects such as `Fix: Pagination Footer Centering With Sidebar` and `Add User-Agent header for stash-box requests`. Keep commits focused and descriptive. Pull requests must address one issue or feature, link an open issue, include sufficient tests, and describe manual verification steps. For UI changes, include screenshots or a short recording when visual behavior changes.

## Security & Configuration Tips

Do not commit local runtime data from `.local`, generated media, credentials, or API keys. Follow `docs/AI_POLICY.md` and `docs/CONTRIBUTING.md` for contribution and AI-use requirements.

## Agent-Specific Instructions

Always answer the user in Chinese.
Write all documentation generated for the user in Chinese as well.
