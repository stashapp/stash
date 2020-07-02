/* eslint-disable react/no-this-in-sfc */

import React, { useEffect, useState } from "react";
import {
  Button,
  Dropdown,
  DropdownButton,
  Form,
  Col,
  Row,
} from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryScrapeScene,
  queryScrapeSceneURL,
  useListSceneScrapers,
  useSceneUpdate,
  mutateReloadScrapers,
} from "src/core/StashService";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
  SceneGallerySelect,
  Icon,
  LoadingIndicator,
  ImageInput,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { ImageUtils, FormUtils, EditableTextUtils } from "src/utils";
import { MovieSelect } from "src/components/Shared/Select";
import { SceneMovieTable, MovieSceneIndexMap } from "./SceneMovieTable";
import { RatingStars } from "./RatingStars";
import { SceneScrapeDialog } from "./SceneScrapeDialog";

interface IProps {
  scene: GQL.SceneDataFragment;
  isVisible: boolean;
  onUpdate: (scene: GQL.SceneDataFragment) => void;
  onDelete: () => void;
}

export const SceneEditPanel: React.FC<IProps> = (props: IProps) => {
  const Toast = useToast();
  const [title, setTitle] = useState<string>();
  const [details, setDetails] = useState<string>();
  const [url, setUrl] = useState<string>();
  const [date, setDate] = useState<string>();
  const [rating, setRating] = useState<number>();
  const [galleryId, setGalleryId] = useState<string>();
  const [studioId, setStudioId] = useState<string>();
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [movieIds, setMovieIds] = useState<string[] | undefined>(undefined);
  const [movieSceneIndexes, setMovieSceneIndexes] = useState<
    MovieSceneIndexMap
  >(new Map());
  const [tagIds, setTagIds] = useState<string[]>();
  const [coverImage, setCoverImage] = useState<string>();

  const Scrapers = useListSceneScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [scrapedScene, setScrapedScene] = useState<GQL.ScrapedScene | null>();

  const [coverImagePreview, setCoverImagePreview] = useState<string>();

  // Network state
  const [isLoading, setIsLoading] = useState(true);

  const [updateScene] = useSceneUpdate(getSceneInput());

  useEffect(() => {
    if (props.isVisible) {
      Mousetrap.bind("s s", () => {
        onSave();
      });
      Mousetrap.bind("d d", () => {
        props.onDelete();
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
  }, [Scrapers]);

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
        movieSceneIndexes.forEach((v, id) => {
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

  function updateSceneEditState(state: Partial<GQL.SceneDataFragment>) {
    const perfIds = state.performers?.map((performer) => performer.id);
    const tIds = state.tags ? state.tags.map((tag) => tag.id) : undefined;
    const moviIds = state.movies
      ? state.movies.map((sceneMovie) => sceneMovie.movie.id)
      : undefined;
    const movieSceneIdx: MovieSceneIndexMap = new Map();
    if (state.movies) {
      state.movies.forEach((m) => {
        movieSceneIdx.set(m.movie.id, m.scene_index ?? undefined);
      });
    }

    setTitle(state.title ?? undefined);
    setDetails(state.details ?? undefined);
    setUrl(state.url ?? undefined);
    setDate(state.date ?? undefined);
    setRating(state.rating === null ? NaN : state.rating);
    setGalleryId(state?.gallery?.id ?? undefined);
    setStudioId(state?.studio?.id ?? undefined);
    setMovieIds(moviIds);
    setMovieSceneIndexes(movieSceneIdx);
    setPerformerIds(perfIds);
    setTagIds(tIds);
  }

  useEffect(() => {
    updateSceneEditState(props.scene);
    setCoverImagePreview(props.scene?.paths?.screenshot ?? undefined);
    setIsLoading(false);
  }, [props.scene]);

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, true);

  function getSceneInput(): GQL.SceneUpdateInput {
    return {
      id: props.scene.id,
      title,
      details,
      url,
      date,
      rating,
      gallery_id: galleryId,
      studio_id: studioId,
      performer_ids: performerIds,
      movies: makeMovieInputs(),
      tag_ids: tagIds,
      cover_image: coverImage,
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
      const result = await updateScene();
      if (result.data?.sceneUpdate) {
        props.onUpdate(result.data.sceneUpdate);
        Toast.success({ content: "Updated scene" });
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

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

  async function onScrapeClicked(scraper: GQL.Scraper) {
    setIsLoading(true);
    try {
      const result = await queryScrapeScene(scraper.id, getSceneInput());
      if (!result.data || !result.data.scrapeScene) {
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

  function onScrapeDialogClosed(scene?: GQL.ScrapedSceneDataFragment) {
    if (scene) {
      updateSceneFromScrapedScene(scene);
    }
    setScrapedScene(undefined);
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedScene) {
      return;
    }

    const currentScene = getSceneInput();
    if (!currentScene.cover_image) {
      currentScene.cover_image = props.scene.paths.screenshot;
    }

    return (
      <SceneScrapeDialog
        scene={currentScene}
        scraped={scrapedScene}
        onClose={(scene) => {
          onScrapeDialogClosed(scene);
        }}
      />
    );
  }

  function renderScraperMenu() {
    return (
      <DropdownButton id="scene-scrape" title="Scrape with...">
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

  function urlScrapable(scrapedUrl: string): boolean {
    return (Scrapers?.data?.listSceneScrapers ?? []).some((s) =>
      (s?.scene?.urls ?? []).some((u) => scrapedUrl.includes(u))
    );
  }

  function updateSceneFromScrapedScene(scene: GQL.ScrapedSceneDataFragment) {
    if (scene.title) {
      setTitle(scene.title);
    }

    if (scene.details) {
      setDetails(scene.details);
    }

    if (scene.date) {
      setDate(scene.date);
    }

    if (scene.url) {
      setUrl(scene.url);
    }

    if (scene.studio && scene.studio.id) {
      setStudioId(scene.studio.id);
    }

    if (scene.performers && scene.performers.length > 0) {
      const idPerfs = scene.performers.filter((p) => {
        return p.id !== undefined && p.id !== null;
      });

      if (idPerfs.length > 0) {
        const newIds = idPerfs.map((p) => p.id);
        setPerformerIds(newIds as string[]);
      }
    }

    if (scene.movies && scene.movies.length > 0) {
      const idMovis = scene.movies.filter((p) => {
        return p.id !== undefined && p.id !== null;
      });

      if (idMovis.length > 0) {
        const newIds = idMovis.map((p) => p.id);
        setMovieIds(newIds as string[]);
      }
    }

    if (scene?.tags?.length) {
      const idTags = scene.tags.filter((p) => {
        return p.id !== undefined && p.id !== null;
      });

      if (idTags.length > 0) {
        const newIds = idTags.map((p) => p.id);
        setTagIds(newIds as string[]);
      }
    }

    if (scene.image) {
      // image is a base64 string
      setCoverImage(scene.image);
      setCoverImagePreview(scene.image);
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
        <div className="col edit-buttons mb-3 pl-0">
          <Button className="edit-button" variant="primary" onClick={onSave}>
            Save
          </Button>
          <Button
            className="edit-button"
            variant="danger"
            onClick={() => props.onDelete()}
          >
            Delete
          </Button>
        </div>
        {renderScraperMenu()}
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
          <Form.Group controlId="gallery" as={Row}>
            {FormUtils.renderLabel({
              title: "Gallery",
            })}
            <Col xs={9}>
              <SceneGallerySelect
                sceneId={props.scene.id}
                initialId={galleryId}
                onSelect={(item) => setGalleryId(item ? item.id : undefined)}
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
              <ImageInput isEditing onImageChange={onCoverImageChange} />
            </Form.Group>
          </div>
        </div>
      </div>
    </div>
  );
};
