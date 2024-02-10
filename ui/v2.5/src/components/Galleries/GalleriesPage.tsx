import React, { useCallback, useEffect, useMemo, useState } from "react";
import { DisplayMode } from "src/models/list-filter/types";
import {
  FilterMode,
  FindGalleriesQueryResult,
} from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { queryFindGalleries, useFindGalleries } from "src/core/StashService";
import { FormattedMessage, useIntl } from "react-intl";
import { useListSelect } from "src/hooks/listSelect";
import { IItemListOperation } from "../List/ItemList";
import { Icon } from "../Shared/Icon";
import { ListOperationButtons } from "../List/ListOperationButtons";
import { ListOperationDropdown } from "../List/ListOperationDropdown";
import { faShuffle } from "@fortawesome/free-solid-svg-icons";
import { useHistory } from "react-router-dom";
import { useModal } from "src/hooks/modal";
import { ExportDialog } from "../Shared/ExportDialog";
import { getFromIds } from "src/utils/data";
import { GalleryCard } from "./GalleryCard";
import GalleryWallCard from "./GalleryWallCard";
import { EditGalleriesDialog } from "./EditGalleriesDialog";
import { DeleteGalleriesDialog } from "./DeleteGalleriesDialog";
import { ListPage } from "../List/ListPage";
import {
  useFilterURL,
  useInitialFilter,
  useLocalFilterState,
  useResultCount,
} from "../List/util";
import { GalleryListTable } from "./GalleryListTable";
import DropdownItem from "react-bootstrap/esm/DropdownItem";
import Mousetrap from "mousetrap";

const filterMode = FilterMode.Galleries;
const pageView = "galleries";

export const GalleriesPageImpl: React.FC<{
  defaultFilter: ListFilterModel;
  defaultSidebarCollapsed?: boolean;
}> = ({ defaultFilter, defaultSidebarCollapsed }) => {
  const intl = useIntl();
  const history = useHistory();

  const [filter, setFilterState] = useState<ListFilterModel>(defaultFilter);

  const { setFilter } = useFilterURL(filter, setFilterState, defaultFilter);

  const result = useFindGalleries(filter);
  const { loading } = result;
  const items = useMemo(
    () => result.data?.findGalleries.galleries ?? [],
    [result.data?.findGalleries.galleries]
  );

  const listSelect = useListSelect(items);
  const { selectedIds, onSelectChange } = listSelect;

  const { modal, showModal, closeModal } = useModal();

  const totalCount = useResultCount(
    filter,
    loading,
    result.data?.findGalleries.count ?? 0
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
        <GalleryListTable
          galleries={items}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <div className="row">
          <div className="GalleryWall">
            {items.map((gallery) => (
              <GalleryWallCard key={gallery.id} gallery={gallery} />
            ))}
          </div>
        </div>
      );
    }
  }

  const viewRandom = useCallback(async () => {
    if (totalCount === 0) return;

    const randomFilter = filter.randomSingle(totalCount);
    const singleResult = await queryFindGalleries(randomFilter);

    if (singleResult.data.findGalleries.galleries.length === 1) {
      const { id } = singleResult.data.findGalleries.galleries[0];
      // navigate to the image player page
      history.push(`/galleries/${id}`);
    }
  }, [totalCount, filter, history]);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      viewRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [viewRandom]);

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
      <>
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
      </>
    );
  }

  return (
    <>
      <ListPage
        id="galleries-page"
        loading={loading}
        filter={filter}
        setFilter={(f) => setFilter(f)}
        listSelect={listSelect}
        pageView={pageView}
        initialSidebarCollapsed={defaultSidebarCollapsed}
        actionButtons={
          <>
            {/* <div>
              <Link to="/galleries/new">
                <Button variant="primary">
                  <FormattedMessage id="new" defaultMessage="New" />
                </Button>
              </Link>
            </div> */}
            {items.length > 0 && (
              <>
                <DropdownItem
                  className="shuffle-galleries-button"
                  onClick={() => viewRandom()}
                >
                  <Icon icon={faShuffle} />
                  <span>
                    <FormattedMessage id="actions.shuffle" />
                  </span>
                </DropdownItem>
              </>
            )}
          </>
        }
        selectedButtons={renderButtons}
        totalCount={totalCount}
      >
        {renderGalleries()}
      </ListPage>
      {modal}
    </>
  );
};

export const GalleriesPage: React.FC = () => {
  const initialFilter = useInitialFilter(filterMode);
  const initialLocalFilterState = useLocalFilterState(pageView, filterMode);

  if (!initialFilter || !initialLocalFilterState) {
    return null;
  }

  const initialSidebarCollapsed = initialLocalFilterState.sidebarCollapsed;

  return (
    <GalleriesPageImpl
      defaultFilter={initialFilter}
      defaultSidebarCollapsed={initialSidebarCollapsed}
    />
  );
};

export default GalleriesPage;
