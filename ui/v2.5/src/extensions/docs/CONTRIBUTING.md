# Contributing to Extensions

This guide explains how to add new features to the extensions system.

## Quick Start

All fork-specific code goes in `src/extensions/`. Follow these patterns:

## Important: Do NOT Modify Upstream Files

The `src/components/List/Filters/` directory contains **clean v0.29.3 upstream code**.

❌ **DON'T:**
- Add facet counts to files in `components/List/Filters/`
- Create `My*` files in `components/` directories
- Modify upstream components directly

✅ **DO:**
- Create/modify filters in `src/extensions/filters/`
- Create/modify lists in `src/extensions/lists/`
- Import from `src/extensions/filters` in extension code

## Adding a Filter Component

### 1. Create the Component

```tsx
// src/extensions/filters/MyFilter.tsx
import React from "react";
import { useFacetCountsContext } from "src/extensions/hooks";

interface ISidebarMyFilterProps {
  filter: ListFilterModel;
  onFilterUpdate: (filter: ListFilterModel) => void;
}

export const SidebarMyFilter: React.FC<ISidebarMyFilterProps> = ({
  filter,
  onFilterUpdate
}) => {
  const { counts, loading } = useFacetCountsContext();
  
  // Your filter implementation
  return (
    <div className="my-filter">
      {/* Filter UI */}
    </div>
  );
};
```

### 2. Export from Index

```typescript
// src/extensions/filters/index.ts
export { SidebarMyFilter } from "./MyFilter";
```

### 3. Use It

```typescript
import { SidebarMyFilter } from "src/extensions/filters";

// In your list component's sidebar
<SidebarMyFilter filter={filter} onFilterUpdate={handleUpdate} />
```

---

## Adding a Hook

### 1. Create the Hook

```typescript
// src/extensions/hooks/useMyHook.ts
import { useState, useEffect } from "react";

export function useMyHook(param: string) {
  const [value, setValue] = useState<string | null>(null);
  
  useEffect(() => {
    // Your logic
  }, [param]);
  
  return value;
}
```

### 2. Export from Index

```typescript
// src/extensions/hooks/index.ts
export { useMyHook } from "./useMyHook";
```

### 3. Use It

```typescript
import { useMyHook } from "src/extensions/hooks";

const value = useMyHook("param");
```

---

## Adding Styles

### 1. Create SCSS Partial

```scss
// src/extensions/styles/_my-feature.scss
.my-feature {
  // Styles here
  
  &__element {
    color: var(--color-brand-accent);
  }
}
```

### 2. Import in Index

```scss
// src/extensions/styles/index.scss
@import "variables";
@import "facets";
@import "sidebar";
@import "filter-tags";
@import "my-feature";  // Add this line
```

---

## Adding a List Feature

List components are in `src/extensions/lists/`. To add a feature:

### 1. Find the Right Component

- `PerformerList.tsx` - Performer grid/table
- `SceneList.tsx` - Scene grid/table
- `GalleryList.tsx` - Gallery grid
- `GroupList.tsx` - Group/movie grid
- `StudioList.tsx` - Studio grid
- `TagList.tsx` - Tag grid

### 2. Add Your Feature

Features are typically added to the sidebar or toolbar:

```tsx
// In the list component
const MyFilterSidebarSections = (props) => (
  <>
    <ExistingFilters {...props} />
    <MySidebarSection {...props} />  {/* Add your section */}
  </>
);
```

### 3. Add Any Supporting Hooks

If your feature needs state management, add a hook:

```typescript
// src/extensions/hooks/useMyFeature.ts
export function useMyFeature(filter: ListFilterModel) {
  // Your logic
}
```

---

## Adding a UI Component

Reusable UI components go in `src/extensions/ui/`:

### 1. Create Component

```tsx
// src/extensions/ui/MyComponent.tsx
import React from "react";

interface IMyComponentProps {
  value: string;
  onChange: (value: string) => void;
}

export const MyComponent: React.FC<IMyComponentProps> = ({
  value,
  onChange
}) => {
  return (
    <div className="my-component">
      {/* UI */}
    </div>
  );
};
```

### 2. Export from Index

```typescript
// src/extensions/ui/index.ts
export * from "./MyComponent";
```

---

## Import Rules

### Always Use Absolute Imports

```typescript
// ✅ Correct
import { SceneCardsGrid } from "src/components/Scenes/SceneCardsGrid";
import { useFindPerformers } from "src/core/StashService";

// ❌ Wrong - will break when files move
import { SceneCardsGrid } from "../Scenes/SceneCardsGrid";
```

### Import from Extensions

```typescript
// From specific module
import { SidebarTagsFilter } from "src/extensions/filters";
import { useSceneFacetCounts } from "src/extensions/hooks";

// Or from main entry
import { SidebarTagsFilter, useSceneFacetCounts } from "src/extensions";
```

---

## Testing Your Changes

### 1. Run Extension Tests

```bash
cd ui/v2.5
yarn test --run
```

Extension tests are in `extensions/__tests__/`:

| Test File | Tests |
|-----------|-------|
| `useFacetCounts.test.ts` | 14 tests - Facet interfaces, lazy loading |
| `facetCandidateUtils.test.ts` | 18 tests - Candidate building logic |
| `GroupsFilter.test.ts` | 8 tests - Groups filter logic |

### 2. TypeScript Check

```bash
yarn tsc --noEmit
```

### 3. Build Check

```bash
yarn build
```

### 4. Manual Testing

1. Start the dev server
2. Navigate to the affected list page
3. Open the sidebar
4. Verify your feature works

### Adding New Tests

Create test files in `extensions/__tests__/`:

```typescript
// extensions/__tests__/MyFeature.test.ts
import { describe, it, expect } from "vitest";
import { myFunction } from "src/extensions/hooks/myHook";

describe("MyFeature", () => {
  it("should do something", () => {
    expect(myFunction()).toBe(expected);
  });
});
```

---

## Checklist

Before submitting changes:

- [ ] Code is in `src/extensions/` (not scattered elsewhere)
- [ ] Uses absolute imports (`src/...`)
- [ ] Exported from the appropriate `index.ts`
- [ ] Styles use CSS variables where possible
- [ ] TypeScript compiles without errors
- [ ] Build succeeds
- [ ] Tested manually in browser

