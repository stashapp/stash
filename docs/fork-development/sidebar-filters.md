# Sidebar Filter Components Documentation

This document describes the sidebar filter components implemented in Stash, their features, and design patterns.

## Overview

The sidebar filter system provides a consistent, user-friendly interface for filtering content across all list pages. Each filter type has been designed with specific UX considerations for its data type.

## Filter Component Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    SidebarFilterSelector                         │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ FilterWrapper (collapsible section)                         ││
│  │  ┌─────────────────────────────────────────────────────────┐││
│  │  │ Filter Component (e.g., SidebarTagsFilter)              │││
│  │  │  ┌─────────────────────────────────────────────────────┐│││
│  │  │  │ SidebarListFilter / SidebarSection                  ││││
│  │  │  │  - Header with title and count                      ││││
│  │  │  │  - Selected items list                              ││││
│  │  │  │  - Search/input                                     ││││
│  │  │  │  - Candidates list with counts                      ││││
│  │  │  └─────────────────────────────────────────────────────┘│││
│  │  └─────────────────────────────────────────────────────────┘││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## Filter Types

### Entity Filters (ID-based)

These filters select entities by ID and display dynamic counts from facets.

| Component | File | Features |
|-----------|------|----------|
| `SidebarTagsFilter` | `TagsFilter.tsx` | Hierarchical, sub-tags support |
| `SidebarPerformersFilter` | `PerformersFilter.tsx` | Search, facet counts |
| `SidebarStudiosFilter` | `StudiosFilter.tsx` | Hierarchical, subsidiary studios |
| `SidebarGroupsFilter` | `GroupsFilter.tsx` | Hierarchical, containing/sub groups |
| `SidebarPerformerTagsFilter` | `PerformerTagsFilter.tsx` | Scene performer tags |
| `SidebarCountryFilter` | `CountryFilter.tsx` | Country selection |

### Enum Filters

These filters select from predefined enum values.

| Component | File | Features |
|-----------|------|----------|
| `SidebarGenderFilter` | `GenderFilter.tsx` | Gender enum with icons |
| `SidebarResolutionFilter` | `ResolutionFilter.tsx` | Video resolution enum with icon |
| `SidebarOrientationFilter` | `OrientationFilter.tsx` | Video orientation with icon |
| `SidebarCircumcisedFilter` | `CircumcisedFilter.tsx` | Cut/Uncut with icon |
| `SidebarCaptionsFilter` | `CaptionsFilter.tsx` | Language codes, sorted alphabetically |

### Boolean Filters

These filters select true/false values with context-specific labels.

| Component | File | Features |
|-----------|------|----------|
| `SidebarBooleanFilter` | `BooleanFilter.tsx` | Context labels, icons, counts |

**Boolean Filter Types:**
- `organized` - "Organized" / "Unorganized"
- `interactive` - "Interactive" / "Non-Interactive"
- `favorite` / `filter_favorites` - "Favorite" / "Not Favorite" (with heart icons)
- `has_markers` - "Has Markers" / "No Markers"
- `duplicated` - "Duplicated" / "Not Duplicated"
- `ignore_auto_tag` - "Ignored" / "Not Ignored"
- `has_chapters` - "Has Chapters" / "No Chapters"

### Number Filters

These filters input numeric values with preset options.

| Component | File | Features |
|-----------|------|----------|
| `SidebarNumberFilter` | `NumberFilter.tsx` | Presets, custom input, cancel |
| `SidebarAgeFilter` | `AgeFilter.tsx` | Age ranges, presets, cancel |

**Number Filter Presets by Type:**

| Filter Type | Presets |
|------------|---------|
| `scene_count` | 0, 1, 5, 10, 25, 50, 100, 250, 500, 1000 |
| `image_count` | 0, 10, 50, 100, 500, 1000, 5000, 10000, 50000, 100000 |
| `gallery_count` | 0, 1, 5, 10, 25, 50, 100, 250, 500, 1000 |
| `performer_count` | 0, 1, 2, 3, 4, 5, 10 |
| `tag_count` | 0, 1, 2, 3, 5, 10, 20, 50 |
| `play_count` | 0, 1, 2, 3, 5, 10, 25, 50, 100 |
| `o_counter` | 0, 1, 2, 3, 5, 10, 25, 50 |
| `height_cm` | 150, 160, 165, 170, 175, 180, 185, 190, 200 |
| `weight` | 45, 50, 55, 60, 65, 70, 75, 80, 90, 100 |
| `age` | 18, 20, 25, 30, 35, 40, 45, 50, 60, 70 |
| `birth_year` | 1950, 1960, 1970, 1980, 1985, 1990, 1995, 2000, 2005 |
| `rating` | 20, 40, 60, 80, 100 (1-5 stars) |
| Default | 0, 1, 5, 10, 25, 50, 100 |

### String Filters

These filters input text values with multi-value support.

| Component | File | Features |
|-----------|------|----------|
| `SidebarStringFilter` | `StringFilter.tsx` | Multi-value OR/AND, regex conversion |

**Multi-Value Feature:**
- Enter multiple values separated by `,` or newlines
- Automatically converted to regex pattern: `value1|value2|value3`
- Supports both OR (default) and AND logic

### Date/Time Filters

These filters input date ranges with preset options.

| Component | File | Features |
|-----------|------|----------|
| `SidebarDateFilter` | `DateFilter.tsx` | Presets, custom picker, cancel |
| `SidebarDurationFilter` | `SidebarDurationFilter.tsx` | Duration presets, custom, cancel |

**Date Filter Presets:**
- Today
- Yesterday
- Last 7 days
- Last 30 days
- Last 90 days
- Last year
- Custom...

**Duration Filter Presets:**
- Under 5 min
- 5-15 min
- 15-30 min
- 30-60 min
- Over 1 hour
- Over 2 hours
- Custom...

### Special Filters

| Component | File | Features |
|-----------|------|----------|
| `SidebarRatingFilter` | `RatingFilter.tsx` | Star rating with counts |
| `SidebarPhashFilter` | `PhashFilter.tsx` | PHash input, distance presets |
| `SidebarStashIDFilter` | `StashIDFilter.tsx` | URL paste auto-parsing |
| `SidebarPathFilter` | `PathFilter.tsx` | Path presets, browse option |
| `SidebarIsMissingFilter` | `IsMissingFilter.tsx` | Missing field selection |

## Feature Details

### 1. Dynamic Facet Counts

Entity and enum filters display counts from facets:

```tsx
// Count badge in candidate list
<span className="facet-count">{count}</span>

// Zero-count items are hidden
if (count === 0) return null;
```

### 2. Loading Indicators

| State | Indicator |
|-------|-----------|
| Candidates loading | Spinner in dropdown |
| Counts loading | Dash (`—`) placeholder |
| Filter active | Count badge |

### 3. Context-Specific Boolean Labels

```tsx
// BooleanFilter.tsx
switch (option.type) {
  case "organized":
    return value 
      ? { label: "Organized", icon: faCheckCircle }
      : { label: "Unorganized", icon: faTimesCircle };
  case "filter_favorites":
    return value
      ? { label: "Favorite", icon: faHeart }
      : { label: "Not Favorite", icon: faHeartBroken };
  // ...
}
```

### 4. Quick Preset Values

Number, date, and duration filters show clickable presets:

```tsx
// NumberFilter.tsx
<div className="filter-presets">
  {presets.map(preset => (
    <Button
      key={preset}
      variant={isSelected ? "primary" : "outline-secondary"}
      onClick={() => selectPreset(preset)}
    >
      {formatPreset(preset)}
    </Button>
  ))}
  <Button onClick={() => setShowCustom(true)}>Custom...</Button>
</div>
```

### 5. Custom Input with Cancel

Filters with presets support custom input with cancel option:

```tsx
// Show custom input
{showCustomInput && (
  <div className="custom-input-wrapper">
    <input type="number" value={customValue} onChange={...} />
    <Button variant="link" onClick={() => setShowCustomInput(false)}>
      Cancel
    </Button>
  </div>
)}
```

### 6. Multi-Value String Input

String filters support multiple values:

```tsx
// StringFilter.tsx
const handleMultiValueChange = (input: string) => {
  // Split by comma or newline
  const values = input.split(/[,\n]+/).map(v => v.trim()).filter(Boolean);
  
  if (values.length > 1) {
    // Convert to regex OR pattern
    const regex = values.join("|");
    setCriterion({ value: regex, modifier: "MATCHES_REGEX" });
  }
};
```

### 7. URL Paste Auto-Parsing (Stash ID)

```tsx
// StashIDFilter.tsx
const handlePaste = (e: ClipboardEvent) => {
  const text = e.clipboardData?.getData("text");
  
  // Detect stash-box URL patterns
  const match = text?.match(/stashdb\.org\/([a-z]+)\/([a-f0-9-]+)/i);
  if (match) {
    const [, type, stashId] = match;
    setStashID(stashId);
    setEndpoint("https://stashdb.org/graphql");
  }
};
```

### 8. PHash Distance Presets

```tsx
// PhashFilter.tsx
const distancePresets = [
  { label: "Exact Match", value: 0 },
  { label: "Very Similar", value: 4 },
  { label: "Similar", value: 8 },
  { label: "Somewhat Similar", value: 12 },
  { label: "Custom...", value: -1 },
];
```

### 9. Captions Language Support

```tsx
// CaptionsFilter.tsx
// Languages displayed with country codes
<span className="language-badge">[{languageCode}]</span>
<span className="language-name">{languageName}</span>

// Sorted alphabetically by name
languages.sort((a, b) => a.name.localeCompare(b.name));
```

### 10. Hierarchical Filters

Tags, studios, and groups support hierarchy:

```tsx
// Include sub-items toggle
<Form.Check
  type="switch"
  label={<FormattedMessage id="sub_tags" />}
  checked={includeSubItems}
  onChange={(e) => setIncludeSubItems(e.target.checked)}
/>
```

## Design Consistency

### Icons

| Filter Type | Icon |
|-------------|------|
| Tags | `faTags` |
| Performers | `faUser` |
| Studios | `faBuilding` |
| Groups | `faFilm` |
| Rating | `faStar` |
| Resolution | `faExpand` |
| Orientation | `faMobileAlt` |
| Circumcised | `faCut` |
| Favorite | `faHeart` / `faHeartBroken` |
| Organized | `faCheckCircle` / `faTimesCircle` |

### Opacity

| Element | Opacity |
|---------|---------|
| Default candidate | 0.7 |
| Zero count (dimmed) | 0.4 |
| Selected item | 1.0 |
| Modifier option (Any/None) | 0.6 |

### Colors

| Element | Color Variable |
|---------|----------------|
| Selected | `--bs-primary` |
| Excluded | `--bs-danger` |
| Count badge | `--bs-secondary` |
| Zero count | `--bs-gray-500` |

## File Reference

### Filter Components
```
ui/v2.5/src/components/List/Filters/
├── BooleanFilter.tsx           # Boolean with context labels
├── CaptionsFilter.tsx          # Language captions
├── CircumcisedFilter.tsx       # Cut/Uncut enum
├── CountryFilter.tsx           # Country selection
├── DateFilter.tsx              # Date with presets
├── GenderFilter.tsx            # Gender enum
├── GroupsFilter.tsx            # Groups with hierarchy
├── IsMissingFilter.tsx         # Missing fields
├── NumberFilter.tsx            # Numbers with presets
├── OrientationFilter.tsx       # Orientation enum
├── PathFilter.tsx              # File paths
├── PerformersFilter.tsx        # Performers
├── PerformerTagsFilter.tsx     # Performer tags
├── PhashFilter.tsx             # PHash with distance
├── RatingFilter.tsx            # Star ratings
├── ResolutionFilter.tsx        # Resolution enum
├── SidebarDurationFilter.tsx   # Duration with presets
├── SidebarAgeFilter.tsx        # Age with presets
├── SidebarListFilter.tsx       # Base list component
├── SidebarFilterSelector.tsx   # Filter visibility manager
├── StashIDFilter.tsx           # Stash ID with URL paste
├── StringFilter.tsx            # Multi-value strings
├── StudiosFilter.tsx           # Studios with hierarchy
├── TagsFilter.tsx              # Tags with hierarchy
├── facetCandidateUtils.ts      # Utility functions
└── LabeledIdFilter.tsx         # ID filter hooks
```

### Shared Components
```
ui/v2.5/src/components/List/
├── MyFilterSidebar.tsx         # Sidebar container
├── MyFilterTags.tsx            # Filter chip display
├── MyListToolbar.tsx           # Toolbar with filters
└── SidebarSection.tsx          # Collapsible section
```

### Hooks
```
ui/v2.5/src/hooks/
├── useFacetCounts.ts           # Facet counts hook
├── useSidebarFilters.ts        # Filter visibility state
└── useSidebarState.ts          # Sidebar open/close state
```

## Adding a New Filter

### 1. Create the component

```tsx
// MyNewFilter.tsx
import { useContext } from "react";
import { FacetCountsContext } from "src/hooks/useFacetCounts";

export const SidebarMyNewFilter: React.FC<Props> = ({ option, filter, setFilter }) => {
  const { counts, loading } = useContext(FacetCountsContext);
  
  // Build candidates with counts
  const candidates = useMemo(() => {
    // ... filter logic
  }, [counts, loading]);
  
  return (
    <SidebarListFilter
      title={<FormattedMessage id="my_filter" />}
      candidates={candidates}
      // ... other props
    />
  );
};
```

### 2. Add to list page

```tsx
// MyListPage.tsx
import { SidebarMyNewFilter } from "../List/Filters/MyNewFilter";
import { MyNewCriterionOption } from "src/models/list-filter/criteria/my-new";

// In filter definitions
const filterDefinitions = [
  { id: "my_filter", messageId: "my_filter", defaultVisible: false },
];

// In render
<FilterWrapper filterId="my_filter">
  <SidebarMyNewFilter
    title={<FormattedMessage id="my_filter" />}
    option={MyNewCriterionOption}
    filter={filter}
    setFilter={setFilter}
    sectionID="my_filter"
  />
</FilterWrapper>
```

### 3. Add translations

```json
// en-GB.json
{
  "my_filter": "My Filter"
}
```

## Testing Filter Components

```typescript
// MyFilter.test.ts
import { describe, it, expect } from "vitest";

describe("MyFilter", () => {
  it("should build candidates from facet counts", () => {
    // Test candidate building logic
  });
  
  it("should exclude zero-count options", () => {
    // Test filtering
  });
  
  it("should preserve labels from facets", () => {
    // Test label preservation
  });
});
```

