import React, { useCallback, useEffect, useState } from "react";
import _ from "lodash";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import {
  FindImagesQueryResult,
  SlimImageDataFragment,
} from "src/core/generated-graphql";
import * as GQL from "src/core/generated-graphql";
import { queryFindImages } from "src/core/StashService";
import { useImagesList, useLightbox } from "src/hooks";
import { TextUtils } from "src/utils";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import {
  IListHookOperation,
  showWhenSelected,
  PersistanceLevel,
  showWhenDisplayModeWall,
} from "src/hooks/ListHook";

import { ImageCard } from "./ImageCard";
import { EditImagesDialog } from "./EditImagesDialog";
import { DeleteImagesDialog } from "./DeleteImagesDialog";
import "flexbin/flexbin.css";
import { ExportDialog } from "../Shared/ExportDialog";

interface IImageWallProps {
  images: GQL.SlimImageDataFragment[];
  onChangePage: (page: number) => void;
  currentPage: number;
  pageCount: number;
  slideshowRunning: boolean;
  onSlideshowClose: () => void;
}

const ImageWall: React.FC<IImageWallProps> = ({
  images,
  onChangePage,
  currentPage,
  pageCount,
  onSlideshowClose,
  slideshowRunning = false,
}) => {
  const handleLightBoxPage = useCallback(
    (direction: number) => {
      if (direction === -1) {
        if (currentPage === 1) return false;
        onChangePage(currentPage - 1);
      } else {
        if (currentPage === pageCount) {
          // if the slideshow is running
          // return to the first page
          if (slideshowRunning) {
            onChangePage(0);
            return true;
          }
          return false;
        }
        onChangePage(currentPage + 1);
      }
      return direction === -1 || direction === 1;
    },
    [onChangePage, currentPage, pageCount, slideshowRunning]
  );

  const showLightbox = useLightbox({
    images,
    showNavigation: false,
    pageCallback: handleLightBoxPage,
    pageHeader: `Page ${currentPage} / ${pageCount}`,
    slideshowEnabled: slideshowRunning,
    onClose: onSlideshowClose,
  });

  useEffect(() => {
    if (slideshowRunning) {
      showLightbox(0, true);
    }
  }, [slideshowRunning, showLightbox]);

  const thumbs = images.map((image, index) => (
    <div
      role="link"
      tabIndex={index}
      key={image.id}
      onClick={() => showLightbox(index)}
      onKeyPress={() => showLightbox(index)}
    >
      <img
        src={image.paths.thumbnail ?? ""}
        loading="lazy"
        className="gallery-image"
        alt={image.title ?? TextUtils.fileNameFromPath(image.path)}
      />
    </div>
  ));

  return (
    <div className="gallery">
      <div className="flexbin">{thumbs}</div>
    </div>
  );
};

interface IImageList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  persistState?: PersistanceLevel;
  persistanceKey?: string;
  extraOperations?: IListHookOperation<FindImagesQueryResult>[];
}

export const ImageList: React.FC<IImageList> = ({
  filterHook,
  persistState,
  persistanceKey,
  extraOperations,
}) => {
  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);
  const [isSlideshowRunning, setIsSlideshowRunning] = useState(false);

  function startSlideshow(
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    result: GQL.FindImagesQueryResult,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    filter: ListFilterModel,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    selectedIds: Set<string>
  ) {
    setIsSlideshowRunning(true);
  }

  const onSlideshowClose = useCallback(() => {
    setIsSlideshowRunning(false);
  }, []);

  const otherOperations = (extraOperations ?? []).concat([
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
    {
      text: "Start Slideshow",
      onClick: startSlideshow,
      isDisplayed: showWhenDisplayModeWall,
      postRefetch: false,
    },
  ]);

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

  const { template, onSelectChange } = useImagesList({
    zoomable: true,
    selectable: true,
    otherOperations,
    renderContent,
    renderEditDialog: renderEditImagesDialog,
    renderDeleteDialog: renderDeleteImagesDialog,
    filterHook,
    addKeybinds,
    persistState,
    persistanceKey,
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
      if (singleResult.data.findImages.images.length === 1) {
        const { id } = singleResult!.data!.findImages!.images[0];
        // navigate to the image player page
        history.push(`/images/${id}`);
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
          <ExportDialog
            exportInput={{
              images: {
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
          onSelectChange(image.id, selected, shiftKey)
        }
      />
    );
  }

  function renderImages(
    result: FindImagesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number,
    onChangePage: (page: number) => void,
    pageCount: number
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
    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <ImageWall
          images={result.data.findImages.images}
          onChangePage={onChangePage}
          currentPage={filter.currentPage}
          pageCount={pageCount}
          slideshowRunning={isSlideshowRunning}
          onSlideshowClose={onSlideshowClose}
        />
      );
    }
  }

  function renderContent(
    result: FindImagesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number,
    onChangePage: (page: number) => void,
    pageCount: number
  ) {
    return (
      <>
        {maybeRenderImageExportDialog(selectedIds)}
        {renderImages(
          result,
          filter,
          selectedIds,
          zoomIndex,
          onChangePage,
          pageCount
        )}
      </>
    );
  }

  return template;
};
