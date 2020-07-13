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
import {
  DetailsEditNavbar,
  LoadingIndicator,
  Modal,
  StudioSelect,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Table, Form, Modal as BSModal, Button } from "react-bootstrap";
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
  const [isImageAlertOpen, setIsImageAlertOpen] = useState<boolean>(false);

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

  const [imageClipboard, setImageClipboard] = useState<string | undefined>(
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

  // set up hotkeys
  useEffect(() => {
    if (isEditing) {
      Mousetrap.bind("r 0", () => setRating(NaN));
      Mousetrap.bind("r 1", () => setRating(1));
      Mousetrap.bind("r 2", () => setRating(2));
      Mousetrap.bind("r 3", () => setRating(3));
      Mousetrap.bind("r 4", () => setRating(4));
      Mousetrap.bind("r 5", () => setRating(5));
      // Mousetrap.bind("u", (e) => {
      //   setStudioFocus()
      //   e.preventDefault();
      // });
      Mousetrap.bind("s s", () => onSave());
    }

    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => onDelete());

    return () => {
      if (isEditing) {
        Mousetrap.unbind("r 0");
        Mousetrap.unbind("r 1");
        Mousetrap.unbind("r 2");
        Mousetrap.unbind("r 3");
        Mousetrap.unbind("r 4");
        Mousetrap.unbind("r 5");
        // Mousetrap.unbind("u");
        Mousetrap.unbind("s s");
      }

      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

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
    setImageClipboard(imageData);
    setIsImageAlertOpen(true);
  }

  function setImageFromClipboard(isFrontImage: boolean) {
    if (isFrontImage) {
      setImagePreview(imageClipboard);
      setFrontImage(imageClipboard);
    } else {
      setBackImagePreview(imageClipboard);
      setBackImage(imageClipboard);
    }

    setImageClipboard(undefined);
    setIsImageAlertOpen(false);
  }

  function onBackImageLoad(imageData: string) {
    setBackImagePreview(imageData);
    setBackImage(imageData);
  }

  const encodingImage = ImageUtils.usePasteImage(onImageLoad, isEditing);

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

  function renderImageAlert() {
    return (
      <BSModal
        show={isImageAlertOpen}
        onHide={() => setIsImageAlertOpen(false)}
      >
        <BSModal.Body>
          <p>Select image to set</p>
        </BSModal.Body>
        <BSModal.Footer>
          <div>
            <Button
              className="mr-2"
              variant="secondary"
              onClick={() => setIsImageAlertOpen(false)}
            >
              Cancel
            </Button>

            <Button
              className="mr-2"
              onClick={() => setImageFromClipboard(false)}
            >
              Back Image
            </Button>
            <Button
              className="mr-2"
              onClick={() => setImageFromClipboard(true)}
            >
              Front Image
            </Button>
          </div>
        </BSModal.Footer>
      </BSModal>
    );
  }

  // TODO: CSS class
  return (
    <div className="row">
      <div className="movie-details col">
        {isNew && <h2>Add Movie</h2>}
        <div className="logo w-100">
          {encodingImage ? (
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
        <div className="col-lg-8 col-md-7">
          <MovieScenesPanel movie={movie} />
        </div>
      )}
      {renderDeleteAlert()}
      {renderImageAlert()}
    </div>
  );
};
