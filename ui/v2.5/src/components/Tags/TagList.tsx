import React, { useState } from "react";
import { Button, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import { FormattedNumber } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  mutateMetadataAutoTag,
  useAllTags,
  useTagUpdate,
  useTagCreate,
  useTagDestroy,
} from "src/core/StashService";
import { NavUtils } from "src/utils";
import { Icon, Modal, LoadingIndicator } from "src/components/Shared";
import { useToast } from "src/hooks";

export const TagList: React.FC = () => {
  const Toast = useToast();
  // Editing / New state
  const [name, setName] = useState("");
  const [editingTag, setEditingTag] = useState<Partial<
    GQL.TagDataFragment
  > | null>(null);
  const [deletingTag, setDeletingTag] = useState<Partial<
    GQL.TagDataFragment
  > | null>(null);

  const { data, error } = useAllTags();
  const [updateTag] = useTagUpdate(getTagInput() as GQL.TagUpdateInput);
  const [createTag] = useTagCreate(getTagInput() as GQL.TagCreateInput);
  const [deleteTag] = useTagDestroy(getDeleteTagInput() as GQL.TagDestroyInput);

  function getTagInput() {
    const tagInput: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = { name };
    if (editingTag)
      (tagInput as Partial<GQL.TagUpdateInput>).id = editingTag.id;
    return tagInput;
  }

  function getDeleteTagInput() {
    const tagInput: Partial<GQL.TagDestroyInput> = {};
    if (deletingTag) {
      tagInput.id = deletingTag.id;
    }
    return tagInput;
  }

  async function onEdit() {
    try {
      if (editingTag && editingTag.id) {
        await updateTag();
        Toast.success({ content: "Updated tag" });
      } else {
        await createTag();
        Toast.success({ content: "Created tag" });
      }
      setEditingTag(null);
    } catch (e) {
      Toast.error(e);
    }
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

  if (!data?.allTags) return <LoadingIndicator />;
  if (error) return <div>{error.message}</div>;

  const tagElements = data.allTags.map((tag) => {
    return (
      <div key={tag.id} className="tag-list-row row">
        <Button variant="link" onClick={() => setEditingTag(tag)}>
          {tag.name}
        </Button>
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
              Markers: <FormattedNumber value={tag.scene_marker_count ?? 0} />
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
      <Button
        variant="primary"
        className="mt-2"
        onClick={() => setEditingTag({})}
      >
        New Tag
      </Button>

      <Modal
        show={!!editingTag}
        header={editingTag && editingTag.id ? "Edit Tag" : "New Tag"}
        onHide={() => setEditingTag(null)}
        accept={{
          onClick: onEdit,
          variant: "danger",
          text: editingTag?.id ? "Update" : "Create",
        }}
      >
        <Form.Group controlId="tag-name">
          <Form.Label>Name</Form.Label>
          <Form.Control
            onChange={(newValue: React.ChangeEvent<HTMLInputElement>) =>
              setName(newValue.currentTarget.value)
            }
            defaultValue={(editingTag && editingTag.name) || ""}
          />
        </Form.Group>
      </Modal>

      {tagElements}
      {deleteAlert}
    </div>
  );
};
