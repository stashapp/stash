import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapeDialog,
  ScrapedInputGroupRow,
  ScrapedTextAreaRow,
  ScrapedImageRow,
  ScrapedStringListRow,
} from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import clone from "lodash-es/clone";
import {
  useStudioCreate,
  usePerformerCreate,
  useMovieCreate,
  useTagCreate,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { useIntl } from "react-intl";
import { uniq } from "lodash-es";
import { scrapedPerformerToCreateInput } from "src/core/performers";
import { scrapedMovieToCreateInput } from "src/core/movies";
import { Performer } from "src/components/Performers/PerformerSelect";
import { IHasStoredID, sortStoredIdObjects } from "src/utils/data";
import { ScrapeResult } from "src/components/Shared/ScrapeDialog/scrapeResult";
import {
  ScrapedMoviesRow,
  ScrapedPerformersRow,
  ScrapedStudioRow,
  ScrapedTagsRow,
} from "src/components/Shared/ScrapeDialog/ScrapedObjectsRow";

interface ISceneScrapeDialogProps {
  scene: Partial<GQL.SceneUpdateInput>;
  scenePerformers: Performer[];
  scraped: GQL.ScrapedScene;
  endpoint?: string;

  onClose: (scrapedScene?: GQL.ScrapedScene) => void;
}

export const SceneScrapeDialog: React.FC<ISceneScrapeDialogProps> = ({
  scene,
  scenePerformers,
  scraped,
  onClose,
  endpoint,
}) => {
  const [title, setTitle] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.title, scraped.title)
  );
  const [code, setCode] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.code, scraped.code)
  );

  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      scene.urls,
      scraped.urls
        ? uniq((scene.urls ?? []).concat(scraped.urls ?? []))
        : undefined
    )
  );

  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.date, scraped.date)
  );
  const [director, setDirector] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.director, scraped.director)
  );
  const [studio, setStudio] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.studio_id, scraped.studio?.stored_id)
  );
  const [newStudio, setNewStudio] = useState<GQL.ScrapedStudio | undefined>(
    scraped.studio && !scraped.studio.stored_id ? scraped.studio : undefined
  );

  const [stashID, setStashID] = useState(
    new ScrapeResult<string>(
      scene.stash_ids?.find((s) => s.endpoint === endpoint)?.stash_id,
      scraped.remote_site_id
    )
  );

  function mapStoredIdObjects(
    scrapedObjects?: IHasStoredID[]
  ): string[] | undefined {
    if (!scrapedObjects) {
      return undefined;
    }
    const ret = scrapedObjects
      .map((p) => p.stored_id)
      .filter((p) => {
        return p !== undefined && p !== null;
      }) as string[];

    if (ret.length === 0) {
      return undefined;
    }

    // sort by id numerically
    ret.sort((a, b) => {
      return parseInt(a, 10) - parseInt(b, 10);
    });

    return ret;
  }

  function sortIdList(idList?: string[] | null) {
    if (!idList) {
      return;
    }

    const ret = clone(idList);
    // sort by id numerically
    ret.sort((a, b) => {
      return parseInt(a, 10) - parseInt(b, 10);
    });

    return ret;
  }

  const [performers, setPerformers] = useState<
    ScrapeResult<GQL.ScrapedPerformer[]>
  >(
    new ScrapeResult<GQL.ScrapedPerformer[]>(
      sortStoredIdObjects(
        scenePerformers.map((p) => ({
          stored_id: p.id,
          name: p.name,
        }))
      ),
      sortStoredIdObjects(scraped.performers ?? undefined)
    )
  );
  const [newPerformers, setNewPerformers] = useState<GQL.ScrapedPerformer[]>(
    scraped.performers?.filter((t) => !t.stored_id) ?? []
  );

  const [movies, setMovies] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(scene.movies?.map((p) => p.movie_id)),
      mapStoredIdObjects(scraped.movies ?? undefined)
    )
  );
  const [newMovies, setNewMovies] = useState<GQL.ScrapedMovie[]>(
    scraped.movies?.filter((t) => !t.stored_id) ?? []
  );

  const [tags, setTags] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(scene.tag_ids),
      mapStoredIdObjects(scraped.tags ?? undefined)
    )
  );
  const [newTags, setNewTags] = useState<GQL.ScrapedTag[]>(
    scraped.tags?.filter((t) => !t.stored_id) ?? []
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.details, scraped.details)
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(scene.cover_image, scraped.image)
  );

  const [createStudio] = useStudioCreate();
  const [createPerformer] = usePerformerCreate();
  const [createMovie] = useMovieCreate();
  const [createTag] = useTagCreate();

  const intl = useIntl();
  const Toast = useToast();

  // don't show the dialog if nothing was scraped
  if (
    [
      title,
      code,
      urls,
      date,
      director,
      studio,
      performers,
      movies,
      tags,
      details,
      image,
      stashID,
    ].every((r) => !r.scraped) &&
    newTags.length === 0 &&
    newPerformers.length === 0 &&
    newMovies.length === 0 &&
    !newStudio
  ) {
    onClose();
    return <></>;
  }

  async function createNewStudio(toCreate: GQL.ScrapedStudio) {
    try {
      const result = await createStudio({
        variables: {
          input: {
            name: toCreate.name,
            url: toCreate.url,
          },
        },
      });

      // set the new studio as the value
      setStudio(studio.cloneWithValue(result.data!.studioCreate!.id));
      setNewStudio(undefined);

      Toast.success({
        content: (
          <span>
            Created studio: <b>{toCreate.name}</b>
          </span>
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function createNewPerformer(toCreate: GQL.ScrapedPerformer) {
    const input = scrapedPerformerToCreateInput(toCreate);

    try {
      const result = await createPerformer({
        variables: { input },
      });

      const newValue = [...(performers.newValue ?? [])];
      if (result.data?.performerCreate)
        newValue.push({
          stored_id: result.data.performerCreate.id,
          name: result.data.performerCreate.name,
        });

      // add the new performer to the new performers value
      const performerClone = performers.cloneWithValue(newValue);
      setPerformers(performerClone);

      // remove the performer from the list
      const newPerformersClone = newPerformers.concat();
      const pIndex = newPerformersClone.findIndex(
        (p) => p.name === toCreate.name
      );
      if (pIndex === -1) throw new Error("Could not find performer to remove");

      newPerformersClone.splice(pIndex, 1);

      setNewPerformers(newPerformersClone);

      Toast.success({
        content: (
          <span>
            Created performer: <b>{toCreate.name}</b>
          </span>
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function createNewMovie(toCreate: GQL.ScrapedMovie) {
    const movieInput = scrapedMovieToCreateInput(toCreate);
    try {
      const result = await createMovie({
        variables: { input: movieInput },
      });

      // add the new movie to the new movies value
      const movieClone = movies.cloneWithValue(movies.newValue);
      if (!movieClone.newValue) {
        movieClone.newValue = [];
      }
      movieClone.newValue.push(result.data!.movieCreate!.id);
      setMovies(movieClone);

      // remove the movie from the list
      const newMoviesClone = newMovies.concat();
      const pIndex = newMoviesClone.findIndex((p) => p.name === toCreate.name);
      if (pIndex === -1) throw new Error("Could not find movie to remove");
      newMoviesClone.splice(pIndex, 1);

      setNewMovies(newMoviesClone);

      Toast.success({
        content: (
          <span>
            Created movie: <b>{toCreate.name}</b>
          </span>
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function createNewTag(toCreate: GQL.ScrapedTag) {
    const tagInput: GQL.TagCreateInput = { name: toCreate.name ?? "" };
    try {
      const result = await createTag({
        variables: {
          input: tagInput,
        },
      });

      const newValue = [...(tags.newValue ?? [])];
      if (result.data?.tagCreate) newValue.push(result.data.tagCreate.id);

      // add the new tag to the new tags value
      const tagClone = tags.cloneWithValue(newValue);
      setTags(tagClone);

      // remove the tag from the list
      const newTagsClone = newTags.concat();
      const pIndex = newTagsClone.indexOf(toCreate);
      if (pIndex === -1) throw new Error("Could not find tag to remove");
      newTagsClone.splice(pIndex, 1);

      setNewTags(newTagsClone);

      Toast.success({
        content: (
          <span>
            Created tag: <b>{toCreate.name}</b>
          </span>
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function makeNewScrapedItem(): GQL.ScrapedSceneDataFragment {
    const newStudioValue = studio.getNewValue();

    return {
      title: title.getNewValue(),
      code: code.getNewValue(),
      urls: urls.getNewValue(),
      date: date.getNewValue(),
      director: director.getNewValue(),
      studio: newStudioValue
        ? {
            stored_id: newStudioValue,
            name: "",
          }
        : undefined,
      performers: performers.getNewValue(),
      movies: movies.getNewValue()?.map((m) => {
        return {
          stored_id: m,
          name: "",
        };
      }),
      tags: tags.getNewValue()?.map((m) => {
        return {
          stored_id: m,
          name: "",
        };
      }),
      details: details.getNewValue(),
      image: image.getNewValue(),
      remote_site_id: stashID.getNewValue(),
    };
  }

  function renderScrapeRows() {
    return (
      <>
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "title" })}
          result={title}
          onChange={(value) => setTitle(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "scene_code" })}
          result={code}
          onChange={(value) => setCode(value)}
        />
        <ScrapedStringListRow
          title={intl.formatMessage({ id: "urls" })}
          result={urls}
          onChange={(value) => setURLs(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "date" })}
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "director" })}
          result={director}
          onChange={(value) => setDirector(value)}
        />
        <ScrapedStudioRow
          title={intl.formatMessage({ id: "studios" })}
          result={studio}
          onChange={(value) => setStudio(value)}
          newStudio={newStudio}
          onCreateNew={createNewStudio}
        />
        <ScrapedPerformersRow
          title={intl.formatMessage({ id: "performers" })}
          result={performers}
          onChange={(value) => setPerformers(value)}
          newObjects={newPerformers}
          onCreateNew={createNewPerformer}
        />
        <ScrapedMoviesRow
          title={intl.formatMessage({ id: "movies" })}
          result={movies}
          onChange={(value) => setMovies(value)}
          newObjects={newMovies}
          onCreateNew={createNewMovie}
        />
        <ScrapedTagsRow
          title={intl.formatMessage({ id: "tags" })}
          result={tags}
          onChange={(value) => setTags(value)}
          newObjects={newTags}
          onCreateNew={createNewTag}
        />
        <ScrapedTextAreaRow
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "stash_id" })}
          result={stashID}
          locked
          onChange={(value) => setStashID(value)}
        />
        <ScrapedImageRow
          title={intl.formatMessage({ id: "cover_image" })}
          className="scene-cover"
          result={image}
          onChange={(value) => setImage(value)}
        />
      </>
    );
  }

  return (
    <ScrapeDialog
      title={intl.formatMessage(
        { id: "dialogs.scrape_entity_title" },
        { entity_type: intl.formatMessage({ id: "scene" }) }
      )}
      renderScrapeRows={renderScrapeRows}
      onClose={(apply) => {
        onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    />
  );
};

export default SceneScrapeDialog;
