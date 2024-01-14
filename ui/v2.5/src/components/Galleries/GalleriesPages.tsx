import React, { useMemo, useState } from "react";
import { PaginationIndex } from "../List/Pagination";
import { DisplayMode } from "src/models/list-filter/types";
import {
  FilterMode,
  FindGalleriesQueryResult,
} from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { queryFindGalleries, useFindGalleries } from "src/core/StashService";
import { useIntl } from "react-intl";
import cx from "classnames";
import { useListSelect } from "src/hooks/listSelect";
import { IItemListOperation } from "../List/ItemList";
import { FilterSidebar } from "../List/FilterSidebar";
import { ListHeader } from "../List/ListHeader";
import { Button, Table } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { ListOperationButtons } from "../List/ListOperationButtons";
import { ListOperationDropdown } from "../List/ListOperationDropdown";
import { faShuffle } from "@fortawesome/free-solid-svg-icons";
import { Link, useHistory } from "react-router-dom";
import { useModal } from "src/hooks/modal";
import { ExportDialog } from "../Shared/ExportDialog";
import { getFromIds } from "src/utils/data";
import { GalleryCard } from "./GalleryCard";
import { galleryTitle } from "src/core/galleries";
import GalleryWallCard from "./GalleryWallCard";
import { EditGalleriesDialog } from "./EditGalleriesDialog";
import { DeleteGalleriesDialog } from "./DeleteGalleriesDialog";
import { CollapseDivider } from "../Shared/CollapseDivider";

export const GalleriesPage: React.FC = ({}) => {
  const intl = useIntl();
  const history = useHistory();

  const [filter, setFilter] = useState<ListFilterModel>(
    () => new ListFilterModel(FilterMode.Galleries)
  );
  const [filterCollapsed, setFilterCollapsed] = useState(false);

  const result = useFindGalleries(filter);
  const items = result.data?.findGalleries.galleries ?? [];
  const { selectedIds, onSelectChange, onSelectAll, onSelectNone } =
    useListSelect(items);

  const { modal, showModal, closeModal } = useModal();

  const totalCount = useMemo(
    () => result.data?.findGalleries.count ?? 0,
    [result.data?.findGalleries.count]
  );

  function renderGalleries() {
    if (!result.data?.findGalleries) return;

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row justify-content-center">
          {items.map((gallery) => (
            <GalleryCard
              key={gallery.id}
              gallery={gallery}
              zoomIndex={filter.zoomIndex}
              selecting={selectedIds.size > 0}
              selected={selectedIds.has(gallery.id)}
              onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
                onSelectChange(gallery.id, selected, shiftKey)
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
            {items.map((gallery) => (
              <tr key={gallery.id}>
                <td>
                  <Link to={`/galleries/${gallery.id}`}>
                    {gallery.cover ? (
                      <img
                        loading="lazy"
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

  async function viewRandom() {
    if (items.length === 0) return;

    // query for a random image
    if (result.data?.findGalleries) {
      const { count } = result.data.findGalleries;

      const index = Math.floor(Math.random() * count);
      const filterCopy = filter.clone();
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

  async function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          galleries: {
            ids: Array.from(selectedIds.values()),
            all,
          },
        }}
        onClose={closeModal}
      />
    );
  }

  const otherOperations: IItemListOperation<FindGalleriesQueryResult>[] = [
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
    o: IItemListOperation<FindGalleriesQueryResult>
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
              <EditGalleriesDialog
                selected={getFromIds(items, selectedIds)}
                onClose={closeModal}
              />
            )
          }
          onDelete={() =>
            showModal(
              <DeleteGalleriesDialog
                selected={getFromIds(items, selectedIds)}
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
    <div id="galleries-page" className="list-page">
      {modal}

      {!filterCollapsed && (
        <FilterSidebar filter={filter} setFilter={(f) => setFilter(f)} />
      )}
      <CollapseDivider
        collapsed={filterCollapsed}
        setCollapsed={(v) => setFilterCollapsed(v)}
      />
      <div className={cx("list-page-results", { expanded: filterCollapsed })}>
        <ListHeader
          filter={filter}
          setFilter={setFilter}
          totalItems={totalCount}
          selectedIds={selectedIds}
          onSelectAll={onSelectAll}
          onSelectNone={onSelectNone}
          actionButtons={
            items.length > 0 && (
              <div>
                <Button
                  className="shuffle-galleries-button"
                  variant="secondary"
                  onClick={() => viewRandom()}
                >
                  <Icon icon={faShuffle} />
                </Button>
              </div>
            )
          }
          selectedButtons={renderButtons}
        />
        <div className="list-page-items">
          <PaginationIndex
            itemsPerPage={filter.itemsPerPage}
            currentPage={filter.currentPage}
            totalItems={totalCount}
          />
          {renderGalleries()}
        </div>
      </div>
    </div>
  );
};

export default GalleriesPage;
