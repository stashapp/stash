/**
 * Faceted Search Hooks
 *
 * This module provides hooks for efficient faceted filtering in the UI.
 * Uses client-side batching and caching with backend aggregation support.
 */

export {
  useSceneFacets,
  useGalleryFacets,
  usePerformerFacets,
  type FacetCounts,
} from "src/extensions/hooks/useSceneFacets";

export {
  useBatchedFilterCounts,
  useOptimizedFilterCounts,
  type BatchedCounts,
} from "src/extensions/hooks/useBatchedFilterCounts";
