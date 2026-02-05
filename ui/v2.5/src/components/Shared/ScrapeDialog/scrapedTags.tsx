import { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { ObjectListScrapeResult } from "./scrapeResult";
import { sortStoredIdObjects } from "src/utils/data";
import { Tag } from "src/components/Tags/TagSelect";
import { useCreateScrapedTag, useLinkScrapedTag } from "./createObjects";
import { ScrapedTagsRow } from "./ScrapedObjectsRow";
import { CreateLinkTagDialog } from "src/components/Shared/ScrapeDialog/CreateLinkTagDialog";
import { useTagCreate, useTagUpdate } from "src/core/StashService";
import { toastOperation, useToast } from "src/hooks/Toast";

export function useScrapedTags(
  existingTags: Tag[],
  scrapedTags?: GQL.Maybe<GQL.ScrapedTag[]>,
  endpoint?: string
) {
  const intl = useIntl();
  const Toast = useToast();

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
  const [linkedTag, setLinkedTag] = useState<GQL.ScrapedTag | null>(null);

  const createNewTag = useCreateScrapedTag({
    scrapeResult: tags,
    setScrapeResult: setTags,
    newObjects: newTags,
    setNewObjects: setNewTags,
    endpoint,
  });

  const [createTag] = useTagCreate();
  const [updateTag] = useTagUpdate();

  const linkScrapedTag = useLinkScrapedTag({
    scrapeResult: tags,
    setScrapeResult: setTags,
    newObjects: newTags,
    setNewObjects: setNewTags,
  });

  async function handleLinkTagResult(tag: {
    create?: GQL.TagCreateInput;
    update?: GQL.TagUpdateInput;
  }) {
    if (tag.create) {
      await toastOperation(
        Toast,
        async () => {
          // create the new tag
          const result = await createTag({ variables: { input: tag.create! } });

          // adjust scrape result
          if (result.data?.tagCreate) {
            linkScrapedTag(
              result.data.tagCreate.id,
              result.data.tagCreate.name,
              linkedTag?.name ?? ""
            );
          }
        },
        intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl.formatMessage({ id: "tag" }).toLocaleLowerCase(),
          }
        )
      )();
    } else if (tag.update) {
      // link existing tag
      await toastOperation(
        Toast,
        async () => {
          const result = await updateTag({ variables: { input: tag.update! } });

          // adjust scrape result
          if (result.data?.tagUpdate) {
            linkScrapedTag(
              result.data.tagUpdate.id,
              result.data.tagUpdate.name,
              linkedTag?.name ?? ""
            );
          }
        },
        intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl.formatMessage({ id: "tag" }).toLocaleLowerCase(),
          }
        )
      )();
    }

    setLinkedTag(null);
  }

  const linkDialog = linkedTag ? (
    <CreateLinkTagDialog
      tag={linkedTag}
      onClose={handleLinkTagResult}
      endpoint={endpoint}
    />
  ) : null;

  const scrapedTagsRow = (
    <ScrapedTagsRow
      field="tags"
      title={intl.formatMessage({ id: "tags" })}
      result={tags}
      onChange={(value) => setTags(value)}
      newObjects={newTags}
      onCreateNew={createNewTag}
      onLinkExisting={(l) => setLinkedTag(l)}
    />
  );

  return {
    tags,
    newTags,
    linkDialog,
    scrapedTagsRow,
  };
}
