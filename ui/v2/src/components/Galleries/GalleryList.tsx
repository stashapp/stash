import { HTMLTable } from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { QueryHookResult } from "react-apollo-hooks";
import { Link } from "react-router-dom";
import { FindGalleriesQuery, FindGalleriesVariables } from "../../core/generated-graphql";
import { ListHook } from "../../hooks/ListHook";
import { IBaseProps } from "../../models/base-props";
import { ListFilterModel } from "../../models/list-filter/filter";
import { DisplayMode, FilterMode } from "../../models/list-filter/types";

interface IProps extends IBaseProps {}

export const GalleryList: FunctionComponent<IProps> = (props: IProps) => {
  const listData = ListHook.useList({
    filterMode: FilterMode.Galleries,
    props,
    renderContent,
  });

  function renderContent(result: QueryHookResult<FindGalleriesQuery, FindGalleriesVariables>, filter: ListFilterModel) {
    if (!result.data || !result.data.findGalleries) { return; }
    if (filter.displayMode === DisplayMode.Grid) {
      return <h1>TODO</h1>;
    } else if (filter.displayMode === DisplayMode.List) {
      return (
        <HTMLTable style={{margin: "0 auto"}}>
          <thead>
            <tr>
              <th>Preview</th>
              <th>Path</th>
            </tr>
          </thead>
          <tbody>
            {result.data.findGalleries.galleries.map((gallery) => (
              <tr key={gallery.id}>
                <td>
                  <Link to={`/galleries/${gallery.id}`}>
                    {gallery.files.length > 0 ? <img alt={gallery.title} src={`${gallery.files[0].path}?thumb=true`} /> : undefined}
                  </Link>
                </td>
                <td><Link to={`/galleries/${gallery.id}`}>{gallery.path} ({gallery.files.length} {gallery.files.length === 1 ? 'image' : 'images'})</Link></td>
              </tr>
            ))}
          </tbody>
        </HTMLTable>
      );
    } else if (filter.displayMode === DisplayMode.Wall) {
      return <h1>TODO</h1>;
    }
  }

  return listData.template;
};
