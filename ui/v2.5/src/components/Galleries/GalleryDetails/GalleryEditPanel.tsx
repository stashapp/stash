import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Prompt } from "react-router-dom";
import {
  Button,
  Dropdown,
  DropdownButton,
  Form,
  Col,
  Row,
} from "react-bootstrap";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import {
  queryScrapeGallery,
  queryScrapeGalleryURL,
  useListGalleryScrapers,
  mutateReloadScrapers,
} from "src/core/StashService";
import {
  TagSelect,
  SceneSelect,
  StudioSelect,
} from "src/components/Shared/Select";
import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { URLField } from "src/components/Shared/URLField";
import { useToast } from "src/hooks/Toast";
import { useFormik } from "formik";
import FormUtils from "src/utils/form";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { GalleryScrapeDialog } from "./GalleryScrapeDialog";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";
import { galleryTitle } from "src/core/galleries";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { ConfigurationContext } from "src/hooks/Config";
import isEqual from "lodash-es/isEqual";
import { DateInput } from "src/components/Shared/DateInput";
import { handleUnsavedChanges } from "src/utils/navigation";
import {
  Performer,
  PerformerSelect,
} from "src/components/Performers/PerformerSelect";

interface IProps {
  gallery: Partial<GQL.GalleryDataFragment>;
  isVisible: boolean;
  onSubmit: (input: GQL.GalleryCreateInput) => Promise<void>;
  onDelete: () => void;
}

export const GalleryEditPanel: React.FC<IProps> = ({
  gallery,
  isVisible,
  onSubmit,
  onDelete,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [scenes, setScenes] = useState<{ id: string; title: string }[]>(
    (gallery?.scenes ?? []).map((s) => ({
      id: s.id,
      title: galleryTitle(s),
    }))
  );

  const [performers, setPerformers] = useState<Performer[]>([]);

  const isNew = gallery.id === undefined;
  const { configuration: stashConfig } = React.useContext(ConfigurationContext);

  const Scrapers = useListGalleryScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [scrapedGallery, setScrapedGallery] =
    useState<GQL.ScrapedGallery | null>();

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const titleRequired =
    isNew || (gallery?.files?.length === 0 && !gallery?.folder);

  const schema = yup.object({
    title: titleRequired ? yup.string().required() : yup.string().ensure(),
    url: yup.string().ensure(),
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
    rating100: yup.number().nullable().defined(),
    studio_id: yup.string().required().nullable(),
    performer_ids: yup.array(yup.string().required()).defined(),
    tag_ids: yup.array(yup.string().required()).defined(),
    scene_ids: yup.array(yup.string().required()).defined(),
    details: yup.string().ensure(),
  });

  const initialValues = {
    title: gallery?.title ?? "",
    url: gallery?.url ?? "",
    date: gallery?.date ?? "",
    rating100: gallery?.rating100 ?? null,
    studio_id: gallery?.studio?.id ?? null,
    performer_ids: (gallery?.performers ?? []).map((p) => p.id),
    tag_ids: (gallery?.tags ?? []).map((t) => t.id),
    scene_ids: (gallery?.scenes ?? []).map((s) => s.id),
    details: gallery?.details ?? "",
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSave(values),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating100", v);
  }

  interface ISceneSelectValue {
    id: string;
    title: string;
  }

  function onSetScenes(items: ISceneSelectValue[]) {
    setScenes(items);
    formik.setFieldValue(
      "scene_ids",
      items.map((i) => i.id)
    );
  }

  function onSetPerformers(items: Performer[]) {
    setPerformers(items);
    formik.setFieldValue(
      "performer_ids",
      items.map((item) => item.id)
    );
  }

  useRatingKeybinds(
    isVisible,
    stashConfig?.ui?.ratingSystemOptions?.type,
    setRating
  );

  useEffect(() => {
    setPerformers(gallery.performers ?? []);
  }, [gallery.performers]);

  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        if (formik.dirty) {
          formik.submitForm();
        }
      });
      Mousetrap.bind("d d", () => {
        onDelete();
      });

      return () => {
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");
      };
    }
  });

  useEffect(() => {
    const newQueryableScrapers = (
      Scrapers?.data?.listGalleryScrapers ?? []
    ).filter((s) =>
      s.gallery?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );

    setQueryableScrapers(newQueryableScrapers);
  }, [Scrapers]);

  async function onSave(input: InputValues) {
    setIsLoading(true);
    try {
      await onSubmit(input);
      formik.resetForm();
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  async function onScrapeClicked(scraper: GQL.Scraper) {
    if (!gallery || !gallery.id) return;

    setIsLoading(true);
    try {
      const result = await queryScrapeGallery(scraper.id, gallery.id);
      if (!result.data || !result.data.scrapeSingleGallery?.length) {
        Toast.success({
          content: "No galleries found",
        });
        return;
      }
      setScrapedGallery(result.data.scrapeSingleGallery[0]);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onReloadScrapers() {
    setIsLoading(true);
    try {
      await mutateReloadScrapers();

      // reload the performer scrapers
      await Scrapers.refetch();
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function onScrapeDialogClosed(data?: GQL.ScrapedGalleryDataFragment) {
    if (data) {
      updateGalleryFromScrapedGallery(data);
    }
    setScrapedGallery(undefined);
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedGallery) {
      return;
    }

    const currentGallery = {
      id: gallery.id!,
      ...formik.values,
    };

    return (
      <GalleryScrapeDialog
        gallery={currentGallery}
        galleryPerformers={performers}
        scraped={scrapedGallery}
        onClose={(data) => {
          onScrapeDialogClosed(data);
        }}
      />
    );
  }

  function renderScraperMenu() {
    if (isNew) {
      return;
    }

    return (
      <DropdownButton
        className="d-inline-block"
        id="gallery-scrape"
        title={intl.formatMessage({ id: "actions.scrape_with" })}
      >
        {queryableScrapers.map((s) => (
          <Dropdown.Item key={s.name} onClick={() => onScrapeClicked(s)}>
            {s.name}
          </Dropdown.Item>
        ))}
        <Dropdown.Item onClick={() => onReloadScrapers()}>
          <span className="fa-icon">
            <Icon icon={faSyncAlt} />
          </span>
          <span>
            <FormattedMessage id="actions.reload_scrapers" />
          </span>
        </Dropdown.Item>
      </DropdownButton>
    );
  }

  function urlScrapable(scrapedUrl: string): boolean {
    return (Scrapers?.data?.listGalleryScrapers ?? []).some((s) =>
      (s?.gallery?.urls ?? []).some((u) => scrapedUrl.includes(u))
    );
  }

  function updateGalleryFromScrapedGallery(
    galleryData: GQL.ScrapedGalleryDataFragment
  ) {
    if (galleryData.title) {
      formik.setFieldValue("title", galleryData.title);
    }

    if (galleryData.details) {
      formik.setFieldValue("details", galleryData.details);
    }

    if (galleryData.date) {
      formik.setFieldValue("date", galleryData.date);
    }

    if (galleryData.url) {
      formik.setFieldValue("url", galleryData.url);
    }

    if (galleryData.studio?.stored_id) {
      formik.setFieldValue("studio_id", galleryData.studio.stored_id);
    }

    if (galleryData.performers?.length) {
      const idPerfs = galleryData.performers.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idPerfs.length > 0) {
        onSetPerformers(
          idPerfs.map((p) => {
            return {
              id: p.stored_id!,
              name: p.name ?? "",
              alias_list: [],
            };
          })
        );
      }
    }

    if (galleryData?.tags?.length) {
      const idTags = galleryData.tags.filter((t) => {
        return t.stored_id !== undefined && t.stored_id !== null;
      });

      if (idTags.length > 0) {
        const newIds = idTags.map((t) => t.stored_id);
        formik.setFieldValue("tag_ids", newIds as string[]);
      }
    }
  }

  async function onScrapeGalleryURL() {
    if (!formik.values.url) {
      return;
    }
    setIsLoading(true);
    try {
      const result = await queryScrapeGalleryURL(formik.values.url);
      if (!result || !result.data || !result.data.scrapeGalleryURL) {
        return;
      }
      setScrapedGallery(result.data.scrapeGalleryURL);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

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

  if (isLoading) return <LoadingIndicator />;

  return (
    <div id="gallery-edit-details">
      <Prompt
        when={formik.dirty}
        message={handleUnsavedChanges(intl, "galleries", gallery?.id)}
      />

      {maybeRenderScrapeDialog()}
      <Form noValidate onSubmit={formik.handleSubmit}>
        <div className="form-container row px-3 pt-3">
          <div className="col edit-buttons mb-3 pl-0">
            <Button
              className="edit-button"
              variant="primary"
              disabled={
                (!isNew && !formik.dirty) || !isEqual(formik.errors, {})
              }
              onClick={() => formik.submitForm()}
            >
              <FormattedMessage id="actions.save" />
            </Button>
            <Button
              className="edit-button"
              variant="danger"
              onClick={() => onDelete()}
            >
              <FormattedMessage id="actions.delete" />
            </Button>
          </div>
          <Col xs={6} className="text-right">
            {renderScraperMenu()}
          </Col>
        </div>
        <div className="form-container row px-3">
          <div className="col-12 col-lg-6 col-xl-12">
            {renderTextField("title", intl.formatMessage({ id: "title" }))}
            <Form.Group controlId="url" as={Row}>
              <Col xs={3} className="pr-0 url-label">
                <Form.Label className="col-form-label">
                  <FormattedMessage id="url" />
                </Form.Label>
              </Col>
              <Col xs={9}>
                <URLField
                  {...formik.getFieldProps("url")}
                  onScrapeClick={onScrapeGalleryURL}
                  urlScrapable={urlScrapable}
                  isInvalid={!!formik.getFieldMeta("url").error}
                />
              </Col>
            </Form.Group>
            <Form.Group controlId="date" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "date" }),
              })}
              <Col xs={9}>
                <DateInput
                  value={formik.values.date}
                  onValueChange={(value) => formik.setFieldValue("date", value)}
                  error={formik.errors.date}
                />
              </Col>
            </Form.Group>
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

            <Form.Group controlId="performers" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "performers" }),
                labelProps: {
                  column: true,
                  sm: 3,
                  xl: 12,
                },
              })}
              <Col sm={9} xl={12}>
                <PerformerSelect
                  isMulti
                  onSelect={onSetPerformers}
                  values={performers}
                />
              </Col>
            </Form.Group>

            <Form.Group controlId="tags" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "tags" }),
                labelProps: {
                  column: true,
                  sm: 3,
                  xl: 12,
                },
              })}
              <Col sm={9} xl={12}>
                <TagSelect
                  isMulti
                  onSelect={(items) =>
                    formik.setFieldValue(
                      "tag_ids",
                      items.map((item) => item.id)
                    )
                  }
                  ids={formik.values.tag_ids}
                  hoverPlacement="right"
                />
              </Col>
            </Form.Group>

            <Form.Group controlId="scenes" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "scenes" }),
                labelProps: {
                  column: true,
                  sm: 3,
                  xl: 12,
                },
              })}
              <Col sm={9} xl={12}>
                <SceneSelect
                  selected={scenes}
                  onSelect={(items) => onSetScenes(items)}
                  isMulti
                />
              </Col>
            </Form.Group>
          </div>
          <div className="col-12 col-lg-6 col-xl-12">
            <Form.Group controlId="details">
              <Form.Label>
                <FormattedMessage id="details" />
              </Form.Label>
              <Form.Control
                as="textarea"
                className="gallery-description text-input"
                onChange={(e) =>
                  formik.setFieldValue("details", e.currentTarget.value)
                }
                value={formik.values.details}
              />
            </Form.Group>
          </div>
        </div>
      </Form>
    </div>
  );
};
