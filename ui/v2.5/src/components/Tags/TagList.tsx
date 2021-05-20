import React, { useState } from "react";
import _ from "lodash";
import Mousetrap from "mousetrap";
import { FindTagsQueryResult } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import {
  showWhenSelected,
  useTagsList,
  PersistanceLevel,
} from "src/hooks/ListHook";
import { Button } from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindTags,
  mutateMetadataAutoTag,
  useTagDestroy,
  useTagsDestroy,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { FormattedNumber } from "react-intl";
import { NavUtils } from "src/utils";
import { Icon, Modal, DeleteEntityDialog } from "src/components/Shared";
import { TagCard } from "./TagCard";
import { ExportDialog } from "../Shared/ExportDialog";

interface ITagList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const TagList: React.FC<ITagList> = ({ filterHook }) => {
  const Toast = useToast();
  const [
    deletingTag,
    setDeletingTag,
  ] = useState<Partial<GQL.TagDataFragment> | null>(null);

  const [deleteTag] = useTagDestroy(getDeleteTagInput() as GQL.TagDestroyInput);

  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const otherOperations = [
    {
      text: "View Random",
      onClick: viewRandom,
    },
    {
      text: "Export...",
      onClick: onExport,
      isDisplayed: showWhenSelected,
    },
    {
      text: "Export all...",
      onClick: onExportAll,
    },
  ];

  const addKeybinds = (
    result: FindTagsQueryResult,
    filter: ListFilterModel
  ) => {
    Mousetrap.bind("p r", () => {
      viewRandom(result, filter);
    });

    return () => {
      Mousetrap.unbind("p r");
    };
  };

  async function viewRandom(
    result: FindTagsQueryResult,
    filter: ListFilterModel
  ) {
    // query for a random tag
    if (result.data && result.data.findTags) {
      const { count } = result.data.findTags;

      const index = Math.floor(Math.random() * count);
      const filterCopy = _.cloneDeep(filter);
      filterCopy.itemsPerPage = 1;
      filterCopy.currentPage = index + 1;
      const singleResult = await queryFindTags(filterCopy);
      if (
        singleResult &&
        singleResult.data &&
        singleResult.data.findTags &&
        singleResult.data.findTags.tags.length === 1
      ) {
        const { id } = singleResult!.data!.findTags!.tags[0];
        // navigate to the tag page
        history.push(`/tags/${id}`);
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

  function maybeRenderExportDialog(selectedIds: Set<string>) {
    if (isExportDialogOpen) {
      return (
        <>
          <ExportDialog
            exportInput={{
              tags: {
                ids: Array.from(selectedIds.values()),
                all: isExportAll,
              },
            }}
            onClose={() => {
              setIsExportDialogOpen(false);
            }}
          />
        </>
      );
    }
  }

  const renderDeleteDialog = (
    selectedTags: GQL.TagDataFragment[],
    onClose: (confirmed: boolean) => void
  ) => (
    <DeleteEntityDialog
      selected={selectedTags}
      onClose={onClose}
      singularEntity="tag"
      pluralEntity="tags"
      destroyMutation={useTagsDestroy}
    />
  );

  const listData = useTagsList({
    renderContent,
    filterHook,
    addKeybinds,
    otherOperations,
    selectable: true,
    zoomable: true,
    defaultZoomIndex: 0,
    persistState: PersistanceLevel.ALL,
    renderDeleteDialog,
  });

  function getDeleteTagInput() {
    const tagInput: Partial<GQL.TagDestroyInput> = {};
    if (deletingTag) {
      tagInput.id = deletingTag.id;
    }
    return tagInput;
  }

  async function onAutoTag(tag: GQL.TagDataFragment) {
    if (!tag) return;
    try {
      await mutateMetadataAutoTag({ tags: [tag.id] });
      Toast.success({ content: "Started auto tagging" });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onDelete() {
    try {
      await deleteTag();
      Toast.success({ content: "Deleted tag" });
      setDeletingTag(null);
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderTags(
    result: FindTagsQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    if (!result.data?.findTags) return;

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row px-xl-5 justify-content-center">
          {result.data.findTags.tags.map((tag) => (
            <TagCard
              key={tag.id}
              tag={tag}
              zoomIndex={zoomIndex}
              selecting={selectedIds.size > 0}
              selected={selectedIds.has(tag.id)}
              onSelectedChanged={(selected: boolean, shiftKey: boolean) =>
                listData.onSelectChange(tag.id, selected, shiftKey)
              }
            />
          ))}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      const deleteAlert = (
        <Modal
          onHide={() => {}}
          show={!!deletingTag}
          icon="trash-alt"
          accept={{ onClick: onDelete, variant: "danger", text: "Delete" }}
          cancel={{ onClick: () => setDeletingTag(null) }}
        >
          <span>
            Are you sure you want to delete {deletingTag && deletingTag.name}?
          </span>
        </Modal>
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
                Auto Tag
              </Button>
              <Button variant="secondary" className="tag-list-button">
                <Link
                  to={NavUtils.makeTagScenesUrl(tag)}
                  className="tag-list-anchor"
                >
                  Scenes: <FormattedNumber value={tag.scene_count ?? 0} />
                </Link>
              </Button>
              <Button variant="secondary" className="tag-list-button">
                <Link
                  to={NavUtils.makeTagSceneMarkersUrl(tag)}
                  className="tag-list-anchor"
                >
                  Markers:{" "}
                  <FormattedNumber value={tag.scene_marker_count ?? 0} />
                </Link>
              </Button>
              <span className="tag-list-count">
                Total:{" "}
                <FormattedNumber
                  value={(tag.scene_count || 0) + (tag.scene_marker_count || 0)}
                />
              </span>
              <Button variant="danger" onClick={() => setDeletingTag(tag)}>
                <Icon icon="trash-alt" color="danger" />
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

  function renderContent(
    result: FindTagsQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>,
    zoomIndex: number
  ) {
    return (
      <>
        {maybeRenderExportDialog(selectedIds)}
        {renderTags(result, filter, selectedIds, zoomIndex)}
      </>
    );
  }

  return listData.template;
};
