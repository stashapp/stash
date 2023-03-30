import React from "react";
import { useFindTagsQuery } from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";

interface ITagsFilter {
  criterion: StudiosCriterion;
  setCriterion: (c: StudiosCriterion) => void;
}

function useStudioQuery(query: string) {
  const results = useFindTagsQuery({
    variables: {
      filter: {
        q: query,
        per_page: 200,
      },
    },
  });

  return (
    results.data?.findTags.tags.map((p) => {
      return {
        id: p.id,
        label: p.name,
      };
    }) ?? []
  );
}

const TagsFilter: React.FC<ITagsFilter> = ({ criterion, setCriterion }) => {
  return (
    <HierarchicalObjectsFilter
      criterion={criterion}
      setCriterion={setCriterion}
      queryHook={useStudioQuery}
    />
  );
};

export default TagsFilter;
