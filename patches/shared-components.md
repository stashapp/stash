# Shared Components Modifications

## Overview

Shared components have been enhanced with additional functionality for better UX.

## Files

| File | Lines | Changes |
|------|-------|---------|
| `ClearableInput.tsx` | +9 / -2 | Search icon, styling |
| `CollapseButton.tsx` | +2 | Minor tweak |
| `DetailItem.tsx` | +45 / -10 | Collapsible support, truncation |
| `GridCard/GridCard.tsx` | +30 / -11 | Card width, mobile detection |
| `Sidebar.tsx` | +7 | Additional props |
| `TagLink.tsx` | +32 | Hierarchical display |
| `styles.scss` | ✅ Extracted | Now in `extensions/styles/_shared-components.scss` |

## Detailed Changes

### ClearableInput.tsx

Added search icon display when input is empty:

```diff
+ import { faMagnifyingGlass, faTimes } from "@fortawesome/free-solid-svg-icons";

// In render, add search icon when not showing clear button:
+ {!queryClearShowing && (
+   <Button variant="secondary" className="search-icon">
+     <Icon icon={faMagnifyingGlass} />
+   </Button>
+ )}
```

### CollapseButton.tsx

Minor styling adjustment:

```diff
+ className={cx(className, "collapse-button")}
```

### DetailItem.tsx

Enhanced with collapsible support and text truncation:

- Added `collapsible` prop for expandable sections
- Added text truncation with "show more" for long values
- Added `valueClassName` prop for custom styling
- Added `truncateValue` prop for controlling truncation

**Key additions:**
```typescript
interface IDetailItemProps {
  // ... existing props
  collapsible?: boolean;
  defaultCollapsed?: boolean;
  truncateValue?: boolean;
  valueClassName?: string;
}
```

### GridCard/GridCard.tsx

Enhanced card width calculation:

- Added `useCardWidth` hook for responsive card sizing
- Better mobile detection using `ScreenUtils.isMobile()`
- Custom width overrides via CSS variables

### Sidebar.tsx

Additional props for extension support:

```diff
+ filterHeader?: React.ReactNode;
+ onClearFilters?: () => void;
```

### TagLink.tsx

Added hierarchical tag display:

- Shows parent tag chain (e.g., "Parent > Child > Tag")
- Hover tooltip with full hierarchy
- Click navigation to parent tags

## Merge Strategy

1. **Accept upstream first** - Get base component updates
2. **Re-apply enhancements** - Add the additional props and features
3. **Test integration** - Verify with detail panels and list views

## Related Styles

All shared component styling is in `extensions/styles/_shared-components.scss`.


