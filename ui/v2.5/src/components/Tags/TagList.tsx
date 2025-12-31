import React, { useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
import Mousetrap from "mousetrap";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { ItemList, ItemListContext, showWhenSelected } from "../List/ItemList";
import { Button } from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindTagsForList,
  mutateMetadataAutoTag,
  useFindTagsForList,
  useTagDestroy,
  useTagsDestroy,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import NavUtils from "src/utils/navigation";
import { Icon } from "../Shared/Icon";
import { ModalComponent } from "../Shared/Modal";
import { DeleteEntityDialog } from "../Shared/DeleteEntityDialog";
import { ExportDialog } from "../Shared/ExportDialog";
import { tagRelationHook } from "../../core/tags";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { TagMergeModal } from "./TagMergeDialog";
import { Tag } from "./TagSelect";
import { TagCardGrid } from "./TagCardGrid";
import { EditTagsDialog } from "./EditTagsDialog";
import { View } from "../List/views";
import { IItemListOperation } from "../List/FilteredListToolbar";
import { PatchComponent } from "src/patch";

function getItems(result: GQL.FindTagsForListQueryResult) {
  return result?.data?.findTags?.tags ?? [];
}

function getCount(result: GQL.FindTagsForListQueryResult) {
  return result?.data?.findTags?.count ?? 0;
}

interface ITagList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
  alterQuery?: boolean;
  extraOperations?: IItemListOperation<GQL.FindTagsForListQueryResult>[];
}

export const TagList: React.FC<ITagList> = PatchComponent(
  "TagList",
  ({ filterHook, alterQuery, extraOperations = [] }) => {
    const Toast = useToast();
    const [deletingTag, setDeletingTag] =
      useState<Partial<GQL.TagListDataFragment> | null>(null);

    const filterMode = GQL.FilterMode.Tags;
    const view = View.Tags;

    function getDeleteTagInput() {
      const tagInput: Partial<GQL.TagDestroyInput> = {};
      if (deletingTag) {
        tagInput.id = deletingTag.id;
      }
      return tagInput as GQL.TagDestroyInput;
    }
    const [deleteTag] = useTagDestroy(getDeleteTagInput());

    const intl = useIntl();
    const history = useHistory();
    const [mergeTags, setMergeTags] = useState<Tag[] | undefined>(undefined);
    const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
    const [isExportAll, setIsExportAll] = useState(false);

    const otherOperations = [
      ...extraOperations,
      {
        text: intl.formatMessage({ id: "actions.view_random" }),
        onClick: viewRandom,
      },
      {
        text: `${intl.formatMessage({ id: "actions.merge" })}…`,
        onClick: merge,
        isDisplayed: showWhenSelected,
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
      result: GQL.FindTagsForListQueryResult,
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
      result: GQL.FindTagsForListQueryResult,
      filter: ListFilterModel
    ) {
      // query for a random tag
      if (result.data?.findTags) {
        const { count } = result.data.findTags;

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
      }
    }

    async function merge(
      result: GQL.FindTagsForListQueryResult,
      filter: ListFilterModel,
      selectedIds: Set<string>
    ) {
      const selected =
        result.data?.findTags.tags.filter((t) => selectedIds.has(t.id)) ?? [];
      setMergeTags(selected);
    }

    async function onExport() {
      setIsExportAll(false);
      setIsExportDialogOpen(true);
    }

    async function onExportAll() {
      setIsExportAll(true);
      setIsExportDialogOpen(true);
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

    async function onDelete() {
      try {
        const oldRelations = {
          parents: deletingTag?.parents ?? [],
          children: deletingTag?.children ?? [],
        };
        await deleteTag();
        tagRelationHook(deletingTag as GQL.TagListDataFragment, oldRelations, {
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
      } catch (e) {
        Toast.error(e);
      }
    }

    function renderContent(
      result: GQL.FindTagsForListQueryResult,
      filter: ListFilterModel,
      selectedIds: Set<string>,
      onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void
    ) {
      function renderMergeDialog() {
        if (mergeTags) {
          return (
            <TagMergeModal
              tags={mergeTags}
              onClose={(mergedId?: string) => {
                setMergeTags(undefined);
                if (mergedId) {
                  history.push(`/tags/${mergedId}`);
                }
              }}
              show
            />
          );
        }
      }

      function maybeRenderExportDialog() {
        if (isExportDialogOpen) {
          return (
            <ExportDialog
              exportInput={{
                tags: {
                  ids: Array.from(selectedIds.values()),
                  all: isExportAll,
                },
              }}
              onClose={() => setIsExportDialogOpen(false)}
            />
          );
        }
      }

      function renderTags() {
        if (!result.data?.findTags) return;

        if (filter.displayMode === DisplayMode.Grid) {
          return (
            <TagCardGrid
              tags={result.data.findTags.tags}
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
                onClick: onDelete,
                variant: "danger",
                text: intl.formatMessage({ id: "actions.delete" }),
              }}
              cancel={{ onClick: () => setDeletingTag(null) }}
            >
              <span>
                <FormattedMessage
                  id="dialogs.delete_confirm"
                  values={{ entityName: deletingTag && deletingTag.name }}
                />
              </span>
            </ModalComponent>
          );

          const tagElements = result.data.findTags.tags.map((tag) => {
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
                  <Button variant="danger" onClick={() => setDeletingTag(tag)}>
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
        if (filter.displayMode === DisplayMode.Wall) {
          return <h1>TODO</h1>;
        }
      }
      return (
        <>
          {renderMergeDialog()}
          {maybeRenderExportDialog()}
          {renderTags()}
        </>
      );
    }

    function renderEditDialog(
      selectedTags: GQL.TagListDataFragment[],
      onClose: (confirmed: boolean) => void
    ) {
      return <EditTagsDialog selected={selectedTags} onClose={onClose} />;
    }

    function renderDeleteDialog(
      selectedTags: GQL.TagListDataFragment[],
      onClose: (confirmed: boolean) => void
    ) {
      return (
        <DeleteEntityDialog
          selected={selectedTags}
          onClose={onClose}
          singularEntity={intl.formatMessage({ id: "tag" })}
          pluralEntity={intl.formatMessage({ id: "tags" })}
          destroyMutation={useTagsDestroy}
          onDeleted={() => {
            selectedTags.forEach((t) =>
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

    return (
      <ItemListContext
        filterMode={filterMode}
        useResult={useFindTagsForList}
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
        />
      </ItemListContext>
    );
  }
);
