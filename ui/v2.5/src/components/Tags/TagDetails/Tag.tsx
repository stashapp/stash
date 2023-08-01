import { Tabs, Tab, Dropdown, Button } from "react-bootstrap";
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
import { CompressedTagDetailsPanel, TagDetailsPanel } from "./TagDetailsPanel";
import { TagEditPanel } from "./TagEditPanel";
import { TagMergeModal } from "./TagMergeDialog";
import {
  faChevronDown,
  faChevronUp,
  faSignInAlt,
  faSignOutAlt,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { IUIConfig } from "src/core/config";
import ImageUtils from "src/utils/image";
import { isPlatfornUniquelyRenderByApple } from "src/utils/apple";

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

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui as IUIConfig | undefined;
  const abbreviateCounter = uiConfig?.abbreviateCounters ?? false;
  const enableBackgroundImage = uiConfig?.enableTagBackgroundImage ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? false;
  const compactExpandedDetails = uiConfig?.compactExpandedDetails ?? false;

  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);
  const [loadStickyHeader, setLoadStickyHeader] = useState<boolean>(false);

  const appleRendering = isPlatfornUniquelyRenderByApple();

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

  const showAllCounts = (configuration?.ui as IUIConfig)?.showChildTagContent;
  const sceneCount =
    (showAllCounts ? tag.scene_count_all : tag.scene_count) ?? 0;
  const imageCount =
    (showAllCounts ? tag.image_count_all : tag.image_count) ?? 0;
  const galleryCount =
    (showAllCounts ? tag.gallery_count_all : tag.gallery_count) ?? 0;
  const sceneMarkerCount =
    (showAllCounts ? tag.scene_marker_count_all : tag.scene_marker_count) ?? 0;
  const performerCount =
    (showAllCounts ? tag.performer_count_all : tag.performer_count) ?? 0;

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
    Mousetrap.bind("e", () => toggleEditing());
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

  useEffect(() => {
    const f = () => {
      if (document.documentElement.scrollTop <= 50) {
        setLoadStickyHeader(false);
      } else {
        setLoadStickyHeader(true);
      }
    };

    window.addEventListener("scroll", f);
    return () => {
      window.removeEventListener("scroll", f);
    };
  });

  async function onSave(input: GQL.TagCreateInput) {
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
      toggleEditing(false);
      const updated = result.data.tagUpdate;
      tagRelationHook(updated, oldRelations, {
        parents: updated.parents,
        children: updated.children,
      });
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "tag" }).toLocaleLowerCase() }
        ),
      });
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

  function getCollapseButtonIcon() {
    return collapsed ? faChevronDown : faChevronUp;
  }

  function maybeRenderShowCollapseButton() {
    if (!isEditing) {
      return (
        <span className="detail-expand-collapse">
          <Button
            className="minimal expand-collapse"
            onClick={() => setCollapsed(!collapsed)}
          >
            <Icon className="fa-fw" icon={getCollapseButtonIcon()} />
          </Button>
        </span>
      );
    }
  }

  function maybeRenderAliases() {
    if (tag?.aliases?.length) {
      return (
        <div>
          <span className="alias-head">{tag?.aliases?.join(", ")}</span>
        </div>
      );
    }
  }

  function toggleEditing(value?: boolean) {
    if (value !== undefined) {
      setIsEditing(value);
    } else {
      setIsEditing((e) => !e);
    }
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
      return (
        <img
          className="logo"
          alt={tag.name}
          src={tagImage}
          onLoad={ImageUtils.verifyImageSize}
        />
      );
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

  function maybeRenderDetails() {
    if (!isEditing) {
      return (
        <TagDetailsPanel
          tag={tag}
          fullWidth={!collapsed && !compactExpandedDetails}
        />
      );
    }
  }

  function maybeRenderEditPanel() {
    if (isEditing) {
      return (
        <TagEditPanel
          tag={tag}
          onSubmit={onSave}
          onCancel={() => toggleEditing()}
          onDelete={onDelete}
          setImage={setImage}
          setEncodingImage={setEncodingImage}
        />
      );
    }
    {
      return (
        <DetailsEditNavbar
          objectName={tag.name}
          isNew={false}
          isEditing={isEditing}
          onToggleEdit={() => toggleEditing()}
          onSave={() => {}}
          onImageChange={() => {}}
          onClearImage={() => {}}
          onAutoTag={onAutoTag}
          onDelete={onDelete}
          classNames="mb-2"
          customButtons={renderMergeButton()}
        />
      );
    }
  }

  const renderTabs = () => (
    <React.Fragment>
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
            <>
              {intl.formatMessage({ id: "scenes" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={sceneCount}
                hideZero
              />
            </>
          }
        >
          <TagScenesPanel active={activeTabKey == "scenes"} tag={tag} />
        </Tab>
        <Tab
          eventKey="images"
          title={
            <>
              {intl.formatMessage({ id: "images" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={imageCount}
                hideZero
              />
            </>
          }
        >
          <TagImagesPanel active={activeTabKey == "images"} tag={tag} />
        </Tab>
        <Tab
          eventKey="galleries"
          title={
            <>
              {intl.formatMessage({ id: "galleries" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={galleryCount}
                hideZero
              />
            </>
          }
        >
          <TagGalleriesPanel active={activeTabKey == "galleries"} tag={tag} />
        </Tab>
        <Tab
          eventKey="markers"
          title={
            <>
              {intl.formatMessage({ id: "markers" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={sceneMarkerCount}
                hideZero
              />
            </>
          }
        >
          <TagMarkersPanel active={activeTabKey == "markers"} tag={tag} />
        </Tab>
        <Tab
          eventKey="performers"
          title={
            <>
              {intl.formatMessage({ id: "performers" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={performerCount}
                hideZero
              />
            </>
          }
        >
          <TagPerformersPanel active={activeTabKey == "performers"} tag={tag} />
        </Tab>
      </Tabs>
    </React.Fragment>
  );

  function maybeRenderHeaderBackgroundImage() {
    let tagImage = tag.image_path;
    if (enableBackgroundImage && !isEditing && tagImage) {
      return (
        <div className="background-image-container">
          <picture>
            <source src={tagImage} />
            <img
              className="background-image"
              src={tagImage}
              alt={`${tag.name} background`}
            />
          </picture>
        </div>
      );
    }
  }

  function maybeRenderTab() {
    if (!isEditing) {
      return renderTabs();
    }
  }

  function maybeRenderCompressedDetails() {
    if (!isEditing && loadStickyHeader) {
      return <CompressedTagDetailsPanel tag={tag} />;
    }
  }

  return (
    <div id="tag-page" className="row">
      <Helmet>
        <title>{tag.name}</title>
      </Helmet>

      <div
        className={`detail-header ${isEditing ? "edit" : ""}  ${
          collapsed ? "collapsed" : !compactExpandedDetails ? "full-width" : ""
        } ${appleRendering ? "apple" : ""}`}
      >
        {maybeRenderHeaderBackgroundImage()}
        <div className="detail-container">
          <div className="detail-header-image">
            {encodingImage ? (
              <LoadingIndicator
                message={`${intl.formatMessage({ id: "encoding_image" })}...`}
              />
            ) : (
              renderImage()
            )}
          </div>
          <div className="row">
            <div className="studio-head col">
              <h2>
                <span className="tag-name">{tag.name}</span>
                {maybeRenderShowCollapseButton()}
              </h2>
              {maybeRenderAliases()}
              {maybeRenderDetails()}
              {maybeRenderEditPanel()}
            </div>
          </div>
        </div>
      </div>
      {maybeRenderCompressedDetails()}
      <div className="detail-body">
        <div className="tag-body">
          <div className="tag-tabs">{maybeRenderTab()}</div>
        </div>
      </div>
      {renderDeleteAlert()}
      {renderMergeDialog()}
    </div>
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
