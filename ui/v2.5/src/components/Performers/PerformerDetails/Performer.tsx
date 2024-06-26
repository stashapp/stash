import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Button, Tabs, Tab, Col, Row } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useHistory, Redirect, RouteComponentProps } from "react-router-dom";
import { Helmet } from "react-helmet";
import cx from "classnames";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  useFindPerformer,
  usePerformerUpdate,
  usePerformerDestroy,
  mutateMetadataAutoTag,
} from "src/core/StashService";
import { Counter } from "src/components/Shared/Counter";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import {
  CompressedPerformerDetailsPanel,
  PerformerDetailsPanel,
} from "./PerformerDetailsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";
import { PerformerGalleriesPanel } from "./PerformerGalleriesPanel";
import { PerformerGroupsPanel } from "./PerformerGroupsPanel";
import { PerformerImagesPanel } from "./PerformerImagesPanel";
import { PerformerAppearsWithPanel } from "./performerAppearsWithPanel";
import { PerformerEditPanel } from "./PerformerEditPanel";
import { PerformerSubmitButton } from "./PerformerSubmitButton";
import {
  faChevronDown,
  faChevronUp,
  faHeart,
  faLink,
} from "@fortawesome/free-solid-svg-icons";
import { faInstagram, faTwitter } from "@fortawesome/free-brands-svg-icons";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { DetailImage } from "src/components/Shared/DetailImage";
import { useLoadStickyHeader } from "src/hooks/detailsPanel";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { ExternalLinksButton } from "src/components/Shared/ExternalLinksButton";

interface IProps {
  performer: GQL.PerformerDataFragment;
  tabKey?: TabKey;
}

interface IPerformerParams {
  id: string;
  tab?: string;
}

const validTabs = [
  "default",
  "scenes",
  "galleries",
  "images",
  "groups",
  "appearswith",
] as const;
type TabKey = (typeof validTabs)[number];

function isTabKey(tab: string): tab is TabKey {
  return validTabs.includes(tab as TabKey);
}

const PerformerPage: React.FC<IProps> = ({ performer, tabKey }) => {
  const Toast = useToast();
  const history = useHistory();
  const intl = useIntl();

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui;
  const abbreviateCounter = uiConfig?.abbreviateCounters ?? false;
  const enableBackgroundImage =
    uiConfig?.enablePerformerBackgroundImage ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? true;
  const compactExpandedDetails = uiConfig?.compactExpandedDetails ?? false;

  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);
  const loadStickyHeader = useLoadStickyHeader();

  // a list of urls to display in the performer details
  const urls = useMemo(() => {
    if (!performer.urls?.length) {
      return [];
    }

    const twitter = performer.urls.filter((u) =>
      u.match(/https?:\/\/(?:www\.)?twitter.com\//)
    );
    const instagram = performer.urls.filter((u) =>
      u.match(/https?:\/\/(?:www\.)?instagram.com\//)
    );
    const others = performer.urls.filter(
      (u) => !twitter.includes(u) && !instagram.includes(u)
    );

    return [
      { icon: faLink, className: "", urls: others },
      { icon: faTwitter, className: "twitter", urls: twitter },
      { icon: faInstagram, className: "instagram", urls: instagram },
    ];
  }, [performer.urls]);

  const activeImage = useMemo(() => {
    const performerImage = performer.image_path;
    if (isEditing) {
      if (image === null && performerImage) {
        const performerImageURL = new URL(performerImage);
        performerImageURL.searchParams.set("default", "true");
        return performerImageURL.toString();
      } else if (image) {
        return image;
      }
    }
    return performerImage;
  }, [image, isEditing, performer.image_path]);

  const lightboxImages = useMemo(
    () => [{ paths: { thumbnail: activeImage, image: activeImage } }],
    [activeImage]
  );

  const showLightbox = useLightbox({
    images: lightboxImages,
  });

  const [updatePerformer] = usePerformerUpdate();
  const [deletePerformer, { loading: isDestroying }] = usePerformerDestroy();

  const populatedDefaultTab = useMemo(() => {
    let ret: TabKey = "scenes";
    if (performer.scene_count == 0) {
      if (performer.gallery_count != 0) {
        ret = "galleries";
      } else if (performer.image_count != 0) {
        ret = "images";
      } else if (performer.group_count != 0) {
        ret = "groups";
      }
    }

    return ret;
  }, [performer]);

  const setTabKey = useCallback(
    (newTabKey: string | null) => {
      if (!newTabKey) newTabKey = populatedDefaultTab;
      if (newTabKey === tabKey) return;

      if (isTabKey(newTabKey)) {
        history.replace(`/performers/${performer.id}/${newTabKey}`);
      }
    },
    [populatedDefaultTab, tabKey, history, performer.id]
  );

  useEffect(() => {
    if (!tabKey) {
      setTabKey(populatedDefaultTab);
    }
  }, [setTabKey, populatedDefaultTab, tabKey]);

  async function onAutoTag() {
    try {
      await mutateMetadataAutoTag({ performers: [performer.id] });
      Toast.success(intl.formatMessage({ id: "toast.started_auto_tagging" }));
    } catch (e) {
      Toast.error(e);
    }
  }

  useRatingKeybinds(
    true,
    configuration?.ui.ratingSystemOptions?.type,
    setRating
  );

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => toggleEditing());
    Mousetrap.bind("c", () => setTabKey("scenes"));
    Mousetrap.bind("g", () => setTabKey("galleries"));
    Mousetrap.bind("m", () => setTabKey("groups"));
    Mousetrap.bind("f", () => setFavorite(!performer.favorite));
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("c");
      Mousetrap.unbind("g");
      Mousetrap.unbind("m");
      Mousetrap.unbind("f");
      Mousetrap.unbind(",");
    };
  });

  async function onSave(input: GQL.PerformerCreateInput) {
    await updatePerformer({
      variables: {
        input: {
          id: performer.id,
          ...input,
        },
      },
    });
    toggleEditing(false);
    Toast.success(
      intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "performer" }).toLocaleLowerCase() }
      )
    );
  }

  async function onDelete() {
    try {
      await deletePerformer({ variables: { id: performer.id } });
    } catch (e) {
      Toast.error(e);
    }

    // redirect to performers page
    history.push("/performers");
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
    if (activeImage) {
      return (
        <Button variant="link" onClick={() => showLightbox()}>
          <DetailImage
            className="performer"
            src={activeImage}
            alt={performer.name}
          />
        </Button>
      );
    }
  }
  const renderTabs = () => (
    <Tabs
      id="performer-tabs"
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
              count={performer.scene_count}
              hideZero
            />
          </>
        }
      >
        <PerformerScenesPanel
          active={tabKey === "scenes"}
          performer={performer}
        />
      </Tab>
      <Tab
        eventKey="galleries"
        title={
          <>
            {intl.formatMessage({ id: "galleries" })}
            <Counter
              abbreviateCounter={abbreviateCounter}
              count={performer.gallery_count}
              hideZero
            />
          </>
        }
      >
        <PerformerGalleriesPanel
          active={tabKey === "galleries"}
          performer={performer}
        />
      </Tab>
      <Tab
        eventKey="images"
        title={
          <>
            {intl.formatMessage({ id: "images" })}
            <Counter
              abbreviateCounter={abbreviateCounter}
              count={performer.image_count}
              hideZero
            />
          </>
        }
      >
        <PerformerImagesPanel
          active={tabKey === "images"}
          performer={performer}
        />
      </Tab>
      <Tab
        eventKey="groups"
        title={
          <>
            {intl.formatMessage({ id: "groups" })}
            <Counter
              abbreviateCounter={abbreviateCounter}
              count={performer.group_count}
              hideZero
            />
          </>
        }
      >
        <PerformerGroupsPanel
          active={tabKey === "groups"}
          performer={performer}
        />
      </Tab>
      <Tab
        eventKey="appearswith"
        title={
          <>
            {intl.formatMessage({ id: "appears_with" })}
            <Counter
              abbreviateCounter={abbreviateCounter}
              count={performer.performer_count}
              hideZero
            />
          </>
        }
      >
        <PerformerAppearsWithPanel
          active={tabKey === "appearswith"}
          performer={performer}
        />
      </Tab>
    </Tabs>
  );

  function maybeRenderHeaderBackgroundImage() {
    if (enableBackgroundImage && !isEditing && activeImage) {
      const activeImageURL = new URL(activeImage);
      let isDefaultImage = activeImageURL.searchParams.get("default");
      if (!isDefaultImage) {
        return (
          <div className="background-image-container">
            <picture>
              <source src={activeImage} />
              <img
                className="background-image"
                src={activeImage}
                alt={`${performer.name} background`}
              />
            </picture>
          </div>
        );
      }
    }
  }

  function maybeRenderEditPanel() {
    if (isEditing) {
      return (
        <PerformerEditPanel
          performer={performer}
          isVisible={isEditing}
          onSubmit={onSave}
          onCancel={() => toggleEditing()}
          setImage={setImage}
          setEncodingImage={setEncodingImage}
        />
      );
    }
    {
      return (
        <Col>
          <Row xs={8}>
            <DetailsEditNavbar
              objectName={
                performer?.name ?? intl.formatMessage({ id: "performer" })
              }
              onToggleEdit={() => toggleEditing()}
              onDelete={onDelete}
              onAutoTag={onAutoTag}
              autoTagDisabled={performer.ignore_auto_tag}
              isNew={false}
              isEditing={false}
              onSave={() => {}}
              onImageChange={() => {}}
              classNames="mb-2"
              customButtons={
                <div>
                  <PerformerSubmitButton performer={performer} />
                </div>
              }
            ></DetailsEditNavbar>
          </Row>
        </Col>
      );
    }
  }

  function getCollapseButtonIcon() {
    return collapsed ? faChevronDown : faChevronUp;
  }

  function maybeRenderDetails() {
    if (!isEditing) {
      return (
        <PerformerDetailsPanel
          performer={performer}
          collapsed={collapsed}
          fullWidth={!collapsed && !compactExpandedDetails}
        />
      );
    }
  }

  function maybeRenderCompressedDetails() {
    if (!isEditing && loadStickyHeader) {
      return <CompressedPerformerDetailsPanel performer={performer} />;
    }
  }

  function maybeRenderTab() {
    if (!isEditing) {
      return renderTabs();
    }
  }

  function maybeRenderAliases() {
    if (performer?.alias_list?.length) {
      return (
        <div>
          <span className="alias-head">{performer.alias_list?.join(", ")}</span>
        </div>
      );
    }
  }

  function setFavorite(v: boolean) {
    if (performer.id) {
      updatePerformer({
        variables: {
          input: {
            id: performer.id,
            favorite: v,
          },
        },
      });
    }
  }

  function setRating(v: number | null) {
    if (performer.id) {
      updatePerformer({
        variables: {
          input: {
            id: performer.id,
            rating100: v,
          },
        },
      });
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

  function renderClickableIcons() {
    return (
      <span className="name-icons">
        <Button
          className={cx(
            "minimal",
            performer.favorite ? "favorite" : "not-favorite"
          )}
          onClick={() => setFavorite(!performer.favorite)}
        >
          <Icon icon={faHeart} />
        </Button>
        {urls.map((url) => (
          <ExternalLinksButton
            key={url.icon.iconName}
            icon={url.icon}
            className={url.className}
            urls={url.urls}
          />
        ))}
      </span>
    );
  }

  if (isDestroying)
    return (
      <LoadingIndicator
        message={`Deleting performer ${performer.id}: ${performer.name}`}
      />
    );

  const headerClassName = cx("detail-header", {
    edit: isEditing,
    collapsed,
    "full-width": !collapsed && !compactExpandedDetails,
  });

  return (
    <div id="performer-page" className="row">
      <Helmet>
        <title>{performer.name}</title>
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
            <div className="performer-head col">
              <h2>
                <span className="performer-name">{performer.name}</span>
                {performer.disambiguation && (
                  <span className="performer-disambiguation">
                    {` (${performer.disambiguation})`}
                  </span>
                )}
                {maybeRenderShowCollapseButton()}
                {renderClickableIcons()}
              </h2>
              {maybeRenderAliases()}
              <RatingSystem
                value={performer.rating100}
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
        <div className="performer-body">
          <div className="performer-tabs">{maybeRenderTab()}</div>
        </div>
      </div>
    </div>
  );
};

const PerformerLoader: React.FC<RouteComponentProps<IPerformerParams>> = ({
  location,
  match,
}) => {
  const { id, tab } = match.params;
  const { data, loading, error } = useFindPerformer(id);

  useScrollToTopOnMount();

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findPerformer)
    return <ErrorMessage error={`No performer found with id ${id}.`} />;

  if (tab && !isTabKey(tab)) {
    return (
      <Redirect
        to={{
          ...location,
          pathname: `/performers/${id}`,
        }}
      />
    );
  }

  return (
    <PerformerPage
      performer={data.findPerformer}
      tabKey={tab as TabKey | undefined}
    />
  );
};

export default PerformerLoader;
