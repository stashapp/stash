# Miscellaneous Patches

This document describes localization, assets, and configuration changes that need to be reapplied after upgrading.

## Overview

| Category | File | Changes |
|----------|------|---------|
| Localization | `locales/en-GB.json` | +55 new translation strings |
| Assets | `public/` | Custom branding icons |
| Config | `package.json`, `vitest.config.ts` | Dependencies and test config |

---

## Localization (en-GB.json)

**File:** `src/locales/en-GB.json`

### New Translation Strings

Add the following new keys:

#### Age Range Presets

```json
"age_ranges": {
  "18_25": "18-25",
  "26_35": "26-35",
  "36_45": "36-45",
  "46_55": "46-55",
  "56_65": "56-65",
  "65_plus": "65+",
  "custom": "Custom age..."
},
```

#### Duration Range Presets

```json
"duration_ranges": {
  "under_5_min": "Under 5 min",
  "5_15_min": "5-15 min",
  "15_30_min": "15-30 min",
  "30_60_min": "30-60 min",
  "1_2_hours": "1-2 hours",
  "over_2_hours": "Over 2 hours",
  "custom": "Custom duration..."
},
```

#### Date Presets

```json
"date_presets": {
  "today": "Today",
  "yesterday": "Yesterday",
  "this_week": "This week",
  "last_7_days": "Last 7 days",
  "this_month": "This month",
  "last_30_days": "Last 30 days",
  "this_year": "This year",
  "last_year": "Last year",
  "custom": "Custom date..."
},
```

#### Criterion Additions

```json
"criterion": {
  "from": "From",
  "to": "To",
  // ... existing keys
},
"criterion_modifier": {
  "is_null_short": "none",
  "not": "NOT",
  "not_null_short": "any",
  // ... existing keys
},
```

#### Boolean Filter Labels

```json
"not_favourite": "Not Favourite",
"not_interactive": "Not Interactive",
"not_organized": "Not Organised",
"has_markers_true": "Has Markers",
"has_markers_false": "No Markers",
"duplicated_phash_true": "Duplicated",
"duplicated_phash_false": "Not Duplicated",
"performer_favorite_true": "Performer Favourite",
"performer_favorite_false": "Performer Not Favourite",
"has_chapters_true": "Has Chapters",
"has_chapters_false": "No Chapters",
"ignore_auto_tag_true": "Ignored",
"ignore_auto_tag_false": "Not Ignored",
```

#### UI Labels

```json
"discover": "Discover",
"customize_filters": "Customize Filters",
"show_all": "Show All",
"reset_to_defaults": "Reset to Defaults",
"similar": "Similar",
"recommended_scenes": "Recommended Scenes ({count})",
```

---

## Public Assets

**Directory:** `public/`

Custom branding files (if customization is desired):

| File | Description |
|------|-------------|
| `apple-touch-icon.png` | iOS home screen icon |
| `favicon.ico` | Browser favicon |
| `favicon.png` | PNG favicon |
| `plexhub_icon.png` | Custom app icon |

**Note:** These are optional branding customizations. Skip if using default stash branding.

---

## Configuration Files

### package.json

Check for any custom dependencies added to `dependencies` or `devDependencies`.

**Common additions:**
- None currently required beyond upstream

### vitest.config.ts

**File:** `vitest.config.ts`

Check if test configuration has custom settings for extension tests.

---

## Application Instructions

After upgrading upstream:

1. **Localization:**
   - Merge new translation keys into `en-GB.json`
   - Keys are additive, should not conflict
   - Check if upstream added similar keys with different values

2. **Assets:**
   - Optional: Replace with custom branding
   - Or keep upstream defaults

3. **Config:**
   - Compare `package.json` for dependency changes
   - Merge any test configuration changes

## Related Files

- `extensions/ui/` - Uses filter customization strings
- `extensions/filters/` - Uses preset strings
- `extensions/components/` - Uses UI label strings

