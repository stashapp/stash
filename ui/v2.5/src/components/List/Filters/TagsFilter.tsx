import React, { ReactNode, useMemo } from "react";
import {
  CriterionModifier,
  TagDataFragment,
  TagFilterType,
  useFindTagsForSelectQuery,
} from "src/core/generated-graphql";
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
import { SidebarListFilter } from "./SidebarListFilter";
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

      // TODO - look for same in AND?
    }

    setObjectFilter(tagFilter, f.mode, filterOutput);
  }

  return makeQueryVariables(query, { tag_filter: tagFilter });
}

function sortResults(
  query: string,
  tags: Pick<TagDataFragment, "id" | "name" | "aliases">[]
) {
  return sortByRelevance(
    query,
    tags ?? [],
    (t) => t.name,
    (t) => t.aliases
  ).map((p) => {
    return {
      id: p.id,
      label: p.name,
    };
  });
}

function useTagQueryFilter(props: IUseQueryHookProps) {
  const { q: query, filter: f, skip, filterHook } = props;
  const appliedFilter = filterHook && f ? filterHook(f.clone()) : f;

  const { data, loading } = useFindTagsForSelectQuery({
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
  const state = useLabeledIdFilterState({
    filter,
    setFilter,
    filterHook,
    option,
    useQuery: useTagQueryFilter,
    hierarchical: true,
    includeSubMessageID: "sub_tags",
  });

  return <SidebarListFilter {...state} title={title} sectionID={sectionID} />;
};

export default TagsFilter;
