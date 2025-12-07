# Extensions

Fork-specific features isolated from upstream Stash code for easier maintenance.

**Baseline:** Stash v0.29.3 | **Upstream:** `stashapp/stash` develop branch

## Quick Start

```tsx
// Import everything you need from one place
import { 
  PerformerList,           // List components
  SidebarTagsFilter,       // Filter components  
  useSceneFacetCounts,     // Hooks
} from "src/extensions";

// Or from specific modules
import { PerformerList } from "src/extensions/lists";
import { SidebarTagsFilter } from "src/extensions/filters";
import { useSceneFacetCounts } from "src/extensions/hooks";
```

## Directory Structure

```
extensions/
├── index.ts            # Main entry - exports everything
├── registry.tsx        # Extension registration
├── README.md           # This file
│
├── lists/              # 6 enhanced list components
├── filters/            # 29 custom filter components
├── hooks/              # Custom React hooks
├── ui/                 # Reusable UI components
├── facets/             # Facets extension registration
├── styles/             # Custom SCSS
├── __tests__/          # Extension tests (4 files, 52 tests)
│
└── docs/               # Documentation
    ├── ARCHITECTURE.md # Full architecture guide
    ├── CHANGELOG.md    # What changed from upstream
    └── CONTRIBUTING.md # How to add features
```

## What's Included

### Lists (6)
Full list page implementations with facet counts, custom sidebars, and extended features.

| Component | Key Features |
|-----------|--------------|
| `PerformerList` | Random performer (`p r`), facets |
| `SceneList` | Play queue, scene stats |
| `GalleryList` | Facets, custom filters |
| `GroupList` | Hierarchical groups |
| `StudioList` | Tagger integration |
| `TagList` | Merge dialog |

### Filters (29)
13 NEW + 16 enhanced filter components with facet counts, quick presets, and better UX.

### Hooks (6)
- `useFacetCounts` - Main facet counting hook
- `useSceneFacets`, `useGalleryFacets`, `usePerformerFacets` - Batch facets
- `useSidebarFilters` - Sidebar state management
- `useBatchedFilterCounts` - Batched counting
- `useFacetsContext` - Facets React context

### Styles (~5,700 lines)
SCSS files loaded last (can override anything):
- Component styles: `_list-`, `_scene-`, `_player-`, `_shared-`, `_gallery-`, `_image-components.scss`
- Feature styles: `_facets.scss`, `_sidebar.scss`, `_filter-tags.scss`
- Optional theme: `_plex-theme*.scss`

## Documentation

| Document | Purpose |
|----------|---------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Full architecture guide |
| [UPGRADE-GUIDE.md](docs/UPGRADE-GUIDE.md) | **How to upgrade from upstream** ⭐ |
| [MIGRATION-PLAN.md](docs/MIGRATION-PLAN.md) | Migration status & patch list |
| [BACKEND-API.md](docs/BACKEND-API.md) | Custom GraphQL endpoints & backend files |
| [CHANGELOG.md](docs/CHANGELOG.md) | What changed from upstream |
| [CONTRIBUTING.md](docs/CONTRIBUTING.md) | How to add features |
| [FACETS-SYSTEM.md](docs/FACETS-SYSTEM.md) | Facets technical reference |
| [FILTER-COMPONENTS.md](docs/FILTER-COMPONENTS.md) | Filter components reference |
| [facets/README.md](facets/README.md) | Facets overview |

## Key Principles

1. **All fork code in `extensions/`** - Clean separation from upstream
2. **Upstream files untouched** - `components/List/Filters/` is clean v0.29.3
3. **Absolute imports** - `src/components/...` not `../components/...`
4. **Index exports** - Clean imports via `index.ts` files
5. **SCSS loads last** - Can override any upstream style

## Merging from Upstream

Since `src/components/List/Filters/` is clean:
- Upstream can overwrite those files without conflict
- All customizations are in `extensions/filters/`
- Extension lists import from `extensions/filters`, not upstream

## Migration Status ✅ Complete

All fork changes are documented. See [MIGRATION-PLAN.md](docs/MIGRATION-PLAN.md).

| Category | Files | Status |
|----------|-------|--------|
| List components | 6 | ✅ In extensions |
| Filter components | 29 | ✅ In extensions |
| UI components | 4 | ✅ In extensions |
| Hooks | 6 | ✅ In extensions |
| SCSS | 10 | ✅ In extensions (~5,700 lines) |
| Tests | 4 | ✅ In extensions (52 tests) |
| Component modifications | ~40 | ✅ Documented (12 patch files) |
| GraphQL | 11 | ✅ Documented |
| Core config | 2 | ✅ Documented |

**Upgrading?** See [UPGRADE-GUIDE.md](docs/UPGRADE-GUIDE.md) for step-by-step instructions.
