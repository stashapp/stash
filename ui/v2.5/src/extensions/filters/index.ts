/**
 * Custom Filter Components
 * 
 * These filter components extend and enhance the upstream filter system
 * with facet counting, better UX, and additional filter types.
 */

// Age filter
export { SidebarAgeFilter } from "./AgeFilter";

// Boolean filter
export { BooleanFilter, SidebarBooleanFilter } from "./BooleanFilter";

// Captions filter
export { SidebarCaptionsFilter } from "./CaptionsFilter";

// Circumcised filter
export { SidebarCircumcisedFilter } from "./CircumcisedFilter";

// Country filter
export { SidebarCountryFilter } from "./CountryFilter";

// Date filter
export { DateFilter, SidebarDateFilter, useDateCriterion } from "./DateFilter";

// Duration filter
export { DurationFilter } from "./DurationFilter";

// Gender filter
export { SidebarGenderFilter } from "./GenderFilter";

// Groups filter
export { SidebarGroupsFilter } from "./GroupsFilter";

// IsMissing filter
export { SidebarIsMissingFilter } from "./IsMissingFilter";

// LabeledId filter
export {
  LabeledIdFilter,
  getModifierCandidates,
  modifierValueToModifier,
  useSelectionState,
  useCriterion,
  useQueryState,
  useCandidates,
  useLabeledIdFilterState,
  makeQueryVariables,
  setObjectFilter,
} from "./LabeledIdFilter";

// MyFilterSidebar - note: FilteredSidebarHeader and useFilteredSidebarKeybinds 
// are exported from extensions/ui/FilterSidebar.tsx to avoid duplicates

// Number filter
export {
  NumberFilter,
  SidebarNumberFilter,
  NumberSelectedItems,
  useNumberCriterion,
} from "./NumberFilter";

// Orientation filter
export { OrientationFilter, SidebarOrientationFilter } from "./OrientationFilter";

// Path filter
export { PathFilter, SidebarPathFilter } from "./PathFilter";

// Performers filter
export { SidebarPerformersFilter } from "./PerformersFilter";

// PerformerTags filter
export { SidebarPerformerTagsFilter } from "./PerformerTagsFilter";

// Phash filter
export { PhashFilter, SidebarPhashFilter } from "./PhashFilter";

// Rating filter
export { RatingFilter, SidebarRatingFilter } from "./RatingFilter";

// Resolution filter
export { SidebarResolutionFilter } from "./ResolutionFilter";

// Selectable filter
export { ObjectsFilter, HierarchicalObjectsFilter } from "./SelectableFilter";

// SidebarDuration filter
export { SidebarDurationFilter } from "./SidebarDurationFilter";

// SidebarFilterSelector
export {
  SidebarFilterEditContext,
  SidebarFilterSelector,
  FilterVisibilityToggle,
  useFilterVisibility,
  FilterWrapper,
} from "./SidebarFilterSelector";

// SidebarListFilter
export {
  SelectedItem,
  SelectedList,
  CandidateList,
  SidebarListFilter,
  useStaticResults,
  type Option,
} from "./SidebarListFilter";

// StashID filter
export {
  StashIDFilter,
  SidebarStashIDFilter,
  useStashIDCriterion,
} from "./StashIDFilter";

// String filter
export {
  StringFilter,
  SidebarStringFilter,
  SidebarTattoosFilter,
  SelectedItems,
  ModifierControls,
  useStringCriterion,
  useModifierCriterion,
} from "./StringFilter";

// Studios filter
export {
  SidebarStudiosFilter,
  SidebarParentStudiosFilter,
} from "./StudiosFilter";

// Tags filter
export { SidebarTagsFilter } from "./TagsFilter";

// Utilities
export * from "./facetCandidateUtils";
