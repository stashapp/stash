import React, { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapeDialog,
  ScrapedInputGroupRow,
  ScrapedTextAreaRow,
} from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import clone from "lodash-es/clone";
import { ScrapeResult } from "src/components/Shared/ScrapeDialog/scrapeResult";
import {
  ScrapedPerformersRow,
  ScrapedStudioRow,
  ScrapedTagsRow,
} from "src/components/Shared/ScrapeDialog/ScrapedObjectsRow";
import { sortStoredIdObjects } from "src/utils/data";
import { Performer } from "src/components/Performers/PerformerSelect";
import {
  useCreateScrapedPerformer,
  useCreateScrapedStudio,
  useCreateScrapedTag,
} from "src/components/Shared/ScrapeDialog/createObjects";

interface IGalleryScrapeDialogProps {
  gallery: Partial<GQL.GalleryUpdateInput>;
  galleryPerformers: Performer[];
  scraped: GQL.ScrapedGallery;

  onClose: (scrapedGallery?: GQL.ScrapedGallery) => void;
}

interface IHasStoredID {
  stored_id?: string | null;
}

export const GalleryScrapeDialog: React.FC<IGalleryScrapeDialogProps> = (
  props: IGalleryScrapeDialogProps
) => {
  const intl = useIntl();
  const [title, setTitle] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.gallery.title, props.scraped.title)
  );
  const [url, setURL] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.gallery.url, props.scraped.url)
  );
  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.gallery.date, props.scraped.date)
  );
  const [studio, setStudio] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.gallery.studio_id,
      props.scraped.studio?.stored_id
    )
  );
  const [newStudio, setNewStudio] = useState<GQL.ScrapedStudio | undefined>(
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
        props.galleryPerformers.map((p) => ({
          stored_id: p.id,
          name: p.name,
        }))
      ),
      sortStoredIdObjects(props.scraped.performers ?? undefined)
    )
  );
  const [newPerformers, setNewPerformers] = useState<GQL.ScrapedPerformer[]>(
    props.scraped.performers?.filter((t) => !t.stored_id) ?? []
  );

  const [tags, setTags] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.gallery.tag_ids),
      mapStoredIdObjects(props.scraped.tags ?? undefined)
    )
  );
  const [newTags, setNewTags] = useState<GQL.ScrapedTag[]>(
    props.scraped.tags?.filter((t) => !t.stored_id) ?? []
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.gallery.details, props.scraped.details)
  );

  const createNewStudio = useCreateScrapedStudio({
    scrapeResult: studio,
    setScrapeResult: setStudio,
    setNewObject: setNewStudio,
  });

  const createNewPerformer = useCreateScrapedPerformer({
    scrapeResult: performers,
    setScrapeResult: setPerformers,
    newObjects: newPerformers,
    setNewObjects: setNewPerformers,
  });

  const createNewTag = useCreateScrapedTag({
    scrapeResult: tags,
    setScrapeResult: setTags,
    newObjects: newTags,
    setNewObjects: setNewTags,
  });

  // don't show the dialog if nothing was scraped
  if (
    [title, url, date, studio, performers, tags, details].every(
      (r) => !r.scraped
    ) &&
    !newStudio &&
    newPerformers.length === 0 &&
    newTags.length === 0
  ) {
    props.onClose();
    return <></>;
  }

  function makeNewScrapedItem(): GQL.ScrapedGalleryDataFragment {
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
      performers: performers.getNewValue(),
      tags: tags.getNewValue()?.map((m) => {
        return {
          stored_id: m,
          name: "",
        };
      }),
      details: details.getNewValue(),
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
      </>
    );
  }

  return (
    <ScrapeDialog
      title={intl.formatMessage(
        { id: "dialogs.scrape_entity_title" },
        { entity_type: intl.formatMessage({ id: "gallery" }) }
      )}
      renderScrapeRows={renderScrapeRows}
      onClose={(apply) => {
        props.onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    />
  );
};
