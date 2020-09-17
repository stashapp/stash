import React, { useState } from "react";
import _ from "lodash";
import { useHistory } from "react-router-dom";
import {
  FindImagesQueryResult,
  SlimImageDataFragment,
} from "src/core/generated-graphql";
import { queryFindImages } from "src/core/StashService";
import { useImagesList } from "src/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { showWhenSelected } from "src/hooks/ListHook";
import { WallPanel } from "../Wall/WallPanel";
import { ImageCard } from "./ImageCard";
import { EditImagesDialog } from "./EditImagesDialog";
import { DeleteImagesDialog } from "./DeleteImagesDialog";
import { ImageExportDialog } from "./ImageExportDialog";

interface IImageList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  persistState?: boolean;
}

export const ImageList: React.FC<IImageList> = ({
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
    result: FindImagesQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      viewRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  const listData = useImagesList({
    zoomable: true,
    selectable: true,
    otherOperations,
    renderContent,
    renderEditDialog: renderEditImagesDialog,
    renderDeleteDialog: renderDeleteImagesDialog,
    filterHook,
    addKeybinds,
    persistState,
  });

  async function viewRandom(
    result: FindImagesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random image
    if (result.data && result.data.findImages) {
      const { count } = result.data.findImages;

      const index = Math.floor(Math.random() * count);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindImages(filterCopy);
      if (
        singleResult &&
        singleResult.data &&
        singleResult.data.findImages &&
        singleResult.data.findImages.images.length === 1
      ) {
        const { id } = singleResult!.data!.findImages!.images[0];
        // navigate to the image player page
        history.push(`/images/${id}?autoplay=true`);
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

  function maybeRenderImageExportDialog(selectedIds: Set<string>) {
    if (isExportDialogOpen) {
      return (
        <>
          <ImageExportDialog
            selectedIds={Array.from(selectedIds.values())}
            all={isExportAll}
            onClose={() => {
              setIsExportDialogOpen(false);
            }}
          />
        </>
      );
    }
  }

  function renderEditImagesDialog(
    selectedImages: SlimImageDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return (
      <>
        <EditImagesDialog selected={selectedImages} onClose={onClose} />
      </>
    );
  }

  function renderDeleteImagesDialog(
    selectedImages: SlimImageDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return (
      <>
        <DeleteImagesDialog selected={selectedImages} onClose={onClose} />
      </>
    );
  }

  function renderImageCard(
    image: SlimImageDataFragment,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    return (
      <ImageCard
        key={image.id}
        image={image}
        zoomIndex={zoomIndex}
        selecting={selectedIds.size > 0}
        selected={selectedIds.has(image.id)}
        onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
          listData.onSelectChange(image.id, selected, shiftKey)
        }
      />
    );
  }

  function renderImages(
    result: FindImagesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    if (!result.data || !result.data.findImages) {
      return;
    }
    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row justify-content-center">
          {result.data.findImages.images.map((image) =>
            renderImageCard(image, selectedIds, zoomIndex)
          )}
        </div>
      );
    }
    // if (filter.displayMode === DisplayMode.List) {
    //   return <ImageListTable images={result.data.findImages.images} />;
    // }
    if (filter.displayMode === DisplayMode.Wall) {
      return <WallPanel images={result.data.findImages.images} />;
    }
  }

  function renderContent(
    result: FindImagesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    return (
      <>
        {maybeRenderImageExportDialog(selectedIds)}
        {renderImages(result, filter, selectedIds, zoomIndex)}
      </>
    );
  }

  return listData.template;
};
