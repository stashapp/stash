/**
 * Hook Extensions
 *
 * Custom React hooks for fork-specific functionality.
 * Re-exports from src/hooks/ for centralized access.
 */

// Focus management
export { default as useFocus, useFocusOnce } from "./useFocus";

// Facet counting system - re-export from src/hooks
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
} from "src/hooks/useFacetCounts";

// Scene-specific facets - re-export from src/hooks
export {
  useSceneFacets,
  useGalleryFacets,
  usePerformerFacets,
} from "src/hooks/useSceneFacets";

// Sidebar filter state management - re-export from src/hooks
export {
  useSidebarFilters,
  type SidebarFilterDefinition,
} from "src/hooks/useSidebarFilters";

// Batched filter counts - re-export from src/hooks
export {
  useBatchedFilterCounts,
  useOptimizedFilterCounts,
  type BatchedCounts,
} from "src/hooks/useBatchedFilterCounts";
