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
import { IStashBox } from "./StudioStashBoxModal";
import { ScrapeResult } from "src/components/Shared/ScrapeDialog/scrapeResult";
import { useScrapedTags } from "src/components/Shared/ScrapeDialog/scrapedTags";
import { Tag } from "src/components/Tags/TagSelect";
import { uniq } from "lodash-es";

interface IStudioScrapeDialogProps {
  studio: Partial<GQL.StudioUpdateInput>;
  studioTags: Tag[];
  scraped: GQL.ScrapedStudio;
  scraper: IStashBox;
  onClose: (scrapedStudio?: GQL.ScrapedStudio) => void;
}

export const StudioScrapeDialog: React.FC<IStudioScrapeDialogProps> = ({
  studio,
  studioTags,
  scraped,
  scraper,
  onClose,
}) => {
  const intl = useIntl();

  const { endpoint } = scraper;

  function getCurrentRemoteSiteID() {
    if (!endpoint) {
      return;
    }

    const stashIDs = (studio.stash_ids ?? []).filter(
      (s) => s.endpoint === endpoint
    );
    if (stashIDs.length > 1 && scraped.remote_site_id) {
      const matchingID = stashIDs.find(
        (s) => s.stash_id === scraped.remote_site_id
      );
      if (matchingID) {
        return matchingID.stash_id;
      }
    }

    return studio.stash_ids?.find((s) => s.endpoint === endpoint)?.stash_id;
  }

  const [name, setName] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(studio.name, scraped.name)
  );

  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      studio.urls,
      scraped.urls
        ? uniq((studio.urls ?? []).concat(scraped.urls ?? []))
        : undefined
    )
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(studio.details, scraped.details)
  );

  const [aliases, setAliases] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      studio.aliases?.join(", "),
      scraped.aliases
    )
  );

  const [remoteSiteID, setRemoteSiteID] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      getCurrentRemoteSiteID(),
      scraped.remote_site_id
    )
  );

  const { tags, newTags, scrapedTagsRow, linkDialog } = useScrapedTags(
    studioTags,
    scraped.tags,
    endpoint
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(studio.image, scraped.image)
  );

  const allFields = [name, urls, details, aliases, tags, image, remoteSiteID];

  // don't show the dialog if nothing was scraped
  if (allFields.every((r) => !r.scraped) && newTags.length === 0) {
    onClose();
    return <></>;
  }

  function makeNewScrapedItem(): GQL.ScrapedStudio {
    return {
      name: name.getNewValue() ?? "",
      urls: urls.getNewValue(),
      details: details.getNewValue(),
      aliases: aliases.getNewValue(),
      tags: tags.getNewValue(),
      image: image.getNewValue(),
      remote_site_id: remoteSiteID.getNewValue(),
      // Include parent from original scraped data (read-only)
      parent: scraped.parent,
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
        <ScrapedStringListRow
          field="urls"
          title={intl.formatMessage({ id: "urls" })}
          result={urls}
          onChange={(value) => setURLs(value)}
        />
        <ScrapedTextAreaRow
          field="details"
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
        />
        <ScrapedTextAreaRow
          field="aliases"
          title={intl.formatMessage({ id: "aliases" })}
          result={aliases}
          onChange={(value) => setAliases(value)}
        />
        {scrapedTagsRow}
        <ScrapedImageRow
          field="image"
          title={intl.formatMessage({ id: "studio_image" })}
          className="studio-image"
          result={image}
          onChange={(value) => setImage(value)}
        />
        <ScrapedInputGroupRow
          field="remote_site_id"
          title={intl.formatMessage({ id: "stash_id" })}
          result={remoteSiteID}
          locked
          onChange={(value) => setRemoteSiteID(value)}
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
        { entity_type: intl.formatMessage({ id: "studio" }) }
      )}
      onClose={(apply) => {
        onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    >
      {renderScrapeRows()}
    </ScrapeDialog>
  );
};
