import React, { useState } from "react";
import { FindTagsQueryResult } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { DisplayMode } from "src/models/list-filter/types";
import { useTagsList } from "src/hooks/ListHook";
import { Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { mutateMetadataAutoTag, useTagDestroy } from "src/core/StashService";
import { useToast } from "src/hooks";
import { FormattedNumber } from "react-intl";
import { NavUtils } from "src/utils";
import { TagCard } from "./TagCard";
import { Icon, Modal } from "../Shared";

interface ITagList {
  filterHook?: (filter: ListFilterModel) => ListFilterModel;
}

export const TagList: React.FC<ITagList> = ({ filterHook }) => {
  const Toast = useToast();
  const [deletingTag, setDeletingTag] = useState<Partial<
    GQL.TagDataFragment
  > | null>(null);

  const [deleteTag] = useTagDestroy(getDeleteTagInput() as GQL.TagDestroyInput);

  const listData = useTagsList({
    renderContent,
    filterHook,
    zoomable: true,
    defaultZoomIndex: 0,
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

  function renderContent(
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
            <TagCard key={tag.id} tag={tag} zoomIndex={zoomIndex} />
          ))}
        </div>
      );
    }
    if (filter.displayMode === DisplayMode.List) {
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

  return listData.template;
};
