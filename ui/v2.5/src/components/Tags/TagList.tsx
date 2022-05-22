import React, { useState } from "react";
import cloneDeep from "lodash-es/cloneDeep";
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
import { FormattedMessage, FormattedNumber, useIntl } from "react-intl";
import { NavUtils } from "src/utils";
import { Icon, Modal, DeleteEntityDialog } from "src/components/Shared";
import { TagCard } from "./TagCard";
import { ExportDialog } from "../Shared/ExportDialog";
import { tagRelationHook } from "../../core/tags";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";

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

  const intl = useIntl();
  const history = useHistory();
  const [isExportDialogOpen, setIsExportDialogOpen] = useState(false);
  const [isExportAll, setIsExportAll] = useState(false);

  const otherOperations = [
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
      const filterCopy = cloneDeep(filter);
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
      Toast.success({
        content: intl.formatMessage({ id: "toast.started_auto_tagging" }),
      });
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
      tagRelationHook(deletingTag as GQL.TagDataFragment, oldRelations, {
        parents: [],
        children: [],
      });
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.delete_past_tense" },
          {
            count: 1,
            singularEntity: intl.formatMessage({ id: "tag" }),
            pluralEntity: intl.formatMessage({ id: "tags" }),
          }
        ),
      });
      setDeletingTag(null);
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderTags(
    result: FindTagsQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    if (!result.data?.findTags) return;

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <div className="row px-xl-5 justify-content-center">
          {result.data.findTags.tags.map((tag) => (
            <TagCard
              key={tag.id}
              tag={tag}
              zoomIndex={filter.zoomIndex}
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

  function renderContent(
    result: FindTagsQueryResult,
    filter: ListFilterModel,
    selectedIds: Set<string>
  ) {
    return (
      <>
        {maybeRenderExportDialog(selectedIds)}
        {renderTags(result, filter, selectedIds)}
      </>
    );
  }

  return listData.template;
};
