import {createPersistedQueryLink} from "apollo-link-persisted-queries";
import {InMemoryCache} from "apollo-cache-inmemory";
import {HttpLink} from "apollo-link-http";
import {WebSocketLink} from "apollo-link-ws";
import {SubscriptionClient} from "subscriptions-transport-ws";
import ws from 'ws';
import {ApolloClient} from "apollo-client";
import fetch from "node-fetch";
import gql from 'graphql-tag';

var uri = process.env.SERVER_URL || 'http://localhost:8080/query';

function test(client) {
    describe('Json', () => {
        it('should follow json escaping rules', async () => {
            let res = await client.query({
                query: gql`{ jsonEncoding }`,
            });

            expect(res.data.jsonEncoding).toBe("󾓭");
            expect(res.errors).toBe(undefined);
        });
    });

    describe('Input defaults', () => {
        it('should pass default values to resolver', async () => {
            let res = await client.query({
                query: gql`{ date(filter:{value: "asdf"}) }`,
            });

            expect(res.data.date).toBe(true);
            expect(res.errors).toBe(undefined);
        });
    });

    describe('Complexity', () => {
        it('should fail when complexity is too high', async () => {
            let res = await client.query({
                query: gql`{ complexity(value: 2000) }`,
            });

            expect(res.errors[0].message).toBe("operation has complexity 2000, which exceeds the limit of 1000");
        });

        it('should succeed when complexity is not too high', async () => {
            let res = await client.query({
                query: gql`{ complexity(value: 100) }`,
            });

            expect(res.data.complexity).toBe(true);
            expect(res.errors).toBe(undefined);
        });
    });

    describe('List Coercion', () => {

        it('should succeed when nested single values are passed', async () => {
            const variable = {
                enumVal: "CUSTOM",
                strVal: "test",
                intVal: 1,
            }
            let res = await client.query({
                variables: {in: variable},
                query: gql`query coercion($in: [ListCoercion!]){ coercion(value: $in ) }`,
            });

            expect(res.data.coercion).toBe(true);
        });

        it('should succeed when single value is passed', async () => {
            let res = await client.query({
                query: gql`{ coercion(value: [{
                    enumVal: CUSTOM
                }]) }`,
            });

            expect(res.data.coercion).toBe(true);
        });

        it('should succeed when single scalar value is passed', async () => {
            let res = await client.query({
                query: gql`{ coercion(value: [{
                    scalarVal: {
                        key : someValue
                    }
                }]) }`,
            });

            expect(res.data.coercion).toBe(true);
        });

        it('should succeed when multiple values are passed', async () => {
            let res = await client.query({
                query: gql`{ coercion(value: [{
                    enumVal: [CUSTOM,NORMAL]
                }]) }`,
            });

            expect(res.data.coercion).toBe(true);
        });

    });

    describe('Errors', () => {
        it('should respond with correct paths', async () => {
            let res = await client.query({
                query: gql`{ path { cc:child { error } } }`,
            });

            expect(res.errors[0].path).toEqual(['path', 0, 'cc', 'error']);
            expect(res.errors[1].path).toEqual(['path', 1, 'cc', 'error']);
            expect(res.errors[2].path).toEqual(['path', 2, 'cc', 'error']);
            expect(res.errors[3].path).toEqual(['path', 3, 'cc', 'error']);
        });

        it('should use the error presenter for custom errors', async () => {
            let res = await client.query({
                query: gql`{ error(type: CUSTOM) }`,
            });

            expect(res.errors[0].message).toEqual('User message');
        });

        it('should pass through for other errors', async () => {
            let res = await client.query({
                query: gql`{ error(type: NORMAL) }`,
            });

            expect(res.errors[0].message).toEqual('normal error');
        });
    });
}

describe('HTTP client', () => {
    const client = new ApolloClient({
        link: createPersistedQueryLink().concat(new HttpLink({uri, fetch})),
        cache: new InMemoryCache(),
        defaultOptions: {
            watchQuery: {
                fetchPolicy: 'network-only',
                errorPolicy: 'ignore',
            },
            query: {
                fetchPolicy: 'network-only',
                errorPolicy: 'all',
            },
        },
    });

   test(client);
});

describe('Websocket client', () => {
    const sc = new SubscriptionClient(uri.replace('http://', 'ws://').replace('https://', 'wss://'), {
        reconnect: false,
        timeout: 1000,
        inactivityTimeout: 100,
    }, ws);

    const client = new ApolloClient({
            link: createPersistedQueryLink().concat(new WebSocketLink(sc)),
            cache: new InMemoryCache(),
            defaultOptions: {
                watchQuery: {
                    fetchPolicy: 'network-only',
                    errorPolicy: 'ignore',
                },
                query: {
                    fetchPolicy: 'network-only',
                    errorPolicy: 'all',
                },
            },
        });

    test(client);

    afterAll(() => {
        client.stop();
        sc.close(true);
    });
});
