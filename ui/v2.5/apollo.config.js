module.exports = {
    client: {
        service: {
            name: 'stashdb',
            url: 'http://stashdb.org/graphql',
        },
        excludes: ['**/queries/**/_*', '**/mutations/**/_*', '**/__tests__/**/*', '**/node_modules']
    }
};
