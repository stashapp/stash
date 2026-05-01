# Validate Changes for PR

Run all checks required before submitting a pull request.

## When to use
- Before submitting a PR
- After making code changes to verify nothing is broken
- During development to catch issues early

## Steps

### Full validation (all checks)

```bash
make validate
```

This runs both `validate-backend` and `validate-ui`.

### Backend-only validation

```bash
make validate-backend
```

This runs `golangci-lint` + integration tests.

### Frontend-only validation

```bash
make validate-ui
```

This runs ESLint, Stylelint, TypeScript check, and Prettier.

### Quick validation (changed files only, experimental)

For quick frontend iteration:

```bash
make fmt-ui-quick          # Format only changed files
make validate-ui-quick     # Lint only changed files (no tsc)
```

### Run specific checks

```bash
make lint                  # golangci-lint only
make test                  # Unit tests only
make it                    # Unit + integration tests
make fmt                   # Format Go code
make fmt-ui                # Format UI code
```

### Run tests for a specific package

```bash
go test ./internal/api/...
go test ./pkg/sqlite/...
```

## CI requirements
The GitHub Actions workflow (`build.yml`) runs:
1. `make generate` — code generation
2. `make validate-ui` — frontend checks
3. `make ui` — frontend build
4. `make validate-backend` — lint + integration tests
5. Cross-platform build verification

All these must pass for a PR to be accepted.
