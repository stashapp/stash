import { useToast } from "src/hooks/Toast";
import * as GQL from "src/core/generated-graphql";
import {
  useMovieCreate,
  usePerformerCreate,
  useStudioCreate,
  useTagCreate,
} from "src/core/StashService";
import { ScrapeResult } from "./scrapeResult";
import { useIntl } from "react-intl";
import { scrapedPerformerToCreateInput } from "src/core/performers";
import { scrapedMovieToCreateInput } from "src/core/movies";

function useCreateObject<T>(
  entityTypeID: string,
  createFunc: (o: T) => Promise<void>
) {
  const Toast = useToast();
  const intl = useIntl();

  async function createNewObject(o: T) {
    try {
      await createFunc(o);

      Toast.success({
        content: intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl
              .formatMessage({ id: entityTypeID })
              .toLocaleLowerCase(),
          }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  return createNewObject;
}

interface IUseCreateNewStudioProps {
  scrapeResult: ScrapeResult<string>;
  setScrapeResult: (scrapeResult: ScrapeResult<string>) => void;
  setNewObject: (newObject: GQL.ScrapedStudio | undefined) => void;
}

export function useCreateScrapedStudio(props: IUseCreateNewStudioProps) {
  const [createStudio] = useStudioCreate();

  const { scrapeResult, setScrapeResult, setNewObject } = props;

  async function createNewStudio(toCreate: GQL.ScrapedStudio) {
    const result = await createStudio({
      variables: {
        input: {
          name: toCreate.name,
          url: toCreate.url,
        },
      },
    });

    // set the new studio as the value
    setScrapeResult(scrapeResult.cloneWithValue(result.data!.studioCreate!.id));
    setNewObject(undefined);
  }

  return useCreateObject("studio", createNewStudio);
}

interface IUseCreateNewPerformerProps {
  scrapeResult: ScrapeResult<GQL.ScrapedPerformer[]>;
  setScrapeResult: (scrapeResult: ScrapeResult<GQL.ScrapedPerformer[]>) => void;
  newObjects: GQL.ScrapedPerformer[];
  setNewObjects: (newObject: GQL.ScrapedPerformer[]) => void;
}

export function useCreateScrapedPerformer(props: IUseCreateNewPerformerProps) {
  const [createPerformer] = usePerformerCreate();

  const { scrapeResult, setScrapeResult, newObjects, setNewObjects } = props;

  async function createNewPerformer(toCreate: GQL.ScrapedPerformer) {
    const input = scrapedPerformerToCreateInput(toCreate);

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

interface IUseCreateNewObjectIDListProps<
  T extends { name?: string | undefined | null }
> {
  scrapeResult: ScrapeResult<string[]>;
  setScrapeResult: (scrapeResult: ScrapeResult<string[]>) => void;
  newObjects: T[];
  setNewObjects: (newObject: T[]) => void;
}

function useCreateNewObjectIDList<
  T extends { name?: string | undefined | null }
>(
  entityTypeID: string,
  props: IUseCreateNewObjectIDListProps<T>,
  createObject: (toCreate: T) => Promise<string>
) {
  const { scrapeResult, setScrapeResult, newObjects, setNewObjects } = props;

  async function createNewObject(toCreate: T) {
    const newID = await createObject(toCreate);

    // add the new object to the new objects value
    const newResult = scrapeResult.cloneWithValue(scrapeResult.newValue);
    if (!newResult.newValue) {
      newResult.newValue = [];
    }
    newResult.newValue.push(newID);
    setScrapeResult(newResult);

    // remove the object from the list
    const newObjectsClone = newObjects.concat();
    const pIndex = newObjectsClone.findIndex((p) => p.name === toCreate.name);
    if (pIndex === -1) throw new Error("Could not find object to remove");
    newObjectsClone.splice(pIndex, 1);

    setNewObjects(newObjectsClone);
  }

  return useCreateObject(entityTypeID, createNewObject);
}

export function useCreateScrapedMovie(
  props: IUseCreateNewObjectIDListProps<GQL.ScrapedMovie>
) {
  const [createMovie] = useMovieCreate();

  async function createNewMovie(toCreate: GQL.ScrapedMovie) {
    const movieInput = scrapedMovieToCreateInput(toCreate);
    const result = await createMovie({
      variables: { input: movieInput },
    });

    return result.data?.movieCreate?.id ?? "";
  }

  return useCreateNewObjectIDList("movie", props, createNewMovie);
}

export function useCreateScrapedTag(
  props: IUseCreateNewObjectIDListProps<GQL.ScrapedTag>
) {
  const [createTag] = useTagCreate();

  async function createNewTag(toCreate: GQL.ScrapedTag) {
    const tagInput: GQL.TagCreateInput = { name: toCreate.name ?? "" };
    const result = await createTag({
      variables: {
        input: tagInput,
      },
    });

    return result.data?.tagCreate?.id ?? "";
  }

  return useCreateNewObjectIDList("tag", props, createNewTag);
}
