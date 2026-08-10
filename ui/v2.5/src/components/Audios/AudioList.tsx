import React, { useCallback, useEffect, useMemo } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { queryFindAudios, useFindAudios } from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { AudioListTable } from "./AudioListTable";
import { EditAudiosDialog } from "./EditAudiosDialog";
import { DeleteAudiosDialog } from "./DeleteAudiosDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { AudioCardGrid } from "./AudioCardGrid";
import { AudioMergeModal } from "./AudioMergeDialog";
import { objectTitle } from "src/core/files";
import TextUtils from "src/utils/text";
import { View } from "../List/views";
import { FileSize } from "../Shared/FileSize";
import { LoadedContent } from "../List/PagedList";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import { ListOperations } from "../List/ListOperationButtons";
import { useFilteredItemList } from "../List/ItemList";
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "../Shared/Sidebar";
import { SidebarPerformersFilter } from "../List/Filters/PerformersFilter";
import { SidebarStudiosFilter } from "../List/Filters/StudiosFilter";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import cx from "classnames";
import { SidebarRatingFilter } from "../List/Filters/RatingFilter";
import { OrganizedCriterionOption } from "src/models/list-filter/criteria/organized";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { AudioPerformerAgeCriterionOption } from "src/models/list-filter/audios";
import { SidebarAgeFilter } from "../List/Filters/SidebarAgeFilter";
import { SidebarDurationFilter } from "../List/Filters/SidebarDurationFilter";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/FilterSidebar";
import { PatchComponent, PatchContainerComponent } from "src/patch";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { Button } from "react-bootstrap";
import useFocus from "src/utils/focus";
import { useZoomKeybinds } from "../List/ZoomSlider";
import { FilteredListToolbar } from "../List/FilteredListToolbar";
import { FilterTags } from "../List/FilterTags";
import { SidebarFolderFilter } from "../List/Filters/FolderFilter";

function renderMetadataByline(result: GQL.FindAudiosQueryResult) {
  const duration = result?.data?.findAudios?.duration;
  const size = result?.data?.findAudios?.filesize;

  if (!duration && !size) {
    return;
  }

  const separator = duration && size ? " - " : "";

  return (
    <span className="audios-stats">
      &nbsp;(
      {duration ? (
        <span className="audios-duration">
          {TextUtils.secondsAsTimeString(duration, 3)}
        </span>
      ) : undefined}
      {separator}
      {size ? (
        <span className="audios-size">
          <FileSize size={size} />
        </span>
      ) : undefined}
      )
    </span>
  );
}

// audio has no play queue - playing simply navigates to the audio's page
function usePlayRandom(filter: ListFilterModel, count: number) {
  const history = useHistory();

  const playRandom = useCallback(async () => {
    if (count === 0) {
      return;
    }

    const pages = Math.ceil(count / filter.itemsPerPage);
    const page = Math.floor(Math.random() * pages) + 1;

    const indexMax = Math.min(filter.itemsPerPage, count);
    const index = Math.floor(Math.random() * indexMax);
    const filterCopy = cloneDeep(filter);
    filterCopy.currentPage = page;
    filterCopy.sortBy = "random";
    const queryResults = await queryFindAudios(filterCopy);
    const audio = queryResults.data.findAudios.audios[index];
    if (audio) {
      history.push(`/audios/${audio.id}`);
    }
  }, [filter, count, history]);

  return playRandom;
}

function useAddKeybinds(filter: ListFilterModel, count: number) {
  const playRandom = usePlayRandom(filter, count);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      playRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [playRandom]);
}

const AudioList: React.FC<{
  audios: GQL.SlimAudioDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
}> = PatchComponent(
  "AudioList",
  ({ audios, filter, selectedIds, onSelectChange, fromGroupId }) => {
    if (audios.length === 0) {
      return null;
    }

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <AudioCardGrid
          audios={audios}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
          fromGroupId={fromGroupId}
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <AudioListTable
          audios={audios}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }

    return null;
  }
);

const AudiosFilterSidebarSections = PatchContainerComponent(
  "FilteredAudioList.SidebarSections"
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

  const hideStudios = view === View.StudioAudios;

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

      <AudiosFilterSidebarSections>
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
        <SidebarDurationFilter filter={filter} setFilter={setFilter} />
        <SidebarFolderFilter
          text={<FormattedMessage id="folder" />}
          filter={filter}
          setFilter={setFilter}
          sectionID="folder"
        />
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
          option={AudioPerformerAgeCriterionOption}
          filter={filter}
          setFilter={setFilter}
          sectionID="performer_age"
        />
      </AudiosFilterSidebarSections>

      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
    </>
  );
};

interface IFilteredAudios {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
  fromGroupId?: string;
}

export const FilteredAudioList = PatchComponent(
  "FilteredAudioList",
  (props: IFilteredAudios) => {
    const intl = useIntl();
    const history = useHistory();

    const searchFocus = useFocus();

    const { filterHook, defaultSort, view, alterQuery, fromGroupId } = props;

    const {
      showSidebar,
      setShowSidebar,
      loading: sidebarStateLoading,
      sectionOpen,
      setSectionOpen,
    } = useSidebarState(view);

    const { filterState, queryResult, modalState, listSelect, showEditFilter } =
      useFilteredItemList({
        filterStateProps: {
          filterMode: GQL.FilterMode.Audios,
          defaultSort,
          view,
          useURL: alterQuery,
        },
        queryResultProps: {
          useResult: useFindAudios,
          getCount: (r) => r.data?.findAudios.count ?? 0,
          getItems: (r) => r.data?.findAudios.audios ?? [],
          filterHook,
        },
      });

    const { filter, setFilter } = filterState;

    const { effectiveFilter, result, cachedResult, items, totalCount } =
      queryResult;

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

    const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
      filter,
      setFilter,
    });

    useAddKeybinds(effectiveFilter, totalCount);
    useFilteredSidebarKeybinds({
      showSidebar,
      setShowSidebar,
    });

    const onCloseEditDelete = useCloseEditDelete({
      closeModal,
      onSelectNone,
      result,
    });

    const onEdit = useCallback(() => {
      showModal(
        <EditAudiosDialog
          selected={selectedItems}
          onClose={onCloseEditDelete}
        />
      );
    }, [showModal, selectedItems, onCloseEditDelete]);

    const onDelete = useCallback(() => {
      showModal(
        <DeleteAudiosDialog
          selected={selectedItems}
          onClose={onCloseEditDelete}
        />
      );
    }, [showModal, selectedItems, onCloseEditDelete]);

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
    }, [hasSelection, onEdit, onDelete]);

    useZoomKeybinds({
      zoomIndex: filter.zoomIndex,
      onChangeZoom: (zoom) => setFilter(filter.setZoom(zoom)),
    });

    const metadataByline = useMemo(() => {
      if (cachedResult.loading) return null;

      return renderMetadataByline(cachedResult) ?? null;
    }, [cachedResult]);

    const playRandom = usePlayRandom(effectiveFilter, totalCount);

    function onPlay() {
      if (items.length === 0) {
        return;
      }

      const audioID = hasSelection
        ? Array.from(selectedIds.values())[0]
        : items[0].id;
      history.push(`/audios/${audioID}`);
    }

    function onExport(all: boolean) {
      showModal(
        <ExportDialog
          exportInput={{
            audios: {
              ids: Array.from(selectedIds.values()),
              all: all,
            },
          }}
          onClose={() => closeModal()}
        />
      );
    }

    function onMerge() {
      const selected =
        selectedItems.map((s) => {
          return {
            id: s.id,
            title: objectTitle(s),
          };
        }) ?? [];
      showModal(
        <AudioMergeModal
          audios={selected}
          onClose={(mergedID?: string) => {
            closeModal();
            if (mergedID) {
              history.push(`/audios/${mergedID}`);
            }
          }}
          show
        />
      );
    }

    const otherOperations = [
      {
        text: intl.formatMessage({ id: "actions.play" }),
        onClick: () => onPlay(),
        isDisplayed: () => items.length > 0,
        className: "play-item",
      },
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
        text: intl.formatMessage({ id: "actions.play_random" }),
        onClick: playRandom,
        isDisplayed: () => totalCount > 1,
      },
      {
        text: `${intl.formatMessage({ id: "actions.merge" })}…`,
        onClick: () => onMerge(),
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
        onPlay={onPlay}
        operationsMenuClassName="audio-list-operations-dropdown"
      />
    );

    return (
      <div
        className={cx("item-list-container audio-list", {
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
                view={view}
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
                <AudioList
                  filter={effectiveFilter}
                  audios={items}
                  selectedIds={selectedIds}
                  onSelectChange={onSelectChange}
                  fromGroupId={fromGroupId}
                />
              </LoadedContent>

              {totalCount > filter.itemsPerPage && (
                <div className="pagination-footer-container">
                  <div className="pagination-footer">
                    <Pagination
                      itemsPerPage={filter.itemsPerPage}
                      currentPage={filter.currentPage}
                      totalItems={totalCount}
                      metadataByline={metadataByline}
                      onChangePage={setPage}
                      pagePopupPlacement="top"
                    />
                  </div>
                </div>
              )}
            </SidebarPaneContent>
          </SidebarPane>
        </SidebarStateContext.Provider>
      </div>
    );
  }
);

export default FilteredAudioList;
