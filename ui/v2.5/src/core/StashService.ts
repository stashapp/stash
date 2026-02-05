import {
  ApolloCache,
  DocumentNode,
  FetchResult,
  NetworkStatus,
  useQuery,
} from "@apollo/client";
import { Modifiers } from "@apollo/client/cache";
import {
  isField,
  getQueryDefinition,
  StoreObject,
} from "@apollo/client/utilities";
import { ListFilterModel } from "../models/list-filter/filter";
import * as GQL from "./generated-graphql";

import { createClient } from "./createClient";
import { Client } from "graphql-ws";
import { useEffect, useState } from "react";

const { client, wsClient, cache: clientCache } = createClient();

export const getClient = () => client;
export const getWSClient = () => wsClient;

export function useWSState(ws: Client) {
  const [state, setState] = useState<"connecting" | "connected" | "error">(
    "connecting"
  );

  useEffect(() => {
    const disposeConnected = ws.on("connected", () => {
      setState("connected");
    });

    const disposeError = ws.on("error", () => {
      setState("error");
    });

    return () => {
      disposeConnected();
      disposeError();
    };
  }, [ws]);

  return { state };
}

// Evicts cached results for the given queries.
// Will also call a cache GC afterwards.
export function evictQueries(
  cache: ApolloCache<unknown>,
  queries: DocumentNode[]
) {
  const fields: Modifiers = {};
  for (const query of queries) {
    const { selections } = getQueryDefinition(query).selectionSet;
    for (const field of selections) {
      if (!isField(field)) continue;
      const keyName = field.name.value;
      fields[keyName] = (_value, { DELETE }) => DELETE;
    }
  }

  cache.modify({ fields });

  // evictQueries is usually called at the end of
  // an update function - so call a GC here
  cache.gc();
}

/**
 * Evicts fields from all objects of a given type.
 *
 * @param input   a map from typename -> list of field names to evict
 * @param ignore  optionally specify a cache id to ignore and not modify
 */
function evictTypeFields(
  cache: ApolloCache<Record<string, StoreObject>>,
  input: Record<string, string[]>,
  ignore?: string
) {
  const data = cache.extract();
  for (const key in data) {
    if (ignore?.includes(key)) continue;

    const obj = data[key];
    const typename = obj.__typename;

    if (typename && input[typename]) {
      const modifiers: Modifiers = {};
      for (const field of input[typename]) {
        modifiers[field] = (_value, { DELETE }) => DELETE;
      }
      cache.modify({
        id: key,
        fields: modifiers,
      });
    }
  }
}

// Deletes obj from the cache, and sets the
// cached result of the given query to null.
// Use with "Destroy" mutations.
function deleteObject(
  cache: ApolloCache<unknown>,
  obj: StoreObject,
  query: DocumentNode
) {
  const field = getQueryDefinition(query).selectionSet.selections[0];
  if (!isField(field)) return;
  const keyName = field.name.value;

  cache.writeQuery({
    query,
    variables: { id: obj.id },
    data: { [keyName]: null },
  });
  cache.evict({ id: cache.identify(obj) });
}

export function isLoading(networkStatus: NetworkStatus) {
  // useQuery hook loading field only returns true when initially loading the query
  // and not during subsequent fetches
  return (
    networkStatus === NetworkStatus.loading ||
    networkStatus === NetworkStatus.fetchMore ||
    networkStatus === NetworkStatus.refetch
  );
}

/// Object queries

export const useFindScene = (id: string) => {
  const skip = id === "new" || id === "";
  return GQL.useFindSceneQuery({ variables: { id }, skip });
};

export const useSceneStreams = (id: string) =>
  GQL.useSceneStreamsQuery({ variables: { id } });

export const useFindScenes = (filter?: ListFilterModel) =>
  GQL.useFindScenesQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      scene_filter: filter?.makeFilter(),
    },
  });

export const queryFindScenes = (filter: ListFilterModel) =>
  client.query<GQL.FindScenesQuery>({
    query: GQL.FindScenesDocument,
    variables: {
      filter: filter.makeFindFilter(),
      scene_filter: filter.makeFilter(),
    },
  });

export const queryFindScenesByID = (sceneIDs: number[]) =>
  client.query<GQL.FindScenesQuery>({
    query: GQL.FindScenesDocument,
    variables: {
      scene_ids: sceneIDs,
    },
  });

export const queryFindScenesForSelect = (filter: ListFilterModel) =>
  client.query<GQL.FindScenesForSelectQuery>({
    query: GQL.FindScenesForSelectDocument,
    variables: {
      filter: filter.makeFindFilter(),
      scene_filter: filter.makeFilter(),
    },
  });

export const queryFindScenesByIDForSelect = (sceneIDs: string[]) =>
  client.query<GQL.FindScenesForSelectQuery>({
    query: GQL.FindScenesForSelectDocument,
    variables: {
      ids: sceneIDs,
    },
  });

export const querySceneByPathRegex = (filter: GQL.FindFilterType) =>
  client.query<GQL.FindScenesByPathRegexQuery>({
    query: GQL.FindScenesByPathRegexDocument,
    variables: { filter },
  });

export const useFindImage = (id: string) =>
  GQL.useFindImageQuery({ variables: { id } });

export const useFindImages = (filter?: ListFilterModel) =>
  GQL.useFindImagesQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      image_filter: filter?.makeFilter(),
    },
  });

export const useFindImagesMetadata = (filter?: ListFilterModel) =>
  GQL.useFindImagesMetadataQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      image_filter: filter?.makeFilter(),
    },
  });

export const queryFindImages = (filter: ListFilterModel) =>
  client.query<GQL.FindImagesQuery>({
    query: GQL.FindImagesDocument,
    variables: {
      filter: filter.makeFindFilter(),
      image_filter: filter.makeFilter(),
    },
  });

export const useFindGroup = (id: string) => {
  const skip = id === "new" || id === "";
  return GQL.useFindGroupQuery({ variables: { id }, skip });
};

export const useFindGroups = (filter?: ListFilterModel) =>
  GQL.useFindGroupsQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      group_filter: filter?.makeFilter(),
    },
  });

export const queryFindGroups = (filter: ListFilterModel) =>
  client.query<GQL.FindGroupsQuery>({
    query: GQL.FindGroupsDocument,
    variables: {
      filter: filter.makeFindFilter(),
      group_filter: filter.makeFilter(),
    },
  });

export const queryFindGroupsByIDForSelect = (groupIDs: string[]) =>
  client.query<GQL.FindGroupsForSelectQuery>({
    query: GQL.FindGroupsForSelectDocument,
    variables: {
      ids: groupIDs,
    },
  });

export const queryFindGroupsForSelect = (filter: ListFilterModel) =>
  client.query<GQL.FindGroupsForSelectQuery>({
    query: GQL.FindGroupsForSelectDocument,
    variables: {
      filter: filter.makeFindFilter(),
      group_filter: filter.makeFilter(),
    },
  });

export const useFindSceneMarkers = (filter?: ListFilterModel) =>
  GQL.useFindSceneMarkersQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      scene_marker_filter: filter?.makeFilter(),
    },
  });

export const queryFindSceneMarkers = (filter: ListFilterModel) =>
  client.query<GQL.FindSceneMarkersQuery>({
    query: GQL.FindSceneMarkersDocument,
    variables: {
      filter: filter.makeFindFilter(),
      scene_marker_filter: filter.makeFilter(),
    },
  });

export const useMarkerStrings = () => GQL.useMarkerStringsQuery();

export const useFindGallery = (id: string) => {
  const skip = id === "new" || id === "";
  return GQL.useFindGalleryQuery({ variables: { id }, skip });
};

export const useFindGalleryImageID = (id: string, index: number) => {
  return GQL.useFindGalleryImageIdQuery({ variables: { id, index } });
};

export const useFindGalleries = (filter?: ListFilterModel) =>
  GQL.useFindGalleriesQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      gallery_filter: filter?.makeFilter(),
    },
  });

export const queryFindGalleries = (filter: ListFilterModel) =>
  client.query<GQL.FindGalleriesQuery>({
    query: GQL.FindGalleriesDocument,
    variables: {
      filter: filter.makeFindFilter(),
      gallery_filter: filter.makeFilter(),
    },
  });

export const queryFindGalleriesForSelect = (filter: ListFilterModel) =>
  client.query<GQL.FindGalleriesForSelectQuery>({
    query: GQL.FindGalleriesForSelectDocument,
    variables: {
      filter: filter.makeFindFilter(),
      gallery_filter: filter.makeFilter(),
    },
  });

export const queryFindGalleriesByIDForSelect = (galleryIDs: string[]) =>
  client.query<GQL.FindGalleriesForSelectQuery>({
    query: GQL.FindGalleriesForSelectDocument,
    variables: {
      ids: galleryIDs,
    },
  });

export const useFindPerformer = (id: string) => {
  const skip = id === "new" || id === "";
  return GQL.useFindPerformerQuery({ variables: { id }, skip });
};

export const queryFindPerformer = (id: string) =>
  client.query<GQL.FindPerformerQuery>({
    query: GQL.FindPerformerDocument,
    variables: { id },
  });

export const useFindPerformers = (filter?: ListFilterModel) =>
  GQL.useFindPerformersQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      performer_filter: filter?.makeFilter(),
    },
  });

export const queryFindPerformers = (filter: ListFilterModel) =>
  client.query<GQL.FindPerformersQuery>({
    query: GQL.FindPerformersDocument,
    variables: {
      filter: filter.makeFindFilter(),
      performer_filter: filter.makeFilter(),
    },
  });

export const queryFindPerformersByID = (performerIDs: number[]) =>
  client.query<GQL.FindPerformersQuery>({
    query: GQL.FindPerformersDocument,
    variables: {
      performer_ids: performerIDs,
    },
  });

export const queryFindPerformersByIDForSelect = (performerIDs: string[]) =>
  client.query<GQL.FindPerformersForSelectQuery>({
    query: GQL.FindPerformersForSelectDocument,
    variables: {
      ids: performerIDs,
    },
  });

export const queryFindPerformersForSelect = (filter: ListFilterModel) =>
  client.query<GQL.FindPerformersForSelectQuery>({
    query: GQL.FindPerformersForSelectDocument,
    variables: {
      filter: filter.makeFindFilter(),
      performer_filter: filter.makeFilter(),
    },
  });

export const useFindStudio = (id: string) => {
  const skip = id === "new" || id === "";
  return GQL.useFindStudioQuery({ variables: { id }, skip });
};

export const queryFindStudio = (id: string) =>
  client.query<GQL.FindStudioQuery>({
    query: GQL.FindStudioDocument,
    variables: { id },
  });

export const useFindStudios = (filter?: ListFilterModel) =>
  GQL.useFindStudiosQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      studio_filter: filter?.makeFilter(),
    },
  });

export const queryFindStudios = (filter: ListFilterModel) =>
  client.query<GQL.FindStudiosQuery>({
    query: GQL.FindStudiosDocument,
    variables: {
      filter: filter.makeFindFilter(),
      studio_filter: filter.makeFilter(),
    },
  });

export const queryFindStudiosByIDForSelect = (studioIDs: string[]) =>
  client.query<GQL.FindStudiosForSelectQuery>({
    query: GQL.FindStudiosForSelectDocument,
    variables: {
      ids: studioIDs,
    },
  });

export const queryFindStudiosForSelect = (filter: ListFilterModel) =>
  client.query<GQL.FindStudiosForSelectQuery>({
    query: GQL.FindStudiosForSelectDocument,
    variables: {
      filter: filter.makeFindFilter(),
      studio_filter: filter.makeFilter(),
    },
  });

export const useFindTag = (id: string) => {
  const skip = id === "new" || id === "";
  return GQL.useFindTagQuery({ variables: { id }, skip });
};

export const queryFindTag = (id: string) =>
  client.query<GQL.FindTagQuery>({
    query: GQL.FindTagDocument,
    variables: { id },
  });

export const useFindTags = (filter?: ListFilterModel) =>
  GQL.useFindTagsQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      tag_filter: filter?.makeFilter(),
    },
  });

// Optimized query for tag list page - excludes expensive recursive *_count_all fields
export const useFindTagsForList = (filter?: ListFilterModel) =>
  GQL.useFindTagsForListQuery({
    skip: filter === undefined,
    variables: {
      filter: filter?.makeFindFilter(),
      tag_filter: filter?.makeFilter(),
    },
  });

export const queryFindTags = (filter: ListFilterModel) =>
  client.query<GQL.FindTagsQuery>({
    query: GQL.FindTagsDocument,
    variables: {
      filter: filter.makeFindFilter(),
      tag_filter: filter.makeFilter(),
    },
  });

// Optimized query for tag list page
export const queryFindTagsForList = (filter: ListFilterModel) =>
  client.query<GQL.FindTagsForListQuery>({
    query: GQL.FindTagsForListDocument,
    variables: {
      filter: filter.makeFindFilter(),
      tag_filter: filter.makeFilter(),
    },
  });

export const queryFindTagsByIDForSelect = (tagIDs: string[]) =>
  client.query<GQL.FindTagsForSelectQuery>({
    query: GQL.FindTagsForSelectDocument,
    variables: {
      ids: tagIDs,
    },
  });

export const queryFindTagsForSelect = (filter: ListFilterModel) =>
  client.query<GQL.FindTagsForSelectQuery>({
    query: GQL.FindTagsForSelectDocument,
    variables: {
      filter: filter.makeFindFilter(),
      tag_filter: filter.makeFilter(),
    },
  });

export const useFindSavedFilter = (id: string) =>
  GQL.useFindSavedFilterQuery({
    variables: { id },
  });

export const useFindSavedFilters = (mode?: GQL.FilterMode) =>
  GQL.useFindSavedFiltersQuery({
    variables: { mode },
  });

/// Object Mutations

// Increases/decreases the given field of the Stats query by diff
function updateStats(cache: ApolloCache<unknown>, field: string, diff: number) {
  cache.modify({
    fields: {
      stats(value) {
        return {
          ...value,
          [field]: value[field] + diff,
        };
      },
    },
  });
}

function updateO(
  cache: ApolloCache<unknown>,
  typename: string,
  id: string,
  updatedOCount: number
) {
  cache.modify({
    id: cache.identify({ __typename: typename, id }),
    fields: {
      o_counter() {
        return updatedOCount;
      },
    },
  });
}

const sceneMutationImpactedTypeFields = {
  Group: ["scenes", "scene_count"],
  Gallery: ["scenes"],
  Performer: [
    "scenes",
    "scene_count",
    "groups",
    "group_count",
    "performer_count",
  ],
  Studio: ["scene_count", "performer_count"],
  Tag: ["scene_count"],
};

const sceneMutationImpactedQueries = [
  GQL.FindScenesDocument, // various filters
  GQL.FindGroupsDocument, // is missing scenes
  GQL.FindGalleriesDocument, // is missing scenes
  GQL.FindPerformersDocument, // filter by scene count
  GQL.FindStudiosDocument, // filter by scene count
  GQL.FindTagsDocument, // filter by scene count
];

export const mutateCreateScene = (input: GQL.SceneCreateInput) =>
  client.mutate<GQL.SceneCreateMutation>({
    mutation: GQL.SceneCreateDocument,
    variables: { input },
    update(cache, result) {
      const scene = result.data?.sceneCreate;
      if (!scene) return;

      // update stats
      updateStats(cache, "scene_count", 1);

      // if we're reassigning files, refetch files from other scenes
      if (input.file_ids?.length) {
        const obj = { __typename: "Scene", id: scene.id };
        evictTypeFields(
          cache,
          {
            ...sceneMutationImpactedTypeFields,
            Scene: ["files"],
          },
          cache.identify(obj) // don't evict this scene
        );
      } else {
        evictTypeFields(cache, sceneMutationImpactedTypeFields);
      }

      evictQueries(cache, sceneMutationImpactedQueries);
    },
  });

export const useSceneUpdate = () =>
  GQL.useSceneUpdateMutation({
    update(cache, result) {
      if (!result.data?.sceneUpdate) return;

      evictTypeFields(cache, sceneMutationImpactedTypeFields);
      evictQueries(cache, sceneMutationImpactedQueries);
    },
  });

export const useBulkSceneUpdate = (input: GQL.BulkSceneUpdateInput) =>
  GQL.useBulkSceneUpdateMutation({
    variables: { input },
    update(cache, result) {
      if (!result.data?.bulkSceneUpdate) return;

      evictTypeFields(cache, sceneMutationImpactedTypeFields);
      evictQueries(cache, sceneMutationImpactedQueries);
    },
  });

export const useScenesUpdate = (input: GQL.SceneUpdateInput[]) =>
  GQL.useScenesUpdateMutation({
    variables: { input },
    update(cache, result) {
      if (!result.data?.scenesUpdate) return;

      evictTypeFields(cache, sceneMutationImpactedTypeFields);
      evictQueries(cache, sceneMutationImpactedQueries);
    },
  });

export const useSceneDestroy = (input: GQL.SceneDestroyInput) =>
  GQL.useSceneDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.sceneDestroy) return;

      const obj = { __typename: "Scene", id: input.id };
      deleteObject(cache, obj, GQL.FindSceneDocument);

      evictTypeFields(cache, sceneMutationImpactedTypeFields);
      evictQueries(cache, [
        ...sceneMutationImpactedQueries,
        GQL.FindSceneMarkersDocument, // filter by scene tags
        GQL.StatsDocument, // scenes size, scene count, etc
      ]);
    },
  });

export const useScenesDestroy = (input: GQL.ScenesDestroyInput) =>
  GQL.useScenesDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.scenesDestroy) return;

      for (const id of input.ids) {
        const obj = { __typename: "Scene", id };
        deleteObject(cache, obj, GQL.FindSceneDocument);
      }

      evictTypeFields(cache, sceneMutationImpactedTypeFields);
      evictQueries(cache, [
        ...sceneMutationImpactedQueries,
        GQL.FindSceneMarkersDocument, // filter by scene tags
        GQL.StatsDocument, // scenes size, scene count, etc
      ]);
    },
  });

export const useSceneIncrementO = (id: string) =>
  GQL.useSceneAddOMutation({
    variables: { id },
    update(cache, result, { variables }) {
      // this is not perfectly accurate, the time is set server-side
      // it isn't even displayed anywhere in the UI anyway
      const at = new Date().toISOString();

      const mutationResult = result.data?.sceneAddO;
      if (!mutationResult || !variables) return;

      const { history } = mutationResult;
      const { times } = variables;
      const timeArray = !times ? [at] : Array.isArray(times) ? times : [times];

      const scene = cache.readFragment<GQL.SlimSceneDataFragment>({
        id: cache.identify({ __typename: "Scene", id }),
        fragment: GQL.SlimSceneDataFragmentDoc,
        fragmentName: "SlimSceneData",
      });

      if (scene) {
        // if we have the scene, update performer o_counters manually
        for (const performer of scene.performers) {
          cache.modify({
            id: cache.identify(performer),
            fields: {
              o_counter(value) {
                return value + timeArray.length;
              },
            },
          });
        }
      } else {
        // else refresh all performer o_counters
        evictTypeFields(cache, {
          Performer: ["o_counter"],
        });
      }

      updateStats(cache, "total_o_count", timeArray.length);

      cache.modify({
        id: cache.identify({ __typename: "Scene", id }),
        fields: {
          o_history() {
            return history;
          },
        },
      });

      updateO(cache, "Scene", id, history.length);
      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by o_counter
        GQL.FindPerformersDocument, // filter by o_counter
      ]);
    },
  });

export const useSceneDecrementO = (id: string) =>
  GQL.useSceneDeleteOMutation({
    variables: { id },
    update(cache, result, { variables }) {
      const mutationResult = result.data?.sceneDeleteO;
      if (!mutationResult || !variables) return;

      const { history } = mutationResult;
      const { times } = variables;
      const timeArray = !times ? null : Array.isArray(times) ? times : [times];

      const scene = cache.readFragment<GQL.SlimSceneDataFragment>({
        id: cache.identify({ __typename: "Scene", id }),
        fragment: GQL.SlimSceneDataFragmentDoc,
        fragmentName: "SlimSceneData",
      });

      if (scene) {
        // if we have the scene, update performer o_counters manually
        for (const performer of scene.performers) {
          cache.modify({
            id: cache.identify(performer),
            fields: {
              o_counter(value) {
                return value - (timeArray?.length ?? 1);
              },
            },
          });
        }
      } else {
        // else refresh all performer o_counters
        evictTypeFields(cache, {
          Performer: ["o_counter"],
        });
      }

      updateStats(cache, "total_o_count", -(timeArray?.length ?? 1));

      cache.modify({
        id: cache.identify({ __typename: "Scene", id }),
        fields: {
          o_history() {
            return history;
          },
        },
      });

      updateO(cache, "Scene", id, history.length);
      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by o_counter
        GQL.FindPerformersDocument, // filter by o_counter
      ]);
    },
  });

export const useSceneResetO = (id: string) =>
  GQL.useSceneResetOMutation({
    variables: { id },
    update(cache, result) {
      const updatedOCount = result.data?.sceneResetO;
      if (updatedOCount === undefined) return;

      const scene = cache.readFragment<GQL.SlimSceneDataFragment>({
        id: cache.identify({ __typename: "Scene", id }),
        fragment: GQL.SlimSceneDataFragmentDoc,
        fragmentName: "SlimSceneData",
      });

      if (scene) {
        // if we have the scene, update performer o_counters manually
        const old_count = scene.o_counter ?? 0;
        for (const performer of scene.performers) {
          cache.modify({
            id: cache.identify(performer),
            fields: {
              o_counter(value) {
                return value - old_count;
              },
            },
          });
        }
        updateStats(cache, "total_o_count", -old_count);
      } else {
        // else refresh all performer o_counters
        evictTypeFields(cache, {
          Performer: ["o_counter"],
        });
        // also refresh stats total_o_count
        cache.modify({
          fields: {
            stats: (value) => ({
              ...value,
              total_o_count: undefined,
            }),
          },
        });
      }

      cache.modify({
        id: cache.identify({ __typename: "Scene", id }),
        fields: {
          o_history() {
            const ret: string[] = [];
            return ret;
          },
        },
      });

      updateO(cache, "Scene", id, updatedOCount);
      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by o_counter
        GQL.FindPerformersDocument, // filter by o_counter
      ]);
    },
  });

export const useSceneResetActivity = (
  id: string,
  reset_resume: boolean,
  reset_duration: boolean
) =>
  GQL.useSceneResetActivityMutation({
    variables: { id, reset_resume, reset_duration },
    update(cache, result) {
      if (!result.data?.sceneResetActivity) return;

      evictTypeFields(cache, sceneMutationImpactedTypeFields);
      evictQueries(cache, sceneMutationImpactedQueries);
    },
  });

export const useSceneGenerateScreenshot = () =>
  GQL.useSceneGenerateScreenshotMutation();

export const mutateSceneSetPrimaryFile = (id: string, fileID: string) =>
  client.mutate<GQL.SceneUpdateMutation>({
    mutation: GQL.SceneUpdateDocument,
    variables: {
      input: {
        id,
        primary_file_id: fileID,
      },
    },
    update(cache, result) {
      if (!result.data?.sceneUpdate) return;

      evictQueries(cache, [
        GQL.FindScenesDocument, // sort by primary basename when missing title
      ]);
    },
  });

export const mutateSceneAssignFile = (sceneID: string, fileID: string) =>
  client.mutate<GQL.SceneAssignFileMutation>({
    mutation: GQL.SceneAssignFileDocument,
    variables: {
      input: {
        scene_id: sceneID,
        file_id: fileID,
      },
    },
    update(cache, result) {
      if (!result.data?.sceneAssignFile) return;

      // refetch target scene
      cache.evict({
        id: cache.identify({ __typename: "Scene", id: sceneID }),
      });

      // refetch files of the scene the file was previously assigned to
      evictTypeFields(cache, { Scene: ["files"] });

      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by file count
      ]);
    },
  });

export const mutateSceneMerge = (
  destination: string,
  source: string[],
  values: GQL.SceneUpdateInput,
  includeViewHistory: boolean,
  includeOHistory: boolean
) =>
  client.mutate<GQL.SceneMergeMutation>({
    mutation: GQL.SceneMergeDocument,
    variables: {
      input: {
        source,
        destination,
        values,
        play_history: includeViewHistory,
        o_history: includeOHistory,
      },
    },
    update(cache, result) {
      if (!result.data?.sceneMerge) return;

      for (const id of source) {
        const obj = { __typename: "Scene", id };
        deleteObject(cache, obj, GQL.FindSceneDocument);
      }

      cache.evict({
        id: cache.identify({ __typename: "Scene", id: destination }),
      });

      evictTypeFields(cache, sceneMutationImpactedTypeFields);
      evictQueries(cache, [
        ...sceneMutationImpactedQueries,
        GQL.StatsDocument, // scenes size, scene count, etc
      ]);
    },
  });

export const useSceneSaveActivity = () =>
  GQL.useSceneSaveActivityMutation({
    update(cache, result, { variables }) {
      if (!result.data?.sceneSaveActivity || !variables) return;

      const { id, playDuration, resume_time: resumeTime } = variables;

      cache.modify({
        id: cache.identify({ __typename: "Scene", id }),
        fields: {
          resume_time() {
            return resumeTime ?? null;
          },
          play_duration(value) {
            return value + playDuration;
          },
        },
      });

      if (playDuration) {
        updateStats(cache, "total_play_duration", playDuration);
      }

      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by play duration
      ]);
    },
  });

export const useSceneIncrementPlayCount = () =>
  GQL.useSceneAddPlayMutation({
    update(cache, result, { variables }) {
      const mutationResult = result.data?.sceneAddPlay;

      if (!mutationResult || !variables) return;

      const { history } = mutationResult;
      const { id } = variables;

      let lastPlayCount = 0;
      const playCount = history.length;

      cache.modify({
        id: cache.identify({ __typename: "Scene", id }),
        fields: {
          play_count(value) {
            lastPlayCount = value;
            return history.length;
          },
          last_played_at() {
            // assume only one entry - or the first is the most recent
            return history[0];
          },
          play_history() {
            return history;
          },
        },
      });

      updateStats(cache, "total_play_count", playCount - lastPlayCount);
      if (lastPlayCount === 0) {
        updateStats(cache, "scenes_played", 1);
      }

      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by play count
      ]);
    },
  });

export const useSceneDecrementPlayCount = () =>
  GQL.useSceneDeletePlayMutation({
    update(cache, result, { variables }) {
      const mutationResult = result.data?.sceneDeletePlay;

      if (!mutationResult || !variables) return;

      const { history } = mutationResult;
      const { id, times } = variables;
      const timeArray = !times ? null : Array.isArray(times) ? times : [times];
      const nRemoved = timeArray?.length ?? 1;

      let lastPlayCount = 0;
      let lastPlayedAt: string | null = null;
      const playCount = history.length;

      cache.modify({
        id: cache.identify({ __typename: "Scene", id }),
        fields: {
          play_count(value) {
            lastPlayCount = value;
            return playCount;
          },
          play_history() {
            if (history.length > 0) {
              lastPlayedAt = history[0];
            }
            return history;
          },
        },
      });

      cache.modify({
        id: cache.identify({ __typename: "Scene", id }),
        fields: {
          last_played_at() {
            return lastPlayedAt;
          },
        },
      });

      if (lastPlayCount > 0) {
        updateStats(
          cache,
          "total_play_count",
          nRemoved > lastPlayCount ? -lastPlayCount : -nRemoved
        );
      }
      if (lastPlayCount - nRemoved <= 0) {
        updateStats(cache, "scenes_played", -1);
      }

      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by play count
      ]);
    },
  });

export const useSceneResetPlayCount = () =>
  GQL.useSceneResetPlayCountMutation({
    update(cache, result, { variables }) {
      if (!variables) return;

      let lastPlayCount = 0;
      cache.modify({
        id: cache.identify({ __typename: "Scene", id: variables.id }),
        fields: {
          play_count(value) {
            lastPlayCount = value;
            return 0;
          },
          play_history() {
            const ret: string[] = [];
            return ret;
          },
          last_played_at() {
            return null;
          },
        },
      });

      if (lastPlayCount > 0) {
        updateStats(cache, "total_play_count", -lastPlayCount);
      }
      if (lastPlayCount > 0) {
        updateStats(cache, "scenes_played", -1);
      }

      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by play count
      ]);
    },
  });

const imageMutationImpactedTypeFields = {
  Gallery: ["images", "image_count"],
  Performer: ["image_count", "performer_count"],
  Studio: ["image_count", "performer_count"],
  Tag: ["image_count"],
};

const imageMutationImpactedQueries = [
  GQL.FindImagesDocument, // various filters
  GQL.FindGalleriesDocument, // filter by image count
  GQL.FindPerformersDocument, // filter by image count
  GQL.FindStudiosDocument, // filter by image count
  GQL.FindTagsDocument, // filter by image count
];

export const useImageUpdate = () =>
  GQL.useImageUpdateMutation({
    update(cache, result) {
      if (!result.data?.imageUpdate) return;

      evictTypeFields(cache, imageMutationImpactedTypeFields);
      evictQueries(cache, imageMutationImpactedQueries);
    },
  });

export const useBulkImageUpdate = () =>
  GQL.useBulkImageUpdateMutation({
    update(cache, result) {
      if (!result.data?.bulkImageUpdate) return;

      evictTypeFields(cache, imageMutationImpactedTypeFields);
      evictQueries(cache, imageMutationImpactedQueries);
    },
  });

export const useImagesDestroy = (input: GQL.ImagesDestroyInput) =>
  GQL.useImagesDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.imagesDestroy) return;

      for (const id of input.ids) {
        const obj = { __typename: "Image", id };
        deleteObject(cache, obj, GQL.FindImageDocument);
      }

      evictTypeFields(cache, imageMutationImpactedTypeFields);
      evictQueries(cache, [
        ...imageMutationImpactedQueries,
        GQL.StatsDocument, // images size, images count
      ]);
    },
  });

function updateImageIncrementO(id: string) {
  return (
    cache: ApolloCache<Record<string, StoreObject>>,
    result: FetchResult<GQL.ImageIncrementOMutation>
  ) => {
    const updatedOCount = result.data?.imageIncrementO;
    if (updatedOCount === undefined) return;

    const image = cache.readFragment<GQL.SlimImageDataFragment>({
      id: cache.identify({ __typename: "Image", id }),
      fragment: GQL.SlimImageDataFragmentDoc,
      fragmentName: "SlimImageData",
    });

    if (image) {
      // if we have the image, update performer o_counters manually
      for (const performer of image.performers) {
        cache.modify({
          id: cache.identify(performer),
          fields: {
            o_counter(value) {
              return value + 1;
            },
          },
        });
      }
    } else {
      // else refresh all performer o_counters
      evictTypeFields(cache, {
        Performer: ["o_counter"],
      });
    }

    updateStats(cache, "total_o_count", 1);
    updateO(cache, "Image", id, updatedOCount);
    evictQueries(cache, [
      GQL.FindImagesDocument, // filter by o_counter
      GQL.FindPerformersDocument, // filter by o_counter
    ]);
  };
}
export const useImageIncrementO = (id: string) =>
  GQL.useImageIncrementOMutation({
    variables: { id },
    update: updateImageIncrementO(id),
  });

export const mutateImageIncrementO = (id: string) =>
  client.mutate<GQL.ImageIncrementOMutation>({
    mutation: GQL.ImageIncrementODocument,
    variables: { id },
    update: updateImageIncrementO(id),
  });

function updateImageDecrementO(id: string) {
  return (
    cache: ApolloCache<Record<string, StoreObject>>,
    result: FetchResult<GQL.ImageDecrementOMutation>
  ) => {
    const updatedOCount = result.data?.imageDecrementO;
    if (updatedOCount === undefined) return;

    const image = cache.readFragment<GQL.SlimImageDataFragment>({
      id: cache.identify({ __typename: "Image", id }),
      fragment: GQL.SlimImageDataFragmentDoc,
      fragmentName: "SlimImageData",
    });

    if (image) {
      // if we have the image, update performer o_counters manually
      for (const performer of image.performers) {
        cache.modify({
          id: cache.identify(performer),
          fields: {
            o_counter(value) {
              return value - 1;
            },
          },
        });
      }
    } else {
      // else refresh all performer o_counters
      evictTypeFields(cache, {
        Performer: ["o_counter"],
      });
    }

    updateStats(cache, "total_o_count", -1);
    updateO(cache, "Image", id, updatedOCount);
    evictQueries(cache, [
      GQL.FindImagesDocument, // filter by o_counter
      GQL.FindPerformersDocument, // filter by o_counter
    ]);
  };
}

export const useImageDecrementO = (id: string) =>
  GQL.useImageDecrementOMutation({
    variables: { id },
    update: updateImageDecrementO(id),
  });

export const mutateImageDecrementO = (id: string) =>
  client.mutate<GQL.ImageDecrementOMutation>({
    mutation: GQL.ImageDecrementODocument,
    variables: { id },
    update: updateImageDecrementO(id),
  });

function updateImageResetO(id: string) {
  return (
    cache: ApolloCache<Record<string, StoreObject>>,
    result: FetchResult<GQL.ImageResetOMutation>
  ) => {
    const updatedOCount = result.data?.imageResetO;
    if (updatedOCount === undefined) return;

    const image = cache.readFragment<GQL.SlimImageDataFragment>({
      id: cache.identify({ __typename: "Image", id }),
      fragment: GQL.SlimImageDataFragmentDoc,
      fragmentName: "SlimImageData",
    });

    if (image) {
      // if we have the image, update performer o_counters manually
      const old_count = image.o_counter ?? 0;
      for (const performer of image.performers) {
        cache.modify({
          id: cache.identify(performer),
          fields: {
            o_counter(value) {
              return value - old_count;
            },
          },
        });
      }
      updateStats(cache, "total_o_count", -old_count);
    } else {
      // else refresh all performer o_counters
      evictTypeFields(cache, {
        Performer: ["o_counter"],
      });
      // also refresh stats total_o_count
      cache.modify({
        fields: {
          stats: (value) => ({
            ...value,
            total_o_count: undefined,
          }),
        },
      });
    }

    updateO(cache, "Image", id, updatedOCount);
    evictQueries(cache, [
      GQL.FindImagesDocument, // filter by o_counter
      GQL.FindPerformersDocument, // filter by o_counter
    ]);
  };
}

export const useImageResetO = (id: string) =>
  GQL.useImageResetOMutation({
    variables: { id },
    update: updateImageResetO(id),
  });

export const mutateImageResetO = (id: string) =>
  client.mutate<GQL.ImageResetOMutation>({
    mutation: GQL.ImageResetODocument,
    variables: { id },
    update: updateImageResetO(id),
  });

export const mutateImageSetPrimaryFile = (id: string, fileID: string) =>
  client.mutate<GQL.ImageUpdateMutation>({
    mutation: GQL.ImageUpdateDocument,
    variables: {
      input: {
        id,
        primary_file_id: fileID,
      },
    },
    update(cache, result) {
      if (!result.data?.imageUpdate) return;

      evictQueries(cache, [
        GQL.FindImagesDocument, // sort by primary basename when missing title
      ]);
    },
  });

const groupMutationImpactedTypeFields = {
  Performer: ["group_count"],
  Studio: ["group_count"],
};

const groupMutationImpactedQueries = [
  GQL.FindGroupsDocument, // various filters
];

export const useGroupCreate = () =>
  GQL.useGroupCreateMutation({
    update(cache, result) {
      const group = result.data?.groupCreate;
      if (!group) return;

      // update stats
      updateStats(cache, "group_count", 1);

      evictTypeFields(cache, groupMutationImpactedTypeFields);
      evictQueries(cache, groupMutationImpactedQueries);
    },
  });

export const useGroupUpdate = () =>
  GQL.useGroupUpdateMutation({
    update(cache, result) {
      if (!result.data?.groupUpdate) return;

      evictTypeFields(cache, groupMutationImpactedTypeFields);
      evictQueries(cache, groupMutationImpactedQueries);
    },
  });

export const useBulkGroupUpdate = (input: GQL.BulkGroupUpdateInput) =>
  GQL.useBulkGroupUpdateMutation({
    variables: { input },
    update(cache, result) {
      if (!result.data?.bulkGroupUpdate) return;

      evictTypeFields(cache, groupMutationImpactedTypeFields);
      evictQueries(cache, groupMutationImpactedQueries);
    },
  });

export const useGroupDestroy = (input: GQL.GroupDestroyInput) =>
  GQL.useGroupDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.groupDestroy) return;

      const obj = { __typename: "Group", id: input.id };
      deleteObject(cache, obj, GQL.FindGroupDocument);

      // update stats
      updateStats(cache, "group_count", -1);

      evictTypeFields(cache, {
        Scene: ["groups"],
        Performer: ["group_count"],
        Studio: ["group_count"],
      });
      evictQueries(cache, [
        ...groupMutationImpactedQueries,
        GQL.FindScenesDocument, // filter by group
      ]);
    },
  });

export const useGroupsDestroy = (input: GQL.GroupsDestroyMutationVariables) =>
  GQL.useGroupsDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.groupsDestroy) return;

      const { ids } = input;

      for (const id of ids) {
        const obj = { __typename: "Group", id };
        deleteObject(cache, obj, GQL.FindGroupDocument);
      }

      // update stats
      updateStats(cache, "group_count", -ids.length);

      evictTypeFields(cache, {
        Scene: ["groups"],
        Performer: ["group_count"],
        Studio: ["group_count"],
      });
      evictQueries(cache, [
        ...groupMutationImpactedQueries,
        GQL.FindScenesDocument, // filter by group
      ]);
    },
  });

export function useReorderSubGroupsMutation() {
  return GQL.useReorderSubGroupsMutation({
    update(cache) {
      evictQueries(cache, [
        GQL.FindGroupsDocument, // various filters
      ]);
    },
  });
}

export const useAddSubGroups = () => {
  const [addSubGroups] = GQL.useAddGroupSubGroupsMutation({
    update(cache, result) {
      if (!result.data?.addGroupSubGroups) return;

      evictTypeFields(cache, groupMutationImpactedTypeFields);
      evictQueries(cache, groupMutationImpactedQueries);
    },
  });

  return (containingGroupId: string, toAdd: GQL.GroupDescriptionInput[]) => {
    return addSubGroups({
      variables: {
        input: {
          containing_group_id: containingGroupId,
          sub_groups: toAdd,
        },
      },
    });
  };
};

export const useRemoveSubGroups = () => {
  const [removeSubGroups] = GQL.useRemoveGroupSubGroupsMutation({
    update(cache, result) {
      if (!result.data?.removeGroupSubGroups) return;

      evictTypeFields(cache, groupMutationImpactedTypeFields);
      evictQueries(cache, groupMutationImpactedQueries);
    },
  });

  return (containingGroupId: string, removeIds: string[]) => {
    return removeSubGroups({
      variables: {
        input: {
          containing_group_id: containingGroupId,
          sub_group_ids: removeIds,
        },
      },
    });
  };
};

const sceneMarkerMutationImpactedTypeFields = {
  Tag: ["scene_marker_count"],
};

const sceneMarkerMutationImpactedQueries = [
  GQL.FindScenesDocument, // has marker filter
  GQL.FindSceneMarkersDocument, // various filters
  GQL.MarkerStringsDocument, // marker list
  GQL.FindSceneMarkerTagsDocument, // marker tag list
  GQL.FindTagsDocument, // filter by marker count
];

export const useSceneMarkerCreate = () =>
  GQL.useSceneMarkerCreateMutation({
    update(cache, result, { variables }) {
      if (!result.data?.sceneMarkerCreate || !variables) return;

      // refetch linked scene's marker list
      cache.evict({
        id: cache.identify({ __typename: "Scene", id: variables.scene_id }),
        fieldName: "scene_markers",
      });

      evictTypeFields(cache, sceneMarkerMutationImpactedTypeFields);
      evictQueries(cache, sceneMarkerMutationImpactedQueries);
    },
  });

export const useSceneMarkerUpdate = () =>
  GQL.useSceneMarkerUpdateMutation({
    update(cache, result, { variables }) {
      if (!result.data?.sceneMarkerUpdate || !variables) return;

      // refetch linked scene's marker list
      cache.evict({
        id: cache.identify({ __typename: "Scene", id: variables.scene_id }),
        fieldName: "scene_markers",
      });

      evictTypeFields(cache, sceneMarkerMutationImpactedTypeFields);
      evictQueries(cache, sceneMarkerMutationImpactedQueries);
    },
  });

export const useBulkSceneMarkerUpdate = () =>
  GQL.useBulkSceneMarkerUpdateMutation({
    update(cache, result) {
      if (!result.data?.bulkSceneMarkerUpdate) return;

      evictTypeFields(cache, sceneMarkerMutationImpactedTypeFields);
      evictQueries(cache, sceneMarkerMutationImpactedQueries);
    },
  });

export const useSceneMarkerDestroy = () =>
  GQL.useSceneMarkerDestroyMutation({
    update(cache, result, { variables }) {
      if (!result.data?.sceneMarkerDestroy || !variables) return;

      const obj = { __typename: "SceneMarker", id: variables.id };
      cache.evict({ id: cache.identify(obj) });

      evictTypeFields(cache, sceneMarkerMutationImpactedTypeFields);
      evictQueries(cache, sceneMarkerMutationImpactedQueries);
    },
  });

export const useSceneMarkersDestroy = (
  input: GQL.SceneMarkersDestroyMutationVariables
) =>
  GQL.useSceneMarkersDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.sceneMarkersDestroy) return;

      for (const id of input.ids) {
        const obj = { __typename: "SceneMarker", id };
        cache.evict({ id: cache.identify(obj) });
      }

      evictTypeFields(cache, sceneMarkerMutationImpactedTypeFields);
      evictQueries(cache, sceneMarkerMutationImpactedQueries);
    },
  });

const galleryMutationImpactedTypeFields = {
  Scene: ["galleries"],
  Performer: ["gallery_count", "performer_count"],
  Studio: ["gallery_count", "performer_count"],
  Tag: ["gallery_count"],
};

const galleryMutationImpactedQueries = [
  GQL.FindScenesDocument, // is missing galleries
  GQL.FindGalleriesDocument, // various filters
  GQL.FindPerformersDocument, // filter by gallery count
  GQL.FindStudiosDocument, // filter by gallery count
  GQL.FindTagsDocument, // filter by gallery count
];

export const useGalleryCreate = () =>
  GQL.useGalleryCreateMutation({
    update(cache, result) {
      if (!result.data?.galleryCreate) return;

      // update stats
      updateStats(cache, "gallery_count", 1);

      evictTypeFields(cache, galleryMutationImpactedTypeFields);
      evictQueries(cache, galleryMutationImpactedQueries);
    },
  });

export const useGalleryUpdate = () =>
  GQL.useGalleryUpdateMutation({
    update(cache, result) {
      if (!result.data?.galleryUpdate) return;

      evictTypeFields(cache, galleryMutationImpactedTypeFields);
      evictQueries(cache, galleryMutationImpactedQueries);
    },
  });

export const useBulkGalleryUpdate = () =>
  GQL.useBulkGalleryUpdateMutation({
    update(cache, result) {
      if (!result.data?.bulkGalleryUpdate) return;

      evictTypeFields(cache, galleryMutationImpactedTypeFields);
      evictQueries(cache, galleryMutationImpactedQueries);
    },
  });

export const useGalleryDestroy = (input: GQL.GalleryDestroyInput) =>
  GQL.useGalleryDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.galleryDestroy) return;

      for (const id of input.ids) {
        const obj = { __typename: "Gallery", id };
        deleteObject(cache, obj, GQL.FindGalleryDocument);
      }

      evictTypeFields(cache, galleryMutationImpactedTypeFields);
      evictQueries(cache, [
        ...galleryMutationImpactedQueries,
        GQL.FindImagesDocument, // filter by gallery
        GQL.StatsDocument, // images size, gallery count, etc
      ]);
    },
  });

export const mutateAddGalleryImages = (input: GQL.GalleryAddInput) =>
  client.mutate<GQL.AddGalleryImagesMutation>({
    mutation: GQL.AddGalleryImagesDocument,
    variables: input,
    update(cache, result) {
      if (!result.data?.addGalleryImages) return;

      // refetch gallery image_count
      cache.evict({
        id: cache.identify({ __typename: "Gallery", id: input.gallery_id }),
        fieldName: "image_count",
      });

      // refetch images galleries field
      for (const id of input.image_ids) {
        cache.evict({
          id: cache.identify({ __typename: "Image", id }),
          fieldName: "galleries",
        });
      }

      evictQueries(cache, [
        GQL.FindGalleriesDocument, // filter by image count
        GQL.FindImagesDocument, // filter by gallery
      ]);
    },
  });

function evictCover(cache: ApolloCache<GQL.Gallery>, gallery_id: string) {
  const fields: Partial<Pick<Modifiers<GQL.Gallery>, "paths" | "cover">> = {};
  fields.paths = (paths) => {
    if (!("cover" in paths)) {
      return paths;
    }
    const coverUrl = new URL(paths.cover);
    coverUrl.search = "?t=" + Math.floor(Date.now() / 1000);
    return { ...paths, cover: coverUrl.toString() };
  };
  fields.cover = (_value, { DELETE }) => DELETE;
  cache.modify({
    id: cache.identify({ __typename: "Gallery", id: gallery_id }),
    fields,
  });
}

export const mutateSetGalleryCover = (input: GQL.GallerySetCoverInput) =>
  client.mutate<GQL.SetGalleryCoverMutation>({
    mutation: GQL.SetGalleryCoverDocument,
    variables: input,
    update(cache, result) {
      if (!result.data?.setGalleryCover) return;
      evictCover(cache, input.gallery_id);
    },
  });

export const mutateResetGalleryCover = (input: GQL.GalleryResetCoverInput) =>
  client.mutate<GQL.ResetGalleryCoverMutation>({
    mutation: GQL.ResetGalleryCoverDocument,
    variables: input,
    update(cache, result) {
      if (!result.data?.resetGalleryCover) return;
      evictCover(cache, input.gallery_id);
    },
  });

export const mutateRemoveGalleryImages = (input: GQL.GalleryRemoveInput) =>
  client.mutate<GQL.RemoveGalleryImagesMutation>({
    mutation: GQL.RemoveGalleryImagesDocument,
    variables: input,
    update(cache, result) {
      if (!result.data?.removeGalleryImages) return;

      // refetch gallery image_count
      cache.evict({
        id: cache.identify({ __typename: "Gallery", id: input.gallery_id }),
        fieldName: "image_count",
      });

      // refetch images galleries field
      for (const id of input.image_ids) {
        cache.evict({
          id: cache.identify({ __typename: "Image", id }),
          fieldName: "galleries",
        });
      }

      evictQueries(cache, [
        GQL.FindGalleriesDocument, // filter by image count
        GQL.FindImagesDocument, // filter by gallery
      ]);
    },
  });

export const mutateGallerySetPrimaryFile = (id: string, fileID: string) =>
  client.mutate<GQL.GalleryUpdateMutation>({
    mutation: GQL.GalleryUpdateDocument,
    variables: {
      input: {
        id,
        primary_file_id: fileID,
      },
    },
    update(cache, result) {
      if (!result.data?.galleryUpdate) return;

      evictQueries(cache, [
        GQL.FindGalleriesDocument, // sort by primary basename when missing title
      ]);
    },
  });

const galleryChapterMutationImpactedTypeFields = {
  Gallery: ["chapters"],
};

const galleryChapterMutationImpactedQueries = [
  GQL.FindGalleriesDocument, // filter by has chapters
];

export const useGalleryChapterCreate = () =>
  GQL.useGalleryChapterCreateMutation({
    update(cache, result) {
      if (!result.data?.galleryChapterCreate) return;

      evictTypeFields(cache, galleryChapterMutationImpactedTypeFields);
      evictQueries(cache, galleryChapterMutationImpactedQueries);
    },
  });

export const useGalleryChapterUpdate = () =>
  GQL.useGalleryChapterUpdateMutation({
    update(cache, result) {
      if (!result.data?.galleryChapterUpdate) return;

      evictTypeFields(cache, galleryChapterMutationImpactedTypeFields);
      evictQueries(cache, galleryChapterMutationImpactedQueries);
    },
  });

export const useGalleryChapterDestroy = () =>
  GQL.useGalleryChapterDestroyMutation({
    update(cache, result, { variables }) {
      if (!result.data?.galleryChapterDestroy || !variables) return;

      const obj = { __typename: "GalleryChapter", id: variables.id };
      cache.evict({ id: cache.identify(obj) });

      evictTypeFields(cache, galleryChapterMutationImpactedTypeFields);
      evictQueries(cache, galleryChapterMutationImpactedQueries);
    },
  });

const performerMutationImpactedTypeFields = {
  Tag: ["performer_count"],
};

export const performerMutationImpactedQueries = [
  GQL.FindScenesDocument, // filter by performer tags
  GQL.FindImagesDocument, // filter by performer tags
  GQL.FindGalleriesDocument, // filter by performer tags
  GQL.FindPerformersDocument, // various filters
  GQL.FindTagsDocument, // filter by performer count
];

export const usePerformerCreate = () =>
  GQL.usePerformerCreateMutation({
    update(cache, result) {
      const performer = result.data?.performerCreate;
      if (!performer) return;

      // update stats
      updateStats(cache, "performer_count", 1);

      evictTypeFields(cache, performerMutationImpactedTypeFields);
      evictQueries(cache, [
        GQL.FindPerformersDocument, // various filters
        GQL.FindTagsDocument, // filter by performer count
      ]);
    },
  });

export const usePerformerUpdate = () =>
  GQL.usePerformerUpdateMutation({
    update(cache, result) {
      if (!result.data?.performerUpdate) return;

      evictTypeFields(cache, performerMutationImpactedTypeFields);
      evictQueries(cache, performerMutationImpactedQueries);
    },
  });

export const useBulkPerformerUpdate = (input: GQL.BulkPerformerUpdateInput) =>
  GQL.useBulkPerformerUpdateMutation({
    variables: { input },
    update(cache, result) {
      if (!result.data?.bulkPerformerUpdate) return;

      evictTypeFields(cache, performerMutationImpactedTypeFields);
      evictQueries(cache, performerMutationImpactedQueries);
    },
  });

export const usePerformerDestroy = () =>
  GQL.usePerformerDestroyMutation({
    update(cache, result, { variables }) {
      if (!result.data?.performerDestroy || !variables) return;

      const obj = { __typename: "Performer", id: variables.id };
      deleteObject(cache, obj, GQL.FindPerformerDocument);

      // update stats
      updateStats(cache, "performer_count", -1);

      evictTypeFields(cache, {
        ...performerMutationImpactedTypeFields,
        Performer: ["performer_count"],
        Studio: ["performer_count"],
      });
      evictQueries(cache, [
        ...performerMutationImpactedQueries,
        GQL.FindGroupsDocument, // filter by performers
        GQL.FindSceneMarkersDocument, // filter by performers
      ]);
    },
  });

export const usePerformersDestroy = (
  input: GQL.PerformersDestroyMutationVariables
) =>
  GQL.usePerformersDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.performersDestroy) return;

      const { ids } = input;

      let count: number;
      if (Array.isArray(ids)) {
        for (const id of ids) {
          const obj = { __typename: "Performer", id };
          deleteObject(cache, obj, GQL.FindPerformerDocument);
        }
        count = ids.length;
      } else {
        const obj = { __typename: "Performer", id: ids };
        deleteObject(cache, obj, GQL.FindPerformerDocument);
        count = 1;
      }

      // update stats
      updateStats(cache, "performer_count", -count);

      evictTypeFields(cache, {
        ...performerMutationImpactedTypeFields,
        Performer: ["performer_count"],
        Studio: ["performer_count"],
      });
      evictQueries(cache, [
        ...performerMutationImpactedQueries,
        GQL.FindGroupsDocument, // filter by performers
        GQL.FindSceneMarkersDocument, // filter by performers
      ]);
    },
  });

export const mutatePerformerMerge = (
  destination: string,
  source: string[],
  values: GQL.PerformerUpdateInput
) =>
  client.mutate<GQL.PerformerMergeMutation>({
    mutation: GQL.PerformerMergeDocument,
    variables: {
      input: {
        source,
        destination,
        values,
      },
    },
    update(cache, result) {
      if (!result.data?.performerMerge) return;

      for (const id of source) {
        const obj = { __typename: "Performer", id };
        deleteObject(cache, obj, GQL.FindPerformerDocument);
      }

      cache.evict({
        id: cache.identify({ __typename: "Performer", id: destination }),
      });

      evictTypeFields(cache, performerMutationImpactedTypeFields);
      evictQueries(cache, [
        ...performerMutationImpactedQueries,
        GQL.FindGroupsDocument, // filter by performers
        GQL.FindSceneMarkersDocument, // filter by performers
        GQL.StatsDocument, // performer count
      ]);
    },
  });

const studioMutationImpactedTypeFields = {
  Studio: ["child_studios"],
};

export const studioMutationImpactedQueries = [
  GQL.FindScenesDocument, // filter by studio
  GQL.FindImagesDocument, // filter by studio
  GQL.FindGroupsDocument, // filter by studio
  GQL.FindGalleriesDocument, // filter by studio
  GQL.FindPerformersDocument, // filter by studio
  GQL.FindStudiosDocument, // various filters
];

export const useStudioCreate = () =>
  GQL.useStudioCreateMutation({
    update(cache, result, { variables }) {
      const studio = result.data?.studioCreate;
      if (!studio || !variables) return;

      // update stats
      updateStats(cache, "studio_count", 1);

      // if new scene has a parent studio,
      // refetch the parent's list of child studios
      const { parent_id } = variables.input;
      if (parent_id !== undefined) {
        cache.evict({
          id: cache.identify({ __typename: "Studio", id: parent_id }),
          fieldName: "child_studios",
        });
      }

      evictQueries(cache, [
        GQL.FindStudiosDocument, // various filters
      ]);
    },
  });

export const useStudioUpdate = () =>
  GQL.useStudioUpdateMutation({
    update(cache, result) {
      const studio = result.data?.studioUpdate;
      if (!studio) return;

      const obj = { __typename: "Studio", id: studio.id };
      evictTypeFields(
        cache,
        studioMutationImpactedTypeFields,
        cache.identify(obj) // don't evict this studio
      );

      evictQueries(cache, studioMutationImpactedQueries);
    },
  });

export const useBulkStudioUpdate = () =>
  GQL.useBulkStudioUpdateMutation({
    update(cache, result) {
      if (!result.data?.bulkStudioUpdate) return;

      evictTypeFields(cache, studioMutationImpactedTypeFields);
      evictQueries(cache, studioMutationImpactedQueries);
    },
  });

export const useStudioDestroy = (input: GQL.StudioDestroyInput) =>
  GQL.useStudioDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.studioDestroy) return;

      const obj = { __typename: "Studio", id: input.id };
      deleteObject(cache, obj, GQL.FindStudioDocument);

      // update stats
      updateStats(cache, "studio_count", -1);

      evictTypeFields(cache, studioMutationImpactedTypeFields);
      evictQueries(cache, studioMutationImpactedQueries);
    },
  });

export const useStudiosDestroy = (input: GQL.StudiosDestroyMutationVariables) =>
  GQL.useStudiosDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.studiosDestroy) return;

      const { ids } = input;

      for (const id of ids) {
        const obj = { __typename: "Studio", id };
        deleteObject(cache, obj, GQL.FindStudioDocument);
      }

      // update stats
      updateStats(cache, "studio_count", -ids.length);

      evictTypeFields(cache, studioMutationImpactedTypeFields);
      evictQueries(cache, studioMutationImpactedQueries);
    },
  });

const tagMutationImpactedTypeFields = {
  Tag: ["parents", "children"],
};

const tagMutationImpactedQueries = [
  GQL.FindGroupsDocument, // filter by tags
  GQL.FindSceneMarkersDocument, // filter by tags
  GQL.FindScenesDocument, // filter by tags
  GQL.FindImagesDocument, // filter by tags
  GQL.FindGalleriesDocument, // filter by tags
  GQL.FindPerformersDocument, // filter by tags
  GQL.FindTagsDocument, // various filters
];

export const useTagCreate = () =>
  GQL.useTagCreateMutation({
    update(cache, result) {
      const tag = result.data?.tagCreate;
      if (!tag) return;

      // update stats
      updateStats(cache, "tag_count", 1);

      const obj = { __typename: "Tag", id: tag.id };
      evictTypeFields(
        cache,
        tagMutationImpactedTypeFields,
        cache.identify(obj) // don't evict this tag
      );

      evictQueries(cache, [
        GQL.FindTagsDocument, // various filters
      ]);
    },
  });

export const useTagUpdate = () =>
  GQL.useTagUpdateMutation({
    update(cache, result) {
      const tag = result.data?.tagUpdate;
      if (!tag) return;

      const obj = { __typename: "Tag", id: tag.id };
      evictTypeFields(
        cache,
        tagMutationImpactedTypeFields,
        cache.identify(obj) // don't evict this tag
      );

      evictQueries(cache, tagMutationImpactedQueries);
    },
  });

export const useBulkTagUpdate = (input: GQL.BulkTagUpdateInput) =>
  GQL.useBulkTagUpdateMutation({
    variables: { input },
    update(cache, result) {
      if (!result.data?.bulkTagUpdate) return;

      evictTypeFields(cache, tagMutationImpactedTypeFields);
      evictQueries(cache, tagMutationImpactedQueries);
    },
  });

export const useTagDestroy = (input: GQL.TagDestroyInput) =>
  GQL.useTagDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.tagDestroy) return;

      const obj = { __typename: "Tag", id: input.id };
      deleteObject(cache, obj, GQL.FindTagDocument);

      // update stats
      updateStats(cache, "tag_count", -1);

      evictTypeFields(cache, tagMutationImpactedTypeFields);
      evictQueries(cache, tagMutationImpactedQueries);
    },
  });

export const useTagsDestroy = (input: GQL.TagsDestroyMutationVariables) =>
  GQL.useTagsDestroyMutation({
    variables: input,
    update(cache, result) {
      if (!result.data?.tagsDestroy) return;

      const { ids } = input;

      for (const id of ids) {
        const obj = { __typename: "Tag", id };
        deleteObject(cache, obj, GQL.FindTagDocument);
      }

      // update stats
      updateStats(cache, "tag_count", -ids.length);

      evictTypeFields(cache, tagMutationImpactedTypeFields);
      evictQueries(cache, tagMutationImpactedQueries);
    },
  });

export const useTagsMerge = () =>
  GQL.useTagsMergeMutation({
    update(cache, result, { variables }) {
      if (!result.data?.tagsMerge || !variables) return;

      const { source, destination } = variables;

      for (const id of source) {
        const obj = { __typename: "Tag", id };
        deleteObject(cache, obj, GQL.FindTagDocument);
      }

      cache.evict({
        id: cache.identify({ __typename: "Tag", id: destination }),
      });

      evictQueries(cache, [
        ...tagMutationImpactedQueries,
        GQL.StatsDocument, // tag count
      ]);
    },
  });

export const useSaveFilter = () => {
  const [saveFilterMutation] = GQL.useSaveFilterMutation({
    update(cache, result) {
      if (!result.data?.saveFilter) return;

      evictQueries(cache, [GQL.FindSavedFiltersDocument]);
    },
  });

  function saveFilter(filter: ListFilterModel, name: string, id?: string) {
    const filterCopy = filter.clone();

    return saveFilterMutation({
      variables: {
        input: {
          id,
          mode: filter.mode,
          name,
          find_filter: filterCopy.makeFindFilter(),
          object_filter: filterCopy.makeSavedFilter(),
          ui_options: filterCopy.makeSavedUIOptions(),
        },
      },
    });
  }

  return saveFilter;
};

export const useSavedFilterDestroy = () =>
  GQL.useDestroySavedFilterMutation({
    update(cache, result, { variables }) {
      if (!result.data?.destroySavedFilter || !variables) return;

      const obj = { __typename: "SavedFilter", id: variables.input.id };
      deleteObject(cache, obj, GQL.FindSavedFilterDocument);
    },
  });

export const mutateDeleteFiles = (ids: string[]) =>
  client.mutate<GQL.DeleteFilesMutation>({
    mutation: GQL.DeleteFilesDocument,
    variables: { ids },
    update(cache, result) {
      if (!result.data?.deleteFiles) return;

      // we don't know which type the files are,
      // so evict all of them
      for (const id of ids) {
        cache.evict({
          id: cache.identify({ __typename: "VideoFile", id }),
        });
        cache.evict({
          id: cache.identify({ __typename: "ImageFile", id }),
        });
        cache.evict({
          id: cache.identify({ __typename: "GalleryFile", id }),
        });
      }

      evictQueries(cache, [
        GQL.FindScenesDocument, // filter by file count
        GQL.FindImagesDocument, // filter by file count
        GQL.FindGalleriesDocument, // filter by file count
        GQL.StatsDocument, // scenes size, images size
      ]);
    },
  });

/// Scrapers

export const useListSceneScrapers = () => GQL.useListSceneScrapersQuery();

export const queryScrapeScene = (
  source: GQL.ScraperSourceInput,
  sceneId: string
) =>
  client.query<GQL.ScrapeSingleSceneQuery>({
    query: GQL.ScrapeSingleSceneDocument,
    variables: {
      source,
      input: {
        scene_id: sceneId,
      },
    },
    fetchPolicy: "network-only",
  });

export const queryScrapeSceneQuery = (
  source: GQL.ScraperSourceInput,
  q: string
) =>
  client.query<GQL.ScrapeSingleSceneQuery>({
    query: GQL.ScrapeSingleSceneDocument,
    variables: {
      source,
      input: {
        query: q,
      },
    },
    fetchPolicy: "network-only",
  });

export const queryScrapeSceneURL = (url: string) =>
  client.query<GQL.ScrapeSceneUrlQuery>({
    query: GQL.ScrapeSceneUrlDocument,
    variables: { url },
    fetchPolicy: "network-only",
  });

export const queryScrapeSceneQueryFragment = (
  source: GQL.ScraperSourceInput,
  input: GQL.ScrapedSceneInput
) =>
  client.query<GQL.ScrapeSingleSceneQuery>({
    query: GQL.ScrapeSingleSceneDocument,
    variables: {
      source,
      input: {
        scene_input: input,
      },
    },
    fetchPolicy: "network-only",
  });

export const stashBoxSceneBatchQuery = (
  sceneIds: string[],
  stashBoxEndpoint: string
) =>
  client.query<GQL.ScrapeMultiScenesQuery, GQL.ScrapeMultiScenesQueryVariables>(
    {
      query: GQL.ScrapeMultiScenesDocument,
      variables: {
        source: {
          stash_box_endpoint: stashBoxEndpoint,
        },
        input: {
          scene_ids: sceneIds,
        },
      },
    }
  );

export const useListPerformerScrapers = () =>
  GQL.useListPerformerScrapersQuery();

export const useScrapePerformerList = (scraperId: string, q: string) =>
  GQL.useScrapeSinglePerformerQuery({
    variables: {
      source: {
        scraper_id: scraperId,
      },
      input: {
        query: q,
      },
    },
    skip: q === "",
  });

export const queryScrapePerformer = (
  scraperId: string,
  scrapedPerformer: GQL.ScrapedPerformerInput
) =>
  client.query<GQL.ScrapeSinglePerformerQuery>({
    query: GQL.ScrapeSinglePerformerDocument,
    variables: {
      source: {
        scraper_id: scraperId,
      },
      input: {
        performer_input: scrapedPerformer,
      },
    },
    fetchPolicy: "network-only",
  });

export const queryScrapePerformerURL = (url: string) =>
  client.query<GQL.ScrapePerformerUrlQuery>({
    query: GQL.ScrapePerformerUrlDocument,
    variables: { url },
    fetchPolicy: "network-only",
  });

export const stashBoxPerformerQuery = (
  searchVal: string,
  stashBoxEndpoint: string
) =>
  client.query<
    GQL.ScrapeSinglePerformerQuery,
    GQL.ScrapeSinglePerformerQueryVariables
  >({
    query: GQL.ScrapeSinglePerformerDocument,
    variables: {
      source: {
        stash_box_endpoint: stashBoxEndpoint,
      },
      input: {
        query: searchVal,
      },
    },
    fetchPolicy: "network-only",
  });

export const stashBoxStudioQuery = (
  query: string | null,
  stashBoxEndpoint: string
) =>
  client.query<
    GQL.ScrapeSingleStudioQuery,
    GQL.ScrapeSingleStudioQueryVariables
  >({
    query: GQL.ScrapeSingleStudioDocument,
    variables: {
      source: {
        stash_box_endpoint: stashBoxEndpoint,
      },
      input: {
        query: query,
      },
    },
    fetchPolicy: "network-only",
  });

export const stashBoxSceneQuery = (query: string, stashBoxEndpoint: string) =>
  client.query<GQL.ScrapeSingleSceneQuery, GQL.ScrapeSingleSceneQueryVariables>(
    {
      query: GQL.ScrapeSingleSceneDocument,
      variables: {
        source: {
          stash_box_endpoint: stashBoxEndpoint,
        },
        input: {
          query: query,
        },
      },
      fetchPolicy: "network-only",
    }
  );

export const stashBoxTagQuery = (
  query: string | null,
  stashBoxEndpoint: string
) =>
  client.query<GQL.ScrapeSingleTagQuery, GQL.ScrapeSingleTagQueryVariables>({
    query: GQL.ScrapeSingleTagDocument,
    variables: {
      source: {
        stash_box_endpoint: stashBoxEndpoint,
      },
      input: {
        query: query,
      },
    },
    fetchPolicy: "network-only",
  });

export const mutateStashBoxBatchPerformerTag = (
  input: GQL.StashBoxBatchTagInput
) =>
  client.mutate<GQL.StashBoxBatchPerformerTagMutation>({
    mutation: GQL.StashBoxBatchPerformerTagDocument,
    variables: { input },
  });

export const mutateStashBoxBatchStudioTag = (
  input: GQL.StashBoxBatchTagInput
) =>
  client.mutate<GQL.StashBoxBatchStudioTagMutation>({
    mutation: GQL.StashBoxBatchStudioTagDocument,
    variables: { input },
  });

export const useListGroupScrapers = () => GQL.useListGroupScrapersQuery();

export const queryScrapeGroupURL = (url: string) =>
  client.query<GQL.ScrapeGroupUrlQuery>({
    query: GQL.ScrapeGroupUrlDocument,
    variables: { url },
    fetchPolicy: "network-only",
  });

export const useListGalleryScrapers = () => GQL.useListGalleryScrapersQuery();

export const useListImageScrapers = () => GQL.useListImageScrapersQuery();

export const queryScrapeGallery = (scraperId: string, galleryId: string) =>
  client.query<GQL.ScrapeSingleGalleryQuery>({
    query: GQL.ScrapeSingleGalleryDocument,
    variables: {
      source: {
        scraper_id: scraperId,
      },
      input: {
        gallery_id: galleryId,
      },
    },
    fetchPolicy: "network-only",
  });

export const queryScrapeGalleryURL = (url: string) =>
  client.query<GQL.ScrapeGalleryUrlQuery>({
    query: GQL.ScrapeGalleryUrlDocument,
    variables: { url },
    fetchPolicy: "network-only",
  });

export const queryScrapeImage = (scraperId: string, imageId: string) =>
  client.query<GQL.ScrapeSingleImageQuery>({
    query: GQL.ScrapeSingleImageDocument,
    variables: {
      source: {
        scraper_id: scraperId,
      },
      input: {
        image_id: imageId,
      },
    },
    fetchPolicy: "network-only",
  });

export const queryScrapeImageURL = (url: string) =>
  client.query<GQL.ScrapeImageUrlQuery>({
    query: GQL.ScrapeImageUrlDocument,
    variables: { url },
    fetchPolicy: "network-only",
  });

export const mutateSubmitStashBoxSceneDraft = (
  input: GQL.StashBoxDraftSubmissionInput
) =>
  client.mutate<GQL.SubmitStashBoxSceneDraftMutation>({
    mutation: GQL.SubmitStashBoxSceneDraftDocument,
    variables: { input },
  });

export const mutateSubmitStashBoxPerformerDraft = (
  input: GQL.StashBoxDraftSubmissionInput
) =>
  client.mutate<GQL.SubmitStashBoxPerformerDraftMutation>({
    mutation: GQL.SubmitStashBoxPerformerDraftDocument,
    variables: { input },
  });

/// Configuration

export const useConfiguration = () => GQL.useConfigurationQuery();

export const usePlugins = () => GQL.usePluginsQuery();

export const usePluginTasks = () => GQL.usePluginTasksQuery();

export const useStats = () => GQL.useStatsQuery();

export const useVersion = () => GQL.useVersionQuery();

export const useLatestVersion = () =>
  GQL.useLatestVersionQuery({
    notifyOnNetworkStatusChange: true,
    errorPolicy: "ignore",
  });

export const useDLNAStatus = () =>
  GQL.useDlnaStatusQuery({
    fetchPolicy: "no-cache",
  });

export const useJobQueue = () =>
  GQL.useJobQueueQuery({
    fetchPolicy: "no-cache",
  });

export const useLogs = () =>
  GQL.useLogsQuery({
    fetchPolicy: "no-cache",
  });

export const queryLogs = () =>
  client.query<GQL.LogsQuery>({
    query: GQL.LogsDocument,
    fetchPolicy: "no-cache",
  });

export const useSystemStatus = () => GQL.useSystemStatusQuery();
export const refetchSystemStatus = () => {
  client.refetchQueries({
    include: [GQL.SystemStatusDocument],
  });
};

export const useJobsSubscribe = () => GQL.useJobsSubscribeSubscription();

export const useLoggingSubscribe = () => GQL.useLoggingSubscribeSubscription();

// all scraper-related queries
export const scraperMutationImpactedQueries = [
  GQL.ListGroupScrapersDocument,
  GQL.ListPerformerScrapersDocument,
  GQL.ListSceneScrapersDocument,
  GQL.ListImageScrapersDocument,
  GQL.InstalledScraperPackagesDocument,
  GQL.InstalledScraperPackagesStatusDocument,
];

export const mutateReloadScrapers = () =>
  client.mutate<GQL.ReloadScrapersMutation>({
    mutation: GQL.ReloadScrapersDocument,
    update(cache, result) {
      if (!result.data?.reloadScrapers) return;

      evictQueries(cache, scraperMutationImpactedQueries);
    },
  });

// all plugin-related queries
export const pluginMutationImpactedQueries = [
  GQL.PluginsDocument,
  GQL.PluginTasksDocument,
  GQL.InstalledPluginPackagesDocument,
  GQL.InstalledPluginPackagesStatusDocument,
];

export const mutateReloadPlugins = () =>
  client.mutate<GQL.ReloadPluginsMutation>({
    mutation: GQL.ReloadPluginsDocument,
    update(cache, result) {
      if (!result.data?.reloadPlugins) return;

      evictQueries(cache, pluginMutationImpactedQueries);
    },
  });

type BoolMap = { [key: string]: boolean };

export const mutateSetPluginsEnabled = (enabledMap: BoolMap) =>
  client.mutate<GQL.SetPluginsEnabledMutation>({
    mutation: GQL.SetPluginsEnabledDocument,
    variables: { enabledMap },
    update(cache, result) {
      if (!result.data?.setPluginsEnabled) return;

      for (const id in enabledMap) {
        cache.modify({
          id: cache.identify({ __typename: "Plugin", id }),
          fields: {
            enabled() {
              return enabledMap[id];
            },
          },
        });
      }
    },
  });

function updateConfiguration(cache: ApolloCache<unknown>, result: FetchResult) {
  if (!result.data) return;

  evictQueries(cache, [GQL.ConfigurationDocument]);
}

export const useConfigureGeneral = () =>
  GQL.useConfigureGeneralMutation({
    update(cache, result) {
      if (!result.data?.configureGeneral) return;

      evictQueries(cache, [
        GQL.ConfigurationDocument,
        ...scraperMutationImpactedQueries,
        ...pluginMutationImpactedQueries,
      ]);
    },
  });

export const useConfigureInterface = () =>
  GQL.useConfigureInterfaceMutation({
    update: updateConfiguration,
  });

export const useGenerateAPIKey = () =>
  GQL.useGenerateApiKeyMutation({
    update: updateConfiguration,
  });

export const useConfigureDefaults = () =>
  GQL.useConfigureDefaultsMutation({
    update: updateConfiguration,
  });

function updateUIConfig(
  cache: ApolloCache<Record<string, StoreObject>>,
  result: GQL.ConfigureUiMutation["configureUI"] | undefined
) {
  if (!result) return;

  const existing = cache.readQuery<GQL.ConfigurationQuery>({
    query: GQL.ConfigurationDocument,
  });

  cache.writeQuery({
    query: GQL.ConfigurationDocument,
    data: {
      configuration: {
        ...existing?.configuration,
        ui: result,
      },
    },
  });
}

export const useConfigureUI = () =>
  GQL.useConfigureUiMutation({
    update: (cache, result) => updateUIConfig(cache, result.data?.configureUI),
  });

export const useConfigureUISetting = () =>
  GQL.useConfigureUiSettingMutation({
    update: (cache, result) =>
      updateUIConfig(cache, result.data?.configureUISetting),
  });

export const useConfigureScraping = () =>
  GQL.useConfigureScrapingMutation({
    update: updateConfiguration,
  });

export const useConfigureDLNA = () =>
  GQL.useConfigureDlnaMutation({
    update: updateConfiguration,
  });

export const useConfigurePlugin = () =>
  GQL.useConfigurePluginMutation({
    update: updateConfiguration,
  });

export const useEnableDLNA = () => GQL.useEnableDlnaMutation();

export const useDisableDLNA = () => GQL.useDisableDlnaMutation();

export const useAddTempDLNAIP = () => GQL.useAddTempDlnaipMutation();

export const useRemoveTempDLNAIP = () => GQL.useRemoveTempDlnaipMutation();

export const mutateStopJob = (jobID: string) =>
  client.mutate<GQL.StopJobMutation>({
    mutation: GQL.StopJobDocument,
    variables: { job_id: jobID },
  });

const setupMutationImpactedQueries = [
  GQL.ConfigurationDocument,
  GQL.SystemStatusDocument,
];

export const mutateSetup = (input: GQL.SetupInput) =>
  client.mutate<GQL.SetupMutation>({
    mutation: GQL.SetupDocument,
    variables: { input },
    update(cache, result) {
      if (!result.data?.setup) return;

      evictQueries(cache, setupMutationImpactedQueries);
    },
  });

export const mutateMigrate = (input: GQL.MigrateInput) =>
  client.mutate<GQL.MigrateMutation>({
    mutation: GQL.MigrateDocument,
    variables: { input },
  });

// migrate now runs asynchronously, so we need to evict queries
// once it successfully completes
export function postMigrate() {
  evictQueries(clientCache, setupMutationImpactedQueries);
}

/// Packages

// Acts like GQL.useInstalledScraperPackagesStatusQuery if loadUpgrades is true,
// and GQL.useInstalledScraperPackagesQuery if it is false
export const useInstalledScraperPackages = <T extends boolean>(
  loadUpgrades: T
) => {
  const query = loadUpgrades
    ? GQL.InstalledScraperPackagesStatusDocument
    : GQL.InstalledScraperPackagesDocument;

  type TData = T extends true
    ? GQL.InstalledScraperPackagesStatusQuery
    : GQL.InstalledScraperPackagesQuery;
  type TVariables = T extends true
    ? GQL.InstalledScraperPackagesStatusQueryVariables
    : GQL.InstalledScraperPackagesQueryVariables;

  return useQuery<TData, TVariables>(query);
};

export const queryAvailableScraperPackages = (source: string) =>
  client.query<GQL.AvailableScraperPackagesQuery>({
    query: GQL.AvailableScraperPackagesDocument,
    variables: {
      source,
    },
    fetchPolicy: "network-only",
  });

export const mutateInstallScraperPackages = (
  packages: GQL.PackageSpecInput[]
) =>
  client.mutate<GQL.InstallScraperPackagesMutation>({
    mutation: GQL.InstallScraperPackagesDocument,
    variables: {
      packages,
    },
  });

export const mutateUpdateScraperPackages = (packages: GQL.PackageSpecInput[]) =>
  client.mutate<GQL.UpdateScraperPackagesMutation>({
    mutation: GQL.UpdateScraperPackagesDocument,
    variables: {
      packages,
    },
  });

export const mutateUninstallScraperPackages = (
  packages: GQL.PackageSpecInput[]
) =>
  client.mutate<GQL.UninstallScraperPackagesMutation>({
    mutation: GQL.UninstallScraperPackagesDocument,
    variables: {
      packages,
    },
  });

// Acts like GQL.useInstalledPluginPackagesStatusQuery if loadUpgrades is true,
// and GQL.useInstalledPluginPackagesQuery if it is false
export const useInstalledPluginPackages = <T extends boolean>(
  loadUpgrades: T
) => {
  const query = loadUpgrades
    ? GQL.InstalledPluginPackagesStatusDocument
    : GQL.InstalledPluginPackagesDocument;

  type TData = T extends true
    ? GQL.InstalledPluginPackagesStatusQuery
    : GQL.InstalledPluginPackagesQuery;
  type TVariables = T extends true
    ? GQL.InstalledPluginPackagesStatusQueryVariables
    : GQL.InstalledPluginPackagesQueryVariables;

  return useQuery<TData, TVariables>(query);
};

export const queryAvailablePluginPackages = (source: string) =>
  client.query<GQL.AvailablePluginPackagesQuery>({
    query: GQL.AvailablePluginPackagesDocument,
    variables: {
      source,
    },
    fetchPolicy: "network-only",
  });

export const mutateInstallPluginPackages = (packages: GQL.PackageSpecInput[]) =>
  client.mutate<GQL.InstallPluginPackagesMutation>({
    mutation: GQL.InstallPluginPackagesDocument,
    variables: {
      packages,
    },
  });

export const mutateUpdatePluginPackages = (packages: GQL.PackageSpecInput[]) =>
  client.mutate<GQL.UpdatePluginPackagesMutation>({
    mutation: GQL.UpdatePluginPackagesDocument,
    variables: {
      packages,
    },
  });

export const mutateUninstallPluginPackages = (
  packages: GQL.PackageSpecInput[]
) =>
  client.mutate<GQL.UninstallPluginPackagesMutation>({
    mutation: GQL.UninstallPluginPackagesDocument,
    variables: {
      packages,
    },
  });

/// Tasks

export const mutateMetadataScan = (input: GQL.ScanMetadataInput) =>
  client.mutate<GQL.MetadataScanMutation>({
    mutation: GQL.MetadataScanDocument,
    variables: { input },
  });

export const mutateMetadataIdentify = (input: GQL.IdentifyMetadataInput) =>
  client.mutate<GQL.MetadataIdentifyMutation>({
    mutation: GQL.MetadataIdentifyDocument,
    variables: { input },
  });

export const mutateMetadataAutoTag = (input: GQL.AutoTagMetadataInput) =>
  client.mutate<GQL.MetadataAutoTagMutation>({
    mutation: GQL.MetadataAutoTagDocument,
    variables: { input },
  });

export const mutateMetadataGenerate = (input: GQL.GenerateMetadataInput) =>
  client.mutate<GQL.MetadataGenerateMutation>({
    mutation: GQL.MetadataGenerateDocument,
    variables: { input },
  });

export const mutateMetadataClean = (input: GQL.CleanMetadataInput) =>
  client.mutate<GQL.MetadataCleanMutation>({
    mutation: GQL.MetadataCleanDocument,
    variables: { input },
  });

export const mutateCleanGenerated = (input: GQL.CleanGeneratedInput) =>
  client.mutate<GQL.MetadataCleanGeneratedMutation>({
    mutation: GQL.MetadataCleanGeneratedDocument,
    variables: { input },
  });

export const mutateRunPluginTask = (
  pluginId: string,
  taskName: string,
  args?: GQL.Scalars["Map"]["input"]
) =>
  client.mutate<GQL.RunPluginTaskMutation>({
    mutation: GQL.RunPluginTaskDocument,
    variables: { plugin_id: pluginId, task_name: taskName, args },
  });

export const mutateMetadataExport = () =>
  client.mutate<GQL.MetadataExportMutation>({
    mutation: GQL.MetadataExportDocument,
  });

export const mutateExportObjects = (input: GQL.ExportObjectsInput) =>
  client.mutate<GQL.ExportObjectsMutation>({
    mutation: GQL.ExportObjectsDocument,
    variables: { input },
  });

export const mutateMetadataImport = () =>
  client.mutate<GQL.MetadataImportMutation>({
    mutation: GQL.MetadataImportDocument,
  });

export const mutateImportObjects = (input: GQL.ImportObjectsInput) =>
  client.mutate<GQL.ImportObjectsMutation>({
    mutation: GQL.ImportObjectsDocument,
    variables: { input },
  });

export const mutateBackupDatabase = (input: GQL.BackupDatabaseInput) =>
  client.mutate<GQL.BackupDatabaseMutation>({
    mutation: GQL.BackupDatabaseDocument,
    variables: { input },
  });

export const mutateAnonymiseDatabase = (input: GQL.AnonymiseDatabaseInput) =>
  client.mutate<GQL.AnonymiseDatabaseMutation>({
    mutation: GQL.AnonymiseDatabaseDocument,
    variables: { input },
  });

export const mutateOptimiseDatabase = () =>
  client.mutate<GQL.OptimiseDatabaseMutation>({
    mutation: GQL.OptimiseDatabaseDocument,
  });

export const mutateMigrateHashNaming = () =>
  client.mutate<GQL.MigrateHashNamingMutation>({
    mutation: GQL.MigrateHashNamingDocument,
  });

export const mutateMigrateSceneScreenshots = (
  input: GQL.MigrateSceneScreenshotsInput
) =>
  client.mutate<GQL.MigrateSceneScreenshotsMutation>({
    mutation: GQL.MigrateSceneScreenshotsDocument,
    variables: { input },
  });

export const mutateMigrateBlobs = (input: GQL.MigrateBlobsInput) =>
  client.mutate<GQL.MigrateBlobsMutation>({
    mutation: GQL.MigrateBlobsDocument,
    variables: { input },
  });

/// Misc

export const useDirectory = (path?: string) =>
  GQL.useDirectoryQuery({ variables: { path } });

export const queryParseSceneFilenames = (
  filter: GQL.FindFilterType,
  config: GQL.SceneParserInput
) =>
  client.query<GQL.ParseSceneFilenamesQuery>({
    query: GQL.ParseSceneFilenamesDocument,
    variables: { filter, config },
    fetchPolicy: "network-only",
  });
