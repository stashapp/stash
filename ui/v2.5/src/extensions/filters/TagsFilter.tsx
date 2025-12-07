import React, { ReactNode, useCallback, useContext, useMemo, useState } from "react";
import {
  CriterionModifier,
  FilterTagDataFragment,
  TagFilterType,
  useFindTagsForFilterQuery,
} from "src/core/generated-graphql";
import { FacetCountsContext } from "src/extensions/hooks/useFacetCounts";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { sortByRelevance } from "src/utils/query";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  IUseQueryHookProps,
  makeQueryVariables,
  setObjectFilter,
  useLabeledIdFilterState,
} from "./LabeledIdFilter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { TagsCriterion } from "src/models/list-filter/criteria/tags";

interface ITagsFilter {
  criterion: TagsCriterion;
  setCriterion: (c: TagsCriterion) => void;
}

interface IHasModifier {
  modifier: CriterionModifier;
}

function queryVariables(query: string, f?: ListFilterModel) {
  const tagFilter: TagFilterType = {};

  if (f) {
    const filterOutput = f.makeFilter();

    // if tag modifier is includes, take it out of the filter
    if (
      (filterOutput.tags as IHasModifier)?.modifier ===
      CriterionModifier.Includes
    ) {
      delete filterOutput.tags;
    }

    setObjectFilter(tagFilter, f.mode, filterOutput, "tags");
  }

  return makeQueryVariables(query, { tag_filter: tagFilter });
}

function sortResults(
  query: string,
  tags: FilterTagDataFragment[]
) {
  return sortByRelevance(
    query,
    tags,
    (t) => t.name,
    (t) => t.aliases
  ).map((p) => ({
    id: p.id,
    label: p.name,
  }));
}

function useTagQueryFilter(props: IUseQueryHookProps) {
  const { q: query, filter: f, skip, filterHook } = props;
  const appliedFilter = filterHook && f ? filterHook(f.clone()) : f;

  const { data, loading } = useFindTagsForFilterQuery({
    variables: queryVariables(query, appliedFilter),
    skip,
  });

  const results = useMemo(
    () => sortResults(query, data?.findTags.tags ?? []),
    [data, query]
  );

  return { results, loading };
}

function useTagQuery(query: string, skip?: boolean) {
  return useTagQueryFilter({ q: query, skip: !!skip });
}

const TagsFilter: React.FC<ITagsFilter> = ({ criterion, setCriterion }) => {
  return (
    <HierarchicalObjectsFilter
      criterion={criterion}
      setCriterion={setCriterion}
      useResults={useTagQuery}
    />
  );
};

export const SidebarTagsFilter: React.FC<{
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
    useQuery: useTagQueryFilter,
    hierarchical: true,
    includeSubMessageID: "sub_tags",
  });

  // Build candidates list
  // Strategy: When facets are loaded and no search query, USE facet results as candidates
  // (since facets returns the TOP N most relevant items by count).
  // When user searches, use search results merged with facet counts.
  const candidatesWithCounts: Option[] = useMemo(() => {
    const hasValidFacets = facetCounts.tags.size > 0 && !facetsLoading;
    const hasSearchQuery = state.query && state.query.length > 0;
    
    // Extract modifier options from candidates
    const modifierOptions = state.candidates.filter(c => c.className === "modifier-object");
    
    // Get selected IDs to exclude from candidates
    const selectedIds = new Set(state.selected.map(s => s.id));
    
    if (hasValidFacets && !hasSearchQuery) {
      // No search query: Use facet results directly as candidates (with labels from facets)
      const facetCandidates: Option[] = [];
      facetCounts.tags.forEach((facetData, id) => {
        if (selectedIds.has(id) || facetData.count === 0) return;
        facetCandidates.push({
          id,
          label: facetData.label,
          count: facetData.count,
        });
      });
      facetCandidates.sort((a, b) => (b.count ?? 0) - (a.count ?? 0));
      return [...modifierOptions, ...facetCandidates];
    } else {
      // With search query OR facets not loaded: Use search results with counts merged
      return state.candidates
        .map((c) => {
          if (c.className === "modifier-object") return c;
          const facetData = hasValidFacets ? facetCounts.tags.get(c.id) : undefined;
          return { ...c, count: facetData?.count };
        })
        .filter((c) => {
          if (c.className === "modifier-object") return true;
          if (!hasValidFacets) return true;
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

export default TagsFilter;
