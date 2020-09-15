import {
  ApolloClient,
  InMemoryCache,
  split,
  from,
  ServerError,
} from "@apollo/client";
import { WebSocketLink } from "@apollo/client/link/ws";
import { onError } from "@apollo/client/link/error";
import { getMainDefinition } from "@apollo/client/utilities";
import { createUploadLink } from "apollo-upload-client";

export const getPlatformURL = (ws?: boolean) => {
  const platformUrl = new URL(window.location.origin);

  if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
    platformUrl.port = "9999"; // TODO: Hack. Development expects port 9999

    if (process.env.REACT_APP_HTTPS === "true") {
      platformUrl.protocol = "https:";
    }
  }

  if (ws) {
    platformUrl.protocol = "ws:";
  }

  return platformUrl;
};

export const createClient = () => {
  const platformUrl = getPlatformURL();
  const wsPlatformUrl = getPlatformURL(true);

  if (platformUrl.protocol === "https:") {
    wsPlatformUrl.protocol = "wss:";
  }

  const url = `${platformUrl.toString().slice(0, -1)}/graphql`;
  const wsUrl = `${wsPlatformUrl.toString().slice(0, -1)}/graphql`;

  const httpLink = createUploadLink({
    uri: url,
  });

  const wsLink = new WebSocketLink({
    uri: wsUrl,
    options: {
      reconnect: true,
    },
  });

  const errorLink = onError(({ networkError }) => {
    // handle unauthorized error by redirecting to the login page
    if (networkError && (networkError as ServerError).statusCode === 401) {
      // redirect to login page
      window.location.href = "/login";
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
    // @ts-ignore
    httpLink
  );

  const link = from([errorLink, splitLink]);

  const cache = new InMemoryCache();
  const client = new ApolloClient({
    link,
    cache,
  });

  return {
    cache,
    client,
  };
};
