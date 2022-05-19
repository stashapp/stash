---
linkTitle: Feature Comparison
title: Comparing Features of Other Go GraphQL Implementations
description: Comparing Features of Other Go GraphQL Implementations
menu: main
weight: -1
---

| | [gqlgen](https://github.com/99designs/gqlgen) | [gophers](https://github.com/graph-gophers/graphql-go) | [graphql-go](https://github.com/graphql-go/graphql) | [thunder](https://github.com/samsarahq/thunder) |
| --------: | :-------- | :-------- | :-------- | :-------- |
| Kind | schema first | schema first | run time types | struct first |
| Boilerplate | less | more | more | some |
| Docs | [docs](https://gqlgen.com) & [examples](https://github.com/99designs/gqlgen/tree/master/_examples) | [examples](https://github.com/graph-gophers/graphql-go/tree/master/example/starwars) | [examples](https://github.com/graphql-go/graphql/tree/master/examples) | [examples](https://github.com/samsarahq/thunder/tree/master/example)|
| Query | 👍 | 👍 | 👍 | 👍 |
| Mutation | 👍 | 🚧 [pr](https://github.com/graph-gophers/graphql-go/pull/182) | 👍 | 👍 |
| Subscription | 👍 | 🚧 [pr](https://github.com/graph-gophers/graphql-go/pull/182) | 👍 | 👍 |
| Type Safety | 👍 | 👍 | ⛔️ | 👍 | 
| Type Binding | 👍 | 👍 | ⛔️ | 👍 |
| Embedding | 👍 | ⛔️ | 🚧 [pr](https://github.com/graphql-go/graphql/pull/371) | ⛔️ |
| Interfaces | 👍 | 👍 | 👍 | ⛔️ [is](https://github.com/samsarahq/thunder/issues/78) |
| Generated Enums | 👍 | ⛔️ | ⛔️ | ⛔️ |
| Generated Inputs | 👍 | ⛔️ | ⛔️ | ⛔️ |
| Federation | 👍 | ⛔️ | ⛔️ | ⛔️ |
| Opentracing | 👍 | 👍 | ⛔️ | ✂️[pr](https://github.com/samsarahq/thunder/pull/77) |
| Hooks for error logging | 👍 | ⛔️ | ⛔️ | ⛔️ |
| Dataloading | 👍 | 👍 | 👍 | ⚠️ |
| Concurrency | 👍 | 👍 | 👍 | 👍 |
| Custom errors & error.path | 👍 | ⛔️ [is](https://github.com/graphql-go/graphql/issues/259) | ⛔️ | ⛔️ |
| Query complexity | 👍 | ⛔️ [is](https://github.com/graphql-go/graphql/issues/231) | ⛔️ | ⛔️ |
