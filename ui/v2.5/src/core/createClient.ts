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
import createUploadLink from "apollo-upload-client/createUploadLink.mjs";
import * as GQL from "src/core/generated-graphql";
import { FieldReadFunction } from "@apollo/client/cache";

// A read function that returns a cache reference with the given
// typename if no valid reference is available.
// Allows to return a cached object rather than fetching.
const readReference = (typename: string): FieldReadFunction => {
  return (existing, { args, canRead, toReference }) =>
    canRead(existing)
      ? existing
      : toReference({
          __typename: typename,
          id: args?.id,
        });
};

// A read function that returns null if a cached reference is invalid.
// Means that a dangling reference implies the object was deleted.
const readDanglingNull: FieldReadFunction = (existing, { canRead }) => {
  if (existing === undefined) return undefined;
  return canRead(existing) ? existing : null;
};

const typePolicies: TypePolicies = {
  Query: {
    fields: {
      findImage: {
        read: readReference("Image"),
      },
      findPerformer: {
        read: readReference("Performer"),
      },
      findStudio: {
        read: readReference("Studio"),
      },
      findGroup: {
        read: readReference("Group"),
      },
      findGallery: {
        read: readReference("Gallery"),
      },
      findScene: {
        read: readReference("Scene"),
      },
      findTag: {
        read: readReference("Tag"),
      },
      findSavedFilter: {
        read: readReference("SavedFilter"),
      },
    },
  },
  Scene: {
    fields: {
      studio: {
        read: readDanglingNull,
      },
    },
  },
  Image: {
    fields: {
      studio: {
        read: readDanglingNull,
      },
      paths: {
        merge: false,
      },
    },
  },
  Group: {
    fields: {
      studio: {
        read: readDanglingNull,
      },
    },
  },
  Gallery: {
    fields: {
      studio: {
        read: readDanglingNull,
      },
    },
  },
  Studio: {
    fields: {
      parent_studio: {
        read: readDanglingNull,
      },
    },
  },
};

const possibleTypes = {
  BaseFile: ["VideoFile", "ImageFile", "GalleryFile"],
  VisualFile: ["VideoFile", "ImageFile"],
};

export const baseURL =
  document.querySelector("base")?.getAttribute("href") ?? "/";

export const getPlatformURL = (path?: string) => {
  let url = new URL(window.location.origin + baseURL);

  if (import.meta.env.DEV) {
    if (import.meta.env.VITE_APP_PLATFORM_URL) {
      url = new URL(import.meta.env.VITE_APP_PLATFORM_URL);
    } else {
      url.port = import.meta.env.VITE_APP_PLATFORM_PORT ?? "9999";
    }
  }

  if (path) {
    url.pathname += path;
  }

  return url;
};

export const createClient = () => {
  const url = getPlatformURL("graphql");

  const wsUrl = getPlatformURL("graphql");
  if (wsUrl.protocol === "https:") {
    wsUrl.protocol = "wss:";
  } else {
    wsUrl.protocol = "ws:";
  }

  const httpLink = createUploadLink({ uri: url.toString() });

  const wsClient = createWSClient({
    url: wsUrl.toString(),
    retryAttempts: Infinity,
    shouldRetry() {
      return true;
    },
  });

  const wsLink = new GraphQLWsLink(wsClient);

  const errorLink = onError(({ networkError }) => {
    // handle graphql unauthorized error
    if (networkError && (networkError as ServerError).statusCode === 401) {
      if (import.meta.env.DEV) {
        alert(`\
GraphQL server error: 401 Unauthorized
Authentication cannot be used with the dev server, since the session authorization cookie cannot be sent cross-origin.
Please disable it on the server and refresh the page.`);
        return;
      }
      // redirect to login page
      const newURL = new URL(
        getPlatformURL("login"),
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

  const cache = new InMemoryCache({
    typePolicies,
    possibleTypes: possibleTypes,
  });
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
    wsClient,
  };
};
