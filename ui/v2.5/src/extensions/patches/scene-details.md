# Scene Details Panel Patches

This document describes the modifications to Scene Details panel components that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `SceneFileInfoPanel.tsx` | Added `id` attribute | +1/-1 |
| `SceneHistoryPanel.tsx` | Added `id` attribute | +1/-1 |
| `SceneMarkersPanel.tsx` | Added `id` attribute | +1/-1 |
| `SceneVideoFilterPanel.tsx` | Added `id` attribute | +1/-1 |

**Purpose:** These changes add `id` attributes to wrapper divs for CSS targeting and styling hooks.

---

## SceneFileInfoPanel.tsx

**File:** `src/components/Scenes/SceneDetails/SceneFileInfoPanel.tsx`

```diff
   return (
-    <div>
+    <div id="scene-file-info-panel">
       <dl className="container scene-file-info details-list">
```

---

## SceneHistoryPanel.tsx

**File:** `src/components/Scenes/SceneDetails/SceneHistoryPanel.tsx`

```diff
   return (
-    <div className="scene-history">
+    <div id="scene-history-panel" className="scene-history">
       <ul className={className}>
```

---

## SceneMarkersPanel.tsx

**File:** `src/components/Scenes/SceneDetails/SceneMarkersPanel.tsx`

```diff
   return (
-    <div className="scene-markers-panel">
+    <div id="scene-markers-panel" className="scene-markers-panel">
       <Button onClick={() => onOpenEditor()}>
```

---

## SceneVideoFilterPanel.tsx

**File:** `src/components/Scenes/SceneDetails/SceneVideoFilterPanel.tsx`

```diff
   return (
-    <div className="container scene-video-filter">
+    <div id="scene-video-filter-panel" className="container scene-video-filter">
       <div className="row form-group">
```

---

## Application Instructions

After upgrading upstream:

1. These are simple id attribute additions
2. Should apply cleanly unless the component structure has changed significantly
3. If upstream adds their own ids, evaluate if our ids are still needed
4. The ids are used by CSS in `extensions/styles/` for custom styling

## Related SCSS

The ids are referenced in `extensions/styles/_scene-components.scss` for panel-specific styling.

