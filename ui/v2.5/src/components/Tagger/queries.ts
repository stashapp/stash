import * as GQL from "src/core/generated-graphql";
import {
  evictQueries,
  getClient,
  studioMutationImpactedQueries,
} from "src/core/StashService";

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
                stash_id_endpoint: {
                  stash_id: id.stash_id,
                  endpoint: id.endpoint,
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
                  stash_id_endpoint: {
                    stash_id: id.stash_id,
                    endpoint: id.endpoint,
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
