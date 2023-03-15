import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import Mousetrap from "mousetrap";
import {
  queryScrapeMovieURL,
  useListMovieScrapers,
} from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { StudioSelect } from "src/components/Shared/Select";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { DurationInput } from "src/components/Shared/DurationInput";
import { URLField } from "src/components/Shared/URLField";
import { useToast } from "src/hooks/Toast";
import { Modal as BSModal, Form, Button, Col, Row } from "react-bootstrap";
import DurationUtils from "src/utils/duration";
import FormUtils from "src/utils/form";
import ImageUtils from "src/utils/image";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { MovieScrapeDialog } from "./MovieScrapeDialog";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { ConfigurationContext } from "src/hooks/Config";
import isEqual from "lodash-es/isEqual";

interface IMovieEditPanel {
  movie: Partial<GQL.MovieDataFragment>;
  onSubmit: (movie: GQL.MovieCreateInput) => void;
  onCancel: () => void;
  onDelete: () => void;
  setFrontImage: (image?: string | null) => void;
  setBackImage: (image?: string | null) => void;
  setEncodingImage: (loading: boolean) => void;
}

export const MovieEditPanel: React.FC<IMovieEditPanel> = ({
  movie,
  onSubmit,
  onCancel,
  onDelete,
  setFrontImage,
  setBackImage,
  setEncodingImage,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const { configuration: stashConfig } = React.useContext(ConfigurationContext);

  const isNew = movie.id === undefined;

  const [isLoading, setIsLoading] = useState(false);
  const [isImageAlertOpen, setIsImageAlertOpen] = useState<boolean>(false);

  const [imageClipboard, setImageClipboard] = useState<string>();

  const Scrapers = useListMovieScrapers();
  const [scrapedMovie, setScrapedMovie] = useState<GQL.ScrapedMovie>();

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yup.string().ensure(),
    duration: yup.number().nullable().defined(),
    date: yup
      .string()
      .ensure()
      .test({
        name: "date",
        test: (value) => {
          if (!value) return true;
          if (!value.match(/^\d{4}-\d{2}-\d{2}$/)) return false;
          if (Number.isNaN(Date.parse(value))) return false;
          return true;
        },
        message: intl.formatMessage({ id: "validation.date_invalid_form" }),
      }),
    studio_id: yup.string().required().nullable(),
    director: yup.string().ensure(),
    rating100: yup.number().nullable().defined(),
    url: yup.string().ensure(),
    synopsis: yup.string().ensure(),
    front_image: yup.string().nullable().optional(),
    back_image: yup.string().nullable().optional(),
  });

  const initialValues = {
    name: movie?.name ?? "",
    aliases: movie?.aliases ?? "",
    duration: movie?.duration ?? null,
    date: movie?.date ?? "",
    studio_id: movie?.studio?.id ?? null,
    director: movie?.director ?? "",
    rating100: movie?.rating100 ?? null,
    url: movie?.url ?? "",
    synopsis: movie?.synopsis ?? "",
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSubmit(values),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating100", v);
  }

  useRatingKeybinds(
    true,
    stashConfig?.ui?.ratingSystemOptions?.type,
    setRating
  );

  function onCancelEditing() {
    setFrontImage(undefined);
    setBackImage(undefined);
    onCancel?.();
  }

  // set up hotkeys
  useEffect(() => {
    // Mousetrap.bind("u", (e) => {
    //   setStudioFocus()
    //   e.preventDefault();
    // });
    Mousetrap.bind("s s", () => formik.handleSubmit());

    return () => {
      // Mousetrap.unbind("u");
      Mousetrap.unbind("s s");
    };
  });

  function updateMovieEditStateFromScraper(
    state: Partial<GQL.ScrapedMovieDataFragment>
  ) {
    if (state.name) {
      formik.setFieldValue("name", state.name);
    }

    if (state.aliases) {
      formik.setFieldValue("aliases", state.aliases);
    }

    if (state.duration) {
      formik.setFieldValue(
        "duration",
        DurationUtils.stringToSeconds(state.duration)
      );
    }

    if (state.date) {
      formik.setFieldValue("date", state.date);
    }

    if (state.studio && state.studio.stored_id) {
      formik.setFieldValue("studio_id", state.studio.stored_id);
    }

    if (state.director) {
      formik.setFieldValue("director", state.director);
    }
    if (state.synopsis) {
      formik.setFieldValue("synopsis", state.synopsis);
    }
    if (state.url) {
      formik.setFieldValue("url", state.url);
    }

    if (state.front_image) {
      // image is a base64 string
      formik.setFieldValue("front_image", state.front_image);
    }
    if (state.back_image) {
      // image is a base64 string
      formik.setFieldValue("back_image", state.back_image);
    }
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

    const currentMovie = {
      id: movie.id!,
      ...formik.values,
    };

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

  const encodingImage = ImageUtils.usePasteImage(showImageAlert);

  useEffect(() => {
    setFrontImage(formik.values.front_image);
  }, [formik.values.front_image, setFrontImage]);

  useEffect(() => {
    setBackImage(formik.values.back_image);
  }, [formik.values.back_image, setBackImage]);

  useEffect(() => {
    setEncodingImage(encodingImage);
  }, [setEncodingImage, encodingImage]);

  function onFrontImageLoad(imageData: string | null) {
    formik.setFieldValue("front_image", imageData);
  }

  function onFrontImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onFrontImageLoad);
  }

  function onBackImageLoad(imageData: string | null) {
    formik.setFieldValue("back_image", imageData);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onBackImageLoad);
  }

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

  function renderTextField(field: string, title: string, placeholder?: string) {
    return (
      <Form.Group controlId={field} as={Row}>
        {FormUtils.renderLabel({
          title,
        })}
        <Col xs={9}>
          <Form.Control
            className="text-input"
            placeholder={placeholder ?? title}
            {...formik.getFieldProps(field)}
            isInvalid={!!formik.getFieldMeta(field).error}
          />
          <Form.Control.Feedback type="invalid">
            {formik.getFieldMeta(field).error}
          </Form.Control.Feedback>
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
              onValueChange={(valueAsNumber) => {
                formik.setFieldValue("duration", valueAsNumber ?? null);
              }}
            />
          </Col>
        </Form.Group>

        {renderTextField(
          "date",
          intl.formatMessage({ id: "date" }),
          "YYYY-MM-DD"
        )}

        <Form.Group controlId="studio" as={Row}>
          {FormUtils.renderLabel({
            title: intl.formatMessage({ id: "studio" }),
          })}
          <Col xs={9}>
            <StudioSelect
              onSelect={(items) =>
                formik.setFieldValue(
                  "studio_id",
                  items.length > 0 ? items[0]?.id : null
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
        onToggleEdit={onCancelEditing}
        onSave={formik.handleSubmit}
        saveDisabled={(!isNew && !formik.dirty) || !isEqual(formik.errors, {})}
        onImageChange={onFrontImageChange}
        onImageChangeURL={onFrontImageLoad}
        onClearImage={() => onFrontImageLoad(null)}
        onBackImageChange={onBackImageChange}
        onBackImageChangeURL={onBackImageLoad}
        onClearBackImage={() => onBackImageLoad(null)}
        onDelete={onDelete}
      />

      {maybeRenderScrapeDialog()}
      {renderImageAlert()}
    </div>
  );
};
