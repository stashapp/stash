import React, { useEffect, useState, useCallback } from "react";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { useParams, useHistory } from "react-router-dom";
import cx from "classnames";
import {
  DetailsEditNavbar,
  LoadingIndicator,
  Modal
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Table } from "react-bootstrap";
import { TableUtils, ImageUtils } from "src/utils";
import { MovieScenesPanel } from "./MovieScenesPanel";

export const Movie: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const { id = "new" } = useParams();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing movie state
  const [front_image, setFrontImage] = useState<string | undefined>(undefined);
  const [back_image, setBackImage] = useState<string | undefined>(undefined);
  const [name, setName] = useState<string | undefined>(undefined);
  const [aliases, setAliases] = useState<string | undefined>(undefined);
  const [duration, setDuration] = useState<string | undefined>(undefined);
  const [date, setDate] = useState<string | undefined>(undefined);
  const [rating, setRating] = useState<string | undefined>(undefined);
  const [director, setDirector] = useState<string | undefined>(undefined);
  const [synopsis, setSynopsis] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);

  // Movie state
  const [movie, setMovie] = useState<Partial<GQL.MovieDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(
    undefined
  );
  const [backimagePreview, setBackImagePreview] = useState<string | undefined>(
    undefined
  );

  // Network state
  const { data, error, loading } = StashService.useFindMovie(id);
  const [updateMovie] = StashService.useMovieUpdate(
    getMovieInput() as GQL.MovieUpdateInput
  );
  const [createMovie] = StashService.useMovieCreate(
    getMovieInput() as GQL.MovieCreateInput
  );
  const [deleteMovie] = StashService.useMovieDestroy(
    getMovieInput() as GQL.MovieDestroyInput
  );

  function updateMovieEditState(state: Partial<GQL.MovieDataFragment>) {
    setName(state.name ?? undefined);
    setAliases(state.aliases ?? undefined);
    setDuration(state.duration ?? undefined);
    setDate(state.date ?? undefined);
    setRating(state.rating ?? undefined);
    setDirector(state.director ?? undefined);
    setSynopsis(state.synopsis ?? undefined);
    setUrl(state.url ?? undefined);
  }

  const updateMovieData = useCallback(
    (movieData: Partial<GQL.MovieDataFragment>) => {
      setFrontImage(undefined);
      setBackImage(undefined);
      updateMovieEditState(movieData);
      setImagePreview(movieData.front_image_path ?? undefined);
      setBackImagePreview(movieData.back_image_path ?? undefined);
      setMovie(movieData);
    },
    []
  );

  useEffect(() => {
    if (data && data.findMovie) {
      updateMovieData(data.findMovie);
    }
  }, [data, updateMovieData]);

  function onImageLoad(this: FileReader) {
    setImagePreview(this.result as string);
    setFrontImage(this.result as string);
  }

  function onBackImageLoad(this: FileReader) {
    setBackImagePreview(this.result as string);
    setBackImage(this.result as string);
  }

  ImageUtils.usePasteImage(onImageLoad);
  ImageUtils.usePasteImage(onBackImageLoad);

  if (!isNew && !isEditing) {
    if (!data || !data.findMovie || loading) return <LoadingIndicator />;
    if (!!error) {
      return <>{error.message}</>;
    }
  }

  function getMovieInput() {
    const input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      name,
      aliases,
      duration,
      date,
      rating,
      director,
      synopsis,
      url,
      front_image,
      back_image
    };

    if (!isNew) {
      (input as GQL.MovieUpdateInput).id = id;
    }
    return input;
  }

  async function onSave() {
    try {
      if (!isNew) {
        const result = await updateMovie();
        if (result.data?.movieUpdate) {
          updateMovieData(result.data.movieUpdate);
          setIsEditing(false);
        }
      } else {
        const result = await createMovie();
        if (result.data?.movieCreate?.id) {
          history.push(`/movies/${result.data.movieCreate.id}`);
          setIsEditing(false);
        }
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

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onBackImageLoad);
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{ text: "Delete", variant: "danger", onClick: onDelete }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>Are you sure you want to delete {movie.name ?? "movie"}?</p>
      </Modal>
    );
  }

  // TODO: CSS class
  return (
    <div className="row">
      <div
        className={cx("movie-details", {
          "col ml-sm-5": !isNew,
          "col-8": isNew
        })}
      >
        {isNew && <h2>Add Movie</h2>}
        <div className="logo w-100">
          <img alt={name} className="logo w-50" src={imagePreview} />
          <img alt={name} className="logo w-50" src={backimagePreview} />
        </div>

        <Table>
          <tbody>
            {TableUtils.renderInputGroup({
              title: "Name",
              value: movie.name ?? "",
              isEditing: !!isEditing,
              onChange: setName
            })}
            {TableUtils.renderInputGroup({
              title: "Aliases",
              value: aliases,
              isEditing,
              onChange: setAliases
            })}
            {TableUtils.renderInputGroup({
              title: "Duration",
              value: duration,
              isEditing,
              onChange: setDuration
            })}
            {TableUtils.renderInputGroup({
              title: "Date (YYYY-MM-DD)",
              value: date,
              isEditing,
              onChange: setDate
            })}
            {TableUtils.renderInputGroup({
              title: "Director",
              value: director,
              isEditing,
              onChange: setDirector
            })}
            {TableUtils.renderHtmlSelect({
              title: "Rating",
              value: rating,
              isEditing,
              onChange: (value: string) => setRating(value),
              selectOptions: ["", "1", "2", "3", "4", "5"]
            })}
            {TableUtils.renderInputGroup({
              title: "URL",
              value: url,
              isEditing,
              onChange: setUrl
            })}
            {TableUtils.renderTextArea({
              title: "Synopsis",
              value: synopsis,
              isEditing,
              onChange: setSynopsis
            })}
          </tbody>
        </Table>
        <DetailsEditNavbar
          objectName={movie.name ?? "movie"}
          isNew={isNew}
          isEditing={isEditing}
          onToggleEdit={() => setIsEditing(!isEditing)}
          onSave={onSave}
          onImageChange={onImageChange}
          onBackImageChange={onBackImageChange}
          onDelete={onDelete}
        />
      </div>
      {!isNew && (
        <div className="col-12 col-sm-8">
          <MovieScenesPanel movie={movie} />
        </div>
      )}
      {renderDeleteAlert()}
    </div>
  );
};
