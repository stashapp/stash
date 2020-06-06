/* eslint-disable react/no-this-in-sfc */
import React, { useEffect, useState, useCallback } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  useFindMovie,
  useMovieUpdate,
  useMovieCreate,
  useMovieDestroy,
} from "src/core/StashService";
import { useParams, useHistory } from "react-router-dom";
import cx from "classnames";
import {
  DetailsEditNavbar,
  LoadingIndicator,
  Modal,
  StudioSelect,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Table, Form } from "react-bootstrap";
import {
  TableUtils,
  ImageUtils,
  EditableTextUtils,
  TextUtils,
} from "src/utils";
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
  const [frontImage, setFrontImage] = useState<string | undefined>(undefined);
  const [backImage, setBackImage] = useState<string | undefined>(undefined);
  const [name, setName] = useState<string | undefined>(undefined);
  const [aliases, setAliases] = useState<string | undefined>(undefined);
  const [duration, setDuration] = useState<number | undefined>(undefined);
  const [date, setDate] = useState<string | undefined>(undefined);
  const [rating, setRating] = useState<number | undefined>(undefined);
  const [studioId, setStudioId] = useState<string>();
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
  const { data, error, loading } = useFindMovie(id);
  const [updateMovie] = useMovieUpdate(getMovieInput() as GQL.MovieUpdateInput);
  const [createMovie] = useMovieCreate(getMovieInput() as GQL.MovieCreateInput);
  const [deleteMovie] = useMovieDestroy(
    getMovieInput() as GQL.MovieDestroyInput
  );

  const intl = useIntl();

  function updateMovieEditState(state: Partial<GQL.MovieDataFragment>) {
    setName(state.name ?? undefined);
    setAliases(state.aliases ?? undefined);
    setDuration(state.duration ?? undefined);
    setDate(state.date ?? undefined);
    setRating(state.rating ?? undefined);
    setStudioId(state?.studio?.id ?? undefined);
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

  function onImageLoad(imageData: string) {
    setImagePreview(imageData);
    setFrontImage(imageData);
  }

  function onBackImageLoad(imageData: string) {
    setBackImagePreview(imageData);
    setBackImage(imageData);
  }

  const encodingFrontImage = ImageUtils.usePasteImage(onImageLoad, isEditing);
  const encodingBackImage = ImageUtils.usePasteImage(
    onBackImageLoad,
    isEditing
  );

  if (!isNew && !isEditing) {
    if (!data || !data.findMovie || loading) return <LoadingIndicator />;
    if (error) {
      return <>{error!.message}</>;
    }
  }

  function getMovieInput() {
    const input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      name,
      aliases,
      duration,
      date,
      rating,
      studio_id: studioId,
      director,
      synopsis,
      url,
      front_image: frontImage,
      back_image: backImage,
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

  function onToggleEdit() {
    setIsEditing(!isEditing);
    updateMovieData(movie);
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{ text: "Delete", variant: "danger", onClick: onDelete }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>Are you sure you want to delete {name ?? "movie"}?</p>
      </Modal>
    );
  }

  // TODO: CSS class
  return (
    <div className="row">
      <div
        className={cx("movie-details", "col", {
          "col ml-sm-5": !isNew,
        })}
      >
        {isNew && <h2>Add Movie</h2>}
        <div className="logo w-100">
          {encodingFrontImage || encodingBackImage ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            <>
              <img alt={name} className="logo w-50" src={imagePreview} />
              <img alt={name} className="logo w-50" src={backimagePreview} />
            </>
          )}
        </div>

        <Table>
          <tbody>
            {TableUtils.renderInputGroup({
              title: "Name",
              value: name ?? "",
              isEditing: !!isEditing,
              onChange: setName,
            })}
            {TableUtils.renderInputGroup({
              title: "Aliases",
              value: aliases,
              isEditing,
              onChange: setAliases,
            })}
            {TableUtils.renderDurationInput({
              title: "Duration",
              value: duration ? duration.toString() : "",
              isEditing,
              onChange: (value: string | undefined) =>
                setDuration(value ? Number.parseInt(value, 10) : undefined),
            })}
            {TableUtils.renderInputGroup({
              title: `Date ${isEditing ? "(YYYY-MM-DD)" : ""}`,
              value: isEditing ? date : TextUtils.formatDate(intl, date),
              isEditing,
              onChange: setDate,
            })}
            <tr>
              <td>Studio</td>
              <td>
                <StudioSelect
                  isDisabled={!isEditing}
                  onSelect={(items) =>
                    setStudioId(items.length > 0 ? items[0]?.id : undefined)
                  }
                  ids={studioId ? [studioId] : []}
                />
              </td>
            </tr>
            {TableUtils.renderInputGroup({
              title: "Director",
              value: director,
              isEditing,
              onChange: setDirector,
            })}
            {TableUtils.renderHtmlSelect({
              title: "Rating",
              value: rating ?? "",
              isEditing,
              onChange: (value: string) =>
                setRating(Number.parseInt(value, 10)),
              selectOptions: ["", "1", "2", "3", "4", "5"],
            })}
          </tbody>
        </Table>

        <Form.Group controlId="url">
          <Form.Label>URL</Form.Label>
          <div>
            {EditableTextUtils.renderInputGroup({
              isEditing,
              onChange: setUrl,
              value: url,
              url: TextUtils.sanitiseURL(url),
            })}
          </div>
        </Form.Group>

        <Form.Group controlId="synopsis">
          <Form.Label>Synopsis</Form.Label>
          <Form.Control
            as="textarea"
            readOnly={!isEditing}
            className="movie-synopsis text-input"
            onChange={(newValue: React.ChangeEvent<HTMLTextAreaElement>) =>
              setSynopsis(newValue.currentTarget.value)
            }
            value={synopsis}
          />
        </Form.Group>

        <DetailsEditNavbar
          objectName={name ?? "movie"}
          isNew={isNew}
          isEditing={isEditing}
          onToggleEdit={onToggleEdit}
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
