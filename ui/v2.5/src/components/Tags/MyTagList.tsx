import React, { useCallback, useEffect, useMemo, useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import { Link, useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindTags,
  mutateMetadataAutoTag,
  useFindTags,
  useTagDestroy,
  useTagsDestroy,
} from "src/core/StashService";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { View } from "../List/views";
import { LoadedContent } from "../List/PagedList";
import { useCloseEditDelete, useFilterOperations } from "../List/util";
import {
  OperationDropdown,
  OperationDropdownItem,
} from "../List/ListOperationButtons";
import {
  FilteredListToolbar2,
  ToolbarFilterSection,
  ToolbarSelectionSection,
} from "../List/MyListToolbar";
import { useFilteredItemList } from "../List/ItemList";
import {
  Sidebar,
  SidebarPane,
  SidebarPaneContent,
  SidebarStateContext,
  useSidebarState,
} from "../Shared/Sidebar";
import cx from "classnames";
import {
  FilteredSidebarHeader,
  useFilteredSidebarKeybinds,
} from "../List/Filters/MyFilterSidebar";
import { Pagination } from "../List/Pagination";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import {
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
  createBooleanCriterionOption,
  MandatoryNumberCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import useFocus from "src/utils/myFocus";
import {
  faFilter,
  faPencil,
  faPlus,
  faTrash,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { SidebarNumberFilter } from "../List/Filters/NumberFilter";
import { PatchContainerComponent } from "src/patch";
import { ListResultsHeader } from "../List/MyListResultsHeader";
import { SidebarBooleanFilter } from "../List/Filters/BooleanFilter";
import { FavoriteTagCriterionOption } from "src/models/list-filter/criteria/favorite";
import { SidebarStringFilter } from "../List/Filters/StringFilter";
import { TagCardGrid } from "./TagCardGrid";
import { EditTagsDialog } from "./EditTagsDialog";
import { tagRelationHook } from "../../core/tags";
import { useToast } from "src/hooks/Toast";
import NavUtils from "src/utils/navigation";
import { ModalComponent } from "../Shared/Modal";
import {
  ParentTagsCriterionOption,
  ChildTagsCriterionOption,
} from "src/models/list-filter/criteria/tags";
import { SidebarTagsFilter } from "../List/Filters/TagsFilter";
import {
  SidebarFilterSelector,
  FilterWrapper,
} from "../List/Filters/SidebarFilterSelector";
import { SidebarFilterDefinition } from "src/hooks/useSidebarFilters";
import { createMandatoryTimestampCriterionOption } from "src/models/list-filter/criteria/criterion";
import { SidebarDateFilter } from "../List/Filters/DateFilter";
import { TagIsMissingCriterionOption } from "src/models/list-filter/criteria/is-missing";
import { SidebarIsMissingFilter } from "../List/Filters/IsMissingFilter";
import {
  useTagFacetCounts,
  FacetCountsContext,
} from "src/hooks/useFacetCounts";

function useViewRandom(
  result: GQL.FindTagsQueryResult,
  filter: ListFilterModel
) {
  const history = useHistory();

  const viewRandom = useCallback(async () => {
    // query for a random tag
    if (result.data?.findTags) {
      const { count } = result.data.findTags;

      const index = Math.floor(Math.random() * count);
      const filterCopy = cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindTags(filterCopy);
      if (singleResult.data.findTags.tags.length === 1) {
        const { id } = singleResult.data.findTags.tags[0];
        // navigate to the tag page
        history.push(`/tags/${id}`);
      }
    }
  }, [result, filter, history]);

  return viewRandom;
}

function useAddKeybinds(
  result: GQL.FindTagsQueryResult,
  filter: ListFilterModel
) {
  const viewRandom = useViewRandom(result, filter);

  useEffect(() => {
    Mousetrap.bind("p r", () => {
      viewRandom();
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  }, [viewRandom]);
}

const TagListContent: React.FC<{
  tags: GQL.TagDataFragment[];
  filter: ListFilterModel;
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
  onAutoTag: (tag: GQL.TagDataFragment) => void;
  onDeleteTag: (tag: GQL.TagDataFragment) => void;
  deletingTag: Partial<GQL.TagDataFragment> | null;
  onConfirmDelete: () => void;
  onCancelDelete: () => void;
}> = ({
  tags,
  filter,
  selectedIds,
  onSelectChange,
  onAutoTag,
  onDeleteTag,
  deletingTag,
  onConfirmDelete,
  onCancelDelete,
}) => {
  const intl = useIntl();

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
    const deleteAlert = (
      <ModalComponent
        onHide={() => {}}
        show={!!deletingTag}
        icon={faTrashAlt}
        accept={{
          onClick: onConfirmDelete,
          variant: "danger",
          text: intl.formatMessage({ id: "actions.delete" }),
        }}
        cancel={{ onClick: onCancelDelete }}
      >
        <span>
          <FormattedMessage
            id="dialogs.delete_confirm"
            values={{ entityName: deletingTag && deletingTag.name }}
          />
        </span>
      </ModalComponent>
    );

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
            <Button variant="danger" onClick={() => onDeleteTag(tag)}>
              <Icon icon={faTrashAlt} color="danger" />
            </Button>
          </div>
        </div>
      );
    });

    return (
      <div className="col col-sm-8 m-auto">
        {tagElements}
        {deleteAlert}
      </div>
    );
  }

  return null;
};

export const MyTagsFilterSidebarSections = PatchContainerComponent(
  "MyFilteredTagList.SidebarSections"
);

// Define available filters for tags sidebar
const tagFilterDefinitions: SidebarFilterDefinition[] = [
  // Tier 1: Primary filters (visible by default)
  { id: "favorite", messageId: "favourite", defaultVisible: true },
  { id: "parent_tags", messageId: "parent_tags", defaultVisible: true },
  { id: "child_tags", messageId: "sub_tags", defaultVisible: true },

  // Tier 2: Library stats
  { id: "scene_count", messageId: "scene_count", defaultVisible: false },
  { id: "image_count", messageId: "image_count", defaultVisible: false },
  { id: "gallery_count", messageId: "gallery_count", defaultVisible: false },
  { id: "performer_count", messageId: "performer_count", defaultVisible: false },
  { id: "studio_count", messageId: "studio_count", defaultVisible: false },
  { id: "group_count", messageId: "group_count", defaultVisible: false },
  { id: "marker_count", messageId: "marker_count", defaultVisible: false },

  // Tier 2b: Count stats
  { id: "parent_tag_count", messageId: "parent_count", defaultVisible: false },
  { id: "sub_tag_count", messageId: "child_count", defaultVisible: false },

  // Tier 3: Metadata
  { id: "name", messageId: "name", defaultVisible: false },
  { id: "sort_name", messageId: "sort_name", defaultVisible: false },
  { id: "aliases", messageId: "aliases", defaultVisible: false },
  { id: "description", messageId: "description", defaultVisible: false },

  // Tier 4: System
  { id: "is_missing", messageId: "isMissing", defaultVisible: false },
  { id: "ignore_auto_tag", messageId: "ignore_auto_tag", defaultVisible: false },
  { id: "created_at", messageId: "created_at", defaultVisible: false },
  { id: "updated_at", messageId: "updated_at", defaultVisible: false },
];

const SidebarContent: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  view?: View;
  showEditFilter: (editingCriterion?: string) => void;
  sidebarOpen: boolean;
  onClose?: () => void;
  count?: number;
  focus?: ReturnType<typeof useFocus>;
  clearAllCriteria: () => void;
  onFilterEditModeChange?: (isEditMode: boolean) => void;
}> = ({
  filter,
  setFilter,
  view,
  showEditFilter,
  sidebarOpen,
  onClose,
  count,
  focus,
  clearAllCriteria,
  onFilterEditModeChange,
}) => {
  const showResultsId =
    count !== undefined ? "actions.show_count_results" : "actions.show_results";

  // Criterion options
  const NameCriterionOption = createStringCriterionOption("name");
  const SortNameCriterionOption = createStringCriterionOption("sort_name");
  const AliasesCriterionOption = createStringCriterionOption("aliases");
  const DescriptionCriterionOption = createStringCriterionOption("description");
  const SceneCountCriterionOption = createMandatoryNumberCriterionOption("scene_count");
  const ImageCountCriterionOption = createMandatoryNumberCriterionOption("image_count");
  const GalleryCountCriterionOption = createMandatoryNumberCriterionOption("gallery_count");
  const PerformerCountCriterionOption = createMandatoryNumberCriterionOption("performer_count");
  const StudioCountCriterionOption = createMandatoryNumberCriterionOption("studio_count");
  const GroupCountCriterionOption = createMandatoryNumberCriterionOption("group_count");
  const MarkerCountCriterionOption = createMandatoryNumberCriterionOption("marker_count");
  const ParentTagCountCriterionOption = new MandatoryNumberCriterionOption("parent_tag_count", "parent_count");
  const SubTagCountCriterionOption = new MandatoryNumberCriterionOption("sub_tag_count", "child_count");
  const IgnoreAutoTagCriterionOption = createBooleanCriterionOption("ignore_auto_tag");
  const CreatedAtCriterionOption = createMandatoryTimestampCriterionOption("created_at");
  const UpdatedAtCriterionOption = createMandatoryTimestampCriterionOption("updated_at");

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

      <MyTagsFilterSidebarSections>
        <div className="sidebar-filters">
          <SidebarFilterSelector
            viewName="tags"
            filterDefinitions={tagFilterDefinitions}
            headerContent={
              <>
                <Icon icon={faFilter} />
                <FormattedMessage id="filters" />
              </>
            }
            onEditModeChange={onFilterEditModeChange}
          >
            {/* Tier 1: Primary Filters */}
            <FilterWrapper filterId="favorite">
              <SidebarBooleanFilter
                title={<FormattedMessage id="favourite" />}
                option={FavoriteTagCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="favourite"
              />
            </FilterWrapper>
            <FilterWrapper filterId="parent_tags">
              <SidebarTagsFilter
                title={<FormattedMessage id="parent_tags" />}
                option={ParentTagsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="parent_tags"
              />
            </FilterWrapper>
            <FilterWrapper filterId="child_tags">
              <SidebarTagsFilter
                title={<FormattedMessage id="sub_tags" />}
                option={ChildTagsCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="child_tags"
              />
            </FilterWrapper>

            {/* Tier 2: Library Stats */}
            <FilterWrapper filterId="scene_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="scene_count" />}
                option={SceneCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="scene_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="image_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="image_count" />}
                option={ImageCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="image_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="gallery_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="gallery_count" />}
                option={GalleryCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="gallery_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="performer_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="performer_count" />}
                option={PerformerCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="performer_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="studio_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="studio_count" />}
                option={StudioCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="studio_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="group_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="group_count" />}
                option={GroupCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="group_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="marker_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="marker_count" />}
                option={MarkerCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="marker_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="parent_tag_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="parent_count" />}
                option={ParentTagCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="parent_tag_count"
              />
            </FilterWrapper>
            <FilterWrapper filterId="sub_tag_count">
              <SidebarNumberFilter
                title={<FormattedMessage id="child_count" />}
                option={SubTagCountCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="sub_tag_count"
              />
            </FilterWrapper>

            {/* Tier 3: Metadata */}
            <FilterWrapper filterId="name">
              <SidebarStringFilter
                title={<FormattedMessage id="name" />}
                option={NameCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="name"
              />
            </FilterWrapper>
            <FilterWrapper filterId="sort_name">
              <SidebarStringFilter
                title={<FormattedMessage id="sort_name" />}
                option={SortNameCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="sort_name"
              />
            </FilterWrapper>
            <FilterWrapper filterId="aliases">
              <SidebarStringFilter
                title={<FormattedMessage id="aliases" />}
                option={AliasesCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="aliases"
              />
            </FilterWrapper>
            <FilterWrapper filterId="description">
              <SidebarStringFilter
                title={<FormattedMessage id="description" />}
                option={DescriptionCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="description"
              />
            </FilterWrapper>

            {/* Tier 4: System */}
            <FilterWrapper filterId="is_missing">
              <SidebarIsMissingFilter
                title={<FormattedMessage id="isMissing" />}
                option={TagIsMissingCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="is_missing"
              />
            </FilterWrapper>
            <FilterWrapper filterId="ignore_auto_tag">
              <SidebarBooleanFilter
                title={<FormattedMessage id="ignore_auto_tag" />}
                option={IgnoreAutoTagCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="ignore_auto_tag"
              />
            </FilterWrapper>
            <FilterWrapper filterId="created_at">
              <SidebarDateFilter
                title={<FormattedMessage id="created_at" />}
                option={CreatedAtCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="created_at"
                isTime
              />
            </FilterWrapper>
            <FilterWrapper filterId="updated_at">
              <SidebarDateFilter
                title={<FormattedMessage id="updated_at" />}
                option={UpdatedAtCriterionOption}
                filter={filter}
                setFilter={setFilter}
                sectionID="updated_at"
                isTime
              />
            </FilterWrapper>
          </SidebarFilterSelector>
        </div>
      </MyTagsFilterSidebarSections>
      <div className="sidebar-footer">
        <Button className="sidebar-close-button" onClick={onClose}>
          <FormattedMessage id={showResultsId} values={{ count }} />
        </Button>
      </div>
      <div className="clear-all-filters">
        <Button
          className="clear-all-filters-button"
          variant="secondary"
          onClick={() => clearAllCriteria()}
          title="Clear All Filters"
        >
          <FormattedMessage id="Clear All Filters" />
        </Button>
      </div>
    </>
  );
};

interface IOperations {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
  className?: string;
}

const TagListOperations: React.FC<{
  items: number;
  hasSelection: boolean;
  operations: IOperations[];
  onEdit: () => void;
  onDelete: () => void;
  onCreateNew: () => void;
}> = ({ items, hasSelection, operations, onEdit, onDelete, onCreateNew }) => {
  const intl = useIntl();

  return (
    <div className="list-operations">
      <ButtonGroup>
        {!hasSelection && (
          <Button
            className="create-new-button"
            variant="secondary"
            onClick={() => onCreateNew()}
            title={intl.formatMessage(
              { id: "actions.create_entity" },
              { entityType: intl.formatMessage({ id: "tag" }) }
            )}
          >
            <Icon icon={faPlus} />
          </Button>
        )}

        {hasSelection && (
          <>
            <Button variant="secondary" onClick={() => onEdit()}>
              <Icon icon={faPencil} />
            </Button>
            <Button
              variant="danger"
              className="btn-danger-minimal"
              onClick={() => onDelete()}
            >
              <Icon icon={faTrash} />
            </Button>
          </>
        )}

        <OperationDropdown
          className="list-operations"
          menuPortalTarget={document.body}
        >
          {operations.map((o) => {
            if (o.isDisplayed && !o.isDisplayed()) {
              return null;
            }

            return (
              <OperationDropdownItem
                key={o.text}
                onClick={o.onClick}
                text={o.text}
                className={o.className}
              />
            );
          })}
        </OperationDropdown>
      </ButtonGroup>
    </div>
  );
};

interface IFilteredTags {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  defaultSort?: string;
  view?: View;
  alterQuery?: boolean;
}

export const MyFilteredTagList: React.FC<IFilteredTags> = (props) => {
  const intl = useIntl();
  const history = useHistory();
  const Toast = useToast();

  const searchFocus = useFocus();
  const [, setSearchFocus] = searchFocus;

  const { filterHook, defaultSort, view, alterQuery } = props;

  // State for individual tag deletion in list view
  const [deletingTag, setDeletingTag] =
    useState<Partial<GQL.TagDataFragment> | null>(null);

  function getDeleteTagInput() {
    const tagInput: Partial<GQL.TagDestroyInput> = {};
    if (deletingTag) {
      tagInput.id = deletingTag.id;
    }
    return tagInput as GQL.TagDestroyInput;
  }
  const [deleteTag] = useTagDestroy(getDeleteTagInput());

  // States
  const {
    showSidebar,
    setShowSidebar,
    loading: sidebarStateLoading,
    sectionOpen: baseSectionOpen,
    setSectionOpen: baseSetSectionOpen,
  } = useSidebarState(view);

  // Track filter customization edit mode
  const [isFilterEditMode, setIsFilterEditMode] = useState(false);

  const sectionOpen = useMemo(() => {
    if (isFilterEditMode) {
      const closedSections: Record<string, boolean> = {};
      Object.keys(baseSectionOpen).forEach((key) => {
        closedSections[key] = false;
      });
      return closedSections;
    }
    return baseSectionOpen;
  }, [isFilterEditMode, baseSectionOpen]);

  const setSectionOpen = useCallback(
    (section: string, open: boolean) => {
      if (isFilterEditMode) return;
      baseSetSectionOpen(section, open);
    },
    [isFilterEditMode, baseSetSectionOpen]
  );

  const { filterState, queryResult, modalState, listSelect, showEditFilter } =
    useFilteredItemList({
      filterStateProps: {
        filterMode: GQL.FilterMode.Tags,
        defaultSort,
        view,
        useURL: alterQuery,
      },
      queryResultProps: {
        useResult: useFindTags,
        getCount: (r) => r.data?.findTags.count ?? 0,
        getItems: (r) => r.data?.findTags.tags ?? [],
        filterHook,
      },
    });

  const { filter, setFilter, loading: filterLoading } = filterState;

  const { effectiveFilter, result, cachedResult, items, totalCount } =
    queryResult;

  const {
    selectedIds,
    selectedItems,
    onSelectChange,
    onSelectAll,
    onSelectNone,
    hasSelection,
  } = listSelect;

  const { modal, showModal, closeModal } = modalState;

  // Utility hooks
  const { setPage, removeCriterion, clearAllCriteria } = useFilterOperations({
    filter,
    setFilter,
  });

  useAddKeybinds(result, filter);
  useFilteredSidebarKeybinds({
    showSidebar,
    setShowSidebar,
  });

  // Fetch facet counts for sidebar filters
  const { counts: facetCounts, loading: facetLoading } = useTagFacetCounts(filter, {
    isOpen: showSidebar ?? false,
    debounceMs: 300,
  });

  useEffect(() => {
    Mousetrap.bind("e", () => {
      if (hasSelection) {
        onEdit();
      }
    });

    Mousetrap.bind("d d", () => {
      if (hasSelection) {
        onDelete();
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

  const viewRandom = useViewRandom(result, filter);

  function onCreateNew() {
    history.push("/tags/new");
  }

  function onExport(all: boolean) {
    showModal(
      <ExportDialog
        exportInput={{
          tags: {
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

  function onDelete() {
    showModal(
      <DeleteEntityDialog
        selected={selectedItems}
        onClose={onCloseEditDelete}
        singularEntity={intl.formatMessage({ id: "tag" })}
        pluralEntity={intl.formatMessage({ id: "tags" })}
        destroyMutation={useTagsDestroy}
        onDeleted={() => {
          selectedItems.forEach((t) =>
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

  async function onAutoTag(tag: GQL.TagDataFragment) {
    if (!tag) return;
    try {
      await mutateMetadataAutoTag({ tags: [tag.id] });
      Toast.success(intl.formatMessage({ id: "toast.started_auto_tagging" }));
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onConfirmDeleteSingleTag() {
    try {
      const oldRelations = {
        parents: deletingTag?.parents ?? [],
        children: deletingTag?.children ?? [],
      };
      await deleteTag();
      tagRelationHook(deletingTag as GQL.TagDataFragment, oldRelations, {
        parents: [],
        children: [],
      });
      Toast.success(
        intl.formatMessage(
          { id: "toast.delete_past_tense" },
          {
            count: 1,
            singularEntity: intl.formatMessage({ id: "tag" }),
            pluralEntity: intl.formatMessage({ id: "tags" }),
          }
        )
      );
      setDeletingTag(null);
      result.refetch();
    } catch (e) {
      Toast.error(e);
    }
  }

  const otherOperations = [
    {
      text: intl.formatMessage({ id: "actions.view_random" }),
      onClick: viewRandom,
      isDisplayed: () => totalCount > 1,
    },
    {
      text: intl.formatMessage(
        { id: "actions.create_entity" },
        { entityType: intl.formatMessage({ id: "tag" }) }
      ),
      onClick: () => onCreateNew(),
      isDisplayed: () => !hasSelection,
      className: "create-new-item",
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
  if (filterLoading || sidebarStateLoading) return null;

  const operations = (
    <TagListOperations
      items={items.length}
      hasSelection={hasSelection}
      operations={otherOperations}
      onEdit={onEdit}
      onDelete={onDelete}
      onCreateNew={onCreateNew}
    />
  );

  return (
    <div
      className={cx("item-list-container tag-list", {
        "hide-sidebar": !showSidebar,
      })}
    >
      {modal}

      <SidebarStateContext.Provider value={{ sectionOpen, setSectionOpen, disabled: isFilterEditMode }}>
        <FacetCountsContext.Provider value={{ counts: facetCounts, loading: facetLoading }}>
        <SidebarPane hideSidebar={!showSidebar}>
          <Sidebar hide={!showSidebar} onHide={() => setShowSidebar(false)}>
            <SidebarContent
              filter={filter}
              setFilter={setFilter}
              showEditFilter={showEditFilter}
              view={view}
              sidebarOpen={showSidebar}
              onFilterEditModeChange={setIsFilterEditMode}
              onClose={() => setShowSidebar(false)}
              count={cachedResult.loading ? undefined : totalCount}
              focus={searchFocus}
              clearAllCriteria={() => clearAllCriteria(true)}
            />
          </Sidebar>
          <SidebarPaneContent>
            <FilteredListToolbar2
              className="tag-list-toolbar"
              hasSelection={hasSelection}
              filterSection={
                <ToolbarFilterSection
                  filter={filter}
                  onSetFilter={setFilter}
                  onToggleSidebar={() => setShowSidebar(!showSidebar)}
                  onEditCriterion={(c) =>
                    showEditFilter(c?.criterionOption.type)
                  }
                  onRemoveCriterion={removeCriterion}
                  onRemoveAllCriterion={() => clearAllCriteria(true)}
                  onEditSearchTerm={() => {
                    setShowSidebar(true);
                    setSearchFocus(true);
                  }}
                  onRemoveSearchTerm={() => setFilter(filter.clearSearchTerm())}
                />
              }
              selectionSection={
                <ToolbarSelectionSection
                  selected={selectedIds.size}
                  onToggleSidebar={() => setShowSidebar(!showSidebar)}
                  onSelectAll={() => onSelectAll()}
                  onSelectNone={() => onSelectNone()}
                  operations={operations}
                />
              }
              operationSection={operations}
            />

            <ListResultsHeader
              loading={cachedResult.loading}
              filter={filter}
              totalCount={totalCount}
              onChangeFilter={(newFilter) => setFilter(newFilter)}
            />

            <LoadedContent loading={result.loading} error={result.error}>
              <TagListContent
                filter={effectiveFilter}
                tags={items}
                selectedIds={selectedIds}
                onSelectChange={onSelectChange}
                onAutoTag={onAutoTag}
                onDeleteTag={setDeletingTag}
                deletingTag={deletingTag}
                onConfirmDelete={onConfirmDeleteSingleTag}
                onCancelDelete={() => setDeletingTag(null)}
              />
            </LoadedContent>

            {totalCount > filter.itemsPerPage && (
              <div className="pagination-footer">
                <Pagination
                  itemsPerPage={filter.itemsPerPage}
                  currentPage={filter.currentPage}
                  totalItems={totalCount}
                  onChangePage={setPage}
                  pagePopupPlacement="top"
                />
              </div>
            )}
          </SidebarPaneContent>
        </SidebarPane>
        </FacetCountsContext.Provider>
      </SidebarStateContext.Provider>
    </div>
  );
};

// Backward compatibility wrapper
export const TagList: React.FC<IFilteredTags> = (props) => {
  return <MyFilteredTagList {...props} />;
};

