import React, { useMemo } from "react";
import { useFindTagsQuery } from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";

interface ITagsFilter {
  criterion: StudiosCriterion;
  setCriterion: (c: StudiosCriterion) => void;
}

function useTagQuery(query: string) {
  const { data, loading } = useFindTagsQuery({
    variables: {
      filter: {
        q: query,
        per_page: 200,
      },
    },
  });

  const results = useMemo(
    () =>
      data?.findTags.tags.map((p) => {
        return {
          id: p.id,
          label: p.name,
        };
      }) ?? [],
    [data]
  );

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

export default TagsFilter;
