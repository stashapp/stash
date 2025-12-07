import React, { ReactNode, useCallback, useContext, useMemo, useState } from "react";
import {
  CriterionModifier,
  FilterStudioDataFragment,
  StudioFilterType,
  useFindStudiosForFilterQuery,
} from "src/core/generated-graphql";
import { FacetCountsContext } from "src/extensions/hooks/useFacetCounts";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import {
  StudiosCriterion,
  ParentStudiosCriterion,
} from "src/models/list-filter/criteria/studios";
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

interface IStudiosFilter {
  criterion: StudiosCriterion;
  setCriterion: (c: StudiosCriterion) => void;
}

function queryVariables(query: string, f?: ListFilterModel) {
  const studioFilter: StudioFilterType = {};

  if (f) {
    const filterOutput = f.makeFilter();
    delete filterOutput.studios;
    setObjectFilter(studioFilter, f.mode, filterOutput, "studios");
  }

  return makeQueryVariables(query, { studio_filter: studioFilter });
}

function sortResults(
  query: string,
  studios: FilterStudioDataFragment[]
) {
  return sortByRelevance(
    query,
    studios ?? [],
    (s) => s.name,
    (s) => s.aliases
  ).map((p) => ({
    id: p.id,
    label: p.name,
  }));
}

function useStudioQueryFilter(props: IUseQueryHookProps) {
  const { q: query, filter: f, skip, filterHook } = props;
  const appliedFilter = filterHook && f ? filterHook(f.clone()) : f;

  const { data, loading } = useFindStudiosForFilterQuery({
    variables: queryVariables(query, appliedFilter),
    skip,
  });

  const results = useMemo(
    () => sortResults(query, data?.findStudios.studios ?? []),
    [data?.findStudios.studios, query]
  );

  return { results, loading };
}

function useStudioQuery(query: string, skip?: boolean) {
  return useStudioQueryFilter({ q: query, skip: !!skip });
}

const StudiosFilter: React.FC<IStudiosFilter> = ({
  criterion,
  setCriterion,
}) => {
  return (
    <HierarchicalObjectsFilter
      criterion={criterion}
      setCriterion={setCriterion}
      useResults={useStudioQuery}
      singleValue
    />
  );
};

export const SidebarStudiosFilter: React.FC<{
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
    useQuery: useStudioQueryFilter,
    singleValue: true,
    hierarchical: true,
    includeSubMessageID: "subsidiary_studios",
  });

  // Build candidates list
  // Strategy: When facets are loaded and no search query, USE facet results as candidates
  // (since facets returns the TOP N most relevant items by count).
  // When user searches, use search results merged with facet counts.
  const candidatesWithCounts: Option[] = useMemo(() => {
    const hasValidFacets = facetCounts.studios.size > 0 && !facetsLoading;
    const hasSearchQuery = state.query && state.query.length > 0;
    
    // Extract modifier options from candidates
    const modifierOptions = state.candidates.filter(c => c.className === "modifier-object");
    
    // Get selected IDs to exclude from candidates
    const selectedIds = new Set(state.selected.map(s => s.id));
    
    if (hasValidFacets && !hasSearchQuery) {
      // No search query: Use facet results directly as candidates (with labels from facets)
      const facetCandidates: Option[] = [];
      facetCounts.studios.forEach((facetData, id) => {
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
          const facetData = hasValidFacets ? facetCounts.studios.get(c.id) : undefined;
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

// Hook for non-hierarchical ILabeledIdCriterion (like ParentStudiosCriterion)
function useParentStudiosFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  const [skip, setSkip] = useState(true);
  const [query, setQuery] = useState("");

  const { results: queryResults } = useCacheResults(
    useStudioQueryFilter({ q: query, skip })
  );

  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as ParentStudiosCriterion;
    return filter.makeCriterion(option.type) as ParentStudiosCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: ParentStudiosCriterion) => {
      const newCriteria = filter.criteria.filter(
        (cc) => cc.criterionOption.type !== option.type
      );
      if (c.isValid()) newCriteria.push(c);
      setFilter(filter.setCriteria(newCriteria));
    },
    [option.type, setFilter, filter]
  );

  const { modifier } = criterion;

  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
    };
  }, [modifier]);

  const selected = useMemo(() => {
    const modifierValues: Option[] = Object.entries(selectedModifiers)
      .filter((v) => v[1])
      .map((v) => ({
        id: v[0],
        label: `(${intl.formatMessage({
          id: `criterion_modifier_values.${v[0]}`,
        })})`,
        className: "modifier-object",
      }));

    return modifierValues.concat(
      (criterion.value || []).map((s: ILabeledId) => ({
        id: s.id,
        label: s.label,
      }))
    );
  }, [intl, selectedModifiers, criterion.value]);

  const candidates = useMemo(() => {
    if (
      !queryResults ||
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      const modifierCandidates: Option[] = [];

      if (
        modifier === CriterionModifier.Includes &&
        criterion.value.length === 0
      ) {
        modifierCandidates.push({
          id: "any",
          label: `(${intl.formatMessage({
            id: "criterion_modifier_values.any",
          })})`,
          className: "modifier-object",
          canExclude: false,
        });
        modifierCandidates.push({
          id: "none",
          label: `(${intl.formatMessage({
            id: "criterion_modifier_values.none",
          })})`,
          className: "modifier-object",
          canExclude: false,
        });
      }

      return modifierCandidates;
    }

    const selectedIds = new Set(criterion.value.map((v: ILabeledId) => v.id));
    const filteredResults = queryResults.filter(
      (p: ILabeledId) => !selectedIds.has(p.id)
    );

    const modifierCandidates: Option[] = [];
    if (
      modifier === CriterionModifier.Includes &&
      criterion.value.length === 0
    ) {
      modifierCandidates.push({
        id: "any",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.any",
        })})`,
        className: "modifier-object",
        canExclude: false,
      });
      modifierCandidates.push({
        id: "none",
        label: `(${intl.formatMessage({
          id: "criterion_modifier_values.none",
        })})`,
        className: "modifier-object",
        canExclude: false,
      });
    }

    return modifierCandidates.concat(
      filteredResults.map((r: ILabeledId) => ({
        id: r.id,
        label: r.label,
      }))
    );
  }, [queryResults, modifier, criterion.value, intl]);

  const onSelect = useCallback(
    (v: Option) => {
      const newCriterion = criterion.clone() as ParentStudiosCriterion;

      if (v.className === "modifier-object") {
        if (v.id === "any") {
          newCriterion.modifier = CriterionModifier.NotNull;
        } else if (v.id === "none") {
          newCriterion.modifier = CriterionModifier.IsNull;
        }
        setCriterion(newCriterion);
        return;
      }

      newCriterion.value = [...criterion.value, { id: v.id, label: v.label }];
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option) => {
      const newCriterion = criterion.clone() as ParentStudiosCriterion;

      if (v.className === "modifier-object") {
        newCriterion.modifier = CriterionModifier.Includes;
        setCriterion(newCriterion);
        return;
      }

      newCriterion.value = criterion.value.filter(
        (i: ILabeledId) => i.id !== v.id
      );
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onOpen = useCallback(() => {
    setSkip(false);
  }, []);

  return {
    candidates,
    onSelect: (item: Option) => onSelect(item),
    onUnselect: (item: Option) => onUnselect(item),
    selected,
    excluded: [] as Option[],
    canExclude: false,
    query,
    setQuery,
    onOpen,
  };
}

export const SidebarParentStudiosFilter: React.FC<{
  title?: ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}> = ({ title, option, filter, setFilter, sectionID }) => {
  const state = useParentStudiosFilterState({
    filter,
    setFilter,
    option,
  });

  return (
    <SidebarListFilter
      {...state}
      title={title}
      sectionID={sectionID}
      onSelect={(item) => state.onSelect(item)}
      onUnselect={(item) => state.onUnselect(item)}
    />
  );
};

export default StudiosFilter;
