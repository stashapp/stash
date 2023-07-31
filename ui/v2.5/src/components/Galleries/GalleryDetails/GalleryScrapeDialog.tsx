import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import {
  StudioSelect,
  PerformerSelect,
  TagSelect,
} from "src/components/Shared/Select";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapeDialog,
  ScrapeDialogRow,
  ScrapeResult,
  ScrapedInputGroupRow,
  ScrapedTextAreaRow,
} from "src/components/Shared/ScrapeDialog";
import clone from "lodash-es/clone";
import {
  useStudioCreate,
  usePerformerCreate,
  useTagCreate,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { scrapedPerformerToCreateInput } from "src/core/performers";

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
  newStudio?: GQL.ScrapedStudio,
  onCreateNew?: (value: GQL.ScrapedStudio) => void
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
      onCreateNew={() => {
        if (onCreateNew && newStudio) onCreateNew(newStudio);
      }}
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
  newPerformers: GQL.ScrapedPerformer[],
  onCreateNew?: (value: GQL.ScrapedPerformer) => void
) {
  const performersCopy = newPerformers.map((p) => {
    const name: string = p.name ?? "";
    return { ...p, name };
  });

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
      newValues={performersCopy}
      onCreateNew={(i) => {
        if (onCreateNew) onCreateNew(newPerformers[i]);
      }}
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
  newTags: GQL.ScrapedTag[],
  onCreateNew?: (value: GQL.ScrapedTag) => void
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
      onCreateNew={(i) => {
        if (onCreateNew) onCreateNew(newTags[i]);
      }}
    />
  );
}

interface IGalleryScrapeDialogProps {
  gallery: Partial<GQL.GalleryUpdateInput>;
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

  const [performers, setPerformers] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.gallery.performer_ids),
      mapStoredIdObjects(props.scraped.performers ?? undefined)
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

  const [createStudio] = useStudioCreate();
  const [createPerformer] = usePerformerCreate();
  const [createTag] = useTagCreate();

  const Toast = useToast();

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
            <FormattedMessage
              id="actions.created_entity"
              values={{
                entity_type: intl.formatMessage({ id: "studio" }),
                entity_name: <b>{toCreate.name}</b>,
              }}
            />
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
            <FormattedMessage
              id="actions.created_entity"
              values={{
                entity_type: intl.formatMessage({ id: "performer" }),
                entity_name: <b>{toCreate.name}</b>,
              }}
            />
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
            <FormattedMessage
              id="actions.created_entity"
              values={{
                entity_type: intl.formatMessage({ id: "tag" }),
                entity_name: <b>{toCreate.name}</b>,
              }}
            />
          </span>
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
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
      performers: performers.getNewValue()?.map((p) => {
        return {
          stored_id: p,
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
