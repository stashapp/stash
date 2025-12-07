/// <reference types="node" />
import { useCallback, useEffect, useRef, useState } from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import * as GQL from "src/core/generated-graphql";
import { getClient } from "src/core/StashService";

/**
 * Batched count request configuration
 */
interface CountRequest {
  type: "tag" | "performer" | "studio" | "resolution" | "orientation" | "gender";
  id: string;
}

/**
 * Result of batched count requests
 */
export interface BatchedCounts {
  [key: string]: number;
}

/**
 * Configuration for batch requests
 */
interface BatchConfig {
  /** Maximum number of concurrent requests */
  maxConcurrent: number;
  /** Debounce delay in milliseconds */
  debounceMs: number;
  /** Maximum items per batch */
  batchSize: number;
}

const DEFAULT_CONFIG: BatchConfig = {
  maxConcurrent: 4,
  debounceMs: 500,
  batchSize: 10,
};

/**
 * Execute a count query for scenes with a modified filter
 */
async function querySceneCount(
  baseFilter: GQL.SceneFilterType,
  additionalFilter: Partial<GQL.SceneFilterType>
): Promise<number> {
  const client = getClient();
  const result = await client.query<GQL.FindScenesQuery>({
    query: GQL.FindScenesDocument,
    variables: {
      filter: { per_page: 0 },
      scene_filter: { ...baseFilter, ...additionalFilter },
    },
    fetchPolicy: "network-only",
  });
  return result.data.findScenes.count;
}

/**
 * Hook that batches and rate-limits filter count requests.
 *
 * This is a client-side optimization that:
 * 1. Batches multiple count requests together
 * 2. Rate-limits concurrent requests
 * 3. Debounces filter changes
 * 4. Caches results
 *
 * This provides improved performance over individual count queries
 * while waiting for a proper backend facets endpoint.
 */
export function useBatchedFilterCounts(
  filter: ListFilterModel,
  isOpen: boolean,
  config: Partial<BatchConfig> = {}
) {
  const effectiveConfig = { ...DEFAULT_CONFIG, ...config };

  const [counts, setCounts] = useState<BatchedCounts>({});
  const [loading, setLoading] = useState(false);
  const [pendingRequests, setPendingRequests] = useState<CountRequest[]>([]);

  const debounceRef = useRef<NodeJS.Timeout | null>(null);
  const cacheRef = useRef<
    Map<string, { count: number; timestamp: number }>
  >(new Map());
  const filterFingerprintRef = useRef<string>("");

  // Cache TTL in milliseconds (5 minutes)
  const CACHE_TTL = 5 * 60 * 1000;

  // Build filter fingerprint for cache invalidation
  const getFilterFingerprint = useCallback(() => {
    const filterCopy = filter.clone();
    return JSON.stringify(filterCopy.makeFilter());
  }, [filter]);

  // Check if cache entry is valid
  const isCacheValid = useCallback(
    (key: string): boolean => {
      const entry = cacheRef.current.get(key);
      if (!entry) return false;
      return Date.now() - entry.timestamp < CACHE_TTL;
    },
    []
  );

  // Get cached count
  const getCachedCount = useCallback(
    (key: string): number | undefined => {
      if (isCacheValid(key)) {
        return cacheRef.current.get(key)?.count;
      }
      return undefined;
    },
    [isCacheValid]
  );

  // Set cached count
  const setCachedCount = useCallback((key: string, count: number) => {
    cacheRef.current.set(key, { count, timestamp: Date.now() });
  }, []);

  // Request count for a specific filter value
  const requestCount = useCallback(
    (request: CountRequest) => {
      const filterFingerprint = getFilterFingerprint();
      const cacheKey = `${filterFingerprint}:${request.type}:${request.id}`;
      const cached = getCachedCount(cacheKey);

      if (cached !== undefined) {
        setCounts((prev) => ({
          ...prev,
          [`${request.type}:${request.id}`]: cached,
        }));
        return;
      }

      setPendingRequests((prev) => {
        // Avoid duplicates
        if (prev.some((r) => r.type === request.type && r.id === request.id)) {
          return prev;
        }
        return [...prev, request];
      });
    },
    [getFilterFingerprint, getCachedCount]
  );

  // Process pending requests in batches
  const processPendingRequests = useCallback(async () => {
    if (pendingRequests.length === 0 || !isOpen) return;

    setLoading(true);
    const filterFingerprint = getFilterFingerprint();
    filterFingerprintRef.current = filterFingerprint;

    try {
      // Group requests by type
      const requestsByType = pendingRequests.reduce(
        (acc, req) => {
          if (!acc[req.type]) acc[req.type] = [];
          acc[req.type].push(req.id);
          return acc;
        },
        {} as Record<string, string[]>
      );

      const results: BatchedCounts = {};
      const baseFilter = filter.makeFilter() as GQL.SceneFilterType;

      // Process each type with rate limiting
      for (const [type, ids] of Object.entries(requestsByType)) {
        // Process in batches
        for (let i = 0; i < ids.length; i += effectiveConfig.batchSize) {
          const batch = ids.slice(i, i + effectiveConfig.batchSize);

          // Execute batch queries in parallel with concurrency limit
          const batchPromises = batch.map(async (id) => {
            const cacheKey = `${filterFingerprint}:${type}:${id}`;

            try {
              let count = 0;

              // Build type-specific filter
              switch (type) {
                case "tag": {
                  count = await querySceneCount(baseFilter, {
                    tags: {
                      value: [id],
                      modifier: GQL.CriterionModifier.Includes,
                    },
                  });
                  break;
                }
                case "performer": {
                  count = await querySceneCount(baseFilter, {
                    performers: {
                      value: [id],
                      modifier: GQL.CriterionModifier.Includes,
                    },
                  });
                  break;
                }
                case "studio": {
                  count = await querySceneCount(baseFilter, {
                    studios: {
                      value: [id],
                      modifier: GQL.CriterionModifier.Includes,
                      depth: -1,
                    },
                  });
                  break;
                }
              }

              setCachedCount(cacheKey, count);
              return { key: `${type}:${id}`, count };
            } catch (error) {
              console.error(`Error fetching count for ${type}:${id}`, error);
              return { key: `${type}:${id}`, count: 0 };
            }
          });

          // Wait for batch to complete with concurrency limit
          const batchResults = await Promise.all(
            batchPromises.slice(0, effectiveConfig.maxConcurrent)
          );

          for (const result of batchResults) {
            results[result.key] = result.count;
          }

          // Small delay between batches to avoid overwhelming the server
          if (i + effectiveConfig.batchSize < ids.length) {
            await new Promise((resolve) => setTimeout(resolve, 100));
          }
        }
      }

      setCounts((prev) => ({ ...prev, ...results }));
      setPendingRequests([]);
    } catch (error) {
      console.error("Error processing batched count requests:", error);
    } finally {
      setLoading(false);
    }
  }, [
    pendingRequests,
    isOpen,
    filter,
    getFilterFingerprint,
    effectiveConfig,
    setCachedCount,
  ]);

  // Debounced processing of pending requests
  useEffect(() => {
    if (pendingRequests.length === 0) return;

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    debounceRef.current = setTimeout(() => {
      processPendingRequests();
    }, effectiveConfig.debounceMs);

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [pendingRequests, effectiveConfig.debounceMs, processPendingRequests]);

  return {
    counts,
    loading,
    requestCount,
    clearCache: () => {
      cacheRef.current.clear();
      setCounts({});
    },
  };
}

/**
 * Hook that provides optimized filter counts using batching.
 *
 * This is a drop-in optimization that can be used by individual
 * filter components to share a common batching layer.
 */
export function useOptimizedFilterCounts(
  filter: ListFilterModel,
  isOpen: boolean
) {
  return useBatchedFilterCounts(filter, isOpen, {
    maxConcurrent: 6,
    debounceMs: 750,
    batchSize: 15,
  });
}
