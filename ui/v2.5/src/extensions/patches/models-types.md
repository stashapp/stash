# Models & Types Patches

This document describes the modifications to model and type files that need to be reapplied after upgrading.

## Overview

| File | Change | Lines | Type |
|------|--------|-------|------|
| `criteria/piercings.ts` | New filter criterion | +13 | New file |
| `criteria/tattoos.ts` | New filter criterion | +13 | New file |
| `types.ts` | Added `count` to ILabeledId | +1 | Modification |
| `sceneQueue.ts` | Added IFileObject interface | +10 | Modification |

---

## New Files

### piercings.ts

**File:** `src/models/list-filter/criteria/piercings.ts`

This is a **new file** that doesn't exist in upstream. Create it:

```typescript
import { StringCriterion, StringCriterionOption } from "./criterion";

export const PiercingsCriterionOption = new StringCriterionOption({
  messageID: "piercings",
  type: "piercings",
  makeCriterion: () => new PiercingsCriterion(),
});

export class PiercingsCriterion extends StringCriterion {
  constructor() {
    super(PiercingsCriterionOption);
  }
}
```

### tattoos.ts

**File:** `src/models/list-filter/criteria/tattoos.ts`

This is a **new file** that doesn't exist in upstream. Create it:

```typescript
import { StringCriterion, StringCriterionOption } from "./criterion";

export const TattoosCriterionOption = new StringCriterionOption({
  messageID: "tattoos",
  type: "tattoos",
  makeCriterion: () => new TattoosCriterion(),
});

export class TattoosCriterion extends StringCriterion {
  constructor() {
    super(TattoosCriterionOption);
  }
}
```

---

## Modified Files

### types.ts

**Purpose:** Add optional `count` property to `ILabeledId` for sidebar filter counts.

**File:** `src/models/list-filter/types.ts`

```diff
 export interface ILabeledId {
   id: string;
   label: string;
+  count?: number; // optional count for sidebar filters
 }
```

---

### sceneQueue.ts

**Purpose:** Add `IFileObject` interface and `files` property to `QueuedScene` for queue file info.

**File:** `src/models/sceneQueue.ts`

```diff
 import { FilterMode, Scene } from "src/core/generated-graphql";
 import { ListFilterModel } from "./list-filter/filter";
 import { INamedObject } from "src/utils/navigation";

+export interface IFileObject {
+  id: string;
+  duration: number;
+  height: number;
+  path: string;
+  width: number;
+  size: number;
+}
+
 export type QueuedScene = Pick<Scene, "id" | "title" | "date" | "paths"> & {
   performers?: INamedObject[] | null;
   studio?: INamedObject | null;
+  files: IFileObject[];
 };
```

---

## Application Instructions

After upgrading upstream:

1. **New files** (piercings.ts, tattoos.ts):
   - Check if upstream has added similar criteria
   - If not, create the files
   - May need to register in criteria factory/index

2. **types.ts**:
   - Simple property addition
   - Check if upstream changed `ILabeledId` interface

3. **sceneQueue.ts**:
   - Check if upstream changed `QueuedScene` type
   - `IFileObject` interface may conflict if upstream adds similar

## Related Files

These changes are used by:
- `extensions/filters/` - Custom filter components
- `extensions/ui/` - Sidebar filter with counts
- Scene queue/playlist functionality
