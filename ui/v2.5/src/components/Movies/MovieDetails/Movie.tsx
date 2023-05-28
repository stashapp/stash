import React, { useEffect, useMemo, useState } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  useFindMovie,
  useMovieUpdate,
  useMovieDestroy,
} from "src/core/StashService";
import { useParams, useHistory } from "react-router-dom";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useLightbox } from "src/hooks/Lightbox/hooks";
import { ModalComponent } from "src/components/Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { MovieScenesPanel } from "./MovieScenesPanel";
import { MovieDetailsPanel } from "./MovieDetailsPanel";
import { MovieEditPanel } from "./MovieEditPanel";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";

interface IProps {
  movie: GQL.MovieDataFragment;
}

const MoviePage: React.FC<IProps> = ({ movie }) => {
  const intl = useIntl();
  const history = useHistory();
  const Toast = useToast();

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
      onDelete();
    });

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

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
    Toast.success({
      content: intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "movie" }).toLocaleLowerCase() }
      ),
    });
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
          <img alt="Front Cover" src={image} />
        </div>
      );
    } else if (image) {
      return (
        <Button
          className="movie-image-container"
          variant="link"
          onClick={() => showLightbox()}
        >
          <img alt="Front Cover" src={image} />
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
          <img alt="Back Cover" src={image} />
        </Button>
      );
    }
  }

  if (updating || deleting) return <LoadingIndicator />;

  // TODO: CSS class
  return (
    <div className="row">
      <Helmet>
        <title>{movie?.name}</title>
      </Helmet>

      <div className="movie-details mb-3 col col-xl-4 col-lg-6">
        <div className="logo w-100">
          {encodingImage ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            <div className="movie-images">
              {renderFrontImage()}
              {renderBackImage()}
            </div>
          )}
        </div>

        {!isEditing ? (
          <>
            <MovieDetailsPanel movie={movie} />
            {/* HACK - this is also rendered in the MovieEditPanel */}
            <DetailsEditNavbar
              objectName={movie.name}
              isNew={false}
              isEditing={isEditing}
              onToggleEdit={() => toggleEditing()}
              onSave={() => {}}
              onImageChange={() => {}}
              onDelete={onDelete}
            />
          </>
        ) : (
          <MovieEditPanel
            movie={movie}
            onSubmit={onSave}
            onCancel={() => toggleEditing()}
            onDelete={onDelete}
            setFrontImage={setFrontImage}
            setBackImage={setBackImage}
            setEncodingImage={setEncodingImage}
          />
        )}
      </div>

      <div className="col-xl-8 col-lg-6">
        <MovieScenesPanel active={true} movie={movie} />
      </div>
      {renderDeleteAlert()}
    </div>
  );
};

const MovieLoader: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const { data, loading, error } = useFindMovie(id ?? "");

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findMovie)
    return <ErrorMessage error={`No movie found with id ${id}.`} />;

  return <MoviePage movie={data.findMovie} />;
};

export default MovieLoader;
