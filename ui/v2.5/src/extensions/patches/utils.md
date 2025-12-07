# Utils Patches

This document describes the modifications to utility files that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `caption.ts` | Refactored language data + added helpers | +30/-13 |
| `screen.ts` | Added `isSmallScreen` function | +3 |

---

## caption.ts

**Purpose:** Refactor language data structure and add helper functions for caption handling.

**File:** `src/utils/caption.ts`

### Replace languageMap with languageData array:

```diff
-export const languageMap = new Map<string, string>([
-  ["de", "Deutsche"],
-  ["en", "English"],
-  ["es", "Español"],
-  ["fr", "Français"],
-  ["it", "Italiano"],
-  ["ja", "日本"],
-  ["ko", "한국인"],
-  ["nl", "Holandés"],
-  ["pt", "Português"],
-  ["ru", "Русский"],
-  ["00", "Unknown"], // stash reserved language code
-]);
+// Language data with code abbreviation (sorted alphabetically, Unknown at end)
+export const languageData: { code: string; name: string }[] = [
+  { code: "de", name: "Deutsche" },
+  { code: "en", name: "English" },
+  { code: "es", name: "Español" },
+  { code: "fr", name: "Français" },
+  { code: "nl", name: "Holandés" },
+  { code: "it", name: "Italiano" },
+  { code: "pt", name: "Português" },
+  { code: "ja", name: "日本" },
+  { code: "ko", name: "한국인" },
+  { code: "ru", name: "Русский" },
+  { code: "00", name: "Unknown" },
+].sort((a, b) => {
+  // Keep "Unknown" at the end
+  if (a.code === "00") return 1;
+  if (b.code === "00") return -1;
+  return a.name.localeCompare(b.name);
+});
+
+// Legacy map for backward compatibility
+export const languageMap = new Map<string, string>(
+  languageData.map((lang) => [lang.code, lang.name])
+);
+
+// Get language code by name (for display as badge)
+export const getLanguageCode = (name: string): string => {
+  const lang = languageData.find((l) => l.name === name);
+  return lang?.code.toUpperCase() ?? "";
+};
```

**Benefits:**
- `languageData` array provides both code and name together
- Sorted alphabetically with "Unknown" at end
- `languageMap` preserved for backward compatibility
- `getLanguageCode()` helper for reverse lookup

---

## screen.ts

**Purpose:** Add helper function to detect small screens.

**File:** `src/utils/screen.ts`

```diff
 };

+const isSmallScreen = () => window.matchMedia("(max-width: 1200px)").matches;
+
 const ScreenUtils = {
   isMobile,
   isTouch,
   matchesMediaQuery,
+  isSmallScreen,
 };
```

**Usage:** Used in detail pages to default sidebar collapsed state on smaller screens.

---

## Application Instructions

After upgrading upstream:

1. **caption.ts**:
   - Check if upstream modified the language map
   - If so, merge new languages into `languageData` array
   - Ensure `languageMap` backward compatibility is preserved

2. **screen.ts**:
   - Simple function addition
   - Should apply cleanly unless upstream restructures ScreenUtils

## Related Files

- `extensions/player/settings-menu.ts` - Uses caption utilities
- `components/Galleries/Gallery.tsx` - Uses `isSmallScreen()`
- `extensions/components/Scene/Scene.tsx` - Uses `isSmallScreen()`
