import React, { useEffect, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import {
  queryScrapeMovieURL,
  useListMovieScrapers,
} from "src/core/StashService";
import {
  LoadingIndicator,
  StudioSelect,
  Icon,
  DetailsEditNavbar,
  DurationInput,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Form, Button, Col, Row, InputGroup } from "react-bootstrap";
import { DurationUtils, ImageUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { MovieScrapeDialog } from "./MovieScrapeDialog";

interface IMovieEditPanel {
  movie?: Partial<GQL.MovieDataFragment>;
  onSubmit: (
    movie: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput>
  ) => void;
  onCancel: () => void;
  onDelete: () => void;
  setFrontImage: (image?: string | null) => void;
  setBackImage: (image?: string | null) => void;
}

export const MovieEditPanel: React.FC<IMovieEditPanel> = ({
  movie,
  onSubmit,
  onCancel,
  onDelete,
  setFrontImage,
  setBackImage,
}) => {
  const Toast = useToast();

  const isNew = movie === undefined;

  const [isLoading, setIsLoading] = useState(false);

  const Scrapers = useListMovieScrapers();
  const [scrapedMovie, setScrapedMovie] = useState<
    GQL.ScrapedMovie | undefined
  >();

  const labelXS = 3;
  const labelXL = 3;
  const fieldXS = 9;
  const fieldXL = 9;

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yup.string().optional().nullable(),
    duration: yup.string().optional().nullable(),
    date: yup
      .string()
      .optional()
      .nullable()
      .matches(/^\d{4}-\d{2}-\d{2}$/),
    rating: yup.number().optional().nullable(),
    studio_id: yup.string().optional().nullable(),
    director: yup.string().optional().nullable(),
    synopsis: yup.string().optional().nullable(),
    url: yup.string().optional().nullable(),
  });

  const initialValues = {
    name: movie?.name,
    aliases: movie?.aliases,
    duration: movie?.duration,
    date: movie?.date,
    rating: movie?.rating ?? null,
    studio_id: movie?.studio?.id,
    director: movie?.director,
    synopsis: movie?.synopsis,
    url: movie?.url,
  };

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSubmit(getMovieInput(values)),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating", v);
  }

  // set up hotkeys
  useEffect(() => {
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
    Mousetrap.bind("s s", () => formik.handleSubmit());

    return () => {
      Mousetrap.unbind("r 0");
      Mousetrap.unbind("r 1");
      Mousetrap.unbind("r 2");
      Mousetrap.unbind("r 3");
      Mousetrap.unbind("r 4");
      Mousetrap.unbind("r 5");
      // Mousetrap.unbind("u");
      Mousetrap.unbind("s s");
    };
  });

  function getMovieInput(values: InputValues) {
    const input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      ...values,
      rating: values.rating ?? null,
      studio_id: values.studio_id ?? null,
    };

    if (movie && movie.id) {
      (input as GQL.MovieUpdateInput).id = movie.id;
    }
    return input;
  }

  function updateMovieEditStateFromScraper(
    state: Partial<GQL.ScrapedMovieDataFragment>
  ) {
    if (state.name) {
      formik.setFieldValue("name", state.name);
    }

    if (state.aliases) {
      formik.setFieldValue("aliases", state.aliases ?? undefined);
    }

    if (state.duration) {
      formik.setFieldValue(
        "duration",
        DurationUtils.stringToSeconds(state.duration) ?? undefined
      );
    }

    if (state.date) {
      formik.setFieldValue("date", state.date ?? undefined);
    }

    if (state.studio && state.studio.id) {
      formik.setFieldValue("studio_id", state.studio.id ?? undefined);
    }

    if (state.director) {
      formik.setFieldValue("director", state.director ?? undefined);
    }
    if (state.synopsis) {
      formik.setFieldValue("synopsis", state.synopsis ?? undefined);
    }
    if (state.url) {
      formik.setFieldValue("url", state.url ?? undefined);
    }

    const imageStr = (state as GQL.ScrapedMovieDataFragment).front_image;
    setFrontImage(imageStr ?? undefined);

    const backImageStr = (state as GQL.ScrapedMovieDataFragment).back_image;
    setBackImage(backImageStr ?? undefined);
  }

  async function onScrapeMovieURL() {
    const { url } = formik.values;
    if (!url) return;
    setIsLoading(true);

    try {
      const result = await queryScrapeMovieURL(url);
      if (!result.data || !result.data.scrapeMovieURL) {
        return;
      }

      // if this is a new movie, just dump the data
      if (isNew) {
        updateMovieEditStateFromScraper(result.data.scrapeMovieURL);
      } else {
        setScrapedMovie(result.data.scrapeMovieURL);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function urlScrapable(scrapedUrl: string) {
    return (
      !!scrapedUrl &&
      (Scrapers?.data?.listMovieScrapers ?? []).some((s) =>
        (s?.movie?.urls ?? []).some((u) => scrapedUrl.includes(u))
      )
    );
  }

  function maybeRenderScrapeButton() {
    const { url } = formik.values;
    if (!url || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button
        className="minimal scrape-url-button"
        onClick={() => onScrapeMovieURL()}
      >
        <Icon icon="file-upload" />
      </Button>
    );
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedMovie) {
      return;
    }

    const currentMovie = getMovieInput(formik.values);

    // Get image paths for scrape gui
    currentMovie.front_image = movie?.front_image_path;
    currentMovie.back_image = movie?.back_image_path;

    return (
      <MovieScrapeDialog
        movie={currentMovie}
        scraped={scrapedMovie}
        onClose={(m) => {
          onScrapeDialogClosed(m);
        }}
      />
    );
  }

  function onScrapeDialogClosed(p?: GQL.ScrapedMovieDataFragment) {
    if (p) {
      updateMovieEditStateFromScraper(p);
    }
    setScrapedMovie(undefined);
  }

  function onFrontImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, setFrontImage);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, setBackImage);
  }

  if (isLoading) return <LoadingIndicator />;

  const isEditing = true;

  function renderTextField(field: string, title: string) {
    return (
      <Form.Group controlId={field} as={Row}>
        <Form.Label column xs={labelXS} xl={labelXL}>
          {title}
        </Form.Label>
        <Col xs={fieldXS} xl={fieldXL}>
          <Form.Control
            className="text-input"
            placeholder={title}
            {...formik.getFieldProps(field)}
            isInvalid={!!formik.getFieldMeta(field).error}
          />
        </Col>
      </Form.Group>
    );
  }

  // TODO: CSS class
  return (
    <div>
      {isNew && <h2>Add Movie</h2>}

      <Prompt
        when={formik.dirty}
        message="Unsaved changes. Are you sure you want to leave?"
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="movie-edit">
        <Form.Group controlId="name" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            Name
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <Form.Control
              className="text-input"
              placeholder="Name"
              {...formik.getFieldProps("name")}
              isInvalid={!!formik.errors.name}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.name}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        {renderTextField("aliases", "Aliases")}

        <Form.Group controlId="duration" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Duration
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <DurationInput
              numericValue={formik.values.duration ?? undefined}
              onValueChange={(valueAsNumber: number) => {
                formik.setFieldValue("duration", valueAsNumber);
              }}
            />
          </Col>
        </Form.Group>

        {renderTextField("date", "Date (YYYY-MM-DD)")}

        <Form.Group controlId="studio" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Studio
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <StudioSelect
              onSelect={(items) =>
                formik.setFieldValue(
                  "studio_id",
                  items.length > 0 ? items[0]?.id : undefined
                )
              }
              ids={formik.values.studio_id ? [formik.values.studio_id] : []}
            />
          </Col>
        </Form.Group>

        {renderTextField("director", "Director")}

        <Form.Group controlId="rating" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Rating
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <RatingStars
              value={formik.values.rating ?? undefined}
              onSetRating={(value) => formik.setFieldValue("rating", value ?? null)}
            />
          </Col>
        </Form.Group>

        <Form.Group controlId="url" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            URL
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <InputGroup>
              <Form.Control
                className="text-input"
                placeholder="URL"
                {...formik.getFieldProps("url")}
              />
              <InputGroup.Append>{maybeRenderScrapeButton()}</InputGroup.Append>
            </InputGroup>
          </Col>
        </Form.Group>

        <Form.Group controlId="synopsis" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Synopsis
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder="Synopsis"
              {...formik.getFieldProps("synopsis")}
            />
          </Col>
        </Form.Group>
      </Form>

      <DetailsEditNavbar
        objectName={movie?.name ?? "movie"}
        isNew={isNew}
        isEditing={isEditing}
        onToggleEdit={onCancel}
        onSave={() => formik.handleSubmit()}
        saveDisabled={!formik.dirty}
        onImageChange={onFrontImageChange}
        onImageChangeURL={setFrontImage}
        onClearImage={() => {
          setFrontImage(null);
        }}
        onBackImageChange={onBackImageChange}
        onBackImageChangeURL={setBackImage}
        onClearBackImage={() => {
          setBackImage(null);
        }}
        onDelete={onDelete}
      />

      {maybeRenderScrapeDialog()}
    </div>
  );
};
