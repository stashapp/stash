/// <reference types="node" />
import React, {
  createContext,
  useContext,
  useCallback,
  useMemo,
  useRef,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import * as GQL from "src/core/generated-graphql";
import { getClient } from "src/core/StashService";

/**
 * Facet counts for various filter dimensions
 */
export interface FacetCounts {
  tags: Map<string, number>;
  performers: Map<string, number>;
  studios: Map<string, number>;
  groups: Map<string, number>;
  galleries: Map<string, number>;
  resolutions: Map<string, number>;
  orientations: Map<string, number>;
  genders: Map<string, number>;
  organized: { true: number; false: number } | null;
}

interface FacetsContextValue {
  /** Current facet counts */
  counts: FacetCounts;
  /** Whether facets are currently loading */
  loading: boolean;
  /** Request a count for a specific entity */
  requestCount: (type: keyof FacetCounts, id: string) => void;
  /** Get a count if available */
  getCount: (type: keyof FacetCounts, id: string) => number | undefined;
  /** Clear all cached counts */
  clearCache: () => void;
  /** Set external counts (e.g., from a future facets endpoint) */
  setExternalCounts: (counts: Partial<FacetCounts>) => void;
}

const FacetsContext = createContext<FacetsContextValue | null>(null);

interface FacetsProviderProps {
  children: ReactNode;
  filter: ListFilterModel;
  isOpen: boolean;
}

/**
 * Pending count request
 */
interface PendingRequest {
  type: keyof FacetCounts;
  id: string;
  resolve: (count: number) => void;
}

/**
 * Provider component that manages facet counts for sidebar filters.
 * 
 * This provider:
 * 1. Centralizes count fetching for all filter types
 * 2. Batches requests to reduce network calls
 * 3. Caches results with TTL-based invalidation
 * 4. Debounces filter changes
 * 
 * When the backend facets endpoint becomes available, this provider
 * will be updated to use it instead of individual queries.
 */
export function FacetsProvider({
  children,
  filter,
  isOpen,
}: FacetsProviderProps) {
  const [counts, setCounts] = useState<FacetCounts>({
    tags: new Map(),
    performers: new Map(),
    studios: new Map(),
    groups: new Map(),
    galleries: new Map(),
    resolutions: new Map(),
    orientations: new Map(),
    genders: new Map(),
    organized: null,
  });
  const [loading, setLoading] = useState(false);

  // Pending requests queue
  const pendingRef = useRef<PendingRequest[]>([]);
  const debounceRef = useRef<NodeJS.Timeout | null>(null);
  const cacheRef = useRef<Map<string, { count: number; timestamp: number }>>(
    new Map()
  );

  // Cache configuration
  const CACHE_TTL = 5 * 60 * 1000; // 5 minutes
  const DEBOUNCE_MS = 500;
  const BATCH_SIZE = 10;
  const MAX_CONCURRENT = 4;

  // Filter fingerprint for cache invalidation
  const filterFingerprint = useMemo(() => {
    const filterCopy = filter.clone();
    return JSON.stringify({
      filter: filterCopy.makeFilter(),
      mode: filter.mode,
    });
  }, [filter]);

  // Build cache key
  const getCacheKey = useCallback(
    (type: string, id: string) => {
      return `${filterFingerprint}:${type}:${id}`;
    },
    [filterFingerprint]
  );

  // Check if cache entry is valid
  const isCacheValid = useCallback((key: string): boolean => {
    const entry = cacheRef.current.get(key);
    if (!entry) return false;
    return Date.now() - entry.timestamp < CACHE_TTL;
  }, []);

  // Get count from cache or counts state
  const getCount = useCallback(
    (type: keyof FacetCounts, id: string): number | undefined => {
      // Check cache first
      const cacheKey = getCacheKey(type, id);
      if (isCacheValid(cacheKey)) {
        return cacheRef.current.get(cacheKey)?.count;
      }

      // Check state
      const typeMap = counts[type];
      if (typeMap instanceof Map) {
        return typeMap.get(id);
      }
      return undefined;
    },
    [counts, getCacheKey, isCacheValid]
  );

  // Set count in cache and state
  const setCount = useCallback(
    (type: keyof FacetCounts, id: string, count: number) => {
      const cacheKey = getCacheKey(type, id);
      cacheRef.current.set(cacheKey, { count, timestamp: Date.now() });

      setCounts((prev) => {
        const typeMap = prev[type];
        if (typeMap instanceof Map) {
          const newMap = new Map(typeMap);
          newMap.set(id, count);
          return { ...prev, [type]: newMap };
        }
        return prev;
      });
    },
    [getCacheKey]
  );

  // Fetch count for a single entity
  const fetchCount = useCallback(
    async (
      type: keyof FacetCounts,
      id: string
    ): Promise<number> => {
      const baseFilter = filter.makeFilter();

      try {
        // Build type-specific filter based on filter mode
        switch (filter.mode) {
          case GQL.FilterMode.Scenes: {
            const sceneFilter = baseFilter as GQL.SceneFilterType;
            let modifiedFilter: GQL.SceneFilterType;

            switch (type) {
              case "tags":
                modifiedFilter = {
                  ...sceneFilter,
                  tags: {
                    value: [id],
                    modifier: GQL.CriterionModifier.Includes,
                  },
                };
                break;
              case "performers":
                modifiedFilter = {
                  ...sceneFilter,
                  performers: {
                    value: [id],
                    modifier: GQL.CriterionModifier.Includes,
                  },
                };
                break;
              case "studios":
                modifiedFilter = {
                  ...sceneFilter,
                  studios: {
                    value: [id],
                    modifier: GQL.CriterionModifier.Includes,
                    depth: -1,
                  },
                };
                break;
              default:
                return 0;
            }

            const client = getClient();
            const result = await client.query<GQL.FindScenesQuery>({
              query: GQL.FindScenesDocument,
              variables: {
                filter: { per_page: 0 },
                scene_filter: modifiedFilter,
              },
              fetchPolicy: "network-only",
            });

            return result.data.findScenes.count;
          }

          // Add other filter modes as needed
          default:
            return 0;
        }
      } catch (error) {
        console.error(`Error fetching count for ${type}:${id}`, error);
        return 0;
      }
    },
    [filter]
  );

  // Process pending requests in batches
  const processPendingRequests = useCallback(async () => {
    const requests = [...pendingRef.current];
    pendingRef.current = [];

    if (requests.length === 0 || !isOpen) return;

    setLoading(true);

    try {
      // Process in batches with concurrency limit
      for (let i = 0; i < requests.length; i += BATCH_SIZE) {
        const batch = requests.slice(i, i + BATCH_SIZE);

        // Execute batch with concurrency limit
        const promises = batch.map(async (req) => {
          const count = await fetchCount(req.type, req.id);
          setCount(req.type, req.id, count);
          req.resolve(count);
          return count;
        });

        // Process with concurrency limit
        for (let j = 0; j < promises.length; j += MAX_CONCURRENT) {
          await Promise.all(promises.slice(j, j + MAX_CONCURRENT));
        }

        // Small delay between batches
        if (i + BATCH_SIZE < requests.length) {
          await new Promise((resolve) => setTimeout(resolve, 100));
        }
      }
    } finally {
      setLoading(false);
    }
  }, [isOpen, fetchCount, setCount]);

  // Request a count (debounced)
  const requestCount = useCallback(
    (type: keyof FacetCounts, id: string) => {
      // Check cache first
      const cached = getCount(type, id);
      if (cached !== undefined) {
        return;
      }

      // Add to pending queue
      return new Promise<number>((resolve) => {
        // Avoid duplicates
        const existing = pendingRef.current.find(
          (r) => r.type === type && r.id === id
        );
        if (!existing) {
          pendingRef.current.push({ type, id, resolve });
        }

        // Debounce processing
        if (debounceRef.current) {
          clearTimeout(debounceRef.current);
        }
        debounceRef.current = setTimeout(() => {
          processPendingRequests();
        }, DEBOUNCE_MS);
      });
    },
    [getCount, processPendingRequests]
  );

  // Clear cache
  const clearCache = useCallback(() => {
    cacheRef.current.clear();
    setCounts({
      tags: new Map(),
      performers: new Map(),
      studios: new Map(),
      groups: new Map(),
      galleries: new Map(),
      resolutions: new Map(),
      orientations: new Map(),
      genders: new Map(),
      organized: null,
    });
  }, []);

  // Set external counts (from future facets endpoint)
  const setExternalCounts = useCallback((external: Partial<FacetCounts>) => {
    setCounts((prev) => {
      const updated = { ...prev };
      for (const [key, value] of Object.entries(external)) {
        if (value instanceof Map) {
          updated[key as keyof FacetCounts] = value as any;
          // Also populate cache
          for (const [id, count] of value.entries()) {
            const cacheKey = `${filterFingerprint}:${key}:${id}`;
            cacheRef.current.set(cacheKey, { count, timestamp: Date.now() });
          }
        }
      }
      return updated;
    });
  }, [filterFingerprint]);

  // Clear pending on unmount
  useEffect(() => {
    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, []);

  const value = useMemo(
    () => ({
      counts,
      loading,
      requestCount,
      getCount,
      clearCache,
      setExternalCounts,
    }),
    [counts, loading, requestCount, getCount, clearCache, setExternalCounts]
  );

  return (
    <FacetsContext.Provider value={value}>{children}</FacetsContext.Provider>
  );
}

/**
 * Hook to access the facets context
 */
export function useFacets() {
  const context = useContext(FacetsContext);
  if (!context) {
    throw new Error("useFacets must be used within a FacetsProvider");
  }
  return context;
}

/**
 * Hook to get facet counts for a specific filter type.
 * Can be used by individual filter components.
 */
export function useFilterFacets(type: keyof FacetCounts) {
  const { counts, loading, requestCount, getCount } = useFacets();

  return {
    counts: counts[type],
    loading,
    requestCount: (id: string) => requestCount(type, id),
    getCount: (id: string) => getCount(type, id),
  };
}

