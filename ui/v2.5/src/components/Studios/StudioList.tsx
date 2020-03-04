import React from "react";
import { FindStudiosQueryResult } from "src/core/generated-graphql";
import { useStudiosList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { StudioCard } from "./StudioCard";

export const StudioList: React.FC = () => {
  const listData = useStudiosList({
    renderContent
  });

  function renderContent(
    result: FindStudiosQueryResult,
    filter: ListFilterModel
  ) {
    if (!result.data?.findStudios) return;

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row px-xl-5 justify-content-center">
          {result.data.findStudios.studios.map(studio => (
            <StudioCard key={studio.id} studio={studio} />
          ))}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return <h1>TODO</h1>;
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return <h1>TODO</h1>;
    }
  }

  return listData.template;
};
