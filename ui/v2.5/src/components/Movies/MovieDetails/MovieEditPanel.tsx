import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
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
  DetailsEditNavbar,
  DurationInput,
  URLField,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { Modal as BSModal, Form, Button, Col, Row } from "react-bootstrap";
import { DurationUtils, FormUtils, ImageUtils } from "src/utils";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
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
  onImageEncoding: (loading?: boolean) => void;
}

export const MovieEditPanel: React.FC<IMovieEditPanel> = ({
  movie,
  onSubmit,
  onCancel,
  onDelete,
  setFrontImage,
  setBackImage,
  onImageEncoding,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  const isNew = movie === undefined;

  const [isLoading, setIsLoading] = useState(false);
  const [isImageAlertOpen, setIsImageAlertOpen] = useState<boolean>(false);

  const [imageClipboard, setImageClipboard] = useState<string | undefined>(
    undefined
  );

  const Scrapers = useListMovieScrapers();
  const [scrapedMovie, setScrapedMovie] = useState<
    GQL.ScrapedMovie | undefined
  >();

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yup.string().optional().nullable(),
    duration: yup.string().optional().nullable(),
    date: yup
      .string()
      .optional()
      .nullable()
      .matches(/^\d{4}-\d{2}-\d{2}$/),
    rating100: yup.number().optional().nullable(),
    studio_id: yup.string().optional().nullable(),
    director: yup.string().optional().nullable(),
    synopsis: yup.string().optional().nullable(),
    url: yup.string().optional().nullable(),
    front_image: yup.string().optional().nullable(),
    back_image: yup.string().optional().nullable(),
  });

  const initialValues = {
    name: movie?.name,
    aliases: movie?.aliases,
    duration: movie?.duration,
    date: movie?.date,
    rating100: movie?.rating100 ?? null,
    studio_id: movie?.studio?.id,
    director: movie?.director,
    synopsis: movie?.synopsis,
    url: movie?.url,
    front_image: undefined,
    back_image: undefined,
  };

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSubmit(getMovieInput(values)),
  });

  const encodingImage = ImageUtils.usePasteImage(showImageAlert);

  useEffect(() => {
    setFrontImage(formik.values.front_image);
  }, [formik.values.front_image, setFrontImage]);

  useEffect(() => {
    setBackImage(formik.values.back_image);
  }, [formik.values.back_image, setBackImage]);

  useEffect(() => onImageEncoding(encodingImage), [
    onImageEncoding,
    encodingImage,
  ]);

  function setRating(v: number) {
    formik.setFieldValue("rating100", v);
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("r 0", () => setRating(NaN));
    Mousetrap.bind("r 1", () => setRating(20));
    Mousetrap.bind("r 2", () => setRating(40));
    Mousetrap.bind("r 3", () => setRating(60));
    Mousetrap.bind("r 4", () => setRating(80));
    Mousetrap.bind("r 5", () => setRating(100));
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

  function showImageAlert(imageData: string) {
    setImageClipboard(imageData);
    setIsImageAlertOpen(true);
  }

  function setImageFromClipboard(isFrontImage: boolean) {
    if (isFrontImage) {
      formik.setFieldValue("front_image", imageClipboard);
    } else {
      formik.setFieldValue("back_image", imageClipboard);
    }

    setImageClipboard(undefined);
    setIsImageAlertOpen(false);
  }

  function getMovieInput(values: InputValues) {
    const input: Partial<GQL.MovieCreateInput | GQL.MovieUpdateInput> = {
      ...values,
      rating100: values.rating100 ?? null,
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

    if (state.studio && state.studio.stored_id) {
      formik.setFieldValue("studio_id", state.studio.stored_id ?? undefined);
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
    formik.setFieldValue("front_image", imageStr ?? undefined);

    const backImageStr = (state as GQL.ScrapedMovieDataFragment).back_image;
    formik.setFieldValue("back_image", backImageStr ?? undefined);
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
    ImageUtils.onImageChange(event, (data) =>
      formik.setFieldValue("front_image", data)
    );
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, (data) =>
      formik.setFieldValue("back_image", data)
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
              <FormattedMessage id="actions.cancel" />
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

  if (isLoading) return <LoadingIndicator />;

  const isEditing = true;

  function renderTextField(field: string, title: string) {
    return (
      <Form.Group controlId={field} as={Row}>
        {FormUtils.renderLabel({
          title,
        })}
        <Col xs={9}>
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
      {isNew && (
        <h2>
          {intl.formatMessage(
            { id: "actions.add_entity" },
            { entityType: intl.formatMessage({ id: "movie" }) }
          )}
        </h2>
      )}

      <Prompt
        when={formik.dirty}
        message={(location, action) => {
          // Check if it's a redirect after movie creation
          if (action === "PUSH" && location.pathname.startsWith("/movies/"))
            return true;
          return intl.formatMessage({ id: "dialogs.unsaved_changes" });
        }}
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="movie-edit">
        <Form.Group controlId="name" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "name" }),
          })}
          <Col xs={9}>
            <Form.Control
              className="text-input"
              placeholder={intl.formatMessage({ id: "name" })}
              {...formik.getFieldProps("name")}
              isInvalid={!!formik.errors.name}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.name}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        {renderTextField("aliases", intl.formatMessage({ id: "aliases" }))}

        <Form.Group controlId="duration" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "duration" }),
          })}
          <Col xs={9}>
            <DurationInput
              numericValue={formik.values.duration ?? undefined}
              onValueChange={(valueAsNumber: number) => {
                formik.setFieldValue("duration", valueAsNumber);
              }}
            />
          </Col>
        </Form.Group>

        {renderTextField("date", intl.formatMessage({ id: "date" }))}

        <Form.Group controlId="studio" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "studio" }),
          })}
          <Col xs={9}>
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

        {renderTextField("director", intl.formatMessage({ id: "director" }))}

        <Form.Group controlId="rating" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "rating" }),
          })}
          <Col xs={9}>
            <RatingSystem
              value={formik.values.rating100 ?? undefined}
              onSetRating={(value) =>
                formik.setFieldValue("rating100", value ?? null)
              }
            />
          </Col>
        </Form.Group>
        <Form.Group controlId="url" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "url" }),
          })}
          <Col xs={9}>
            <URLField
              {...formik.getFieldProps("url")}
              onScrapeClick={onScrapeMovieURL}
              urlScrapable={urlScrapable}
            />
          </Col>
        </Form.Group>

        <Form.Group controlId="synopsis" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "synopsis" }),
          })}
          <Col xs={9}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder={intl.formatMessage({ id: "synopsis" })}
              {...formik.getFieldProps("synopsis")}
            />
          </Col>
        </Form.Group>
      </Form>

      <DetailsEditNavbar
        objectName={movie?.name ?? intl.formatMessage({ id: "movie" })}
        isNew={isNew}
        isEditing={isEditing}
        onToggleEdit={onCancel}
        onSave={() => formik.handleSubmit()}
        saveDisabled={!formik.dirty}
        onImageChange={onFrontImageChange}
        onImageChangeURL={(i) => formik.setFieldValue("front_image", i)}
        onClearImage={() => {
          formik.setFieldValue("front_image", null);
        }}
        onBackImageChange={onBackImageChange}
        onBackImageChangeURL={(i) => formik.setFieldValue("back_image", i)}
        onClearBackImage={() => {
          formik.setFieldValue("back_image", null);
        }}
        onDelete={onDelete}
      />

      {maybeRenderScrapeDialog()}
      {renderImageAlert()}
    </div>
  );
};
