import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { ListFilterModel } from "src/models/list-filter/filter";
import * as GQL from "src/core/generated-graphql";

/**
 * Labeled facet count with id, label, and count
 */
export interface LabeledFacetCount {
  count: number;
  label: string;
}

/**
 * Facet count data organized by filter type
 */
export interface FacetCounts {
  tags: Map<string, LabeledFacetCount>;
  performers: Map<string, LabeledFacetCount>;
  studios: Map<string, LabeledFacetCount>;
  groups: Map<string, LabeledFacetCount>;
  performerTags: Map<string, LabeledFacetCount>;
  resolutions: Map<GQL.ResolutionEnum, number>;
  orientations: Map<GQL.OrientationEnum, number>;
  genders: Map<GQL.GenderEnum, number>;
  countries: Map<string, LabeledFacetCount>;
  circumcised: Map<GQL.CircumisedEnum, number>;
  ratings: Map<number, number>;
  captions: Map<string, number>;
  booleans: {
    organized: { true: number; false: number };
    interactive: { true: number; false: number };
    favorite: { true: number; false: number };
  };
  parents: Map<string, LabeledFacetCount>;
  children: Map<string, LabeledFacetCount>;
}

const EMPTY_COUNTS: FacetCounts = {
  tags: new Map<string, LabeledFacetCount>(),
  performers: new Map<string, LabeledFacetCount>(),
  studios: new Map<string, LabeledFacetCount>(),
  groups: new Map<string, LabeledFacetCount>(),
  performerTags: new Map<string, LabeledFacetCount>(),
  resolutions: new Map<GQL.ResolutionEnum, number>(),
  orientations: new Map<GQL.OrientationEnum, number>(),
  genders: new Map<GQL.GenderEnum, number>(),
  countries: new Map<string, LabeledFacetCount>(),
  circumcised: new Map<GQL.CircumisedEnum, number>(),
  ratings: new Map<number, number>(),
  captions: new Map<string, number>(),
  booleans: {
    organized: { true: 0, false: 0 },
    interactive: { true: 0, false: 0 },
    favorite: { true: 0, false: 0 },
  },
  parents: new Map<string, LabeledFacetCount>(),
  children: new Map<string, LabeledFacetCount>(),
};

/**
 * Convert facet count array to Map with labels
 */
function toMap(counts: { id: string; label: string; count: number }[]): Map<string, LabeledFacetCount> {
  return new Map(counts.map((c) => [c.id, { count: c.count, label: c.label }]));
}

/**
 * Convert resolution facet counts to Map
 */
function toResolutionMap(
  counts: { resolution: GQL.ResolutionEnum; count: number }[]
): Map<GQL.ResolutionEnum, number> {
  return new Map(counts.map((c) => [c.resolution, c.count]));
}

/**
 * Convert orientation facet counts to Map
 */
function toOrientationMap(
  counts: { orientation: GQL.OrientationEnum; count: number }[]
): Map<GQL.OrientationEnum, number> {
  return new Map(counts.map((c) => [c.orientation, c.count]));
}

/**
 * Convert gender facet counts to Map
 */
function toGenderMap(
  counts: { gender: GQL.GenderEnum; count: number }[]
): Map<GQL.GenderEnum, number> {
  return new Map(counts.map((c) => [c.gender, c.count]));
}

/**
 * Convert boolean facet counts to object
 */
function toBooleanCounts(
  counts: { value: boolean; count: number }[]
): { true: number; false: number } {
  const result = { true: 0, false: 0 };
  for (const c of counts) {
    if (c.value) {
      result.true = c.count;
    } else {
      result.false = c.count;
    }
  }
  return result;
}

/**
 * Convert circumcised facet counts to Map
 */
function toCircumcisedMap(
  counts: { value: GQL.CircumisedEnum; count: number }[]
): Map<GQL.CircumisedEnum, number> {
  return new Map(counts.map((c) => [c.value, c.count]));
}

/**
 * Convert rating facet counts to Map
 */
function toRatingMap(
  counts: { rating: number; count: number }[]
): Map<number, number> {
  return new Map(counts.map((c) => [c.rating, c.count]));
}

/**
 * Convert caption facet counts to Map
 */
function toCaptionMap(
  counts: { language: string; count: number }[]
): Map<string, number> {
  return new Map(counts.map((c) => [c.language, c.count]));
}

interface UseFacetCountsOptions {
  /** Only fetch when sidebar is open */
  isOpen?: boolean;
  /** Debounce delay in ms (default 500) */
  debounceMs?: number;
  /** Limit number of facets per category (default 100) */
  limit?: number;
  /** Include performer_tags facet (expensive, lazy-loaded) */
  includePerformerTags?: boolean;
  /** Include captions facet (expensive, lazy-loaded) */
  includeCaptions?: boolean;
}

/**
 * Hook that fetches aggregated facet counts using the facets endpoint.
 * This is much more efficient than making individual count queries.
 * 
 * Expensive facets (performer_tags, captions) are lazy-loaded - only fetched
 * when their respective filter sections are expanded.
 */
export function useSceneFacetCounts(
  filter: ListFilterModel,
  options: UseFacetCountsOptions = {}
) {
  const { 
    isOpen = true, 
    debounceMs = 500, 
    limit = 100,
    includePerformerTags = false,
    includeCaptions = false,
  } = options;

  const [counts, setCounts] = useState<FacetCounts>(EMPTY_COUNTS);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastFilterRef = useRef<string>("");
  // Track last options to detect when expensive facets should be fetched
  const lastOptionsRef = useRef({ includePerformerTags: false, includeCaptions: false });

  const [fetchFacets] = GQL.useSceneFacetsLazyQuery({
    fetchPolicy: "network-only",
  });

  const filterFingerprint = useMemo(() => {
    return JSON.stringify(filter.makeFilter());
  }, [filter]);

  const doFetch = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);
    try {
      const result = await fetchFacets({
        variables: {
          scene_filter: filter.makeFilter() as GQL.SceneFilterType,
          limit,
          include_performer_tags: includePerformerTags,
          include_captions: includeCaptions,
        },
      });

      if (result.data?.sceneFacets) {
        const facets = result.data.sceneFacets;
        
        // Check if this is a lazy-load update (only updating expensive facets)
        const isLazyLoadUpdate = 
          (includePerformerTags && !lastOptionsRef.current.includePerformerTags) ||
          (includeCaptions && !lastOptionsRef.current.includeCaptions);
        
        if (isLazyLoadUpdate) {
          // Only update the lazy-loaded facets, preserve everything else
          // This prevents React rendering glitches during partial updates
          setCounts((prev) => ({
            ...prev,
            performerTags: includePerformerTags ? toMap(facets.performer_tags ?? []) : prev.performerTags,
            captions: includeCaptions ? toCaptionMap(facets.captions ?? []) : prev.captions,
          }));
        } else {
          // Full update - update all facets
          setCounts((prev) => ({
            tags: toMap(facets.tags),
            performers: toMap(facets.performers),
            studios: toMap(facets.studios),
            groups: toMap(facets.groups),
            performerTags: includePerformerTags ? toMap(facets.performer_tags ?? []) : prev.performerTags,
            resolutions: toResolutionMap(facets.resolutions),
            orientations: toOrientationMap(facets.orientations),
            genders: new Map(),
            countries: new Map(),
            circumcised: new Map(),
            ratings: toRatingMap(facets.ratings),
            captions: includeCaptions ? toCaptionMap(facets.captions ?? []) : prev.captions,
            booleans: {
              organized: toBooleanCounts(facets.organized),
              interactive: toBooleanCounts(facets.interactive),
              favorite: { true: 0, false: 0 },
            },
            parents: new Map(),
            children: new Map(),
          }));
        }
        
        // Update last options
        lastOptionsRef.current = { includePerformerTags, includeCaptions };
      }
    } catch (error) {
      console.error("Error fetching scene facets:", error);
    } finally {
      setLoading(false);
    }
  }, [fetchFacets, filter, isOpen, limit, includePerformerTags, includeCaptions]);

  // Fetch when filter changes, sidebar opens, or expensive facets are requested
  useEffect(() => {
    if (!isOpen) return;

    // Check if this is the first fetch or filter changed
    const isFirstFetch = lastFilterRef.current === "";
    const filterChanged = filterFingerprint !== lastFilterRef.current;
    
    // Check if new expensive facets are being requested
    const newExpensiveFacetsRequested = 
      (includePerformerTags && !lastOptionsRef.current.includePerformerTags) ||
      (includeCaptions && !lastOptionsRef.current.includeCaptions);

    if (!isFirstFetch && !filterChanged && !newExpensiveFacetsRequested) return;

    // Mark as loading
    setLoading(true);

    // Clear any pending debounce
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    // First fetch or new expensive facets requested = immediate, filter changes are debounced
    if (isFirstFetch || newExpensiveFacetsRequested) {
      lastFilterRef.current = filterFingerprint;
      doFetch();
    } else {
      debounceRef.current = setTimeout(() => {
        lastFilterRef.current = filterFingerprint;
        doFetch();
      }, debounceMs);
    }

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [filterFingerprint, isOpen, debounceMs, doFetch, includePerformerTags, includeCaptions]);

  return { counts, loading, refetch: doFetch };
}

/**
 * Hook for performer facet counts
 */
export function usePerformerFacetCounts(
  filter: ListFilterModel,
  options: UseFacetCountsOptions = {}
) {
  const { isOpen = true, debounceMs = 500, limit = 100 } = options;

  const [counts, setCounts] = useState<FacetCounts>(EMPTY_COUNTS);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastFilterRef = useRef<string>("");

  const [fetchFacets] = GQL.usePerformerFacetsLazyQuery({
    fetchPolicy: "network-only",
  });

  const filterFingerprint = useMemo(() => {
    return JSON.stringify(filter.makeFilter());
  }, [filter]);

  const doFetch = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);
    try {
      const result = await fetchFacets({
        variables: {
          performer_filter: filter.makeFilter() as GQL.PerformerFilterType,
          limit,
        },
      });

      if (result.data?.performerFacets) {
        const facets = result.data.performerFacets;
        setCounts({
          ...EMPTY_COUNTS,
          tags: toMap(facets.tags),
          studios: toMap(facets.studios),
          genders: toGenderMap(facets.genders),
          countries: toMap(facets.countries),
          circumcised: toCircumcisedMap(facets.circumcised),
          ratings: toRatingMap(facets.ratings),
          booleans: {
            ...EMPTY_COUNTS.booleans,
            favorite: toBooleanCounts(facets.favorite),
          },
        });
      }
    } catch (error) {
      console.error("Error fetching performer facets:", error);
    } finally {
      setLoading(false);
    }
  }, [fetchFacets, filter, isOpen, limit]);

  useEffect(() => {
    if (!isOpen) return;

    const isFirstFetch = lastFilterRef.current === "";
    const filterChanged = filterFingerprint !== lastFilterRef.current;

    if (!isFirstFetch && !filterChanged) return;

    setLoading(true);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    if (isFirstFetch) {
      lastFilterRef.current = filterFingerprint;
      doFetch();
    } else {
      debounceRef.current = setTimeout(() => {
        lastFilterRef.current = filterFingerprint;
        doFetch();
      }, debounceMs);
    }

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [filterFingerprint, isOpen, debounceMs, doFetch]);

  return { counts, loading, refetch: doFetch };
}

/**
 * Hook for gallery facet counts
 */
export function useGalleryFacetCounts(
  filter: ListFilterModel,
  options: UseFacetCountsOptions = {}
) {
  const { isOpen = true, debounceMs = 500, limit = 100 } = options;

  const [counts, setCounts] = useState<FacetCounts>(EMPTY_COUNTS);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastFilterRef = useRef<string>("");

  const [fetchFacets] = GQL.useGalleryFacetsLazyQuery({
    fetchPolicy: "network-only",
  });

  const filterFingerprint = useMemo(() => {
    return JSON.stringify(filter.makeFilter());
  }, [filter]);

  const doFetch = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);
    try {
      const result = await fetchFacets({
        variables: {
          gallery_filter: filter.makeFilter() as GQL.GalleryFilterType,
          limit,
        },
      });

      if (result.data?.galleryFacets) {
        const facets = result.data.galleryFacets;
        setCounts({
          ...EMPTY_COUNTS,
          tags: toMap(facets.tags),
          performers: toMap(facets.performers),
          studios: toMap(facets.studios),
          ratings: toRatingMap(facets.ratings),
          booleans: {
            ...EMPTY_COUNTS.booleans,
            organized: toBooleanCounts(facets.organized),
          },
        });
      }
    } catch (error) {
      console.error("Error fetching gallery facets:", error);
    } finally {
      setLoading(false);
    }
  }, [fetchFacets, filter, isOpen, limit]);

  useEffect(() => {
    if (!isOpen) return;

    const isFirstFetch = lastFilterRef.current === "";
    const filterChanged = filterFingerprint !== lastFilterRef.current;

    if (!isFirstFetch && !filterChanged) return;

    setLoading(true);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    if (isFirstFetch) {
      lastFilterRef.current = filterFingerprint;
      doFetch();
    } else {
      debounceRef.current = setTimeout(() => {
        lastFilterRef.current = filterFingerprint;
        doFetch();
      }, debounceMs);
    }

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [filterFingerprint, isOpen, debounceMs, doFetch]);

  return { counts, loading, refetch: doFetch };
}

/**
 * Hook for group facet counts
 */
export function useGroupFacetCounts(
  filter: ListFilterModel,
  options: UseFacetCountsOptions = {}
) {
  const { isOpen = true, debounceMs = 500, limit = 100 } = options;

  const [counts, setCounts] = useState<FacetCounts>(EMPTY_COUNTS);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastFilterRef = useRef<string>("");

  const [fetchFacets] = GQL.useGroupFacetsLazyQuery({
    fetchPolicy: "network-only",
  });

  const filterFingerprint = useMemo(() => {
    return JSON.stringify(filter.makeFilter());
  }, [filter]);

  const doFetch = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);
    try {
      const result = await fetchFacets({
        variables: {
          group_filter: filter.makeFilter() as GQL.GroupFilterType,
          limit,
        },
      });

      if (result.data?.groupFacets) {
        const facets = result.data.groupFacets;
        setCounts({
          ...EMPTY_COUNTS,
          tags: toMap(facets.tags),
          performers: toMap(facets.performers),
          studios: toMap(facets.studios),
        });
      }
    } catch (error) {
      console.error("Error fetching group facets:", error);
    } finally {
      setLoading(false);
    }
  }, [fetchFacets, filter, isOpen, limit]);

  useEffect(() => {
    if (!isOpen) return;

    const isFirstFetch = lastFilterRef.current === "";
    const filterChanged = filterFingerprint !== lastFilterRef.current;

    if (!isFirstFetch && !filterChanged) return;

    setLoading(true);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    if (isFirstFetch) {
      lastFilterRef.current = filterFingerprint;
      doFetch();
    } else {
      debounceRef.current = setTimeout(() => {
        lastFilterRef.current = filterFingerprint;
        doFetch();
      }, debounceMs);
    }

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [filterFingerprint, isOpen, debounceMs, doFetch]);

  return { counts, loading, refetch: doFetch };
}

/**
 * Hook for studio facet counts
 */
export function useStudioFacetCounts(
  filter: ListFilterModel,
  options: UseFacetCountsOptions = {}
) {
  const { isOpen = true, debounceMs = 500, limit = 100 } = options;

  const [counts, setCounts] = useState<FacetCounts>(EMPTY_COUNTS);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastFilterRef = useRef<string>("");

  const [fetchFacets] = GQL.useStudioFacetsLazyQuery({
    fetchPolicy: "network-only",
  });

  const filterFingerprint = useMemo(() => {
    return JSON.stringify(filter.makeFilter());
  }, [filter]);

  const doFetch = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);
    try {
      const result = await fetchFacets({
        variables: {
          studio_filter: filter.makeFilter() as GQL.StudioFilterType,
          limit,
        },
      });

      if (result.data?.studioFacets) {
        const facets = result.data.studioFacets;
        setCounts({
          ...EMPTY_COUNTS,
          tags: toMap(facets.tags),
          parents: toMap(facets.parents),
          booleans: {
            ...EMPTY_COUNTS.booleans,
            favorite: toBooleanCounts(facets.favorite),
          },
        });
      }
    } catch (error) {
      console.error("Error fetching studio facets:", error);
    } finally {
      setLoading(false);
    }
  }, [fetchFacets, filter, isOpen, limit]);

  useEffect(() => {
    if (!isOpen) return;

    const isFirstFetch = lastFilterRef.current === "";
    const filterChanged = filterFingerprint !== lastFilterRef.current;

    if (!isFirstFetch && !filterChanged) return;

    setLoading(true);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    if (isFirstFetch) {
      lastFilterRef.current = filterFingerprint;
      doFetch();
    } else {
      debounceRef.current = setTimeout(() => {
        lastFilterRef.current = filterFingerprint;
        doFetch();
      }, debounceMs);
    }

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [filterFingerprint, isOpen, debounceMs, doFetch]);

  return { counts, loading, refetch: doFetch };
}

/**
 * Hook for tag facet counts
 */
export function useTagFacetCounts(
  filter: ListFilterModel,
  options: UseFacetCountsOptions = {}
) {
  const { isOpen = true, debounceMs = 500, limit = 100 } = options;

  const [counts, setCounts] = useState<FacetCounts>(EMPTY_COUNTS);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastFilterRef = useRef<string>("");

  const [fetchFacets] = GQL.useTagFacetsLazyQuery({
    fetchPolicy: "network-only",
  });

  const filterFingerprint = useMemo(() => {
    return JSON.stringify(filter.makeFilter());
  }, [filter]);

  const doFetch = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);
    try {
      const result = await fetchFacets({
        variables: {
          tag_filter: filter.makeFilter() as GQL.TagFilterType,
          limit,
        },
      });

      if (result.data?.tagFacets) {
        const facets = result.data.tagFacets;
        setCounts({
          ...EMPTY_COUNTS,
          parents: toMap(facets.parents),
          children: toMap(facets.children),
          booleans: {
            ...EMPTY_COUNTS.booleans,
            favorite: toBooleanCounts(facets.favorite),
          },
        });
      }
    } catch (error) {
      console.error("Error fetching tag facets:", error);
    } finally {
      setLoading(false);
    }
  }, [fetchFacets, filter, isOpen, limit]);

  useEffect(() => {
    if (!isOpen) return;

    const isFirstFetch = lastFilterRef.current === "";
    const filterChanged = filterFingerprint !== lastFilterRef.current;

    if (!isFirstFetch && !filterChanged) return;

    setLoading(true);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    if (isFirstFetch) {
      lastFilterRef.current = filterFingerprint;
      doFetch();
    } else {
      debounceRef.current = setTimeout(() => {
        lastFilterRef.current = filterFingerprint;
        doFetch();
      }, debounceMs);
    }

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [filterFingerprint, isOpen, debounceMs, doFetch]);

  return { counts, loading, refetch: doFetch };
}

/**
 * Universal hook that selects the appropriate facet hook based on filter mode
 */
export function useFacetCounts(
  filter: ListFilterModel,
  options: UseFacetCountsOptions = {}
) {
  const mode = filter.mode;

  // We call all hooks but only one will be active
  const sceneCounts = useSceneFacetCounts(filter, {
    ...options,
    isOpen: options.isOpen && mode === GQL.FilterMode.Scenes,
  });
  const performerCounts = usePerformerFacetCounts(filter, {
    ...options,
    isOpen: options.isOpen && mode === GQL.FilterMode.Performers,
  });
  const galleryCounts = useGalleryFacetCounts(filter, {
    ...options,
    isOpen: options.isOpen && mode === GQL.FilterMode.Galleries,
  });
  const groupCounts = useGroupFacetCounts(filter, {
    ...options,
    isOpen: options.isOpen && mode === GQL.FilterMode.Groups,
  });
  const studioCounts = useStudioFacetCounts(filter, {
    ...options,
    isOpen: options.isOpen && mode === GQL.FilterMode.Studios,
  });
  const tagCounts = useTagFacetCounts(filter, {
    ...options,
    isOpen: options.isOpen && mode === GQL.FilterMode.Tags,
  });

  switch (mode) {
    case GQL.FilterMode.Scenes:
      return sceneCounts;
    case GQL.FilterMode.Performers:
      return performerCounts;
    case GQL.FilterMode.Galleries:
      return galleryCounts;
    case GQL.FilterMode.Groups:
      return groupCounts;
    case GQL.FilterMode.Studios:
      return studioCounts;
    case GQL.FilterMode.Tags:
      return tagCounts;
    default:
      return sceneCounts;
  }
}

/**
 * Context for sharing facet counts with child filter components
 */
export const FacetCountsContext = React.createContext<{
  counts: FacetCounts;
  loading: boolean;
}>({
  counts: EMPTY_COUNTS,
  loading: false,
});

/**
 * Hook for filter components to access facet counts from context
 */
export function useFacetCountsContext() {
  return useContext(FacetCountsContext);
}
