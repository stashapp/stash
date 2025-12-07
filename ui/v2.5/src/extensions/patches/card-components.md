# Card Components Patches

This document describes the modifications to upstream Card components that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `GalleryCard.tsx` | Portrait image detection + `titleOnImage` prop | +15 |
| `GroupCard.tsx` | `titleOnImage` prop passthrough | +4 |

---

## GalleryCard.tsx

**Purpose:** 
1. Detect portrait vs landscape images and apply appropriate CSS class
2. Support `titleOnImage` prop for alternative card layout

**File:** `src/components/Galleries/GalleryCard.tsx`

### Change 1: Portrait Image Detection

Add state and handler in `GalleryPreview` component:

```diff
@@ -27,6 +27,13 @@ export const GalleryPreview: React.FC<IGalleryPreviewProps> = ({
   gallery,
   onScrubberClick,
 }) => {
+  const [isPortrait, setIsPortrait] = React.useState(false);
+
+  function identifyPortaitImage(e: React.UIEvent<HTMLImageElement>) {
+    const img = e.target as HTMLImageElement;
+    setIsPortrait(img.width < img.height);
+  }
+
   const [imgSrc, setImgSrc] = useState<string | undefined>(
```

Update the img element:

```diff
@@ -36,9 +43,10 @@ export const GalleryPreview: React.FC<IGalleryPreviewProps> = ({
       {!!imgSrc && (
         <img
           loading="lazy"
-          className="gallery-card-image"
+          className={`gallery-card-image ${isPortrait ? "portrait-image" : ""}`}
           alt={gallery.title ?? ""}
           src={imgSrc}
+          onLoad={identifyPortaitImage}
         />
       )}
```

### Change 2: titleOnImage Prop

Add prop to interface:

```diff
@@ -60,6 +68,7 @@ interface IGalleryCardProps {
   selecting?: boolean;
   selected?: boolean | undefined;
   zoomIndex?: number;
+  titleOnImage?: boolean;
   onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
 }
```

Pass prop to GridCard:

```diff
@@ -235,6 +244,7 @@ export const GalleryCard = PatchComponent(
         selected={props.selected}
         selecting={props.selecting}
         onSelectedChanged={props.onSelectedChanged}
+        titleOnImage={props.titleOnImage}
       />
     );
   }
```

---

## GroupCard.tsx

**Purpose:** Support `titleOnImage` prop for alternative card layout.

**File:** `src/components/Groups/GroupCard.tsx`

### Add prop to interface:

```diff
@@ -42,6 +42,7 @@ interface IProps {
   selected?: boolean;
   zoomIndex?: number;
   onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
+  titleOnImage?: boolean;
   fromGroupId?: string;
   onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
 }
```

### Destructure prop:

```diff
@@ -54,6 +55,7 @@ export const GroupCard: React.FC<IProps> = ({
   selected,
   zoomIndex,
   onSelectedChanged,
+  titleOnImage,
   fromGroupId,
   onMove,
 }) => {
```

### Pass to GridCard:

```diff
@@ -170,6 +172,7 @@ export const GroupCard: React.FC<IProps> = ({
       selecting={selecting}
       onSelectedChanged={onSelectedChanged}
       popovers={maybeRenderPopoverButtonGroup()}
+      titleOnImage={titleOnImage}
     />
   );
 };
```

---

## Related Components

These changes work in conjunction with:
- `GridCard.tsx` - Must support the `titleOnImage` prop
- CSS in `extensions/styles/` - Portrait image and titleOnImage styling

## Application Instructions

After upgrading upstream:

1. Check if these files have changed in the new version
2. If unchanged, apply the patches directly
3. If changed, manually merge the changes
4. Verify GridCard still supports `titleOnImage` prop
5. Run `yarn build` and `yarn test` to verify

