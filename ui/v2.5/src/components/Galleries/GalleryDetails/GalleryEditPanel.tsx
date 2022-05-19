import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory, Prompt } from "react-router-dom";
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
  useGalleryCreate,
  useGalleryUpdate,
  useListGalleryScrapers,
  mutateReloadScrapers,
} from "src/core/StashService";
import {
  PerformerSelect,
  TagSelect,
  SceneSelect,
  StudioSelect,
  Icon,
  LoadingIndicator,
  URLField,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { useFormik } from "formik";
import { FormUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { GalleryScrapeDialog } from "./GalleryScrapeDialog";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";
import { galleryTitle } from "src/core/galleries";

interface IProps {
  isVisible: boolean;
  onDelete: () => void;
}

interface INewProps {
  isNew: true;
  gallery: undefined;
}

interface IExistingProps {
  isNew: false;
  gallery: GQL.GalleryDataFragment;
}

export const GalleryEditPanel: React.FC<
  IProps & (INewProps | IExistingProps)
> = ({ gallery, isNew, isVisible, onDelete }) => {
  const intl = useIntl();
  const Toast = useToast();
  const history = useHistory();
  const [scenes, setScenes] = useState<{ id: string; title: string }[]>(
    (gallery?.scenes ?? []).map((s) => ({
      id: s.id,
      title: galleryTitle(s),
    }))
  );

  const Scrapers = useListGalleryScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [
    scrapedGallery,
    setScrapedGallery,
  ] = useState<GQL.ScrapedGallery | null>();

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [createGallery] = useGalleryCreate();
  const [updateGallery] = useGalleryUpdate();

  const schema = yup.object({
    title: yup.string().required(),
    details: yup.string().optional().nullable(),
    url: yup.string().optional().nullable(),
    date: yup.string().optional().nullable(),
    rating: yup.number().optional().nullable(),
    studio_id: yup.string().optional().nullable(),
    performer_ids: yup.array(yup.string().required()).optional().nullable(),
    tag_ids: yup.array(yup.string().required()).optional().nullable(),
    scene_ids: yup.array(yup.string().required()).optional().nullable(),
  });

  const initialValues = {
    title: gallery?.title ?? "",
    details: gallery?.details ?? "",
    url: gallery?.url ?? "",
    date: gallery?.date ?? "",
    rating: gallery?.rating ?? null,
    studio_id: gallery?.studio?.id,
    performer_ids: (gallery?.performers ?? []).map((p) => p.id),
    tag_ids: (gallery?.tags ?? []).map((t) => t.id),
    scene_ids: (gallery?.scenes ?? []).map((s) => s.id),
  };

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSave(getGalleryInput(values)),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating", v);
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

  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        formik.handleSubmit();
      });
      Mousetrap.bind("d d", () => {
        onDelete();
      });

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
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");

        Mousetrap.unbind("r");
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

  function getGalleryInput(
    input: InputValues
  ): GQL.GalleryCreateInput | GQL.GalleryUpdateInput {
    return {
      id: isNew ? undefined : gallery?.id ?? "",
      ...input,
    };
  }

  async function onSave(
    input: GQL.GalleryCreateInput | GQL.GalleryUpdateInput
  ) {
    setIsLoading(true);
    try {
      if (isNew) {
        const result = await createGallery({
          variables: {
            input: input as GQL.GalleryCreateInput,
          },
        });
        if (result.data?.galleryCreate) {
          history.push(`/galleries/${result.data.galleryCreate.id}`);
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.created_entity" },
              {
                entity: intl
                  .formatMessage({ id: "gallery" })
                  .toLocaleLowerCase(),
              }
            ),
          });
        }
      } else {
        const result = await updateGallery({
          variables: {
            input: input as GQL.GalleryUpdateInput,
          },
        });
        if (result.data?.galleryUpdate) {
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "gallery" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          formik.resetForm({ values: formik.values });
        }
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  async function onScrapeClicked(scraper: GQL.Scraper) {
    if (!gallery) return;

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

    const currentGallery = getGalleryInput(formik.values);

    return (
      <GalleryScrapeDialog
        gallery={currentGallery}
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
        const newIds = idPerfs.map((p) => p.stored_id);
        formik.setFieldValue("performer_ids", newIds as string[]);
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
      <Form.Group controlId={title} as={Row}>
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
        message={intl.formatMessage({ id: "dialogs.unsaved_changes" })}
      />

      {maybeRenderScrapeDialog()}
      <Form noValidate onSubmit={formik.handleSubmit}>
        <div className="form-container row px-3 pt-3">
          <div className="col edit-buttons mb-3 pl-0">
            <Button
              className="edit-button"
              variant="primary"
              disabled={!formik.dirty}
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
            {renderTextField(
              "date",
              intl.formatMessage({ id: "date" }),
              "YYYY-MM-DD"
            )}
            <Form.Group controlId="rating" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "rating" }),
              })}
              <Col xs={9}>
                <RatingStars
                  value={formik.values.rating ?? undefined}
                  onSetRating={(value) =>
                    formik.setFieldValue("rating", value ?? null)
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
                  onSelect={(items) =>
                    formik.setFieldValue(
                      "performer_ids",
                      items.map((item) => item.id)
                    )
                  }
                  ids={formik.values.performer_ids}
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
                  scenes={scenes}
                  onSelect={(items) => onSetScenes(items)}
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
                onChange={(newValue: React.ChangeEvent<HTMLTextAreaElement>) =>
                  formik.setFieldValue("details", newValue.currentTarget.value)
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
