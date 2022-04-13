import React, { useEffect, useState, useMemo } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import {
  Button,
  Dropdown,
  DropdownButton,
  Form,
  Col,
  Row,
  ButtonGroup,
} from "react-bootstrap";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import {
  queryScrapeScene,
  queryScrapeSceneURL,
  useListSceneScrapers,
  useSceneUpdate,
  mutateReloadScrapers,
  queryScrapeSceneQueryFragment,
} from "src/core/StashService";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
  GallerySelect,
  Icon,
  LoadingIndicator,
  ImageInput,
  URLField,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { ImageUtils, FormUtils, TextUtils, getStashIDs } from "src/utils";
import { MovieSelect } from "src/components/Shared/Select";
import { useFormik } from "formik";
import { Prompt } from "react-router-dom";
import { ConfigurationContext } from "src/hooks/Config";
import { stashboxDisplayName } from "src/utils/stashbox";
import { SceneMovieTable } from "./SceneMovieTable";
import { RatingStars } from "./RatingStars";
import { SceneScrapeDialog } from "./SceneScrapeDialog";
import { SceneQueryModal } from "./SceneQueryModal";

interface IProps {
  scene: GQL.SceneDataFragment;
  isVisible: boolean;
  onDelete: () => void;
  onUpdate?: () => void;
}

export const SceneEditPanel: React.FC<IProps> = ({
  scene,
  isVisible,
  onDelete,
}) => {
  const intl = useIntl();
  const Toast = useToast();
  const [galleries, setGalleries] = useState<{ id: string; title: string }[]>(
    scene.galleries.map((g) => ({
      id: g.id,
      title: g.title ?? TextUtils.fileNameFromPath(g.path ?? ""),
    }))
  );

  const Scrapers = useListSceneScrapers();
  const [fragmentScrapers, setFragmentScrapers] = useState<GQL.Scraper[]>([]);
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [scraper, setScraper] = useState<GQL.ScraperSourceInput | undefined>();
  const [
    isScraperQueryModalOpen,
    setIsScraperQueryModalOpen,
  ] = useState<boolean>(false);
  const [scrapedScene, setScrapedScene] = useState<GQL.ScrapedScene | null>();
  const [endpoint, setEndpoint] = useState<string | undefined>();

  const [coverImagePreview, setCoverImagePreview] = useState<
    string | undefined
  >(scene.paths.screenshot ?? undefined);

  const { configuration: stashConfig } = React.useContext(ConfigurationContext);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [updateScene] = useSceneUpdate();

  const schema = yup.object({
    title: yup.string().optional().nullable(),
    details: yup.string().optional().nullable(),
    url: yup.string().optional().nullable(),
    date: yup.string().optional().nullable(),
    rating: yup.number().optional().nullable(),
    gallery_ids: yup.array(yup.string().required()).optional().nullable(),
    studio_id: yup.string().optional().nullable(),
    performer_ids: yup.array(yup.string().required()).optional().nullable(),
    movies: yup
      .object({
        movie_id: yup.string().required(),
        scene_index: yup.string().optional().nullable(),
      })
      .optional()
      .nullable(),
    tag_ids: yup.array(yup.string().required()).optional().nullable(),
    cover_image: yup.string().optional().nullable(),
    stash_ids: yup.mixed<GQL.StashIdInput>().optional().nullable(),
  });

  const initialValues = useMemo(
    () => ({
      title: scene.title ?? "",
      details: scene.details ?? "",
      url: scene.url ?? "",
      date: scene.date ?? "",
      rating: scene.rating ?? null,
      gallery_ids: (scene.galleries ?? []).map((g) => g.id),
      studio_id: scene.studio?.id,
      performer_ids: (scene.performers ?? []).map((p) => p.id),
      movies: (scene.movies ?? []).map((m) => {
        return { movie_id: m.movie.id, scene_index: m.scene_index };
      }),
      tag_ids: (scene.tags ?? []).map((t) => t.id),
      cover_image: undefined,
      stash_ids: getStashIDs(scene.stash_ids),
    }),
    [scene]
  );

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSave(getSceneInput(values)),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating", v);
  }

  interface IGallerySelectValue {
    id: string;
    title: string;
  }

  function onSetGalleries(items: IGallerySelectValue[]) {
    setGalleries(items);
    formik.setFieldValue(
      "gallery_ids",
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
    const toFilter = Scrapers?.data?.listSceneScrapers ?? [];

    const newFragmentScrapers = toFilter.filter((s) =>
      s.scene?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );
    const newQueryableScrapers = toFilter.filter((s) =>
      s.scene?.supported_scrapes.includes(GQL.ScrapeType.Name)
    );

    setFragmentScrapers(newFragmentScrapers);
    setQueryableScrapers(newQueryableScrapers);
  }, [Scrapers, stashConfig]);

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, true);

  function getSceneInput(input: InputValues): GQL.SceneUpdateInput {
    return {
      id: scene.id,
      ...input,
    };
  }

  function setMovieIds(movieIds: string[]) {
    const existingMovies = formik.values.movies;

    const newMovies = movieIds.map((m) => {
      const existing = existingMovies.find((mm) => mm.movie_id === m);
      if (existing) {
        return existing;
      }

      return {
        movie_id: m,
      };
    });

    formik.setFieldValue("movies", newMovies);
  }

  async function onSave(input: GQL.SceneUpdateInput) {
    setIsLoading(true);
    try {
      const result = await updateScene({
        variables: {
          input: {
            ...input,
            rating: input.rating ?? null,
          },
        },
      });
      if (result.data?.sceneUpdate) {
        Toast.success({
          content: intl.formatMessage(
            { id: "toast.updated_entity" },
            { entity: intl.formatMessage({ id: "scene" }).toLocaleLowerCase() }
          ),
        });
        // clear the cover image so that it doesn't appear dirty
        formik.resetForm({ values: formik.values });
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  const removeStashID = (stashID: GQL.StashIdInput) => {
    formik.setFieldValue(
      "stash_ids",
      formik.values.stash_ids.filter(
        (s) =>
          !(s.endpoint === stashID.endpoint && s.stash_id === stashID.stash_id)
      )
    );
  };

  function renderTableMovies() {
    return (
      <SceneMovieTable
        movieScenes={formik.values.movies}
        onUpdate={(items) => {
          formik.setFieldValue("movies", items);
        }}
      />
    );
  }

  function onImageLoad(imageData: string) {
    setCoverImagePreview(imageData);
    formik.setFieldValue("cover_image", imageData);
  }

  function onCoverImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  async function onScrapeClicked(s: GQL.ScraperSourceInput) {
    setIsLoading(true);
    try {
      const result = await queryScrapeScene(s, scene.id);
      if (!result.data || !result.data.scrapeSingleScene?.length) {
        Toast.success({
          content: "No scenes found",
        });
        return;
      }
      // assume one returned scene
      setScrapedScene(result.data.scrapeSingleScene[0]);
      setEndpoint(s.stash_box_endpoint ?? undefined);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function scrapeFromQuery(
    s: GQL.ScraperSourceInput,
    fragment: GQL.ScrapedSceneDataFragment
  ) {
    setIsLoading(true);
    try {
      const input: GQL.ScrapedSceneInput = {
        date: fragment.date,
        details: fragment.details,
        remote_site_id: fragment.remote_site_id,
        title: fragment.title,
        url: fragment.url,
      };

      const result = await queryScrapeSceneQueryFragment(s, input);
      if (!result.data || !result.data.scrapeSingleScene?.length) {
        Toast.success({
          content: "No scenes found",
        });
        return;
      }
      // assume one returned scene
      setScrapedScene(result.data.scrapeSingleScene[0]);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function onScrapeQueryClicked(s: GQL.ScraperSourceInput) {
    setScraper(s);
    setEndpoint(s.stash_box_endpoint ?? undefined);
    setIsScraperQueryModalOpen(true);
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

  function onScrapeDialogClosed(sceneData?: GQL.ScrapedSceneDataFragment) {
    if (sceneData) {
      updateSceneFromScrapedScene(sceneData);
    }
    setScrapedScene(undefined);
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedScene) {
      return;
    }

    const currentScene = getSceneInput(formik.values);
    if (!currentScene.cover_image) {
      currentScene.cover_image = scene.paths.screenshot;
    }

    return (
      <SceneScrapeDialog
        scene={currentScene}
        scraped={scrapedScene}
        endpoint={endpoint}
        onClose={(s) => onScrapeDialogClosed(s)}
      />
    );
  }

  function renderScrapeQueryMenu() {
    const stashBoxes = stashConfig?.general.stashBoxes ?? [];

    if (stashBoxes.length === 0 && queryableScrapers.length === 0) return;

    return (
      <Dropdown title={intl.formatMessage({ id: "actions.scrape_query" })}>
        <Dropdown.Toggle variant="secondary">
          <Icon icon="search" />
        </Dropdown.Toggle>

        <Dropdown.Menu>
          {stashBoxes.map((s, index) => (
            <Dropdown.Item
              key={s.endpoint}
              onClick={() =>
                onScrapeQueryClicked({
                  stash_box_index: index,
                  stash_box_endpoint: s.endpoint,
                })
              }
            >
              {stashboxDisplayName(s.name, index)}
            </Dropdown.Item>
          ))}
          {queryableScrapers.map((s) => (
            <Dropdown.Item
              key={s.name}
              onClick={() => onScrapeQueryClicked({ scraper_id: s.id })}
            >
              {s.name}
            </Dropdown.Item>
          ))}
          <Dropdown.Item onClick={() => onReloadScrapers()}>
            <span className="fa-icon">
              <Icon icon="sync-alt" />
            </span>
            <span>
              <FormattedMessage id="actions.reload_scrapers" />
            </span>
          </Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    );
  }

  function onSceneSelected(s: GQL.ScrapedSceneDataFragment) {
    if (!scraper) return;

    if (scraper?.stash_box_index !== undefined) {
      // must be stash-box - assume full scene
      setScrapedScene(s);
    } else {
      // must be scraper
      scrapeFromQuery(scraper, s);
    }
  }

  const renderScrapeQueryModal = () => {
    if (!isScraperQueryModalOpen || !scraper) return;

    return (
      <SceneQueryModal
        scraper={scraper}
        onHide={() => setScraper(undefined)}
        onSelectScene={(s) => {
          setIsScraperQueryModalOpen(false);
          setScraper(undefined);
          onSceneSelected(s);
        }}
        name={formik.values.title || ""}
      />
    );
  };

  function renderScraperMenu() {
    const stashBoxes = stashConfig?.general.stashBoxes ?? [];

    return (
      <DropdownButton
        className="d-inline-block"
        id="scene-scrape"
        title={intl.formatMessage({ id: "actions.scrape_with" })}
      >
        {stashBoxes.map((s, index) => (
          <Dropdown.Item
            key={s.endpoint}
            onClick={() =>
              onScrapeClicked({
                stash_box_index: index,
                stash_box_endpoint: s.endpoint,
              })
            }
          >
            {stashboxDisplayName(s.name, index)}
          </Dropdown.Item>
        ))}
        {fragmentScrapers.map((s) => (
          <Dropdown.Item
            key={s.name}
            onClick={() => onScrapeClicked({ scraper_id: s.id })}
          >
            {s.name}
          </Dropdown.Item>
        ))}
        <Dropdown.Item onClick={() => onReloadScrapers()}>
          <span className="fa-icon">
            <Icon icon="sync-alt" />
          </span>
          <span>
            <FormattedMessage id="actions.reload_scrapers" />
          </span>
        </Dropdown.Item>
      </DropdownButton>
    );
  }

  function urlScrapable(scrapedUrl: string): boolean {
    return (Scrapers?.data?.listSceneScrapers ?? []).some((s) =>
      (s?.scene?.urls ?? []).some((u) => scrapedUrl.includes(u))
    );
  }

  function updateSceneFromScrapedScene(
    updatedScene: GQL.ScrapedSceneDataFragment
  ) {
    if (updatedScene.title) {
      formik.setFieldValue("title", updatedScene.title);
    }

    if (updatedScene.details) {
      formik.setFieldValue("details", updatedScene.details);
    }

    if (updatedScene.date) {
      formik.setFieldValue("date", updatedScene.date);
    }

    if (updatedScene.url) {
      formik.setFieldValue("url", updatedScene.url);
    }

    if (updatedScene.studio && updatedScene.studio.stored_id) {
      formik.setFieldValue("studio_id", updatedScene.studio.stored_id);
    }

    if (updatedScene.performers && updatedScene.performers.length > 0) {
      const idPerfs = updatedScene.performers.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idPerfs.length > 0) {
        const newIds = idPerfs.map((p) => p.stored_id);
        formik.setFieldValue("performer_ids", newIds as string[]);
      }
    }

    if (updatedScene.movies && updatedScene.movies.length > 0) {
      const idMovis = updatedScene.movies.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idMovis.length > 0) {
        const newIds = idMovis.map((p) => p.stored_id);
        setMovieIds(newIds as string[]);
      }
    }

    if (updatedScene?.tags?.length) {
      const idTags = updatedScene.tags.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idTags.length > 0) {
        const newIds = idTags.map((p) => p.stored_id);
        formik.setFieldValue("tag_ids", newIds as string[]);
      }
    }

    if (updatedScene.image) {
      // image is a base64 string
      formik.setFieldValue("cover_image", updatedScene.image);
      setCoverImagePreview(updatedScene.image);
    }

    if (updatedScene.remote_site_id && endpoint) {
      let found = false;
      formik.setFieldValue(
        "stash_ids",
        formik.values.stash_ids.map((s) => {
          if (s.endpoint === endpoint) {
            found = true;
            return {
              endpoint,
              stash_id: updatedScene.remote_site_id,
            };
          }

          return s;
        })
      );

      if (!found) {
        formik.setFieldValue(
          "stash_ids",
          formik.values.stash_ids.concat({
            endpoint,
            stash_id: updatedScene.remote_site_id,
          })
        );
      }
    }
  }

  async function onScrapeSceneURL() {
    if (!formik.values.url) {
      return;
    }
    setIsLoading(true);
    try {
      const result = await queryScrapeSceneURL(formik.values.url);
      if (!result.data || !result.data.scrapeSceneURL) {
        return;
      }
      setScrapedScene(result.data.scrapeSceneURL);
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
        </Col>
      </Form.Group>
    );
  }

  if (isLoading) return <LoadingIndicator />;

  return (
    <div id="scene-edit-details">
      <Prompt
        when={formik.dirty}
        message={intl.formatMessage({ id: "dialogs.unsaved_changes" })}
      />

      {renderScrapeQueryModal()}
      {maybeRenderScrapeDialog()}
      <Form noValidate onSubmit={formik.handleSubmit}>
        <div className="form-container edit-buttons-container row px-3 pt-3">
          <div className="edit-buttons mb-3 pl-0">
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
          <div className="ml-auto pr-3 text-right d-flex">
            <ButtonGroup className="scraper-group">
              {renderScraperMenu()}
              {renderScrapeQueryMenu()}
            </ButtonGroup>
          </div>
        </div>
        <div className="form-container row px-3">
          <div className="col-12 col-lg-7 col-xl-12">
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
                  onScrapeClick={onScrapeSceneURL}
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
            <Form.Group controlId="galleries" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "galleries" }),
                labelProps: {
                  column: true,
                  sm: 3,
                },
              })}
              <Col sm={9}>
                <GallerySelect
                  galleries={galleries}
                  onSelect={(items) => onSetGalleries(items)}
                />
              </Col>
            </Form.Group>

            <Form.Group controlId="studio" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "studio" }),
                labelProps: {
                  column: true,
                  sm: 3,
                },
              })}
              <Col sm={9}>
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

            <Form.Group controlId="moviesScenes" as={Row}>
              {FormUtils.renderLabel({
                title: `${intl.formatMessage({
                  id: "movies",
                })}/${intl.formatMessage({ id: "scenes" })}`,
                labelProps: {
                  column: true,
                  sm: 3,
                  xl: 12,
                },
              })}
              <Col sm={9} xl={12}>
                <MovieSelect
                  isMulti
                  onSelect={(items) =>
                    setMovieIds(items.map((item) => item.id))
                  }
                  ids={formik.values.movies.map((m) => m.movie_id)}
                />
                {renderTableMovies()}
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
            {formik.values.stash_ids.length ? (
              <Form.Group controlId="stashIDs">
                <Form.Label>
                  <FormattedMessage id="stash_ids" />
                </Form.Label>
                <ul className="pl-0">
                  {formik.values.stash_ids.map((stashID) => {
                    const base = stashID.endpoint.match(
                      /https?:\/\/.*?\//
                    )?.[0];
                    const link = base ? (
                      <a
                        href={`${base}scenes/${stashID.stash_id}`}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {stashID.stash_id}
                      </a>
                    ) : (
                      stashID.stash_id
                    );
                    return (
                      <li key={stashID.stash_id} className="row no-gutters">
                        <Button
                          variant="danger"
                          className="mr-2 py-0"
                          title={intl.formatMessage(
                            { id: "actions.delete_entity" },
                            {
                              entityType: intl.formatMessage({
                                id: "stash_id",
                              }),
                            }
                          )}
                          onClick={() => removeStashID(stashID)}
                        >
                          <Icon icon="trash-alt" />
                        </Button>
                        {link}
                      </li>
                    );
                  })}
                </ul>
              </Form.Group>
            ) : undefined}
          </div>
          <div className="col-12 col-lg-5 col-xl-12">
            <Form.Group controlId="details">
              <Form.Label>
                <FormattedMessage id="details" />
              </Form.Label>
              <Form.Control
                as="textarea"
                className="scene-description text-input"
                onChange={(newValue: React.ChangeEvent<HTMLTextAreaElement>) =>
                  formik.setFieldValue("details", newValue.currentTarget.value)
                }
                value={formik.values.details}
              />
            </Form.Group>
            <div>
              <Form.Group controlId="cover">
                <Form.Label>
                  <FormattedMessage id="cover_image" />
                </Form.Label>
                {imageEncoding ? (
                  <LoadingIndicator message="Encoding image..." />
                ) : (
                  <img
                    className="scene-cover"
                    src={coverImagePreview}
                    alt={intl.formatMessage({ id: "cover_image" })}
                  />
                )}
                <ImageInput
                  isEditing
                  onImageChange={onCoverImageChange}
                  onImageURL={onImageLoad}
                />
              </Form.Group>
            </div>
          </div>
        </div>
      </Form>
    </div>
  );
};
