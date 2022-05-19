---
title: 'Using Apollo federation gqlgen'
description: How federate many services into a single graph using Apollo
linkTitle: Apollo Federation
menu: { main: { parent: 'recipes' } }
---

In this quick guide we are going to implement the example [Apollo Federation](https://www.apollographql.com/docs/apollo-server/federation/introduction/)
server in gqlgen. You can find the finished result in the [examples directory](https://github.com/99designs/gqlgen/tree/master/_examples/federation).

## Enable federation

Uncomment federation configuration in your `gqlgen.yml`

```yml
# Uncomment to enable federation
federation:
  filename: graph/generated/federation.go
  package: generated
```

## Create the federated servers

For each server to be federated we will create a new gqlgen project.

```bash
go run github.com/99designs/gqlgen
```

Update the schema to reflect the federated example
```graphql
type Review {
  body: String
  author: User @provides(fields: "username")
  product: Product
}

extend type User @key(fields: "id") {
  id: ID! @external
  reviews: [Review]
}

extend type Product @key(fields: "upc") {
  upc: String! @external
  reviews: [Review]
}
```


and regenerate
```bash
go run github.com/99designs/gqlgen
```

then implement the resolvers
```go
// These two methods are required for gqlgen to resolve the internal id-only wrapper structs.
// This boilerplate might be removed in a future version of gqlgen that can no-op id only nodes.
func (r *entityResolver) FindProductByUpc(ctx context.Context, upc string) (*model.Product, error) {
	return &model.Product{
		Upc: upc,
	}, nil
}

func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	return &model.User{
		ID: id,
	}, nil
}

// Here we implement the stitched part of this service, returning reviews for a product. Of course normally you would
// go back to the database, but we are just making some data up here.
func (r *productResolver) Reviews(ctx context.Context, obj *model.Product) ([]*model.Review, error) {
	switch obj.Upc {
	case "top-1":
		return []*model.Review{{
			Body: "A highly effective form of birth control.",
		}}, nil

	case "top-2":
		return []*model.Review{{
			Body: "Fedoras are one of the most fashionable hats around and can look great with a variety of outfits.",
		}}, nil

	case "top-3":
		return []*model.Review{{
			Body: "This is the last straw. Hat you will wear. 11/10",
		}}, nil

	}
	return nil, nil
}

func (r *userResolver) Reviews(ctx context.Context, obj *model.User) ([]*model.Review, error) {
	if obj.ID == "1234" {
		return []*model.Review{{
			Body: "Has an odd fascination with hats.",
		}}, nil
	}
	return nil, nil
}
```

> Note
>
> Repeat this step for each of the services in the apollo doc (accounts, products, reviews)

## Create the federation gateway

```bash
npm install --save @apollo/gateway apollo-server graphql
```

```typescript
const { ApolloServer } = require('apollo-server');
const { ApolloGateway } = require("@apollo/gateway");

const gateway = new ApolloGateway({
    serviceList: [
        { name: 'accounts', url: 'http://localhost:4001/query' },
        { name: 'products', url: 'http://localhost:4002/query' },
        { name: 'reviews', url: 'http://localhost:4003/query' }
    ],
});

const server = new ApolloServer({
    gateway,

    subscriptions: false,
});

server.listen().then(({ url }) => {
    console.log(`🚀 Server ready at ${url}`);
});
```

## Start all the services

In separate terminals:
```bash
go run accounts/server.go
go run products/server.go
go run reviews/server.go
node gateway/index.js
```

## Query the federated gateway

The examples from the apollo doc should all work, eg

```graphql
query {
  me {
    username
    reviews {
      body
      product {
        name
        upc
      }
    }
  }
}
```

should return

```json
{
  "data": {
    "me": {
      "username": "Me",
      "reviews": [
        {
          "body": "A highly effective form of birth control.",
          "product": {
            "name": "Trilby",
            "upc": "top-1"
          }
        },
        {
          "body": "Fedoras are one of the most fashionable hats around and can look great with a variety of outfits.",
          "product": {
            "name": "Trilby",
            "upc": "top-1"
          }
        }
      ]
    }
  }
}
```
