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
  result: ScrapeResult<string>,
  onChange: (value: ScrapeResult<string>) => void
) {
  return (
    <ScrapeDialogRow
      title="Studio"
      result={result}
      renderOriginalField={() => renderScrapedStudio(result)}
      renderNewField={() =>
        renderScrapedStudio(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
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
  result: ScrapeResult<string[]>,
  onChange: (value: ScrapeResult<string[]>) => void
) {
  return (
    <ScrapeDialogRow
      title="Performers"
      result={result}
      renderOriginalField={() => renderScrapedPerformers(result)}
      renderNewField={() =>
        renderScrapedPerformers(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
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
  result: ScrapeResult<string[]>,
  onChange: (value: ScrapeResult<string[]>) => void
) {
  return (
    <ScrapeDialogRow
      title="Movies"
      result={result}
      renderOriginalField={() => renderScrapedMovies(result)}
      renderNewField={() =>
        renderScrapedMovies(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
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
  result: ScrapeResult<string[]>,
  onChange: (value: ScrapeResult<string[]>) => void
) {
  return (
    <ScrapeDialogRow
      title="Tags"
      result={result}
      renderOriginalField={() => renderScrapedTags(result)}
      renderNewField={() =>
        renderScrapedTags(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
    />
  );
}

interface ISceneScrapeDialogProps {
  scene: Partial<GQL.SceneUpdateInput>;
  scraped: GQL.ScrapedScene;

  onClose: (scrapedScene?: GQL.ScrapedScene) => void;
}

interface IHasID {
  id?: string | null;
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
    new ScrapeResult<string>(props.scene.studio_id, props.scraped.studio?.id)
  );

  function mapIdObjects(scrapedObjects?: IHasID[]): string[] | undefined {
    if (!scrapedObjects) {
      return undefined;
    }
    const ret = scrapedObjects
      .map((p) => p.id)
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
      mapIdObjects(props.scraped.performers ?? undefined)
    )
  );
  const [movies, setMovies] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.scene.movies?.map((p) => p.movie_id)),
      mapIdObjects(props.scraped.movies ?? undefined)
    )
  );
  const [tags, setTags] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.scene.tag_ids),
      mapIdObjects(props.scraped.tags ?? undefined)
    )
  );
  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.scene.details, props.scraped.details)
  );
  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.scene.cover_image, props.scraped.image)
  );

  // don't show the dialog if nothing was scraped
  if (
    [title, url, date, studio, performers, movies, tags, details, image].every(
      (r) => !r.scraped
    )
  ) {
    props.onClose();
    return <></>;
  }

  function makeNewScrapedItem() {
    const newStudio = studio.getNewValue();

    return {
      title: title.getNewValue(),
      url: url.getNewValue(),
      date: date.getNewValue(),
      studio: newStudio
        ? {
            id: newStudio,
            name: "",
          }
        : undefined,
      performers: performers.getNewValue()?.map((p) => {
        return {
          id: p,
          name: "",
        };
      }),
      movies: movies.getNewValue()?.map((m) => {
        return {
          id: m,
          name: "",
        };
      }),
      tags: tags.getNewValue()?.map((m) => {
        return {
          id: m,
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
          title="Title"
          result={title}
          onChange={(value) => setTitle(value)}
        />
        <ScrapedInputGroupRow
          title="URL"
          result={url}
          onChange={(value) => setURL(value)}
        />
        <ScrapedInputGroupRow
          title="Date"
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        {renderScrapedStudioRow(studio, (value) => setStudio(value))}
        {renderScrapedPerformersRow(performers, (value) =>
          setPerformers(value)
        )}
        {renderScrapedMoviesRow(movies, (value) => setMovies(value))}
        {renderScrapedTagsRow(tags, (value) => setTags(value))}
        <ScrapedTextAreaRow
          title="Details"
          result={details}
          onChange={(value) => setDetails(value)}
        />
        <ScrapedImageRow
          title="Cover Image"
          className="scene-cover"
          result={image}
          onChange={(value) => setImage(value)}
        />
      </>
    );
  }

  return (
    <ScrapeDialog
      title="Scene Scrape Results"
      renderScrapeRows={renderScrapeRows}
      onClose={(apply) => {
        props.onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    />
  );
};
