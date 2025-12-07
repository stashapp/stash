/**
 * Hook Extensions
 *
 * Custom React hooks for fork-specific functionality.
 */

// Focus management - re-export from upstream utils
export { default as useFocus, useFocusOnce } from "src/utils/focus";

// Facet counting system
export {
  useSceneFacetCounts,
  usePerformerFacetCounts,
  useGalleryFacetCounts,
  useGroupFacetCounts,
  useStudioFacetCounts,
  useTagFacetCounts,
  useFacetCounts,
  FacetCountsContext,
  useFacetCountsContext,
  type FacetCounts,
} from "./useFacetCounts";

// Scene-specific facets
export {
  useSceneFacets,
  useGalleryFacets,
  usePerformerFacets,
} from "./useSceneFacets";

// Sidebar filter state management
export {
  useSidebarFilters,
  type SidebarFilterDefinition,
} from "./useSidebarFilters";

// Batched filter counts
export {
  useBatchedFilterCounts,
  useOptimizedFilterCounts,
  type BatchedCounts,
} from "./useBatchedFilterCounts";

// Facets context
export { FacetsProvider, useFacets } from "./useFacetsContext";
