import React, { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapedInputGroupRow,
  ScrapedStringListRow,
  ScrapedTextAreaRow,
} from "src/components/Shared/ScrapeDialog/ScrapeDialogRow";
import { ScrapeDialog } from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import {
  ObjectListScrapeResult,
  ObjectScrapeResult,
  ScrapeResult,
} from "src/components/Shared/ScrapeDialog/scrapeResult";
import { sortStoredIdObjects } from "src/utils/data";
import { Performer } from "src/components/Performers/PerformerSelect";
import { uniq } from "lodash-es";
import { Tag } from "src/components/Tags/TagSelect";
import { Studio } from "src/components/Studios/StudioSelect";
import { useScrapedTags } from "src/components/Shared/ScrapeDialog/scrapedTags";
import { useScrapedStudios } from "src/components/Shared/ScrapeDialog/scrapedStudios";
import { useScrapedPerformers } from "src/components/Shared/ScrapeDialog/scrapedPerformers";

interface IImageScrapeDialogProps {
  image: Partial<GQL.ImageUpdateInput>;
  imageStudio: Studio | null;
  imageTags: Tag[];
  imagePerformers: Performer[];
  scraped: GQL.ScrapedImage;

  onClose: (scrapedImage?: GQL.ScrapedImage) => void;
}

export const ImageScrapeDialog: React.FC<IImageScrapeDialogProps> = ({
  image,
  imageStudio,
  imageTags,
  imagePerformers,
  scraped,
  onClose,
}) => {
  const intl = useIntl();
  const [title, setTitle] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(image.title, scraped.title)
  );
  const [code, setCode] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(image.code, scraped.code)
  );
  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      image.urls,
      scraped.urls
        ? uniq((image.urls ?? []).concat(scraped.urls ?? []))
        : undefined
    )
  );
  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(image.date, scraped.date)
  );

  const [photographer, setPhotographer] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(image.photographer, scraped.photographer)
  );

  const {
    studio,
    newStudio,
    linkDialog: studioLinkDialog,
    scrapedStudioRow,
  } = useScrapedStudios(imageStudio, scraped.studio);

  const {
    performers,
    newPerformers,
    linkDialog: performersLinkDialog,
    scrapedPerformersRow,
  } = useScrapedPerformers(
    imagePerformers,
    scraped.performers,
    undefined,
    date.useNewValue ? date.newValue : date.originalValue
  );

  const { tags, newTags, scrapedTagsRow, linkDialog } = useScrapedTags(
    imageTags,
    scraped.tags
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(image.details, scraped.details)
  );

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

  // render link dialogs if any are active
  if (studioLinkDialog) {
    return studioLinkDialog;
  }
  if (performersLinkDialog) {
    return performersLinkDialog;
  }

  function makeNewScrapedItem(): GQL.ScrapedImageDataFragment {
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
          field="title"
          title={intl.formatMessage({ id: "title" })}
          result={title}
          onChange={(value) => setTitle(value)}
        />
        <ScrapedInputGroupRow
          field="code"
          title={intl.formatMessage({ id: "scene_code" })}
          result={code}
          onChange={(value) => setCode(value)}
        />
        <ScrapedStringListRow
          field="urls"
          title={intl.formatMessage({ id: "urls" })}
          result={urls}
          onChange={(value) => setURLs(value)}
        />
        <ScrapedInputGroupRow
          field="date"
          title={intl.formatMessage({ id: "date" })}
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        <ScrapedInputGroupRow
          field="photographer"
          title={intl.formatMessage({ id: "photographer" })}
          result={photographer}
          onChange={(value) => setPhotographer(value)}
        />
        {scrapedStudioRow}
        {scrapedPerformersRow}
        {scrapedTagsRow}
        <ScrapedTextAreaRow
          field="details"
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
        />
      </>
    );
  }

  if (linkDialog) {
    return linkDialog;
  }

  return (
    <ScrapeDialog
      title={intl.formatMessage(
        { id: "dialogs.scrape_entity_title" },
        { entity_type: intl.formatMessage({ id: "image" }) }
      )}
      onClose={(apply) => {
        onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    >
      {renderScrapeRows()}
    </ScrapeDialog>
  );
};
