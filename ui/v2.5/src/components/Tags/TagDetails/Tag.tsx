import { Tabs, Tab, Dropdown, Badge } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import Mousetrap from "mousetrap";

import * as GQL from "src/core/generated-graphql";
import {
  useFindTag,
  useTagUpdate,
  useTagDestroy,
  mutateMetadataAutoTag,
} from "src/core/StashService";
import { ImageUtils } from "src/utils";
import {
  DetailsEditNavbar,
  ErrorMessage,
  Modal,
  LoadingIndicator,
  Icon,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { tagRelationHook } from "src/core/tags";
import { TagScenesPanel } from "./TagScenesPanel";
import { TagMarkersPanel } from "./TagMarkersPanel";
import { TagImagesPanel } from "./TagImagesPanel";
import { TagPerformersPanel } from "./TagPerformersPanel";
import { TagGalleriesPanel } from "./TagGalleriesPanel";
import { TagDetailsPanel } from "./TagDetailsPanel";
import { TagEditPanel } from "./TagEditPanel";
import { TagMergeModal } from "./TagMergeDialog";
import renderNonZero from "src/utils/render";

interface IProps {
  tag: GQL.TagDataFragment;
}

interface ITabParams {
  tab?: string;
}

const TagPage: React.FC<IProps> = ({ tag }) => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();
  const { tab = "scenes" } = useParams<ITabParams>();

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [mergeType, setMergeType] = useState<"from" | "into" | undefined>();

  // Editing tag state
  const [image, setImage] = useState<string | null>();

  const [updateTag] = useTagUpdate();
  const [deleteTag] = useTagDestroy({ id: tag.id });

  const defaultTab =
    tag?.scene_count ?? 0 > 0
      ? "scenes"
      : tag?.image_count ?? 0 > 0
      ? "images"
      : tag?.gallery_count ?? 0 > 0
      ? "galleries"
      : tag?.scene_marker_count ?? 0 > 0
      ? "markers"
      : "performers";

  const activeTabKey =
    tab === "markers" ||
    tab === "images" ||
    tab === "performers" ||
    tab === "galleries"
      ? tab
      : defaultTab;
  const setActiveTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "scenes" ? "" : `/${newTab}`;
      history.replace(`/tags/${tag.id}${tabParam}`);
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

  function onImageLoad(imageData: string) {
    setImage(imageData);
  }

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, isEditing);

  function getTagInput(
    input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput>
  ) {
    const ret: Partial<GQL.TagCreateInput | GQL.TagUpdateInput> = {
      ...input,
      image,
      id: tag.id,
    };

    return ret;
  }

  async function onSave(
    input: Partial<GQL.TagCreateInput | GQL.TagUpdateInput>
  ) {
    try {
      const oldRelations = {
        parents: tag.parents ?? [],
        children: tag.children ?? [],
      };
      const result = await updateTag({
        variables: {
          input: getTagInput(input) as GQL.TagUpdateInput,
        },
      });
      if (result.data?.tagUpdate) {
        setIsEditing(false);
        const updated = result.data.tagUpdate;
        tagRelationHook(updated, oldRelations, {
          parents: updated.parents,
          children: updated.children,
        });
        return updated.id;
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onAutoTag() {
    if (!tag.id) return;
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
        parents: tag.parents ?? [],
        children: tag.children ?? [],
      };
      await deleteTag();
      tagRelationHook(tag as GQL.TagDataFragment, oldRelations, {
        parents: [],
        children: [],
      });
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
                tag.name ??
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
    let tagImage = tag.image_path;
    if (isEditing) {
      if (image === null) {
        tagImage = `${tagImage}&default=true`;
      } else if (image) {
        tagImage = image;
      }
    }

    if (tagImage) {
      return <img className="logo" alt={tag.name} src={tagImage} />;
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
    <>
      <Helmet>
        <title>{tag.name}</title>
      </Helmet>
      <div className="row">
        <div className="tag-details col-md-4">
          <div className="text-center logo-container">
            {imageEncoding ? (
              <LoadingIndicator message="Encoding image..." />
            ) : (
              renderImage()
            )}
            <h2>{tag.name}</h2>
          </div>
          {!isEditing ? (
            <>
              <TagDetailsPanel tag={tag} />
              {/* HACK - this is also rendered in the TagEditPanel */}
              <DetailsEditNavbar
                objectName={tag.name}
                isNew={false}
                isEditing={isEditing}
                onToggleEdit={onToggleEdit}
                onSave={() => {}}
                onImageChange={() => {}}
                onClearImage={() => {}}
                onAutoTag={onAutoTag}
                onDelete={onDelete}
                classNames="mb-2"
                customButtons={renderMergeButton()}
              />
            </>
          ) : (
            <TagEditPanel
              tag={tag}
              onSubmit={onSave}
              onCancel={onToggleEdit}
              onDelete={onDelete}
              setImage={setImage}
            />
          )}
        </div>
        <div className="col col-md-8">
          <Tabs
            id="tag-tabs"
            mountOnEnter
            activeKey={activeTabKey}
            onSelect={setActiveTabKey}
          >
            {renderNonZero(
              tag.scene_count,
              <Tab
                eventKey="scenes"
                title={
                  <React.Fragment>
                    {intl.formatMessage({ id: "scenes" })}
                    <Badge className="left-spacing" pill variant="secondary">
                      {intl.formatNumber(tag.scene_count ?? 0)}
                    </Badge>
                  </React.Fragment>
                }
              >
                <TagScenesPanel tag={tag} />
              </Tab>
            )}
            {renderNonZero(
              tag.image_count,
              <Tab
                eventKey="images"
                title={
                  <React.Fragment>
                    {intl.formatMessage({ id: "images" })}
                    <Badge className="left-spacing" pill variant="secondary">
                      {intl.formatNumber(tag.image_count ?? 0)}
                    </Badge>
                  </React.Fragment>
                }
              >
                <TagImagesPanel tag={tag} />
              </Tab>
            )}
            {renderNonZero(
              tag.gallery_count,
              <Tab
                eventKey="galleries"
                title={
                  <React.Fragment>
                    {intl.formatMessage({ id: "galleries" })}
                    <Badge className="left-spacing" pill variant="secondary">
                      {intl.formatNumber(tag.gallery_count ?? 0)}
                    </Badge>
                  </React.Fragment>
                }
              >
                <TagGalleriesPanel tag={tag} />
              </Tab>
            )}
            {renderNonZero(
              tag.scene_marker_count,
              <Tab
                eventKey="markers"
                title={
                  <React.Fragment>
                    {intl.formatMessage({ id: "markers" })}
                    <Badge className="left-spacing" pill variant="secondary">
                      {intl.formatNumber(tag.scene_marker_count ?? 0)}
                    </Badge>
                  </React.Fragment>
                }
              >
                <TagMarkersPanel tag={tag} />
              </Tab>
            )}
            {renderNonZero(
              tag.performer_count,
              <Tab
                eventKey="performers"
                title={
                  <React.Fragment>
                    {intl.formatMessage({ id: "performers" })}
                    <Badge className="left-spacing" pill variant="secondary">
                      {intl.formatNumber(tag.performer_count ?? 0)}
                    </Badge>
                  </React.Fragment>
                }
              >
                <TagPerformersPanel tag={tag} />
              </Tab>
            )}
          </Tabs>
        </div>
        {renderDeleteAlert()}
        {renderMergeDialog()}
      </div>
    </>
  );
};

const TagLoader: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const { data, loading, error } = useFindTag(id ?? "");

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findTag)
    return <ErrorMessage error={`No tag found with id ${id}.`} />;

  return <TagPage tag={data.findTag} />;
};

export default TagLoader;
