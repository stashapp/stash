import { Table, Tabs, Tab } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import cx from "classnames";
import Mousetrap from "mousetrap";

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
import { TagImagesPanel } from "./TagImagesPanel";
import { TagPerformersPanel } from "./TagPerformersPanel";

interface ITabParams {
  id?: string;
  tab?: string;
}

export const Tag: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const { tab = "scenes", id = "new" } = useParams<ITabParams>();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing tag state
  const [image, setImage] = useState<string | null>();
  const [name, setName] = useState<string>();

  // Tag state
  const [tag, setTag] = useState<GQL.TagDataFragment | undefined>();
  const [imagePreview, setImagePreview] = useState<string>();

  const { data, error, loading } = useFindTag(id);
  const [updateTag] = useTagUpdate();
  const [createTag] = useTagCreate(getTagInput() as GQL.TagUpdateInput);
  const [deleteTag] = useTagDestroy(getTagInput() as GQL.TagUpdateInput);

  const activeTabKey =
    tab === "markers" || tab === "images" || tab === "performers"
      ? tab
      : "scenes";
  const setActiveTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "scenes" ? "" : `/${newTab}`;
      history.replace(`/tags/${id}${tabParam}`);
    }
  };

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

  function updateTagEditState(state: GQL.TagDataFragment) {
    setName(state.name);
  }

  function updateTagData(tagData: GQL.TagDataFragment) {
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
    if (!isNew) {
      return {
        id,
        name,
        image,
      };
    }
    return {
      name,
      image,
    };
  }

  async function onSave() {
    try {
      if (!isNew) {
        const result = await updateTag({
          variables: {
            input: getTagInput() as GQL.TagUpdateInput,
          },
        });
        if (result.data?.tagUpdate) {
          if (result.data.tagUpdate.image_path)
            await fetch(result.data.tagUpdate.image_path, { cache: "reload" });
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
    if (!tag?.id) return;
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
    if (tag) {
      updateTagData(tag);
    }
  }

  function onClearImage() {
    setImage(null);
    setImagePreview(
      tag?.image_path ? `${tag.image_path}?default=true` : undefined
    );
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
          onClearImage={() => {
            onClearImage();
          }}
          onAutoTag={onAutoTag}
          onDelete={onDelete}
          acceptSVG
        />
      </div>
      {!isNew && tag && (
        <div className="col col-md-8">
          <Tabs
            id="tag-tabs"
            mountOnEnter
            activeKey={activeTabKey}
            onSelect={setActiveTabKey}
          >
            <Tab eventKey="scenes" title="Scenes">
              <TagScenesPanel tag={tag} />
            </Tab>
            <Tab eventKey="images" title="Images">
              <TagImagesPanel tag={tag} />
            </Tab>
            <Tab eventKey="markers" title="Markers">
              <TagMarkersPanel tag={tag} />
            </Tab>
            <Tab eventKey="performers" title="Performers">
              <TagPerformersPanel tag={tag} />
            </Tab>
          </Tabs>
        </div>
      )}
      {renderDeleteAlert()}
    </div>
  );
};
