module.exports = {
    client: {
        service: {
            name: 'stashdb',
            //url: 'http://stashdb.org/graphql',
            url: 'http://localhost:9998/graphql',
        },
        excludes: ['**/queries/**/_*', '**/mutations/**/_*', '**/__tests__/**/*', '**/node_modules']
    }
};
