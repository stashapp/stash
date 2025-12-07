# Filter Components Reference

This document describes the sidebar filter components, their features, and design patterns.

## Filter Types

### Entity Filters (ID-based)

| Component | Features |
|-----------|----------|
| `SidebarTagsFilter` | Hierarchical, sub-tags support |
| `SidebarPerformersFilter` | Search, facet counts |
| `SidebarStudiosFilter` | Hierarchical, subsidiary studios |
| `SidebarGroupsFilter` | Hierarchical, containing/sub groups |
| `SidebarPerformerTagsFilter` | Scene performer tags |
| `SidebarCountryFilter` | Country selection |

### Enum Filters

| Component | Features |
|-----------|----------|
| `SidebarGenderFilter` | Gender enum with icons |
| `SidebarResolutionFilter` | Video resolution with icon |
| `SidebarOrientationFilter` | Video orientation with icon |
| `SidebarCircumcisedFilter` | Cut/Uncut with icon |
| `SidebarCaptionsFilter` | Language codes, alphabetical |

### Boolean Filters

Context-specific labels instead of "True/False":

| Type | True Label | False Label |
|------|------------|-------------|
| `organized` | Organized | Unorganized |
| `interactive` | Interactive | Non-Interactive |
| `filter_favorites` | Favorite ❤️ | Not Favorite 💔 |
| `has_markers` | Has Markers | No Markers |
| `duplicated` | Duplicated | Not Duplicated |
| `has_chapters` | Has Chapters | No Chapters |

### Number Filters

Quick preset values:

| Filter Type | Presets |
|------------|---------|
| `scene_count` | 0, 1, 5, 10, 25, 50, 100, 250, 500, 1000 |
| `image_count` | 0, 10, 50, 100, 500, 1000, 5000, 10000, 50000, 100000 |
| `performer_count` | 0, 1, 2, 3, 4, 5, 10 |
| `tag_count` | 0, 1, 2, 3, 5, 10, 20, 50 |
| `height_cm` | 150, 160, 165, 170, 175, 180, 185, 190, 200 |
| `weight` | 45, 50, 55, 60, 65, 70, 75, 80, 90, 100 |
| `age` | 18, 20, 25, 30, 35, 40, 45, 50, 60, 70 |

### String Filters

Multi-value support:
```
Input: "value1, value2, value3"
Result: Regex "value1|value2|value3"
```

### Date Filters

Quick presets:
- Today
- Yesterday
- Last 7 days
- Last 30 days
- Last 90 days
- Last year
- Custom...

### Duration Filters

Presets:
- Under 5 min
- 5-15 min
- 15-30 min
- 30-60 min
- Over 1 hour
- Over 2 hours
- Custom...

### Special Filters

| Component | Features |
|-----------|----------|
| `SidebarRatingFilter` | Star rating with counts |
| `SidebarPhashFilter` | PHash input, distance presets |
| `SidebarStashIDFilter` | URL paste auto-parsing |
| `SidebarPathFilter` | Path presets, browse option |
| `SidebarIsMissingFilter` | Missing field selection |

## Feature Details

### Dynamic Facet Counts

```tsx
// Count badge in candidate list
<span className="facet-count">{count}</span>

// Zero-count items are hidden
if (count === 0) return null;
```

### Loading Indicators

| State | Indicator |
|-------|-----------|
| Candidates loading | Spinner |
| Counts loading | Pulsing dots (`···`) |
| Filter active | Count badge |

### Quick Presets with Cancel

```tsx
<div className="filter-presets">
  {presets.map(preset => (
    <Button onClick={() => selectPreset(preset)}>
      {formatPreset(preset)}
    </Button>
  ))}
  <Button onClick={() => setShowCustom(true)}>Custom...</Button>
</div>

{showCustomInput && (
  <div className="custom-input-wrapper">
    <input type="number" value={customValue} />
    <Button variant="link" onClick={() => setShowCustomInput(false)}>
      Cancel
    </Button>
  </div>
)}
```

### URL Auto-Parsing (Stash ID)

```tsx
const handlePaste = (e: ClipboardEvent) => {
  const text = e.clipboardData?.getData("text");
  const match = text?.match(/stashdb\.org\/([a-z]+)\/([a-f0-9-]+)/i);
  if (match) {
    const [, type, stashId] = match;
    setStashID(stashId);
    setEndpoint("https://stashdb.org/graphql");
  }
};
```

### PHash Distance Presets

| Label | Distance |
|-------|----------|
| Exact Match | 0 |
| Very Similar | 4 |
| Similar | 8 |
| Somewhat Similar | 12 |
| Custom... | -1 |

### Hierarchical Filters

```tsx
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
| Modifier option | 0.6 |

## Adding a New Filter

### 1. Create the component

```tsx
// extensions/filters/MyFilter.tsx
import { useFacetCountsContext } from "src/extensions/hooks";

export const SidebarMyFilter: React.FC<Props> = ({ option, filter, setFilter }) => {
  const { counts, loading } = useFacetCountsContext();
  
  const candidates = useMemo(() => {
    // Build candidates with counts
  }, [counts, loading]);
  
  return (
    <SidebarListFilter
      title={<FormattedMessage id="my_filter" />}
      candidates={candidates}
    />
  );
};
```

### 2. Export from index

```typescript
// extensions/filters/index.ts
export { SidebarMyFilter } from "./MyFilter";
```

### 3. Add to list page

```tsx
<FilterWrapper filterId="my_filter">
  <SidebarMyFilter
    option={MyNewCriterionOption}
    filter={filter}
    setFilter={setFilter}
    sectionID="my_filter"
  />
</FilterWrapper>
```

### 4. Add translations

```json
// en-GB.json
{
  "my_filter": "My Filter"
}
```

