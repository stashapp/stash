import React, { useEffect, useMemo, useState } from "react";
import { Button, Tabs, Tab, Badge, Col, Row } from "react-bootstrap";
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
import {
  CountryFlag,
  DetailsEditNavbar,
  ErrorMessage,
  Icon,
  LoadingIndicator,
} from "src/components/Shared";
import { useLightbox, useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { PerformerDetailsPanel } from "./PerformerDetailsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";
import { PerformerGalleriesPanel } from "./PerformerGalleriesPanel";
import { PerformerMoviesPanel } from "./PerformerMoviesPanel";
import { PerformerImagesPanel } from "./PerformerImagesPanel";
import { PerformerEditPanel } from "./PerformerEditPanel";
import { PerformerSubmitButton } from "./PerformerSubmitButton";
import GenderIcon from "../GenderIcon";

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

  const [imagePreview, setImagePreview] = useState<string | null>();
  const [imageEncoding, setImageEncoding] = useState<boolean>(false);
  const [isEditing, setIsEditing] = useState<boolean>(false);

  // if undefined then get the existing image
  // if null then get the default (no) image
  // otherwise get the set image
  const activeImage =
    imagePreview === undefined
      ? performer.image_path ?? ""
      : imagePreview ?? `${performer.image_path}&default=true`;
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

  const onImageChange = (image?: string | null) => setImagePreview(image);

  const onImageEncoding = (isEncoding = false) => setImageEncoding(isEncoding);

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

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("details"));
    Mousetrap.bind("e", () => setIsEditing(!isEditing));
    Mousetrap.bind("c", () => setActiveTabKey("scenes"));
    Mousetrap.bind("g", () => setActiveTabKey("galleries"));
    Mousetrap.bind("m", () => setActiveTabKey("movies"));
    Mousetrap.bind("f", () => setFavorite(!performer.favorite));

    // numeric keypresses get caught by jwplayer, so blur the element
    // if the rating sequence is started
    Mousetrap.bind("r", () => {
      if (document.activeElement instanceof HTMLElement) {
        document.activeElement.blur();
      }

      Mousetrap.bind("0", () => setRating(NaN));
      Mousetrap.bind("1", () => setRating(1));
      Mousetrap.bind("2", () => setRating(2));
      Mousetrap.bind("3", () => setRating(3));
      Mousetrap.bind("4", () => setRating(4));
      Mousetrap.bind("5", () => setRating(5));

      setTimeout(() => {
        Mousetrap.unbind("0");
        Mousetrap.unbind("1");
        Mousetrap.unbind("2");
        Mousetrap.unbind("3");
        Mousetrap.unbind("4");
        Mousetrap.unbind("5");
      }, 1000);
    });

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("e");
      Mousetrap.unbind("c");
      Mousetrap.unbind("f");
      Mousetrap.unbind("o");
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
              <Badge className="left-spacing" pill variant="secondary">
                {intl.formatNumber(performer.scene_count ?? 0)}
              </Badge>
            </React.Fragment>
          }
        >
          <PerformerScenesPanel performer={performer} />
        </Tab>
        <Tab
          eventKey="galleries"
          title={
            <React.Fragment>
              {intl.formatMessage({ id: "galleries" })}
              <Badge className="left-spacing" pill variant="secondary">
                {intl.formatNumber(performer.gallery_count ?? 0)}
              </Badge>
            </React.Fragment>
          }
        >
          <PerformerGalleriesPanel performer={performer} />
        </Tab>
        <Tab
          eventKey="images"
          title={
            <React.Fragment>
              {intl.formatMessage({ id: "images" })}
              <Badge className="left-spacing" pill variant="secondary">
                {intl.formatNumber(performer.image_count ?? 0)}
              </Badge>
            </React.Fragment>
          }
        >
          <PerformerImagesPanel performer={performer} />
        </Tab>
        <Tab
          eventKey="movies"
          title={
            <React.Fragment>
              {intl.formatMessage({ id: "movies" })}
              <Badge className="left-spacing" pill variant="secondary">
                {intl.formatNumber(performer.movie_count ?? 0)}
              </Badge>
            </React.Fragment>
          }
        >
          <PerformerMoviesPanel performer={performer} />
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
          isNew={false}
          onImageChange={onImageChange}
          onImageEncoding={onImageEncoding}
          onCancelEditing={() => {
            setIsEditing(false);
          }}
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
    if (performer?.aliases) {
      return (
        <div>
          <span className="alias-head">
            <FormattedMessage id="also_known_as" />{" "}
          </span>
          <span className="alias">{performer.aliases}</span>
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
            rating: v,
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
        <Icon icon="heart" />
      </Button>
      {performer.url && (
        <Button className="minimal icon-link">
          <a
            href={TextUtils.sanitiseURL(performer.url)}
            className="link"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Icon icon="link" />
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
            <Icon icon="dove" />
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
            <Icon icon="camera" />
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

  return (
    <div id="performer-page" className="row">
      <Helmet>
        <title>{performer.name}</title>
      </Helmet>

      <div className="performer-image-container col-md-4 text-center">
        {imageEncoding ? (
          <LoadingIndicator message="Encoding image..." />
        ) : (
          <Button variant="link" onClick={() => showLightbox()}>
            <img className="performer" src={activeImage} alt="Performer" />
          </Button>
        )}
      </div>
      <div className="col-md-8">
        <div className="row">
          <div className="performer-head col">
            <h2>
              <GenderIcon
                gender={performer.gender}
                className="gender-icon mr-2 flag-icon"
              />
              <CountryFlag country={performer.country} className="mr-2" />
              {performer.name}
              {renderClickableIcons()}
            </h2>
            <RatingStars
              value={performer.rating ?? undefined}
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
