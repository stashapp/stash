import React, { ReactNode, useMemo } from "react";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import {
  CriterionModifier,
  FindPerformersForSelectQueryVariables,
  PerformerDataFragment,
  PerformerFilterType,
  useFindPerformersForSelectQuery,
} from "src/core/generated-graphql";
import { ObjectsFilter } from "./SelectableFilter";
import { sortByRelevance } from "src/utils/query";
import { ListFilterModel } from "src/models/list-filter/filter";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import {
  makeQueryVariables,
  setObjectFilter,
  useLabeledIdFilterState,
} from "./LabeledIdFilter";
import { SidebarListFilter } from "./SidebarListFilter";

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
): FindPerformersForSelectQueryVariables {
  const performerFilter: PerformerFilterType = {};

  if (f) {
    const filterOutput = f.makeFilter();

    // if performer modifier is includes, take it out of the filter
    if (
      (filterOutput.performers as IHasModifier)?.modifier ===
      CriterionModifier.Includes
    ) {
      delete filterOutput.performers;

      // TODO - look for same in AND?
    }

    setObjectFilter(performerFilter, f.mode, filterOutput);
  }

  return makeQueryVariables(query, { performer_filter: performerFilter });
}

function sortResults(
  query: string,
  performers?: Pick<PerformerDataFragment, "name" | "alias_list" | "id">[]
) {
  return sortByRelevance(
    query,
    performers ?? [],
    (p) => p.name,
    (p) => p.alias_list
  ).map((p) => {
    return {
      id: p.id,
      label: p.name,
    };
  });
}

function usePerformerQueryFilter(
  query: string,
  f?: ListFilterModel,
  skip?: boolean
) {
  const { data, loading } = useFindPerformersForSelectQuery({
    variables: queryVariables(query, f),
    skip,
  });

  const results = useMemo(
    () => sortResults(query, data?.findPerformers.performers),
    [data, query]
  );

  return { results, loading };
}

function usePerformerQuery(query: string, skip?: boolean) {
  return usePerformerQueryFilter(query, undefined, skip);
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
}> = ({ title, option, filter, setFilter }) => {
  const state = useLabeledIdFilterState({
    filter,
    setFilter,
    option,
    useQuery: usePerformerQueryFilter,
  });

  return <SidebarListFilter {...state} title={title} />;
};

export default PerformersFilter;
