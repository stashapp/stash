import { LabeledFacetCount } from "src/hooks/useFacetCounts";
import { Option } from "./SidebarListFilter";

export interface FacetCandidateOptions {
  /** Candidates from the search query */
  searchCandidates: Option[];
  /** Selected item IDs to exclude */
  selectedIds: Set<string>;
  /** Current search query */
  searchQuery: string;
  /** Facet counts map (id -> { count, label }) */
  facetCounts: Map<string, LabeledFacetCount>;
  /** Whether facets are currently loading */
  facetsLoading: boolean;
}

/**
 * Builds the candidate list for entity filters (performers, tags, studios, etc.)
 * 
 * Strategy:
 * - When facets are loaded and no search query, USE facet results as candidates
 *   (since facets returns the TOP N most relevant items by count)
 * - When user searches, use search results merged with facet counts
 *   (to allow finding specific items not in top N)
 * 
 * This prevents the mismatch bug where search returns different items than facets,
 * causing items to be incorrectly filtered out when counts are applied.
 */
export function buildFacetCandidates(options: FacetCandidateOptions): Option[] {
  const {
    searchCandidates,
    selectedIds,
    searchQuery,
    facetCounts,
    facetsLoading,
  } = options;

  const hasValidFacets = facetCounts.size > 0 && !facetsLoading;
  const hasSearchQuery = searchQuery && searchQuery.length > 0;

  // Extract modifier options (Any, None, etc.)
  const modifierOptions = searchCandidates.filter(
    (c) => c.className === "modifier-object"
  );

  if (hasValidFacets && !hasSearchQuery) {
    // No search query: Use facet results directly as candidates (with labels from facets)
    const facetCandidates: Option[] = [];
    
    facetCounts.forEach((facetData, id) => {
      // Skip if already selected or has zero count
      if (selectedIds.has(id) || facetData.count === 0) return;

      facetCandidates.push({
        id,
        label: facetData.label,
        count: facetData.count,
      });
    });

    // Sort by count descending
    facetCandidates.sort((a, b) => (b.count ?? 0) - (a.count ?? 0));

    return [...modifierOptions, ...facetCandidates];
  } else {
    // With search query OR facets not loaded: Use search results with counts merged
    return searchCandidates
      .map((c) => {
        if (c.className === "modifier-object") return c;
        const facetData = hasValidFacets ? facetCounts.get(c.id) : undefined;
        return { ...c, count: facetData?.count };
      })
      .filter((c) => {
        if (c.className === "modifier-object") return true;
        if (!hasValidFacets) return true;
        // Filter out zero counts, keep undefined (not in top N) and positive
        return c.count !== 0;
      });
  }
}

/**
 * Builds candidates for enum filters (resolution, orientation, gender, etc.)
 * 
 * For enum filters, facets returns counts for ALL possible values (finite set),
 * so we filter out both undefined AND zero counts.
 */
export function buildEnumCandidates(
  options: Option[],
  selectedIds: Set<string>,
  counts: Map<string, number> | undefined,
  countsLoading: boolean
): Option[] {
  const hasValidCounts = counts && counts.size > 0 && !countsLoading;

  return options.filter((opt) => {
    // Skip already selected
    if (selectedIds.has(opt.id)) return false;
    // If counts not loaded, show all
    if (!hasValidCounts) return true;
    // Filter out undefined and zero counts
    const count = counts?.get(opt.id);
    return count !== undefined && count > 0;
  });
}

