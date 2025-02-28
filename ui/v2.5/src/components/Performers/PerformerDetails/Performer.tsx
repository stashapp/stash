import React, { useEffect, useMemo, useState } from "react";
import { Tabs, Tab, Col, Row } from "react-bootstrap";
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
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
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
import { useRatingKeybinds } from "src/hooks/keybinds";
import { DetailImage } from "src/components/Shared/DetailImage";
import { useLoadStickyHeader } from "src/hooks/detailsPanel";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { ExternalLinkButtons } from "src/components/Shared/ExternalLinksButton";
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
import { LightboxLink } from "src/hooks/Lightbox/LightboxLink";

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

const PerformerTabs: React.FC<{
  tabKey?: TabKey;
  performer: GQL.PerformerDataFragment;
  abbreviateCounter: boolean;
}> = ({ tabKey, performer, abbreviateCounter }) => {
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

  const { setTabKey } = useTabKey({
    tabKey,
    validTabs,
    defaultTabKey: populatedDefaultTab,
    baseURL: `/performers/${performer.id}`,
  });

  useEffect(() => {
    Mousetrap.bind("c", () => setTabKey("scenes"));
    Mousetrap.bind("g", () => setTabKey("galleries"));
    Mousetrap.bind("m", () => setTabKey("groups"));

    return () => {
      Mousetrap.unbind("c");
      Mousetrap.unbind("g");
      Mousetrap.unbind("m");
    };
  });

  return (
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
          <TabTitleCounter
            messageID="scenes"
            count={performer.scene_count}
            abbreviateCounter={abbreviateCounter}
          />
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
          <TabTitleCounter
            messageID="galleries"
            count={performer.gallery_count}
            abbreviateCounter={abbreviateCounter}
          />
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
          <TabTitleCounter
            messageID="images"
            count={performer.image_count}
            abbreviateCounter={abbreviateCounter}
          />
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
          <TabTitleCounter
            messageID="groups"
            count={performer.group_count}
            abbreviateCounter={abbreviateCounter}
          />
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
          <TabTitleCounter
            messageID="appears_with"
            count={performer.performer_count}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <PerformerAppearsWithPanel
          active={tabKey === "appearswith"}
          performer={performer}
        />
      </Tab>
    </Tabs>
  );
};

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

  const [updatePerformer] = usePerformerUpdate();
  const [deletePerformer, { loading: isDestroying }] = usePerformerDestroy();

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
    Mousetrap.bind("f", () => setFavorite(!performer.favorite));
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      Mousetrap.unbind("e");
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
        <BackgroundImage
          imagePath={activeImage ?? undefined}
          show={enableBackgroundImage && !isEditing}
        />
        <div className="detail-container">
          <HeaderImage encodingImage={encodingImage}>
            {!!activeImage && (
              <LightboxLink images={lightboxImages}>
                <DetailImage
                  className="performer"
                  src={activeImage}
                  alt={performer.name}
                />
              </LightboxLink>
            )}
          </HeaderImage>

          <div className="row">
            <div className="performer-head col">
              <DetailTitle
                name={performer.name}
                disambiguation={performer.disambiguation ?? undefined}
                classNamePrefix="performer"
              >
                {!isEditing && (
                  <ExpandCollapseButton
                    collapsed={collapsed}
                    setCollapsed={(v) => setCollapsed(v)}
                  />
                )}
                <span className="name-icons">
                  <FavoriteIcon
                    favorite={performer.favorite}
                    onToggleFavorite={(v) => setFavorite(v)}
                  />
                  <ExternalLinkButtons urls={performer.urls ?? undefined} />
                </span>
              </DetailTitle>
              <AliasList aliases={performer.alias_list} />
              <RatingSystem
                value={performer.rating100}
                onSetRating={(value) => setRating(value)}
                clickToRate
                withoutContext
              />
              {!isEditing && (
                <PerformerDetailsPanel
                  performer={performer}
                  collapsed={collapsed}
                  fullWidth={!collapsed && !compactExpandedDetails}
                />
              )}
              {isEditing ? (
                <PerformerEditPanel
                  performer={performer}
                  isVisible={isEditing}
                  onSubmit={onSave}
                  onCancel={() => toggleEditing()}
                  setImage={setImage}
                  setEncodingImage={setEncodingImage}
                />
              ) : (
                <Col>
                  <Row xs={8}>
                    <DetailsEditNavbar
                      objectName={
                        performer?.name ??
                        intl.formatMessage({ id: "performer" })
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
              )}
            </div>
          </div>
        </div>
      </div>

      {!isEditing && loadStickyHeader && (
        <CompressedPerformerDetailsPanel performer={performer} />
      )}

      <div className="detail-body">
        <div className="performer-body">
          <div className="performer-tabs">
            {!isEditing && (
              <PerformerTabs
                tabKey={tabKey}
                performer={performer}
                abbreviateCounter={abbreviateCounter}
              />
            )}
          </div>
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
