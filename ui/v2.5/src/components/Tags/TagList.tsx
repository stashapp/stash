import React, { useCallback, useEffect } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import Mousetrap from "mousetrap";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { useFilteredItemList } from "../List/ItemList";
import { Button } from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindTagsForList,
  mutateMetadataAutoTag,
  useFindTagsForList,
  useTagsDestroy,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import NavUtils from "src/utils/navigation";
import { Icon } from "../Shared/Icon";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { tagRelationHook } from "../../core/tags";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { TagMergeModal } from "./TagMergeDialog";
import { TagCardGrid } from "./TagCardGrid";
import { EditTagsDialog } from "./EditTagsDialog";
import { View } from "../List/views";
import {
  FilteredListToolbar,
  IItemListOperation,
} from "../List/FilteredListToolbar";
import { PatchComponent, PatchContainerComponent } from "src/patch";
import { TagTagger } from "../Tagger/tags/TagTagger";
import useFocus from "src/utils/focus";
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
import { ListOperations } from "../List/ListOperationButtons";
import cx from "classnames";
import { FilterTags } from "../List/FilterTags";
import { Pagination, PaginationIndex } from "../List/Pagination";
import { LoadedContent } from "../List/PagedList";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { FavoriteTagCriterionOption } from "src/models/list-filter/criteria/favorite";

const TagList: React.FC<{
  tags: GQL.TagListDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  onDelete: (tag: GQL.TagListDataFragment) => void;
  onAutoTag: (tag: GQL.TagListDataFragment) => void;
}> = PatchComponent(
  "TagList",
  ({ tags, filter, selectedIds, onSelectChange, onDelete, onAutoTag }) => {
    if (tags.length === 0 && filter.displayMode !== DisplayMode.Tagger) {
      return null;
    }

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <TagCardGrid
          tags={tags}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={onSelectChange}
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      const tagElements = tags.map((tag) => {
        return (
          <div key={tag.id} className="tag-list-row row">
            <Link to={`/tags/${tag.id}`}>{tag.name}</Link>

            <div className="ml-auto">
              <Button
                variant="secondary"
                className="tag-list-button"
                onClick={() => onAutoTag(tag)}
              >
                <FormattedMessage id="actions.auto_tag" />
              </Button>
              <Button variant="secondary" className="tag-list-button">
                <Link
                  to={NavUtils.makeTagScenesUrl(tag)}
                  className="tag-list-anchor"
                >
                  <FormattedMessage
                    id="countables.scenes"
                    values={{
                      count: tag.scene_count ?? 0,
                    }}
                  />
                  : <FormattedNumber value={tag.scene_count ?? 0} />
                </Link>
              </Button>
              <Button variant="secondary" className="tag-list-button">
                <Link
                  to={NavUtils.makeTagImagesUrl(tag)}
                  className="tag-list-anchor"
                >
                  <FormattedMessage
                    id="countables.images"
                    values={{
                      count: tag.image_count ?? 0,
                    }}
                  />
                  : <FormattedNumber value={tag.image_count ?? 0} />
                </Link>
              </Button>
              <Button variant="secondary" className="tag-list-button">
                <Link
                  to={NavUtils.makeTagGalleriesUrl(tag)}
                  className="tag-list-anchor"
                >
                  <FormattedMessage
                    id="countables.galleries"
                    values={{
                      count: tag.gallery_count ?? 0,
                    }}
                  />
                  : <FormattedNumber value={tag.gallery_count ?? 0} />
                </Link>
              </Button>
              <Button variant="secondary" className="tag-list-button">
                <Link
                  to={NavUtils.makeTagSceneMarkersUrl(tag)}
                  className="tag-list-anchor"
                >
                  <FormattedMessage
                    id="countables.markers"
                    values={{
                      count: tag.scene_marker_count ?? 0,
                    }}
                  />
                  : <FormattedNumber value={tag.scene_marker_count ?? 0} />
                </Link>
              </Button>
              <span className="tag-list-count">
                <FormattedMessage id="total" />:{" "}
                <FormattedNumber
                  value={
                    (tag.scene_count || 0) +
                    (tag.scene_marker_count || 0) +
                    (tag.image_count || 0) +
                    (tag.gallery_count || 0)
                  }
                />
              </span>
              <Button variant="danger" onClick={() => onDelete(tag)}>
                <Icon icon={faTrashAlt} color="danger" />
              </Button>
            </div>
          </div>
        );
      });

      return <div className="col col-sm-8 m-auto">{tagElements}</div>;
    }
    if (filter.displayMode === DisplayMode.Tagger) {
      return <TagTagger tags={tags} />;
    }

    return null;
  }
);

const TagFilterSidebarSections = PatchContainerComponent(
  "FilteredTagList.SidebarSections"
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
  // filterHook,
  view,
  showEditFilter,
  sidebarOpen,
  onClose,
  count,
  focus,
}) => {
  const showResultsId =
    count !== undefined ? "actions.show_count_results" : "actions.show_results";

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

      <TagFilterSidebarSections>
        {/* <SidebarTagsFilter
          filter={filter}
          setFilter={setFilter}
          filterHook={filterHook}
        /> */}
        <SidebarBooleanFilter
          title={<FormattedMessage id="favourite" />}
          filter={filter}
          setFilter={setFilter}
          option={FavoriteTagCriterionOption}
          sectionID="favourite"
        />
      </TagFilterSidebarSections>

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
    // query for a random tag
    if (count === 0) {
      return;
    }

    const index = Math.floor(Math.random() * count);
    const filterCopy = cloneDeep(filter);
    filterCopy.itemsPerPage = 1;
    filterCopy.currentPage = index + 1;
    const singleResult = await queryFindTagsForList(filterCopy);
    if (singleResult.data.findTags.tags.length === 1) {
      const { id } = singleResult.data.findTags.tags[0];
      // navigate to the tag page
      history.push(`/tags/${id}`);
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

interface ITagList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  alterQuery?: boolean;
  extraOperations?: IItemListOperation<GQL.FindTagsForListQueryResult>[];
}

export const FilteredTagList = PatchComponent(
  "FilteredTagList",
  (props: ITagList) => {
    const intl = useIntl();
    const history = useHistory();
    const Toast = useToast();

    const searchFocus = useFocus();

    const { filterHook, alterQuery, extraOperations = [] } = props;

    const view = View.Tags;

    // States
    const {
      showSidebar,
      setShowSidebar,
      sectionOpen,
      setSectionOpen,
      loading: sidebarStateLoading,
    } = useSidebarState(view);

    const { filterState, queryResult, modalState, listSelect, showEditFilter } =
      useFilteredItemList({
        filterStateProps: {
          filterMode: GQL.FilterMode.Tags,
          view,
          useURL: alterQuery,
        },
        queryResultProps: {
          useResult: useFindTagsForList,
          getCount: (r) => r.data?.findTags.count ?? 0,
          getItems: (r) => r.data?.findTags.tags ?? [],
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

    // Utility hooks
    const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
      filter,
      setFilter,
    });

    useAddKeybinds(effectiveFilter, totalCount);
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
            studios: {
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
        <EditTagsDialog selected={selectedItems} onClose={onCloseEditDelete} />
      );
    }

    function onDelete(tag?: GQL.TagListDataFragment) {
      const itemsToDelete = tag ? [tag] : selectedItems;

      showModal(
        <DeleteEntityDialog
          selected={itemsToDelete}
          onClose={onCloseEditDelete}
          singularEntity={intl.formatMessage({ id: "tag" })}
          pluralEntity={intl.formatMessage({ id: "tags" })}
          destroyMutation={useTagsDestroy}
          onDeleted={() => {
            itemsToDelete.forEach((t) =>
              tagRelationHook(
                t,
                { parents: t.parents ?? [], children: t.children ?? [] },
                { parents: [], children: [] }
              )
            );
          }}
        />
      );
    }

    function onMerge() {
      showModal(
        <TagMergeModal
          tags={selectedItems}
          onClose={(mergedId?: string) => {
            onCloseEditDelete();
            if (mergedId) {
              history.push(`/tags/${mergedId}`);
            }
          }}
          show
        />
      );
    }

    async function onAutoTag(tag: GQL.TagListDataFragment) {
      if (!tag) return;
      try {
        await mutateMetadataAutoTag({ tags: [tag.id] });
        Toast.success(intl.formatMessage({ id: "toast.started_auto_tagging" }));
      } catch (e) {
        Toast.error(e);
      }
    }

    const convertedExtraOperations = extraOperations.map((op) => ({
      text: op.text,
      onClick: () => op.onClick(result, filter, selectedIds),
      isDisplayed: () => op.isDisplayed?.(result, filter, selectedIds) ?? true,
    }));

    const otherOperations = [
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
        operationsMenuClassName="tag-list-operations-dropdown"
      />
    );

    return (
      <div
        className={cx("item-list-container tag-list", {
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
                />
              </div>

              <LoadedContent loading={result.loading} error={result.error}>
                <TagList
                  filter={effectiveFilter}
                  tags={items}
                  selectedIds={selectedIds}
                  onSelectChange={onSelectChange}
                  onDelete={(tag) => onDelete(tag)}
                  onAutoTag={(tag) => onAutoTag(tag)}
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
