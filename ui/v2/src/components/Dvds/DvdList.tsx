import _ from "lodash";
import React, { FunctionComponent } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { FindDvdsQuery, FindDvdsVariables } from "../../core/generated-graphql";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";
import { DvdCard } from "./DvdCard";

interface IProps extends IBaseProps {}

export const DvdList: FunctionComponent<IProps> = (props: IProps) => {
  const listData = ListHook.useList({
    filterMode: FilterMode.Dvds,
    props,
    renderContent,
  });

  function renderContent(result: QueryHookResult<FindDvdsQuery, FindDvdsVariables>, filter: ListFilterModel) {
    if (!result.data || !result.data.findDvds) { return; }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="grid">
          {result.data.findDvds.dvds.map((dvd) => (<DvdCard key={dvd.id} dvd={dvd} />))}
        </div>
      );
    } else if (filter.displayMode === DisplayMode.List) {
      return <h1>TODO</h1>;
    } else if (filter.displayMode === DisplayMode.Wall) {
      return <h1>TODO</h1>;
    }
  }

  return listData.template;
};
