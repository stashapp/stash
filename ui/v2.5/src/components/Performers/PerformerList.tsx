import _ from "lodash";
import React from "react";
import { useHistory } from "react-router-dom";
import { FindPerformersQueryResult } from "src/core/generated-graphql";
import { queryFindPerformers } from "src/core/StashService";
import { usePerformersList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { PerformerCard } from "./PerformerCard";
import { PerformerListTable } from "./PerformerListTable";

export const PerformerList: React.FC = () => {
  const history = useHistory();
  const otherOperations = [
    {
      text: "Open Random",
      onClick: getRandom,
    },
  ];

  const addKeybinds = (
    result: FindPerformersQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      getRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  const listData = usePerformersList({
    otherOperations,
    renderContent,
    addKeybinds,
  });

  async function getRandom(
    result: FindPerformersQueryResult,
    filter: ListFilterModel
  ) {
    if (result.data?.findPerformers) {
      const { count } = result.data.findPerformers;
      const index = Math.floor(Math.random() * count);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindPerformers(filterCopy);
      if (
        singleResult &&
        singleResult.data &&
        singleResult.data.findPerformers &&
        singleResult.data.findPerformers.performers.length === 1
      ) {
        const { id } = singleResult!.data!.findPerformers!.performers[0]!;
        history.push(`/performers/${id}`);
      }
    }
  }

  function renderContent(
    result: FindPerformersQueryResult,
    filter: ListFilterModel
  ) {
    if (!result.data?.findPerformers) {
      return;
    }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row justify-content-center">
          {result.data.findPerformers.performers.map((p) => (
            <PerformerCard key={p.id} performer={p} />
          ))}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <PerformerListTable
          performers={result.data.findPerformers.performers}
        />
      );
    }
  }

  return listData.template;
};
