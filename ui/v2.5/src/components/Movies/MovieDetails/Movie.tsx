import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import cx from "classnames";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  useFindMovie,
  useMovieUpdate,
  useMovieCreate,
  useMovieDestroy,
} from "src/core/StashService";
import { useParams, useHistory } from "react-router-dom";
import {
  DetailsEditNavbar,
  LoadingIndicator,
  Modal,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { MovieScenesPanel } from "./MovieScenesPanel";
import { MovieDetailsPanel } from "./MovieDetailsPanel";
import { MovieEditPanel } from "./MovieEditPanel";

interface IMovieParams {
  id?: string;
}

export const Movie: React.FC = () => {
  const intl = useIntl();
  const history = useHistory();
  const Toast = useToast();
  const { id = "new" } = useParams<IMovieParams>();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing movie state
  const [frontImage, setFrontImage] = useState<string | undefined | null>(
    undefined
  );
  const [backImage, setBackImage] = useState<string | undefined | null>(
    undefined
  );
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  // Network state
  const { data, error, loading } = useFindMovie(id);
  const movie = data?.findMovie;

  const [isLoading, setIsLoading] = useState(false);
  const [updateMovie] = useMovieUpdate();
  const [createMovie] = useMovieCreate();
  const [deleteMovie] = useMovieDestroy({ id });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => onDelete());

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
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");

      Mousetrap.unbind("r");
    };
  });

  const onImageEncoding = (isEncoding = false) => setEncodingImage(isEncoding);

  if (!isNew && !isEditing) {
    if (!data || !data.findMovie || loading) return <LoadingIndicator />;
    if (error) {
      return <>{error!.message}</>;
    }
  }

  function getMovieInput(
    input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>
  ) {
    const ret: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      ...input,
    };

    if (!isNew) {
      (ret as GQL.MovieUpdateInput).id = id;
    }
    return ret;
  }

  async function onSave(
    input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>
  ) {
    try {
      setIsLoading(true);

      if (!isNew) {
        const result = await updateMovie({
          variables: {
            input: getMovieInput(input) as GQL.MovieUpdateInput,
          },
        });
        if (result.data?.movieUpdate) {
          setIsEditing(false);
          history.push(`/movies/${result.data.movieUpdate.id}`);
        }
      } else {
        const result = await createMovie({
          variables: getMovieInput(input) as GQL.MovieCreateInput,
        });
        if (result.data?.movieCreate?.id) {
          history.push(`/movies/${result.data.movieCreate.id}`);
          setIsEditing(false);
        }
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onDelete() {
    try {
      setIsLoading(true);
      await deleteMovie();
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }

    // redirect to movies page
    history.push(`/movies`);
  }

  function onToggleEdit() {
    setIsEditing(!isEditing);
    setFrontImage(undefined);
    setBackImage(undefined);
  }

  function setRating(v: number | null) {
    if (movie?.id) {
      updateMovie({
        variables: {
          input: {
            id: movie?.id,
            rating: v,
          },
        },
      });
    }
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
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
                movie?.name ??
                intl.formatMessage({ id: "movie" }).toLocaleLowerCase(),
            }}
          />
        </p>
      </Modal>
    );
  }

  function renderFrontImage() {
    let image = movie?.front_image_path;
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
    let image = movie?.back_image_path;
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

  if (isLoading) return <LoadingIndicator />;

  // TODO: CSS class
  return (
    <div className="row">
      <div
        className={cx("movie-details mb-3 col", {
          "col-xl-4 col-lg-6": !isNew,
        })}
      >
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

        <RatingStars
          value={movie?.rating ?? undefined}
          onSetRating={(value) => setRating(value ?? null)}
        />

        {!isEditing && movie ? (
          <>
            <MovieDetailsPanel movie={movie} />
            {/* HACK - this is also rendered in the MovieEditPanel */}
            <DetailsEditNavbar
              objectName={movie?.name ?? "movie"}
              isNew={isNew}
              isEditing={isEditing}
              onToggleEdit={onToggleEdit}
              onSave={() => {}}
              onImageChange={() => {}}
              onDelete={onDelete}
            />
          </>
        ) : (
          <MovieEditPanel
            movie={movie ?? undefined}
            onSubmit={onSave}
            onCancel={onToggleEdit}
            onDelete={onDelete}
            setFrontImage={setFrontImage}
            setBackImage={setBackImage}
            onImageEncoding={onImageEncoding}
          />
        )}
      </div>

      {!isNew && movie && (
        <div className="col-xl-8 col-lg-6">
          <MovieScenesPanel movie={movie} />
        </div>
      )}
      {renderDeleteAlert()}
    </div>
  );
};
