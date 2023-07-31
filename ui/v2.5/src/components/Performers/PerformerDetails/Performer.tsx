import React, { useEffect, useMemo, useState } from "react";
import { Button, Tabs, Tab, Col, Row } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useParams, useHistory } from "react-router-dom";
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
import TextUtils from "src/utils/text";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import {
  CompressedPerformerDetailsPanel,
  PerformerDetailsPanel,
} from "./PerformerDetailsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";
import { PerformerGalleriesPanel } from "./PerformerGalleriesPanel";
import { PerformerMoviesPanel } from "./PerformerMoviesPanel";
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
import { IUIConfig } from "src/core/config";
import { useRatingKeybinds } from "src/hooks/keybinds";
import ImageUtils from "src/utils/image";

interface IProps {
  performer: GQL.PerformerDataFragment;
}
interface IPerformerParams {
  tab?: string;
}

const PerformerPage: React.FC<IProps> = ({ performer }) => {
  const Toast = useToast();
  const history = useHistory();
  const intl = useIntl();
  const { tab = "details" } = useParams<IPerformerParams>();

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui as IUIConfig | undefined;
  const abbreviateCounter = uiConfig?.abbreviateCounters ?? false;
  const enableBackgroundImage =
    uiConfig?.enablePerformerBackgroundImage ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? false;

  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);
  const [loadStickyHeader, setLoadStickyHeader] = useState<boolean>(false);

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

  const activeTabKey =
    tab === "scenes" ||
    tab === "galleries" ||
    tab === "images" ||
    tab === "movies" ||
    tab == "appearswith"
      ? tab
      : "scenes";
  const setActiveTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "scenes" ? "" : `/${newTab}`;
      history.replace(`/performers/${performer.id}${tabParam}`);
    }
  };

  async function onAutoTag() {
    try {
      await mutateMetadataAutoTag({ performers: [performer.id] });
      Toast.success({
        content: intl.formatMessage({ id: "toast.started_auto_tagging" }),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  useRatingKeybinds(
    true,
    configuration?.ui?.ratingSystemOptions?.type,
    setRating
  );

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => toggleEditing());
    Mousetrap.bind("c", () => setActiveTabKey("scenes"));
    Mousetrap.bind("g", () => setActiveTabKey("galleries"));
    Mousetrap.bind("m", () => setActiveTabKey("movies"));
    Mousetrap.bind("f", () => setFavorite(!performer.favorite));
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("e");
      Mousetrap.unbind("c");
      Mousetrap.unbind("f");
      Mousetrap.unbind("o");
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
    Toast.success({
      content: intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "performer" }).toLocaleLowerCase() }
      ),
    });
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
          <img
            className="performer"
            src={activeImage}
            alt={performer.name}
            onLoad={ImageUtils.verifyImageSize}
          />
        </Button>
      );
    }
  }
  const renderTabs = () => (
    <React.Fragment>
      <Tabs
        activeKey={activeTabKey}
        onSelect={setActiveTabKey}
        id="performer-details"
        unmountOnExit
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
            active={activeTabKey == "scenes"}
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
            active={activeTabKey == "galleries"}
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
            active={activeTabKey == "images"}
            performer={performer}
          />
        </Tab>
        <Tab
          eventKey="movies"
          title={
            <>
              {intl.formatMessage({ id: "movies" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={performer.movie_count}
                hideZero
              />
            </>
          }
        >
          <PerformerMoviesPanel
            active={activeTabKey == "movies"}
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
            active={activeTabKey == "appearswith"}
            performer={performer}
          />
        </Tab>
      </Tabs>
    </React.Fragment>
  );

  function maybeRenderHeaderBackgroundImage() {
    if (enableBackgroundImage && !isEditing && activeImage) {
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

  function maybeRenderDetails() {
    if (!isEditing) {
      return (
        <PerformerDetailsPanel
          performer={performer}
          collapsed={collapsed}
          fullWidth={!collapsed && showAllDetails}
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
    /* Collect urls adding into details */
    /* This code can be removed once multple urls are supported for performers */
    const detailURLsRegex = /\[((?:http|www\.)[^\n\]]+)\]/gm;
    let urls = performer?.details?.match(detailURLsRegex);

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
        {performer.url && (
          <Button className="minimal icon-link" title={performer.url}>
            <a
              href={TextUtils.sanitiseURL(performer.url)}
              className="link"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Icon icon={faLink} />
            </a>
          </Button>
        )}
        {(urls ?? []).map((url, index) => (
          <Button key={index} className="minimal icon-link" title={url}>
            <a
              href={TextUtils.sanitiseURL(url)}
              className={`detail-link ${index}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              <Icon icon={faLink} />
            </a>
          </Button>
        ))}
        {performer.twitter && (
          <Button className="minimal icon-link" title={performer.twitter}>
            <a
              href={TextUtils.sanitiseURL(
                performer.twitter,
                TextUtils.twitterURL
              )}
              className="twitter"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Icon icon={faTwitter} />
            </a>
          </Button>
        )}
        {performer.instagram && (
          <Button className="minimal icon-link" title={performer.instagram}>
            <a
              href={TextUtils.sanitiseURL(
                performer.instagram,
                TextUtils.instagramURL
              )}
              className="instagram"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Icon icon={faInstagram} />
            </a>
          </Button>
        )}
      </span>
    );
  }

  if (isDestroying)
    return (
      <LoadingIndicator
        message={`Deleting performer ${performer.id}: ${performer.name}`}
      />
    );

  return (
    <div id="performer-page" className="row">
      <Helmet>
        <title>{performer.name}</title>
      </Helmet>

      <div
        className={`detail-header ${isEditing ? "edit" : ""}  ${
          collapsed ? "collapsed" : showAllDetails ? "full-width" : ""
        }`}
      >
        {maybeRenderHeaderBackgroundImage()}
        <div className="detail-container">
          <div className="detail-header-image">
            {encodingImage ? (
              <LoadingIndicator message="Encoding image..." />
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
                value={performer.rating100 ?? undefined}
                onSetRating={(value) => setRating(value ?? null)}
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

const PerformerLoader: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const { data, loading, error } = useFindPerformer(id ?? "");

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findPerformer)
    return <ErrorMessage error={`No performer found with id ${id}.`} />;

  return <PerformerPage performer={data.findPerformer} />;
};

export default PerformerLoader;
