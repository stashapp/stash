import React from "react";
import { Table } from 'react-bootstrap';
import { QueryHookResult } from "react-apollo-hooks";
import { Link } from "react-router-dom";
import { FindGalleriesQuery, FindGalleriesVariables } from "src/core/generated-graphql";
import { ListHook } from "src/hooks";
import { IBaseProps } from "src/models/base-props";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode, FilterMode } from "src/models/list-filter/types";

interface IProps extends IBaseProps {}

export const GalleryList: React.FC<IProps> = (props: IProps) => {
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
        <Table style={{margin: "0 auto"}}>
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
                    {gallery.files.length > 0 ? <img alt="" src={`${gallery.files[0].path}?thumb=true`} /> : undefined}
                  </Link>
                </td>
                <td><Link to={`/galleries/${gallery.id}`}>{gallery.path}</Link></td>
              </tr>
            ))}
          </tbody>
        </Table>
      );
    } else if (filter.displayMode === DisplayMode.Wall) {
      return <h1>TODO</h1>;
    }
  }

  return listData.template;
};
