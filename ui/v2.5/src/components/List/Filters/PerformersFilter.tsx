import React, { ReactNode, useMemo } from "react";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import { useFindPerformersForSelectQuery } from "src/core/generated-graphql";
import { ObjectsFilter } from "./SelectableFilter";
import { sortByRelevance } from "src/utils/query";
import { ListFilterModel } from "src/models/list-filter/filter";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { useLabeledIdFilterState } from "./LabeledIdFilter";
import { SidebarListFilter } from "./SidebarListFilter";

interface IPerformersFilter {
  criterion: PerformersCriterion;
  setCriterion: (c: PerformersCriterion) => void;
}

function usePerformerQuery(query: string, skip?: boolean) {
  const { data, loading } = useFindPerformersForSelectQuery({
    variables: {
      filter: {
        q: query,
        per_page: 200,
      },
    },
    skip,
  });

  const results = useMemo(() => {
    return sortByRelevance(
      query,
      data?.findPerformers.performers ?? [],
      (p) => p.name,
      (p) => p.alias_list
    ).map((p) => {
      return {
        id: p.id,
        label: p.name,
      };
    });
  }, [data, query]);

  return { results, loading };
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
    useQuery: usePerformerQuery,
  });

  return <SidebarListFilter {...state} title={title} />;
};

export default PerformersFilter;
