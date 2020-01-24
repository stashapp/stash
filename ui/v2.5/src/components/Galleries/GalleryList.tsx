import React from "react";
import { Table } from "react-bootstrap";
import { Link } from "react-router-dom";
import { FindGalleriesQueryResult } from "src/core/generated-graphql";
import { useGalleriesList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";

export const GalleryList: React.FC = () => {
  const listData = useGalleriesList({
    renderContent
  });

  function renderContent(
    result: FindGalleriesQueryResult,
    filter: ListFilterModel
  ) {
    if (!result.data || !result.data.findGalleries) {
      return;
    }
    if (filter.displayMode === DisplayMode.Grid) {
      return <h1>TODO</h1>;
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <Table style={{ margin: "0 auto" }}>
          <thead>
            <tr>
              <th>Preview</th>
              <th>Path</th>
            </tr>
          </thead>
          <tbody>
            {result.data.findGalleries.galleries.map(gallery => (
              <tr key={gallery.id}>
                <td>
                  <Link to={`/galleries/${gallery.id}`}>
                    {gallery.files.length > 0 ? (
                      <img
                        alt={gallery.title ?? ""}
                        src={`${gallery.files[0].path}?thumb=true`}
                      />
                    ) : (
                      undefined
                    )}
                  </Link>
                </td>
                <td>
                  <Link to={`/galleries/${gallery.id}`}>{gallery.path}</Link>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return <h1>TODO</h1>;
    }
  }

  return listData.template;
};
