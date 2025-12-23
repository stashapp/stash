import { useToast } from "src/hooks/Toast";
import * as GQL from "src/core/generated-graphql";
import {
  useGroupCreate,
  usePerformerCreate,
  useStudioCreate,
  useTagCreate,
} from "src/core/StashService";
import { ObjectScrapeResult, ScrapeResult } from "./scrapeResult";
import { useIntl } from "react-intl";
import { scrapedPerformerToCreateInput } from "src/core/performers";
import { scrapedGroupToCreateInput } from "src/core/groups";

function useCreateObject<T>(
  entityTypeID: string,
  createFunc: (o: T) => Promise<void>
) {
  const Toast = useToast();
  const intl = useIntl();

  async function createNewObject(o: T) {
    try {
      await createFunc(o);

      Toast.success(
        intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl
              .formatMessage({ id: entityTypeID })
              .toLocaleLowerCase(),
          }
        )
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  return createNewObject;
}

interface IUseCreateNewStudioProps {
  scrapeResult: ObjectScrapeResult<GQL.ScrapedStudio>;
  setScrapeResult: (
    scrapeResult: ObjectScrapeResult<GQL.ScrapedStudio>
  ) => void;
  setNewObject: (newObject: GQL.ScrapedStudio | undefined) => void;
  endpoint?: string;
}

export function useCreateScrapedStudio(props: IUseCreateNewStudioProps) {
  const [createStudio] = useStudioCreate();

  const { scrapeResult, setScrapeResult, setNewObject } = props;

  async function createNewStudio(toCreate: GQL.ScrapedStudio) {
    const input: GQL.StudioCreateInput = {
      name: toCreate.name,
      urls: toCreate.urls,
      aliases:
        toCreate.aliases
          ?.split(",")
          .map((a) => a.trim())
          .filter((a) => a) || [],
      details: toCreate.details,
      image: toCreate.image,
      tag_ids: (toCreate.tags ?? [])
        .filter((t) => t.stored_id)
        .map((t) => t.stored_id!),
    };

    if (props.endpoint && toCreate.remote_site_id) {
      input.stash_ids = [
        {
          endpoint: props.endpoint,
          stash_id: toCreate.remote_site_id,
        },
      ];
    }

    const result = await createStudio({
      variables: {
        input,
      },
    });

    // set the new studio as the value
    setScrapeResult(
      scrapeResult.cloneWithValue({
        stored_id: result.data!.studioCreate!.id,
        name: toCreate.name,
      })
    );
    setNewObject(undefined);
  }

  return useCreateObject("studio", createNewStudio);
}

interface IUseCreateNewObjectProps<T> {
  scrapeResult: ScrapeResult<T[]>;
  setScrapeResult: (scrapeResult: ScrapeResult<T[]>) => void;
  newObjects: T[];
  setNewObjects: (newObject: T[]) => void;
  endpoint?: string;
}

export function useCreateScrapedPerformer(
  props: IUseCreateNewObjectProps<GQL.ScrapedPerformer>
) {
  const [createPerformer] = usePerformerCreate();

  const { scrapeResult, setScrapeResult, newObjects, setNewObjects } = props;

  async function createNewPerformer(toCreate: GQL.ScrapedPerformer) {
    const input = scrapedPerformerToCreateInput(toCreate, props.endpoint);

    const result = await createPerformer({
      variables: { input },
    });

    const newValue = [...(scrapeResult.newValue ?? [])];
    if (result.data?.performerCreate)
      newValue.push({
        stored_id: result.data.performerCreate.id,
        name: result.data.performerCreate.name,
      });

    // add the new performer to the new performers value
    const performerClone = scrapeResult.cloneWithValue(newValue);
    setScrapeResult(performerClone);

    // remove the performer from the list
    const newPerformersClone = newObjects.concat();
    const pIndex = newPerformersClone.findIndex(
      (p) => p.name === toCreate.name
    );
    if (pIndex === -1) throw new Error("Could not find performer to remove");

    newPerformersClone.splice(pIndex, 1);

    setNewObjects(newPerformersClone);
  }

  return useCreateObject("performer", createNewPerformer);
}

export function useCreateScrapedGroup(
  props: IUseCreateNewObjectProps<GQL.ScrapedGroup>
) {
  const { scrapeResult, setScrapeResult, newObjects, setNewObjects } = props;
  const [createGroup] = useGroupCreate();

  async function createNewGroup(toCreate: GQL.ScrapedGroup) {
    const input = scrapedGroupToCreateInput(toCreate);

    const result = await createGroup({
      variables: { input: input },
    });

    const newValue = [...(scrapeResult.newValue ?? [])];
    if (result.data?.groupCreate)
      newValue.push({
        stored_id: result.data.groupCreate.id,
        name: result.data.groupCreate.name,
      });

    // add the new object to the new object value
    const resultClone = scrapeResult.cloneWithValue(newValue);
    setScrapeResult(resultClone);

    // remove the object from the list
    const newObjectsClone = newObjects.concat();
    const pIndex = newObjectsClone.findIndex((p) => p.name === toCreate.name);
    if (pIndex === -1) throw new Error("Could not find group to remove");

    newObjectsClone.splice(pIndex, 1);

    setNewObjects(newObjectsClone);
  }

  return useCreateObject("group", createNewGroup);
}

export function useLinkScrapedTag(
  props: IUseCreateNewObjectProps<GQL.ScrapedTag>
) {
  const { scrapeResult, setScrapeResult, newObjects, setNewObjects } = props;

  function linkTag(id: string, matchedName: string, scrapedName: string) {
    const newValue = [...(scrapeResult.newValue ?? [])];
    newValue.push({
      stored_id: id,
      name: matchedName,
    });

    // add the new tag to the new tags value
    const tagClone = scrapeResult.cloneWithValue(newValue);
    setScrapeResult(tagClone);

    // remove the tag from the list
    const newTagsClone = newObjects.concat();
    const pIndex = newTagsClone.findIndex((p) => p.name === scrapedName);
    if (pIndex === -1) throw new Error("Could not find tag to remove");

    newTagsClone.splice(pIndex, 1);

    setNewObjects(newTagsClone);
  }

  return linkTag;
}

export function useCreateScrapedTag(
  props: IUseCreateNewObjectProps<GQL.ScrapedTag>
) {
  const [createTag] = useTagCreate();
  const linkTag = useLinkScrapedTag(props);

  async function createNewTag(toCreate: GQL.ScrapedTag) {
    const input: GQL.TagCreateInput = {
      name: toCreate.name ?? "",
    };

    if (props.endpoint && toCreate.remote_site_id) {
      input.stash_ids = [
        {
          endpoint: props.endpoint,
          stash_id: toCreate.remote_site_id,
        },
      ];
    }

    const result = await createTag({
      variables: { input },
    });

    if (result.data?.tagCreate)
      linkTag(
        result.data.tagCreate.id,
        result.data.tagCreate.name,
        toCreate.name ?? ""
      );
  }

  return useCreateObject("tag", createNewTag);
}
