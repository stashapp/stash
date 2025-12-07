# Fork Migration Plan

This document outlines how custom fork changes are integrated into the extension system.

**Upstream:** `stashapp/stash` (develop branch) | **Baseline:** v0.29.3

## Current State Summary

| Category | Files | Lines Changed | Status |
|----------|-------|---------------|--------|
| **Extensions (Complete)** | ~95 | ~18,500 | ✅ Done |
| **Upstream Modifications** | ~40 | ~5,000 | ✅ Documented |
| **Total** | 143 | +24,675 / -1,122 | |

**Status:** Phase 1-3 complete. All component modifications documented in patch files.

### What's Been Migrated

| Component | Status | Location |
|-----------|--------|----------|
| List components (6) | ✅ | `extensions/lists/` |
| Filter components (29) | ✅ | `extensions/filters/` |
| Custom hooks (7) | ✅ | `extensions/hooks/` |
| UI components (4) | ✅ | `extensions/ui/` |
| SCSS styles (~5,700 lines) | ✅ | `extensions/styles/` |
| Tests (3) | ✅ | `extensions/__tests__/` |
| Scene.tsx | ✅ | `extensions/components/Scene/` (upstream preserved) |
| SceneDetailPanel.tsx | ✅ | `extensions/components/` |
| QueueViewer.tsx | ✅ | `extensions/components/` |
| GalleryDetailPanel.tsx | ✅ | `extensions/components/` |
| ImageDetailPanel.tsx | ✅ | `extensions/components/` |
| AISceneRecommendationRow.tsx | ✅ | `extensions/components/` |
| settings-menu.ts | ✅ | `extensions/player/` |
| GalleryPopover.tsx | ✅ | `extensions/components/` |
| FilterTags | ✅ | `extensions/ui/` (upstream preserved) |

---

## Phase 1: Quick Wins ✅ COMPLETE

### 1.1 Revert Filter Components ✅
**Files:** 31 files in `src/components/List/Filters/`  
**Status:** ✅ Complete - reverted to v0.29.3

### 1.2 Move Hooks to Extensions ✅
**Files:** 7 files in `src/hooks/`  
**Status:** ✅ Complete - moved to `extensions/hooks/`

### 1.3 Fix Import-Only Files ✅
**Files:** ~22 files  
**Status:** ✅ Complete - imports updated

---

## Phase 2: SCSS Extraction ✅ COMPLETE

All custom SCSS is now in `extensions/styles/`. The upstream SCSS files have been reverted to clean v0.29.3. The extension styles load LAST (via `src/index.scss`) and override upstream styles.

| Original File | Extension File | Lines |
|---------------|----------------|-------|
| `List/styles.scss` | `_list-components.scss` | 1,726 |
| `Scenes/styles.scss` | `_scene-components.scss` | 1,305 |
| `ScenePlayer/styles.scss` | `_player-components.scss` | 830 |
| `Shared/styles.scss` | `_shared-components.scss` | 1,086 |
| `Galleries/styles.scss` | `_gallery-components.scss` | 528 |
| `Images/styles.scss` | `_image-components.scss` | 197 |
| **Total** | | **~5,700 lines** |

---

## Phase 3: Component Modifications ✅ COMPLETE

### 3.1 Migrated to Extensions ✅

**Important:** We use the "update imports" pattern, NOT re-exports. Upstream files remain unchanged; importing files are updated to import from extensions.

| Component | Extension Location | Importing File | Status |
|-----------|-------------------|----------------|--------|
| `Scene.tsx` | `extensions/components/Scene/` | `Scenes.tsx` imports from extensions | ✅ |
| `SceneDetailPanel.tsx` | `extensions/components/` | — | ✅ |
| `QueueViewer.tsx` | `extensions/components/` | — | ✅ |
| `GalleryDetailPanel.tsx` | `extensions/components/` | — | ✅ |
| `ImageDetailPanel.tsx` | `extensions/components/` | — | ✅ |
| `AISceneRecommendationRow.tsx` | `extensions/components/` | — | ✅ |
| `settings-menu.ts` | `extensions/player/` | — | ✅ |
| `GalleryPopover.tsx` | `extensions/components/` | `TagLink.tsx` imports from extensions | ✅ |
| `FilterTags` | `extensions/ui/` | 3 List files import from extensions | ✅ |

### 3.2 New Components ✅ COMPLETE

| Component | Status | Location |
|-----------|--------|----------|
| `GalleryPopover.tsx` | ✅ Moved | `extensions/components/GalleryPopover.tsx` |

### 3.3 List Infrastructure Changes ✅ COMPLETE

**Import-only files** - Updated to import from `src/extensions/ui`:

| File | Status |
|------|--------|
| `List/EditFilterDialog.tsx` | ✅ Import updated |
| `List/ItemList.tsx` | ✅ Import updated |
| `List/ListToolbar.tsx` | ✅ Import updated |
| `List/FilterTags.tsx` | ✅ Upstream preserved |

**Small feature changes** - Documented in `patches/list-components.md`:

| File | Change | Lines | Action |
|------|--------|-------|--------|
| `List/ListFilter.tsx` | Added `placeholder` prop | ~10 | Document in patches |
| `List/ListTable.tsx` | Checkbox wrapper for styling | ~15 | Document in patches |
| `List/Pagination.tsx` | Bold formatting in pagination text | ~5 | Document in patches |

### 3.4 Card Components ✅ DOCUMENTED

See `patches/card-components.md`

| File | Change | Lines | Status |
|------|--------|-------|--------|
| `Galleries/GalleryCard.tsx` | Portrait image detection, `titleOnImage` prop | ~15 | ✅ Documented |
| `Groups/GroupCard.tsx` | `titleOnImage` prop passthrough | ~4 | ✅ Documented |

### 3.5 ScenePlayer ✅ DOCUMENTED

See `patches/scene-player.md`

| File | Change | Lines | Status |
|------|--------|-------|--------|
| `ScenePlayer/ScenePlayer.tsx` | Settings-menu integration | +6/-5 | ✅ Documented |
| ~~`settings-menu.ts`~~ | ~~Menu customizations~~ | — | ✅ Migrated to extensions |

### 3.6 SceneDetails ✅ DOCUMENTED

See `patches/scene-details.md`

| File | Change | Status |
|------|--------|--------|
| ~~`Scene.tsx`~~ | ✅ Migrated | `extensions/components/Scene/` |
| ~~`SceneDetailPanel.tsx`~~ | ✅ Migrated | `extensions/components/` |
| ~~`QueueViewer.tsx`~~ | ✅ Migrated | `extensions/components/` |
| `SceneFileInfoPanel.tsx` | Added `id` attribute | ✅ Documented |
| `SceneHistoryPanel.tsx` | Added `id` attribute | ✅ Documented |
| `SceneMarkersPanel.tsx` | Added `id` attribute | ✅ Documented |
| `SceneVideoFilterPanel.tsx` | Added `id` attribute | ✅ Documented |

### 3.7 Other Detail Panels ✅ DOCUMENTED

See `patches/detail-panels.md`

| Group | Files | Status |
|-------|-------|--------|
| `Galleries/GalleryDetails/` | Gallery.tsx, GalleryChaptersPanel.tsx, GalleryFileInfoPanel.tsx, GalleryScenesPanel.tsx | ✅ Documented |
| `Groups/GroupDetails/` | GroupPerformersPanel.tsx, GroupScenesPanel.tsx | ✅ Documented |
| `Images/ImageDetails/` | Image.tsx | ✅ Documented |
| `Performers/PerformerDetails/` | 4 panel files | ✅ Documented |
| `Studios/StudioDetails/` | 5 panel files | ✅ Documented |
| `Tags/TagDetails/` | 5 panel files | ✅ Documented |

### 3.8 Shared Components ✅ DOCUMENTED

See `patches/shared-components.md`

| File | Change | Lines | Status |
|------|--------|-------|--------|
| `ClearableInput.tsx` | Search icon button | +11 | ✅ Documented |
| `CollapseButton.tsx` | `disabled` prop | +2 | ✅ Documented |
| `DetailItem.tsx` | Show more/less + props | +55 | ✅ Documented |
| `GridCard/GridCard.tsx` | Checkbox wrapper + `titleOnImage` | +41 | ✅ Documented |
| `Sidebar.tsx` | `disabled` prop passthrough | +7 | ✅ Documented |
| `TagLink.tsx` | GalleryDetailedLink component | +32 | ✅ Documented |

### 3.9 FrontPage Components ✅ DOCUMENTED

See `patches/frontpage.md`

| File | Change | Lines | Status |
|------|--------|-------|--------|
| `Control.tsx` | AI recommendations rendering | +41 | ✅ Documented |
| `FrontPageConfig.tsx` | AI filter type handling | +9 | ✅ Documented |

### 3.10 Navigation & Settings ✅ DOCUMENTED

See `patches/navigation-settings.md`

| File | Change | Lines | Status |
|------|--------|-------|--------|
| `MainNavbar.tsx` | Removed `userCreatable` from galleries/performers | -2 | ✅ Documented |
| `Scenes/SceneListTable.tsx` | Play button, queue integration | +30 | ✅ Documented |
| `Settings/SettingsInterfacePanel.tsx` | Background image settings | +18 | ✅ Documented |

---

## Phase 4: Core/Config Changes ✅ COMPLETE

### 4.1 Core Files ✅ DOCUMENTED

> Note: Core config changes are backend-specific and documented in BACKEND-API.md

| File | Changes | Status |
|------|---------|--------|
| `core/config.ts` | AI recommendations, sidebar filters, background images | ✅ Backend |
| `core/StashService.ts` | Recommendations query | ✅ Backend |

### 4.2 Models/Types ✅ DOCUMENTED

See `patches/models-types.md`

| File | Changes | Status |
|------|---------|--------|
| `models/list-filter/types.ts` | Added `count` to ILabeledId | ✅ Documented |
| `models/list-filter/criteria/piercings.ts` | New criteria file | ✅ Documented |
| `models/list-filter/criteria/tattoos.ts` | New criteria file | ✅ Documented |
| `models/sceneQueue.ts` | Added IFileObject interface | ✅ Documented |

### 4.3 Utils ✅ DOCUMENTED

See `patches/utils.md`

| File | Changes | Status |
|------|---------|--------|
| `utils/caption.ts` | Refactored language data | ✅ Documented |
| `utils/screen.ts` | Added isSmallScreen function | ✅ Documented |

---

## Phase 5: GraphQL Changes ✅ COMPLETE

See `patches/graphql-frontend.md`

### 5.1 New Files

| File | Status |
|------|--------|
| `graphql/data/facets.graphql` | ✅ Documented (177 lines) |

### 5.2 Modified Files

| File | Changes | Status |
|------|---------|--------|
| `graphql/data/gallery.graphql` | +image_count field | ✅ Documented |
| `graphql/data/group-slim.graphql` | +counts, FilterGroupData | ✅ Documented |
| `graphql/data/performer-slim.graphql` | +counts, FilterPerformerData | ✅ Documented |
| `graphql/data/studio.graphql` | +counts, FilterStudioData | ✅ Documented |
| `graphql/data/tag.graphql` | +counts, FilterTagData | ✅ Documented |
| `graphql/queries/movie.graphql` | FindGroupsForFilter | ✅ Documented |
| `graphql/queries/performer.graphql` | FindPerformersForFilter | ✅ Documented |
| `graphql/queries/scene.graphql` | SceneRecommendations queries | ✅ Documented |
| `graphql/queries/studio.graphql` | FindStudiosForFilter | ✅ Documented |
| `graphql/queries/tag.graphql` | FindTagsForFilter | ✅ Documented |

---

## Phase 6: Miscellaneous ✅ COMPLETE

See `patches/miscellaneous.md`

### 6.1 Localization ✅ DOCUMENTED

| File | Changes | Status |
|------|---------|--------|
| `locales/en-GB.json` | +55 translation strings | ✅ Documented |

### 6.2 Public Assets

| File | Description | Status |
|------|-------------|--------|
| `public/apple-touch-icon.png` | Custom branding | Optional |
| `public/favicon.ico` | Custom branding | Optional |
| `public/favicon.png` | Custom branding | Optional |
| `public/plexhub_icon.png` | Custom icon | Optional |

### 6.3 Config Files

| File | Status |
|------|--------|
| `package.json` | No custom deps |
| `vitest.config.ts` | Standard config |

---

## Migration Checklist

### Phase 1: Quick Wins ✅ COMPLETE
- [x] Revert 31 filter files to v0.29.3
- [x] Delete fork-created filter files
- [x] Move 7 hook files to extensions
- [x] Update hook imports throughout codebase
- [x] Fix ~22 import-only files
- [x] Verify build passes
- [x] Verify tests pass

### Phase 2: SCSS Extraction ✅ COMPLETE
- [x] Extract all component SCSS to extensions/styles/
- [x] Revert upstream SCSS files to v0.29.3
- [x] Verify build passes
- [x] Verify tests pass

### Phase 3: Component Migrations ✅ COMPLETE
- [x] Move `AISceneRecommendationRow.tsx` to extensions
- [x] Move `settings-menu.ts` to extensions
- [x] Move `QueueViewer.tsx` to extensions
- [x] Move `SceneDetailPanel.tsx` to extensions
- [x] Move `GalleryDetailPanel.tsx` to extensions
- [x] Move `ImageDetailPanel.tsx` to extensions
- [x] Move `Scene.tsx` to extensions (correct approach: update Scenes.tsx import)
- [x] Move `GalleryPopover.tsx` to extensions
- [x] Handle import-only List files (EditFilterDialog, ItemList, ListToolbar)
- [x] Create `patches/list-components.md`
- [x] Create `patches/card-components.md`
- [x] Create `patches/navigation-settings.md`
- [x] Create `patches/scene-player.md`
- [x] Create `patches/scene-details.md`
- [x] Create `patches/detail-panels.md`
- [x] Create `patches/shared-components.md`
- [x] Create `patches/frontpage.md`

### Phase 4: Core/Config Patches ✅ COMPLETE
- [x] Core config changes documented (see BACKEND-API.md)
- [x] Create `patches/models-types.md`
- [x] Create `patches/utils.md`

### Phase 5: GraphQL Patches ✅ COMPLETE
- [x] Create `patches/graphql-frontend.md`

### Phase 6: Miscellaneous ✅ COMPLETE
- [x] Create `patches/miscellaneous.md`
- [x] Document localization additions
- [x] Document asset changes (optional branding)
- [x] Document config file changes

---

## Remaining Modified Files Summary

After Phase 3 completion, **~40 component files** still show as modified vs v0.29.3:

| Category | Count | Status |
|----------|-------|--------|
| Import-only changes | 6 | ✅ Imports updated to extensions |
| Detail panels | 22 | ✅ Documented in `patches/detail-panels.md` |
| List infrastructure | 3 | ✅ Documented in `patches/list-components.md` |
| Card components | 2 | ✅ Documented in `patches/card-components.md` |
| Shared components | 6 | ✅ Documented in `patches/shared-components.md` |
| Scene panels | 4 | ✅ Documented in `patches/scene-details.md` |
| ScenePlayer | 1 | ✅ Documented in `patches/scene-player.md` |
| FrontPage | 2 | ✅ Documented in `patches/frontpage.md` |
| Navigation/Settings | 3 | ✅ Documented in `patches/navigation-settings.md` |

**All component modifications are now documented.**

---

## Current Outcome (All Phases Complete) ✅

| Metric | Before | Current |
|--------|--------|---------|
| Modified upstream files | ~90 | ~40 |
| Files in extensions/ | ~50 | ~95 |
| Documented patches | 1 | **12** |
| Merge conflict risk | High | **Low** |

### Upgrade Strategy

**📚 See `UPGRADE-GUIDE.md` for detailed upgrade instructions and checklist.**

Quick steps:
1. `git fetch upstream && git merge upstream/develop`
2. Check which patched files changed (use guide's checklist)
3. Re-apply patches from `extensions/patches/`
4. `yarn generate` (if GraphQL changed)
5. `yarn build && yarn test`

### Patch Categories

| Category | Patch File | Priority |
|----------|------------|----------|
| Components | card-components, detail-panels, frontpage, list-components, navigation-settings, scene-details, scene-player, shared-components | High |
| Core | config-extensions, models-types, utils | Medium |
| GraphQL | graphql-frontend | Medium |
| Misc | miscellaneous | Low |

---

## Notes

- **Phases 1-3 Complete:** All component modifications documented
- **Pattern Used:** "Update imports" pattern (NOT re-exports) - upstream files preserved
- **Testing:** Build and tests pass after each migration step (52 tests)
- **All Phases Complete:** Patches documented in `extensions/patches/`

### Files in `extensions/components/`

```
extensions/components/
├── Scene/
│   └── Scene.tsx           # Full scene page (Scenes.tsx imports from here)
├── AISceneRecommendationRow.tsx
├── GalleryDetailPanel.tsx
├── GalleryPopover.tsx      # Moved from Galleries/ (TagLink.tsx imports from here)
├── ImageDetailPanel.tsx
├── QueueViewer.tsx
└── SceneDetailPanel.tsx
```

### Patch Documentation in `extensions/patches/`

```
extensions/patches/
├── card-components.md      # GalleryCard, GroupCard changes
├── detail-panels.md        # Gallery, Image, other detail panels
├── frontpage.md            # Control.tsx, FrontPageConfig.tsx
├── graphql-frontend.md     # GraphQL schema and queries
├── list-components.md      # ListFilter, ListTable, Pagination changes
├── miscellaneous.md        # Localization, assets, config
├── models-types.md         # Models and type definitions
├── navigation-settings.md  # MainNavbar, SceneListTable, SettingsInterfacePanel
├── scene-details.md        # Scene panel id attributes
├── scene-player.md         # ScenePlayer settings-menu integration
├── shared-components.md    # ClearableInput, CollapseButton, DetailItem, etc.
└── utils.md                # Utility function changes
```
