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
import { ImageUtils, FormUtils, EditableTextUtils, TextUtils } from "src/utils";
import { MovieSelect } from "src/components/Shared/Select";
import { SceneMovieTable, MovieSceneIndexMap } from "./SceneMovieTable";
import { RatingStars } from "./RatingStars";
import { SceneScrapeDialog } from "./SceneScrapeDialog";

interface IProps {
  scene: GQL.SceneDataFragment;
  isVisible: boolean;
  onDelete: () => void;
}

export const SceneEditPanel: React.FC<IProps> = ({
  scene,
  isVisible,
  onDelete,
}) => {
  const Toast = useToast();
  const [title, setTitle] = useState<string>(scene.title ?? "");
  const [details, setDetails] = useState<string>(scene.details ?? "");
  const [url, setUrl] = useState<string>(scene.url ?? "");
  const [date, setDate] = useState<string>(scene.date ?? "");
  const [rating, setRating] = useState<number | undefined>(
    scene.rating ?? undefined
  );
  const [galleries, setGalleries] = useState<{ id: string; title: string }[]>(
    scene.galleries.map((g) => ({
      id: g.id,
      title: g.title ?? TextUtils.fileNameFromPath(g.path ?? ""),
    }))
  );
  const [studioId, setStudioId] = useState<string | undefined>(
    scene.studio?.id
  );
  const [performerIds, setPerformerIds] = useState<string[]>(
    scene.performers.map((p) => p.id)
  );
  const [movieIds, setMovieIds] = useState<string[]>(
    scene.movies.map((m) => m.movie.id)
  );
  const [
    movieSceneIndexes,
    setMovieSceneIndexes,
  ] = useState<MovieSceneIndexMap>(
    new Map(scene.movies.map((m) => [m.movie.id, m.scene_index ?? undefined]))
  );
  const [tagIds, setTagIds] = useState<string[]>(scene.tags.map((t) => t.id));
  const [coverImage, setCoverImage] = useState<string>();
  const [stashIDs, setStashIDs] = useState<GQL.StashIdInput[]>(scene.stash_ids);

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

  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        onSave();
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

  useEffect(() => {
    let changed = false;
    const newMap: MovieSceneIndexMap = new Map();
    if (movieIds) {
      movieIds.forEach((id) => {
        if (!movieSceneIndexes.has(id)) {
          changed = true;
          newMap.set(id, undefined);
        } else {
          newMap.set(id, movieSceneIndexes.get(id));
        }
      });

      if (!changed) {
        movieSceneIndexes.forEach((_v, id) => {
          if (!newMap.has(id)) {
            // id was removed
            changed = true;
          }
        });
      }

      if (changed) {
        setMovieSceneIndexes(newMap);
      }
    }
  }, [movieIds, movieSceneIndexes]);

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, true);

  function getSceneInput(): GQL.SceneUpdateInput {
    return {
      id: scene.id,
      title,
      details,
      url,
      date,
      rating: rating ?? null,
      gallery_ids: galleries.map((g) => g.id),
      studio_id: studioId ?? null,
      performer_ids: performerIds,
      movies: makeMovieInputs(),
      tag_ids: tagIds,
      cover_image: coverImage,
      stash_ids: stashIDs.map((s) => ({
        stash_id: s.stash_id,
        endpoint: s.endpoint,
      })),
    };
  }

  function makeMovieInputs(): GQL.SceneMovieInput[] | undefined {
    if (!movieIds) {
      return undefined;
    }

    let ret = movieIds.map((id) => {
      const r: GQL.SceneMovieInput = {
        movie_id: id,
      };
      return r;
    });

    ret = ret.map((r) => {
      return { scene_index: movieSceneIndexes.get(r.movie_id), ...r };
    });

    return ret;
  }

  async function onSave() {
    setIsLoading(true);
    try {
      const result = await updateScene({
        variables: {
          input: getSceneInput(),
        },
      });
      if (result.data?.sceneUpdate) {
        Toast.success({ content: "Updated scene" });
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  const removeStashID = (stashID: GQL.StashIdInput) => {
    setStashIDs(
      stashIDs.filter(
        (s) =>
          !(s.endpoint === stashID.endpoint && s.stash_id === stashID.stash_id)
      )
    );
  };

  function renderTableMovies() {
    return (
      <SceneMovieTable
        movieSceneIndexes={movieSceneIndexes}
        onUpdate={(items) => {
          setMovieSceneIndexes(items);
        }}
      />
    );
  }

  function onImageLoad(imageData: string) {
    setCoverImagePreview(imageData);
    setCoverImage(imageData);
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
      const result = await queryScrapeScene(scraper.id, getSceneInput());
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

    const currentScene = getSceneInput();
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
            { s.name ?? "Stash-Box" }
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
      setTitle(updatedScene.title);
    }

    if (updatedScene.details) {
      setDetails(updatedScene.details);
    }

    if (updatedScene.date) {
      setDate(updatedScene.date);
    }

    if (updatedScene.url) {
      setUrl(updatedScene.url);
    }

    if (updatedScene.studio && updatedScene.studio.stored_id) {
      setStudioId(updatedScene.studio.stored_id);
    }

    if (updatedScene.performers && updatedScene.performers.length > 0) {
      const idPerfs = updatedScene.performers.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idPerfs.length > 0) {
        const newIds = idPerfs.map((p) => p.stored_id);
        setPerformerIds(newIds as string[]);
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
        setTagIds(newIds as string[]);
      }
    }

    if (updatedScene.image) {
      // image is a base64 string
      setCoverImage(updatedScene.image);
      setCoverImagePreview(updatedScene.image);
    }
  }

  async function onScrapeSceneURL() {
    if (!url) {
      return;
    }
    setIsLoading(true);
    try {
      const result = await queryScrapeSceneURL(url);
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
    if (!url || !urlScrapable(url)) {
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

  if (isLoading) return <LoadingIndicator />;

  return (
    <div id="scene-edit-details">
      {maybeRenderScrapeDialog()}
      <div className="form-container row px-3 pt-3">
        <div className="col-6 edit-buttons mb-3 pl-0">
          <Button className="edit-button" variant="primary" onClick={onSave}>
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
          {FormUtils.renderInputGroup({
            title: "Title",
            value: title,
            onChange: setTitle,
            isEditing: true,
          })}
          <Form.Group controlId="url" as={Row}>
            <Col xs={3} className="pr-0 url-label">
              <Form.Label className="col-form-label">URL</Form.Label>
              <div className="float-right scrape-button-container">
                {maybeRenderScrapeButton()}
              </div>
            </Col>
            <Col xs={9}>
              {EditableTextUtils.renderInputGroup({
                title: "URL",
                value: url,
                onChange: setUrl,
                isEditing: true,
              })}
            </Col>
          </Form.Group>
          {FormUtils.renderInputGroup({
            title: "Date",
            value: date,
            isEditing: true,
            onChange: setDate,
            placeholder: "YYYY-MM-DD",
          })}
          <Form.Group controlId="rating" as={Row}>
            {FormUtils.renderLabel({
              title: "Rating",
            })}
            <Col xs={9}>
              <RatingStars
                value={rating}
                onSetRating={(value) => setRating(value)}
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
                  setStudioId(items.length > 0 ? items[0]?.id : undefined)
                }
                ids={studioId ? [studioId] : []}
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
                  setPerformerIds(items.map((item) => item.id))
                }
                ids={performerIds}
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
                onSelect={(items) => setMovieIds(items.map((item) => item.id))}
                ids={movieIds}
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
                onSelect={(items) => setTagIds(items.map((item) => item.id))}
                ids={tagIds}
              />
            </Col>
          </Form.Group>
          <Form.Group controlId="details">
            <Form.Label>StashIDs</Form.Label>
            <ul className="pl-0">
              {stashIDs.map((stashID) => {
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
                setDetails(newValue.currentTarget.value)
              }
              value={details}
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
    </div>
  );
};
