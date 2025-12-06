# Facets System

## Overview

The facets system provides dynamic filter counts in list sidebars. When you filter by a tag, studio, or other criteria, the sidebar shows counts for how many items match each additional filter option.

## Architecture

```
src/extensions/
├── hooks/                      # Facet hooks
│   ├── useFacetCounts.ts       # Core hooks (useXxxFacetCounts)
│   ├── useSceneFacets.ts       # Scene-specific batching
│   ├── useBatchedFilterCounts.ts
│   └── useSidebarFilters.ts
│
├── filters/                    # 29 filter components with facet support
│   ├── TagsFilter.tsx
│   ├── StudiosFilter.tsx
│   ├── PerformersFilter.tsx
│   ├── RatingFilter.tsx
│   └── ... (25 more)
│
├── lists/                      # 6 list components
│   ├── PerformerList.tsx       # Uses usePerformerFacetCounts
│   ├── SceneList.tsx           # Uses useSceneFacetCounts
│   ├── GalleryList.tsx         # Uses useGalleryFacetCounts
│   ├── GroupList.tsx           # Uses useGroupFacetCounts
│   ├── StudioList.tsx          # Uses useStudioFacetCounts
│   └── TagList.tsx             # Uses useTagFacetCounts
│
└── facets/                     # Extension registration
    ├── index.ts
    ├── README.md               # This file
    └── enhanced/index.ts       # EnhancedXxxList aliases
```

## How It Works

### 1. List Component Fetches Counts

```typescript
// In extensions/lists/PerformerList.tsx
const { counts, loading } = usePerformerFacetCounts(filter, {
  isOpen: sidebarOpen
});
```

### 2. Counts Passed via Context

```typescript
<FacetCountsContext.Provider value={{ counts, loading }}>
  <Sidebar>
    <SidebarTagsFilter filter={filter} />
  </Sidebar>
</FacetCountsContext.Provider>
```

### 3. Filter Components Consume Counts

```typescript
// In extensions/filters/TagsFilter.tsx
const { counts, loading } = useFacetCountsContext();

// Display count badge
<span className="badge">{counts?.tags?.[tag.id] ?? 0}</span>
```

## Supported Facets

| Entity | Facets Available |
|--------|------------------|
| Scenes | tags, performers, studios, groups, performer_tags*, captions*, resolutions, orientations, organized, interactive, ratings |
| Performers | tags, studios, genders, countries, circumcised, favorite, ratings |
| Galleries | tags, performers, studios, organized, ratings |
| Groups | tags, performers, studios, containing_groups, sub_groups |
| Studios | tags, parents, favorite |
| Tags | parents, children, favorite |

\* = Lazy loaded on-demand when section expands

## Performance Optimizations

### Lazy Loading
- Facets only fetched when sidebar is open
- Expensive facets (performer tags, captions) load on demand

### Batching
- `useBatchedFilterCounts` groups multiple facet requests
- `useSceneFacets` pre-fetches related facets

### Debouncing
- Filter changes are debounced (300ms) before fetching
- Prevents API spam during rapid filter changes

## Usage Example

```typescript
import { 
  PerformerList,
  usePerformerFacetCounts 
} from "src/extensions";

// The list components handle everything automatically
// Just use the enhanced list:
<PerformerList />
```

## Notes

- List components in `extensions/lists/` are complete implementations, not wrappers
- Filter components in `extensions/filters/` are fork-specific enhancements
- All files use absolute imports (`src/...`) for portability

## See Also

- [ARCHITECTURE.md](../docs/ARCHITECTURE.md) - Full architecture guide
- [CHANGELOG.md](../docs/CHANGELOG.md) - What changed from upstream
- [CONTRIBUTING.md](../docs/CONTRIBUTING.md) - How to add features

