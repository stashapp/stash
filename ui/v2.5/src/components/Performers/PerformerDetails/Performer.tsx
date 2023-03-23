import React, { useEffect, useMemo, useState } from "react";
import { Button, Tabs, Tab, Col, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
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
import { CountryFlag } from "src/components/Shared/CountryFlag";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import TextUtils from "src/utils/text";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { PerformerDetailsPanel } from "./PerformerDetailsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";
import { PerformerGalleriesPanel } from "./PerformerGalleriesPanel";
import { PerformerMoviesPanel } from "./PerformerMoviesPanel";
import { PerformerImagesPanel } from "./PerformerImagesPanel";
import { PerformerEditPanel } from "./PerformerEditPanel";
import { PerformerSubmitButton } from "./PerformerSubmitButton";
import GenderIcon from "../GenderIcon";
import { faHeart, faLink } from "@fortawesome/free-solid-svg-icons";
import { faInstagram, faTwitter } from "@fortawesome/free-brands-svg-icons";
import { IUIConfig } from "src/core/config";
import { useRatingKeybinds } from "src/hooks/keybinds";

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

  const [collapsed, setCollapsed] = useState(false);

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const abbreviateCounter =
    (configuration?.ui as IUIConfig)?.abbreviateCounters ?? false;

  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  // if undefined then get the existing image
  // if null then get the default (no) image
  // otherwise get the set image
  const activeImage =
    image === undefined
      ? performer.image_path ?? ""
      : image ?? `${performer.image_path}&default=true`;
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
    tab === "movies"
      ? tab
      : "details";
  const setActiveTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "details" ? "" : `/${newTab}`;
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
    Mousetrap.bind("a", () => setActiveTabKey("details"));
    Mousetrap.bind("e", () => setIsEditing(!isEditing));
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

  async function onDelete() {
    try {
      await deletePerformer({ variables: { id: performer.id } });
    } catch (e) {
      Toast.error(e);
    }

    // redirect to performers page
    history.push("/performers");
  }

  const renderTabs = () => (
    <React.Fragment>
      <Col>
        <Row xs={8}>
          <DetailsEditNavbar
            objectName={
              performer?.name ?? intl.formatMessage({ id: "performer" })
            }
            onToggleEdit={() => {
              setIsEditing(!isEditing);
            }}
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
      <Tabs
        activeKey={activeTabKey}
        onSelect={setActiveTabKey}
        id="performer-details"
        unmountOnExit
      >
        <Tab eventKey="details" title={intl.formatMessage({ id: "details" })}>
          <PerformerDetailsPanel performer={performer} />
        </Tab>
        <Tab
          eventKey="scenes"
          title={
            <React.Fragment>
              {intl.formatMessage({ id: "scenes" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={performer.scene_count ?? 0}
              />
            </React.Fragment>
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
            <React.Fragment>
              {intl.formatMessage({ id: "galleries" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={performer.gallery_count ?? 0}
              />
            </React.Fragment>
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
            <React.Fragment>
              {intl.formatMessage({ id: "images" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={performer.image_count ?? 0}
              />
            </React.Fragment>
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
            <React.Fragment>
              {intl.formatMessage({ id: "movies" })}
              <Counter
                abbreviateCounter={abbreviateCounter}
                count={performer.movie_count ?? 0}
              />
            </React.Fragment>
          }
        >
          <PerformerMoviesPanel
            active={activeTabKey == "movies"}
            performer={performer}
          />
        </Tab>
      </Tabs>
    </React.Fragment>
  );

  function renderTabsOrEditPanel() {
    if (isEditing) {
      return (
        <PerformerEditPanel
          performer={performer}
          isVisible={isEditing}
          onCancel={() => setIsEditing(false)}
          setImage={setImage}
          setEncodingImage={setEncodingImage}
        />
      );
    } else {
      return renderTabs();
    }
  }

  function maybeRenderAge() {
    if (performer?.birthdate) {
      // calculate the age from birthdate. In future, this should probably be
      // provided by the server
      return (
        <div>
          <span className="age">
            {TextUtils.age(performer.birthdate, performer.death_date)}
          </span>
          <span className="age-tail">
            {" "}
            <FormattedMessage id="years_old" />
          </span>
        </div>
      );
    }
  }

  function maybeRenderAliases() {
    if (performer?.alias_list?.length) {
      return (
        <div>
          <span className="alias-head">
            <FormattedMessage id="also_known_as" />{" "}
          </span>
          <span className="alias">{performer.alias_list?.join(", ")}</span>
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

  const renderClickableIcons = () => (
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
        <Button className="minimal icon-link">
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
      {performer.twitter && (
        <Button className="minimal icon-link">
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
        <Button className="minimal icon-link">
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

  if (isDestroying)
    return (
      <LoadingIndicator
        message={`Deleting performer ${performer.id}: ${performer.name}`}
      />
    );

  function getCollapseButtonText() {
    return collapsed ? ">" : "<";
  }

  return (
    <div id="performer-page" className="row">
      <Helmet>
        <title>{performer.name}</title>
      </Helmet>

      <div
        className={`performer-image-container details-tab text-center text-center ${
          collapsed ? "collapsed" : ""
        }`}
      >
        {encodingImage ? (
          <LoadingIndicator message="Encoding image..." />
        ) : (
          <Button variant="link" onClick={() => showLightbox()}>
            <img
              className="performer"
              src={activeImage}
              alt={intl.formatMessage({ id: "performer" })}
            />
          </Button>
        )}
      </div>
      <div className="details-divider d-none d-xl-block">
        <Button onClick={() => setCollapsed(!collapsed)}>
          {getCollapseButtonText()}
        </Button>
      </div>
      <div className={`content-container ${collapsed ? "expanded" : ""}`}>
        <div className="row">
          <div className="performer-head col">
            <h2>
              <GenderIcon
                gender={performer.gender}
                className="gender-icon mr-2 fi"
              />
              <CountryFlag country={performer.country} className="mr-2" />
              <span className="performer-name">{performer.name}</span>
              {performer.disambiguation && (
                <span className="performer-disambiguation">
                  {` (${performer.disambiguation})`}
                </span>
              )}
              {renderClickableIcons()}
            </h2>
            <RatingSystem
              value={performer.rating100 ?? undefined}
              onSetRating={(value) => setRating(value ?? null)}
            />
            {maybeRenderAliases()}
            {maybeRenderAge()}
          </div>
        </div>
        <div className="performer-body">
          <div className="performer-tabs">{renderTabsOrEditPanel()}</div>
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
