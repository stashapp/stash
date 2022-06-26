import React, { useEffect, useState } from "react";
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
import {
  DetailsEditNavbar,
  ErrorMessage,
  LoadingIndicator,
  Modal,
} from "src/components/Shared";
import { useToast } from "src/hooks";
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
  const [frontImage, setFrontImage] = useState<string | undefined | null>(
    undefined
  );
  const [backImage, setBackImage] = useState<string | undefined | null>(
    undefined
  );
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [updateMovie, { loading: updating }] = useMovieUpdate();
  const [deleteMovie, { loading: deleting }] = useMovieDestroy({
    id: movie.id,
  });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => onDelete());

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  const onImageEncoding = (isEncoding = false) => setEncodingImage(isEncoding);

  function getMovieInput(
    input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>
  ) {
    const ret: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      ...input,
      id: movie.id,
    };

    return ret;
  }

  async function onSave(
    input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>
  ) {
    try {
      const result = await updateMovie({
        variables: {
          input: getMovieInput(input) as GQL.MovieUpdateInput,
        },
      });
      if (result.data?.movieUpdate) {
        setIsEditing(false);
        history.push(`/movies/${result.data.movieUpdate.id}`);
      }
    } catch (e) {
      Toast.error(e);
    }
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

  function onToggleEdit() {
    setIsEditing(!isEditing);
    setFrontImage(undefined);
    setBackImage(undefined);
  }

  function renderDeleteAlert() {
    return (
      <Modal
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
      </Modal>
    );
  }

  function renderFrontImage() {
    let image = movie.front_image_path;
    if (isEditing) {
      if (frontImage === null) {
        image = `${image}&default=true`;
      } else if (frontImage) {
        image = frontImage;
      }
    }

    if (image) {
      return (
        <div className="movie-image-container">
          <img alt="Front Cover" src={image} />
        </div>
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
        <div className="movie-image-container">
          <img alt="Back Cover" src={image} />
        </div>
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
              onToggleEdit={onToggleEdit}
              onSave={() => {}}
              onImageChange={() => {}}
              onDelete={onDelete}
            />
          </>
        ) : (
          <MovieEditPanel
            movie={movie}
            onSubmit={onSave}
            onCancel={onToggleEdit}
            onDelete={onDelete}
            setFrontImage={setFrontImage}
            setBackImage={setBackImage}
            onImageEncoding={onImageEncoding}
          />
        )}
      </div>

      <div className="col-xl-8 col-lg-6">
        <MovieScenesPanel movie={movie} />
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
