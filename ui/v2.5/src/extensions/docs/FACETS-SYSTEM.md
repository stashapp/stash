# Facets System Technical Reference

This document provides in-depth technical details about the facets aggregation system.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend                                 │
│  ┌─────────────────┐    ┌──────────────────┐                   │
│  │ List Page       │───▶│ FacetCountsContext│                   │
│  │ (SceneList)     │    │ Provider          │                   │
│  └─────────────────┘    └────────┬─────────┘                   │
│                                  │                              │
│  ┌─────────────────┐    ┌────────▼─────────┐                   │
│  │ Filter Component│◀───│ useFacetCounts   │                   │
│  │ (SidebarTags)   │    │ Hook             │                   │
│  └─────────────────┘    └────────┬─────────┘                   │
│                                  │                              │
└──────────────────────────────────┼──────────────────────────────┘
                                   │ GraphQL Query
                                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Backend                                  │
│  ┌─────────────────┐    ┌──────────────────┐                   │
│  │ GraphQL Resolver│───▶│ Repository       │                   │
│  │ (sceneFacets)   │    │ GetFacets()      │                   │
│  └─────────────────┘    └────────┬─────────┘                   │
│                                  │                              │
│                         ┌────────▼─────────┐                   │
│                         │ SQLite CTE Query │                   │
│                         │ (UNION ALL)      │                   │
│                         └──────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
```

## Backend Implementation

### GraphQL Schema

```graphql
# graphql/schema/types/facets.graphql

type FacetCount {
  id: ID!
  label: String!
  count: Int!
}

type BooleanFacetCount {
  value: Boolean!
  count: Int!
}

type ResolutionFacetCount {
  resolution: ResolutionEnum!
  count: Int!
}
# ... etc
```

### Query Endpoints

```graphql
type Query {
  sceneFacets(
    scene_filter: SceneFilterType
    limit: Int
    include_performer_tags: Boolean  # Lazy loading
    include_captions: Boolean         # Lazy loading
  ): SceneFacetsResult!
  
  performerFacets(performer_filter: PerformerFilterType, limit: Int): PerformerFacetsResult!
  galleryFacets(gallery_filter: GalleryFilterType, limit: Int): GalleryFacetsResult!
  groupFacets(group_filter: GroupFilterType, limit: Int): GroupFacetsResult!
  studioFacets(studio_filter: StudioFilterType, limit: Int): StudioFacetsResult!
  tagFacets(tag_filter: TagFilterType, limit: Int): TagFacetsResult!
}
```

### SQLite CTE Implementation

Facets are computed using a single CTE query:

```sql
WITH filtered_scenes AS (
    SELECT DISTINCT scenes.id FROM scenes
    WHERE ... -- base filter applied once
)

-- Tag counts
SELECT 'tag' as facet_type, t.id, t.name as label, COUNT(DISTINCT st.scene_id) as count
FROM filtered_scenes fs
INNER JOIN scenes_tags st ON fs.id = st.scene_id
INNER JOIN tags t ON st.tag_id = t.id
GROUP BY t.id
ORDER BY count DESC
LIMIT ?

UNION ALL

-- Performer counts
SELECT 'performer' as facet_type, p.id, p.name as label, COUNT(DISTINCT sp.scene_id) as count
FROM filtered_scenes fs
INNER JOIN performers_scenes sp ON fs.id = sp.scene_id
INNER JOIN performers p ON sp.performer_id = p.id
GROUP BY p.id
ORDER BY count DESC
LIMIT ?

-- ... more facet types
```

### Lazy Loading

Some facets are expensive (multiple joins):

| Facet | Joins | Lazy Loaded |
|-------|-------|-------------|
| `performer_tags` | 3 | Yes |
| `captions` | 2 | Yes |
| `tags`, `performers`, `studios` | 2 | No |
| `organized`, `interactive`, `rating` | 1 | No |

```go
// pkg/sqlite/scene_facets.go
func (qb *SceneStore) GetFacets(ctx context.Context, filter *SceneFilterType, limit int, options SceneFacetOptions) (*SceneFacets, error) {
    // Core facets always run
    go qb.getCoreFacets(ctx, baseSQL, baseArgs, limit, result, &mu)
    
    // Expensive facets only if requested
    if options.IncludePerformerTags {
        go qb.getPerformerTagsFacet(ctx, baseSQL, baseArgs, limit, result, &mu)
    }
}
```

## Frontend Implementation

### Hooks

```typescript
// src/extensions/hooks/useFacetCounts.ts

interface UseFacetCountsOptions {
  isOpen: boolean;           // Only fetch when sidebar is open
  debounceMs?: number;       // Debounce filter changes (default: 300ms)
  includePerformerTags?: boolean;  // Lazy load
  includeCaptions?: boolean;       // Lazy load
}

const { counts, loading } = useSceneFacetCounts(filter, {
  isOpen: showSidebar,
  debounceMs: 300,
  includePerformerTags: sectionOpen["performer_tags"],
});
```

### Context Provider

```tsx
<FacetCountsContext.Provider value={{ counts: facetCounts, loading: facetLoading }}>
  <SidebarPane>
    {/* Filter components can access counts via useContext */}
  </SidebarPane>
</FacetCountsContext.Provider>
```

### State Preservation During Lazy Loading

When lazy-loaded facets are fetched, use partial state updates:

```typescript
if (isLazyLoadUpdate) {
  // Only update the lazy-loaded facets, preserve everything else
  setCounts((prev) => ({
    ...prev,  // Keep existing tags, performers, studios, etc.
    performerTags: includePerformerTags ? toMap(facets.performer_tags) : prev.performerTags,
    captions: includeCaptions ? toCaptionMap(facets.captions) : prev.captions,
  }));
}
```

## Supported Facets by Entity

### Scene Facets
| Facet | Type | Lazy |
|-------|------|------|
| `tags` | FacetCount | No |
| `performers` | FacetCount | No |
| `studios` | FacetCount | No |
| `groups` | FacetCount | No |
| `performer_tags` | FacetCount | **Yes** |
| `resolutions` | ResolutionFacetCount | No |
| `orientations` | OrientationFacetCount | No |
| `organized` | BooleanFacetCount | No |
| `interactive` | BooleanFacetCount | No |
| `ratings` | RatingFacetCount | No |
| `captions` | CaptionFacetCount | **Yes** |

### Performer Facets
| Facet | Type |
|-------|------|
| `tags` | FacetCount |
| `studios` | FacetCount |
| `genders` | GenderFacetCount |
| `countries` | FacetCount |
| `circumcised` | CircumcisedFacetCount |
| `favorite` | BooleanFacetCount |
| `ratings` | RatingFacetCount |

### Gallery Facets
| Facet | Type |
|-------|------|
| `tags` | FacetCount |
| `performers` | FacetCount |
| `studios` | FacetCount |
| `organized` | BooleanFacetCount |
| `ratings` | RatingFacetCount |

### Group Facets
| Facet | Type |
|-------|------|
| `tags` | FacetCount |
| `performers` | FacetCount |
| `studios` | FacetCount |

### Studio Facets
| Facet | Type |
|-------|------|
| `tags` | FacetCount |
| `parents` | FacetCount |
| `favorite` | BooleanFacetCount |

### Tag Facets
| Facet | Type |
|-------|------|
| `parents` | FacetCount |
| `children` | FacetCount |
| `favorite` | BooleanFacetCount |

## Performance Considerations

### Large Databases (>100k items)

1. **Lazy Loading**: Only load expensive facets when their filter section is expanded
2. **Debouncing**: Delay facet fetches when filter changes rapidly (300ms default)
3. **Caching**: Facet results are cached until filter changes
4. **Limit**: Default limit of 100 facets per category

### Query Optimization

The CTE-based query structure ensures:
- Base filter executes only once
- Each facet dimension uses the same filtered set
- Results are limited per category

## Testing

### Backend Tests

```bash
# Run facet integration tests
go test -v -tags=integration ./pkg/sqlite/... -run Facet
```

| File | Tests |
|------|-------|
| `scene_facets_test.go` | 14 |
| `performer_facets_test.go` | 9 |
| `gallery_facets_test.go` | 8 |
| `group_facets_test.go` | 6 |
| `studio_facets_test.go` | 6 |
| `tag_facets_test.go` | 5 |

### Frontend Tests

```bash
cd ui/v2.5
yarn test
```

| File | Tests |
|------|-------|
| `useFacetCounts.test.ts` | 14 |
| `facetCandidateUtils.test.ts` | 18 |
| `GroupsFilter.test.ts` | 8 |

## Known Issues & Solutions

### Labels Showing as IDs
**Cause**: `toMap()` discarded labels
**Fix**: Store both count and label in `LabeledFacetCount`

### Stale Counts Filtering Candidates
**Cause**: Missing loading state check
**Fix**: Added `!facetsLoading` check

### Labels Jumping Between Filters
**Cause**: Full state replacement during lazy load
**Fix**: Partial state update + unique key prefixes

