# Draft PR: Tagger Performer Hover Preview Enhancements

## PR Title
Enhance Tagger performer match UX with hover previews and mismatch warnings

## Base / Head
- Base: `stashapp/stash` `develop`
- Head: `Stash-KennyG/stash` `feature/6853-performer-hover-preview`

## Summary
- Adds a hover preview for matched performer selection in Scene Tagger by attaching the preview to the picker control, not the matched text label.
- Adds a scraped performer hover preview on the left performer label using scraped/remote data, while rendering through shared performer card UI for visual consistency.
- Adds mismatch context in the matched performer hover card for non-matching performer fields (including aliases/URLs as counts), to make potential cross-linkage easier to spot before save.
- Adds stash ID conflict warning inside the hover card when the local selected performer already has a stash ID for the same endpoint but with a different GUID.
- Uses a compact red endpoint chip for stash ID conflict warning, with full conflicting GUID available via chip tooltip.

## Files Changed
- `ui/v2.5/src/components/Performers/PerformerPopover.tsx`
- `ui/v2.5/src/components/Tagger/scenes/MatchedPerformerPreview.tsx`
- `ui/v2.5/src/components/Tagger/scenes/PerformerResult.tsx`
- `ui/v2.5/src/components/Tagger/scenes/ScrapedPerformerPreview.tsx` (new)
- `ui/v2.5/src/components/Tagger/styles.scss`

## Behavior Details
- **Picker hover target**: hover card now opens from the selected picker chip/control.
- **Left-side scraped hover**: hover on scraped performer text shows card built from scraped data (image/name/etc.) but styled with shared `PerformerCard`.
- **Mismatch rows**: shown only for fields where scraped has a value and differs from selected local value.
- **Aliases / URLs**: reported as counts only.
- **Stash ID conflict warning**:
  - shown only when `selected stash_id(endpoint) != incoming remote_site_id`
  - displayed inside card, above mismatch rows
  - compact red endpoint chip; GUID is tooltip text.

## Mismatch Fields Included
- Birthdate
- Death Date
- Ethnicity
- Hair Color
- Eye Color
- Height
- Weight
- Penis Length
- Circumcised
- Measurements
- Fake Tits
- Tattoos
- Piercings
- Career Start
- Career End
- Aliases (count)
- URLs (count)

## Test Plan
- [ ] Hover selected performer picker in Scene Tagger matched result and confirm card opens below picker.
- [ ] Hover left scraped performer name and confirm scraped preview card appears and matches shared performer card styling.
- [ ] Confirm mismatch rows render only for differing values and are hidden for exact matches.
- [ ] Confirm aliases/URLs mismatch rows show count only.
- [ ] Confirm stash ID warning appears only when endpoint matches and GUID differs.
- [ ] Confirm warning chip is compact red endpoint chip and GUID is visible in tooltip.
- [ ] Confirm no stash ID warning appears when GUIDs match.
- [ ] Confirm no regressions in create/skip/select/link flow for performer rows.
- [ ] Confirm lint/build checks pass.

## Notes
- This feature is designed to improve confidence in match correctness before save, especially in larger datasets where name collisions are more common.

## Compatibility Note: CSS Class Shape
- `PerformerPopover` now supports an optional `cardClassName` prop and applies it by appending to the existing wrapper class:
  - before: `className="tag-popover-card"`
  - after: `className={\`tag-popover-card ${cardClassName ?? ""}\`.trim()}`
- The base class `tag-popover-card` is preserved, so existing styling and selectors that target it continue to work.
- This is expected to be low-risk for downstream usage because:
  - the prop is optional
  - existing callers do not need to change
  - DOM structure is unchanged (same wrapper element).
