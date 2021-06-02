import React, { useEffect, useState } from "react";
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
  queryScrapeScene,
  queryScrapeSceneURL,
  useListSceneScrapers,
  useSceneUpdate,
  mutateReloadScrapers,
  useConfiguration,
  queryStashBoxScene,
} from "src/core/StashService";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
  GallerySelect,
  Icon,
  LoadingIndicator,
  ImageInput,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { ImageUtils, FormUtils, TextUtils } from "src/utils";
import { MovieSelect } from "src/components/Shared/Select";
import { useFormik } from "formik";
import { Prompt } from "react-router";
import { SceneMovieTable } from "./SceneMovieTable";
import { RatingStars } from "./RatingStars";
import { SceneScrapeDialog } from "./SceneScrapeDialog";

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
  const Toast = useToast();
  const [galleries, setGalleries] = useState<{ id: string; title: string }[]>(
    scene.galleries.map((g) => ({
      id: g.id,
      title: g.title ?? TextUtils.fileNameFromPath(g.path ?? ""),
    }))
  );

  const Scrapers = useListSceneScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [scrapedScene, setScrapedScene] = useState<GQL.ScrapedScene | null>();

  const [coverImagePreview, setCoverImagePreview] = useState<
    string | undefined
  >(scene.paths.screenshot ?? undefined);

  const stashConfig = useConfiguration();

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

  const initialValues = {
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
    stash_ids: (scene.stash_ids ?? []).map(s => ({
        stash_id: s.stash_id,
        endpoint: s.endpoint,
    })),
  };

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSave(getSceneInput(values)),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating", v);
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
      Scrapers?.data?.listSceneScrapers ?? []
    ).filter((s) =>
      s.scene?.supported_scrapes.includes(GQL.ScrapeType.Fragment)
    );

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
        Toast.success({ content: "Updated scene" });
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

  async function onScrapeStashBoxClicked(stashBoxIndex: number) {
    setIsLoading(true);
    try {
      const result = await queryStashBoxScene(stashBoxIndex, scene.id);
      if (!result.data || !result.data.queryStashBoxScene) {
        return;
      }

      if (result.data.queryStashBoxScene.length > 0) {
        setScrapedScene(result.data.queryStashBoxScene[0]);
      } else {
        Toast.success({
          content: "No scenes found",
        });
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onScrapeClicked(scraper: GQL.Scraper) {
    setIsLoading(true);
    try {
      const result = await queryScrapeScene(
        scraper.id,
        getSceneInput(formik.values)
      );
      if (!result.data || !result.data.scrapeScene) {
        Toast.success({
          content: "No scenes found",
        });
        return;
      }
      setScrapedScene(result.data.scrapeScene);
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
        onClose={(s) => onScrapeDialogClosed(s)}
      />
    );
  }

  function renderScraperMenu() {
    const stashBoxes = stashConfig.data?.configuration.general.stashBoxes ?? [];

    // TODO - change name based on stashbox configuration
    return (
      <DropdownButton
        className="d-inline-block"
        id="scene-scrape"
        title="Scrape with..."
      >
        {stashBoxes.map((s, index) => (
          <Dropdown.Item
            key={s.endpoint}
            onClick={() => onScrapeStashBoxClicked(index)}
          >
            {s.name ?? "Stash-Box"}
          </Dropdown.Item>
        ))}
        {queryableScrapers.map((s) => (
          <Dropdown.Item key={s.name} onClick={() => onScrapeClicked(s)}>
            {s.name}
          </Dropdown.Item>
        ))}
        <Dropdown.Item onClick={() => onReloadScrapers()}>
          <span className="fa-icon">
            <Icon icon="sync-alt" />
          </span>
          <span>Reload scrapers</span>
        </Dropdown.Item>
      </DropdownButton>
    );
  }

  function maybeRenderStashboxQueryButton() {
    // const stashBoxes = stashConfig.data?.configuration.general.stashBoxes ?? [];
    // if (stashBoxes.length === 0) {
    //   return;
    // }
    // TODO - hide this button for now, with the view to add it when we get
    // the query dialog going
    // if (stashBoxes.length === 1) {
    //   return (
    //     <Button
    //       className="mr-1"
    //       onClick={() => onStashBoxQueryClicked(0)}
    //       title="Query"
    //     >
    //       <Icon className="fa-fw" icon="search" />
    //     </Button>
    //   );
    // }
    // // TODO - change name based on stashbox configuration
    // return (
    //   <Dropdown className="d-inline-block mr-1">
    //     <Dropdown.Toggle id="stashbox-query-dropdown">
    //       <Icon className="fa-fw" icon="search" />
    //     </Dropdown.Toggle>
    //     <Dropdown.Menu>
    //       {stashBoxes.map((s, index) => (
    //         <Dropdown.Item
    //           key={s.endpoint}
    //           onClick={() => onStashBoxQueryClicked(index)}
    //         >
    //           stash-box
    //         </Dropdown.Item>
    //       ))}
    //     </Dropdown.Menu>
    //   </Dropdown>
    // );
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

  function maybeRenderScrapeButton() {
    if (!formik.values.url || !urlScrapable(formik.values.url)) {
      return undefined;
    }
    return (
      <Button
        className="minimal scrape-url-button"
        onClick={onScrapeSceneURL}
        title="Scrape"
      >
        <Icon className="fa-fw" icon="file-download" />
      </Button>
    );
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
        message="Unsaved changes. Are you sure you want to leave?"
      />

      {maybeRenderScrapeDialog()}
      <Form noValidate onSubmit={formik.handleSubmit}>
        <div className="form-container row px-3 pt-3">
          <div className="col-6 edit-buttons mb-3 pl-0">
            <Button
              className="edit-button"
              variant="primary"
              disabled={!formik.dirty}
              onClick={() => formik.submitForm()}
            >
              Save
            </Button>
            <Button
              className="edit-button"
              variant="danger"
              onClick={() => onDelete()}
            >
              Delete
            </Button>
          </div>
          <Col xs={6} className="text-right">
            {maybeRenderStashboxQueryButton()}
            {renderScraperMenu()}
          </Col>
        </div>
        <div className="form-container row px-3">
          <div className="col-12 col-lg-6 col-xl-12">
            {renderTextField("title", "Title")}
            <Form.Group controlId="url" as={Row}>
              <Col xs={3} className="pr-0 url-label">
                <Form.Label className="col-form-label">URL</Form.Label>
                <div className="float-right scrape-button-container">
                  {maybeRenderScrapeButton()}
                </div>
              </Col>
              <Col xs={9}>
                <Form.Control
                  className="text-input"
                  placeholder="URL"
                  {...formik.getFieldProps("url")}
                  isInvalid={!!formik.getFieldMeta("url").error}
                />
              </Col>
            </Form.Group>
            {renderTextField("date", "Date", "YYYY-MM-DD")}
            <Form.Group controlId="rating" as={Row}>
              {FormUtils.renderLabel({
                title: "Rating",
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
                title: "Galleries",
              })}
              <Col xs={9}>
                <GallerySelect
                  galleries={galleries}
                  onSelect={(items) => setGalleries(items)}
                />
              </Col>
            </Form.Group>

            <Form.Group controlId="studio" as={Row}>
              {FormUtils.renderLabel({
                title: "Studio",
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

            <Form.Group controlId="performers" as={Row}>
              {FormUtils.renderLabel({
                title: "Performers",
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
                title: "Movies/Scenes",
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
                title: "Tags",
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
            <Form.Group controlId="details">
              <Form.Label>StashIDs</Form.Label>
              <ul className="pl-0">
                {formik.values.stash_ids.map((stashID) => {
                  const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
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
                        title="Delete StashID"
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
          </div>
          <div className="col-12 col-lg-6 col-xl-12">
            <Form.Group controlId="details">
              <Form.Label>Details</Form.Label>
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
                <Form.Label>Cover Image</Form.Label>
                {imageEncoding ? (
                  <LoadingIndicator message="Encoding image..." />
                ) : (
                  <img
                    className="scene-cover"
                    src={coverImagePreview}
                    alt="Scene cover"
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
