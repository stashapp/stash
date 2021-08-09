import React, { useState } from "react";
import { StudioSelect, PerformerSelect } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { MovieSelect, TagSelect } from "src/components/Shared/Select";
import {
  ScrapeDialog,
  ScrapeDialogRow,
  ScrapeResult,
  ScrapedInputGroupRow,
  ScrapedTextAreaRow,
  ScrapedImageRow,
} from "src/components/Shared/ScrapeDialog";
import _ from "lodash";
import {
  useStudioCreate,
  usePerformerCreate,
  useMovieCreate,
  useTagCreate,
  makePerformerCreateInput,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { DurationUtils } from "src/utils";
import { useIntl } from "react-intl";

function renderScrapedStudio(
  result: ScrapeResult<string>,
  isNew?: boolean,
  onChange?: (value: string) => void
) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ? [resultValue] : [];

  return (
    <StudioSelect
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items[0]?.id);
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedStudioRow(
  title: string,
  result: ScrapeResult<string>,
  onChange: (value: ScrapeResult<string>) => void,
  newStudio?: GQL.ScrapedSceneStudio,
  onCreateNew?: (value: GQL.ScrapedSceneStudio) => void
) {
  return (
    <ScrapeDialogRow
      title={title}
      result={result}
      renderOriginalField={() => renderScrapedStudio(result)}
      renderNewField={() =>
        renderScrapedStudio(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
      newValues={newStudio ? [newStudio] : undefined}
      onCreateNew={onCreateNew}
    />
  );
}

function renderScrapedPerformers(
  result: ScrapeResult<string[]>,
  isNew?: boolean,
  onChange?: (value: string[]) => void
) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ?? [];

  return (
    <PerformerSelect
      isMulti
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items.map((i) => i.id));
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedPerformersRow(
  title: string,
  result: ScrapeResult<string[]>,
  onChange: (value: ScrapeResult<string[]>) => void,
  newPerformers: GQL.ScrapedScenePerformer[],
  onCreateNew?: (value: GQL.ScrapedScenePerformer) => void
) {
  return (
    <ScrapeDialogRow
      title={title}
      result={result}
      renderOriginalField={() => renderScrapedPerformers(result)}
      renderNewField={() =>
        renderScrapedPerformers(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
      newValues={newPerformers}
      onCreateNew={onCreateNew}
    />
  );
}

function renderScrapedMovies(
  result: ScrapeResult<string[]>,
  isNew?: boolean,
  onChange?: (value: string[]) => void
) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ?? [];

  return (
    <MovieSelect
      isMulti
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items.map((i) => i.id));
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedMoviesRow(
  title: string,
  result: ScrapeResult<string[]>,
  onChange: (value: ScrapeResult<string[]>) => void,
  newMovies: GQL.ScrapedSceneMovie[],
  onCreateNew?: (value: GQL.ScrapedSceneMovie) => void
) {
  return (
    <ScrapeDialogRow
      title={title}
      result={result}
      renderOriginalField={() => renderScrapedMovies(result)}
      renderNewField={() =>
        renderScrapedMovies(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
      newValues={newMovies}
      onCreateNew={onCreateNew}
    />
  );
}

function renderScrapedTags(
  result: ScrapeResult<string[]>,
  isNew?: boolean,
  onChange?: (value: string[]) => void
) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ?? [];

  return (
    <TagSelect
      isMulti
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items.map((i) => i.id));
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedTagsRow(
  title: string,
  result: ScrapeResult<string[]>,
  onChange: (value: ScrapeResult<string[]>) => void,
  newTags: GQL.ScrapedSceneTag[],
  onCreateNew?: (value: GQL.ScrapedSceneTag) => void
) {
  return (
    <ScrapeDialogRow
      title={title}
      result={result}
      renderOriginalField={() => renderScrapedTags(result)}
      renderNewField={() =>
        renderScrapedTags(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      newValues={newTags}
      onChange={onChange}
      onCreateNew={onCreateNew}
    />
  );
}

interface ISceneScrapeDialogProps {
  scene: Partial<GQL.SceneUpdateInput>;
  scraped: GQL.ScrapedScene;

  onClose: (scrapedScene?: GQL.ScrapedScene) => void;
}

interface IHasStoredID {
  stored_id?: string | null;
}

export const SceneScrapeDialog: React.FC<ISceneScrapeDialogProps> = (
  props: ISceneScrapeDialogProps
) => {
  const [title, setTitle] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.scene.title, props.scraped.title)
  );
  const [url, setURL] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.scene.url, props.scraped.url)
  );
  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.scene.date, props.scraped.date)
  );
  const [studio, setStudio] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.scene.studio_id,
      props.scraped.studio?.stored_id
    )
  );
  const [newStudio, setNewStudio] = useState<
    GQL.ScrapedSceneStudio | undefined
  >(
    props.scraped.studio && !props.scraped.studio.stored_id
      ? props.scraped.studio
      : undefined
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

    const ret = _.clone(idList);
    // sort by id numerically
    ret.sort((a, b) => {
      return parseInt(a, 10) - parseInt(b, 10);
    });

    return ret;
  }

  const [performers, setPerformers] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.scene.performer_ids),
      mapStoredIdObjects(props.scraped.performers ?? undefined)
    )
  );
  const [newPerformers, setNewPerformers] = useState<
    GQL.ScrapedScenePerformer[]
  >(props.scraped.performers?.filter((t) => !t.stored_id) ?? []);

  const [movies, setMovies] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.scene.movies?.map((p) => p.movie_id)),
      mapStoredIdObjects(props.scraped.movies ?? undefined)
    )
  );
  const [newMovies, setNewMovies] = useState<GQL.ScrapedSceneMovie[]>(
    props.scraped.movies?.filter((t) => !t.stored_id) ?? []
  );

  const [tags, setTags] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.scene.tag_ids),
      mapStoredIdObjects(props.scraped.tags ?? undefined)
    )
  );
  const [newTags, setNewTags] = useState<GQL.ScrapedSceneTag[]>(
    props.scraped.tags?.filter((t) => !t.stored_id) ?? []
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.scene.details, props.scraped.details)
  );
  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.scene.cover_image, props.scraped.image)
  );

  const [createStudio] = useStudioCreate();
  const [createPerformer] = usePerformerCreate();
  const [createMovie] = useMovieCreate();
  const [createTag] = useTagCreate();

  const intl = useIntl();
  const Toast = useToast();

  // don't show the dialog if nothing was scraped
  if (
    [title, url, date, studio, performers, movies, tags, details, image].every(
      (r) => !r.scraped
    )
  ) {
    props.onClose();
    return <></>;
  }

  async function createNewStudio(toCreate: GQL.ScrapedSceneStudio) {
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

  async function createNewPerformer(toCreate: GQL.ScrapedScenePerformer) {
    const input = makePerformerCreateInput(toCreate);

    try {
      const result = await createPerformer({
        variables: { input },
      });

      // add the new performer to the new performers value
      const performerClone = performers.cloneWithValue(performers.newValue);
      if (!performerClone.newValue) {
        performerClone.newValue = [];
      }
      performerClone.newValue.push(result.data!.performerCreate!.id);
      setPerformers(performerClone);

      // remove the performer from the list
      const newPerformersClone = newPerformers.concat();
      const pIndex = newPerformersClone.indexOf(toCreate);
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

  async function createNewMovie(toCreate: GQL.ScrapedSceneMovie) {
    let movieInput: GQL.MovieCreateInput = { name: "" };
    try {
      movieInput = Object.assign(movieInput, toCreate);

      // #788 - convert duration and rating to the correct type
      movieInput.duration = DurationUtils.stringToSeconds(
        toCreate.duration ?? undefined
      );
      if (!movieInput.duration) {
        movieInput.duration = undefined;
      }

      movieInput.rating = parseInt(toCreate.rating ?? "0", 10);
      if (!movieInput.rating || Number.isNaN(movieInput.rating)) {
        movieInput.rating = undefined;
      }

      const result = await createMovie({
        variables: movieInput,
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
      const pIndex = newMoviesClone.indexOf(toCreate);
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

  async function createNewTag(toCreate: GQL.ScrapedSceneTag) {
    const tagInput: GQL.TagCreateInput = { name: toCreate.name ?? "" };
    try {
      const result = await createTag({
        variables: {
          input: tagInput,
        },
      });

      // add the new tag to the new tags value
      const tagClone = tags.cloneWithValue(tags.newValue);
      if (!tagClone.newValue) {
        tagClone.newValue = [];
      }
      tagClone.newValue.push(result.data!.tagCreate!.id);
      setTags(tagClone);

      // remove the tag from the list
      const newTagsClone = newTags.concat();
      const pIndex = newTagsClone.indexOf(toCreate);
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
      url: url.getNewValue(),
      date: date.getNewValue(),
      studio: newStudioValue
        ? {
            stored_id: newStudioValue,
            name: "",
          }
        : undefined,
      performers: performers.getNewValue()?.map((p) => {
        return {
          stored_id: p,
          name: "",
        };
      }),
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
          title={intl.formatMessage({ id: "url" })}
          result={url}
          onChange={(value) => setURL(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "date" })}
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        {renderScrapedStudioRow(
          intl.formatMessage({ id: "studios" }),
          studio,
          (value) => setStudio(value),
          newStudio,
          createNewStudio
        )}
        {renderScrapedPerformersRow(
          intl.formatMessage({ id: "performers" }),
          performers,
          (value) => setPerformers(value),
          newPerformers,
          createNewPerformer
        )}
        {renderScrapedMoviesRow(
          intl.formatMessage({ id: "movies" }),
          movies,
          (value) => setMovies(value),
          newMovies,
          createNewMovie
        )}
        {renderScrapedTagsRow(
          intl.formatMessage({ id: "tags" }),
          tags,
          (value) => setTags(value),
          newTags,
          createNewTag
        )}
        <ScrapedTextAreaRow
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
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
        props.onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    />
  );
};
