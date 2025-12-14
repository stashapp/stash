import React, { useCallback, useState, useMemo, MouseEvent } from "react";
import { FormattedNumber, useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindImages,
  useFindImages,
  useFindImagesMetadata,
} from "src/core/StashService";
import { ItemList, ItemListContext, showWhenSelected } from "../List/ItemList";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";

import { ImageWallItem } from "./ImageWallItem";
import { EditImagesDialog } from "./EditImagesDialog";
import { DeleteImagesDialog } from "./DeleteImagesDialog";
import "flexbin/flexbin.css";
import Gallery, { RenderImageProps } from "react-photo-gallery";
import { ExportDialog } from "../Shared/ExportDialog";
import { objectTitle } from "src/core/files";
import { useConfigurationContext } from "src/hooks/Config";
import { ImageGridCard } from "./ImageGridCard";
import { View } from "../List/views";
import { IItemListOperation } from "../List/FilteredListToolbar";
import { FileSize } from "../Shared/FileSize";
import { PatchComponent } from "src/patch";

interface IImageWallProps {
  images: GQL.SlimImageDataFragment[];
  onChangePage: (page: number) => void;
  currentPage: number;
  pageCount: number;
  handleImageOpen: (index: number) => void;
  zoomIndex: number;
}

const zoomWidths = [280, 340, 480, 640];
const breakpointZoomHeights = [
  { minWidth: 576, heights: [100, 120, 240, 360] },
  { minWidth: 768, heights: [120, 160, 240, 480] },
  { minWidth: 1200, heights: [120, 160, 240, 300] },
  { minWidth: 1400, heights: [160, 240, 300, 480] },
];

const ImageWall: React.FC<IImageWallProps> = ({
  images,
  zoomIndex,
  handleImageOpen,
}) => {
  const { configuration } = useConfigurationContext();
  const uiConfig = configuration?.ui;

  const containerRef = React.useRef<HTMLDivElement>(null);

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
      width: image.visual_files?.[0]?.width ?? 0,
      height: image.visual_files?.[0]?.height ?? 0,
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
    let preferredSize = zoomWidths[zoomIndex];
    let columnCount = containerWidth / preferredSize;
    return Math.round(columnCount);
  }

  const targetRowHeight = useCallback(
    (containerWidth: number) => {
      let zoomHeight = 280;
      breakpointZoomHeights.forEach((e) => {
        if (containerWidth >= e.minWidth) {
          zoomHeight = e.heights[zoomIndex];
        }
      });
      return zoomHeight;
    },
    [zoomIndex]
  );

  // set the max height as a factor of the targetRowHeight
  // this allows some images to be taller than the target row height
  // but prevents images from becoming too tall when there is a small number of items
  const maxHeightFactor = 1.3;

  const renderImage = useCallback(
    (props: RenderImageProps) => {
      // #6165 - only use targetRowHeight in row direction
      const maxHeight =
        props.direction === "column"
          ? props.photo.height
          : targetRowHeight(containerRef.current?.offsetWidth ?? 0) *
            maxHeightFactor;
      return <ImageWallItem {...props} maxHeight={maxHeight} />;
    },
    [targetRowHeight]
  );

  return (
    <div className="gallery" ref={containerRef}>
      {photos.length ? (
        <Gallery
          photos={photos}
          renderImage={renderImage}
          onClick={showLightboxOnClick}
          margin={uiConfig?.imageWallOptions?.margin!}
          direction={uiConfig?.imageWallOptions?.direction!}
          columns={columns}
          targetRowHeight={targetRowHeight}
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
      showLightbox({ initialIndex: index, slideshowEnabled: true });
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
        zoomIndex={filter.zoomIndex}
      />
    );
  }

  // should not happen
  return <></>;
};

function getItems(result: GQL.FindImagesQueryResult) {
  return result?.data?.findImages?.images ?? [];
}

function getCount(result: GQL.FindImagesQueryResult) {
  return result?.data?.findImages?.count ?? 0;
}

function renderMetadataByline(
  result: GQL.FindImagesQueryResult,
  metadataInfo?: GQL.FindImagesMetadataQueryResult
) {
  const megapixels = metadataInfo?.data?.findImages?.megapixels;
  const size = metadataInfo?.data?.findImages?.filesize;

  if (metadataInfo?.loading) {
    // return ellipsis
    return <span className="images-stats">&nbsp;(...)</span>;
  }

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
      {size ? (
        <span className="images-size">
          <FileSize size={size} />
        </span>
      ) : undefined}
      )
    </span>
  );
}

interface IImageList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
  extraOperations?: IItemListOperation<GQL.FindImagesQueryResult>[];
  chapters?: GQL.GalleryChapterDataFragment[];
}

export const ImageList: React.FC<IImageList> = PatchComponent(
  "ImageList",
  ({ filterHook, view, alterQuery, extraOperations = [], chapters = [] }) => {
    const intl = useIntl();
    const history = useHistory();
    const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
    const [isExportAll, setIsExportAll] = useState(false);
    const [slideshowRunning, setSlideshowRunning] = useState<boolean>(false);

    const filterMode = GQL.FilterMode.Images;

    const otherOperations = [
      ...extraOperations,
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
      onSelectChange: (
        id: string,
        selected: boolean,
        shiftKey: boolean
      ) => void,
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
      <ItemListContext
        filterMode={filterMode}
        useResult={useFindImages}
        useMetadataInfo={useFindImagesMetadata}
        getItems={getItems}
        getCount={getCount}
        alterQuery={alterQuery}
        filterHook={filterHook}
        view={view}
        selectable
      >
        <ItemList
          view={view}
          otherOperations={otherOperations}
          addKeybinds={addKeybinds}
          renderContent={renderContent}
          renderEditDialog={renderEditDialog}
          renderDeleteDialog={renderDeleteDialog}
          renderMetadataByline={renderMetadataByline}
        />
      </ItemListContext>
    );
  }
);
