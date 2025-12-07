import React, { ReactNode, useCallback, useContext, useMemo } from "react";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import {
  CriterionModifier,
  FilterPerformerDataFragment,
  FindPerformersForFilterQueryVariables,
  PerformerFilterType,
  useFindPerformersForFilterQuery,
} from "src/core/generated-graphql";
import { FacetCountsContext } from "src/extensions/hooks/useFacetCounts";
import { ObjectsFilter } from "./SelectableFilter";
import { sortByRelevance } from "src/utils/query";
import { ListFilterModel } from "src/models/list-filter/filter";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import {
  IUseQueryHookProps,
  makeQueryVariables,
  setObjectFilter,
  useLabeledIdFilterState,
} from "./LabeledIdFilter";
import { Option, SidebarListFilter } from "./SidebarListFilter";

interface IPerformersFilter {
  criterion: PerformersCriterion;
  setCriterion: (c: PerformersCriterion) => void;
}

interface IHasModifier {
  modifier: CriterionModifier;
}

function queryVariables(
  query: string,
  f?: ListFilterModel
): FindPerformersForFilterQueryVariables {
  const performerFilter: PerformerFilterType = {};

  if (f) {
    const filterOutput = f.makeFilter();

    // if performer modifier is includes, take it out of the filter
    if (
      (filterOutput.performers as IHasModifier)?.modifier ===
      CriterionModifier.Includes
    ) {
      delete filterOutput.performers;
    }

    setObjectFilter(performerFilter, f.mode, filterOutput);
  }

  return makeQueryVariables(query, { performer_filter: performerFilter });
}

function sortResults(
  query: string,
  performers?: FilterPerformerDataFragment[]
) {
  return sortByRelevance(
    query,
    performers ?? [],
    (p) => p.name,
    (p) => p.alias_list
  ).map((p) => ({
    id: p.id,
    label: p.name,
  }));
}

function usePerformerQueryFilter(props: IUseQueryHookProps) {
  const { q: query, filter: f, skip, filterHook } = props;
  const appliedFilter = filterHook && f ? filterHook(f.clone()) : f;

  const { data, loading } = useFindPerformersForFilterQuery({
    variables: queryVariables(query, appliedFilter),
    skip,
  });

  const results = useMemo(
    () => sortResults(query, data?.findPerformers.performers),
    [data, query]
  );

  return { results, loading };
}

function usePerformerQuery(query: string, skip?: boolean) {
  return usePerformerQueryFilter({ q: query, skip: !!skip });
}

const PerformersFilter: React.FC<IPerformersFilter> = ({
  criterion,
  setCriterion,
}) => {
  return (
    <ObjectsFilter
      criterion={criterion}
      setCriterion={setCriterion}
      useResults={usePerformerQuery}
    />
  );
};

export const SidebarPerformersFilter: React.FC<{
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  filterHook?: (f: ListFilterModel) => ListFilterModel;
  sectionID?: string;
}> = ({ title, option, filter, setFilter, filterHook, sectionID }) => {
  // Get facet counts from context
  const { counts: facetCounts, loading: facetsLoading } = useContext(FacetCountsContext);
  
  const state = useLabeledIdFilterState({
    filter,
    setFilter,
    filterHook,
    option,
    useQuery: usePerformerQueryFilter,
  });

  // Build candidates list
  // Strategy: When facets are loaded and no search query, USE facet results as candidates
  // (since facets returns the TOP N most relevant items by count).
  // When user searches, use search results merged with facet counts.
  const candidatesWithCounts: Option[] = useMemo(() => {
    const hasValidFacets = facetCounts.performers.size > 0 && !facetsLoading;
    const hasSearchQuery = state.query && state.query.length > 0;
    
    // Extract modifier options from candidates
    const modifierOptions = state.candidates.filter(c => c.className === "modifier-object");
    
    // Get selected IDs to exclude from candidates
    const selectedIds = new Set(state.selected.map(s => s.id));
    
    if (hasValidFacets && !hasSearchQuery) {
      // No search query: Use facet results directly as candidates (with labels from facets)
      const facetCandidates: Option[] = [];
      facetCounts.performers.forEach((facetData, id) => {
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
      return state.candidates
        .map((c) => {
          if (c.className === "modifier-object") return c;
          const facetData = hasValidFacets 
            ? facetCounts.performers.get(c.id) 
            : undefined;
          return { ...c, count: facetData?.count };
        })
        .filter((c) => {
          if (c.className === "modifier-object") return true;
          if (!hasValidFacets) return true;
          // Filter out zero counts, keep undefined (not in top N) and positive
          return c.count !== 0;
        });
    }
  }, [state.candidates, state.selected, state.query, facetCounts, facetsLoading]);

  const onOpen = useCallback(() => {
    state.onOpen?.();
  }, [state.onOpen]);

  return (
    <SidebarListFilter
      {...state}
      candidates={candidatesWithCounts}
      title={title}
      sectionID={sectionID}
      onOpen={onOpen}
      loading={state.loading}
      countsLoading={facetsLoading}
    />
  );
};

export default PerformersFilter;
