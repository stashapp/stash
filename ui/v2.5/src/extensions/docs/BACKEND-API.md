# Backend API Reference

This document defines all custom backend endpoints and modifications required for the extension system to function.

> ⚠️ **Important**: These backend changes must be preserved when merging upstream updates.
>
> 📁 **See also**: `/patches/` directory for merge instructions and patch files.

## Custom GraphQL Queries

### Facets Endpoints

The facets system adds 6 new GraphQL queries for fetching aggregated filter counts:

| Query | Purpose | Frontend Hook |
|-------|---------|---------------|
| `sceneFacets` | Scene filter counts | `useSceneFacetCounts` |
| `performerFacets` | Performer filter counts | `usePerformerFacetCounts` |
| `galleryFacets` | Gallery filter counts | `useGalleryFacetCounts` |
| `groupFacets` | Group filter counts | `useGroupFacetCounts` |
| `studioFacets` | Studio filter counts | `useStudioFacetCounts` |
| `tagFacets` | Tag filter counts | `useTagFacetCounts` |

### Query Signatures

```graphql
# graphql/schema/schema.graphql

type Query {
  # Scene facets with lazy loading options
  sceneFacets(
    scene_filter: SceneFilterType
    limit: Int
    include_performer_tags: Boolean  # Expensive, load on-demand
    include_captions: Boolean         # Expensive, load on-demand
  ): SceneFacetsResult!
  
  # Standard facets queries
  performerFacets(performer_filter: PerformerFilterType, limit: Int): PerformerFacetsResult!
  galleryFacets(gallery_filter: GalleryFilterType, limit: Int): GalleryFacetsResult!
  groupFacets(group_filter: GroupFilterType, limit: Int): GroupFacetsResult!
  studioFacets(studio_filter: StudioFilterType, limit: Int): StudioFacetsResult!
  tagFacets(tag_filter: TagFilterType, limit: Int): TagFacetsResult!
}
```

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `*_filter` | FilterType | `null` | Filter criteria (same as list queries) |
| `limit` | Int | `100` | Max facets per category |
| `include_performer_tags` | Boolean | `false` | Load performer tags (scene only) |
| `include_captions` | Boolean | `false` | Load caption languages (scene only) |

---

## Custom GraphQL Types

### Facet Count Types

```graphql
# graphql/schema/types/facets.graphql

# Entity facet (tags, performers, studios, etc.)
type FacetCount {
  id: ID!
  label: String!   # Display name (critical for UI)
  count: Int!
}

# Boolean facet (organized, favorite, etc.)
type BooleanFacetCount {
  value: Boolean!
  count: Int!
}

# Video resolution facet
type ResolutionFacetCount {
  resolution: ResolutionEnum!
  count: Int!
}

# Video orientation facet
type OrientationFacetCount {
  orientation: OrientationEnum!
  count: Int!
}

# Gender facet
type GenderFacetCount {
  gender: GenderEnum!
  count: Int!
}

# Rating facet (1-5 stars as 20-100)
type RatingFacetCount {
  rating: Int!
  count: Int!
}

# Caption language facet
type CaptionFacetCount {
  language: String!
  count: Int!
}

# Circumcised status facet
type CircumcisedFacetCount {
  value: CircumisedEnum!
  count: Int!
}
```

### Result Types

```graphql
# Scene facets result
type SceneFacetsResult {
  tags: [FacetCount!]!
  performers: [FacetCount!]!
  studios: [FacetCount!]!
  groups: [FacetCount!]!
  performer_tags: [FacetCount!]        # Only if include_performer_tags=true
  resolutions: [ResolutionFacetCount!]!
  orientations: [OrientationFacetCount!]!
  organized: [BooleanFacetCount!]!
  interactive: [BooleanFacetCount!]!
  ratings: [RatingFacetCount!]!
  captions: [CaptionFacetCount!]       # Only if include_captions=true
}

# Performer facets result
type PerformerFacetsResult {
  tags: [FacetCount!]!
  studios: [FacetCount!]!
  genders: [GenderFacetCount!]!
  countries: [FacetCount!]!
  circumcised: [CircumcisedFacetCount!]!
  favorite: [BooleanFacetCount!]!
  ratings: [RatingFacetCount!]!
}

# Gallery facets result
type GalleryFacetsResult {
  tags: [FacetCount!]!
  performers: [FacetCount!]!
  studios: [FacetCount!]!
  organized: [BooleanFacetCount!]!
  ratings: [RatingFacetCount!]!
}

# Group facets result
type GroupFacetsResult {
  tags: [FacetCount!]!
  performers: [FacetCount!]!
  studios: [FacetCount!]!
}

# Studio facets result
type StudioFacetsResult {
  tags: [FacetCount!]!
  parents: [FacetCount!]!
  favorite: [BooleanFacetCount!]!
}

# Tag facets result
type TagFacetsResult {
  parents: [FacetCount!]!
  children: [FacetCount!]!
  favorite: [BooleanFacetCount!]!
}
```

---

## Backend File Reference

### New Files (Must Preserve)

| File | Purpose | Lines |
|------|---------|-------|
| `graphql/schema/types/facets.graphql` | GraphQL type definitions | ~100 |
| `pkg/models/facets.go` | Go model structs | ~150 |
| `pkg/sqlite/scene_facets.go` | Scene facets SQLite implementation | ~400 |
| `pkg/sqlite/performer_facets.go` | Performer facets SQLite implementation | ~250 |
| `pkg/sqlite/gallery_facets.go` | Gallery facets SQLite implementation | ~200 |
| `pkg/sqlite/group_facets.go` | Group facets SQLite implementation | ~150 |
| `pkg/sqlite/studio_facets.go` | Studio facets SQLite implementation | ~150 |
| `pkg/sqlite/tag_facets.go` | Tag facets SQLite implementation | ~100 |
| `internal/api/resolver_query_facets.go` | GraphQL resolvers | ~200 |
| `internal/api/types_facets.go` | API type mappings | ~100 |

### Modified Files (Merge Carefully)

| File | Changes |
|------|---------|
| `graphql/schema/schema.graphql` | Added facet + recommendation queries |
| `pkg/models/repository_scene.go` | Added `SceneFaceter` interface |
| `pkg/models/repository_performer.go` | Added `PerformerFaceter` interface |
| `pkg/models/repository_gallery.go` | Added `GalleryFaceter` interface |
| `pkg/models/repository_group.go` | Added `GroupFaceter` interface |
| `pkg/models/repository_studio.go` | Added `StudioFaceter` interface |
| `pkg/models/repository_tag.go` | Added `TagFaceter` interface + `FindFavoriteTagIDs` |
| `graphql/schema/types/filters.graphql` | Added tag filter fields |
| `pkg/models/tag.go` | Added `PerformersFilter`, `GroupsFilter` |
| `pkg/sqlite/tag.go` | Added join repos + `FindFavoriteTagIDs` |
| `pkg/sqlite/tag_filter.go` | Added performer/group filter handlers |
| `pkg/models/resolution.go` | Added `ResolutionFromHeight` |

### Test Files

| File | Tests |
|------|-------|
| `pkg/sqlite/scene_facets_test.go` | 14 integration tests |
| `pkg/sqlite/performer_facets_test.go` | 9 integration tests |
| `pkg/sqlite/gallery_facets_test.go` | 8 integration tests |
| `pkg/sqlite/group_facets_test.go` | 6 integration tests |
| `pkg/sqlite/studio_facets_test.go` | 6 integration tests |
| `pkg/sqlite/tag_facets_test.go` | 5 integration tests |
| `internal/api/resolver_query_facets_test.go` | 18 unit tests |

---

## Recommendations System

Provides personalized recommendations for scenes and performers based on user viewing history and preferences.

### GraphQL Queries

```graphql
# Scene recommendations
sceneRecommendations(limit: Int): SceneRecommendationsResultType!
sceneRecommendationsForScene(scene_id: ID!, limit: Int): SceneRecommendationsResultType!

# Performer recommendations  
performerRecommendations(limit: Int): PerformerRecommendationsResultType!
performerRecommendationsForPerformer(performer_id: ID!, limit: Int): PerformerRecommendationsResultType!
```

### Implementation Files

| File | Purpose | Lines |
|------|---------|-------|
| `pkg/recommendation/scene.go` | Scene recommendation logic | ~900 |
| `pkg/recommendation/performer.go` | Performer recommendation logic | ~525 |
| `internal/api/resolver_query_scene_recommendations.go` | Scene resolver | - |
| `internal/api/resolver_query_performer_recommendations.go` | Performer resolver | - |
| `internal/api/resolver_scene_recommendations_result_type.go` | Result type | - |
| `internal/api/resolver_performer_recommendations_result_type.go` | Result type | - |

### Algorithm Overview

The recommendation system analyzes:
- User viewing history (play counts, O counts)
- Favorite performers and tags
- Similar content based on shared attributes
- Recency weighting

---

## Tag Filter Extensions

Adds ability to filter tags by their related performers and groups.

### New Filter Fields

```graphql
input TagFilterType {
  # ... existing fields ...
  performers_filter: PerformerFilterType
  groups_filter: GroupFilterType
}
```

### Implementation Files

| File | Changes |
|------|---------|
| `graphql/schema/types/filters.graphql` | Added filter fields |
| `pkg/models/tag.go` | Added Go struct fields |
| `pkg/sqlite/tag.go` | Added join repositories |
| `pkg/sqlite/tag_filter.go` | Added filter handlers |

### Use Cases

- Find tags used by performers from a specific country
- Find tags associated with groups from a specific studio
- Complex multi-level filtering

---

## Utility Functions

### ResolutionFromHeight

**File:** `pkg/models/resolution.go`

```go
func ResolutionFromHeight(height int) ResolutionEnum
```

Converts a video height to the matching resolution enum. Used by facets system.

### FindFavoriteTagIDs

**File:** `pkg/sqlite/tag.go`

```go
func (qb *TagStore) FindFavoriteTagIDs(ctx context.Context) ([]int, error)
```

Efficiently fetches IDs of all favorite tags. Used by tag facets.

---

## Repository Interfaces

Each entity type has a Faceter interface added to its repository:

```go
// pkg/models/repository_scene.go

type SceneFaceter interface {
    GetFacets(ctx context.Context, filter *SceneFilterType, limit int, options SceneFacetOptions) (*SceneFacets, error)
}

// SceneFacetOptions controls lazy loading of expensive facets
type SceneFacetOptions struct {
    IncludePerformerTags bool
    IncludeCaptions      bool
}
```

```go
// pkg/models/repository_performer.go

type PerformerFaceter interface {
    GetFacets(ctx context.Context, filter *PerformerFilterType, limit int) (*PerformerFacets, error)
}
```

---

## SQLite Implementation

### Query Strategy

Facets use CTE (Common Table Expression) queries for efficiency:

```sql
-- Base filter applied ONCE
WITH filtered_scenes AS (
    SELECT DISTINCT scenes.id FROM scenes
    WHERE ... -- all filter criteria
)

-- Tag counts
SELECT 'tag' as facet_type, t.id, t.name as label, 
       COUNT(DISTINCT st.scene_id) as count
FROM filtered_scenes fs
INNER JOIN scenes_tags st ON fs.id = st.scene_id
INNER JOIN tags t ON st.tag_id = t.id
GROUP BY t.id
ORDER BY count DESC
LIMIT ?

UNION ALL

-- Performer counts (same pattern)
SELECT 'performer' as facet_type, p.id, p.name as label,
       COUNT(DISTINCT sp.scene_id) as count
FROM filtered_scenes fs
INNER JOIN performers_scenes sp ON fs.id = sp.scene_id
INNER JOIN performers p ON sp.performer_id = p.id
GROUP BY p.id
ORDER BY count DESC
LIMIT ?

-- ... more facet types
```

### Lazy Loading Implementation

Expensive facets run in goroutines:

```go
// pkg/sqlite/scene_facets.go

func (qb *SceneStore) GetFacets(ctx context.Context, filter *SceneFilterType, 
    limit int, options SceneFacetOptions) (*SceneFacets, error) {
    
    var wg sync.WaitGroup
    var mu sync.Mutex
    result := &SceneFacets{}
    
    // Core facets always run
    wg.Add(1)
    go func() {
        defer wg.Done()
        qb.getCoreFacets(ctx, baseSQL, baseArgs, limit, result, &mu)
    }()
    
    // Expensive facets only if requested
    if options.IncludePerformerTags {
        wg.Add(1)
        go func() {
            defer wg.Done()
            qb.getPerformerTagsFacet(ctx, baseSQL, baseArgs, limit, result, &mu)
        }()
    }
    
    if options.IncludeCaptions {
        wg.Add(1)
        go func() {
            defer wg.Done()
            qb.getCaptionsFacet(ctx, baseSQL, baseArgs, result, &mu)
        }()
    }
    
    wg.Wait()
    return result, nil
}
```

---

## Frontend GraphQL Queries

The frontend queries are defined in:

```graphql
# ui/v2.5/graphql/data/facets.graphql

query SceneFacets(
  $scene_filter: SceneFilterType
  $limit: Int
  $includePerformerTags: Boolean!
  $includeCaptions: Boolean!
) {
  sceneFacets(
    scene_filter: $scene_filter
    limit: $limit
    include_performer_tags: $includePerformerTags
    include_captions: $includeCaptions
  ) {
    tags { id label count }
    performers { id label count }
    studios { id label count }
    groups { id label count }
    performer_tags @include(if: $includePerformerTags) { id label count }
    resolutions { resolution count }
    orientations { orientation count }
    organized { value count }
    interactive { value count }
    ratings { rating count }
    captions @include(if: $includeCaptions) { language count }
  }
}

# Similar queries for other entity types...
```

---

## Testing

### Run Backend Tests

```bash
# Integration tests (requires test database)
go test -v -tags=integration ./pkg/sqlite/... -run Facet

# Unit tests
go test -v ./internal/api/... -run Facet

# All facet tests
go test -v -tags=integration ./... -run Facet
```

### Test Coverage

| Component | Tests | Coverage |
|-----------|-------|----------|
| Scene Facets | 14 | All facet types, filters, lazy loading |
| Performer Facets | 9 | Tags, genders, studios, countries |
| Gallery Facets | 8 | Tags, performers, studios, organized |
| Group Facets | 6 | Tags, performers, studios |
| Studio Facets | 6 | Tags, parents, favorite |
| Tag Facets | 5 | Parents, children, favorite |
| API Resolvers | 18 | Type conversions, error handling |

---

## Merge Guide

When merging upstream changes, use the `/patches/` directory for detailed instructions.

### Quick Reference

| Step | Resource |
|------|----------|
| Schema queries | `/patches/schema-queries.md` |
| Repository interfaces | `/patches/repository-interfaces.md` |
| Overview | `/patches/README.md` |

### 1. Preserve New Files

These files don't exist upstream - they won't conflict:
```
graphql/schema/types/facets.graphql
pkg/models/facets.go
pkg/models/facets_interfaces.go      # NEW: Interface definitions
pkg/sqlite/*_facets.go
pkg/sqlite/*_facets_test.go
internal/api/resolver_query_facets.go
internal/api/types_facets.go
```

### 2. Re-add Schema Queries

If `graphql/schema/schema.graphql` conflicts:

```bash
# See full query definitions
cat patches/schema-queries.md
```

Add after `findTags` query:
```graphql
  sceneFacets(...): SceneFacetsResult!
  performerFacets(...): PerformerFacetsResult!
  galleryFacets(...): GalleryFacetsResult!
  groupFacets(...): GroupFacetsResult!
  studioFacets(...): StudioFacetsResult!
  tagFacets(...): TagFacetsResult!
```

### 3. Re-add Repository Interfaces

If `pkg/models/repository_*.go` files conflict:

```bash
# See exact lines to add
cat patches/repository-interfaces.md
```

Add Faceter interface to each Reader:
```go
type SceneReader interface {
    // ... existing interfaces ...
    SceneFaceter  // <-- Add this line
}
```

### 4. Regenerate GraphQL

After resolving conflicts:
```bash
go generate ./...
```

### 5. Run Tests

```bash
# Backend
go build ./...
go test -v -tags=integration ./pkg/sqlite/... -run Facet
go test -v ./internal/api/...

# Frontend
cd ui/v2.5
yarn build
yarn test
```

---

## Troubleshooting

### "Unknown field 'sceneFacets'"

**Cause**: GraphQL schema not regenerated after merge
**Fix**: Run `go generate ./...`

### "No such table: performers_galleries"

**Cause**: Wrong join table name (was `galleries_performers`)
**Fix**: Use correct table name `performers_galleries`

### Slow facets on large database

**Cause**: Missing indexes or expensive filter
**Mitigations**:
- Use `limit` parameter (default 100)
- Enable lazy loading for expensive facets
- Add database indexes on frequently filtered columns

### Empty facet counts

**Cause**: Filter too restrictive or no matching items
**Check**: Try with no filter to verify endpoint works

