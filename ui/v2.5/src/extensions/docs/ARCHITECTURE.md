# Extensions Architecture for Fork Maintenance

This document describes how fork-specific features are organized to minimize conflicts with upstream changes.

## Overview

All fork-specific code lives in `ui/v2.5/src/extensions/`. This isolation provides:

1. **Clean upstream merges** - Upstream changes don't touch `extensions/`
2. **Clear boundaries** - Easy to see what's custom vs upstream
3. **Portable features** - Extensions can be shared or disabled
4. **Absolute imports** - No relative path breakage on restructuring

## Upstream Baseline

This fork is based on **Stash v0.29.3**. Files in `src/components/List/Filters/` have been reverted to the clean v0.29.3 versions - all fork modifications now live in `src/extensions/filters/`.

This means:
- **Upstream filters** = Clean v0.29.3 (no fork code)
- **Extension filters** = All customizations (facet counts, new filters, enhanced UX)
- **Future merges** = Filter folder has no conflicts

## Directory Structure

```
ui/v2.5/src/
├── components/              # Upstream components (CLEAN v0.29.3)
│   └── List/Filters/        # Upstream filter components (v0.29.3 baseline)
│                            # These are NOT modified - all changes in extensions/
│
├── hooks/                   # Shared hooks (some fork additions remain here)
│
└── extensions/              # ALL FORK-SPECIFIC CODE
    ├── index.ts             # Main entry point - exports everything
    ├── registry.tsx         # Extension registration system
    │
    ├── lists/               # Enhanced list components (6 files)
    │   ├── index.ts         # Exports all lists
    │   ├── PerformerList.tsx
    │   ├── SceneList.tsx
    │   ├── GalleryList.tsx
    │   ├── GroupList.tsx
    │   ├── StudioList.tsx
    │   └── TagList.tsx
    │
    ├── filters/             # Custom filter components (29 files)
    │   ├── index.ts         # Exports all filters
    │   ├── AgeFilter.tsx
    │   ├── BooleanFilter.tsx
    │   ├── NumberFilter.tsx
    │   ├── StringFilter.tsx
    │   ├── TagsFilter.tsx
    │   └── ... (24 more)
    │
    ├── hooks/               # Custom React hooks
    │   ├── index.ts         # Exports all hooks
    │   ├── useFacetCounts.ts
    │   ├── useSceneFacets.ts
    │   ├── useSidebarFilters.ts
    │   ├── useBatchedFilterCounts.ts
    │   ├── useFocus.ts
    │   └── facets/
    │       └── index.ts     # Facet-specific exports
    │
    ├── facets/              # Facets extension registration
    │   ├── index.ts         # Extension registration
    │   ├── README.md        # Facets documentation
    │   └── enhanced/
    │       └── index.ts     # Enhanced* component aliases
    │
    ├── __tests__/           # Extension tests
    │   ├── useFacetCounts.test.ts
    │   ├── facetCandidateUtils.test.ts
    │   └── GroupsFilter.test.ts
    │
    ├── ui/                  # Reusable UI components
    │   ├── index.ts
    │   ├── FilterTags.tsx
    │   ├── ListToolbar.tsx
    │   ├── ListResultsHeader.tsx
    │   └── FilterSidebar.tsx
    │
    ├── styles/              # Custom SCSS
    │   ├── index.scss
    │   ├── _variables.scss
    │   ├── _facets.scss
    │   ├── _sidebar.scss
    │   ├── _filter-tags.scss
    │   └── _plex-theme.scss (+ extended, desktop)
    │
    └── docs/                # Documentation
        ├── ARCHITECTURE.md  # This file
        ├── CHANGELOG.md     # What changed from upstream
        └── CONTRIBUTING.md  # How to add features
```

## Index Files

Each folder has an `index.ts` that exports its contents. This enables clean imports:

### Purpose of Index Files

| Index File | Purpose |
|------------|---------|
| `extensions/index.ts` | Main entry - one import for everything |
| `extensions/lists/index.ts` | All 6 list components |
| `extensions/filters/index.ts` | All 29 filter components |
| `extensions/hooks/index.ts` | All custom hooks |
| `extensions/ui/index.ts` | Shared UI components |
| `extensions/facets/index.ts` | Extension registration |
| `extensions/facets/enhanced/index.ts` | Backward-compatible aliases |

### Benefits

```typescript
// Without index files - verbose
import { SidebarTagsFilter } from "src/extensions/filters/TagsFilter";
import { useSceneFacetCounts } from "src/extensions/hooks/useFacetCounts";
import { PerformerList } from "src/extensions/lists/PerformerList";

// With index files - clean
import { SidebarTagsFilter, useSceneFacetCounts, PerformerList } from "src/extensions";

// Or from specific modules
import { SidebarTagsFilter } from "src/extensions/filters";
import { useSceneFacetCounts } from "src/extensions/hooks";
```

## Module Details

### Lists (`extensions/lists/`)

Full list page implementations with:
- Custom sidebar with filter sections
- Facet counts integration
- Random navigation
- Export dialogs
- ~800-1300 lines each

**Exports:**
```typescript
export { PerformerList, MyFilteredPerformerList, MyPerformersFilterSidebarSections } from "./PerformerList";
export { FilteredSceneList as SceneList, ScenesFilterSidebarSections } from "./SceneList";
// ... etc
```

### Filters (`extensions/filters/`)

29 custom/enhanced filter components:

| Category | Components |
|----------|------------|
| **NEW** (13) | AgeFilter, CaptionsFilter, CircumcisedFilter, CountryFilter, GenderFilter, GroupsFilter, IsMissingFilter, MyFilterSidebar, OrientationFilter, PerformerTagsFilter, ResolutionFilter, SidebarFilterSelector, StringFilter |
| **Enhanced** (16) | BooleanFilter, DateFilter, DurationFilter, LabeledIdFilter, NumberFilter, PathFilter, PerformersFilter, PhashFilter, RatingFilter, SelectableFilter, SidebarDurationFilter, SidebarListFilter, StashIDFilter, StudiosFilter, TagsFilter, + utilities |

**Key enhancements:**
- Facet counts display
- Quick presets (dates, numbers, ratings)
- Better UX (icons, loading states)
- Context-specific labels

### Hooks (`extensions/hooks/`)

| Hook | Purpose |
|------|---------|
| `useFacetCounts` | Main facet counting hook |
| `useSceneFacetCounts` | Scene-specific counts |
| `usePerformerFacetCounts` | Performer-specific counts |
| `useGalleryFacetCounts` | Gallery-specific counts |
| `useGroupFacetCounts` | Group-specific counts |
| `useStudioFacetCounts` | Studio-specific counts |
| `useTagFacetCounts` | Tag-specific counts |
| `useSceneFacets` | Batch scene facets |
| `useSidebarFilters` | Sidebar state management |
| `useBatchedFilterCounts` | Batched counting |
| `useFocus` | Focus management |

### Styles (`extensions/styles/`)

SCSS files loaded LAST to override upstream:

```scss
// extensions/styles/index.scss
@import "variables";
@import "facets";
@import "sidebar";
@import "filter-tags";

// Optional Plex theme (currently commented out)
// @import "plex-theme";
// @import "plex-theme-extended";
// @import "plex-theme-desktop";
```

**Wired in `src/index.scss`:**
```scss
// ... upstream imports ...
@import "src/extensions/styles/index";  // Added at end
```

## Import Patterns

### Absolute Imports (Required)

All code in `extensions/` uses absolute imports:

```typescript
// ✅ Correct - works from any location
import { SceneCardsGrid } from "src/components/Scenes/SceneCardsGrid";
import { useFindPerformers } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";

// ❌ Wrong - breaks when files move
import { SceneCardsGrid } from "../Scenes/SceneCardsGrid";
```

### Extension Imports

```typescript
// Everything from one import
import { 
  PerformerList, 
  SidebarTagsFilter, 
  useSceneFacetCounts 
} from "src/extensions";

// Or from specific modules
import { PerformerList } from "src/extensions/lists";
import { SidebarTagsFilter } from "src/extensions/filters";
import { useSceneFacetCounts } from "src/extensions/hooks";

// Enhanced exports (backward compatible)
import { EnhancedPerformerList } from "src/extensions/facets/enhanced";
```

## Extension Registry

The registry allows providers to wrap the app:

```typescript
// registry.tsx
export interface Extension {
  id: string;
  name: string;
  version: string;
  enabled: boolean;
  Provider?: React.ComponentType<{ children: ReactNode }>;
  initialize?: () => void;
}

// Register an extension
registerExtension({
  id: "facets",
  name: "Facets System",
  version: "1.0.0",
  enabled: true,
  initialize: () => console.log("[Facets] Initialized"),
});
```

```typescript
// App.tsx
import { ExtensionRegistryProvider } from "./extensions";

function App() {
  return (
    <ExtensionRegistryProvider>
      {/* ... app content ... */}
    </ExtensionRegistryProvider>
  );
}
```

## Backend Dependencies

The extension system requires custom backend GraphQL endpoints. These are **not** part of upstream Stash.

### Required Backend Changes

| Component | Files |
|-----------|-------|
| GraphQL Schema | `graphql/schema/types/facets.graphql`, `schema.graphql` |
| Go Models | `pkg/models/facets.go` |
| SQLite Implementation | `pkg/sqlite/*_facets.go` (6 files) |
| API Resolvers | `internal/api/resolver_query_facets.go`, `types_facets.go` |
| Repository Interfaces | `pkg/models/repository_*.go` (6 files modified) |

### GraphQL Endpoints

| Query | Purpose |
|-------|---------|
| `sceneFacets` | Scene filter counts |
| `performerFacets` | Performer filter counts |
| `galleryFacets` | Gallery filter counts |
| `groupFacets` | Group filter counts |
| `studioFacets` | Studio filter counts |
| `tagFacets` | Tag filter counts |

> 📖 **Full details**: See [BACKEND-API.md](BACKEND-API.md) for complete backend reference.

## Updating from Upstream

### Baseline Tag

This fork is based on **v0.29.3**. When merging from upstream:

```bash
# 1. Fetch upstream changes
git fetch upstream

# 2. Merge (or rebase)
git checkout develop-local
git merge upstream/develop

# 3. Resolve conflicts (see table below)

# 4. Regenerate GraphQL (if schema changed)
go generate ./...

# 5. Test
yarn build
yarn test
go test -v -tags=integration ./pkg/sqlite/... -run Facet
```

### Expected Conflicts

| Location | Conflict? | Resolution |
|----------|-----------|------------|
| `src/components/List/Filters/` | **NONE** | Clean v0.29.3 - upstream can overwrite |
| `src/extensions/` | **NONE** | Upstream doesn't touch this folder |
| `App.tsx` | Yes | Keep our ExtensionRegistryProvider import |
| `src/index.scss` | Yes | Keep our extensions import at end |
| Route files | Yes | Keep our Enhanced* imports |
| `graphql/schema/schema.graphql` | Yes | Re-add our facet queries |
| `pkg/models/repository_*.go` | Yes | Re-add our Faceter interfaces |
| `pkg/sqlite/*_facets.go` | **NONE** | New files - upstream won't have them |

### Key Points

1. **Filter folder is clean** - `src/components/List/Filters/` can be overwritten by upstream without conflict
2. **Extensions are isolated** - All our code is in `src/extensions/` which upstream doesn't modify
3. **Backend additions** - Our GraphQL types and Go implementations are additive, easy to re-add

## Best Practices

### DO ✅

- Keep all fork code in `src/extensions/`
- Use absolute imports (`src/components/...`)
- Export from index files
- Document what each file does
- Keep list components self-contained
- Use CSS variables for theming

### DON'T ❌

- Modify upstream components directly
- Use relative imports in extensions
- Scatter fork code across the codebase
- Create `My*` prefixed files outside extensions
- Use `export *` when it could cause conflicts

## Troubleshooting

### Build fails with "module not found"

- Check that imports use `src/` prefix
- Verify the file exists at the import path
- Check exports in the relevant index.ts

### Duplicate export error

- Check for conflicting `export *` statements
- Use explicit named exports instead

### Styles not applying

- Check that `extensions/styles/index.scss` imports your file
- Verify `src/index.scss` has the extensions import at the end

### Component not rendering

- Check browser console for errors
- Verify all required exports exist
- Run `yarn build` to check for TypeScript errors

---

## Migration Status

Not all fork changes have been migrated to the extension system yet. See **[MIGRATION-PLAN.md](./MIGRATION-PLAN.md)** for:

- Remaining upstream modifications (~90 files)
- SCSS extraction strategy (~2,000 lines)
- Component patch documentation
- Phase-by-phase migration checklist

### Current Migration Progress

| Category | Status |
|----------|--------|
| List components (`lists/`) | ✅ Complete |
| Filter components (`filters/`) | ✅ Complete |
| UI components (`ui/`) | ✅ Complete |
| Hooks (`hooks/`) | 🔄 Partial (re-exports from `src/hooks/`) |
| Styles (`styles/`) | 🔄 Partial (Plex theme, facets) |
| Component patches | 📝 Needs documentation |
| SCSS modifications | 📝 Needs extraction |
