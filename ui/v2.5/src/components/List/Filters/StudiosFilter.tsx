import React, { ReactNode, useMemo } from "react";
import {
  StudioDataFragment,
  StudioFilterType,
  useFindStudiosForSelectQuery,
} from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { sortByRelevance } from "src/utils/query";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  makeQueryVariables,
  setObjectFilter,
  useLabeledIdFilterState,
} from "./LabeledIdFilter";
import { SidebarListFilter } from "./SidebarListFilter";

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

    setObjectFilter(studioFilter, f.mode, filterOutput);
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

function useStudioQueryFilter(
  query: string,
  filter?: ListFilterModel,
  skip?: boolean
) {
  const { data, loading } = useFindStudiosForSelectQuery({
    variables: queryVariables(query, filter),
    skip,
  });

  const results = useMemo(
    () => sortResults(query, data?.findStudios.studios ?? []),
    [data?.findStudios.studios, query]
  );

  return { results, loading };
}

function useStudioQuery(query: string, skip?: boolean) {
  return useStudioQueryFilter(query, undefined, skip);
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
}> = ({ title, option, filter, setFilter }) => {
  const state = useLabeledIdFilterState({
    filter,
    setFilter,
    option,
    useQuery: useStudioQueryFilter,
    singleValue: true,
    hierarchical: true,
    includeSubMessageID: "subsidiary_studios",
  });

  return <SidebarListFilter {...state} title={title} />;
};

export default StudiosFilter;
