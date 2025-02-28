import React, { useMemo } from "react";
import { useFindStudiosQuery } from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { sortByRelevance } from "src/utils/query";

interface IStudiosFilter {
  criterion: StudiosCriterion;
  setCriterion: (c: StudiosCriterion) => void;
}

function useStudioQuery(query: string) {
  const { data, loading } = useFindStudiosQuery({
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

export default StudiosFilter;
