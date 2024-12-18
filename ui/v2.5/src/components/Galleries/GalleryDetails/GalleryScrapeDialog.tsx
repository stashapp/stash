import React, { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapeDialog,
  ScrapedInputGroupRow,
  ScrapedStringListRow,
  ScrapedTextAreaRow,
} from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import {
  ObjectListScrapeResult,
  ObjectScrapeResult,
  ScrapeResult,
} from "src/components/Shared/ScrapeDialog/scrapeResult";
import {
  ScrapedPerformersRow,
  ScrapedStudioRow,
} from "src/components/Shared/ScrapeDialog/ScrapedObjectsRow";
import { sortStoredIdObjects } from "src/utils/data";
import { Performer } from "src/components/Performers/PerformerSelect";
import {
  useCreateScrapedPerformer,
  useCreateScrapedStudio,
} from "src/components/Shared/ScrapeDialog/createObjects";
import { uniq } from "lodash-es";
import { Tag } from "src/components/Tags/TagSelect";
import { Studio } from "src/components/Studios/StudioSelect";
import { useScrapedTags } from "src/components/Shared/ScrapeDialog/scrapedTags";

interface IGalleryScrapeDialogProps {
  gallery: Partial<GQL.GalleryUpdateInput>;
  galleryStudio: Studio | null;
  galleryTags: Tag[];
  galleryPerformers: Performer[];
  scraped: GQL.ScrapedGallery;

  onClose: (scrapedGallery?: GQL.ScrapedGallery) => void;
}

export const GalleryScrapeDialog: React.FC<IGalleryScrapeDialogProps> = ({
  gallery,
  galleryStudio,
  galleryTags,
  galleryPerformers,
  scraped,
  onClose,
}) => {
  const intl = useIntl();
  const [title, setTitle] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(gallery.title, scraped.title)
  );
  const [code, setCode] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(gallery.code, scraped.code)
  );
  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      gallery.urls,
      scraped.urls
        ? uniq((gallery.urls ?? []).concat(scraped.urls ?? []))
        : undefined
    )
  );
  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(gallery.date, scraped.date)
  );
  const [photographer, setPhotographer] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(gallery.photographer, scraped.photographer)
  );
  const [studio, setStudio] = useState<ObjectScrapeResult<GQL.ScrapedStudio>>(
    new ObjectScrapeResult<GQL.ScrapedStudio>(
      galleryStudio
        ? {
            stored_id: galleryStudio.id,
            name: galleryStudio.name,
          }
        : undefined,
      scraped.studio
    )
  );
  const [newStudio, setNewStudio] = useState<GQL.ScrapedStudio | undefined>(
    scraped.studio && !scraped.studio.stored_id ? scraped.studio : undefined
  );

  const [performers, setPerformers] = useState<
    ObjectListScrapeResult<GQL.ScrapedPerformer>
  >(
    new ObjectListScrapeResult<GQL.ScrapedPerformer>(
      sortStoredIdObjects(
        galleryPerformers.map((p) => ({
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

  const { tags, newTags, scrapedTagsRow } = useScrapedTags(
    galleryTags,
    scraped.tags
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(gallery.details, scraped.details)
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

  // don't show the dialog if nothing was scraped
  if (
    [
      title,
      code,
      urls,
      date,
      photographer,
      studio,
      performers,
      tags,
      details,
    ].every((r) => !r.scraped) &&
    !newStudio &&
    newPerformers.length === 0 &&
    newTags.length === 0
  ) {
    onClose();
    return <></>;
  }

  function makeNewScrapedItem(): GQL.ScrapedGalleryDataFragment {
    const newStudioValue = studio.getNewValue();

    return {
      title: title.getNewValue(),
      code: code.getNewValue(),
      urls: urls.getNewValue(),
      date: date.getNewValue(),
      photographer: photographer.getNewValue(),
      studio: newStudioValue,
      performers: performers.getNewValue(),
      tags: tags.getNewValue(),
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
          title={intl.formatMessage({ id: "photographer" })}
          result={photographer}
          onChange={(value) => setPhotographer(value)}
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
        {scrapedTagsRow}
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
        onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    />
  );
};
