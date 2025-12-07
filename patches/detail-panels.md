# Detail Panel Modifications

## Overview

Detail panels across entity types have been modified to use the enhanced list components from extensions and to add collapsible sections.

## Import Changes (Minor)

These files have a single import change to use `FilteredSceneList` etc. from extensions:

```diff
- import { FilteredSceneList } from "src/components/Scenes/SceneList";
+ import { FilteredSceneList } from "src/extensions/facets/enhanced";
```

### Files with Import-Only Changes

| Entity | Files |
|--------|-------|
| **Performers** | `PerformerScenesPanel.tsx`, `PerformerGalleriesPanel.tsx`, `PerformerGroupsPanel.tsx`, `performerAppearsWithPanel.tsx` |
| **Studios** | `StudioScenesPanel.tsx`, `StudioPerformersPanel.tsx`, `StudioGalleriesPanel.tsx`, `StudioGroupsPanel.tsx`, `StudioChildrenPanel.tsx` |
| **Tags** | `TagScenesPanel.tsx`, `TagPerformersPanel.tsx`, `TagGalleriesPanel.tsx`, `TagGroupsPanel.tsx`, `TagStudiosPanel.tsx` |
| **Groups** | `GroupScenesPanel.tsx`, `GroupPerformersPanel.tsx` |
| **Galleries** | `GalleryScenesPanel.tsx`, `GalleryChaptersPanel.tsx` |

**Note:** These import changes enable facet counts in the list views within these panels.

## Major Changes

### Gallery Details

| File | Lines | Changes |
|------|-------|---------|
| `Gallery.tsx` | +90 / -53 | Layout, collapsible header |
| `GalleryDetailPanel.tsx` | +130 / -63 | Collapsible sections, performer cards |
| `GalleryFileInfoPanel.tsx` | +3 / -1 | ID attribute |

**Key Features:**
- Collapsible detail header
- Performer grid with hover effects
- Tag categories display
- Improved responsive layout

### Image Details

| File | Lines | Changes |
|------|-------|---------|
| `Image.tsx` | +100 / -53 | Layout, collapsible header |
| `ImageDetailPanel.tsx` | +150 / -85 | Collapsible sections, performer cards |

**Key Features:**
- Collapsible detail header
- Performer grid display
- Tag categories
- Background image support

## Merge Strategy

### For Import-Only Files

These should merge cleanly or with minimal conflicts:

1. **Accept upstream changes**
2. **Update imports** to use:
   ```typescript
   import { FilteredSceneList } from "src/extensions/facets/enhanced";
   import { FilteredPerformerList } from "src/extensions/facets/enhanced";
   import { FilteredGalleryList } from "src/extensions/facets/enhanced";
   // etc.
   ```

### For Major Changes (Gallery/Image)

1. **Accept upstream first** - Get base functionality
2. **Re-apply collapsible sections** - Use upstream CollapseButton if available
3. **Re-apply layout changes** - Follow the detail header pattern
4. **Test thoroughly** - Verify collapsible behavior and responsive layout

## Related Files

- **Styles:** `extensions/styles/_gallery-components.scss`, `extensions/styles/_image-components.scss`
- **Shared Components:** `CollapseButton.tsx`, `DetailItem.tsx` (see `patches/shared-components.md`)


