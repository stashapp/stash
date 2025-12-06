# Utility Function Additions

Small utility functions added to support fork features.

---

## 1. ResolutionFromHeight

**File:** `pkg/models/resolution.go`

Add this function at the end of the file:

```go
// ResolutionFromHeight returns the ResolutionEnum that matches the given height.
// Returns empty string if no matching resolution is found.
func ResolutionFromHeight(height int) ResolutionEnum {
	for res, r := range resolutionRanges {
		if height >= r.min && height <= r.max {
			return res
		}
	}
	return ""
}
```

**Purpose:** Used by facets system to categorize videos by resolution.

---

## 2. FindFavoriteTagIDs

**File:** `pkg/sqlite/tag.go`

Add this method to `TagStore`:

```go
func (qb *TagStore) FindFavoriteTagIDs(ctx context.Context) ([]int, error) {
	query := `SELECT id FROM tags WHERE favorite = 1`
	var ret []int
	if err := tagRepository.queryFunc(ctx, query, nil, false, func(r *sqlx.Rows) error {
		var id int
		if err := r.Scan(&id); err != nil {
			return err
		}
		ret = append(ret, id)
		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}
```

**Purpose:** Efficiently fetch IDs of all favorite tags for facet filtering.

**Note:** Also add to the `TagFinder` interface in `pkg/models/repository_tag.go`:

```go
type TagFinder interface {
	// ... existing methods ...
	
	// FindFavoriteTagIDs returns IDs of all favorited tags
	FindFavoriteTagIDs(ctx context.Context) ([]int, error)
}
```

---

## 3. Random Sort Testing Helper

**File:** `pkg/sqlite/sql.go`

Add this function (used for development/testing only):

```go
func getRandomSortForDevTesting(tableName string, direction string, seed uint64) string {
	// For seeded random (reproducible), use a hash-based approach with smaller multipliers
	// to avoid integer overflow in SQLite which would cause precision loss
	if seed != 0 {
		// cap seed at 10^6 for smaller numbers
		seed %= 1e6

		colName := getColumn(tableName, "id")

		// Use smaller prime multipliers to avoid overflow
		// This provides pseudo-random ordering that's reproducible with the same seed
		// The formula: ((id * seed) % largePrime + id) % anotherPrime
		return fmt.Sprintf(" ORDER BY ((%[1]s * %[2]d) %% 999983 + %[1]s) %% 999979 %[3]s", colName, seed, direction)
	}

	// For unseeded random (non-reproducible), use SQLite's native RANDOM() function
	// which provides true randomness without overflow issues
	return fmt.Sprintf(" ORDER BY RANDOM() %s", direction)
}
```

**Purpose:** Alternative random sort implementation for testing. Currently commented out in production.

**Note:** There's also a commented line in `getSort()`:
```go
case strings.Compare(sort, "random") == 0:
	return getRandomSort(tableName, direction, rand.Uint64())
	// return getRandomSortNative(tableName, direction, 0)
```

---

## Summary

| Function | File | Used By |
|----------|------|---------|
| `ResolutionFromHeight` | `pkg/models/resolution.go` | Facets system |
| `FindFavoriteTagIDs` | `pkg/sqlite/tag.go` | Tag facets |
| `getRandomSortForDevTesting` | `pkg/sqlite/sql.go` | Dev/testing |

These are minor additions that can be easily identified in a diff and re-added if lost during merge.

