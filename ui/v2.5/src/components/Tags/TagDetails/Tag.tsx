/* eslint-disable react/no-this-in-sfc */

import { Table, Tabs, Tab } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import {
  useFindTag,
  useTagUpdate,
  useTagCreate,
  useTagDestroy,
  mutateMetadataAutoTag,
} from "src/core/StashService";
import { ImageUtils, TableUtils } from "src/utils";
import {
  DetailsEditNavbar,
  Modal,
  LoadingIndicator,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { TagScenesPanel } from "./TagScenesPanel";
import { TagMarkersPanel } from "./TagMarkersPanel";

export const Tag: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const { id = "new" } = useParams();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing tag state
  const [image, setImage] = useState<string>();
  const [name, setName] = useState<string>();

  // Tag state
  const [tag, setTag] = useState<Partial<GQL.TagDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string>();

  const { data, error, loading } = useFindTag(id);
  const [updateTag] = useTagUpdate(getTagInput() as GQL.TagUpdateInput);
  const [createTag] = useTagCreate(getTagInput() as GQL.TagUpdateInput);
  const [deleteTag] = useTagDestroy(getTagInput() as GQL.TagUpdateInput);

  // set up hotkeys
  useEffect(() => {
    if (isEditing) {
      Mousetrap.bind("s s", () => onSave());
    }

    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => onDelete());

    return () => {
      if (isEditing) {
        Mousetrap.unbind("s s");
      }

      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  function updateTagEditState(state: Partial<GQL.TagDataFragment>) {
    setName(state.name);
  }

  function updateTagData(tagData: Partial<GQL.TagDataFragment>) {
    setImage(undefined);
    updateTagEditState(tagData);
    setImagePreview(tagData.image_path ?? undefined);
    setTag(tagData);
  }

  useEffect(() => {
    if (data && data.findTag) {
      setImage(undefined);
      updateTagEditState(data.findTag);
      setImagePreview(data.findTag.image_path ?? undefined);
      setTag(data.findTag);
    }
  }, [data]);

  function onImageLoad(imageData: string) {
    setImagePreview(imageData);
    setImage(imageData);
  }

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, isEditing);

  if (!isNew && !isEditing) {
    if (!data?.findTag || loading) return <LoadingIndicator />;
    if (error) return <div>{error.message}</div>;
  }

  function getTagInput() {
    const input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = {
      name,
      image,
    };

    if (!isNew) {
      (input as GQL.TagUpdateInput).id = id;
    }
    return input;
  }

  async function onSave() {
    try {
      if (!isNew) {
        const result = await updateTag();
        if (result.data?.tagUpdate) {
          updateTagData(result.data.tagUpdate);
          setIsEditing(false);
        }
      } else {
        const result = await createTag();
        if (result.data?.tagCreate?.id) {
          history.push(`/tags/${result.data.tagCreate.id}`);
          setIsEditing(false);
        }
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onAutoTag() {
    if (!tag.id) return;
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
    } catch (e) {
      Toast.error(e);
    }

    // redirect to tags page
    history.push(`/tags`);
  }

  function onImageChangeHandler(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{ text: "Delete", variant: "danger", onClick: onDelete }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>Are you sure you want to delete {name ?? "tag"}?</p>
      </Modal>
    );
  }

  function onToggleEdit() {
    setIsEditing(!isEditing);
    updateTagData(tag);
  }

  return (
    <div className="row">
      <div
        className={cx("tag-details", {
          "col-md-4": !isNew,
          "col-8": isNew,
        })}
      >
        {isNew && <h2>Add Tag</h2>}
        <div className="text-center">
          {imageEncoding ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            <img className="logo" alt={name} src={imagePreview} />
          )}
        </div>
        <Table>
          <tbody>
            {TableUtils.renderInputGroup({
              title: "Name",
              value: name ?? "",
              isEditing: !!isEditing,
              onChange: setName,
            })}
          </tbody>
        </Table>
        <DetailsEditNavbar
          objectName={name ?? "tag"}
          isNew={isNew}
          isEditing={isEditing}
          onToggleEdit={onToggleEdit}
          onSave={onSave}
          onImageChange={onImageChangeHandler}
          onAutoTag={onAutoTag}
          onDelete={onDelete}
          acceptSVG
        />
      </div>
      {!isNew && (
        <div className="col col-md-8">
          <Tabs id="tag-tabs" mountOnEnter>
            <Tab eventKey="tag-scenes-panel" title="Scenes">
              <TagScenesPanel tag={tag} />
            </Tab>
            <Tab eventKey="tag-markers-panel" title="Markers">
              <TagMarkersPanel tag={tag} />
            </Tab>
          </Tabs>
        </div>
      )}
      {renderDeleteAlert()}
    </div>
  );
};
