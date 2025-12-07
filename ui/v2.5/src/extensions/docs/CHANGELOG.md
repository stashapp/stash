# Extensions Changelog

This document tracks what has been added/modified from the upstream Stash codebase.

---

## Upstream Files Reverted (Latest)

**Baseline:** Stash v0.29.3

The upstream filter files in `src/components/List/Filters/` have been reverted to clean v0.29.3 versions:

### Reverted Files (14)
- `BooleanFilter.tsx`
- `CustomFieldsFilter.tsx`
- `DateFilter.tsx`
- `HierarchicalLabelValueFilter.tsx`
- `LabeledIdFilter.tsx`
- `NumberFilter.tsx`
- `PathFilter.tsx`
- `PerformersFilter.tsx`
- `PhashFilter.tsx`
- `RatingFilter.tsx`
- `SidebarListFilter.tsx`
- `StashIDFilter.tsx`
- `StudiosFilter.tsx`
- `TagsFilter.tsx`

### Moved to Extensions (18)
Fork-created files moved from `components/List/Filters/` to `extensions/filters/`:
- `AgeFilter.tsx`, `CaptionsFilter.tsx`, `CircumcisedFilter.tsx`
- `CountryFilter.tsx`, `GenderFilter.tsx`, `GroupsFilter.tsx`
- `IsMissingFilter.tsx`, `MyFilterSidebar.tsx`, `OrientationFilter.tsx`
- `PerformerTagsFilter.tsx`, `ResolutionFilter.tsx`, `SidebarDurationFilter.tsx`
- `SidebarFilterSelector.tsx`, `StringFilter.tsx`, `facetCandidateUtils.ts`
- Test files moved to `extensions/__tests__/`

### Import Changes
Extension lists (`extensions/lists/`) now import from `extensions/filters/` instead of `components/List/Filters/`.

**Result:** The `src/components/List/Filters/` directory can be cleanly overwritten by upstream merges with no conflicts.

---

## Facets System

A comprehensive facets aggregation system providing dynamic filter counts in sidebar filters.

### Backend Changes

**New GraphQL endpoint** returning aggregated counts for multiple filter dimensions:

- Single query efficiency using CTE (Common Table Expression)
- Lazy loading for expensive facets (performer_tags, captions)
- Parallel execution with goroutines
- 300ms debouncing to prevent API spam

**Supported Entities:**

| Entity | Facets Available |
|--------|------------------|
| Scenes | tags, performers, studios, groups, performer_tags*, captions*, resolutions, orientations, organized, interactive, ratings |
| Performers | tags, studios, genders, countries, circumcised, favorite, ratings |
| Galleries | tags, performers, studios, organized, ratings |
| Groups | tags, performers, studios, containing_groups, sub_groups |
| Studios | tags, parents, favorite |
| Tags | parents, children, favorite |

\* = Lazy loaded on-demand

### Frontend Changes

**New files:**
- `graphql/data/facets.graphql` - GraphQL queries
- `src/hooks/useFacetCounts.ts` - React hooks and context
- `src/extensions/hooks/useFacetCounts.ts` - Extended hooks
- `src/extensions/hooks/useSceneFacets.ts` - Batched facets
- `src/extensions/hooks/useSidebarFilters.ts` - Sidebar state

---

## Filter Components (29 total)

### NEW Filter Types (13)

| Component | Purpose |
|-----------|---------|
| `AgeFilter` | Age range filter with presets |
| `CaptionsFilter` | Caption language filter with country codes |
| `CircumcisedFilter` | Cut/Uncut with icons (was incorrectly using StringFilter) |
| `CountryFilter` | Country selector with flags |
| `GenderFilter` | Gender filter with icons |
| `GroupsFilter` | Hierarchical groups (containing_groups, sub_groups) |
| `IsMissingFilter` | Missing metadata filter |
| `MyFilterSidebar` | Enhanced sidebar wrapper |
| `OrientationFilter` | Video orientation (landscape, portrait, square) |
| `PerformerTagsFilter` | Tags filter for performers specifically |
| `ResolutionFilter` | Quality resolution presets |
| `SidebarFilterSelector` | Filter visibility selector |
| `StringFilter` | Multi-value with OR logic, comma/newline to regex |

### Enhanced Filter Types (16)

| Component | Enhancements |
|-----------|--------------|
| `BooleanFilter` | Context-specific labels, icons, zero-count dimming |
| `DateFilter` | Quick presets (Today, Last 7/30/90 days, Last year), cancel button |
| `DurationFilter` | Range presets, cancel button |
| `LabeledIdFilter` | Better labels, facet counts |
| `NumberFilter` | Quick preset values, "Custom..." with cancel, filter-specific ranges |
| `PathFilter` | Autocomplete, path validation |
| `PerformersFilter` | Facet counts |
| `PhashFilter` | Distance presets (Exact, Very Similar, Similar, etc.) |
| `RatingFilter` | Star preset buttons |
| `SelectableFilter` | Hierarchical support |
| `SidebarDurationFilter` | Range presets |
| `SidebarListFilter` | Multi-select improvements |
| `StashIDFilter` | URL paste auto-detection, automatic parsing |
| `StudiosFilter` | Facet counts, parent studios |
| `TagsFilter` | Facet counts |
| `facetCandidateUtils` | Utilities for facet processing |

---

## List Components (6)

Complete list page implementations extracted from `My*List` components:

| Component | Features |
|-----------|----------|
| `PerformerList` | Custom sidebar, facet counts, random performer (`p r`) |
| `SceneList` | Facets, play queue, scene stats |
| `GalleryList` | Facets, custom filters |
| `GroupList` | Facets, hierarchical groups |
| `StudioList` | Facets, tagger integration |
| `TagList` | Facets, merge dialog |

---

## UI Components

| Component | Purpose |
|-----------|---------|
| `FilterTags` | Visual filter criteria tags |
| `ListToolbar` | Enhanced list toolbar |
| `ListResultsHeader` | Pagination & sort controls |
| `FilterSidebar` | Sidebar header with search |

---

## Custom Styles

Located in `extensions/styles/` (~5,700 lines total):

### Core Styles
| File | Purpose |
|------|---------|
| `_variables.scss` | CSS custom properties |
| `_facets.scss` | Facets feature styles |
| `_sidebar.scss` | Sidebar styles |
| `_filter-tags.scss` | Filter tag styles |

### Component Styles (Extracted from Upstream)
| File | Lines | Source |
|------|-------|--------|
| `_list-components.scss` | 1,726 | `List/styles.scss` |
| `_scene-components.scss` | 1,305 | `Scenes/styles.scss` |
| `_player-components.scss` | 830 | `ScenePlayer/styles.scss` |
| `_shared-components.scss` | 1,086 | `Shared/styles.scss` |
| `_gallery-components.scss` | 528 | `Galleries/styles.scss` |
| `_image-components.scss` | 197 | `Images/styles.scss` |

### Optional Theme
| File | Purpose |
|------|---------|
| `_plex-theme.scss` | Plex-inspired theme (disabled) |
| `_plex-theme-extended.scss` | Extended theme components |
| `_plex-theme-desktop.scss` | Desktop responsive styles |

---

## Loading Indicators

- **Candidate loading**: Spinner while candidates load
- **Count loading**: Pulsing dots (`···`) while facet counts load
- **Unique key prefixes**: Prevents React DOM recycling issues

---

## Bug Fixes

| Issue | Cause | Fix |
|-------|-------|-----|
| Labels showing as IDs | `toMap()` discarded labels | Store both count and label in `LabeledFacetCount` |
| Stale counts filtering candidates | Missing loading state check | Added `!facetsLoading` check |
| Search/facet results mismatch | Different result sets merged incorrectly | Use facet results directly when no search query |
| Gallery facets SQL error | Wrong table name | Changed `galleries_performers` to `performers_galleries` |
| Groups filter wrong component | Used `SidebarStudiosFilter` | Created `SidebarGroupsFilter` |
| Labels jumping between filters | Full state replacement during lazy load | Partial state update preserving other facets |

---

## Performance Optimizations

- **CTE-based queries**: Base filter executes only once
- **Lazy loading**: performer_tags and captions only fetched when needed
- **Parallel execution**: Expensive facets run concurrently
- **Debouncing**: Reduces API calls during rapid filter changes
- **State preservation**: Lazy loads only update relevant facets

---

## Migration History

### Phase 1: List Components
Moved `My*List` files from `src/components/` to `src/extensions/lists/`

### Phase 2: Filter Components
Extracted 29 filter components to `src/extensions/filters/`

### Phase 3: Hooks
Moved facet hooks to `src/extensions/hooks/`

### Phase 4: Documentation
Consolidated docs into `src/extensions/docs/`

