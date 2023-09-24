import React, { ChangeEvent, useEffect, useMemo, useState } from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  DropdownButton,
  Form,
  Col,
  Row,
} from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import isEqual from "lodash-es/isEqual";
import { useBulkSceneUpdate } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import { StudioSelect } from "../Shared/Select";
import { ModalComponent } from "../Shared/Modal";
import { ImageInput } from "../Shared/ImageInput";
import { URLListInput } from "../Shared/URLField";
import { MultiSet } from "../Shared/MultiSet";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import { stashboxDisplayName } from "src/utils/stashbox";
import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import ImageUtils from "src/utils/image";
import FormUtils from "src/utils/form";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import {
  getAggregateInputIDs,
  getAggregateInputValue,
  getAggregateMovieIds,
  getAggregatePerformerIds,
  getAggregateRating,
  getAggregateStudioId,
  getAggregateTagIds,
  getAggregateGalleries,
} from "src/utils/bulkUpdate";
import {
  queryScrapeScene,
  queryScrapeSceneURL,
  useListSceneScrapers,
  mutateReloadScrapers,
  queryScrapeSceneQueryFragment,
} from "src/core/StashService";
import {
  faPencilAlt,
  faSearch,
  faSyncAlt,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { objectTitle } from "src/core/files";

import { lazyComponent } from "src/utils/lazyComponent";

const SceneScrapeDialog = lazyComponent(
  () => import("src/components/Scenes/SceneDetails/SceneScrapeDialog")
);
const SceneQueryModal = lazyComponent(
  () => import("src/components/Scenes/SceneDetails/SceneQueryModal")
);

interface IListOperationProps {
  selected: GQL.SceneDataFragment[];
  onClose: (applied: boolean) => void;
}

export const EditScenesDialog: React.FC<IListOperationProps> = (
  props: IListOperationProps
) => {
  const intl = useIntl();
  const Toast = useToast();

  // Single scene state
  const [isLoading, setIsLoading] = useState(false);
  const { configuration: stashConfig } = React.useContext(ConfigurationContext);
  const Scrapers = useListSceneScrapers();
  const [scrapedScene, setScrapedScene] = useState<GQL.ScrapedScene | null>();
  const [fragmentScrapers, setFragmentScrapers] = useState<GQL.Scraper[]>([]);
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);
  const [scraper, setScraper] = useState<GQL.ScraperSourceInput>();
  const [isScraperQueryModalOpen, setIsScraperQueryModalOpen] =
    useState<boolean>(false);
  const [endpoint, setEndpoint] = useState<string>();
  const isSingleScene = props.selected && props.selected.length === 1;
  const [title, setTitle] = useState<string>("");
  const [studioCode, setStudioCode] = useState<string>("");
  const [date, setDate] = useState<string>("");
  const [details, setDetails] = useState<string>("");
  const [director, setDirector] = useState<string>("");
  const [galleryMode, setGalleryMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [galleryIds, setGalleryIds] = useState<string[]>();
  const [existingGalleryIds, setExistingGalleryIds] = useState<string[]>();
  const [coverImage, setCoverImage] = React.useState<string | undefined>();
  const [stashIds, setStashIds] = React.useState<
    { endpoint: string; stash_id: string | null | undefined }[]
  >([]);

  const [urls, setUrls] = React.useState<GQL.BulkUpdateStrings>({
    mode: GQL.BulkUpdateIdMode.Set,
    values: [],
  });
  const [urlsErrorMsg, setUrlsErrorMsg] = useState<string | undefined>(
    undefined
  );
  const [urlsErrorIdx, setUrlsErrorIdx] = useState<number[]>([]);
  type ValidationErrors = {
    [key: string]: string | undefined;
  };
  const [validationErrors, setValidationErrors] = useState<ValidationErrors>(
    {}
  );

  // > 1 scene state
  const [rating100, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [performerMode, setPerformerMode] =
    React.useState<GQL.BulkUpdateIdMode>(GQL.BulkUpdateIdMode.Add);
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [existingPerformerIds, setExistingPerformerIds] = useState<string[]>();
  const [tagMode, setTagMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [tagIds, setTagIds] = useState<string[]>();
  const [existingTagIds, setExistingTagIds] = useState<string[]>();
  const [movieMode, setMovieMode] = React.useState<GQL.BulkUpdateIdMode>(
    GQL.BulkUpdateIdMode.Add
  );
  const [movieIds, setMovieIds] = useState<string[]>();
  const [existingMovieIds, setExistingMovieIds] = useState<string[]>();
  const [organized, setOrganized] = useState<boolean | undefined>();

  const [updateScenes] = useBulkSceneUpdate(getSceneInput());

  // Network state
  const [isUpdating, setIsUpdating] = useState(false);

  const checkboxRef = React.createRef<HTMLInputElement>();

  const [formValues, setFormValues] = useState(getSceneInput());

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

  const schema = yup.object({
    title: yup.string().ensure(),
    code: yup.string().ensure(),
    urls: yup
      .array(yup.string().required())
      .defined()
      .test({
        name: "unique",
        test: (value) => {
          const dupes = value
            .map((e, i, a) => {
              if (a.indexOf(e) !== i) {
                return String(i - 1);
              } else {
                return null;
              }
            })
            .filter((e) => e !== null) as string[];
          if (dupes.length === 0) return true;
          return new yup.ValidationError(dupes.join(" "), value, "urls");
        },
      }),
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
    director: yup.string().ensure(),
    rating100: yup.number().nullable().defined(),
    gallery_ids: yup.array(yup.string().required()).defined(),
    studio_id: yup.string().required().nullable(),
    performer_ids: yup.array(yup.string().required()).defined(),
    movies: yup
      .array(
        yup.object({
          movie_id: yup.string().required(),
          scene_index: yup.number().nullable().defined(),
        })
      )
      .defined(),
    tag_ids: yup.array(yup.string().required()).defined(),
    stash_ids: yup.mixed<GQL.StashIdInput[]>().defined(),
    details: yup.string().ensure(),
    cover_image: yup.string().nullable().optional(),
  });

  async function onScrapeClicked(s: GQL.ScraperSourceInput) {
    setIsLoading(true);
    try {
      const result = await queryScrapeScene(s, props.selected[0].id!);
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
        code: fragment.code,
        details: fragment.details,
        director: fragment.director,
        remote_site_id: fragment.remote_site_id,
        title: fragment.title,
        urls: fragment.urls,
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

    type CurrentScene = {
      id?: string;
      title?: string;
      code?: string;
      date?: string;
      details?: string;
      director?: string;
      cover_image?: string | null | undefined;
      gallery_ids?: string[];
      urls?: string[];
      rating100?: number;
      studio_id?: string;
      performer_ids?: string[];
      tag_ids?: string[];
      movie_ids?: string[];
      organized?: boolean;
    };

    let currentScene: CurrentScene = {};

    if (props.selected && props.selected.length === 1) {
      currentScene = {
        id: props.selected[0].id!,
        title: title,
        code: studioCode,
        date: date,
        details: details,
        director: director,
        cover_image: coverImage,
        gallery_ids: galleryIds,
        urls: Array.isArray(urls) ? urls : [],
        rating100: rating100,
        studio_id: studioId,
        performer_ids: performerIds,
        tag_ids: tagIds,
        movie_ids: movieIds,
        organized: organized,
      };
    }

    if (!currentScene.cover_image) {
      currentScene.cover_image = props.selected[0].paths?.screenshot;
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
          <Icon icon={faSearch} />
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
              <Icon icon={faSyncAlt} />
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

    if (scraper?.stash_box_endpoint !== undefined) {
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
        name={title || objectTitle(props.selected[0]) || ""}
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

    if (updatedScene.code) {
      setStudioCode(updatedScene.code);
    }

    if (updatedScene.details) {
      setDetails(updatedScene.details);
    }

    if (updatedScene.director) {
      setDirector(updatedScene.director);
    }

    if (updatedScene.date) {
      setDate(updatedScene.date);
    }

    if (updatedScene.urls) {
      setUrls({
        mode: GQL.BulkUpdateIdMode.Set,
        values: updatedScene.urls,
      });
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
    }

    if (updatedScene.remote_site_id && endpoint) {
      let found = false;
      setStashIds((stashIds) =>
        stashIds.map((s) => {
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
        setStashIds((stashIds) => [
          ...stashIds,
          { endpoint, stash_id: updatedScene.remote_site_id },
        ]);
      }
    }
  }

  async function onScrapeSceneURL(url: string) {
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

  function getSceneInput(): GQL.BulkSceneUpdateInput {
    // need to determine what we are actually setting on each scene
    const aggregateRating = getAggregateRating(props.selected);
    const aggregateStudioId = getAggregateStudioId(props.selected);
    const aggregatePerformerIds = getAggregatePerformerIds(props.selected);
    const aggregateTagIds = getAggregateTagIds(props.selected);
    const aggregateMovieIds = getAggregateMovieIds(props.selected);
    const aggregateGalleryIds = getAggregateGalleries(props.selected);

    const newCode = studioCode ? studioCode : undefined;
    const newTitle = title ? title : undefined;
    const newDate = date ? date : undefined;
    const newDirector = director ? director : undefined;
    const newDetails = details ? details : undefined;
    const newUrls = urls ? urls : undefined;
    const newCoverImage = coverImage ? coverImage : undefined;
    const newStashIds: GQL.StashIdInput[] = stashIds.map((stash) => ({
      endpoint: stash.endpoint,
      stash_id: stash.stash_id || "",
    }));

    const sceneInput: GQL.BulkSceneUpdateInput = {
      ids: props.selected.map((scene) => {
        return scene.id;
      }),
    };

    sceneInput.rating100 = getAggregateInputValue(rating100, aggregateRating);
    sceneInput.studio_id = getAggregateInputValue(studioId, aggregateStudioId);

    sceneInput.performer_ids = getAggregateInputIDs(
      performerMode,
      performerIds,
      aggregatePerformerIds
    );
    sceneInput.tag_ids = getAggregateInputIDs(tagMode, tagIds, aggregateTagIds);
    sceneInput.movie_ids = getAggregateInputIDs(
      movieMode,
      movieIds,
      aggregateMovieIds
    );

    if (organized !== undefined) {
      sceneInput.organized = organized;
    }

    // need to determine what we are actually setting if single scene
    if (props?.selected?.length === 1) {
      sceneInput.title = newTitle || "";
      sceneInput.code = newCode || "";
      sceneInput.date = newDate || "";
      sceneInput.director = newDirector || "";
      sceneInput.gallery_ids = getAggregateInputIDs(
        galleryMode,
        galleryIds,
        aggregateGalleryIds
      );
      sceneInput.details = newDetails || "";
      sceneInput.urls = newUrls;
      sceneInput.cover_image = newCoverImage || "";
      sceneInput.stash_ids = newStashIds || "";
    }
    return sceneInput;
  }

  async function onSave() {
    setIsUpdating(true);
    try {
      await updateScenes();
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          { entity: intl.formatMessage({ id: "scenes" }).toLocaleLowerCase() }
        ),
      });
      props.onClose(true);
    } catch (e) {
      Toast.error(e);
    }
    setIsUpdating(false);
  }

  useEffect(() => {
    const state = props.selected;
    let updateRating: number | undefined;
    let updateStudioID: string | undefined;
    let updatePerformerIds: string[] = [];
    let updateTagIds: string[] = [];
    let updateMovieIds: string[] = [];
    let updateOrganized: boolean | undefined;
    let first = true;

    // For single scene state
    let updateGalleryIds: string[] = [];
    let updateTitle: string | undefined;
    let updateStudioCode: string | undefined;
    let updateUrls: string[] = [];
    let updateDate: string | undefined;
    let updateDirector: string | undefined;
    let updateDetails: string | undefined;
    let updateCoverImage: string | undefined;
    let updatestashIds: { stash_id: string; endpoint: string }[] = [];

    state.forEach((scene: GQL.SceneDataFragment) => {
      if (state.length === 1) {
        const sceneTitle = scene.title;
        const sceneStudioCode = scene.code;
        const sceneurls = (scene?.urls ?? []).map((u) => u).sort();
        const sceneDate = scene?.date;
        const sceneDirector = scene?.director;
        const sceneGalleries = (scene?.galleries).map((g) => g.id).sort();
        const sceneDetails = scene?.details;
        const sceneCoverImage = scene?.paths?.screenshot;
        const sceneStashIds = (scene?.stash_ids ?? [])
          .map((s) => ({
            stash_id: s.stash_id,
            endpoint: s.endpoint,
          }))
          .sort();

        updateGalleryIds = sceneGalleries;
        updateTitle = sceneTitle ?? undefined;
        updateStudioCode = sceneStudioCode ?? undefined;
        updateUrls = sceneurls;
        updateDate = sceneDate ?? undefined;
        updateDirector = sceneDirector ?? undefined;
        updateDetails = sceneDetails ?? undefined;
        updateCoverImage = sceneCoverImage ?? undefined;
        updatestashIds = sceneStashIds;
      }

      const sceneRating = scene.rating100;
      const sceneStudioID = scene?.studio?.id;
      const scenePerformerIDs = (scene.performers ?? [])
        .map((p) => p.id)
        .sort();
      const sceneTagIDs = (scene.tags ?? []).map((p) => p.id).sort();
      const sceneMovieIDs = (scene.movies ?? []).map((m) => m.movie.id).sort();

      if (first) {
        updateRating = sceneRating ?? undefined;
        updateStudioID = sceneStudioID;
        updatePerformerIds = scenePerformerIDs;
        updateTagIds = sceneTagIDs;
        updateMovieIds = sceneMovieIDs;
        first = false;
        updateOrganized = scene.organized;
      } else {
        if (sceneRating !== updateRating) {
          updateRating = undefined;
        }
        if (sceneStudioID !== updateStudioID) {
          updateStudioID = undefined;
        }
        if (!isEqual(scenePerformerIDs, updatePerformerIds)) {
          updatePerformerIds = [];
        }
        if (!isEqual(sceneTagIDs, updateTagIds)) {
          updateTagIds = [];
        }
        if (!isEqual(sceneMovieIDs, updateMovieIds)) {
          updateMovieIds = [];
        }
        if (scene.organized !== updateOrganized) {
          updateOrganized = undefined;
        }
      }
    });

    if (state.length === 1) {
      setTitle(updateTitle || "");
      setStudioCode(updateStudioCode || "");
      setUrls({
        mode: GQL.BulkUpdateIdMode.Add,
        values: updateUrls,
      });
      setDate(updateDate || "");
      setDirector(updateDirector || "");
      setDetails(updateDetails || "");
      setExistingGalleryIds(updateGalleryIds);
      setCoverImage(updateCoverImage);
    }

    setRating(updateRating);
    setStudioId(updateStudioID);
    setExistingPerformerIds(updatePerformerIds);
    setExistingTagIds(updateTagIds);
    setExistingMovieIds(updateMovieIds);
    setOrganized(updateOrganized);
    setStashIds(updatestashIds);
  }, [props.selected, galleryMode, performerMode, tagMode, movieMode]);

  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = organized === undefined;
    }
  }, [organized, checkboxRef]);

  const encodingImage = ImageUtils.usePasteImage(onImageLoad);

  const coverImagePreview = useMemo(() => {
    if (isSingleScene) {
      const sceneImage = props.selected[0].paths?.screenshot;
      const formImage = coverImage;

      if (formImage === null && sceneImage) {
        const sceneImageURL = new URL(sceneImage);
        sceneImageURL.searchParams.set("default", "true");
        return sceneImageURL.toString();
      } else if (formImage) {
        return formImage;
      }
      return sceneImage;
    }

    return null;
  }, [coverImage]);

  function onImageLoad(imageData: string) {
    setCoverImage(imageData);
  }

  function onCoverImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, (imageData) => {
      onImageLoad(imageData);
    });
  }

  const image = useMemo(() => {
    if (isSingleScene) {
      if (encodingImage) {
        return (
          <LoadingIndicator
            message={`${intl.formatMessage({ id: "encoding_image" })}...`}
          />
        );
      }

      if (coverImagePreview) {
        return (
          <img
            className="scene-cover"
            src={coverImagePreview}
            alt={intl.formatMessage({ id: "cover_image" })}
          />
        );
      }
    }

    // Return null or a default component when not a single scene is selected
    return null;
  }, [encodingImage, coverImagePreview, intl]);

  function renderMultiSelect(
    type: "galleries" | "performers" | "tags" | "movies",
    ids: string[] | undefined
  ) {
    let mode = GQL.BulkUpdateIdMode.Add;
    let existingIds: string[] | undefined = [];
    switch (type) {
      case "galleries":
        mode = galleryMode;
        existingIds = existingGalleryIds;
        break;
      case "performers":
        mode = performerMode;
        existingIds = existingPerformerIds;
        break;
      case "tags":
        mode = tagMode;
        existingIds = existingTagIds;
        break;
      case "movies":
        mode = movieMode;
        existingIds = existingMovieIds;
        break;
    }

    return (
      <MultiSet
        type={type}
        disabled={isUpdating}
        onUpdate={(itemIDs) => {
          switch (type) {
            case "galleries":
              setGalleryIds(itemIDs);
              break;
            case "performers":
              setPerformerIds(itemIDs);
              break;
            case "tags":
              setTagIds(itemIDs);
              break;
            case "movies":
              setMovieIds(itemIDs);
              break;
          }
        }}
        onSetMode={(newMode) => {
          switch (type) {
            case "galleries":
              setGalleryMode(newMode);
              break;
            case "performers":
              setPerformerMode(newMode);
              break;
            case "tags":
              setTagMode(newMode);
              break;
            case "movies":
              setMovieMode(newMode);
              break;
          }
        }}
        ids={ids ?? []}
        existingIds={existingIds ?? []}
        mode={mode}
      />
    );
  }

  function cycleOrganized() {
    if (organized) {
      setOrganized(undefined);
    } else if (organized === undefined) {
      setOrganized(false);
    } else {
      setOrganized(true);
    }
  }

  const removeStashID = (stashID: GQL.StashIdInput) => {
    const updatedStashIds = stashIds.filter(
      (s) =>
        !(s.endpoint === stashID.endpoint && s.stash_id === stashID.stash_id)
    );
    setStashIds(updatedStashIds);
  };

  type updateVariableType = "title" | "code" | "date" | "director";

  function updateFormValuesVars(type: updateVariableType, value: string) {
    switch (type) {
      case "title":
        setTitle(value);
        break;
      case "code":
        setStudioCode(value);
        break;
      case "date":
        setDate(value);
        break;
      case "director":
        setDirector(value);
        break;
    }
  }

  const variableMap: Record<updateVariableType, any> = {
    title: title,
    code: studioCode,
    date: date,
    director: director,
  };

  async function handleUrlInputChange(value: string[]) {
    // Update the URL input values
    setUrls({
      mode: GQL.BulkUpdateIdMode.Set,
      values: value,
    });

    try {
      await schema.validateAt("urls", { ...formValues, ["urls"]: value });
      setUrlsErrorMsg(undefined);
      setUrlsErrorIdx([]);
      setValidationErrors((prevErrors) => ({
        ...prevErrors,
        ["urls"]: undefined,
      }));
    } catch (error: any) {
      setUrlsErrorMsg(
        intl.formatMessage({ id: "validation.urls_must_be_unique" })
      );
      const errorIndices = error.message
        .split(" ")
        .map((e: string) => parseInt(e));
      setUrlsErrorIdx(errorIndices);
      setValidationErrors((prevErrors) => ({
        ...prevErrors,
        ["urls"]: error.message,
      }));
    }
  }

  function renderTextField(
    field: updateVariableType,
    title: string,
    placeholder?: string
  ) {
    const handleChanges = async (e: ChangeEvent<HTMLInputElement>) => {
      const { name, value } = e.target;
      updateFormValuesVars(field, value);

      setFormValues((prevFormValues) => ({
        ...prevFormValues,
        [name]: "value",
      }));

      try {
        await schema.validateAt(name, { ...formValues, [name]: value });
        setValidationErrors((prevErrors) => ({
          ...prevErrors,
          [name]: undefined,
        }));
      } catch (error: any) {
        setValidationErrors((prevErrors) => ({
          ...prevErrors,
          [name]: error.message,
        }));
      }
    };
    return (
      <Form.Group controlId={field} as={Row}>
        {FormUtils.renderLabel({
          title,
        })}
        <Col xs={9}>
          <Form.Control
            className="text-input"
            name={field}
            placeholder={placeholder ?? title}
            onChange={handleChanges}
            isInvalid={!!validationErrors[field]}
            value={variableMap[field] as string}
          />
          <Form.Control.Feedback type="invalid">
            {validationErrors[field]}
          </Form.Control.Feedback>
        </Col>
      </Form.Group>
    );
  }

  function renderAdditional(mode: number) {
    switch (mode) {
      case 0:
        return (
          <>
            {
              <div className="scrape-div">
                <ButtonGroup className="scraper-group">
                  {renderScraperMenu()}
                  {renderScrapeQueryMenu()}
                </ButtonGroup>
              </div>
            }
          </>
        );
      case 1:
        return renderTextField("title", intl.formatMessage({ id: "title" }));
      case 2:
        return renderTextField(
          "code",
          intl.formatMessage({ id: "scene_code" })
        );
      case 3:
        return (
          <Form.Group controlId="urls" as={Row}>
            <Col xs={3} className="pr-0 url-label">
              <Form.Label className="col-form-label">
                <FormattedMessage id="urls" />
              </Form.Label>
            </Col>
            <Col xs={9}>
              <URLListInput
                value={urls?.values || []}
                setValue={(value) => handleUrlInputChange(value)}
                errors={urlsErrorMsg}
                errorIdx={urlsErrorIdx}
                onScrapeClick={(url) => onScrapeSceneURL(url)}
                urlScrapable={urlScrapable}
              />
            </Col>
          </Form.Group>
        );
      case 4:
        return renderTextField("date", intl.formatMessage({ id: "date" }));
      case 5:
        return renderTextField(
          "director",
          intl.formatMessage({ id: "director" })
        );
      case 6:
        return (
          <Form.Group controlId="galleries">
            <Form.Label>
              <FormattedMessage id="galleries" />
            </Form.Label>
            {renderMultiSelect("galleries", galleryIds)}
          </Form.Group>
        );
      case 7:
        return (
          <>
            {props.selected[0].stash_ids.length > 0 && (
              <Form.Group controlId="stashIDs">
                <Form.Label>
                  <FormattedMessage id="stash_ids" />
                </Form.Label>
                <ul className="pl-0">
                  {props.selected[0].stash_ids.map((stashID) => {
                    const base =
                      stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
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
                          <Icon icon={faTrashAlt} />
                        </Button>
                        {link}
                      </li>
                    );
                  })}
                </ul>
              </Form.Group>
            )}
          </>
        );
      case 8:
        return (
          <Form.Group controlId="details">
            <Form.Label>
              <FormattedMessage id="details" />
            </Form.Label>
            <Form.Control
              as="textarea"
              className="scene-description text-input"
              onChange={(e) => setDetails(e.currentTarget.value)}
              value={details as string}
            />
          </Form.Group>
        );
      case 9:
        return (
          <Form.Group controlId="cover">
            <Form.Label>
              <FormattedMessage id="cover_image" />
            </Form.Label>
            {image}
            <ImageInput
              isEditing
              onImageChange={onCoverImageChange}
              onImageURL={onImageLoad}
            />
          </Form.Group>
        );
    }
  }

  function render() {
    const additionalComponents = [];

    if (isSingleScene) {
      additionalComponents.push([
        renderScrapeQueryModal(),
        maybeRenderScrapeDialog(),
      ]);
      additionalComponents.push([
        renderAdditional(0),
        renderAdditional(1),
        renderAdditional(2),
        renderAdditional(3),
        renderAdditional(4),
        renderAdditional(5),
      ]);
      additionalComponents.push(renderAdditional(6));
      additionalComponents.push([
        renderAdditional(7),
        renderAdditional(8),
        renderAdditional(9),
      ]);
    }

    return (
      <ModalComponent
        show
        icon={faPencilAlt}
        header={intl.formatMessage(
          { id: "dialogs.edit_entity_title" },
          {
            count: props?.selected?.length ?? 1,
            singularEntity: intl.formatMessage({ id: "scene" }),
            pluralEntity: `${intl.formatMessage({ id: "scenes" })} ${
              " (" + props?.selected?.length + ")"
            }`,
          }
        )}
        accept={{
          onClick: onSave,
          text: intl.formatMessage({ id: "actions.apply" }),
        }}
        cancel={{
          onClick: () => props.onClose(false),
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "secondary",
        }}
        isRunning={isUpdating}
        disabled={Object.values(validationErrors).some((error) => !!error)}
        modalProps={{
          scrollable: props.selected.length === 1,
          dialogClassName: "SceneEditModal",
        }}
      >
        {isLoading && <LoadingIndicator />}

        {isSingleScene && additionalComponents[0]}

        <div style={{ display: isLoading ? "none" : "block" }}>
          <Form>
            {isSingleScene && additionalComponents[1]}

            <Form.Group controlId="rating" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "rating" }),
              })}
              <Col xs={9}>
                <RatingSystem
                  value={rating100}
                  onSetRating={(value) => setRating(value)}
                  disabled={isUpdating}
                />
              </Col>
            </Form.Group>

            {isSingleScene && additionalComponents[2]}

            <Form.Group controlId="studio" as={Row}>
              {FormUtils.renderLabel({
                title: intl.formatMessage({ id: "studio" }),
              })}
              <Col xs={9}>
                <StudioSelect
                  onSelect={(items) =>
                    setStudioId(items.length > 0 ? items[0]?.id : undefined)
                  }
                  ids={studioId ? [studioId] : []}
                  isDisabled={isUpdating}
                />
              </Col>
            </Form.Group>

            <Form.Group controlId="performers">
              <Form.Label>
                <FormattedMessage id="performers" />
              </Form.Label>
              {renderMultiSelect("performers", performerIds)}
            </Form.Group>

            <Form.Group controlId="movies">
              <Form.Label>
                <FormattedMessage id="movies" />
              </Form.Label>
              {renderMultiSelect("movies", movieIds)}
            </Form.Group>

            <Form.Group controlId="tags">
              <Form.Label>
                <FormattedMessage id="tags" />
              </Form.Label>
              {renderMultiSelect("tags", tagIds)}
            </Form.Group>

            {isSingleScene && additionalComponents[3]}

            <Form.Group controlId="organized">
              <Form.Check
                type="checkbox"
                label={intl.formatMessage({ id: "organized" })}
                checked={organized}
                ref={checkboxRef}
                onChange={() => cycleOrganized()}
              />
            </Form.Group>
          </Form>
        </div>
      </ModalComponent>
    );
  }

  return render();
};
