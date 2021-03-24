import React, { useEffect, useState, useCallback } from "react";
import { useIntl } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  useFindMovie,
  useMovieUpdate,
  useMovieCreate,
  useMovieDestroy,
  queryScrapeMovieURL,
  useListMovieScrapers,
} from "src/core/StashService";
import { useParams, useHistory } from "react-router-dom";
import {
  DetailsEditNavbar,
  LoadingIndicator,
  Modal,
  StudioSelect,
  Icon,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Table, Form, Modal as BSModal, Button } from "react-bootstrap";
import {
  TableUtils,
  ImageUtils,
  EditableTextUtils,
  TextUtils,
  DurationUtils,
} from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { MovieScenesPanel } from "./MovieScenesPanel";
import { MovieScrapeDialog } from "./MovieScrapeDialog";
import { MovieDetailsPanel } from "./MovieDetailsPanel";
import { MovieEditPanel } from "./MovieEditPanel";

interface IMovieParams {
  id?: string;
}

export const Movie: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const { id = "new" } = useParams<IMovieParams>();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [isImageAlertOpen, setIsImageAlertOpen] = useState<boolean>(false);
  const [movieEditState, setMovieEditState] = useState<Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>>();

  // Editing movie state
  const [frontImage, setFrontImage] = useState<string | undefined | null>(
    undefined
  );
  const [backImage, setBackImage] = useState<string | undefined | null>(
    undefined
  );

  // Movie state
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
  const movie = data?.findMovie;

  const [isLoading, setIsLoading] = useState(false);
  const [updateMovie] = useMovieUpdate();
  const [createMovie] = useMovieCreate();
  const [deleteMovie] = useMovieDestroy({id});

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => onDelete());

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  function showImageAlert(imageData: string) {
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

  function onFrontImageLoad(imageData: string) {
    setImagePreview(imageData);
    setFrontImage(imageData);
  }

  const encodingImage = ImageUtils.usePasteImage(showImageAlert, isEditing);

  if (!isNew && !isEditing) {
    if (!data || !data.findMovie || loading) return <LoadingIndicator />;
    if (error) {
      return <>{error!.message}</>;
    }
  }

  function getMovieInput(input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>) {
    const ret: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      ...input,
      front_image: frontImage,
      back_image: backImage,
    };

    if (!isNew) {
      (ret as GQL.MovieUpdateInput).id = id;
    }
    return ret;
  }

  async function onSave(input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>) {
    try {
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

  function onFrontImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onFrontImageLoad);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onBackImageLoad);
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
        icon="trash-alt"
        accept={{ text: "Delete", variant: "danger", onClick: onDelete }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>Are you sure you want to delete {movie?.name ?? "movie"}?</p>
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

  function renderImage(originalImage: string | undefined | null, editedImage: string | undefined | null, altText: string) {
    let image = originalImage;
    if (isEditing) {
      if (editedImage === null) {
        image = `${image}&default=true`;
      } else if (editedImage) {
        image = editedImage;
      }
    }

    if (image) {
      return <img alt={altText} className="logo w-50" src={image} />
    }
  }

  if (isLoading) return <LoadingIndicator />;

  // TODO: CSS class
  return (
    <div className="row">
      <div className="movie-details col">
        <div className="logo w-100">
          {encodingImage ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            <>
              {renderImage(movie?.front_image_path, frontImage, "Front cover")}
              {renderImage(movie?.back_image_path, backImage, "Back cover")}
            </>
          )}
        </div>

        {!isEditing && movie ? (
          <>
            <MovieDetailsPanel
              movie={movie}
            />
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
          />
        )}
      </div>

      {!isNew && movie && (
        <div className="col-lg-8 col-md-7">
          <MovieScenesPanel movie={movie} />
        </div>
      )}
      {renderDeleteAlert()}
      {renderImageAlert()}
    </div>
  );
};
