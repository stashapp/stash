# Create a Database Migration

Guide for creating and managing SQLite database migrations.

## When to use
- Adding new tables or columns
- Modifying existing schema
- Adding indexes or constraints

## Steps

### 1. Understand the migration system

- Migrations are in `pkg/sqlite/migrations/`
- Each migration is a Go file with an upward function
- Migrations are numbered sequentially
- SQLite is the database — be aware of SQLite-specific limitations

### 2. Create a new migration

Study existing migration files for the pattern. Each migration:
- Has a unique version number
- Implements a `up` function that executes SQL
- Cannot be rolled back in production

### 3. Write the migration SQL

SQLite limitations to keep in mind:
- No `ALTER TABLE DROP COLUMN` (before SQLite 3.35.0)
- Limited `ALTER TABLE` support — often requires table rebuild pattern
- Use `IF NOT EXISTS` / `IF EXISTS` for safety
- Prefer additive changes (add columns, add indexes)

### 4. Update the repository layer

- Add the new field to the model in `internal/models/`
- Update the repository in `pkg/sqlite/` to read/write the new column
- Update any relevant query builders

### 5. Update tests

- Add test cases for new fields in the relevant test file
- Integration tests are in the same package, gated by `integration` build tag

### 6. Validate

```bash
make it   # Run integration tests
```

## Common patterns

### Adding a nullable column
```sql
ALTER TABLE scenes ADD COLUMN my_column TEXT;
```

### Adding a column with default
```sql
ALTER TABLE scenes ADD COLUMN my_column TEXT NOT NULL DEFAULT '';
```

### Adding an index
```sql
CREATE INDEX IF NOT EXISTS index_scenes_my_column ON scenes(my_column);
```

## Notes
- SQLite uses CGO bindings — `CGO_ENABLED=1` is required
- The `sqlite_stat4` and `sqlite_math_functions` build tags are enabled by default
- Test migrations by running `make it` (integration tests)
