import React, {
  useCallback,
  useState,
  useMemo,
  MouseEvent,
  useEffect,
} from "react";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import cloneDeep from "lodash-es/cloneDeep";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindImages,
  useFindImages,
  useFindImagesMetadata,
} from "src/core/StashService";
import { useFilteredItemList } from "../List/ItemList";
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
import { ImageCardGrid } from "./ImageCardGrid";
import { View } from "../List/views";
import {
  FilteredListToolbar,
  IItemListOperation,
} from "../List/FilteredListToolbar";
import { FileSize } from "../Shared/FileSize";
import { PatchComponent, PatchContainerComponent } from "src/patch";
import { GenerateDialog } from "../Dialogs/GenerateDialog";
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "../Shared/Sidebar";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/FilterSidebar";
import {
  IListFilterOperation,
  ListOperations,
} from "../List/ListOperationButtons";
import { FilterTags } from "../List/FilterTags";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { LoadedContent } from "../List/PagedList";
import useFocus from "src/utils/focus";
import cx from "classnames";
import { SidebarStudiosFilter } from "../List/Filters/StudiosFilter";
import { SidebarPerformersFilter } from "../List/Filters/PerformersFilter";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { Button } from "react-bootstrap";
import { OrganizedCriterionOption } from "src/models/list-filter/criteria/organized";
import { SidebarAgeFilter } from "../List/Filters/SidebarAgeFilter";
import { PerformerAgeCriterionOption } from "src/models/list-filter/images";

interface IImageWallProps {
  images: GQL.SlimImageDataFragment[];
  onChangePage: (page: number) => void;
  currentPage: number;
  pageCount: number;
  handleImageOpen: (index: number) => void;
  zoomIndex: number;
  selectedIds?: Set<string>;
  onSelectChange?: (id: string, selected: boolean, shiftKey: boolean) => void;
  selecting?: boolean;
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
  selectedIds,
  onSelectChange,
  selecting,
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
      const imageId = props.photo.key;
      if (!imageId) {
        return null;
      }
      return (
        <ImageWallItem
          {...props}
          maxHeight={maxHeight}
          selected={selectedIds?.has(imageId)}
          onSelectedChanged={
            onSelectChange
              ? (selected, shiftKey) =>
                  onSelectChange(imageId, selected, shiftKey)
              : undefined
          }
          selecting={selecting}
        />
      );
    },
    [targetRowHeight, selectedIds, onSelectChange, selecting]
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

const ImageList: React.FC<IImageListImages> = PatchComponent(
  "ImageList",
  ({
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
        <ImageCardGrid
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
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
          selecting={!!selectedIds && selectedIds.size > 0}
        />
      );
    }

    // should not happen
    return <></>;
  }
);

function renderMetadataByline(
  metadataInfo: GQL.FindImagesMetadataQueryResult | undefined
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

const ImageFilterSidebarSections = PatchContainerComponent(
  "FilteredImageList.SidebarSections"
);

const SidebarContent: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  sidebarOpen: boolean;
  onClose?: () => void;
  showEditFilter: (editingCriterion?: string) => void;
  count?: number;
  focus?: ReturnType<typeof useFocus>;
}> = ({
  filter,
  setFilter,
  filterHook,
  view,
  showEditFilter,
  sidebarOpen,
  onClose,
  count,
  focus,
}) => {
  const showResultsId =
    count !== undefined ? "actions.show_count_results" : "actions.show_results";

  const hideStudios = view === View.StudioScenes;

  return (
    <>
      <FilteredSidebarHeader
        sidebarOpen={sidebarOpen}
        showEditFilter={showEditFilter}
        filter={filter}
        setFilter={setFilter}
        view={view}
        focus={focus}
      />

      <ImageFilterSidebarSections>
        {!hideStudios && (
          <SidebarStudiosFilter
            filter={filter}
            setFilter={setFilter}
            filterHook={filterHook}
          />
        )}
        <SidebarPerformersFilter
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarTagsFilter
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        />
        <SidebarRatingFilter filter={filter} setFilter={setFilter} />
        <SidebarBooleanFilter
          title={<FormattedMessage id="organized" />}
          data-type={OrganizedCriterionOption.type}
          option={OrganizedCriterionOption}
          filter={filter}
          setFilter={setFilter}
          sectionID="organized"
        />
        <SidebarAgeFilter
          title={<FormattedMessage id="performer_age" />}
          option={PerformerAgeCriterionOption}
          filter={filter}
          setFilter={setFilter}
          sectionID="performer_age"
        />
      </ImageFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
    </>
  );
};

function useViewRandom(filter: ListFilterModel, count: number) {
  const history = useHistory();

  const viewRandom = useCallback(async () => {
    // query for a random image
    if (count === 0) {
      return;
    }

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
  }, [history, filter, count]);

  return viewRandom;
}

function useAddKeybinds(filter: ListFilterModel, count: number) {
  const viewRandom = useViewRandom(filter, count);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      viewRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [viewRandom]);
}

interface IImageList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  view?: View;
  alterQuery?: boolean;
  extraOperations?: IItemListOperation<GQL.FindImagesQueryResult>[];
  chapters?: GQL.GalleryChapterDataFragment[];
}

export const FilteredImageList = PatchComponent(
  "FilteredImageList",
  (props: IImageList) => {
    const intl = useIntl();

    const [slideshowRunning, setSlideshowRunning] = useState<boolean>(false);

    const searchFocus = useFocus();

    const withSidebar = props.view !== View.GalleryImages;

    const {
      filterHook,
      view,
      alterQuery,
      extraOperations: providedOperations = [],
      chapters,
    } = props;

    // States
    const {
      showSidebar,
      setShowSidebar,
      sectionOpen,
      setSectionOpen,
      loading: sidebarStateLoading,
    } = useSidebarState(view);

    const {
      filterState,
      queryResult,
      metadataInfo,
      modalState,
      listSelect,
      showEditFilter,
    } = useFilteredItemList({
      filterStateProps: {
        filterMode: GQL.FilterMode.Images,
        view,
        useURL: alterQuery,
      },
      queryResultProps: {
        useResult: useFindImages,
        useMetadataInfo: useFindImagesMetadata,
        getCount: (r) => r.data?.findImages.count ?? 0,
        getItems: (r) => r.data?.findImages.images ?? [],
        filterHook,
      },
    });

    const { filter, setFilter } = filterState;

    const { effectiveFilter, result, cachedResult, items, totalCount } =
      queryResult;

    const metadataByline = useMemo(() => {
      if (cachedResult.loading) return null;

      return renderMetadataByline(metadataInfo) ?? null;
    }, [cachedResult.loading, metadataInfo]);

    const {
      selectedIds,
      selectedItems,
      onSelectChange,
      onSelectAll,
      onSelectNone,
      onInvertSelection,
      hasSelection,
    } = listSelect;

    const { modal, showModal, closeModal } = modalState;

    // Utility hooks
    const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
      filter,
      setFilter,
    });

    useAddKeybinds(filter, totalCount);
    useFilteredSidebarKeybinds({
      showSidebar,
      setShowSidebar,
    });

    useEffect(() => {
      Mousetrap.bind("e", () => {
        if (hasSelection) {
          onEdit?.();
        }
      });

      Mousetrap.bind("d d", () => {
        if (hasSelection) {
          onDelete?.();
        }
      });

      return () => {
        Mousetrap.unbind("e");
        Mousetrap.unbind("d d");
      };
    });

    const onCloseEditDelete = useCloseEditDelete({
      closeModal,
      onSelectNone,
      result,
    });

    const viewRandom = useViewRandom(effectiveFilter, totalCount);

    function onExport(all: boolean) {
      showModal(
        <ExportDialog
          exportInput={{
            images: {
              ids: Array.from(selectedIds.values()),
              all: all,
            },
          }}
          onClose={() => closeModal()}
        />
      );
    }

    function onEdit() {
      showModal(
        <EditImagesDialog
          selected={selectedItems}
          onClose={onCloseEditDelete}
        />
      );
    }

    function onDelete() {
      showModal(
        <DeleteImagesDialog
          selected={selectedItems}
          onClose={onCloseEditDelete}
        />
      );
    }

    const convertedExtraOperations: IListFilterOperation[] =
      providedOperations.map((o) => ({
        ...o,
        isDisplayed: o.isDisplayed
          ? () => o.isDisplayed!(result, filter, selectedIds)
          : undefined,
        onClick: () => {
          o.onClick(result, filter, selectedIds);
        },
      }));

    const otherOperations: IListFilterOperation[] = [
      ...convertedExtraOperations,
      {
        text: intl.formatMessage({ id: "actions.select_all" }),
        onClick: () => onSelectAll(),
        isDisplayed: () => totalCount > 0,
      },
      {
        text: intl.formatMessage({ id: "actions.select_none" }),
        onClick: () => onSelectNone(),
        isDisplayed: () => hasSelection,
      },
      {
        text: intl.formatMessage({ id: "actions.invert_selection" }),
        onClick: () => onInvertSelection(),
        isDisplayed: () => totalCount > 0,
      },
      {
        text: intl.formatMessage({ id: "actions.view_random" }),
        onClick: viewRandom,
      },
      {
        text: `${intl.formatMessage({ id: "actions.generate" })}…`,
        onClick: () => {
          showModal(
            <GenerateDialog
              type="image"
              selectedIds={Array.from(selectedIds.values())}
              onClose={() => closeModal()}
            />
          );
        },
        isDisplayed: () => hasSelection,
      },
      {
        text: intl.formatMessage({ id: "actions.export" }),
        onClick: () => onExport(false),
        isDisplayed: () => hasSelection,
      },
      {
        text: intl.formatMessage({ id: "actions.export_all" }),
        onClick: () => onExport(true),
      },
    ];

    // render
    if (sidebarStateLoading) return null;

    const operations = (
      <ListOperations
        items={items.length}
        hasSelection={hasSelection}
        operations={otherOperations}
        onEdit={onEdit}
        onDelete={onDelete}
        operationsMenuClassName="image-list-operations-dropdown"
      />
    );

    const pageCount = Math.ceil(totalCount / filter.itemsPerPage);

    const content = (
      <>
        <FilteredListToolbar
          filter={filter}
          listSelect={listSelect}
          setFilter={setFilter}
          showEditFilter={showEditFilter}
          onDelete={onDelete}
          onEdit={onEdit}
          operationComponent={operations}
          view={view}
          zoomable
        />

        <FilterTags
          criteria={filter.criteria}
          onEditCriterion={(c) => showEditFilter(c.criterionOption.type)}
          onRemoveCriterion={removeCriterion}
          onRemoveAll={clearAllCriteria}
        />

        <div className="pagination-index-container">
          <Pagination
            currentPage={filter.currentPage}
            itemsPerPage={filter.itemsPerPage}
            totalItems={totalCount}
            onChangePage={(page) => setFilter(filter.changePage(page))}
          />
          <PaginationIndex
            loading={cachedResult.loading}
            itemsPerPage={filter.itemsPerPage}
            currentPage={filter.currentPage}
            totalItems={totalCount}
            metadataByline={metadataByline}
          />
        </div>

        <LoadedContent loading={result.loading} error={result.error}>
          <ImageList
            filter={filter}
            images={items}
            onChangePage={(page) => setFilter(filter.changePage(page))}
            onSelectChange={onSelectChange}
            pageCount={pageCount}
            selectedIds={selectedIds}
            slideshowRunning={slideshowRunning}
            setSlideshowRunning={setSlideshowRunning}
            chapters={chapters}
          />
        </LoadedContent>

        {totalCount > filter.itemsPerPage && (
          <div className="pagination-footer-container">
            <div className="pagination-footer">
              <Pagination
                itemsPerPage={filter.itemsPerPage}
                currentPage={filter.currentPage}
                totalItems={totalCount}
                onChangePage={setPage}
                pagePopupPlacement="top"
                metadataByline={metadataByline}
              />
            </div>
          </div>
        )}
      </>
    );

    if (!withSidebar) {
      return content;
    }

    return (
      <div
        className={cx("item-list-container image-list", {
          "hide-sidebar": !showSidebar,
        })}
      >
        {modal}

        <SidebarStateContext.Provider value={{ sectionOpen, setSectionOpen }}>
          <SidebarPane hideSidebar={!showSidebar}>
            <Sidebar hide={!showSidebar} onHide={() => setShowSidebar(false)}>
              <SidebarContent
                filter={filter}
                setFilter={setFilter}
                filterHook={filterHook}
                showEditFilter={showEditFilter}
                view={view}
                sidebarOpen={showSidebar}
                onClose={() => setShowSidebar(false)}
                count={cachedResult.loading ? undefined : totalCount}
                focus={searchFocus}
              />
            </Sidebar>
            <SidebarPaneContent
              onSidebarToggle={() => setShowSidebar(!showSidebar)}
            >
              {content}
            </SidebarPaneContent>
          </SidebarPane>
        </SidebarStateContext.Provider>
      </div>
    );
  }
);
