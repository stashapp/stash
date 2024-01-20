import React, { useCallback, useMemo, useState } from "react";
import { DisplayMode } from "src/models/list-filter/types";
import {
  FilterMode,
  FindImagesQueryResult,
  SlimImageDataFragment,
} from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { queryFindImages, useFindImages } from "src/core/StashService";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import TextUtils from "src/utils/text";
import { useListSelect } from "src/hooks/listSelect";
import { IItemListOperation } from "../List/ItemList";
import { ListOperationButtons } from "../List/ListOperationButtons";
import { ListOperationDropdown } from "../List/ListOperationDropdown";
import { useModal } from "src/hooks/modal";
import { ExportDialog } from "../Shared/ExportDialog";
import { getFromIds } from "src/utils/data";
import { ImageCard } from "./ImageCard";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { ImageWall } from "./ImageList";
import { useHistory } from "react-router-dom";
import { EditImagesDialog } from "./EditImagesDialog";
import { DeleteImagesDialog } from "./DeleteImagesDialog";
import { Icon } from "../Shared/Icon";
import { faShuffle } from "@fortawesome/free-solid-svg-icons";
import { ListPage } from "../List/ListPage";
import { useFilterURL, useInitialFilter, useResultCount } from "../List/util";
import DropdownItem from "react-bootstrap/esm/DropdownItem";

const filterMode = FilterMode.Images;

export const ImagesPageImpl: React.FC<{
  defaultFilter: ListFilterModel;
}> = ({ defaultFilter }) => {
  const intl = useIntl();
  const history = useHistory();

  const [slideshowRunning, setSlideshowRunning] = useState<boolean>(false);

  const [filter, setFilterState] = useState<ListFilterModel>(defaultFilter);

  const { setFilter } = useFilterURL(filter, setFilterState, defaultFilter);

  const result = useFindImages(filter);
  const { loading } = result;
  const images = useMemo(
    () => result.data?.findImages.images ?? [],
    [result.data?.findImages.images]
  );

  const totalCount = useResultCount(
    filter,
    loading,
    result.data?.findImages.count ?? 0
  );

  const listSelect = useListSelect(images);
  const { selectedIds, onSelectChange } = listSelect;

  const { modal, showModal, closeModal } = useModal();

  const pageCount = Math.ceil(totalCount / filter.itemsPerPage);

  const onChangePage = useCallback(
    (page: number) => {
      const newFilter = filter.clone();
      newFilter.currentPage = page;
      setFilter(newFilter);

      // if the current page has a detail-header, then
      // scroll up relative to that rather than 0, 0
      const detailHeader = document.querySelector(".detail-header");
      if (detailHeader) {
        window.scrollTo(0, detailHeader.scrollHeight - 50);
      } else {
        window.scrollTo(0, 0);
      }
    },
    [filter, setFilter]
  );

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

  const showLightbox = useLightbox(lightboxState, []);

  const metadataByline = useMemo(() => {
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
  }, [result]);

  const handleImageOpen = useCallback(
    (index) => {
      setSlideshowRunning(true);
      showLightbox(index, true);
    },
    [showLightbox, setSlideshowRunning]
  );

  function onPreview(index: number, ev: React.MouseEvent<Element, MouseEvent>) {
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

  function renderImages() {
    if (!result.data?.findImages) return;

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
  }

  async function viewRandom() {
    if (images.length === 0) return;

    // query for a random image
    if (result.data?.findImages) {
      const { count } = result.data.findImages;

      const index = Math.floor(Math.random() * count);
      const filterCopy = filter.clone();
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

  async function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          images: {
            ids: Array.from(selectedIds.values()),
            all,
          },
        }}
        onClose={closeModal}
      />
    );
  }

  const otherOperations: IItemListOperation<FindImagesQueryResult>[] = [
    {
      text: intl.formatMessage({ id: "actions.export" }),
      onClick: () => onExport(false),
    },
    {
      text: intl.formatMessage({ id: "actions.export_all" }),
      onClick: () => onExport(true),
    },
  ];

  async function onOperationClicked(
    o: IItemListOperation<FindImagesQueryResult>
  ) {
    await o.onClick(result, filter, selectedIds);
    if (o.postRefetch) {
      result.refetch();
    }
  }

  const operations = otherOperations?.map((o) => ({
    text: o.text,
    onClick: () => {
      onOperationClicked(o);
    },
    isDisplayed: () => {
      if (o.isDisplayed) {
        return o.isDisplayed(result, filter, selectedIds);
      }

      return true;
    },
    icon: o.icon,
    buttonVariant: o.buttonVariant,
  }));

  function renderButtons() {
    return (
      <div>
        <ListOperationButtons
          itemsSelected
          onEdit={() =>
            showModal(
              <EditImagesDialog
                selected={getFromIds(images, selectedIds)}
                onClose={closeModal}
              />
            )
          }
          onDelete={() =>
            showModal(
              <DeleteImagesDialog
                selected={getFromIds(images, selectedIds)}
                onClose={closeModal}
              />
            )
          }
        />

        <ListOperationDropdown operations={operations} />
      </div>
    );
  }

  return (
    <>
      <ListPage
        id="images-page"
        loading={loading}
        filter={filter}
        setFilter={(f) => setFilter(f)}
        listSelect={listSelect}
        actionButtons={
          images.length > 0 && (
            <>
              <DropdownItem
                className="shuffle-images-button"
                onClick={() => viewRandom()}
              >
                <Icon icon={faShuffle} />
                <span>
                  <FormattedMessage id="actions.shuffle" />
                </span>
              </DropdownItem>
            </>
          )
        }
        selectedButtons={renderButtons}
        metadataByline={metadataByline}
        totalCount={totalCount}
      >
        {renderImages()}
      </ListPage>
      {modal}
    </>
  );
};

export const ImagesPage: React.FC = () => {
  const initialFilter = useInitialFilter(filterMode);

  if (!initialFilter) {
    return null;
  }

  return <ImagesPageImpl defaultFilter={initialFilter} />;
};

export default ImagesPage;
