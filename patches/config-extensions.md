# Configuration Extensions

Fork additions to `ui/v2.5/src/core/config.ts` for new features.

## Changes Overview

| Addition | Purpose |
|----------|---------|
| `IAIRecommendationFilter` | Type for AI-powered recommendations |
| `FrontPageContent` union | Include AI recommendations |
| Background image options | More entity types support |
| `sidebarFilters` | Sidebar filter visibility settings |

---

## 1. IAIRecommendationFilter Type

Add after `ICustomFilter`:

```typescript
export interface IAIRecommendationFilter extends ITypename {
  __typename: "AIRecommendation";
  message?: IMessage;
  title?: string;
  mode: FilterMode;
  limit: number;
}
```

---

## 2. FrontPageContent Union

Update the type:

```typescript
// Before
export type FrontPageContent = ISavedFilterRow | ICustomFilter;

// After
export type FrontPageContent = ISavedFilterRow | ICustomFilter | IAIRecommendationFilter;
```

---

## 3. IUIConfig Background Images

Add these optional properties:

```typescript
export interface IUIConfig {
  // ... existing properties ...
  
  // ADD THESE:
  // if true a background image will be displayed on header
  enableGalleryBackgroundImage?: boolean;
  // if true a background image will be displayed on header
  enableImageBackgroundImage?: boolean;
  // if true a background image will be displayed on header
  enableSceneBackgroundImage?: boolean;
  
  // ... existing background image options ...
}
```

---

## 4. sidebarFilters Property

Add to `IUIConfig`:

```typescript
export interface IUIConfig {
  // ... existing properties ...
  
  pinnedFilters?: Record<string, string[]>;
  tableColumns?: Record<string, string[]>;
  sidebarFilters?: Record<string, string[]>;  // ADD THIS
  
  // ... rest of properties ...
}
```

---

## 5. AI Recommendations Helper

Add after `generateDefaultFrontPageContent`:

```typescript
function aiRecommendedScenes(
  intl: IntlShape,
  limit: number
): IAIRecommendationFilter {
  return {
    __typename: "AIRecommendation",
    message: {
      id: "recommended_scenes",
      values: { count: limit.toString() },
    },
    mode: FilterMode.Scenes,
    limit,
  };
}
```

---

## 6. Update generatePremadeFrontPageContent

Add AI recommendations to the array:

```typescript
export function generatePremadeFrontPageContent(intl: IntlShape) {
  return [
    aiRecommendedScenes(intl, 100),  // ADD THIS AT THE START
    recentlyReleased(intl, FilterMode.Scenes, "scenes"),
    recentlyAdded(intl, FilterMode.Scenes, "scenes"),
    // ... rest of existing content ...
  ];
}
```

---

## Related Files

| File | Change |
|------|--------|
| `core/StashService.ts` | Added `querySceneRecommendationsForScene` |
| `Scenes/AISceneRecommendationRow.tsx` | New component using recommendations |
| `FrontPage/Control.tsx` | Renders AI recommendations |
| `FrontPage/FrontPageConfig.tsx` | Config for AI recommendations |

---

## Note

These config changes support the recommendations feature. If the backend recommendations queries are not available, the AI recommendation rows will not function.

