import React from "react";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import { useFindPerformersQuery } from "src/core/generated-graphql";
import { ObjectsFilter } from "./SelectableFilter";

interface IPerformersFilter {
  criterion: PerformersCriterion;
  setCriterion: (c: PerformersCriterion) => void;
}

function usePerformerQuery(query: string) {
  const results = useFindPerformersQuery({
    variables: {
      filter: {
        q: query,
        per_page: 200,
      },
    },
  });

  return (
    results.data?.findPerformers.performers.map((p) => {
      return {
        id: p.id,
        label: p.name,
      };
    }) ?? []
  );
}

const PerformersFilter: React.FC<IPerformersFilter> = ({
  criterion,
  setCriterion,
}) => {
  return (
    <ObjectsFilter
      criterion={criterion}
      setCriterion={setCriterion}
      queryHook={usePerformerQuery}
    />
  );
};

export default PerformersFilter;
