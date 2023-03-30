import React from "react";
import { useFindStudiosQuery } from "src/core/generated-graphql";
import { HierarchicalObjectsFilter } from "./SelectableFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";

interface IStudiosFilter {
  criterion: StudiosCriterion;
  setCriterion: (c: StudiosCriterion) => void;
}

function useStudioQuery(query: string) {
  const results = useFindStudiosQuery({
    variables: {
      filter: {
        q: query,
        per_page: 200,
      },
    },
  });

  return (
    results.data?.findStudios.studios.map((p) => {
      return {
        id: p.id,
        label: p.name,
      };
    }) ?? []
  );
}

const StudiosFilter: React.FC<IStudiosFilter> = ({
  criterion,
  setCriterion,
}) => {
  return (
    <HierarchicalObjectsFilter
      single
      criterion={criterion}
      setCriterion={setCriterion}
      queryHook={useStudioQuery}
    />
  );
};

export default StudiosFilter;
