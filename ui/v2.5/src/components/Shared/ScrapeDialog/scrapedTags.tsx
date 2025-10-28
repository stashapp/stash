import { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { ObjectListScrapeResult } from "./scrapeResult";
import { sortStoredIdObjects } from "src/utils/data";
import { Tag } from "src/components/Tags/TagSelect";
import { useCreateScrapedTag } from "./createObjects";
import { ScrapedTagsRow } from "./ScrapedObjectsRow";

export function useScrapedTags(
  existingTags: Tag[],
  scrapedTags?: GQL.Maybe<GQL.ScrapedTag[]>
) {
  const intl = useIntl();
  const [tags, setTags] = useState<ObjectListScrapeResult<GQL.ScrapedTag>>(
    new ObjectListScrapeResult<GQL.ScrapedTag>(
      sortStoredIdObjects(
        existingTags.map((t) => ({
          stored_id: t.id,
          name: t.name,
        }))
      ),
      sortStoredIdObjects(scrapedTags ?? undefined)
    )
  );

  const [newTags, setNewTags] = useState<GQL.ScrapedTag[]>(
    scrapedTags?.filter((t) => !t.stored_id) ?? []
  );

  const createNewTag = useCreateScrapedTag({
    scrapeResult: tags,
    setScrapeResult: setTags,
    newObjects: newTags,
    setNewObjects: setNewTags,
  });

  const scrapedTagsRow = (
    <ScrapedTagsRow
      title={intl.formatMessage({ id: "tags" })}
      result={tags}
      onChange={(value) => setTags(value)}
      newObjects={newTags}
      onCreateNew={createNewTag}
    />
  );

  return {
    tags,
    newTags,
    scrapedTagsRow,
  };
}
