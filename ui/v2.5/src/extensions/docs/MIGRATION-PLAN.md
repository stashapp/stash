# Fork Migration Plan

This document outlines the remaining work to fully integrate custom fork changes into the extension system.

## Current State Summary

| Category | Files | Lines Changed | Status |
|----------|-------|---------------|--------|
| **Extensions (Complete)** | ~50 | ~8,000 | ✅ Done |
| **Upstream Modifications** | ~90 | ~16,000 | 🔄 Needs Migration |
| **Total** | 143 | +24,675 / -1,122 | |

---

## Phase 1: Quick Wins (Est. 1-2 hours)

### 1.1 Revert Filter Components
**Files:** 31 files in `src/components/List/Filters/`  
**Lines:** +8,967 / -236  
**Status:** Already copied to `extensions/filters/` - need to revert upstream

```bash
# Revert to upstream
git checkout v0.29.3 -- ui/v2.5/src/components/List/Filters/

# Delete fork-created files that don't exist in upstream
git rm ui/v2.5/src/components/List/Filters/AgeFilter.tsx
git rm ui/v2.5/src/components/List/Filters/CaptionsFilter.tsx
# ... etc (see list below)
```

**Fork-created files to delete:**
- `AgeFilter.tsx`
- `CaptionsFilter.tsx`
- `CircumcisedFilter.tsx`
- `CountryFilter.tsx`
- `GenderFilter.tsx`
- `GroupsFilter.tsx`
- `IsMissingFilter.tsx`
- `MyFilterSidebar.tsx`
- `OrientationFilter.tsx`
- `PerformerTagsFilter.tsx`
- `ResolutionFilter.tsx`
- `SidebarDurationFilter.tsx`
- `SidebarFilterSelector.tsx`
- `StringFilter.tsx`
- `facetCandidateUtils.ts`
- `*.test.ts` files

### 1.2 Move Hooks to Extensions
**Files:** 7 files in `src/hooks/`  
**Lines:** +2,186  
**Status:** New files - can move entirely

| File | Action |
|------|--------|
| `hooks/facets/index.ts` | Move to `extensions/hooks/facets/` |
| `hooks/useBatchedFilterCounts.ts` | Move to `extensions/hooks/` |
| `hooks/useFacetCounts.ts` | Move to `extensions/hooks/` |
| `hooks/useFacetCounts.test.ts` | Move to `extensions/__tests__/` |
| `hooks/useFacetsContext.tsx` | Move to `extensions/hooks/` |
| `hooks/useSceneFacets.ts` | Move to `extensions/hooks/` |
| `hooks/useSidebarFilters.ts` | Move to `extensions/hooks/` |

After moving, update `extensions/hooks/index.ts` to export directly instead of re-exporting.

### 1.3 Fix Import-Only Files
**Files:** ~22 files  
**Lines:** ~50 total  
**Status:** Just need import path updates

Files that only changed imports from `PerformerList` → `MyPerformerList`:
- `components/Performers/Performers.tsx`
- `components/Scenes/Scenes.tsx`
- `components/Studios/Studios.tsx`
- `components/Tags/Tags.tsx`
- `components/Groups/Groups.tsx`
- `components/Galleries/Galleries.tsx`
- All `*Panel.tsx` files in detail views

**Action:** Update imports to use `extensions/lists/` directly, then revert to upstream.

---

## Phase 2: SCSS Extraction (Est. 2-3 hours)

### 2.1 Component SCSS Changes
**Files:** 6 files  
**Lines:** +1,990 / -253

| File | Lines | Priority |
|------|-------|----------|
| `List/styles.scss` | +1,270 | High - sidebar/filter styling |
| `Scenes/styles.scss` | +453 | High - scene card/detail styling |
| `ScenePlayer/styles.scss` | +226 | Medium - player customizations |
| `Shared/styles.scss` | +141 | Medium - shared component tweaks |
| `Galleries/styles.scss` | +125 | Low - gallery styling |
| `Images/styles.scss` | +28 | Low - image styling |

### 2.2 Migration Strategy

**Option A: Extract All to Extensions (Recommended)**
1. Create category files in `extensions/styles/`:
   - `_list-components.scss` - List/filter styles
   - `_scene-components.scss` - Scene card/detail styles
   - `_player-components.scss` - Player customizations
   - `_shared-components.scss` - Shared tweaks
   - `_gallery-components.scss` - Gallery styles
   - `_image-components.scss` - Image styles

2. Copy diff content to extension files
3. Revert upstream SCSS files to v0.29.3
4. Import extension styles in `extensions/styles/index.scss`

**Option B: Document as Patches**
- Keep modifications in upstream files
- Document exact changes in `patches/scss-modifications.md`
- Higher merge conflict risk

---

## Phase 3: Component Modifications (Est. 4-6 hours)

### 3.1 New Components (Can Move)
| Component | Location | Lines | Action |
|-----------|----------|-------|--------|
| `AISceneRecommendationRow.tsx` | `Scenes/` | ~200 | Move to `extensions/components/` |

### 3.2 Modified Components (Need Patches)

#### ScenePlayer (3 files, +774/-24 lines)
- `ScenePlayer.tsx` - Playback enhancements
- `settings-menu.ts` - Menu customizations
- `styles.scss` - Player styling

**Strategy:** Document in `patches/scene-player.md`

#### SceneDetails (7 files, +849/-274 lines)
- `Scene.tsx` - Layout changes
- `SceneDetailPanel.tsx` - Collapsible sections, galleries
- `SceneFileInfoPanel.tsx` - File info display
- `SceneHistoryPanel.tsx` - Play/O history
- `SceneMarkersPanel.tsx` - Marker UI
- `SceneVideoFilterPanel.tsx` - Video filters
- `QueueViewer.tsx` - Queue management

**Strategy:** Document in `patches/scene-details.md`

#### Other Detail Panels (~15 files)
- `GalleryDetails/*.tsx`
- `ImageDetails/*.tsx`
- `PerformerDetails/*Panel.tsx`
- `StudioDetails/*Panel.tsx`
- `TagDetails/*Panel.tsx`
- `GroupDetails/*Panel.tsx`

**Strategy:** Document in `patches/detail-panels.md`

#### Shared Components (6 files, +141 lines)
- `ClearableInput.tsx`
- `CollapseButton.tsx`
- `DetailItem.tsx`
- `GridCard/GridCard.tsx`
- `Sidebar.tsx`
- `TagLink.tsx`

**Strategy:** Document in `patches/shared-components.md`

### 3.3 FrontPage Components (2 files)
- `Control.tsx` - AI recommendations rendering
- `FrontPageConfig.tsx` - Recommendations config

**Strategy:** Document in `patches/frontpage.md`

---

## Phase 4: Core/Config Changes (Est. 1 hour)

### 4.1 Core Files
| File | Changes | Strategy |
|------|---------|----------|
| `core/config.ts` | AI recommendations, sidebar filters | Document in `patches/config-extensions.md` ✅ |
| `core/StashService.ts` | Recommendations query | Document in `patches/config-extensions.md` ✅ |

### 4.2 Models/Types
| File | Changes | Strategy |
|------|---------|----------|
| `models/list-filter/types.ts` | Filter type extensions | Document in patches |
| `models/list-filter/criteria/piercings.ts` | New criteria | Move to extensions or patch |
| `models/list-filter/criteria/tattoos.ts` | New criteria | Move to extensions or patch |
| `models/sceneQueue.ts` | Queue modifications | Document in patches |

### 4.3 Utils
| File | Changes | Strategy |
|------|---------|----------|
| `utils/caption.ts` | Caption handling | Document in patches |
| `utils/screen.ts` | Screen utilities | Document in patches |

---

## Phase 5: GraphQL Changes (Est. 1 hour)

### 5.1 New Files (No Conflict)
- `graphql/data/facets.graphql` - Facet type definitions

### 5.2 Modified Files (Need Patches)
| File | Changes |
|------|---------|
| `graphql/data/gallery.graphql` | Additional fields |
| `graphql/data/group-slim.graphql` | Additional fields |
| `graphql/data/performer-slim.graphql` | Additional fields |
| `graphql/data/studio.graphql` | Additional fields |
| `graphql/data/tag.graphql` | Additional fields |
| `graphql/queries/movie.graphql` | Recommendation queries |
| `graphql/queries/performer.graphql` | Facet queries |
| `graphql/queries/scene.graphql` | Facet/recommendation queries |
| `graphql/queries/studio.graphql` | Facet queries |
| `graphql/queries/tag.graphql` | Facet queries |

**Strategy:** Document in `patches/graphql-frontend.md`

---

## Phase 6: Miscellaneous (Est. 30 min)

### 6.1 Localization
- `locales/en-GB.json` - New translation strings

### 6.2 Public Assets
- `public/apple-touch-icon.png` - Custom branding
- `public/favicon.ico` - Custom branding
- `public/favicon.png` - Custom branding
- `public/plexhub_icon.png` - Custom icon

### 6.3 Config Files
- `package.json` - Dependencies (if any added)
- `vitest.config.ts` - Test configuration
- `yarn.lock` - Lock file

---

## Migration Checklist

### Phase 1: Quick Wins
- [ ] Revert 31 filter files to v0.29.3
- [ ] Delete fork-created filter files
- [ ] Move 7 hook files to extensions
- [ ] Update hook imports throughout codebase
- [ ] Fix ~22 import-only files
- [ ] Verify build passes
- [ ] Verify tests pass

### Phase 2: SCSS Extraction
- [ ] Create `_list-components.scss`
- [ ] Create `_scene-components.scss`
- [ ] Create `_player-components.scss`
- [ ] Create `_shared-components.scss`
- [ ] Create `_gallery-components.scss`
- [ ] Create `_image-components.scss`
- [ ] Revert upstream SCSS files
- [ ] Update `extensions/styles/index.scss`
- [ ] Verify styling works correctly

### Phase 3: Component Patches
- [ ] Create `patches/scene-player.md`
- [ ] Create `patches/scene-details.md`
- [ ] Create `patches/detail-panels.md`
- [ ] Create `patches/shared-components.md`
- [ ] Create `patches/frontpage.md`
- [ ] Move `AISceneRecommendationRow.tsx` to extensions

### Phase 4: Core/Config Patches
- [ ] Update `patches/config-extensions.md` (done)
- [ ] Create `patches/models-types.md`
- [ ] Create `patches/utils.md`

### Phase 5: GraphQL Patches
- [ ] Create `patches/graphql-frontend.md`

### Phase 6: Miscellaneous
- [ ] Document localization additions
- [ ] Document asset changes
- [ ] Document config file changes

---

## Expected Outcome

After completing all phases:

| Metric | Before | After |
|--------|--------|-------|
| Modified upstream files | ~90 | ~40 |
| Files in extensions/ | ~50 | ~80 |
| Documented patches | 5 | ~15 |
| Merge conflict risk | High | Low |

---

## Notes

- **Priority:** Phases 1-2 provide the biggest reduction in upstream modifications
- **Time estimate:** Total ~10-12 hours of work
- **Testing:** Run build and tests after each phase
- **Backups:** Create a branch before starting each phase

