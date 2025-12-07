/**
 * Faceted Search Hooks
 *
 * This module provides hooks for efficient faceted filtering in the UI.
 * Currently uses client-side batching and caching while awaiting
 * a backend aggregation endpoint.
 *
 * Future enhancement: When the backend `sceneFacets` endpoint is implemented,
 * these hooks will be updated to use a single GraphQL query that returns
 * all facet counts at once, significantly improving performance.
 *
 * @see docs/proposals/facets-aggregation-endpoint.md (if exists)
 */

export {
  useSceneFacets,
  useGalleryFacets,
  usePerformerFacets,
  type FacetCounts,
} from "../useSceneFacets";

export {
  useBatchedFilterCounts,
  useOptimizedFilterCounts,
  type BatchedCounts,
} from "../useBatchedFilterCounts";

export {
  FacetsProvider,
  useFacets,
  useFilterFacets,
} from "../useFacetsContext";

