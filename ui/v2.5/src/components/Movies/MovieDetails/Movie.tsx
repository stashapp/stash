import React, { useEffect, useMemo, useState } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import cx from "classnames";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  useFindMovie,
  useMovieUpdate,
  useMovieDestroy,
} from "src/core/StashService";
import { useHistory, RouteComponentProps } from "react-router-dom";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { ModalComponent } from "src/components/Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { MovieScenesPanel } from "./MovieScenesPanel";
import {
  CompressedMovieDetailsPanel,
  MovieDetailsPanel,
} from "./MovieDetailsPanel";
import { MovieEditPanel } from "./MovieEditPanel";
import {
  faChevronDown,
  faChevronUp,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { ConfigurationContext } from "src/hooks/Config";
import { DetailImage } from "src/components/Shared/DetailImage";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { useLoadStickyHeader } from "src/hooks/detailsPanel";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { ExternalLinksButton } from "src/components/Shared/ExternalLinksButton";

interface IProps {
  movie: GQL.MovieDataFragment;
}

interface IMovieParams {
  id: string;
}

const MoviePage: React.FC<IProps> = ({ movie }) => {
  const intl = useIntl();
  const history = useHistory();
  const Toast = useToast();

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui;
  const enableBackgroundImage = uiConfig?.enableMovieBackgroundImage ?? false;
  const compactExpandedDetails = uiConfig?.compactExpandedDetails ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? true;

  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);
  const loadStickyHeader = useLoadStickyHeader();

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing movie state
  const [frontImage, setFrontImage] = useState<string | null>();
  const [backImage, setBackImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const defaultImage =
    movie.front_image_path && movie.front_image_path.includes("default=true")
      ? true
      : false;

  const lightboxImages = useMemo(() => {
    const covers = [
      ...(movie.front_image_path && !defaultImage
        ? [
            {
              paths: {
                thumbnail: movie.front_image_path,
                image: movie.front_image_path,
              },
            },
          ]
        : []),
      ...(movie.back_image_path
        ? [
            {
              paths: {
                thumbnail: movie.back_image_path,
                image: movie.back_image_path,
              },
            },
          ]
        : []),
    ];
    return covers;
  }, [movie.front_image_path, movie.back_image_path, defaultImage]);

  const index = lightboxImages.length;

  const showLightbox = useLightbox({
    images: lightboxImages,
  });

  const [updateMovie, { loading: updating }] = useMovieUpdate();
  const [deleteMovie, { loading: deleting }] = useMovieDestroy({
    id: movie.id,
  });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => toggleEditing());
    Mousetrap.bind("d d", () => {
      setIsDeleteAlertOpen(true);
    });
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  useRatingKeybinds(
    true,
    configuration?.ui.ratingSystemOptions?.type,
    setRating
  );

  async function onSave(input: GQL.MovieCreateInput) {
    await updateMovie({
      variables: {
        input: {
          id: movie.id,
          ...input,
        },
      },
    });
    toggleEditing(false);
    Toast.success(
      intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "movie" }).toLocaleLowerCase() }
      )
    );
  }

  async function onDelete() {
    try {
      await deleteMovie();
    } catch (e) {
      Toast.error(e);
    }

    // redirect to movies page
    history.push(`/movies`);
  }

  function toggleEditing(value?: boolean) {
    if (value !== undefined) {
      setIsEditing(value);
    } else {
      setIsEditing((e) => !e);
    }
    setFrontImage(undefined);
    setBackImage(undefined);
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
                movie.name ??
                intl.formatMessage({ id: "movie" }).toLocaleLowerCase(),
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

  function renderFrontImage() {
    let image = movie.front_image_path;
    if (isEditing) {
      if (frontImage === null && image) {
        const imageURL = new URL(image);
        imageURL.searchParams.set("default", "true");
        image = imageURL.toString();
      } else if (frontImage) {
        image = frontImage;
      }
    }

    if (image && defaultImage) {
      return (
        <div className="movie-image-container">
          <DetailImage alt="Front Cover" src={image} />
        </div>
      );
    } else if (image) {
      return (
        <Button
          className="movie-image-container"
          variant="link"
          onClick={() => showLightbox()}
        >
          <DetailImage alt="Front Cover" src={image} />
        </Button>
      );
    }
  }

  function renderBackImage() {
    let image = movie.back_image_path;
    if (isEditing) {
      if (backImage === null) {
        image = undefined;
      } else if (backImage) {
        image = backImage;
      }
    }

    if (image) {
      return (
        <Button
          className="movie-image-container"
          variant="link"
          onClick={() => showLightbox(index - 1)}
        >
          <DetailImage alt="Back Cover" src={image} />
        </Button>
      );
    }
  }

  const renderClickableIcons = () => (
    <span className="name-icons">
      {movie.urls.length > 0 && <ExternalLinksButton urls={movie.urls} />}
    </span>
  );

  function maybeRenderAliases() {
    if (movie?.aliases) {
      return (
        <div>
          <span className="alias-head">{movie?.aliases}</span>
        </div>
      );
    }
  }

  function setRating(v: number | null) {
    if (movie.id) {
      updateMovie({
        variables: {
          input: {
            id: movie.id,
            rating100: v,
          },
        },
      });
    }
  }

  const renderTabs = () => <MovieScenesPanel active={true} movie={movie} />;

  function maybeRenderDetails() {
    if (!isEditing) {
      return (
        <MovieDetailsPanel
          movie={movie}
          collapsed={collapsed}
          fullWidth={!collapsed && !compactExpandedDetails}
        />
      );
    }
  }

  function maybeRenderEditPanel() {
    if (isEditing) {
      return (
        <MovieEditPanel
          movie={movie}
          onSubmit={onSave}
          onCancel={() => toggleEditing()}
          onDelete={onDelete}
          setFrontImage={setFrontImage}
          setBackImage={setBackImage}
          setEncodingImage={setEncodingImage}
        />
      );
    }
    {
      return (
        <DetailsEditNavbar
          objectName={movie.name}
          isNew={false}
          isEditing={isEditing}
          onToggleEdit={() => toggleEditing()}
          onSave={() => {}}
          onImageChange={() => {}}
          onDelete={onDelete}
        />
      );
    }
  }

  function maybeRenderCompressedDetails() {
    if (!isEditing && loadStickyHeader) {
      return <CompressedMovieDetailsPanel movie={movie} />;
    }
  }

  function maybeRenderHeaderBackgroundImage() {
    let image = movie.front_image_path;
    if (enableBackgroundImage && !isEditing && image) {
      const imageURL = new URL(image);
      let isDefaultImage = imageURL.searchParams.get("default");
      if (!isDefaultImage) {
        return (
          <div className="background-image-container">
            <picture>
              <source src={image} />
              <img
                className="background-image"
                src={image}
                alt={`${movie.name} background`}
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

  if (updating || deleting) return <LoadingIndicator />;

  const headerClassName = cx("detail-header", {
    edit: isEditing,
    collapsed,
    "full-width": !collapsed && !compactExpandedDetails,
  });

  return (
    <div id="movie-page" className="row">
      <Helmet>
        <title>{movie?.name}</title>
      </Helmet>

      <div className={headerClassName}>
        {maybeRenderHeaderBackgroundImage()}
        <div className="detail-container">
          <div className="detail-header-image">
            <div className="logo w-100">
              {encodingImage ? (
                <LoadingIndicator
                  message={intl.formatMessage({ id: "actions.encoding_image" })}
                />
              ) : (
                <div className="movie-images">
                  {renderFrontImage()}
                  {renderBackImage()}
                </div>
              )}
            </div>
          </div>
          <div className="row">
            <div className="movie-head col">
              <h2>
                <span className="movie-name">{movie.name}</span>
                {maybeRenderShowCollapseButton()}
                {renderClickableIcons()}
              </h2>
              {maybeRenderAliases()}
              <RatingSystem
                value={movie.rating100}
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
        <div className="movie-body">
          <div className="movie-tabs">{maybeRenderTab()}</div>
        </div>
      </div>
      {renderDeleteAlert()}
    </div>
  );
};

const MovieLoader: React.FC<RouteComponentProps<IMovieParams>> = ({
  match,
}) => {
  const { id } = match.params;
  const { data, loading, error } = useFindMovie(id);

  useScrollToTopOnMount();

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findMovie)
    return <ErrorMessage error={`No movie found with id ${id}.`} />;

  return <MoviePage movie={data.findMovie} />;
};

export default MovieLoader;
