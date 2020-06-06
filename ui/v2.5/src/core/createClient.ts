import ApolloClient from "apollo-client";
import { InMemoryCache } from "apollo-cache-inmemory";
import { WebSocketLink } from "apollo-link-ws";
import { HttpLink } from "apollo-link-http";
import { onError } from "apollo-link-error";
import { ServerError } from "apollo-link-http-common";
import { split, from } from "apollo-link";
import { getMainDefinition } from "apollo-utilities";

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

  const httpLink = new HttpLink({
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
