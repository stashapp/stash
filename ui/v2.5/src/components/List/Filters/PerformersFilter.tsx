import React, { useMemo } from "react";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import { useFindPerformersQuery } from "src/core/generated-graphql";
import { ObjectsFilter } from "./SelectableFilter";

interface IPerformersFilter {
  criterion: PerformersCriterion;
  setCriterion: (c: PerformersCriterion) => void;
}

function usePerformerQuery(query: string) {
  const { data, loading } = useFindPerformersQuery({
    variables: {
      filter: {
        q: query,
        per_page: 200,
      },
    },
  });

  const results = useMemo(
    () =>
      data?.findPerformers.performers.map((p) => {
        return {
          id: p.id,
          label: p.name,
        };
      }) ?? [],
    [data]
  );

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

export default PerformersFilter;
