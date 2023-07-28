import * as GQL from "src/core/generated-graphql";
import sortBy from "lodash-es/sortBy";
import {
  evictQueries,
  getClient,
  studioMutationImpactedQueries,
} from "src/core/StashService";

export const useUpdatePerformerStashID = () => {
  const [updatePerformer] = GQL.usePerformerUpdateMutation({
    onError: (errors) => errors,
  });

  const updatePerformerHandler = (
    performerID: string,
    stashIDs: GQL.StashIdInput[]
  ) =>
    updatePerformer({
      variables: {
        input: {
          id: performerID,
          stash_ids: stashIDs.map((s) => ({
            stash_id: s.stash_id,
            endpoint: s.endpoint,
          })),
        },
      },
      update: (store, updatedPerformer) => {
        if (!updatedPerformer.data?.performerUpdate) return;
        const newStashID = stashIDs[stashIDs.length - 1].stash_id;

        store.writeQuery<
          GQL.FindPerformersQuery,
          GQL.FindPerformersQueryVariables
        >({
          query: GQL.FindPerformersDocument,
          variables: {
            performer_filter: {
              stash_id: {
                value: newStashID,
                modifier: GQL.CriterionModifier.Equals,
              },
            },
          },
          data: {
            findPerformers: {
              count: 1,
              performers: [updatedPerformer.data.performerUpdate],
              __typename: "FindPerformersResultType",
            },
          },
        });
      },
    });

  return updatePerformerHandler;
};

export const useUpdatePerformer = () => {
  const [updatePerformer] = GQL.usePerformerUpdateMutation({
    onError: (errors) => errors,
    errorPolicy: "all",
  });

  const updatePerformerHandler = (input: GQL.PerformerUpdateInput) =>
    updatePerformer({
      variables: {
        input,
      },
      update: (store, updatedPerformer) => {
        if (!updatedPerformer.data?.performerUpdate) return;

        updatedPerformer.data.performerUpdate.stash_ids.forEach((id) => {
          store.writeQuery<
            GQL.FindPerformersQuery,
            GQL.FindPerformersQueryVariables
          >({
            query: GQL.FindPerformersDocument,
            variables: {
              performer_filter: {
                stash_id: {
                  value: id.stash_id,
                  modifier: GQL.CriterionModifier.Equals,
                },
              },
            },
            data: {
              findPerformers: {
                count: 1,
                performers: [updatedPerformer.data!.performerUpdate!],
                __typename: "FindPerformersResultType",
              },
            },
          });
        });
      },
    });

  return updatePerformerHandler;
};

export const useCreatePerformer = () => {
  const [createPerformer] = GQL.usePerformerCreateMutation({
    onError: (errors) => errors,
  });

  const handleCreate = (performer: GQL.PerformerCreateInput, stashID: string) =>
    createPerformer({
      variables: { input: performer },
      update: (store, newPerformer) => {
        if (!newPerformer?.data?.performerCreate) return;

        const currentQuery = store.readQuery<
          GQL.AllPerformersForFilterQuery,
          GQL.AllPerformersForFilterQueryVariables
        >({
          query: GQL.AllPerformersForFilterDocument,
        });
        const allPerformers = sortBy(
          [
            ...(currentQuery?.allPerformers ?? []),
            newPerformer.data.performerCreate,
          ],
          ["name"]
        );
        if (allPerformers.length > 1) {
          store.writeQuery<
            GQL.AllPerformersForFilterQuery,
            GQL.AllPerformersForFilterQueryVariables
          >({
            query: GQL.AllPerformersForFilterDocument,
            data: {
              allPerformers,
            },
          });
        }

        store.writeQuery<
          GQL.FindPerformersQuery,
          GQL.FindPerformersQueryVariables
        >({
          query: GQL.FindPerformersDocument,
          variables: {
            performer_filter: {
              stash_id: {
                value: stashID,
                modifier: GQL.CriterionModifier.Equals,
              },
            },
          },
          data: {
            findPerformers: {
              count: 1,
              performers: [newPerformer.data.performerCreate],
              __typename: "FindPerformersResultType",
            },
          },
        });
      },
    });

  return handleCreate;
};

export const useUpdateStudioStashID = () => {
  const [updateStudio] = GQL.useStudioUpdateMutation({
    onError: (errors) => errors,
  });

  const handleUpdate = (
    studio: GQL.SlimStudioDataFragment,
    stashIDs: GQL.StashIdInput[]
  ) =>
    updateStudio({
      variables: {
        input: {
          id: studio.id,
          stash_ids: stashIDs.map((s) => ({
            stash_id: s.stash_id,
            endpoint: s.endpoint,
          })),
        },
      },
      update: (store, result) => {
        if (!result.data?.studioUpdate) return;
        const newStashID = stashIDs[stashIDs.length - 1].stash_id;

        store.writeQuery<GQL.FindStudiosQuery, GQL.FindStudiosQueryVariables>({
          query: GQL.FindStudiosDocument,
          variables: {
            studio_filter: {
              stash_id: {
                value: newStashID,
                modifier: GQL.CriterionModifier.Equals,
              },
            },
          },
          data: {
            findStudios: {
              count: 1,
              studios: [result.data.studioUpdate],
              __typename: "FindStudiosResultType",
            },
          },
        });
      },
    });

  return handleUpdate;
};

export const useUpdateStudio = () => {
  const [updateStudio] = GQL.useStudioUpdateMutation({
    onError: (errors) => errors,
    errorPolicy: "all",
  });

  const updateStudioHandler = (input: GQL.StudioUpdateInput) =>
    updateStudio({
      variables: {
        input,
      },
      update: (store, updatedStudio) => {
        if (!updatedStudio.data?.studioUpdate) return;

        if (updatedStudio.data?.studioUpdate?.parent_studio) {
          const ac = getClient();
          evictQueries(ac.cache, studioMutationImpactedQueries);
        } else {
          updatedStudio.data.studioUpdate.stash_ids.forEach((id) => {
            store.writeQuery<
              GQL.FindStudiosQuery,
              GQL.FindStudiosQueryVariables
            >({
              query: GQL.FindStudiosDocument,
              variables: {
                studio_filter: {
                  stash_id: {
                    value: id.stash_id,
                    modifier: GQL.CriterionModifier.Equals,
                  },
                },
              },
              data: {
                findStudios: {
                  count: 1,
                  studios: [updatedStudio.data!.studioUpdate!],
                  __typename: "FindStudiosResultType",
                },
              },
            });
          });
        }
      },
    });

  return updateStudioHandler;
};

export const useCreateStudio = () => {
  const [createStudio] = GQL.useStudioCreateMutation({
    onError: (errors) => errors,
  });

  const handleCreate = (studio: GQL.StudioCreateInput, stashID: string) =>
    createStudio({
      variables: { input: studio },
      update: (store, result) => {
        if (!result?.data?.studioCreate) return;

        const currentQuery = store.readQuery<
          GQL.AllStudiosForFilterQuery,
          GQL.AllStudiosForFilterQueryVariables
        >({
          query: GQL.AllStudiosForFilterDocument,
        });
        const allStudios = sortBy(
          [...(currentQuery?.allStudios ?? []), result.data.studioCreate],
          ["name"]
        );
        if (allStudios.length > 1) {
          store.writeQuery<
            GQL.AllStudiosForFilterQuery,
            GQL.AllStudiosForFilterQueryVariables
          >({
            query: GQL.AllStudiosForFilterDocument,
            data: {
              allStudios,
            },
          });
        }

        store.writeQuery<GQL.FindStudiosQuery, GQL.FindStudiosQueryVariables>({
          query: GQL.FindStudiosDocument,
          variables: {
            studio_filter: {
              stash_id: {
                value: stashID,
                modifier: GQL.CriterionModifier.Equals,
              },
            },
          },
          data: {
            findStudios: {
              count: 1,
              studios: [result.data.studioCreate],
              __typename: "FindStudiosResultType",
            },
          },
        });
      },
    });

  return handleCreate;
};

export const useCreateTag = () => {
  const [createTag] = GQL.useTagCreateMutation({
    onError: (errors) => errors,
  });

  const handleCreate = (tag: string) =>
    createTag({
      variables: {
        input: {
          name: tag,
        },
      },
      update: (store, result) => {
        if (!result.data?.tagCreate) return;

        const currentQuery = store.readQuery<
          GQL.AllTagsForFilterQuery,
          GQL.AllTagsForFilterQueryVariables
        >({
          query: GQL.AllTagsForFilterDocument,
        });
        const allTags = sortBy(
          [...(currentQuery?.allTags ?? []), result.data.tagCreate],
          ["name"]
        );

        store.writeQuery<
          GQL.AllTagsForFilterQuery,
          GQL.AllTagsForFilterQueryVariables
        >({
          query: GQL.AllTagsForFilterDocument,
          data: {
            allTags,
          },
        });
      },
    });

  return handleCreate;
};
