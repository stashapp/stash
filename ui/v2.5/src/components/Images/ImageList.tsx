import React, {
  useCallback,
  useState,
  useMemo,
  MouseEvent,
  useContext,
} from "react";
import { FormattedNumber, useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindImages, useFindImages } from "src/core/StashService";
import {
  makeItemList,
  IItemListOperation,
  showWhenSelected,
} from "../List/ItemList";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";

import { ImageWallItem } from "./ImageWallItem";
import { EditImagesDialog } from "./EditImagesDialog";
import { DeleteImagesDialog } from "./DeleteImagesDialog";
import "flexbin/flexbin.css";
import Gallery from "react-photo-gallery";
import { ExportDialog } from "../Shared/ExportDialog";
import { objectTitle } from "src/core/files";
import TextUtils from "src/utils/text";
import { ConfigurationContext } from "src/hooks/Config";
import { ImageGridCard } from "./ImageGridCard";
import { View } from "../List/views";

interface IImageWallProps {
  images: GQL.SlimImageDataFragment[];
  onChangePage: (page: number) => void;
  currentPage: number;
  pageCount: number;
  handleImageOpen: (index: number) => void;
}

const ImageWall: React.FC<IImageWallProps> = ({ images, handleImageOpen }) => {
  const { configuration } = useContext(ConfigurationContext);
  const uiConfig = configuration?.ui;

  let photos: {
    src: string;
    srcSet?: string | string[] | undefined;
    sizes?: string | string[] | undefined;
    width: number;
    height: number;
    alt?: string | undefined;
    key?: string | undefined;
  }[] = [];

  images.forEach((image, index) => {
    let imageData = {
      src:
        image.paths.preview != ""
          ? image.paths.preview!
          : image.paths.thumbnail!,
      width: image.visual_files[0].width,
      height: image.visual_files[0].height,
      tabIndex: index,
      key: image.id,
      loading: "lazy",
      className: "gallery-image",
      alt: objectTitle(image),
    };
    photos.push(imageData);
  });

  const showLightboxOnClick = useCallback(
    (event, { index }) => {
      handleImageOpen(index);
    },
    [handleImageOpen]
  );

  function columns(containerWidth: number) {
    let preferredSize = 300;
    let columnCount = containerWidth / preferredSize;
    return Math.round(columnCount);
  }

  return (
    <div className="gallery">
      {photos.length ? (
        <Gallery
          photos={photos}
          renderImage={ImageWallItem}
          onClick={showLightboxOnClick}
          margin={uiConfig?.imageWallOptions?.margin!}
          direction={uiConfig?.imageWallOptions?.direction!}
          columns={columns}
        />
      ) : null}
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
  chapters?: GQL.GalleryChapterDataFragment[];
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
  chapters = [],
}) => {
  const handleLightBoxPage = useCallback(
    (props: { direction?: number; page?: number }) => {
      const { direction, page: newPage } = props;

      if (direction !== undefined) {
        if (direction < 0) {
          if (filter.currentPage === 1) {
            onChangePage(pageCount);
          } else {
            onChangePage(filter.currentPage + direction);
          }
        } else if (direction > 0) {
          if (filter.currentPage === pageCount) {
            // return to the first page
            onChangePage(1);
          } else {
            onChangePage(filter.currentPage + direction);
          }
        }
      } else if (newPage !== undefined) {
        onChangePage(newPage);
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
      page: filter.currentPage,
      pages: pageCount,
      pageSize: filter.itemsPerPage,
      slideshowEnabled: slideshowRunning,
      onClose: handleClose,
    };
  }, [
    images,
    pageCount,
    filter.currentPage,
    filter.itemsPerPage,
    slideshowRunning,
    handleClose,
    handleLightBoxPage,
  ]);

  const showLightbox = useLightbox(
    lightboxState,
    filter.sortBy === "path" &&
      filter.sortDirection === GQL.SortDirectionEnum.Asc
      ? chapters
      : []
  );

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

  if (filter.displayMode === DisplayMode.Grid) {
    return (
      <ImageGridCard
        images={images}
        selectedIds={selectedIds}
        zoomIndex={filter.zoomIndex}
        onSelectChange={onSelectChange}
        onPreview={onPreview}
      />
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
  view?: View;
  alterQuery?: boolean;
  extraOperations?: IItemListOperation<GQL.FindImagesQueryResult>[];
  chapters?: GQL.GalleryChapterDataFragment[];
}

export const ImageList: React.FC<IImageList> = ({
  filterHook,
  view,
  alterQuery,
  extraOperations,
  chapters = [],
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
          chapters={chapters}
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
      view={view}
      alterQuery={alterQuery}
      otherOperations={otherOperations}
      addKeybinds={addKeybinds}
      renderContent={renderContent}
      renderEditDialog={renderEditDialog}
      renderDeleteDialog={renderDeleteDialog}
    />
  );
};
