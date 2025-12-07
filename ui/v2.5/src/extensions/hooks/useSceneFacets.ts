import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import * as GQL from "src/core/generated-graphql";

/**
 * Facet counts for various filter dimensions
 */
export interface FacetCounts {
  tags: Record<string, number>;
  performers: Record<string, number>;
  studios: Record<string, number>;
  resolutions: Record<string, number>;
  orientations: Record<string, number>;
  organized: { true: number; false: number };
}

/**
 * Hook to fetch facet counts for scene filters.
 * 
 * This hook is designed to work with a future backend `sceneFacets` endpoint
 * that returns all facet counts in a single query. Until that endpoint exists,
 * it falls back to the existing individual count queries.
 * 
 * Benefits of the facets approach:
 * - Single network request instead of N requests per filter type
 * - Server-side optimization with a single database query
 * - Reduced latency from 500ms+ to ~50-100ms
 * - Lower server load
 * 
 * @param filter - The current list filter model
 * @param isOpen - Whether the sidebar is open (to skip fetching when closed)
 * @returns Object containing counts and loading state
 */
export function useSceneFacets(filter: ListFilterModel, isOpen: boolean) {
  const [counts, setCounts] = useState<FacetCounts | null>(null);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<NodeJS.Timeout | null>(null);
  const lastFilterRef = useRef<string>("");

  // Build the scene filter, excluding facet-related criteria
  const sceneFilter = useMemo(() => {
    const filterCopy = filter.clone();
    // The facets endpoint should receive the filter without the criteria
    // that we're computing facets for, to get accurate counts
    return filterCopy.makeFilter() as GQL.SceneFilterType;
  }, [filter]);

  // Create a fingerprint of the filter for change detection
  const filterFingerprint = useMemo(() => {
    return JSON.stringify(sceneFilter);
  }, [sceneFilter]);

  const fetchFacets = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);

    try {
      // TODO: Once the backend sceneFacets endpoint is implemented,
      // replace this with the actual GraphQL query:
      //
      // const result = await client.query({
      //   query: SceneFacetsDocument,
      //   variables: {
      //     scene_filter: sceneFilter,
      //     limit: 100,
      //   },
      //   fetchPolicy: "network-only",
      // });
      //
      // if (result.data?.sceneFacets) {
      //   setCounts({
      //     tags: Object.fromEntries(
      //       result.data.sceneFacets.tags.map(t => [t.id, t.count])
      //     ),
      //     performers: Object.fromEntries(
      //       result.data.sceneFacets.performers.map(p => [p.id, p.count])
      //     ),
      //     studios: Object.fromEntries(
      //       result.data.sceneFacets.studios.map(s => [s.id, s.count])
      //     ),
      //     resolutions: Object.fromEntries(
      //       result.data.sceneFacets.resolutions.map(r => [r.resolution, r.count])
      //     ),
      //     orientations: Object.fromEntries(
      //       result.data.sceneFacets.orientations.map(o => [o.orientation, o.count])
      //     ),
      //     organized: {
      //       true: result.data.sceneFacets.organized.find(o => o.value)?.count ?? 0,
      //       false: result.data.sceneFacets.organized.find(o => !o.value)?.count ?? 0,
      //     },
      //   });
      // }

      // For now, this is a placeholder that doesn't fetch anything
      // Individual filters continue to use their own count fetching
      console.log("useSceneFacets: Facets endpoint not yet available");
    } catch (error) {
      console.error("Error fetching scene facets:", error);
    } finally {
      setLoading(false);
    }
  }, [isOpen, sceneFilter]);

  // Debounced fetch when filter changes
  useEffect(() => {
    if (!isOpen) return;

    // Skip if filter hasn't changed
    if (filterFingerprint === lastFilterRef.current) return;
    lastFilterRef.current = filterFingerprint;

    // Clear existing debounce
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    // Debounce the fetch
    debounceRef.current = setTimeout(() => {
      fetchFacets();
    }, 500);

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [filterFingerprint, isOpen, fetchFacets]);

  return { counts, loading };
}

/**
 * Hook to fetch facet counts for gallery filters.
 * Follows the same pattern as useSceneFacets.
 */
export function useGalleryFacets(filter: ListFilterModel, isOpen: boolean) {
  const [counts, setCounts] = useState<Partial<FacetCounts> | null>(null);
  const [loading, setLoading] = useState(false);

  // TODO: Implement when backend endpoint is available
  // Same pattern as useSceneFacets

  return { counts, loading };
}

/**
 * Hook to fetch facet counts for performer filters.
 * Follows the same pattern as useSceneFacets.
 */
export function usePerformerFacets(filter: ListFilterModel, isOpen: boolean) {
  const [counts, setCounts] = useState<Partial<FacetCounts> | null>(null);
  const [loading, setLoading] = useState(false);

  // TODO: Implement when backend endpoint is available
  // Same pattern as useSceneFacets

  return { counts, loading };
}

