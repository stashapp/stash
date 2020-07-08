import { useEffect, useState } from "react";
import ApolloClient from "apollo-client";
import { ApolloLink } from "apollo-link";
import { InMemoryCache, NormalizedCacheObject } from "apollo-cache-inmemory";
import { HttpLink } from "apollo-link-http";

export const useStashBoxClient = (uri: string, ApiKey: string) => {
  const [client, setClient] = useState<ApolloClient<NormalizedCacheObject>>();

  useEffect(() => {
    const httpLink = new HttpLink({
      uri,
      fetch,
    });

    const middlewareLink = new ApolloLink((operation, forward) => {
      operation.setContext({ headers: { ApiKey } });
      if (forward) return forward(operation);
      return null;
    });

    const link = middlewareLink.concat(httpLink);
    const cache = new InMemoryCache();

    setClient(
      new ApolloClient({
        name: "stashdb",
        connectToDevTools: true,
        link,
        cache,
      })
    );
  }, [uri, ApiKey]);

  return client;
};
