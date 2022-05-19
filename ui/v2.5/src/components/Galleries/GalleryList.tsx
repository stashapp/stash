import React, { useState } from "react";
import { useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { Table } from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import {
  FindGalleriesQueryResult,
  SlimGalleryDataFragment,
} from "src/core/generated-graphql";
import { useGalleriesList } from "src/hooks";
import { showWhenSelected, PersistanceLevel } from "src/hooks/ListHook";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { queryFindGalleries } from "src/core/StashService";
import { GalleryCard } from "./GalleryCard";
import GalleryWallCard from "./GalleryWallCard";
import { EditGalleriesDialog } from "./EditGalleriesDialog";
import { DeleteGalleriesDialog } from "./DeleteGalleriesDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { galleryTitle } from "src/core/galleries";

interface IGalleryList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  persistState?: PersistanceLevel;
}

export const GalleryList: React.FC<IGalleryList> = ({
  filterHook,
  persistState,
}) => {
  const intl = useIntl();
  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.view_random" }),
      onClick: viewRandom,
    },
    {
      text: intl.formatMessage({ id: "actions.export" }),
      onClick: onExport,
      isDisplayed: showWhenSelected,
    },
    {
      text: intl.formatMessage({ id: "actions.export_all" }),
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
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindGalleries(filterCopy);
      if (singleResult.data.findGalleries.galleries.length === 1) {
        const { id } = singleResult.data.findGalleries.galleries[0];
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
    selectedImages: SlimGalleryDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return (
      <>
        <EditGalleriesDialog selected={selectedImages} onClose={onClose} />
      </>
    );
  }

  function renderDeleteGalleriesDialog(
    selectedImages: SlimGalleryDataFragment[],
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
    selectedIds: Set<string>
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
              zoomIndex={filter.zoomIndex}
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
              <th>{intl.formatMessage({ id: "actions.preview" })}</th>
              <th className="d-none d-sm-none">
                {intl.formatMessage({ id: "title" })}
              </th>
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
                    {galleryTitle(gallery)} ({gallery.image_count}{" "}
                    {gallery.image_count === 1 ? "image" : "images"})
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <div className="row">
          <div className="GalleryWall">
            {result.data.findGalleries.galleries.map((gallery) => (
              <GalleryWallCard key={gallery.id} gallery={gallery} />
            ))}
          </div>
        </div>
      );
    }
  }

  function renderContent(
    result: FindGalleriesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    return (
      <>
        {maybeRenderGalleryExportDialog(selectedIds)}
        {renderGalleries(result, filter, selectedIds)}
      </>
    );
  }

  return listData.template;
};
