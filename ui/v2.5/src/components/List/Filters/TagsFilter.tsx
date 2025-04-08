import React, { ReactNode, useMemo } from "react";
import { useFindTagsForSelectQuery } from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { sortByRelevance } from "src/utils/query";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useLabeledIdFilterState } from "./LabeledIdFilter";
import { SidebarListFilter } from "./SidebarListFilter";

interface ITagsFilter {
  criterion: StudiosCriterion;
  setCriterion: (c: StudiosCriterion) => void;
}

function useTagQuery(query: string) {
  const { data, loading } = useFindTagsForSelectQuery({
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
      data?.findTags.tags ?? [],
      (t) => t.name,
      (t) => t.aliases
    ).map((p) => {
      return {
        id: p.id,
        label: p.name,
      };
    });
  }, [data, query]);

  return { results, loading };
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
}> = ({ title, option, filter, setFilter }) => {
  const state = useLabeledIdFilterState({
    filter,
    setFilter,
    option,
    useQuery: useTagQuery,
    hierarchical: true,
    includeSubMessageID: "sub_tags",
  });

  return <SidebarListFilter {...state} title={title} />;
};

export default TagsFilter;
