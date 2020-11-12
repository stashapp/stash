import React, { useState } from "react";
import _ from "lodash";
import { Table } from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import {
  FindGalleriesQueryResult,
  GalleryDataFragment,
} from "src/core/generated-graphql";
import { useGalleriesList } from "src/hooks";
import { showWhenSelected } from "src/hooks/ListHook";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { queryFindGalleries } from "src/core/StashService";
import { GalleryCard } from "./GalleryCard";
import { EditGalleriesDialog } from "./EditGalleriesDialog";
import { DeleteGalleriesDialog } from "./DeleteGalleriesDialog";
import { ExportDialog } from "../Shared/ExportDialog";

interface IGalleryList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  persistState?: boolean;
}

export const GalleryList: React.FC<IGalleryList> = ({
  filterHook,
  persistState,
}) => {
  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const otherOperations = [
    {
      text: "View Random",
      onClick: viewRandom,
    },
    {
      text: "Export...",
      onClick: onExport,
      isDisplayed: showWhenSelected,
    },
    {
      text: "Export all...",
      onClick: onExportAll,
    },
  ];

  const addKeybinds = (
    result: FindGalleriesQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      viewRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  const listData = useGalleriesList({
    zoomable: true,
    selectable: true,
    otherOperations,
    renderContent,
    renderEditDialog: renderEditGalleriesDialog,
    renderDeleteDialog: renderDeleteGalleriesDialog,
    filterHook,
    addKeybinds,
    persistState,
  });

  async function viewRandom(
    result: FindGalleriesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random image
    if (result.data && result.data.findGalleries) {
      const { count } = result.data.findGalleries;

      const index = Math.floor(Math.random() * count);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindGalleries(filterCopy);
      if (
        singleResult &&
        singleResult.data &&
        singleResult.data.findGalleries &&
        singleResult.data.findGalleries.galleries.length === 1
      ) {
        const { id } = singleResult!.data!.findGalleries!.galleries[0];
        // navigate to the image player page
        history.push(`/galleries/${id}`);
      }
    }
  }

  async function onExport() {
    setIsExportAll(false);
    setIsExportDialogOpen(true);
  }

  async function onExportAll() {
    setIsExportAll(true);
    setIsExportDialogOpen(true);
  }

  function maybeRenderGalleryExportDialog(selectedIds: Set<string>) {
    if (isExportDialogOpen) {
      return (
        <>
          <ExportDialog
            exportInput={{
              galleries: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => {
              setIsExportDialogOpen(false);
            }}
          />
        </>
      );
    }
  }

  function renderEditGalleriesDialog(
    selectedImages: GalleryDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return (
      <>
        <EditGalleriesDialog selected={selectedImages} onClose={onClose} />
      </>
    );
  }

  function renderDeleteGalleriesDialog(
    selectedImages: GalleryDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return (
      <>
        <DeleteGalleriesDialog selected={selectedImages} onClose={onClose} />
      </>
    );
  }

  function renderGalleries(
    result: FindGalleriesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    if (!result.data || !result.data.findGalleries) {
      return;
    }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row justify-content-center">
          {result.data.findGalleries.galleries.map((gallery) => (
            <GalleryCard
              key={gallery.id}
              gallery={gallery}
              zoomIndex={zoomIndex}
              selecting={selectedIds.size > 0}
              selected={selectedIds.has(gallery.id)}
              onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
                listData.onSelectChange(gallery.id, selected, shiftKey)
              }
            />
          ))}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <Table className="col col-sm-6 mx-auto">
          <thead>
            <tr>
              <th>Preview</th>
              <th className="d-none d-sm-none">Title</th>
            </tr>
          </thead>
          <tbody>
            {result.data.findGalleries.galleries.map((gallery) => (
              <tr key={gallery.id}>
                <td>
                  <Link to={`/galleries/${gallery.id}`}>
                    {gallery.cover ? (
                      <img
                        alt={gallery.title ?? ""}
                        className="w-100 w-sm-auto"
                        src={`${gallery.cover.paths.thumbnail}`}
                      />
                    ) : undefined}
                  </Link>
                </td>
                <td className="d-none d-sm-block">
                  <Link to={`/galleries/${gallery.id}`}>
                    {gallery.title ?? gallery.path} ({gallery.images.length}{" "}
                    {gallery.images.length === 1 ? "image" : "images"})
                  </Link>
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

  function renderContent(
    result: FindGalleriesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    return (
      <>
        {maybeRenderGalleryExportDialog(selectedIds)}
        {renderGalleries(result, filter, selectedIds, zoomIndex)}
      </>
    );
  }

  return listData.template;
};
