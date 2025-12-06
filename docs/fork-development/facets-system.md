# Facets System Documentation

This document describes the facets aggregation system implemented in Stash, which provides dynamic filter counts for sidebar filters across all list pages.

## Overview

The facets system enables the UI to show how many results each filter option would return. For example, when viewing the scenes list, the "Performers" filter can show that "Performer A" appears in 50 scenes, "Performer B" in 30 scenes, etc.

### Key Benefits
- **Performance**: Single query returns counts for multiple filter dimensions
- **User Experience**: Users can see which filters have results before clicking
- **Zero-Count Filtering**: Options with zero results are hidden from the filter list

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend                                 │
│  ┌─────────────────┐    ┌──────────────────┐                   │
│  │ List Page       │───▶│ FacetCountsContext│                   │
│  │ (MySceneList)   │    │ Provider          │                   │
│  └─────────────────┘    └────────┬─────────┘                   │
│                                  │                              │
│  ┌─────────────────┐    ┌────────▼─────────┐                   │
│  │ Filter Component│◀───│ useFacetCounts   │                   │
│  │ (SidebarTagsFilter)  │ Hook             │                   │
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

The facets system defines several GraphQL types for different count types:

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

type OrientationFacetCount {
  orientation: OrientationEnum!
  count: Int!
}

type GenderFacetCount {
  gender: GenderEnum!
  count: Int!
}

type RatingFacetCount {
  rating: Int!
  count: Int!
}

type CaptionFacetCount {
  language: String!
  count: Int!
}

type CircumcisedFacetCount {
  value: CircumisedEnum!
  count: Int!
}
```

### Query Endpoints

Each entity type has a facets query:

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

### Repository Interface

Each entity repository implements a `GetFacets` method:

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

### SQLite Implementation

Facets are computed using a single CTE (Common Table Expression) query that executes the base filter once and then computes counts for each dimension:

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

### Lazy Loading for Expensive Facets

Some facets are expensive to compute (multiple joins). These can be loaded on-demand:

| Facet | Joins | Lazy Loaded |
|-------|-------|-------------|
| `performer_tags` | 3 | Yes |
| `captions` | 2 | Yes |
| `tags`, `performers`, `studios` | 2 | No (always loaded) |
| `organized`, `interactive`, `rating` | 1 | No (always loaded) |

```go
// pkg/sqlite/scene_facets.go
func (qb *SceneStore) GetFacets(ctx context.Context, filter *SceneFilterType, limit int, options SceneFacetOptions) (*SceneFacets, error) {
    // Core facets always run
    go qb.getCoreFacets(ctx, baseSQL, baseArgs, limit, result, &mu)
    
    // Expensive facets only if requested
    if options.IncludePerformerTags {
        go qb.getPerformerTagsFacet(ctx, baseSQL, baseArgs, limit, result, &mu)
    }
    if options.IncludeCaptions {
        go qb.getCaptionsFacet(ctx, baseSQL, baseArgs, result, &mu)
    }
}
```

### State Preservation During Lazy Loading

When lazy-loaded facets (performer_tags, captions) are fetched, the frontend uses a **partial state update** to prevent UI glitches:

```typescript
// src/hooks/useFacetCounts.ts

// Detect if this is a lazy-load update
const isLazyLoadUpdate = 
  (includePerformerTags && !lastOptionsRef.current.includePerformerTags) ||
  (includeCaptions && !lastOptionsRef.current.includeCaptions);

if (isLazyLoadUpdate) {
  // Only update the lazy-loaded facets, preserve everything else
  setCounts((prev) => ({
    ...prev,  // Keep existing tags, performers, studios, etc.
    performerTags: includePerformerTags ? toMap(facets.performer_tags) : prev.performerTags,
    captions: includeCaptions ? toCaptionMap(facets.captions) : prev.captions,
  }));
} else {
  // Full update - update all facets
  setCounts((prev) => ({
    tags: toMap(facets.tags),
    performers: toMap(facets.performers),
    // ... all facets updated
  }));
}
```

**Why this matters**: Without partial updates, a full state replacement triggers React re-renders across all filter components simultaneously. During these re-renders, labels can briefly appear in wrong filter sections due to React's concurrent rendering. The partial update ensures only the relevant filter component re-renders.

## Frontend Implementation

### Hooks

#### `useFacetCounts` Hook

The main hook for fetching facet counts:

```typescript
// src/hooks/useFacetCounts.ts

interface UseFacetCountsOptions {
  isOpen: boolean;           // Only fetch when sidebar is open
  debounceMs?: number;       // Debounce filter changes (default: 300ms)
  includePerformerTags?: boolean;  // Lazy load performer tags
  includeCaptions?: boolean;       // Lazy load captions
}

// Usage
const { counts, loading } = useSceneFacetCounts(filter, {
  isOpen: showSidebar,
  debounceMs: 300,
  includePerformerTags: sectionOpen["performer_tags"],
  includeCaptions: sectionOpen["captions"],
});
```

#### `LabeledFacetCount` Interface

Facet counts include both ID and label to avoid additional lookups:

```typescript
interface LabeledFacetCount {
  count: number;
  label: string;
}

interface FacetCounts {
  tags: Map<string, LabeledFacetCount>;
  performers: Map<string, LabeledFacetCount>;
  studios: Map<string, LabeledFacetCount>;
  groups: Map<string, LabeledFacetCount>;
  // ... enum-based facets use Map<EnumType, number>
}
```

### Context Provider

Facet counts are provided via React Context to all filter components:

```tsx
// In list page component
<FacetCountsContext.Provider value={{ counts: facetCounts, loading: facetLoading }}>
  <SidebarPane>
    {/* Filter components can access counts via useContext */}
  </SidebarPane>
</FacetCountsContext.Provider>
```

### Filter Components

Filter components consume facet counts to:
1. Display count badges next to filter options
2. Hide options with zero counts
3. Sort options by count

```tsx
// src/components/List/Filters/PerformersFilter.tsx

export const SidebarPerformersFilter: React.FC<Props> = ({ ... }) => {
  const { counts: facetCounts, loading: facetsLoading } = useContext(FacetCountsContext);
  
  const candidatesWithCounts = useMemo(() => {
    const hasValidFacets = facetCounts.performers.size > 0 && !facetsLoading;
    const hasSearchQuery = state.query && state.query.length > 0;
    
    if (hasValidFacets && !hasSearchQuery) {
      // Use facet results directly as candidates (TOP N by count)
      const facetCandidates: Option[] = [];
      facetCounts.performers.forEach((facet, id) => {
        if (selectedIds.has(id) || facet.count === 0) return;
        facetCandidates.push({
          id,
          label: facet.label,  // Use label from facet
          count: facet.count,
        });
      });
      facetCandidates.sort((a, b) => (b.count ?? 0) - (a.count ?? 0));
      return facetCandidates;
    }
    // ... fallback for search queries
  }, [state.candidates, facetCounts, facetsLoading]);
};
```

### Utility Functions

The `facetCandidateUtils.ts` module provides reusable logic:

```typescript
// For entity filters (performers, tags, studios, groups)
function buildFacetCandidates(options: FacetCandidateOptions): Option[]

// For enum filters (resolution, orientation, gender)
function buildEnumCandidates(
  options: Option[],
  selectedIds: Set<string>,
  counts: Map<string, number> | undefined,
  countsLoading: boolean
): Option[]
```

## Supported Facets by Entity

### Scene Facets
| Facet | Type | Lazy Loaded |
|-------|------|-------------|
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

## Filter Components with Facet Support

| Component | Entity Types | Facet Key |
|-----------|--------------|-----------|
| `SidebarTagsFilter` | All | `tags` |
| `SidebarPerformersFilter` | Scene, Gallery | `performers` |
| `SidebarStudiosFilter` | Scene, Gallery, Performer, Group | `studios` |
| `SidebarGroupsFilter` | Scene, Performer, Group | `groups` |
| `SidebarPerformerTagsFilter` | Scene | `performerTags` |
| `SidebarGenderFilter` | Performer | `genders` |
| `SidebarResolutionFilter` | Scene | `resolutions` |
| `SidebarOrientationFilter` | Scene | `orientations` |
| `SidebarCaptionsFilter` | Scene | `captions` |
| `SidebarRatingFilter` | All | `ratings` |
| `SidebarCircumcisedFilter` | Performer | `circumcised` |
| `SidebarBooleanFilter` | Various | `organized`, `interactive`, `favorite` |

## Performance Considerations

### Large Databases (>100k items)

For very large databases, facet queries can be slow. Optimizations include:

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

# Run API resolver tests
go test -v ./internal/api/...
```

### Frontend Tests

```bash
cd ui/v2.5

# Run all tests
npm run test

# Run specific test file
npm run test -- facetCandidateUtils.test.ts
```

### Test Coverage

| Area | Tests | Coverage |
|------|-------|----------|
| Scene Facets | 14 | All facet types, filters, lazy loading |
| Performer Facets | 9 | Tags, genders, studios, countries, filters |
| Gallery Facets | 8 | Tags, performers, studios, organized, ratings |
| Group Facets | 6 | Tags, performers, studios, filters |
| Studio Facets | 6 | Tags, parents, favorite, filters |
| Tag Facets | 5 | Parents, children, favorite |
| API Resolver | 18 | All conversion functions |
| Frontend Utils | 26 | Candidate building, regression tests |
| GroupsFilter | 8 | Candidate building, label preservation |

## Known Issues & Solutions

### Issue: Labels Showing as IDs

**Cause**: The `toMap` function was discarding the `label` field from facet results.

**Solution**: Updated `toMap` to store both count and label:
```typescript
function toMap(counts): Map<string, LabeledFacetCount> {
  return new Map(counts.map((c) => [c.id, { count: c.count, label: c.label }]));
}
```

### Issue: Stale Facet Counts Filtering Candidates

**Cause**: When filter changes, old facet counts were used to filter new search results.

**Solution**: Check `facetsLoading` state before applying count-based filtering:
```typescript
const hasValidFacets = facetCounts.size > 0 && !facetsLoading;
```

### Issue: Mismatch Between Search and Facet Results

**Cause**: Search returns items by relevance, facets return TOP N by count. Merging these caused valid items to be filtered out.

**Solution**: When no search query, use facet results directly as candidates:
```typescript
if (hasValidFacets && !hasSearchQuery) {
  // Use facet results directly (TOP N by count)
  facetCounts.forEach((facet, id) => { ... });
} else {
  // Merge search results with facet counts
}
```

### Issue: Labels Jumping Between Filters During Lazy Loading

**Cause**: When performer_tags or captions were lazily loaded, a full state replacement triggered React re-renders across all filter components. During these concurrent re-renders, labels briefly appeared in wrong filter sections.

**Solution**: Lazy load updates now only modify the specific facet field while preserving other facets:
```typescript
if (isLazyLoadUpdate) {
  // Partial update - only update lazy-loaded facets
  setCounts((prev) => ({
    ...prev,  // Preserve tags, performers, studios, etc.
    performerTags: includePerformerTags ? toMap(facets.performer_tags) : prev.performerTags,
    captions: includeCaptions ? toCaptionMap(facets.captions) : prev.captions,
  }));
}
```

**Additional fix**: Added unique key prefixes to `CandidateList` and `SelectedList` items using `sectionID` to prevent React DOM node recycling across filter sections.

## Adding Facets to a New Entity

1. **Define GraphQL types** in `graphql/schema/types/facets.graphql`
2. **Add query** to `graphql/schema/schema.graphql`
3. **Implement repository method** in `pkg/sqlite/{entity}_facets.go`
4. **Add resolver** in `internal/api/resolver_query_facets.go`
5. **Create frontend hook** in `src/hooks/useFacetCounts.ts`
6. **Update filter components** to use `FacetCountsContext`
7. **Add tests** for backend and frontend

## File Reference

### Backend
- `graphql/schema/types/facets.graphql` - GraphQL type definitions
- `graphql/schema/schema.graphql` - Query definitions
- `pkg/models/facets.go` - Go model structs
- `pkg/sqlite/*_facets.go` - SQLite implementations
- `internal/api/resolver_query_facets.go` - GraphQL resolvers
- `internal/api/types_facets.go` - API type definitions

### Frontend
- `ui/v2.5/graphql/data/facets.graphql` - GraphQL fragments and queries
- `ui/v2.5/src/hooks/useFacetCounts.ts` - React hooks and context
- `ui/v2.5/src/components/List/Filters/*Filter.tsx` - Filter components
- `ui/v2.5/src/components/List/Filters/facetCandidateUtils.ts` - Utility functions

### Tests
- `pkg/sqlite/*_facets_test.go` - Backend integration tests
- `internal/api/resolver_query_facets_test.go` - API resolver tests
- `ui/v2.5/src/hooks/useFacetCounts.test.ts` - Hook tests
- `ui/v2.5/src/components/List/Filters/facetCandidateUtils.test.ts` - Utility tests
- `ui/v2.5/src/components/List/Filters/GroupsFilter.test.ts` - Component tests

