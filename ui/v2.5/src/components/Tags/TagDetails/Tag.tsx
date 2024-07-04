import { Tabs, Tab, Dropdown } from "react-bootstrap";
import React, { useEffect, useMemo, useState } from "react";
import { useHistory, Redirect, RouteComponentProps } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import cx from "classnames";
import Mousetrap from "mousetrap";

import * as GQL from "src/core/generated-graphql";
import {
  useFindTag,
  useTagUpdate,
  useTagDestroy,
  mutateMetadataAutoTag,
} from "src/core/StashService";
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
import { TagStudiosPanel } from "./TagStudiosPanel";
import { TagGalleriesPanel } from "./TagGalleriesPanel";
import { CompressedTagDetailsPanel, TagDetailsPanel } from "./TagDetailsPanel";
import { TagEditPanel } from "./TagEditPanel";
import { TagMergeModal } from "./TagMergeDialog";
import {
  faSignInAlt,
  faSignOutAlt,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { DetailImage } from "src/components/Shared/DetailImage";
import { useLoadStickyHeader } from "src/hooks/detailsPanel";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { TagGroupsPanel } from "./TagGroupsPanel";
import { BackgroundImage } from "src/components/Shared/DetailsPage/BackgroundImage";
import {
  TabTitleCounter,
  useTabKey,
} from "src/components/Shared/DetailsPage/Tabs";
import { DetailTitle } from "src/components/Shared/DetailsPage/DetailTitle";
import { ExpandCollapseButton } from "src/components/Shared/CollapseButton";
import { FavoriteIcon } from "src/components/Shared/FavoriteIcon";
import { AliasList } from "src/components/Shared/DetailsPage/AliasList";
import { HeaderImage } from "src/components/Shared/DetailsPage/HeaderImage";

interface IProps {
  tag: GQL.TagDataFragment;
  tabKey?: TabKey;
}

interface ITagParams {
  id: string;
  tab?: string;
}

const validTabs = [
  "default",
  "scenes",
  "images",
  "galleries",
  "groups",
  "markers",
  "performers",
  "studios",
] as const;
type TabKey = (typeof validTabs)[number];

function isTabKey(tab: string): tab is TabKey {
  return validTabs.includes(tab as TabKey);
}

const TagTabs: React.FC<{
  tabKey?: TabKey;
  tag: GQL.TagDataFragment;
  abbreviateCounter: boolean;
  showAllCounts?: boolean;
}> = ({ tabKey, tag, abbreviateCounter, showAllCounts = false }) => {
  const sceneCount =
    (showAllCounts ? tag.scene_count_all : tag.scene_count) ?? 0;
  const imageCount =
    (showAllCounts ? tag.image_count_all : tag.image_count) ?? 0;
  const galleryCount =
    (showAllCounts ? tag.gallery_count_all : tag.gallery_count) ?? 0;
  const groupCount =
    (showAllCounts ? tag.group_count_all : tag.group_count) ?? 0;
  const sceneMarkerCount =
    (showAllCounts ? tag.scene_marker_count_all : tag.scene_marker_count) ?? 0;
  const performerCount =
    (showAllCounts ? tag.performer_count_all : tag.performer_count) ?? 0;
  const studioCount =
    (showAllCounts ? tag.studio_count_all : tag.studio_count) ?? 0;

  const populatedDefaultTab = useMemo(() => {
    let ret: TabKey = "scenes";
    if (sceneCount == 0) {
      if (imageCount != 0) {
        ret = "images";
      } else if (galleryCount != 0) {
        ret = "galleries";
      } else if (groupCount != 0) {
        ret = "groups";
      } else if (sceneMarkerCount != 0) {
        ret = "markers";
      } else if (performerCount != 0) {
        ret = "performers";
      } else if (studioCount != 0) {
        ret = "studios";
      }
    }

    return ret;
  }, [
    sceneCount,
    imageCount,
    galleryCount,
    sceneMarkerCount,
    performerCount,
    studioCount,
    groupCount,
  ]);

  const { setTabKey } = useTabKey({
    tabKey,
    validTabs,
    defaultTabKey: populatedDefaultTab,
    baseURL: `/tags/${tag.id}`,
  });

  return (
    <Tabs
      id="tag-tabs"
      mountOnEnter
      unmountOnExit
      activeKey={tabKey}
      onSelect={setTabKey}
    >
      <Tab
        eventKey="scenes"
        title={
          <TabTitleCounter
            messageID="scenes"
            count={sceneCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <TagScenesPanel active={tabKey === "scenes"} tag={tag} />
      </Tab>
      <Tab
        eventKey="images"
        title={
          <TabTitleCounter
            messageID="images"
            count={imageCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <TagImagesPanel active={tabKey === "images"} tag={tag} />
      </Tab>
      <Tab
        eventKey="galleries"
        title={
          <TabTitleCounter
            messageID="galleries"
            count={galleryCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <TagGalleriesPanel active={tabKey === "galleries"} tag={tag} />
      </Tab>
      <Tab
        eventKey="groups"
        title={
          <TabTitleCounter
            messageID="groups"
            count={groupCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <TagGroupsPanel active={tabKey === "groups"} tag={tag} />
      </Tab>
      <Tab
        eventKey="markers"
        title={
          <TabTitleCounter
            messageID="markers"
            count={sceneMarkerCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <TagMarkersPanel active={tabKey === "markers"} tag={tag} />
      </Tab>
      <Tab
        eventKey="performers"
        title={
          <TabTitleCounter
            messageID="performers"
            count={performerCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <TagPerformersPanel active={tabKey === "performers"} tag={tag} />
      </Tab>
      <Tab
        eventKey="studios"
        title={
          <TabTitleCounter
            messageID="studios"
            count={studioCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <TagStudiosPanel active={tabKey === "studios"} tag={tag} />
      </Tab>
    </Tabs>
  );
};

const TagPage: React.FC<IProps> = ({ tag, tabKey }) => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui;
  const abbreviateCounter = uiConfig?.abbreviateCounters ?? false;
  const enableBackgroundImage = uiConfig?.enableTagBackgroundImage ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? true;
  const compactExpandedDetails = uiConfig?.compactExpandedDetails ?? false;

  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);
  const loadStickyHeader = useLoadStickyHeader();

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [mergeType, setMergeType] = useState<"from" | "into" | undefined>();

  // Editing tag state
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [updateTag] = useTagUpdate();
  const [deleteTag] = useTagDestroy({ id: tag.id });

  const showAllCounts = uiConfig?.showChildTagContent;

  const tagImage = useMemo(() => {
    let existingImage = tag.image_path;
    if (isEditing) {
      if (image === null && existingImage) {
        const tagImageURL = new URL(existingImage);
        tagImageURL.searchParams.set("default", "true");
        return tagImageURL.toString();
      } else if (image) {
        return image;
      }
    }

    return existingImage;
  }, [isEditing, tag.image_path, image]);

  function setFavorite(v: boolean) {
    if (tag.id) {
      updateTag({
        variables: {
          input: {
            id: tag.id,
            favorite: v,
          },
        },
      });
    }
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => toggleEditing());
    Mousetrap.bind("d d", () => {
      setIsDeleteAlertOpen(true);
    });
    Mousetrap.bind(",", () => setCollapsed(!collapsed));
    Mousetrap.bind("f", () => setFavorite(!tag.favorite));

    return () => {
      if (isEditing) {
        Mousetrap.unbind("s s");
      }

      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
      Mousetrap.unbind(",");
      Mousetrap.unbind("f");
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
      Toast.success(
        intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "tag" }).toLocaleLowerCase() }
        )
      );
    }
  }

  async function onAutoTag() {
    if (!tag.id) return;
    try {
      await mutateMetadataAutoTag({ tags: [tag.id] });
      Toast.success(intl.formatMessage({ id: "toast.started_auto_tagging" }));
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

  function toggleEditing(value?: boolean) {
    if (value !== undefined) {
      setIsEditing(value);
    } else {
      setIsEditing((e) => !e);
    }
    setImage(undefined);
  }

  function renderMergeButton() {
    return (
      <Dropdown>
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

  const headerClassName = cx("detail-header", {
    edit: isEditing,
    collapsed,
    "full-width": !collapsed && !compactExpandedDetails,
  });

  return (
    <div id="tag-page" className="row">
      <Helmet>
        <title>{tag.name}</title>
      </Helmet>

      <div className={headerClassName}>
        <BackgroundImage
          imagePath={tag.image_path ?? undefined}
          show={enableBackgroundImage && !isEditing}
        />
        <div className="detail-container">
          <HeaderImage encodingImage={encodingImage}>
            {tagImage && (
              <DetailImage className="logo" alt={tag.name} src={tagImage} />
            )}
          </HeaderImage>
          <div className="row">
            <div className="tag-head col">
              <DetailTitle name={tag.name} classNamePrefix="tag">
                {!isEditing && (
                  <ExpandCollapseButton
                    collapsed={collapsed}
                    setCollapsed={(v) => setCollapsed(v)}
                  />
                )}
                <span className="name-icons">
                  <FavoriteIcon
                    favorite={tag.favorite}
                    onToggleFavorite={(v) => setFavorite(v)}
                  />
                </span>
              </DetailTitle>

              <AliasList aliases={tag.aliases} />
              {!isEditing && (
                <TagDetailsPanel
                  tag={tag}
                  fullWidth={!collapsed && !compactExpandedDetails}
                />
              )}
              {isEditing ? (
                <TagEditPanel
                  tag={tag}
                  onSubmit={onSave}
                  onCancel={() => toggleEditing()}
                  onDelete={onDelete}
                  setImage={setImage}
                  setEncodingImage={setEncodingImage}
                />
              ) : (
                <DetailsEditNavbar
                  objectName={tag.name}
                  isNew={false}
                  isEditing={isEditing}
                  onToggleEdit={() => toggleEditing()}
                  onSave={() => {}}
                  onImageChange={() => {}}
                  onClearImage={() => {}}
                  onAutoTag={onAutoTag}
                  autoTagDisabled={tag.ignore_auto_tag}
                  onDelete={onDelete}
                  classNames="mb-2"
                  customButtons={renderMergeButton()}
                />
              )}
            </div>
          </div>
        </div>
      </div>

      {!isEditing && loadStickyHeader && (
        <CompressedTagDetailsPanel tag={tag} />
      )}

      <div className="detail-body">
        <div className="tag-body">
          <div className="tag-tabs">
            {!isEditing && (
              <TagTabs
                tabKey={tabKey}
                tag={tag}
                abbreviateCounter={abbreviateCounter}
                showAllCounts={showAllCounts}
              />
            )}
          </div>
        </div>
      </div>
      {renderDeleteAlert()}
      {renderMergeDialog()}
    </div>
  );
};

const TagLoader: React.FC<RouteComponentProps<ITagParams>> = ({
  location,
  match,
}) => {
  const { id, tab } = match.params;
  const { data, loading, error } = useFindTag(id);

  useScrollToTopOnMount();

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findTag)
    return <ErrorMessage error={`No tag found with id ${id}.`} />;

  if (tab && !isTabKey(tab)) {
    return (
      <Redirect
        to={{
          ...location,
          pathname: `/tags/${id}`,
        }}
      />
    );
  }

  return <TagPage tag={data.findTag} tabKey={tab as TabKey | undefined} />;
};

export default TagLoader;
