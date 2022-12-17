import React, { useCallback, useState, useMemo, MouseEvent } from "react";
import { FormattedNumber, useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindImages, useFindImages } from "src/core/StashService";
import {
  makeItemList,
  IItemListOperation,
  PersistanceLevel,
  showWhenSelected,
} from "../List/ItemList";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";

import { ImageCard } from "./ImageCard";
import { EditImagesDialog } from "./EditImagesDialog";
import { DeleteImagesDialog } from "./DeleteImagesDialog";
import "flexbin/flexbin.css";
import { ExportDialog } from "../Shared/ExportDialog";
import { objectTitle } from "src/core/files";
import TextUtils from "src/utils/text";

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
        alt={objectTitle(image)}
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
  images: GQL.SlimImageDataFragment[];
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
    image: GQL.SlimImageDataFragment,
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

const ImageItemList = makeItemList({
  filterMode: GQL.FilterMode.Images,
  useResult: useFindImages,
  getItems(result: GQL.FindImagesQueryResult) {
    return result?.data?.findImages?.images ?? [];
  },
  getCount(result: GQL.FindImagesQueryResult) {
    return result?.data?.findImages?.count ?? 0;
  },
  renderMetadataByline(result: GQL.FindImagesQueryResult) {
    const megapixels = result?.data?.findImages?.megapixels;
    const size = result?.data?.findImages?.filesize;
    const filesize = size ? TextUtils.fileSize(size) : undefined;

    if (!megapixels && !size) {
      return;
    }

    const separator = megapixels && size ? " - " : "";

    return (
      <span className="images-stats">
        &nbsp;(
        {megapixels ? (
          <span className="images-megapixels">
            <FormattedNumber value={megapixels} /> Megapixels
          </span>
        ) : undefined}
        {separator}
        {size && filesize ? (
          <span className="images-size">
            <FormattedNumber
              value={filesize.size}
              maximumFractionDigits={TextUtils.fileSizeFractionalDigits(
                filesize.unit
              )}
            />
            {` ${TextUtils.formatFileSizeUnit(filesize.unit)}`}
          </span>
        ) : undefined}
        )
      </span>
    );
  },
});

interface IImageList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  persistState?: PersistanceLevel;
  persistanceKey?: string;
  extraOperations?: IItemListOperation<GQL.FindImagesQueryResult>[];
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

  const otherOperations = [
    ...(extraOperations ?? []),
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

  function addKeybinds(
    result: GQL.FindImagesQueryResult,
    filter: ListFilterModel
  ) {
    Mousetrap.bind("p r", () => {
      viewRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }

  async function viewRandom(
    result: GQL.FindImagesQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random image
    if (result.data?.findImages) {
      const { count } = result.data.findImages;

      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindImages(filterCopy);
      if (singleResult.data.findImages.images.length === 1) {
        const { id } = singleResult.data.findImages.images[0];
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

  function renderContent(
    result: GQL.FindImagesQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void,
    onChangePage: (page: number) => void,
    pageCount: number
  ) {
    function maybeRenderImageExportDialog() {
      if (isExportDialogOpen) {
        return (
          <ExportDialog
            exportInput={{
              images: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => setIsExportDialogOpen(false)}
          />
        );
      }
    }

    function renderImages() {
      if (!result.data?.findImages) return;

      return (
        <ImageListImages
          filter={filter}
          images={result.data.findImages.images}
          onChangePage={onChangePage}
          onSelectChange={onSelectChange}
          pageCount={pageCount}
          selectedIds={selectedIds}
          slideshowRunning={slideshowRunning}
          setSlideshowRunning={setSlideshowRunning}
        />
      );
    }

    return (
      <>
        {maybeRenderImageExportDialog()}
        {renderImages()}
      </>
    );
  }

  function renderEditDialog(
    selectedImages: GQL.SlimImageDataFragment[],
    onClose: (applied: boolean) => void
  ) {
    return <EditImagesDialog selected={selectedImages} onClose={onClose} />;
  }

  function renderDeleteDialog(
    selectedImages: GQL.SlimImageDataFragment[],
    onClose: (confirmed: boolean) => void
  ) {
    return <DeleteImagesDialog selected={selectedImages} onClose={onClose} />;
  }

  return (
    <ImageItemList
      zoomable
      selectable
      filterHook={filterHook}
      persistState={persistState}
      persistanceKey={persistanceKey}
      otherOperations={otherOperations}
      addKeybinds={addKeybinds}
      renderContent={renderContent}
      renderEditDialog={renderEditDialog}
      renderDeleteDialog={renderDeleteDialog}
    />
  );
};
