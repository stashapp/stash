import { useCallback, useEffect, useRef, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { getClient, useSceneUpdate, useBulkSceneUpdate } from "src/core/StashService";

export const WATCH_QUEUE_TAG_NAME = "Watch Queue";

// Percentage of video that must be watched to auto-remove from queue
export const WATCH_QUEUE_REMOVE_PERCENT = 5;

export function useWatchQueue() {
  const [tagId, setTagId] = useState<string | null>(null);
  const resolveAttempted = useRef(false);

  const [updateScene] = useSceneUpdate();
  const [bulkUpdateScenes] = useBulkSceneUpdate();

  useEffect(() => {
    if (resolveAttempted.current) return;
    resolveAttempted.current = true;

    const client = getClient();

    client
      .query<GQL.FindTagsQuery>({
        query: GQL.FindTagsDocument,
        variables: {
          filter: { q: WATCH_QUEUE_TAG_NAME, per_page: 10 },
        },
      })
      .then((result) => {
        const tags = result.data?.findTags?.tags ?? [];
        const existing = tags.find((t) => t.name === WATCH_QUEUE_TAG_NAME);
        if (existing) {
          setTagId(existing.id);
          return;
        }

        return client
          .mutate<GQL.TagCreateMutation>({
            mutation: GQL.TagCreateDocument,
            variables: { input: { name: WATCH_QUEUE_TAG_NAME } },
          })
          .then((createResult) => {
            const newTag = createResult.data?.tagCreate;
            if (newTag) setTagId(newTag.id);
          });
      })
      .catch(console.error);
  }, []);

  const isInQueue = useCallback(
    (scene: { tags: Array<{ id: string }> }) => {
      if (!tagId) return false;
      return scene.tags.some((t) => t.id === tagId);
    },
    [tagId]
  );

  const addToQueue = useCallback(
    (sceneId: string, currentTagIds: string[]) => {
      if (!tagId) return Promise.resolve();
      return updateScene({
        variables: {
          input: {
            id: sceneId,
            tag_ids: [...currentTagIds, tagId],
          },
        },
      });
    },
    [tagId, updateScene]
  );

  const removeFromQueue = useCallback(
    (sceneId: string, currentTagIds: string[]) => {
      if (!tagId) return Promise.resolve();
      return updateScene({
        variables: {
          input: {
            id: sceneId,
            tag_ids: currentTagIds.filter((id) => id !== tagId),
          },
        },
      });
    },
    [tagId, updateScene]
  );

  const bulkAddToQueue = useCallback(
    (sceneIds: string[]) => {
      if (!tagId) return Promise.resolve();
      return bulkUpdateScenes({
        variables: {
          input: {
            ids: sceneIds,
            tag_ids: {
              ids: [tagId],
              mode: GQL.BulkUpdateIdMode.Add,
            },
          },
        },
      });
    },
    [tagId, bulkUpdateScenes]
  );

  const bulkRemoveFromQueue = useCallback(
    (sceneIds: string[]) => {
      if (!tagId) return Promise.resolve();
      return bulkUpdateScenes({
        variables: {
          input: {
            ids: sceneIds,
            tag_ids: {
              ids: [tagId],
              mode: GQL.BulkUpdateIdMode.Remove,
            },
          },
        },
      });
    },
    [tagId, bulkUpdateScenes]
  );

  return {
    tagId,
    isInQueue,
    addToQueue,
    removeFromQueue,
    bulkAddToQueue,
    bulkRemoveFromQueue,
  };
}
