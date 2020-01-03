import React, { useState } from "react";
import { Button, Form, Modal, Spinner } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { ErrorUtils } from "../../utils/errors";
import { NavigationUtils } from "../../utils/navigation";
import { ToastUtils } from "../../utils/toasts";

export const TagList: React.FC = () => {
  // Editing / New state
  const [name, setName] = useState('');
  const [editingTag, setEditingTag] = useState<Partial<GQL.TagDataFragment> | null>(null);
  const [deletingTag, setDeletingTag] = useState<Partial<GQL.TagDataFragment> | null>(null);

  const { data, error } = StashService.useAllTags();
  const updateTag = StashService.useTagUpdate(getTagInput() as GQL.TagUpdateInput);
  const createTag = StashService.useTagCreate(getTagInput() as GQL.TagCreateInput);
  const deleteTag = StashService.useTagDestroy(getDeleteTagInput() as GQL.TagDestroyInput);

  function getTagInput() {
    const tagInput: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = { name };
    if (!!editingTag) { (tagInput as Partial<GQL.TagUpdateInput>).id = editingTag.id; }
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
        ToastUtils.success("Updated tag");
      } else {
        await createTag();
        ToastUtils.success("Created tag");
      }
      setEditingTag(null);
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  async function onAutoTag(tag : GQL.TagDataFragment) {
    if (!tag) {
      return;
    }
    try {
      await StashService.queryMetadataAutoTag({ tags: [tag.id]});
      ToastUtils.success("Started auto tagging");
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  async function onDelete() {
    try {
      await deleteTag();
      ToastUtils.success("Deleted tag");
      setDeletingTag(null);
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  const deleteAlert = (
    <Modal
        onHide={() => {}}
        show={!!deletingTag}
    >
      <Modal.Body>
        <FontAwesomeIcon icon="trash-alt" color="danger" />
        <span>Are you sure you want to delete {deletingTag && deletingTag.name}?</span>
      </Modal.Body>
      <Modal.Footer>
        <div>
          <Button variant="danger" onClick={onDelete}>Delete</Button>
          <Button onClick={() => setDeletingTag(null)}>Cancel</Button>
        </div>
      </Modal.Footer>
    </Modal>
  );

  if (!data || !data.allTags) { return <Spinner animation="border" variant="light" />; }
  if (!!error) { return <>{error.message}</>; }

  const tagElements = data.allTags.map((tag) => {
    return (
      <>
      {deleteAlert}
      <div key={tag.id} className="tag-list-row">
        <span onClick={() => setEditingTag(tag)}>{tag.name}</span>
        <div style={{float: "right"}}>
          <Button onClick={() => onAutoTag(tag)}>Auto Tag</Button>
          <Link to={NavigationUtils.makeTagScenesUrl(tag)}>Scenes: {tag.scene_count}</Link>
          <Link to={NavigationUtils.makeTagSceneMarkersUrl(tag)}>
            Markers: {tag.scene_marker_count}
          </Link>
          <span>Total: {(tag.scene_count || 0) + (tag.scene_marker_count || 0)}</span>
          <Button variant="danger" onClick={() => setDeletingTag(tag)}>
            <FontAwesomeIcon icon="trash-alt" color="danger" />
          </Button>
        </div>
      </div>
      </>
    );
  });

  return (
    <div id="tag-list-container">
      <Button variant="primary" style={{marginTop: "20px"}} onClick={() => setEditingTag({})}>New Tag</Button>

      <Modal
          onHide={() => {setEditingTag(null)}}
          show={!!editingTag}
      >
        <Modal.Header>
          { editingTag && editingTag.id ? "Edit Tag" : "New Tag" }
        </Modal.Header>
        <Modal.Body>
          <Form.Group controlId="tag-name">
            <Form.Label>Name</Form.Label>
            <Form.Control
              onChange={(newValue: any) => setName(newValue.target.value)}
              defaultValue={(editingTag && editingTag.name) || ''}
            />
          </Form.Group>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={() => onEdit()}>{editingTag && editingTag.id ? "Update" : "Create"}</Button>
        </Modal.Footer>
      </Modal>
      {tagElements}
    </div>
  );
};
