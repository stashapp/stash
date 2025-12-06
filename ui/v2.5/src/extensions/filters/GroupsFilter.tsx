import React, { ReactNode, useCallback, useContext, useMemo, useState } from "react";
import {
  CriterionModifier,
  FilterGroupDataFragment,
  GroupFilterType,
  useFindGroupsForFilterQuery,
} from "src/core/generated-graphql";
import { FacetCountsContext } from "src/hooks/useFacetCounts";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { GroupsCriterion } from "src/models/list-filter/criteria/groups";
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
import { ILabeledId } from "src/models/list-filter/types";
import { useCacheResults } from "src/hooks/data";
import { useIntl } from "react-intl";

interface IGroupsFilter {
  criterion: GroupsCriterion;
  setCriterion: (c: GroupsCriterion) => void;
}

function queryVariables(query: string, f?: ListFilterModel) {
  const groupFilter: GroupFilterType = {};

  if (f) {
    const filterOutput = f.makeFilter();
    delete filterOutput.groups;
    setObjectFilter(groupFilter, f.mode, filterOutput, "groups");
  }

  return makeQueryVariables(query, { group_filter: groupFilter });
}

function sortResults(
  query: string,
  groups: FilterGroupDataFragment[]
) {
  return sortByRelevance(
    query,
    groups ?? [],
    (g) => g.name,
    (g) => g.aliases ? [g.aliases] : undefined
  ).map((g) => ({
    id: g.id,
    label: g.name,
  }));
}

function useGroupQueryFilter(props: IUseQueryHookProps) {
  const { q: query, filter: f, skip, filterHook } = props;
  const appliedFilter = filterHook && f ? filterHook(f.clone()) : f;

  const { data, loading } = useFindGroupsForFilterQuery({
    variables: queryVariables(query, appliedFilter),
    skip,
  });

  const results = useMemo(
    () => sortResults(query, data?.findGroups.groups ?? []),
    [data?.findGroups.groups, query]
  );

  return { results, loading };
}

function useGroupQuery(query: string, skip?: boolean) {
  return useGroupQueryFilter({ q: query, skip: !!skip });
}

const GroupsFilter: React.FC<IGroupsFilter> = ({
  criterion,
  setCriterion,
}) => {
  return (
    <HierarchicalObjectsFilter
      criterion={criterion}
      setCriterion={setCriterion}
      useResults={useGroupQuery}
      singleValue
    />
  );
};

export const SidebarGroupsFilter: React.FC<{
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
    useQuery: useGroupQueryFilter,
    singleValue: true,
    hierarchical: true,
    includeSubMessageID: "sub_groups",
  });

  // Build candidates list
  // Strategy: When facets are loaded and no search query, USE facet results as candidates
  // (since facets returns the TOP N most relevant items by count).
  // When user searches, use search results merged with facet counts.
  const candidatesWithCounts: Option[] = useMemo(() => {
    const hasValidFacets = facetCounts.groups.size > 0 && !facetsLoading;
    const hasSearchQuery = state.query && state.query.length > 0;
    
    // Extract modifier options from candidates
    const modifierOptions = state.candidates.filter(c => c.className === "modifier-object");
    
    // Get selected IDs to exclude from candidates
    const selectedIds = new Set(state.selected.map(s => s.id));
    
    if (hasValidFacets && !hasSearchQuery) {
      // No search query: Use facet results directly as candidates (with labels from facets)
      const facetCandidates: Option[] = [];
      facetCounts.groups.forEach((facetData, id) => {
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
          const facetData = hasValidFacets ? facetCounts.groups.get(c.id) : undefined;
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

export default GroupsFilter;

