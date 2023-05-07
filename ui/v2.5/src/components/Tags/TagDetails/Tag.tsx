import { Button, Tabs, Tab, Dropdown } from "react-bootstrap";
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
import { Counter } from "src/components/Shared/Counter";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { Icon } from "src/components/Shared/Icon";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import { tagRelationHook } from "src/core/tags";
import { TagScenesPanel } from "./TagScenesPanel";
import { TagMarkersPanel } from "./TagMarkersPanel";
import { TagImagesPanel } from "./TagImagesPanel";
import { TagPerformersPanel } from "./TagPerformersPanel";
import { TagGalleriesPanel } from "./TagGalleriesPanel";
import { TagDetailsPanel } from "./TagDetailsPanel";
import { TagEditPanel } from "./TagEditPanel";
import { TagMergeModal } from "./TagMergeDialog";
import {
  faSignInAlt,
  faSignOutAlt,
  faTrashAlt,
  faChevronRight,
  faChevronLeft,
} from "@fortawesome/free-solid-svg-icons";
import { IUIConfig } from "src/core/config";

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

  const [collapsed, setCollapsed] = useState(false);

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const abbreviateCounter =
    (configuration?.ui as IUIConfig)?.abbreviateCounters ?? false;

  const { tab = "scenes" } = useParams<ITabParams>();

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [mergeType, setMergeType] = useState<"from" | "into" | undefined>();

  // Editing tag state
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [updateTag] = useTagUpdate();
  const [deleteTag] = useTagDestroy({ id: tag.id });

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
      history.replace(`/tags/${tag.id}${tabParam}`);
    }
  };

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => {
      onDelete();
    });
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      if (isEditing) {
        Mousetrap.unbind("s s");
      }

      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
      Mousetrap.unbind(",");
    };
  });

  async function onSave(input: GQL.TagCreateInput) {
    try {
      const oldRelations = {
        parents: tag.parents ?? [],
        children: tag.children ?? [],
      };
      const result = await updateTag({
        variables: {
          input: {
            id: tag.id,
            ...input,
          },
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
      <ModalComponent
        show={isDeleteAlertOpen}
        icon={faTrashAlt}
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
      </ModalComponent>
    );
  }

  function onToggleEdit() {
    setIsEditing(!isEditing);
    setImage(undefined);
  }

  function renderImage() {
    let tagImage = tag.image_path;
    if (isEditing) {
      if (image === null && tagImage) {
        const tagImageURL = new URL(tagImage);
        tagImageURL.searchParams.set("default", "true");
        tagImage = tagImageURL.toString();
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
            <Icon icon={faSignInAlt} />
            <FormattedMessage id="actions.merge_from" />
            ...
          </Dropdown.Item>
          <Dropdown.Item
            className="bg-secondary text-white"
            onClick={() => setMergeType("into")}
          >
            <Icon icon={faSignOutAlt} />
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

  function getCollapseButtonIcon() {
    return collapsed ? faChevronRight : faChevronLeft;
  }

  return (
    <>
      <Helmet>
        <title>{tag.name}</title>
      </Helmet>
      <div className="row">
        <div
          className={`tag-details details-tab ${collapsed ? "collapsed" : ""}`}
        >
          <div className="text-center logo-container">
            {encodingImage ? (
              <LoadingIndicator message="Encoding image..." />
            ) : (
              renderImage()
            )}
            <h2>{tag.name}</h2>
            <p>{tag.description}</p>
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
              setEncodingImage={setEncodingImage}
            />
          )}
        </div>
        <div className="details-divider d-none d-xl-block">
          <Button onClick={() => setCollapsed(!collapsed)}>
            <Icon className="fa-fw" icon={getCollapseButtonIcon()} />
          </Button>
        </div>
        <div className={`col content-container ${collapsed ? "expanded" : ""}`}>
          <Tabs
            id="tag-tabs"
            mountOnEnter
            unmountOnExit
            activeKey={activeTabKey}
            onSelect={setActiveTabKey}
          >
            <Tab
              eventKey="scenes"
              title={
                <React.Fragment>
                  {intl.formatMessage({ id: "scenes" })}
                  <Counter
                    abbreviateCounter={abbreviateCounter}
                    count={tag.scene_count ?? 0}
                  />
                </React.Fragment>
              }
            >
              <TagScenesPanel active={activeTabKey == "scenes"} tag={tag} />
            </Tab>
            <Tab
              eventKey="images"
              title={
                <React.Fragment>
                  {intl.formatMessage({ id: "images" })}
                  <Counter
                    abbreviateCounter={abbreviateCounter}
                    count={tag.image_count ?? 0}
                  />
                </React.Fragment>
              }
            >
              <TagImagesPanel active={activeTabKey == "images"} tag={tag} />
            </Tab>
            <Tab
              eventKey="galleries"
              title={
                <React.Fragment>
                  {intl.formatMessage({ id: "galleries" })}
                  <Counter
                    abbreviateCounter={abbreviateCounter}
                    count={tag.gallery_count ?? 0}
                  />
                </React.Fragment>
              }
            >
              <TagGalleriesPanel
                active={activeTabKey == "galleries"}
                tag={tag}
              />
            </Tab>
            <Tab
              eventKey="markers"
              title={
                <React.Fragment>
                  {intl.formatMessage({ id: "markers" })}
                  <Counter
                    abbreviateCounter={abbreviateCounter}
                    count={tag.scene_marker_count ?? 0}
                  />
                </React.Fragment>
              }
            >
              <TagMarkersPanel active={activeTabKey == "markers"} tag={tag} />
            </Tab>
            <Tab
              eventKey="performers"
              title={
                <React.Fragment>
                  {intl.formatMessage({ id: "performers" })}
                  <Counter
                    abbreviateCounter={abbreviateCounter}
                    count={tag.performer_count ?? 0}
                  />
                </React.Fragment>
              }
            >
              <TagPerformersPanel
                active={activeTabKey == "performers"}
                tag={tag}
              />
            </Tab>
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
