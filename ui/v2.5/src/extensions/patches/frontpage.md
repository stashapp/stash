# FrontPage Patches

This document describes the modifications to FrontPage components that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `Control.tsx` | AI recommendations support | +41 |
| `FrontPageConfig.tsx` | AI filter type handling | +9 |

---

## Control.tsx

**Purpose:** Add support for AI-powered scene recommendations on the front page.

**File:** `src/components/FrontPage/Control.tsx`

### Add import for IAIRecommendationFilter:

```diff
-import { FrontPageContent, ICustomFilter } from "src/core/config";
+import { FrontPageContent, ICustomFilter, IAIRecommendationFilter } from "src/core/config";
```

### Add import for AISceneRecommendationRow:

```diff
 import { SceneRecommendationRow } from "../Scenes/SceneRecommendationRow";
+import { AISceneRecommendationRow } from "src/extensions/components";
 import { StudioRecommendationRow } from "../Studios/StudioRecommendationRow";
```

### Add AIRecommendationResults component:

```typescript
interface IAIRecommendationProps {
  aiFilter: IAIRecommendationFilter;
}

const AIRecommendationResults: React.FC<IAIRecommendationProps> = ({
  aiFilter,
}) => {
  const intl = useIntl();

  function isTouchEnabled() {
    return "ontouchstart" in window || navigator.maxTouchPoints > 0;
  }

  const isTouch = isTouchEnabled();

  const header = aiFilter.message
    ? intl.formatMessage(
        { id: aiFilter.message.id },
        aiFilter.message.values
      )
    : aiFilter.title ?? "";

  // Currently only supports scenes
  if (aiFilter.mode === GQL.FilterMode.Scenes) {
    return (
      <AISceneRecommendationRow
        isTouch={isTouch}
        limit={aiFilter.limit}
        header={header}
      />
    );
  }

  return null;
};
```

### Add case in Control switch:

```diff
     case "CustomFilter":
       return <CustomFilterResults customFilter={content} />;
+    case "AIRecommendation":
+      return <AIRecommendationResults aiFilter={content as IAIRecommendationFilter} />;
     default:
```

---

## FrontPageConfig.tsx

**Purpose:** Handle AI recommendation filter type in front page configuration UI.

**File:** `src/components/FrontPage/FrontPageConfig.tsx`

### Add import:

```diff
 import {
   ISavedFilterRow,
   ICustomFilter,
+  IAIRecommendationFilter,
   FrontPageContent,
```

### Add case for AI filter title:

```diff
         return asCustomFilter.title ?? "";
+      case "AIRecommendation":
+        const asAIFilter = props.content as IAIRecommendationFilter;
+        if (asAIFilter.message)
+          return intl.formatMessage(
+            { id: asAIFilter.message.id },
+            asAIFilter.message.values
+          );
+        return asAIFilter.title ?? "";
     }
```

---

## Related Files

These changes depend on:
- `extensions/components/AISceneRecommendationRow.tsx` - The AI recommendation row component
- `core/config.ts` - Contains `IAIRecommendationFilter` type definition
- `core/StashService.ts` - GraphQL queries for AI recommendations

## Application Instructions

After upgrading upstream:

1. Check if upstream has added similar AI recommendation functionality
2. Ensure `IAIRecommendationFilter` type is still defined in `core/config.ts`
3. Verify `AISceneRecommendationRow` component exists in extensions
4. Apply the patches and test front page functionality
5. Run `yarn build` and test the front page displays correctly

