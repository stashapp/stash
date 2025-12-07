# Scene Player Modifications

## Overview

The scene player has been enhanced with a custom settings menu plugin that provides:
- Quality/source selection with automatic fallback on errors
- Playback speed control
- Subtitle/CC track selection

## Files

| File | Status | Details |
|------|--------|---------|
| `settings-menu.ts` | ✅ Moved | Now in `extensions/player/settings-menu.ts` |
| `ScenePlayer.tsx` | 📝 Minor patch | 2 small changes |
| `styles.scss` | ✅ Extracted | Now in `extensions/styles/_player-components.scss` |

## Changes to ScenePlayer.tsx

After merging from upstream, apply these changes:

### 1. Import the settings menu plugin

Add this import after the other plugin imports:

```typescript
import "./live";
import "./PlaylistButtons";
import "./source-selector";
import "src/extensions/player/settings-menu";  // ADD THIS LINE
import "./persist-volume";
```

### 2. Volume panel inline mode

Change the volume panel configuration from `inline: false` to `inline: true`:

```typescript
controlBar: {
  pictureInPictureToggle: false,
  volumePanel: {
    inline: true,  // Changed from false
  },
  chaptersButton: false,
},
```

## Why These Changes?

1. **Settings Menu Plugin**: Provides a unified settings menu instead of scattered controls. The plugin auto-registers when imported.

2. **Volume Panel Inline**: Makes the volume slider appear inline with the control bar instead of as a popup, providing better UX.

## Related Styles

Player styling is in `extensions/styles/_player-components.scss` which includes:
- Scrubber styling
- Volume panel styling
- Settings menu styling
- Video.js customizations


