import _ from "lodash";
import React from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindPerformersQuery, FindPerformersVariables } from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { ListHook } from "src/hooks";
import { IBaseProps } from "src/models/base-props";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode, FilterMode } from "src/models/list-filter/types";
import { PerformerCard } from "./PerformerCard";
import { PerformerListTable } from "./PerformerListTable";

interface IPerformerListProps extends IBaseProps {}

export const PerformerList: React.FC<IPerformerListProps> = (props: IPerformerListProps) => {
  const otherOperations = [
    {
      text: "Open Random",
      onClick: getRandom,
    }
  ];

  const listData = ListHook.useList({
    filterMode: FilterMode.Performers,
    props,
    otherOperations: otherOperations,
    renderContent,
  });

  async function getRandom(result: QueryHookResult<FindPerformersQuery, FindPerformersVariables>, filter: ListFilterModel) {
    if (result.data && result.data.findPerformers) {
      let count = result.data.findPerformers.count;
      let index = Math.floor(Math.random() * count);
      let filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await StashService.queryFindPerformers(filterCopy);
      if (singleResult && singleResult.data && singleResult.data.findPerformers && singleResult.data.findPerformers.performers.length === 1) {
        let id = singleResult!.data!.findPerformers!.performers[0]!.id;
        props.history.push("/performers/" + id);
      }
    }
  }

  function renderContent(
    result: QueryHookResult<FindPerformersQuery, FindPerformersVariables>, filter: ListFilterModel) {
    if (!result.data || !result.data.findPerformers) { return; }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="grid">
          {result.data.findPerformers.performers.map((p) => (<PerformerCard key={p.id} performer={p} />))}
        </div>
      );
    } else if (filter.displayMode === DisplayMode.List) {
      return <PerformerListTable performers={result.data.findPerformers.performers}/>;
    } else if (filter.displayMode === DisplayMode.Wall) {
      return;
    }
  }

  return listData.template;
};
