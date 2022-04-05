import React, { useCallback, useState, useMemo, MouseEvent } from "react";
import { useIntl } from "react-intl";
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
  handleImageOpen: (index: number) => void;
}

const ImageWall: React.FC<IImageWallProps> = ({ images, handleImageOpen }) => {
  const thumbs = images.map((image, index) => (
    <div
      role="link"
      tabIndex={index}
      key={image.id}
      onClick={() => handleImageOpen(index)}
      onKeyPress={() => handleImageOpen(index)}
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

interface IImageListImages {
  images: SlimImageDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onChangePage: (page: number) => void;
  pageCount: number;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  slideshowRunning: boolean;
  setSlideshowRunning: (running: boolean) => void;
}

const ImageListImages: React.FC<IImageListImages> = ({
  images,
  filter,
  selectedIds,
  onChangePage,
  pageCount,
  onSelectChange,
  slideshowRunning,
  setSlideshowRunning,
}) => {
  const handleLightBoxPage = useCallback(
    (direction: number) => {
      if (direction === -1) {
        if (filter.currentPage === 1) {
          onChangePage(pageCount);
        } else {
          onChangePage(filter.currentPage - 1);
        }
      } else if (direction === 1) {
        if (filter.currentPage === pageCount) {
          // return to the first page
          onChangePage(1);
        } else {
          onChangePage(filter.currentPage + 1);
        }
      }
    },
    [onChangePage, filter.currentPage, pageCount]
  );

  const handleClose = useCallback(() => {
    setSlideshowRunning(false);
  }, [setSlideshowRunning]);

  const lightboxState = useMemo(() => {
    return {
      images,
      showNavigation: false,
      pageCallback: pageCount > 1 ? handleLightBoxPage : undefined,
      pageHeader: `Page ${filter.currentPage} / ${pageCount}`,
      slideshowEnabled: slideshowRunning,
      onClose: handleClose,
    };
  }, [
    images,
    pageCount,
    filter.currentPage,
    slideshowRunning,
    handleClose,
    handleLightBoxPage,
  ]);

  const showLightbox = useLightbox(lightboxState);

  const handleImageOpen = useCallback(
    (index) => {
      setSlideshowRunning(true);
      showLightbox(index, true);
    },
    [showLightbox, setSlideshowRunning]
  );

  function onPreview(index: number, ev: MouseEvent) {
    handleImageOpen(index);
    ev.preventDefault();
  }

  function renderImageCard(
    index: number,
    image: SlimImageDataFragment,
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
        onPreview={
          selectedIds.size < 1 ? (ev) => onPreview(index, ev) : undefined
        }
      />
    );
  }

  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <div className="row justify-content-center">
        {images.map((image, index) =>
          renderImageCard(index, image, filter.zoomIndex)
        )}
      </div>
    );
  }
  if (filter.displayMode === DisplayMode.Wall) {
    return (
      <ImageWall
        images={images}
        onChangePage={onChangePage}
        currentPage={filter.currentPage}
        pageCount={pageCount}
        handleImageOpen={handleImageOpen}
      />
    );
  }

  // should not happen
  return <></>;
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
  const intl = useIntl();
  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);
  const [slideshowRunning, setSlideshowRunning] = useState<boolean>(false);

  const otherOperations = (extraOperations ?? []).concat([
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

  function selectChange(id: string, selected: boolean, shiftKey: boolean) {
    onSelectChange(id, selected, shiftKey);
  }

  function renderImages(
    result: FindImagesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onChangePage: (page: number) => void,
    pageCount: number
  ) {
    if (!result.data || !result.data.findImages) {
      return;
    }

    return (
      <ImageListImages
        filter={filter}
        images={result.data.findImages.images}
        onChangePage={onChangePage}
        onSelectChange={selectChange}
        pageCount={pageCount}
        selectedIds={selectedIds}
        slideshowRunning={slideshowRunning}
        setSlideshowRunning={setSlideshowRunning}
      />
    );
  }

  function renderContent(
    result: FindImagesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onChangePage: (page: number) => void,
    pageCount: number
  ) {
    return (
      <>
        {maybeRenderImageExportDialog(selectedIds)}
        {renderImages(result, filter, selectedIds, onChangePage, pageCount)}
      </>
    );
  }

  return template;
};
