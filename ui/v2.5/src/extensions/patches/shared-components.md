# Shared Components Patches

This document describes the modifications to shared UI components that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `ClearableInput.tsx` | Added search icon button | +11 |
| `CollapseButton.tsx` | Added `disabled` prop | +2 |
| `DetailItem.tsx` | Added show more/less + `messageId`/`heading` props | +55 |
| `GridCard/GridCard.tsx` | Checkbox wrapper + `titleOnImage` prop | +41 |
| `Sidebar.tsx` | Added `disabled` prop passthrough | +7 |
| `TagLink.tsx` | Added `GalleryDetailedLink` + `GalleryPopover` import | +32 |

---

## ClearableInput.tsx

**Purpose:** Add a search icon when the input is empty.

**File:** `src/components/Shared/ClearableInput.tsx`

### Add import:

```diff
-import { faTimes } from "@fortawesome/free-solid-svg-icons";
+import { faMagnifyingGlass, faTimes } from "@fortawesome/free-solid-svg-icons";
```

### Add search icon button:

```diff
       )}
+      {!queryClearShowing && (
+        <Button
+          variant="secondary"
+          title={intl.formatMessage({ id: "actions.clear" })}
+          className="clearable-text-field-search"
+        >
+          <Icon icon={faMagnifyingGlass} />
+        </Button>
+      )}
     </div>
```

---

## CollapseButton.tsx

**Purpose:** Support disabling the collapse toggle.

**File:** `src/components/Shared/CollapseButton.tsx`

```diff
   open?: boolean;
+  disabled?: boolean;
 }

...

   function toggleOpen() {
+    if (props.disabled) return;
     const nv = !open;
```

---

## DetailItem.tsx

**Purpose:** Add show more/less functionality and additional props.

**File:** `src/components/Shared/DetailItem.tsx`

### Add helper function:

```typescript
export function maybeRenderShowMoreLess(
  height: number,
  limit: number,
  ref: React.MutableRefObject<HTMLDivElement | null>,
  setCollapsed: React.Dispatch<React.SetStateAction<boolean>>,
  collapsed: boolean
) {
  if (height < limit) {
    return;
  }
  return (
    <span
      className={`show-${collapsed ? "more" : "less"}`}
      onClick={() => {
        const container = ref.current;
        if (container == null) {
          return;
        }
        if (container.style.maxHeight) {
          container.style.maxHeight = "";
        } else {
          container.style.maxHeight = container.scrollHeight + "px";
        }
        setCollapsed(!collapsed);
      }}
    >
      {collapsed ? "Show more" : "Show less"}
      <Icon className="fa-solid" icon={collapsed ? faCaretDown : faCaretUp} />
    </span>
  );
}
```

### Add props to interface:

```diff
 interface IDetailItem {
   id?: string | null;
   label?: React.ReactNode;
+  messageId?: string;
+  heading?: React.ReactNode;
   value?: React.ReactNode;
```

### Update component to use new props:

```diff
-  const message = label ?? <FormattedMessage id={id} />;
+  const message = label ?? <FormattedMessage id={messageId ?? id} />;
```

---

## GridCard/GridCard.tsx

**Purpose:** Custom checkbox styling and `titleOnImage` prop support.

**File:** `src/components/Shared/GridCard/GridCard.tsx`

### Add `titleOnImage` prop:

```diff
   onMove?: (srcIds: string[], targetId: string, after: boolean) => void;
+  titleOnImage?: boolean;
 }
```

### Wrap checkbox in div with label:

```diff
 const Checkbox: React.FC<{
+  cardID: string;
   selected?: boolean;
   onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
-}> = ({ selected = false, onSelectedChanged }) => {
+}> = ({ cardID, selected = false, onSelectedChanged }) => {
   let shiftKey = false;
   return (
-    <Form.Control
-      type="checkbox"
-      className="card-check mousetrap"
-      ...
-    />
+    <div className="checkbox-wrapper">
+      <Form.Control
+        type="checkbox"
+        id={cardID}
+        className="card-check mousetrap"
+        ...
+      />
+      <label htmlFor={cardID}>
+        <span />
+      </label>
+    </div>
   );
 };
```

### Generate cardID and pass to Checkbox:

```diff
+  const cardID = props.url.substring(1).split("?")[0].replace("/", "-");
...
           <Checkbox
+            cardID={cardID}
             selected={props.selected}
```

### Use titleOnImage for link title:

```diff
         <Link
+          title={props.titleOnImage ? props.title.toString() : ""}
           to={props.url}
```

### Add data-value attribute:

```diff
-          <h5 className="card-section-title flex-aligned">
+          <h5
+            data-value={props.title}
+            className="card-section-title flex-aligned"
+          >
```

---

## Sidebar.tsx

**Purpose:** Support disabled state for sidebar sections.

**File:** `src/components/Shared/Sidebar.tsx`

### Add to context interface:

```diff
 interface IContext {
   sectionOpen: SidebarSectionStates;
   setSectionOpen: (section: string, open: boolean) => void;
+  disabled?: boolean;
 }
```

### Add prop and logic:

```diff
     sectionID?: string;
+    disabled?: boolean;
   }>
 > = ({
     ...
+    disabled: disabledProp,
     children,
   }) => {
+  const disabled = disabledProp ?? contextState?.disabled ?? false;
   ...
       <CollapseButton
         ...
+        disabled={disabled}
       >
```

---

## TagLink.tsx

**Purpose:** Add GalleryDetailedLink component with hover popover.

**File:** `src/components/Shared/TagLink.tsx`

### Add import:

```diff
+import { GalleryPopover } from "src/extensions/components/GalleryPopover";
```

### Add new component:

```typescript
interface IGalleryDetailedLinkProps {
  gallery: GQL.SlimGalleryDataFragment;
  linkType?: "gallery";
  className?: string;
  hoverPlacement?: Placement;
}

export const GalleryDetailedLink: React.FC<IGalleryDetailedLinkProps> = ({
  gallery,
  linkType = "gallery",
  className,
  hoverPlacement,
}) => {
  const link = useMemo(() => {
    switch (linkType) {
      case "gallery":
        return `/galleries/${gallery.id}`;
    }
  }, [gallery, linkType]);

  const title = galleryTitle(gallery);

  return (
    <CommonLinkComponent link={link} className={className}>
      <GalleryPopover gallery={gallery ?? ""} placement={hoverPlacement}>
        {title}
      </GalleryPopover>
    </CommonLinkComponent>
  );
};
```

---

## Application Instructions

After upgrading upstream:

1. **ClearableInput.tsx** - Check if upstream added similar search icon
2. **CollapseButton.tsx** - Simple prop addition, should apply cleanly
3. **DetailItem.tsx** - Most complex change; check if upstream API changed
4. **GridCard.tsx** - Checkbox styling may conflict if upstream changes checkbox implementation
5. **Sidebar.tsx** - Simple prop passthrough
6. **TagLink.tsx** - Import and new component; ensure GalleryPopover exists

## Related Files

- `extensions/components/GalleryPopover.tsx`
- `extensions/styles/` - CSS for checkbox wrappers, icons

