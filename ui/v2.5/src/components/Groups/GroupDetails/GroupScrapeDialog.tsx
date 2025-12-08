import React, { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapedInputGroupRow,
  ScrapedImageRow,
  ScrapedTextAreaRow,
  ScrapedStringListRow,
} from "src/components/Shared/ScrapeDialog/ScrapeDialogRow";
import { ScrapeDialog } from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import TextUtils from "src/utils/text";
import {
  ObjectScrapeResult,
  ScrapeResult,
} from "src/components/Shared/ScrapeDialog/scrapeResult";
import { Studio } from "src/components/Studios/StudioSelect";
import { useCreateScrapedStudio } from "src/components/Shared/ScrapeDialog/createObjects";
import { ScrapedStudioRow } from "src/components/Shared/ScrapeDialog/ScrapedObjectsRow";
import { uniq } from "lodash-es";
import { Tag } from "src/components/Tags/TagSelect";
import { useScrapedTags } from "src/components/Shared/ScrapeDialog/scrapedTags";

interface IGroupScrapeDialogProps {
  group: Partial<GQL.GroupUpdateInput>;
  groupStudio: Studio | null;
  groupTags: Tag[];
  scraped: GQL.ScrapedGroup;

  onClose: (scrapedGroup?: GQL.ScrapedGroup) => void;
}

export const GroupScrapeDialog: React.FC<IGroupScrapeDialogProps> = ({
  group,
  groupStudio: groupStudio,
  groupTags: groupTags,
  scraped,
  onClose,
}) => {
  const intl = useIntl();

  const [name, setName] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(group.name, scraped.name)
  );
  const [aliases, setAliases] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(group.aliases, scraped.aliases)
  );
  const [duration, setDuration] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      TextUtils.secondsToTimestamp(group.duration || 0),
      // convert seconds to string if it's a number
      scraped.duration && !isNaN(+scraped.duration)
        ? TextUtils.secondsToTimestamp(parseInt(scraped.duration, 10))
        : scraped.duration
    )
  );
  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(group.date, scraped.date)
  );
  const [director, setDirector] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(group.director, scraped.director)
  );
  const [synopsis, setSynopsis] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(group.synopsis, scraped.synopsis)
  );
  const [studio, setStudio] = useState<ObjectScrapeResult<GQL.ScrapedStudio>>(
    new ObjectScrapeResult<GQL.ScrapedStudio>(
      groupStudio
        ? {
            stored_id: groupStudio.id,
            name: groupStudio.name,
          }
        : undefined,
      scraped.studio?.stored_id ? scraped.studio : undefined
    )
  );
  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      group.urls,
      scraped.urls
        ? uniq((group.urls ?? []).concat(scraped.urls ?? []))
        : undefined
    )
  );
  const [frontImage, setFrontImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(group.front_image, scraped.front_image)
  );
  const [backImage, setBackImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(group.back_image, scraped.back_image)
  );

  const [newStudio, setNewStudio] = useState<GQL.ScrapedStudio | undefined>(
    scraped.studio && !scraped.studio.stored_id ? scraped.studio : undefined
  );

  const createNewStudio = useCreateScrapedStudio({
    scrapeResult: studio,
    setScrapeResult: setStudio,
    setNewObject: setNewStudio,
  });

  const { tags, newTags, scrapedTagsRow, linkDialog } = useScrapedTags(
    groupTags,
    scraped.tags
  );

  const allFields = [
    name,
    aliases,
    duration,
    date,
    director,
    synopsis,
    studio,
    tags,
    urls,
    frontImage,
    backImage,
  ];
  // don't show the dialog if nothing was scraped
  if (
    allFields.every((r) => !r.scraped) &&
    !newStudio &&
    newTags.length === 0
  ) {
    onClose();
    return <></>;
  }

  function makeNewScrapedItem(): GQL.ScrapedGroup {
    const newStudioValue = studio.getNewValue();
    const durationString = duration.getNewValue();

    return {
      name: name.getNewValue() ?? "",
      aliases: aliases.getNewValue(),
      duration: durationString,
      date: date.getNewValue(),
      director: director.getNewValue(),
      synopsis: synopsis.getNewValue(),
      studio: newStudioValue,
      tags: tags.getNewValue(),
      urls: urls.getNewValue(),
      front_image: frontImage.getNewValue(),
      back_image: backImage.getNewValue(),
    };
  }

  function renderScrapeRows() {
    return (
      <>
        <ScrapedInputGroupRow
          field="name"
          title={intl.formatMessage({ id: "name" })}
          result={name}
          onChange={(value) => setName(value)}
        />
        <ScrapedInputGroupRow
          field="aliases"
          title={intl.formatMessage({ id: "aliases" })}
          result={aliases}
          onChange={(value) => setAliases(value)}
        />
        <ScrapedInputGroupRow
          field="duration"
          title={intl.formatMessage({ id: "duration" })}
          result={duration}
          onChange={(value) => setDuration(value)}
        />
        <ScrapedInputGroupRow
          field="date"
          title={intl.formatMessage({ id: "date" })}
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        <ScrapedInputGroupRow
          field="director"
          title={intl.formatMessage({ id: "director" })}
          result={director}
          onChange={(value) => setDirector(value)}
        />
        <ScrapedTextAreaRow
          field="synopsis"
          title={intl.formatMessage({ id: "synopsis" })}
          result={synopsis}
          onChange={(value) => setSynopsis(value)}
        />
        <ScrapedStudioRow
          field="studio"
          title={intl.formatMessage({ id: "studios" })}
          result={studio}
          onChange={(value) => setStudio(value)}
          newStudio={newStudio}
          onCreateNew={createNewStudio}
        />
        <ScrapedStringListRow
          field="urls"
          title={intl.formatMessage({ id: "urls" })}
          result={urls}
          onChange={(value) => setURLs(value)}
        />
        {scrapedTagsRow}
        <ScrapedImageRow
          field="front_image"
          title="Front Image"
          className="group-image"
          result={frontImage}
          onChange={(value) => setFrontImage(value)}
        />
        <ScrapedImageRow
          field="back_image"
          title="Back Image"
          className="group-image"
          result={backImage}
          onChange={(value) => setBackImage(value)}
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
        { entity_type: intl.formatMessage({ id: "group" }) }
      )}
      onClose={(apply) => {
        onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    >
      {renderScrapeRows()}
    </ScrapeDialog>
  );
};
