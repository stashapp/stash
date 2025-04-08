import React, { ReactNode, useMemo } from "react";
import { useFindStudiosForSelectQuery } from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { sortByRelevance } from "src/utils/query";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useLabeledIdFilterState } from "./LabeledIdFilter";
import { SidebarListFilter } from "./SidebarListFilter";

interface IStudiosFilter {
  criterion: StudiosCriterion;
  setCriterion: (c: StudiosCriterion) => void;
}

function useStudioQuery(query: string) {
  const { data, loading } = useFindStudiosForSelectQuery({
    variables: {
      filter: {
        q: query,
        per_page: 200,
      },
    },
  });

  const results = useMemo(() => {
    return sortByRelevance(
      query,
      data?.findStudios.studios ?? [],
      (s) => s.name,
      (s) => s.aliases
    ).map((p) => {
      return {
        id: p.id,
        label: p.name,
      };
    });
  }, [data, query]);

  return { results, loading };
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
    useQuery: useStudioQuery,
    singleValue: true,
    hierarchical: true,
    includeSubMessageID: "subsidiary_studios",
  });

  return <SidebarListFilter {...state} title={title} />;
};

export default StudiosFilter;
