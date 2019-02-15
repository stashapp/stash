import _ from "lodash";
import React, { FunctionComponent } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindPerformersQuery, FindPerformersVariables } from "../../core/generated-graphql";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";
import { PerformerCard } from "./PerformerCard";

interface IPerformerListProps extends IBaseProps {}

export const PerformerList: FunctionComponent<IPerformerListProps> = (props: IPerformerListProps) => {
  const listData = ListHook.useList({
    filterMode: FilterMode.Performers,
    props,
    renderContent,
  });

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
      return <h1>TODO</h1>;
    } else if (filter.displayMode === DisplayMode.Wall) {
      return;
    }
  }

  return listData.template;
};
