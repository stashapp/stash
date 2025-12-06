# Sidebar Filters & Facets System - Implementation Changelog

## Summary

This document summarizes all sidebar filter improvements and the facets system implementation.

## Filter Component Improvements

### Boolean Filter Enhancements (`BooleanFilter.tsx`)

**Features Added:**
- Context-specific labels (e.g., "Organized" / "Unorganized" instead of "True" / "False")
- Icons for each boolean type (checkmarks, hearts, etc.)
- Dynamic counts from facets
- Show zero counts as dimmed (0.4 opacity)

**Supported Types:**
- `organized` - Organized / Unorganized
- `interactive` - Interactive / Non-Interactive
- `filter_favorites` - Favorite / Not Favorite (heart icons)
- `has_markers` - Has Markers / No Markers
- `duplicated` - Duplicated / Not Duplicated
- `ignore_auto_tag` - Ignored / Not Ignored
- `has_chapters` - Has Chapters / No Chapters

### Number Filter Enhancements (`NumberFilter.tsx`)

**Features Added:**
- Quick preset values (clickable buttons)
- "Custom..." option for manual input
- Cancel button to exit custom mode
- Filter-specific presets (scene_count, performer_count, etc.)

**Preset Ranges:**
- Scene/Gallery count: 0 to 1000
- Image count: 0 to 100,000
- Performer/Tag count: 0 to 50
- Height: 150cm to 200cm
- Weight: 45kg to 100kg

### String Filter Enhancements (`StringFilter.tsx`)

**Features Added:**
- Multi-value support (OR logic)
- Comma or newline separated values
- Automatic conversion to regex pattern
- Tags-like behavior for string filters

**Example:**
```
Input: "value1, value2, value3"
Result: Regex "value1|value2|value3"
```

### Date Filter Enhancements (`DateFilter.tsx`)

**Features Added:**
- Quick preset options (Today, Yesterday, Last 7/30/90 days, Last year)
- "Custom..." option for date picker
- Cancel button to exit custom mode

### Duration Filter Enhancements (`SidebarDurationFilter.tsx`)

**Features Added:**
- Quick presets (Under 5 min, 5-15 min, etc.)
- "Custom..." option for manual input
- Cancel button to exit custom mode

### Age Filter Enhancements (`SidebarAgeFilter.tsx`)

**Features Added:**
- Age range presets (18-25, 25-35, etc.)
- Cancel button for custom input

### PHash Filter Enhancements (`PhashFilter.tsx`)

**Features Added:**
- Distance presets (Exact Match, Very Similar, Similar, etc.)
- "Custom..." option for manual distance
- Cancel button to exit custom mode
- Clear input explanation

### Stash ID Filter Enhancements (`StashIDFilter.tsx`)

**Features Added:**
- URL paste auto-detection
- Automatic parsing of stash-box URLs
- Extracts endpoint and stash ID from URL

**Supported URL Patterns:**
- `https://stashdb.org/scenes/{id}`
- `https://stashdb.org/performers/{id}`
- Similar patterns for other stash-box instances

### Path Filter Enhancements (`PathFilter.tsx`)

**Features Added:**
- Common path presets
- "Browse..." option for path input
- Cancel button to exit browse mode

### Captions Filter Enhancements (`CaptionsFilter.tsx`)

**Features Added:**
- Text-based country codes (e.g., `[EN]`, `[DE]`)
- Alphabetical sorting of language list
- Dynamic counts from facets
- Removed excluded-list modifier styling

### Rating Filter (`RatingFilter.tsx`)

**Features Added:**
- Star rating display
- Dynamic counts from facets
- Click to filter by rating

### Circumcised Filter (`CircumcisedFilter.tsx`)

**Features Added:**
- New dedicated component (was using StringFilter incorrectly)
- Cut / Uncut options with icons
- Dynamic counts from facets

### Groups Filter (`GroupsFilter.tsx`)

**Features Added:**
- New component for groups filtering
- Support for hierarchical groups (containing_groups, sub_groups)
- Dynamic counts from facets
- Minimal GraphQL fragment for fast loading

### Design Consistency Updates

**Icons Added:**
- Resolution filter: expand icon
- Orientation filter: mobile icon
- Circumcised filter: scissors icon

**Opacity Standardized:**
- Base candidate opacity: 0.7
- Zero count (dimmed): 0.4
- Selected items: 1.0

**Loading Indicators:**
- Spinner for candidate loading
- Dash (`—`) for count loading
- Consistent across all filters

### Filter Customization

**Features Added:**
- Hide chevron-right SVG during customization mode
- Collapsible filter sections
- Visibility toggle for each filter
- Per-view filter preferences

## Facets System Implementation

### 1. Backend Facets Endpoint

**GraphQL Queries Added:**
- `sceneFacets` - with lazy loading support for performer_tags and captions
- `performerFacets`
- `galleryFacets`
- `groupFacets`
- `studioFacets`
- `tagFacets`

**Files Modified/Created:**
- `graphql/schema/types/facets.graphql` - Type definitions
- `graphql/schema/schema.graphql` - Query endpoints
- `pkg/models/facets.go` - Go model structs
- `pkg/models/repository_*.go` - Repository interfaces
- `pkg/sqlite/*_facets.go` - SQLite implementations
- `internal/api/resolver_query_facets.go` - GraphQL resolvers
- `internal/api/types_facets.go` - API types

### 2. Frontend Integration

**Hooks & Context:**
- `useFacetCounts` hook with debouncing and lazy loading support
- `FacetCountsContext` for sharing counts across filter components
- `LabeledFacetCount` interface for preserving labels

**Files Created:**
- `ui/v2.5/graphql/data/facets.graphql` - GraphQL queries
- `ui/v2.5/src/hooks/useFacetCounts.ts` - Hooks and context

### 3. Filter Components Updated

All sidebar filter components were updated to consume facet counts:

| Component | File |
|-----------|------|
| Tags | `TagsFilter.tsx` |
| Performers | `PerformersFilter.tsx` |
| Studios | `StudiosFilter.tsx` |
| Groups | `GroupsFilter.tsx` (NEW) |
| Performer Tags | `PerformerTagsFilter.tsx` |
| Gender | `GenderFilter.tsx` |
| Resolution | `ResolutionFilter.tsx` |
| Orientation | `OrientationFilter.tsx` |
| Captions | `CaptionsFilter.tsx` |
| Rating | `RatingFilter.tsx` |
| Circumcised | `CircumcisedFilter.tsx` |
| Boolean | `BooleanFilter.tsx` |

### 4. List Pages Updated

| Page | Facets Hook | Groups Filter Added |
|------|-------------|---------------------|
| Scenes | `useSceneFacetCounts` | ✓ |
| Performers | `usePerformerFacetCounts` | ✓ (fixed) |
| Galleries | `useGalleryFacetCounts` | - |
| Groups | `useGroupFacetCounts` | ✓ (containing/sub) |
| Studios | `useStudioFacetCounts` | - |
| Tags | `useTagFacetCounts` | - |

### 5. Performance Optimizations

- **CTE-based queries** - Base filter executes once
- **Lazy loading** - Expensive facets load on-demand
- **Debouncing** - 300ms delay on filter changes
- **Parallel execution** - Expensive facets run in goroutines

## Bug Fixes

### Fix: Labels Showing as IDs
- **Issue**: Studio/performer filter showed IDs instead of names
- **Cause**: `toMap()` function discarded labels
- **Fix**: Store both count and label in `LabeledFacetCount`

### Fix: Stale Facet Counts Filtering Candidates
- **Issue**: Old counts filtered new search results incorrectly
- **Cause**: Missing check for loading state
- **Fix**: Added `!facetsLoading` check to `hasValidFacets`

### Fix: Search/Facet Results Mismatch
- **Issue**: Items disappeared when counts loaded
- **Cause**: Search returns by relevance, facets by count (different sets)
- **Fix**: Use facet results directly when no search query

### Fix: Gallery Facets SQL Error
- **Issue**: `no such table: galleries_performers`
- **Cause**: Typo in table name
- **Fix**: Changed to `performers_galleries`

### Fix: Groups Filter Using Wrong Component
- **Issue**: Performers page used `SidebarStudiosFilter` for groups
- **Fix**: Created `SidebarGroupsFilter` component

### Fix: Labels Jumping Between Filters During Lazy Loading
- **Issue**: When performer_tags loaded lazily, labels briefly appeared in wrong filters (performers, studios, tags)
- **Cause**: Full state replacement during lazy load triggered React re-renders with inconsistent data
- **Fix**: Lazy load updates now only update the specific facet (performer_tags or captions) while preserving other facets
- **Additional**: Added unique key prefixes (`sectionID`) to prevent React DOM recycling issues

### Fix: Count Loading Indicator for Candidates
- **Issue**: No visual indication that facet counts were loading
- **Fix**: Added `countsLoading` prop to show pulsing dots (`···`) while counts load

## Test Coverage

### Backend Tests (49 facet tests)
- `pkg/sqlite/scene_facets_test.go` - 14 tests
- `pkg/sqlite/performer_facets_test.go` - 9 tests
- `pkg/sqlite/gallery_facets_test.go` - 8 tests
- `pkg/sqlite/group_facets_test.go` - 6 tests
- `pkg/sqlite/studio_facets_test.go` - 6 tests
- `pkg/sqlite/tag_facets_test.go` - 5 tests
- `internal/api/resolver_query_facets_test.go` - 18 tests

### Frontend Tests (40 tests)
- `useFacetCounts.test.ts` - 14 tests (including 4 lazy loading tests)
- `facetCandidateUtils.test.ts` - 18 tests
- `GroupsFilter.test.ts` - 8 tests

## Running Tests

```bash
# Backend integration tests
go test -v -tags=integration ./pkg/sqlite/... -run Facet

# Backend API tests
go test -v ./internal/api/...

# Frontend tests
cd ui/v2.5 && npm run test
```

## Configuration

### Default Limits
- Facets per category: 100
- Debounce delay: 300ms

### Lazy Loading Triggers
- `performer_tags`: Loads when section is expanded
- `captions`: Loads when section is expanded

## Breaking Changes

None. The facets system is additive and backward compatible.

## Complete File List

### New Files Created

**Backend:**
- `graphql/schema/types/facets.graphql` - GraphQL type definitions
- `pkg/models/facets.go` - Go model structs
- `pkg/sqlite/scene_facets.go` - Scene facets implementation
- `pkg/sqlite/performer_facets.go` - Performer facets implementation
- `pkg/sqlite/gallery_facets.go` - Gallery facets implementation
- `pkg/sqlite/group_facets.go` - Group facets implementation
- `pkg/sqlite/studio_facets.go` - Studio facets implementation
- `pkg/sqlite/tag_facets.go` - Tag facets implementation
- `internal/api/resolver_query_facets.go` - GraphQL resolvers
- `internal/api/types_facets.go` - API type definitions

**Backend Tests:**
- `pkg/sqlite/scene_facets_test.go`
- `pkg/sqlite/performer_facets_test.go`
- `pkg/sqlite/gallery_facets_test.go`
- `pkg/sqlite/group_facets_test.go`
- `pkg/sqlite/studio_facets_test.go`
- `pkg/sqlite/tag_facets_test.go`
- `internal/api/resolver_query_facets_test.go`

**Frontend:**
- `ui/v2.5/graphql/data/facets.graphql` - GraphQL queries
- `ui/v2.5/src/hooks/useFacetCounts.ts` - Facet hooks and context
- `ui/v2.5/src/components/List/Filters/GroupsFilter.tsx` - Groups filter
- `ui/v2.5/src/components/List/Filters/CircumcisedFilter.tsx` - Circumcised filter
- `ui/v2.5/src/components/List/Filters/facetCandidateUtils.ts` - Utility functions

**Frontend Tests:**
- `ui/v2.5/src/hooks/useFacetCounts.test.ts`
- `ui/v2.5/src/components/List/Filters/facetCandidateUtils.test.ts`
- `ui/v2.5/src/components/List/Filters/GroupsFilter.test.ts`

**Documentation:**
- `docs/development/facets-system.md` - Main documentation
- `docs/development/facets-changelog.md` - This changelog
- `docs/development/facets-quick-reference.md` - Developer quick reference
- `docs/development/sidebar-filters.md` - Filter components documentation

### Modified Files

**GraphQL Schema:**
- `graphql/schema/schema.graphql` - Added facet queries

**Backend Models:**
- `pkg/models/repository_scene.go` - Added SceneFaceter interface
- `pkg/models/repository_performer.go` - Added PerformerFaceter interface
- `pkg/models/repository_gallery.go` - Added GalleryFaceter interface
- `pkg/models/repository_group.go` - Added GroupFaceter interface
- `pkg/models/repository_studio.go` - Added StudioFaceter interface
- `pkg/models/repository_tag.go` - Added TagFaceter interface

**Frontend Filter Components:**
- `ui/v2.5/src/components/List/Filters/BooleanFilter.tsx` - Context labels, icons
- `ui/v2.5/src/components/List/Filters/NumberFilter.tsx` - Presets, custom input
- `ui/v2.5/src/components/List/Filters/StringFilter.tsx` - Multi-value support
- `ui/v2.5/src/components/List/Filters/DateFilter.tsx` - Presets, cancel button
- `ui/v2.5/src/components/List/Filters/PhashFilter.tsx` - Distance presets
- `ui/v2.5/src/components/List/Filters/StashIDFilter.tsx` - URL paste support
- `ui/v2.5/src/components/List/Filters/PathFilter.tsx` - Presets, browse option
- `ui/v2.5/src/components/List/Filters/CaptionsFilter.tsx` - Country codes, sorting
- `ui/v2.5/src/components/List/Filters/SidebarDurationFilter.tsx` - Presets, cancel
- `ui/v2.5/src/components/List/Filters/SidebarAgeFilter.tsx` - Presets, cancel
- `ui/v2.5/src/components/List/Filters/PerformersFilter.tsx` - Facet integration
- `ui/v2.5/src/components/List/Filters/TagsFilter.tsx` - Facet integration
- `ui/v2.5/src/components/List/Filters/StudiosFilter.tsx` - Facet integration
- `ui/v2.5/src/components/List/Filters/PerformerTagsFilter.tsx` - Facet integration
- `ui/v2.5/src/components/List/Filters/GenderFilter.tsx` - Facet integration
- `ui/v2.5/src/components/List/Filters/ResolutionFilter.tsx` - Facet integration, icon
- `ui/v2.5/src/components/List/Filters/OrientationFilter.tsx` - Facet integration, icon
- `ui/v2.5/src/components/List/Filters/RatingFilter.tsx` - Facet integration
- `ui/v2.5/src/components/List/Filters/SidebarListFilter.tsx` - Loading indicators
- `ui/v2.5/src/components/List/Filters/LabeledIdFilter.tsx` - Groups filter support

**Frontend List Pages:**
- `ui/v2.5/src/components/Scenes/MySceneList.tsx` - Facets, groups filter
- `ui/v2.5/src/components/Performers/MyPerformerList.tsx` - Facets, groups filter fix
- `ui/v2.5/src/components/Galleries/MyGalleryList.tsx` - Facets
- `ui/v2.5/src/components/Groups/MyGroupList.tsx` - Facets, containing/sub groups
- `ui/v2.5/src/components/Studios/MyStudioList.tsx` - Facets
- `ui/v2.5/src/components/Tags/MyTagList.tsx` - Facets

**GraphQL Queries:**
- `ui/v2.5/graphql/queries/movie.graphql` - Added FindGroupsForFilter
- `ui/v2.5/graphql/data/group-slim.graphql` - Added FilterGroupData fragment

**Styles:**
- `ui/v2.5/src/components/List/styles.scss` - Filter styling updates

## Future Improvements

Potential enhancements for large databases:
1. **Background caching** - Pre-compute facets for common filters
2. **Incremental updates** - Update counts on entity changes
3. **Composite indexes** - Add indexes for common facet queries
4. **Query splitting** - Run facet groups in parallel batches

