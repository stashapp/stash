import { Button, Tabs, Tab } from "react-bootstrap";
import React, { useCallback, useEffect, useMemo, useState } from "react";
import { useHistory, Redirect, RouteComponentProps } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import cx from "classnames";
import Mousetrap from "mousetrap";

import * as GQL from "src/core/generated-graphql";
import {
  useFindStudio,
  useStudioUpdate,
  useStudioDestroy,
  mutateMetadataAutoTag,
} from "src/core/StashService";
import { Counter } from "src/components/Shared/Counter";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import { Icon } from "src/components/Shared/Icon";
import { StudioScenesPanel } from "./StudioScenesPanel";
import { StudioGalleriesPanel } from "./StudioGalleriesPanel";
import { StudioImagesPanel } from "./StudioImagesPanel";
import { StudioChildrenPanel } from "./StudioChildrenPanel";
import { StudioPerformersPanel } from "./StudioPerformersPanel";
import { StudioEditPanel } from "./StudioEditPanel";
import {
  CompressedStudioDetailsPanel,
  StudioDetailsPanel,
} from "./StudioDetailsPanel";
import { StudioGroupsPanel } from "./StudioMoviesPanel";
import {
  faTrashAlt,
  faLink,
  faChevronDown,
  faChevronUp,
  faHeart,
} from "@fortawesome/free-solid-svg-icons";
import TextUtils from "src/utils/text";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { DetailImage } from "src/components/Shared/DetailImage";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { useLoadStickyHeader } from "src/hooks/detailsPanel";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { ExternalLink } from "src/components/Shared/ExternalLink";

interface IProps {
  studio: GQL.StudioDataFragment;
  tabKey?: TabKey;
}

interface IStudioParams {
  id: string;
  tab?: string;
}

const validTabs = [
  "default",
  "scenes",
  "galleries",
  "images",
  "performers",
  "groups",
  "childstudios",
] as const;
type TabKey = (typeof validTabs)[number];

function isTabKey(tab: string): tab is TabKey {
  return validTabs.includes(tab as TabKey);
}

const StudioPage: React.FC<IProps> = ({ studio, tabKey }) => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui;
  const abbreviateCounter = uiConfig?.abbreviateCounters ?? false;
  const enableBackgroundImage = uiConfig?.enableStudioBackgroundImage ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? true;
  const compactExpandedDetails = uiConfig?.compactExpandedDetails ?? false;

  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);
  const loadStickyHeader = useLoadStickyHeader();

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing studio state
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [updateStudio] = useStudioUpdate();
  const [deleteStudio] = useStudioDestroy({ id: studio.id });

  const showAllCounts = uiConfig?.showChildStudioContent;
  const sceneCount =
    (showAllCounts ? studio.scene_count_all : studio.scene_count) ?? 0;
  const galleryCount =
    (showAllCounts ? studio.gallery_count_all : studio.gallery_count) ?? 0;
  const imageCount =
    (showAllCounts ? studio.image_count_all : studio.image_count) ?? 0;
  const performerCount =
    (showAllCounts ? studio.performer_count_all : studio.performer_count) ?? 0;
  const groupCount =
    (showAllCounts ? studio.movie_count_all : studio.movie_count) ?? 0;

  const populatedDefaultTab = useMemo(() => {
    let ret: TabKey = "scenes";
    if (sceneCount == 0) {
      if (galleryCount != 0) {
        ret = "galleries";
      } else if (imageCount != 0) {
        ret = "images";
      } else if (performerCount != 0) {
        ret = "performers";
      } else if (groupCount != 0) {
        ret = "groups";
      } else if (studio.child_studios.length != 0) {
        ret = "childstudios";
      }
    }

    return ret;
  }, [
    sceneCount,
    galleryCount,
    imageCount,
    performerCount,
    groupCount,
    studio,
  ]);

  const setTabKey = useCallback(
    (newTabKey: string | null) => {
      if (!newTabKey) newTabKey = populatedDefaultTab;
      if (newTabKey === tabKey) return;

      if (isTabKey(newTabKey)) {
        history.replace(`/studios/${studio.id}/${newTabKey}`);
      }
    },
    [populatedDefaultTab, tabKey, history, studio.id]
  );

  useEffect(() => {
    if (!tabKey) {
      setTabKey(populatedDefaultTab);
    }
  }, [setTabKey, populatedDefaultTab, tabKey]);

  function setFavorite(v: boolean) {
    if (studio.id) {
      updateStudio({
        variables: {
          input: {
            id: studio.id,
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
    Mousetrap.bind("f", () => setFavorite(!studio.favorite));

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
      Mousetrap.unbind(",");
      Mousetrap.unbind("f");
    };
  });

  useRatingKeybinds(
    true,
    configuration?.ui.ratingSystemOptions?.type,
    setRating
  );

  async function onSave(input: GQL.StudioCreateInput) {
    await updateStudio({
      variables: {
        input: {
          id: studio.id,
          ...input,
        },
      },
    });
    toggleEditing(false);
    Toast.success(
      intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "studio" }).toLocaleLowerCase() }
      )
    );
  }

  async function onAutoTag() {
    if (!studio.id) return;
    try {
      await mutateMetadataAutoTag({ studios: [studio.id] });
      Toast.success(intl.formatMessage({ id: "toast.started_auto_tagging" }));
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onDelete() {
    try {
      await deleteStudio();
    } catch (e) {
      Toast.error(e);
    }

    // redirect to studios page
    history.push(`/studios`);
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
                studio.name ??
                intl.formatMessage({ id: "studio" }).toLocaleLowerCase(),
            }}
          />
        </p>
      </ModalComponent>
    );
  }

  function maybeRenderAliases() {
    if (studio?.aliases?.length) {
      return (
        <div>
          <span className="alias-head">{studio?.aliases?.join(", ")}</span>
        </div>
      );
    }
  }

  function getCollapseButtonIcon() {
    return collapsed ? faChevronDown : faChevronUp;
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
    let studioImage = studio.image_path;
    if (isEditing) {
      if (image === null && studioImage) {
        const studioImageURL = new URL(studioImage);
        studioImageURL.searchParams.set("default", "true");
        studioImage = studioImageURL.toString();
      } else if (image) {
        studioImage = image;
      }
    }

    if (studioImage) {
      return (
        <DetailImage className="logo" alt={studio.name} src={studioImage} />
      );
    }
  }

  const renderClickableIcons = () => (
    <span className="name-icons">
      <Button
        className={cx("minimal", studio.favorite ? "favorite" : "not-favorite")}
        onClick={() => setFavorite(!studio.favorite)}
      >
        <Icon icon={faHeart} />
      </Button>
      {studio.url && (
        <Button
          as={ExternalLink}
          href={TextUtils.sanitiseURL(studio.url)}
          className="minimal link"
          title={studio.url}
        >
          <Icon icon={faLink} />
        </Button>
      )}
    </span>
  );

  function setRating(v: number | null) {
    if (studio.id) {
      updateStudio({
        variables: {
          input: {
            id: studio.id,
            rating100: v,
          },
        },
      });
    }
  }

  function maybeRenderDetails() {
    if (!isEditing) {
      return (
        <StudioDetailsPanel
          studio={studio}
          collapsed={collapsed}
          fullWidth={!collapsed && !compactExpandedDetails}
        />
      );
    }
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

  function maybeRenderCompressedDetails() {
    if (!isEditing && loadStickyHeader) {
      return <CompressedStudioDetailsPanel studio={studio} />;
    }
  }

  const renderTabs = () => (
    <Tabs
      id="studio-tabs"
      mountOnEnter
      unmountOnExit
      activeKey={tabKey}
      onSelect={setTabKey}
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
        <StudioScenesPanel active={tabKey === "scenes"} studio={studio} />
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
        <StudioGalleriesPanel active={tabKey === "galleries"} studio={studio} />
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
        <StudioImagesPanel active={tabKey === "images"} studio={studio} />
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
        <StudioPerformersPanel
          active={tabKey === "performers"}
          studio={studio}
        />
      </Tab>
      <Tab
        eventKey="groups"
        title={
          <>
            {intl.formatMessage({ id: "groups" })}
            <Counter
              abbreviateCounter={abbreviateCounter}
              count={groupCount}
              hideZero
            />
          </>
        }
      >
        <StudioGroupsPanel active={tabKey === "groups"} studio={studio} />
      </Tab>
      <Tab
        eventKey="childstudios"
        title={
          <>
            {intl.formatMessage({ id: "subsidiary_studios" })}
            <Counter
              abbreviateCounter={false}
              count={studio.child_studios.length}
              hideZero
            />
          </>
        }
      >
        <StudioChildrenPanel
          active={tabKey === "childstudios"}
          studio={studio}
        />
      </Tab>
    </Tabs>
  );

  function maybeRenderHeaderBackgroundImage() {
    let studioImage = studio.image_path;
    if (enableBackgroundImage && !isEditing && studioImage) {
      const studioImageURL = new URL(studioImage);
      let isDefaultImage = studioImageURL.searchParams.get("default");
      if (!isDefaultImage) {
        return (
          <div className="background-image-container">
            <picture>
              <source src={studioImage} />
              <img
                className="background-image"
                src={studioImage}
                alt={`${studio.name} background`}
              />
            </picture>
          </div>
        );
      }
    }
  }

  function maybeRenderTab() {
    if (!isEditing) {
      return renderTabs();
    }
  }

  function maybeRenderEditPanel() {
    if (isEditing) {
      return (
        <StudioEditPanel
          studio={studio}
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
          objectName={studio.name ?? intl.formatMessage({ id: "studio" })}
          isNew={false}
          isEditing={isEditing}
          onToggleEdit={() => toggleEditing()}
          onSave={() => {}}
          onImageChange={() => {}}
          onClearImage={() => {}}
          onAutoTag={onAutoTag}
          autoTagDisabled={studio.ignore_auto_tag}
          onDelete={onDelete}
        />
      );
    }
  }

  const headerClassName = cx("detail-header", {
    edit: isEditing,
    collapsed,
    "full-width": !collapsed && !compactExpandedDetails,
  });

  return (
    <div id="studio-page" className="row">
      <Helmet>
        <title>{studio.name ?? intl.formatMessage({ id: "studio" })}</title>
      </Helmet>

      <div className={headerClassName}>
        {maybeRenderHeaderBackgroundImage()}
        <div className="detail-container">
          <div className="detail-header-image">
            {encodingImage ? (
              <LoadingIndicator
                message={intl.formatMessage({ id: "actions.encoding_image" })}
              />
            ) : (
              renderImage()
            )}
          </div>
          <div className="row">
            <div className="studio-head col">
              <h2>
                <span className="studio-name">{studio.name}</span>
                {maybeRenderShowCollapseButton()}
                {renderClickableIcons()}
              </h2>
              {maybeRenderAliases()}
              <RatingSystem
                value={studio.rating100}
                onSetRating={(value) => setRating(value)}
                clickToRate
                withoutContext
              />
              {maybeRenderDetails()}
              {maybeRenderEditPanel()}
            </div>
          </div>
        </div>
      </div>
      {maybeRenderCompressedDetails()}
      <div className="detail-body">
        <div className="studio-body">
          <div className="studio-tabs">{maybeRenderTab()}</div>
        </div>
      </div>
      {renderDeleteAlert()}
    </div>
  );
};

const StudioLoader: React.FC<RouteComponentProps<IStudioParams>> = ({
  location,
  match,
}) => {
  const { id, tab } = match.params;
  const { data, loading, error } = useFindStudio(id);

  useScrollToTopOnMount();

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findStudio)
    return <ErrorMessage error={`No studio found with id ${id}.`} />;

  if (tab && !isTabKey(tab)) {
    return (
      <Redirect
        to={{
          ...location,
          pathname: `/studios/${id}`,
        }}
      />
    );
  }

  return (
    <StudioPage studio={data.findStudio} tabKey={tab as TabKey | undefined} />
  );
};

export default StudioLoader;
