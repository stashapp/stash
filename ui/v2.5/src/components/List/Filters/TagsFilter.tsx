import React, { useMemo } from "react";
import { useFindTagsQuery } from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { sortByRelevance } from "src/utils/query";

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

  const results = useMemo(() => {
    return sortByRelevance(
      query,
      data?.findTags.tags ?? [],
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

export default TagsFilter;
