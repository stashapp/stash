# Navigation & Settings Patches

This document describes the modifications to navigation and settings components that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `MainNavbar.tsx` | Removed `userCreatable` from galleries/performers | -2 |
| `SceneListTable.tsx` | Added play button with queue integration | +30 |
| `SettingsInterfacePanel.tsx` | Added background image settings | +18 |

---

## MainNavbar.tsx

**Purpose:** Remove the "+" quick-create buttons from galleries and performers in the navbar.

**File:** `src/components/MainNavbar.tsx`

```diff
@@ -132,7 +132,6 @@ const allMenuItems: IMenuItem[] = [
     href: "/galleries",
     icon: faImages,
     hotkey: "g l",
-    userCreatable: true,
   },
   {
     name: "performers",
@@ -140,7 +139,6 @@ const allMenuItems: IMenuItem[] = [
     href: "/performers",
     icon: faUser,
     hotkey: "g p",
-    userCreatable: true,
   },
```

**Note:** This removes the quick-create functionality for these items. Consider if this is still desired behavior.

---

## SceneListTable.tsx

**Purpose:** Add a play button to scene list table rows that integrates with the queue system.

**File:** `src/components/Scenes/SceneListTable.tsx`

### Add useHistory import:

```diff
 import React from "react";
-import { Link } from "react-router-dom";
+import { Link, useHistory } from "react-router-dom";
```

### Add history hook and onPlayClick handler:

```diff
 ) => {
   const intl = useIntl();
+  const history = useHistory();

   const [updateScene] = useSceneUpdate();

+  function onPlayClick(
+    scene: GQL.SlimSceneDataFragment,
+    timestamp: number,
+    index: number
+  ) {
+    const link = props.queue
+      ? props.queue.makeLink(scene.id, {
+          sceneIndex: index,
+          continue: false,
+          start: 1.911,
+        })
+      : `/scenes/${scene.id}?t=${timestamp}`;
+
+    history.push(link);
+  }
```

### Add play button to title cell:

```diff
     return (
-      <Link to={sceneLink} title={title}>
-        <span className="ellips-data">{title}</span>
-      </Link>
+      <>
+        <button onClick={() => onPlayClick(scene, 0, index)}>
+          <svg
+            className="circular-playbutton"
+            viewBox="0 0 560 560"
+            xmlns="http://www.w3.org/2000/svg"
+          >
+            <path d="M216 170l190.5 110L216 390z"></path>
+          </svg>
+        </button>
+        <Link to={sceneLink} title={title}>
+          <span className="ellips-data">{title}</span>
+        </Link>
+      </>
     );
```

**Note:** The play button uses an inline SVG for the play icon. Styling is in `extensions/styles/`.

---

## SettingsInterfacePanel.tsx

**Purpose:** Add settings for enabling background images on gallery, image, and scene detail pages.

**File:** `src/components/Settings/SettingsInterfacePanel/SettingsInterfacePanel.tsx`

### Add gallery and image background settings (around line 600):

```diff
             </div>
             <div />
           </div>
+          <BooleanSetting
+            id="enableGalleryBackgroundImage"
+            headingID="gallery"
+            checked={ui.enableGalleryBackgroundImage ?? undefined}
+            onChange={(v) => saveUI({ enableGalleryBackgroundImage: v })}
+          />
+          <BooleanSetting
+            id="enableImageBackgroundImage"
+            headingID="image"
+            checked={ui.enableImageBackgroundImage ?? undefined}
+            onChange={(v) => saveUI({ enableImageBackgroundImage: v })}
+          />
           <BooleanSetting
             id="enableMovieBackgroundImage"
```

### Add scene background setting (after performer setting):

```diff
             checked={ui.enablePerformerBackgroundImage ?? undefined}
             onChange={(v) => saveUI({ enablePerformerBackgroundImage: v })}
           />
+          <BooleanSetting
+            id="enableSceneBackgroundImage"
+            headingID="scene"
+            checked={ui.enableSceneBackgroundImage ?? undefined}
+            onChange={(v) => saveUI({ enableSceneBackgroundImage: v })}
+          />
           <BooleanSetting
             id="enableStudioBackgroundImage"
```

**Note:** These settings require corresponding changes in:
- `core/config.ts` - UI configuration types
- Locale files - Translation strings for settings labels

---

## Related Files

These patches work in conjunction with:
- `core/config.ts` - See `patches/config-extensions.md`
- `locales/en-GB.json` - Translation strings
- `extensions/styles/` - CSS for play button and other UI elements

## Application Instructions

After upgrading upstream:

1. Check if these files have changed in the new version
2. Apply MainNavbar.tsx patch if still desired
3. Apply SceneListTable.tsx patch - check for breaking changes in queue API
4. Apply SettingsInterfacePanel.tsx patch - ensure BooleanSetting component API is compatible
5. Run `yarn build` and `yarn test` to verify

