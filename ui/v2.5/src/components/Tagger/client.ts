import ApolloClient from 'apollo-client';
import { ApolloLink } from 'apollo-link';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpLink } from 'apollo-link-http';

const ApiKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhZTA1NmQ0ZC0wYjRmLTQzNmMtYmVhMy0zNjNjMTQ2MmZlNjMiLCJpYXQiOjE1ODYwNDAzOTUsInN1YiI6IkFQSUtleSJ9.5VENvrLtJXTGcdOhA0QC1SyPQ59padh1XiQRDQelzA4';

const createClient = () => {
  const httpLink = new HttpLink({
    uri: 'https://stashdb.org/graphql',
    fetch
  });

	const middlewareLink = new ApolloLink((operation, forward) => {
		operation.setContext({ headers: { ApiKey } });
		if (forward)
				return forward(operation);
		return null;
	});

	const link = middlewareLink.concat(httpLink);
  const cache = new InMemoryCache();

	return new ApolloClient({
		name: 'stashdb',
		connectToDevTools: true,
		link,
		cache
	});
};

export const client = createClient();
