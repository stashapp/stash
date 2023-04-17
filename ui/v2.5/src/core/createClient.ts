import {
  ApolloClient,
  InMemoryCache,
  split,
  from,
  ServerError,
  TypePolicies,
} from "@apollo/client";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { createClient as createWSClient } from "graphql-ws";
import { onError } from "@apollo/client/link/error";
import { getMainDefinition } from "@apollo/client/utilities";
import { createUploadLink } from "apollo-upload-client";
import * as GQL from "src/core/generated-graphql";

// Policies that tell apollo what the type of the returned object will be.
// In many cases this allows it to return from cache immediately rather than fetching.
const typePolicies: TypePolicies = {
  Query: {
    fields: {
      findImage: {
        read: (_, { args, toReference }) =>
          toReference({
            __typename: "Image",
            id: args?.id,
          }),
      },
      findPerformer: {
        read: (_, { args, toReference }) =>
          toReference({
            __typename: "Performer",
            id: args?.id,
          }),
      },
      findStudio: {
        read: (_, { args, toReference }) =>
          toReference({
            __typename: "Studio",
            id: args?.id,
          }),
      },
      findMovie: {
        read: (_, { args, toReference }) =>
          toReference({
            __typename: "Movie",
            id: args?.id,
          }),
      },
      findGallery: {
        read: (_, { args, toReference }) =>
          toReference({
            __typename: "Gallery",
            id: args?.id,
          }),
      },
      findScene: {
        read: (_, { args, toReference }) =>
          toReference({
            __typename: "Scene",
            id: args?.id,
          }),
      },
      findTag: {
        read: (_, { args, toReference }) =>
          toReference({
            __typename: "Tag",
            id: args?.id,
          }),
      },
    },
  },
  Scene: {
    fields: {
      scene_markers: {
        merge: false,
      },
    },
  },
  Tag: {
    fields: {
      parents: {
        merge: false,
      },
      children: {
        merge: false,
      },
    },
  },
};

export const getBaseURL = () => {
  const baseURL = window.STASH_BASE_URL;
  if (baseURL === "/%BASE_URL%/") return "/";
  return baseURL;
};

export const getPlatformURL = (ws?: boolean) => {
  const platformUrl = new URL(window.location.origin + getBaseURL());

  if (import.meta.env.DEV) {
    platformUrl.port = import.meta.env.VITE_APP_PLATFORM_PORT ?? "9999";

    if (import.meta.env.VITE_APP_HTTPS === "true") {
      platformUrl.protocol = "https:";
    }
  }

  if (ws) {
    if (platformUrl.protocol === "https:") {
      platformUrl.protocol = "wss:";
    } else {
      platformUrl.protocol = "ws:";
    }
  }

  return platformUrl;
};

export const createClient = () => {
  const platformUrl = getPlatformURL();
  const wsPlatformUrl = getPlatformURL(true);

  const url = `${platformUrl}graphql`;
  const wsUrl = `${wsPlatformUrl}graphql`;

  const httpLink = createUploadLink({ uri: url });

  const wsLink = new GraphQLWsLink(
    createWSClient({
      url: wsUrl,
      retryAttempts: Infinity,
      shouldRetry() {
        return true;
      },
    })
  );

  const errorLink = onError(({ networkError }) => {
    // handle unauthorized error by redirecting to the login page
    if (networkError && (networkError as ServerError).statusCode === 401) {
      // redirect to login page
      const newURL = new URL(
        `${getBaseURL()}login`,
        window.location.toString()
      );
      newURL.searchParams.append("returnURL", window.location.href);
      window.location.href = newURL.toString();
    }
  });

  const splitLink = split(
    ({ query }) => {
      const definition = getMainDefinition(query);
      return (
        definition.kind === "OperationDefinition" &&
        definition.operation === "subscription"
      );
    },
    wsLink,
    httpLink
  );

  const link = from([errorLink, splitLink]);

  const cache = new InMemoryCache({ typePolicies });
  const client = new ApolloClient({
    link,
    cache,
  });

  // Watch for scan/clean tasks and reset cache when they complete
  client
    .subscribe<GQL.ScanCompleteSubscribeSubscription>({
      query: GQL.ScanCompleteSubscribeDocument,
    })
    .subscribe({
      next: () => {
        client.resetStore();
      },
    });

  return {
    cache,
    client,
  };
};
