# List Components Patches

This document describes the modifications to upstream List components that need to be reapplied after upgrading.

## Overview

| File | Change | Lines |
|------|--------|-------|
| `ListFilter.tsx` | Added optional `placeholder` prop to SearchTermInput | +5 |
| `ListTable.tsx` | Checkbox wrapper for custom styling | +8 |
| `Pagination.tsx` | Bold formatting in pagination text | +4 |

---

## ListFilter.tsx

**Purpose:** Allow custom placeholder text in the search input.

**File:** `src/components/List/ListFilter.tsx`

```diff
@@ -63,7 +63,8 @@ export const SearchTermInput: React.FC<{
   filter: ListFilterModel;
   onFilterUpdate: (newFilter: ListFilterModel) => void;
   focus?: ReturnType<typeof useFocus>;
-}> = ({ filter, onFilterUpdate, focus: providedFocus }) => {
+  placeholder?: string;
+}> = ({ filter, onFilterUpdate, focus: providedFocus, placeholder }) => {
   const intl = useIntl();
   const [localInput, setLocalInput] = useState(filter.searchTerm);

@@ -107,7 +108,9 @@ export const SearchTermInput: React.FC<{
       focus={focus}
       value={localInput}
       setValue={onSetQuery}
-      placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
+      placeholder={
+        placeholder ?? `${intl.formatMessage({ id: "actions.search" })}…`
+      }
     />
   );
 };
```

---

## ListTable.tsx

**Purpose:** Wrap checkbox in a div for custom styling (allows custom checkbox appearance via CSS).

**File:** `src/components/List/ListTable.tsx`

```diff
@@ -75,9 +75,11 @@ export const ListTable = <T extends { id: string }>(
     return (
       <tr key={item.id}>
         <td className="select-col">
-          <label>
+          <div className="checkbox-wrapper">
             <Form.Control
               type="checkbox"
+              id={item.id}
+              // #2750 - add mousetrap class to ensure keyboard shortcuts work
               checked={selectedIds.has(item.id)}
               onChange={() =>
                 onSelectChange(item.id, !selectedIds.has(item.id), shiftKey)
@@ -89,7 +91,10 @@ export const ListTable = <T extends { id: string }>(
                 event.stopPropagation();
               }}
             />
-          </label>
+            <label htmlFor={item.id}>
+              <span />
+            </label>
+          </div>
         </td>
```

**Note:** The `id` attribute and associated `label htmlFor` enable clicking the label to toggle the checkbox. The `<span />` inside the label is used as a styling hook for custom checkbox appearance.

---

## Pagination.tsx

**Purpose:** Make the count numbers bold for better visual hierarchy.

**File:** `src/components/List/Pagination.tsx`

```diff
@@ -263,7 +263,10 @@ export const PaginationIndex: React.FC<IPaginationIndexProps> = PatchComponent(

     return (
       <span className="filter-container text-muted paginationIndex center-text">
-        {indexText}
+        <b>
+          {intl.formatNumber(firstItemCount)}-{intl.formatNumber(lastItemCount)}
+        </b>{" "}
+        of <b>{intl.formatNumber(totalItems)}</b>
         <br />
         {metadataByline}
       </span>
```

---

## Application Instructions

After upgrading upstream:

1. Check if these files have changed in the new version
2. If unchanged, apply the patches directly
3. If changed, manually merge the changes
4. Run `yarn build` and `yarn test` to verify

## Related SCSS

The checkbox wrapper styling is defined in `extensions/styles/_list-components.scss`. Ensure this SCSS is still being imported after upstream upgrades.

