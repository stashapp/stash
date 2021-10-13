import React, { useEffect, useMemo, useState } from "react";
import { Button, Tabs, Tab } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { useParams, useHistory } from "react-router-dom";
import cx from "classnames";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  useFindPerformer,
  usePerformerUpdate,
  usePerformerDestroy,
} from "src/core/StashService";
import {
  CountryFlag,
  ErrorMessage,
  Icon,
  LoadingIndicator,
} from "src/components/Shared";
import { useLightbox, useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { PerformerDetailsPanel } from "./PerformerDetailsPanel";
import { PerformerOperationsPanel } from "./PerformerOperationsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";
import { PerformerGalleriesPanel } from "./PerformerGalleriesPanel";
import { PerformerMoviesPanel } from "./PerformerMoviesPanel";
import { PerformerImagesPanel } from "./PerformerImagesPanel";
import { PerformerEditPanel } from "./PerformerEditPanel";

interface IPerformerParams {
  id?: string;
  tab?: string;
}

export const Performer: React.FC = () => {
  const Toast = useToast();
  const history = useHistory();
  const intl = useIntl();
  const { tab = "details", id = "new" } = useParams<IPerformerParams>();
  const isNew = id === "new";

  // Performer state
  const [imagePreview, setImagePreview] = useState<string | null>();
  const [imageEncoding, setImageEncoding] = useState<boolean>(false);
  const { data, loading: performerLoading, error } = useFindPerformer(id);
  const performer = data?.findPerformer || ({} as Partial<GQL.Performer>);

  // if undefined then get the existing image
  // if null then get the default (no) image
  // otherwise get the set image
  const activeImage =
    imagePreview === undefined
      ? performer.image_path ?? ""
      : imagePreview ?? (isNew ? "" : `${performer.image_path}&default=true`);
  const lightboxImages = useMemo(
    () => [{ paths: { thumbnail: activeImage, image: activeImage } }],
    [activeImage]
  );

  const showLightbox = useLightbox({
    images: lightboxImages,
  });

  // Network state
  const [loading, setIsLoading] = useState(false);
  const isLoading = performerLoading || loading;

  const [updatePerformer] = usePerformerUpdate();
  const [deletePerformer] = usePerformerDestroy();

  const activeTabKey =
    tab === "scenes" ||
    tab === "galleries" ||
    tab === "images" ||
    tab === "movies" ||
    tab === "edit" ||
    tab === "operations"
      ? tab
      : "details";
  const setActiveTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "details" ? "" : `/${newTab}`;
      history.replace(`/performers/${id}${tabParam}`);
    }
  };

  const onImageChange = (image?: string | null) => setImagePreview(image);

  const onImageEncoding = (isEncoding = false) => setImageEncoding(isEncoding);

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("details"));
    Mousetrap.bind("e", () => setActiveTabKey("edit"));
    Mousetrap.bind("c", () => setActiveTabKey("scenes"));
    Mousetrap.bind("g", () => setActiveTabKey("galleries"));
    Mousetrap.bind("m", () => setActiveTabKey("movies"));
    Mousetrap.bind("o", () => setActiveTabKey("operations"));
    Mousetrap.bind("f", () => setFavorite(!performer.favorite));

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("e");
      Mousetrap.unbind("c");
      Mousetrap.unbind("f");
      Mousetrap.unbind("o");
    };
  });

  if (isLoading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!performer.id && !isNew)
    return <ErrorMessage error={`No performer found with id ${id}.`} />;

  async function onDelete() {
    setIsLoading(true);
    try {
      await deletePerformer({ variables: { id } });
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);

    // redirect to performers page
    history.push("/performers");
  }

  const renderTabs = () => (
    <Tabs
      activeKey={activeTabKey}
      onSelect={setActiveTabKey}
      id="performer-details"
      unmountOnExit
    >
      <Tab eventKey="details" title={intl.formatMessage({ id: "details" })}>
        <PerformerDetailsPanel performer={performer} />
      </Tab>
      <Tab eventKey="scenes" title={intl.formatMessage({ id: "scenes" })}>
        <PerformerScenesPanel performer={performer} />
      </Tab>
      <Tab eventKey="galleries" title={intl.formatMessage({ id: "galleries" })}>
        <PerformerGalleriesPanel performer={performer} />
      </Tab>
      <Tab eventKey="images" title={intl.formatMessage({ id: "images" })}>
        <PerformerImagesPanel performer={performer} />
      </Tab>
      <Tab eventKey="movies" title={intl.formatMessage({ id: "movies" })}>
        <PerformerMoviesPanel performer={performer} />
      </Tab>
      <Tab eventKey="edit" title={intl.formatMessage({ id: "actions.edit" })}>
        <PerformerEditPanel
          performer={performer}
          isVisible={activeTabKey === "edit"}
          isNew={isNew}
          onDelete={onDelete}
          onImageChange={onImageChange}
          onImageEncoding={onImageEncoding}
        />
      </Tab>
      <Tab
        eventKey="operations"
        title={intl.formatMessage({ id: "operations" })}
      >
        <PerformerOperationsPanel performer={performer} />
      </Tab>
    </Tabs>
  );

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

  const renderIcons = () => (
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

  function renderPerformerImage() {
    if (imageEncoding) {
      return <LoadingIndicator message="Encoding image..." />;
    }
    if (activeImage) {
      return <img className="performer" src={activeImage} alt="Performer" />;
    }
  }

  if (isNew)
    return (
      <div className="row new-view" id="performer-page">
        <div className="performer-image-container col-md-4 text-center">
          {renderPerformerImage()}
        </div>
        <div className="col-md-8">
          <h2>Create Performer</h2>
          <PerformerEditPanel
            performer={performer}
            isVisible
            isNew
            onDelete={onDelete}
            onImageChange={onImageChange}
            onImageEncoding={onImageEncoding}
          />
        </div>
      </div>
    );

  if (!performer.id) {
    return <LoadingIndicator />;
  }

  return (
    <div id="performer-page" className="row">
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
              <CountryFlag country={performer.country} className="mr-2" />
              {performer.name}
              {renderIcons()}
            </h2>
            <h4>
              <RatingStars
                value={performer.rating ?? undefined}
                onSetRating={(value) => setRating(value ?? null)}
              />
            </h4>
            {maybeRenderAliases()}
            {maybeRenderAge()}
          </div>
        </div>
        <div className="performer-body">
          <div className="performer-tabs">{renderTabs()}</div>
        </div>
      </div>
    </div>
  );
};
