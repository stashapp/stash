# Generate GraphQL Code

Regenerate all GraphQL backend and frontend code after schema changes.

## When to use
- After modifying any file in `graphql/schema/`
- After changing resolver signatures in `internal/api/`
- When starting development on a fresh clone
- When switching branches that have schema changes

## Steps

1. Ensure UI dependencies are installed:
   ```bash
   cd ui/v2.5 && pnpm install --frozen-lockfile
   ```

2. Run full generation (backend + frontend):
   ```bash
   make generate
   ```

   Or run separately:
   ```bash
   make generate-backend   # Go code generation (go generate ./cmd/stash)
   make generate-ui        # Frontend TypeScript codegen (npm run gqlgen)
   ```

3. If you only need to regenerate dataloaders:
   ```bash
   make generate-dataloaders
   ```

4. If you need to regenerate stash-box client:
   ```bash
   make generate-stash-box-client
   ```

5. Verify no compilation errors:
   ```bash
   go build ./cmd/stash
   cd ui/v2.5 && npx tsc --noEmit
   ```

## Common issues
- If `make generate-backend` fails with "no such package", run `go mod tidy` first
- If `make generate-ui` fails, ensure `pnpm install` has been run in `ui/v2.5/`
- Generated files should NOT be manually edited — they are overwritten on regeneration
