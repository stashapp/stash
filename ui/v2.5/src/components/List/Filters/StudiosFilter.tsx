import React, { ReactNode, useCallback, useMemo, useState } from "react";
import {
  CriterionModifier,
  StudioDataFragment,
  StudioFilterType,
  useFindStudiosForSelectQuery,
} from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import {
  StudiosCriterion,
  ParentStudiosCriterion,
} from "src/models/list-filter/criteria/studios";
import { sortByRelevance } from "src/utils/query";
import {
  CriterionOption,
  ModifierCriterion,
} from "src/models/list-filter/criteria/criterion";
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

    // always remove studio filter from the filter
    // since modifier is includes
    delete filterOutput.studios;

    // TODO - look for same in AND?

    setObjectFilter(studioFilter, f.mode, filterOutput, "studios");
  }

  return makeQueryVariables(query, { studio_filter: studioFilter });
}

function sortResults(
  query: string,
  studios: Pick<StudioDataFragment, "id" | "name" | "aliases">[]
) {
  return sortByRelevance(
    query,
    studios ?? [],
    (s) => s.name,
    (s) => s.aliases
  ).map((p) => {
    return {
      id: p.id,
      label: p.name,
    };
  });
}

function useStudioQueryFilter(props: IUseQueryHookProps) {
  const { q: query, filter: f, skip, filterHook } = props;
  const appliedFilter = filterHook && f ? filterHook(f.clone()) : f;

  const { data, loading } = useFindStudiosForSelectQuery({
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

  return <SidebarListFilter {...state} title={title} sectionID={sectionID} />;
};

// Hook for non-hierarchical ILabeledIdCriterion (like ParentStudiosCriterion)
function useParentStudiosFilterState(props: {
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}) {
  const intl = useIntl();
  const { option, filter, setFilter } = props;

  // defer querying until the user opens the filter
  const [skip, setSkip] = useState(true);
  const [query, setQuery] = useState("");

  const { results: queryResults } = useCacheResults(
    useStudioQueryFilter({ q: query, skip })
  );

  // Get or create criterion
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

  // Build selected modifiers
  const selectedModifiers = useMemo(() => {
    return {
      any: modifier === CriterionModifier.NotNull,
      none: modifier === CriterionModifier.IsNull,
    };
  }, [modifier]);

  // Build selected items list
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

    // ILabeledIdCriterion stores value as ILabeledId[] directly
    return modifierValues.concat(
      (criterion.value || []).map((s: ILabeledId) => ({
        id: s.id,
        label: s.label,
      }))
    );
  }, [intl, selectedModifiers, criterion.value]);

  // Build candidates list
  const candidates = useMemo(() => {
    if (
      !queryResults ||
      modifier === CriterionModifier.IsNull ||
      modifier === CriterionModifier.NotNull
    ) {
      // Show modifier options when no items selected
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

    // Filter out already selected items
    const selectedIds = new Set(criterion.value.map((v: ILabeledId) => v.id));
    const filteredResults = queryResults.filter(
      (p: ILabeledId) => !selectedIds.has(p.id)
    );

    // Add modifier options if no items selected yet
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

      // Add to value array
      newCriterion.value = [...criterion.value, { id: v.id, label: v.label }];
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onUnselect = useCallback(
    (v: Option) => {
      const newCriterion = criterion.clone() as ParentStudiosCriterion;

      if (v.className === "modifier-object") {
        // Reset to default modifier
        newCriterion.modifier = CriterionModifier.Includes;
        setCriterion(newCriterion);
        return;
      }

      // Remove from value array
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

// Sidebar filter for parent studios (non-hierarchical ILabeledIdCriterion)
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
