import { Tabs, Tab, Dropdown } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
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
import { ImageUtils } from "src/utils";
import {
  DetailsEditNavbar,
  Modal,
  LoadingIndicator,
  Icon,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { TagScenesPanel } from "./TagScenesPanel";
import { TagMarkersPanel } from "./TagMarkersPanel";
import { TagImagesPanel } from "./TagImagesPanel";
import { TagPerformersPanel } from "./TagPerformersPanel";
import { TagGalleriesPanel } from "./TagGalleriesPanel";
import { TagDetailsPanel } from "./TagDetailsPanel";
import { TagEditPanel } from "./TagEditPanel";
import { TagMergeModal } from "./TagMergeDialog";

interface ITabParams {
  id?: string;
  tab?: string;
}

export const Tag: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();
  const { tab = "scenes", id = "new" } = useParams<ITabParams>();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [mergeType, setMergeType] = useState<"from" | "into" | undefined>();

  // Editing tag state
  const [image, setImage] = useState<string | null>();

  // Tag state
  const { data, error, loading } = useFindTag(id);
  const tag = data?.findTag;

  const [updateTag] = useTagUpdate();
  const [createTag] = useTagCreate();
  const [deleteTag] = useTagDestroy({ id });

  const activeTabKey =
    tab === "markers" ||
    tab === "images" ||
    tab === "performers" ||
    tab === "galleries"
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

  useEffect(() => {
    if (data && data.findTag) {
      setImage(undefined);
    }
  }, [data]);

  function onImageLoad(imageData: string) {
    setImage(imageData);
  }

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, isEditing);

  if (!isNew && !isEditing) {
    if (!data?.findTag || loading) return <LoadingIndicator />;
    if (error) return <div>{error.message}</div>;
  }

  function getTagInput(
    input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput>
  ) {
    const ret: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = {
      ...input,
      image,
    };

    if (!isNew) {
      (ret as GQL.TagUpdateInput).id = id;
    }

    return ret;
  }

  async function onSave(
    input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput>
  ) {
    try {
      if (!isNew) {
        const result = await updateTag({
          variables: {
            input: getTagInput(input) as GQL.TagUpdateInput,
          },
        });
        if (result.data?.tagUpdate) {
          setIsEditing(false);
          return result.data.tagUpdate.id;
        }
      } else {
        const result = await createTag({
          variables: {
            input: getTagInput(input) as GQL.TagCreateInput,
          },
        });
        if (result.data?.tagCreate?.id) {
          setIsEditing(false);
          return result.data.tagCreate.id;
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

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{
          text: intl.formatMessage({ id: "actions.delete" }),
          variant: "danger",
          onClick: onDelete,
        }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>
          <FormattedMessage
            id="dialogs.delete_confirm"
            values={{
              entityName:
                tag?.name ??
                intl.formatMessage({ id: "tag" }).toLocaleLowerCase(),
            }}
          />
        </p>
      </Modal>
    );
  }

  function onToggleEdit() {
    setIsEditing(!isEditing);
    setImage(undefined);
  }

  function renderImage() {
    let tagImage = tag?.image_path;
    if (isEditing) {
      if (image === null) {
        tagImage = `${tagImage}&default=true`;
      } else if (image) {
        tagImage = image;
      }
    }

    if (tagImage) {
      return <img className="logo" alt={tag?.name ?? ""} src={tagImage} />;
    }
  }

  function renderMergeButton() {
    return (
      <Dropdown drop="up">
        <Dropdown.Toggle variant="secondary">
          <FormattedMessage id="actions.merge" />
          ...
        </Dropdown.Toggle>
        <Dropdown.Menu className="bg-secondary text-white" id="tag-merge-menu">
          <Dropdown.Item
            className="bg-secondary text-white"
            onClick={() => setMergeType("from")}
          >
            <Icon icon="sign-in-alt" />
            <FormattedMessage id="actions.merge_from" />
            ...
          </Dropdown.Item>
          <Dropdown.Item
            className="bg-secondary text-white"
            onClick={() => setMergeType("into")}
          >
            <Icon icon="sign-out-alt" />
            <FormattedMessage id="actions.merge_into" />
            ...
          </Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    );
  }

  function renderMergeDialog() {
    if (!tag || !mergeType) return;
    return (
      <TagMergeModal
        tag={tag}
        onClose={() => setMergeType(undefined)}
        show={!!mergeType}
        mergeType={mergeType}
      />
    );
  }

  return (
    <div className="row">
      <div
        className={cx("tag-details", {
          "col-md-4": !isNew,
          "col-md-8": isNew,
        })}
      >
        <div className="text-center logo-container">
          {imageEncoding ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            renderImage()
          )}
          {!isNew && tag && <h2>{tag.name}</h2>}
        </div>
        {!isEditing && !isNew && tag ? (
          <>
            <TagDetailsPanel tag={tag} />
            {/* HACK - this is also rendered in the TagEditPanel */}
            <DetailsEditNavbar
              objectName={tag.name ?? "tag"}
              isNew={isNew}
              isEditing={isEditing}
              onToggleEdit={onToggleEdit}
              onSave={() => {}}
              onImageChange={() => {}}
              onClearImage={() => {}}
              onAutoTag={onAutoTag}
              onDelete={onDelete}
              customButtons={renderMergeButton()}
            />
          </>
        ) : (
          <TagEditPanel
            tag={tag ?? undefined}
            onSubmit={onSave}
            onCancel={onToggleEdit}
            onDelete={onDelete}
            setImage={setImage}
          />
        )}
      </div>
      {!isNew && tag && (
        <div className="col col-md-8">
          <Tabs
            id="tag-tabs"
            mountOnEnter
            activeKey={activeTabKey}
            onSelect={setActiveTabKey}
          >
            <Tab eventKey="scenes" title={intl.formatMessage({ id: "scenes" })}>
              <TagScenesPanel tag={tag} />
            </Tab>
            <Tab eventKey="images" title={intl.formatMessage({ id: "images" })}>
              <TagImagesPanel tag={tag} />
            </Tab>
            <Tab
              eventKey="galleries"
              title={intl.formatMessage({ id: "galleries" })}
            >
              <TagGalleriesPanel tag={tag} />
            </Tab>
            <Tab
              eventKey="markers"
              title={intl.formatMessage({ id: "markers" })}
            >
              <TagMarkersPanel tag={tag} />
            </Tab>
            <Tab
              eventKey="performers"
              title={intl.formatMessage({ id: "performers" })}
            >
              <TagPerformersPanel tag={tag} />
            </Tab>
          </Tabs>
        </div>
      )}
      {renderDeleteAlert()}
      {renderMergeDialog()}
    </div>
  );
};
