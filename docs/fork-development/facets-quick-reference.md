# Facets Quick Reference

## Adding Facets to a Filter Component

### 1. Import and use context

```tsx
import { FacetCountsContext } from "src/hooks/useFacetCounts";

const MyFilter: React.FC<Props> = () => {
  const { counts: facetCounts, loading: facetsLoading } = useContext(FacetCountsContext);
  
  // Access counts
  const tagCounts = facetCounts.tags;  // Map<string, LabeledFacetCount>
  const count = tagCounts.get("123")?.count;
  const label = tagCounts.get("123")?.label;
};
```

### 2. Build candidates with counts

```tsx
const candidatesWithCounts = useMemo(() => {
  const hasValidFacets = facetCounts.tags.size > 0 && !facetsLoading;
  
  if (hasValidFacets && !searchQuery) {
    // Use facets directly
    const candidates: Option[] = [];
    facetCounts.tags.forEach((facet, id) => {
      if (facet.count === 0) return;
      candidates.push({ id, label: facet.label, count: facet.count });
    });
    return candidates.sort((a, b) => (b.count ?? 0) - (a.count ?? 0));
  }
  
  // Merge search results with counts
  return searchResults.map(item => ({
    ...item,
    count: facetCounts.tags.get(item.id)?.count
  }));
}, [searchResults, facetCounts, facetsLoading, searchQuery]);
```

## Adding Facets to a List Page

### 1. Import the hook

```tsx
import { useSceneFacetCounts, FacetCountsContext } from "src/hooks/useFacetCounts";
```

### 2. Call the hook

```tsx
const { counts: facetCounts, loading: facetLoading } = useSceneFacetCounts(filter, {
  isOpen: showSidebar,
  debounceMs: 300,
  includePerformerTags: sectionOpen["performer_tags"],
  includeCaptions: sectionOpen["captions"],
});
```

### 3. Provide context

```tsx
<FacetCountsContext.Provider value={{ counts: facetCounts, loading: facetLoading }}>
  <SidebarPane>
    {/* Filter components */}
  </SidebarPane>
</FacetCountsContext.Provider>
```

## Available Hooks

| Hook | Entity | Lazy Loading |
|------|--------|--------------|
| `useSceneFacetCounts` | Scenes | performer_tags, captions |
| `usePerformerFacetCounts` | Performers | None |
| `useGalleryFacetCounts` | Galleries | None |
| `useGroupFacetCounts` | Groups | None |
| `useStudioFacetCounts` | Studios | None |
| `useTagFacetCounts` | Tags | None |

## FacetCounts Structure

```typescript
interface FacetCounts {
  // Entity facets (Map<id, {count, label}>)
  tags: Map<string, LabeledFacetCount>;
  performers: Map<string, LabeledFacetCount>;
  studios: Map<string, LabeledFacetCount>;
  groups: Map<string, LabeledFacetCount>;
  performerTags: Map<string, LabeledFacetCount>;
  countries: Map<string, LabeledFacetCount>;
  parents: Map<string, LabeledFacetCount>;
  children: Map<string, LabeledFacetCount>;
  
  // Enum facets (Map<enum, count>)
  resolutions: Map<ResolutionEnum, number>;
  orientations: Map<OrientationEnum, number>;
  genders: Map<GenderEnum, number>;
  circumcised: Map<CircumisedEnum, number>;
  ratings: Map<number, number>;
  captions: Map<string, number>;
  
  // Boolean facets
  booleans: {
    organized: { true: number; false: number };
    interactive: { true: number; false: number };
    favorite: { true: number; false: number };
  };
}
```

## GraphQL Queries

### Scene Facets

```graphql
query SceneFacets($scene_filter: SceneFilterType, $limit: Int, 
                  $includePerformerTags: Boolean, $includeCaptions: Boolean) {
  sceneFacets(scene_filter: $scene_filter, limit: $limit,
              include_performer_tags: $includePerformerTags,
              include_captions: $includeCaptions) {
    tags { id label count }
    performers { id label count }
    studios { id label count }
    groups { id label count }
    performer_tags @include(if: $includePerformerTags) { id label count }
    resolutions { resolution count }
    orientations { orientation count }
    organized { value count }
    interactive { value count }
    ratings { rating count }
    captions @include(if: $includeCaptions) { language count }
  }
}
```

## Common Patterns

### Filtering zero counts

```typescript
// For entity filters - keep undefined, filter zero
const filtered = candidates.filter(c => c.count !== 0);

// For enum filters - filter undefined AND zero
const filtered = options.filter(o => {
  const count = counts.get(o.id);
  return count !== undefined && count > 0;
});
```

### Checking loading state

```typescript
// ALWAYS check loading state before filtering
const hasValidFacets = facetCounts.size > 0 && !facetsLoading;
if (!hasValidFacets) {
  return allCandidates;  // Don't filter
}
```

### Preserving modifier options

```typescript
// Modifier options (Any, None) should always be preserved
const modifiers = candidates.filter(c => c.className === "modifier-object");
const items = candidates.filter(c => c.className !== "modifier-object");
return [...modifiers, ...filteredItems];
```

### State preservation during lazy loading

When lazy-loading expensive facets, only update the specific facet to prevent React rendering glitches:

```typescript
// Detect lazy-load update
const isLazyLoadUpdate = 
  (includePerformerTags && !lastOptionsRef.current.includePerformerTags);

if (isLazyLoadUpdate) {
  // Partial update - preserves other facets
  setCounts((prev) => ({
    ...prev,
    performerTags: toMap(facets.performer_tags),
  }));
} else {
  // Full update
  setCounts((prev) => ({ ...allNewFacets }));
}
```

### Unique keys to prevent DOM recycling

Use `sectionID` prefix for candidate keys to prevent React from confusing items across filter sections:

```tsx
<CandidateList
  items={candidates}
  keyPrefix={sectionID}  // e.g., "performers", "tags"
/>

// In CandidateList
{items.map((p) => (
  <CandidateItem
    key={keyPrefix ? `${keyPrefix}-${p.id}` : p.id}
    // ...
  />
))}
```

## Testing

### Frontend test example

```typescript
import { describe, it, expect } from "vitest";
import { LabeledFacetCount, FacetCounts } from "src/hooks/useFacetCounts";

describe("MyFilter", () => {
  it("should filter zero counts", () => {
    const facets = new Map<string, LabeledFacetCount>([
      ["1", { count: 10, label: "Item 1" }],
      ["2", { count: 0, label: "Item 2" }],  // Zero
    ]);
    
    const candidates: Option[] = [];
    facets.forEach((f, id) => {
      if (f.count > 0) candidates.push({ id, label: f.label, count: f.count });
    });
    
    expect(candidates.map(c => c.id)).toEqual(["1"]);
  });
  
  it("should preserve existing facets during lazy load", () => {
    const existingState: Partial<FacetCounts> = {
      performers: new Map([["1", { count: 50, label: "Performer A" }]]),
      performerTags: new Map(), // Empty before lazy load
    };
    
    // Simulate lazy load update
    const lazyLoadUpdate = (prev: typeof existingState) => ({
      ...prev,
      performerTags: new Map([["100", { count: 500, label: "Tag X" }]]),
    });
    
    const newState = lazyLoadUpdate(existingState);
    
    // Verify performer_tags was updated
    expect(newState.performerTags?.size).toBe(1);
    // Verify performers was PRESERVED (same reference)
    expect(newState.performers).toBe(existingState.performers);
  });
});
```

### Backend test example

```go
func TestMyFacets(t *testing.T) {
  withRollbackTxn(func(ctx context.Context) error {
    facets, err := db.MyEntity.GetFacets(ctx, nil, 100)
    assert.NoError(t, err)
    assert.Greater(t, len(facets.Tags), 0)
    return nil
  })
}
```

