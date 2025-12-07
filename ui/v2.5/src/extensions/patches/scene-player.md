# ScenePlayer Patches

This document describes the modifications to the ScenePlayer component that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `ScenePlayer.tsx` | Settings menu integration | +6/-5 |

---

## ScenePlayer.tsx

**Purpose:** Replace the upstream `source-selector` plugin with our custom `settings-menu` extension plugin.

**File:** `src/components/ScenePlayer/ScenePlayer.tsx`

### Add settings-menu import:

```diff
 import "./live";
 import "./PlaylistButtons";
 import "./source-selector";
+import "src/extensions/player/settings-menu";
 import "./persist-volume";
```

### Change volume panel to inline:

```diff
         controlBar: {
           pictureInPictureToggle: false,
           volumePanel: {
-            inline: false,
+            inline: true,
           },
```

### Replace sourceSelector with settingsMenu in plugins config:

```diff
           markers: {},
-          sourceSelector: {},
+          settingsMenu: {},
           persistVolume: {},
```

### Update method calls (setSources):

```diff
       const { duration } = file;
-      const sourceSelector = player.sourceSelector();
-      sourceSelector.setSources(
+      const settingsMenu = player.settingsMenu();
+      settingsMenu.setSources(
```

### Update method calls (addTextTrack):

```diff
-          sourceSelector.addTextTrack(
+          settingsMenu.addTextTrack(
             {
               src: `${scene.paths.caption}?lang=${lang}&type=${caption.caption_type}`,
               kind: "captions",
```

---

## Related Files

- `extensions/player/settings-menu.ts` - The custom settings menu plugin implementation

## Application Instructions

After upgrading upstream:

1. Check if ScenePlayer.tsx video.js configuration has changed
2. Ensure the settings-menu plugin is still compatible with the video.js version
3. Apply the patches carefully, checking for any API changes
4. Run `yarn build` and test video playback functionality

