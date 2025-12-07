# Detail Panels Patches

This document describes the modifications to various detail panel components that need to be reapplied after upgrading.

## Overview

### Import Changes (use extension components)

| File | Import Change |
|------|---------------|
| `Galleries/GalleryDetails/Gallery.tsx` | `GalleryDetailPanel` from `src/extensions/components` |
| `Images/ImageDetails/Image.tsx` | `ImageDetailPanel` from `src/extensions/components` |
| `Groups/GroupDetails/GroupPerformersPanel.tsx` | `EnhancedPerformerList` from `src/extensions/facets/enhanced` |
| `Groups/GroupDetails/GroupScenesPanel.tsx` | `FilteredSceneList` from `src/extensions/facets/enhanced` |

### ID Attribute Additions (for CSS styling)

| File | Element |
|------|---------|
| `GalleryChaptersPanel.tsx` | `id="gallery-chapter-panel"` |
| `GalleryFileInfoPanel.tsx` | `id="gallery-file-info-panel"` |
| `GalleryScenesPanel.tsx` | `id="gallery-scenes-panel"` |
| `PerformerGalleriesPanel.tsx` | `id` attribute |
| `PerformerGroupsPanel.tsx` | `id` attribute |
| `PerformerScenesPanel.tsx` | `id` attribute |
| `performerAppearsWithPanel.tsx` | `id` attribute |
| `StudioChildrenPanel.tsx` | `id` attribute |
| `StudioGalleriesPanel.tsx` | `id` attribute |
| `StudioGroupsPanel.tsx` | `id` attribute |
| `StudioPerformersPanel.tsx` | `id` attribute |
| `StudioScenesPanel.tsx` | `id` attribute |
| `TagGalleriesPanel.tsx` | `id` attribute |
| `TagGroupsPanel.tsx` | `id` attribute |
| `TagPerformersPanel.tsx` | `id` attribute |
| `TagScenesPanel.tsx` | `id` attribute |
| `TagStudiosPanel.tsx` | `id` attribute |

### Major Changes

| File | Change | Lines |
|------|--------|-------|
| `Gallery.tsx` | Background image + import + structure | ~145 |
| `Image.tsx` | Background image + import + structure | ~155 |

---

## Gallery.tsx

**File:** `src/components/Galleries/GalleryDetails/Gallery.tsx`

### Import change:

```diff
-import { GalleryDetailPanel } from "./GalleryDetailPanel";
+import { GalleryDetailPanel } from "src/extensions/components";
```

### Add ScreenUtils import:

```diff
+import ScreenUtils from "src/utils/screen";
```

### Add background image config:

```diff
   const { configuration } = useContext(ConfigurationContext);
+  const uiConfig = configuration?.ui;
+  const enableBackgroundImage = uiConfig?.enableGalleryBackgroundImage ?? false;
```

### Change default collapsed state:

```diff
-  const [collapsed, setCollapsed] = useState(false);
+  const [collapsed, setCollapsed] = useState(ScreenUtils.isSmallScreen());
```

### Add background image render function:

```typescript
function maybeRenderHeaderBackgroundImage() {
  if (enableBackgroundImage && gallery != null && gallery.studio != null) {
    let image = gallery.studio.image_path;
    if (image) {
      const imageURL = new URL(image);
      let isDefaultImage = imageURL.searchParams.get("default");
      if (!isDefaultImage) {
        return (
          <div className="background-image-container">
            <picture>
              <source src={image} />
              <img
                className="background-image"
                src={image}
                alt={`${gallery.studio.name} background`}
              />
            </picture>
          </div>
        );
      }
    }
  }
}
```

### Add id to wrapper:

```diff
-    <div className="row">
+    <div id="gallery-page" className="row">
```

---

## Image.tsx

**File:** `src/components/Images/ImageDetails/Image.tsx`

Similar changes to Gallery.tsx:
- Import `ImageDetailPanel` from `src/extensions/components`
- Add `enableImageBackgroundImage` config
- Add `maybeRenderHeaderBackgroundImage()` function
- Add `id="image-page"` to wrapper
- Add `detail-header` and `detail-container` wrapper divs

---

## ID Attribute Changes (Example Pattern)

All panel files follow this pattern:

```diff
   return (
-    <div className="container panel-name">
+    <div id="panel-name-panel" className="container panel-name">
```

---

## Application Instructions

After upgrading upstream:

1. **Import changes** - Update imports to point to extension components
2. **ID attributes** - Apply simple id additions for CSS targeting
3. **Gallery.tsx/Image.tsx** - These have significant structural changes:
   - Check if upstream has added similar background image functionality
   - May need manual merging if component structure changed
4. Verify CSS in `extensions/styles/` still targets the correct elements

## Related Files

- `extensions/components/GalleryDetailPanel.tsx`
- `extensions/components/ImageDetailPanel.tsx`
- `extensions/facets/enhanced/` - Enhanced list components
- `extensions/styles/` - CSS that targets these panel ids
- `utils/screen.ts` - ScreenUtils helper

