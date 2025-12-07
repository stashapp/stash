# Fork Patches

This directory contains documentation for merging fork changes after upstream updates.

## Overview

When merging from upstream, most fork code is isolated in:
- **Frontend:** `ui/v2.5/src/extensions/`
- **Backend:** New packages like `pkg/recommendation/`, new files like `*_facets.go`

However, some changes modify upstream files. This directory documents those changes.

## Baseline

This fork is based on **Stash v0.29.3**.

---

## Files

| File | Purpose | Priority |
|------|---------|----------|
| `schema-queries.md` | GraphQL query additions (facets + recommendations) | **High** |
| `repository-interfaces.md` | Faceter interface embeddings | **High** |
| `tag-filter-extensions.md` | Tag filter additions (performers_filter, groups_filter) | Medium |
| `config-extensions.md` | Frontend config changes (recommendations, sidebar) | Medium |
| `utility-additions.md` | Small utility functions | Low |

---

## Quick Merge Guide

### Step 1: Check for Conflicts

```bash
git status
```

### Step 2: Handle Schema Conflicts

If `graphql/schema/schema.graphql` conflicts:
- See `schema-queries.md` for all queries to add

### Step 3: Handle Repository Conflicts

If `pkg/models/repository_*.go` files conflict:
- See `repository-interfaces.md` for Faceter interface lines

### Step 4: Handle Tag Filter Conflicts

If `pkg/sqlite/tag*.go` or `graphql/schema/types/filters.graphql` conflict:
- See `tag-filter-extensions.md` for tag filter additions

### Step 5: Regenerate & Build

```bash
go generate ./...
go build ./...
```

### Step 6: Test

```bash
# Backend tests
go test -v -tags=integration ./pkg/sqlite/... -run Facet

# Frontend
cd ui/v2.5
yarn build
yarn test
```

---

## New Files (Won't Conflict)

These files don't exist upstream - they'll merge cleanly:

### Facets System
```
graphql/schema/types/facets.graphql
pkg/models/facets.go
pkg/sqlite/*_facets.go (6 files + tests)
internal/api/resolver_query_facets.go
internal/api/types_facets.go
```

### Recommendations System
```
pkg/recommendation/performer.go
pkg/recommendation/scene.go
internal/api/resolver_query_performer_recommendations.go
internal/api/resolver_query_scene_recommendations.go
internal/api/resolver_*_recommendations_result_type.go
```

### Frontend Extensions
```
ui/v2.5/src/extensions/           # All custom frontend code (~80 files)
├── lists/                        # List components (6 files)
├── filters/                      # Filter components (29 files)
├── hooks/                        # Custom hooks (7 files)
├── ui/                           # Shared UI components
├── styles/                       # All custom SCSS (~5,700 lines)
│   ├── _list-components.scss
│   ├── _scene-components.scss
│   ├── _player-components.scss
│   ├── _shared-components.scss
│   ├── _gallery-components.scss
│   └── _image-components.scss
└── docs/                         # Documentation
```

---

## Modified Files (May Conflict)

### Backend
| File | Changes | Patch Doc |
|------|---------|-----------|
| `graphql/schema/schema.graphql` | Facet + recommendation queries | `schema-queries.md` |
| `pkg/models/repository_*.go` | Faceter interfaces | `repository-interfaces.md` |
| `graphql/schema/types/filters.graphql` | Tag filter fields | `tag-filter-extensions.md` |
| `pkg/models/tag.go` | TagFilterType fields | `tag-filter-extensions.md` |
| `pkg/sqlite/tag.go` | Join repos + FindFavoriteTagIDs | `tag-filter-extensions.md`, `utility-additions.md` |
| `pkg/sqlite/tag_filter.go` | Filter handlers | `tag-filter-extensions.md` |
| `pkg/models/resolution.go` | ResolutionFromHeight | `utility-additions.md` |
| `pkg/sqlite/sql.go` | Random sort helper | `utility-additions.md` |

### Frontend (Minimal - Most in Extensions)
| File | Changes | Status |
|------|---------|--------|
| `src/core/config.ts` | AI recommendations, sidebar filters | See `config-extensions.md` |
| `src/core/StashService.ts` | Scene recommendations query | See `config-extensions.md` |
| `App.tsx` | ExtensionRegistryProvider wrapper | Keep during merge |
| `src/index.scss` | Extensions import at end | Keep during merge |
| Route files (`Scenes.tsx`, etc.) | Import Enhanced* components | Keep during merge |

> **Note:** All SCSS modifications, filter components, list components, and hooks have been fully extracted to `extensions/` - these files are clean v0.29.3 in upstream locations.
