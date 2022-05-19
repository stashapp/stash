---
title: "Automatic persisted queries"
description:
linkTitle: "APQ"
menu: { main: { parent: 'reference', weight: 10 } }
---

When you work with GraphQL by default your queries are transferred with every request. That can waste significant
bandwidth. To avoid that you can use Automatic Persisted Queries (APQ).

With APQ you send only query hash to the server. If hash is not found on a server then client makes a second request
to register query hash with original query on a server.

## Usage

In order to enable Automatic Persisted Queries you need to change your client. For more information see
[Automatic Persisted Queries Link](https://www.apollographql.com/docs/resources/graphql-glossary/#automatic-persisted-queries-apq) documentation.

For the server you need to implement the `graphql.Cache` interface and pass an instance to
the `extension.AutomaticPersistedQuery` type. Make sure the extension is applied to your GraphQL handler.

See example using [go-redis](https://github.com/go-redis/redis) package below:
```go
import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-redis/redis"
)

type Cache struct {
	client redis.UniversalClient
	ttl    time.Duration
}

const apqPrefix = "apq:"

func NewCache(redisAddress string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
	})

	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("could not create cache: %w", err)
	}

	return &Cache{client: client, ttl: ttl}, nil
}

func (c *Cache) Add(ctx context.Context, key string, value interface{}) {
	c.client.Set(apqPrefix+key, value, c.ttl)
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	s, err := c.client.Get(apqPrefix + key).Result()
	if err != nil {
		return struct{}{}, false
	}
	return s, true
}

func main() {
	cache, err := NewCache(cfg.RedisAddress, 24*time.Hour)
	if err != nil {
		log.Fatalf("cannot create APQ redis cache: %v", err)
	}

	c := Config{ Resolvers: &resolvers{} }
	gqlHandler := handler.New(
		generated.NewExecutableSchema(c),
	)
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.Use(extension.AutomaticPersistedQuery{Cache: cache})
	http.Handle("/query", gqlHandler)
}
```
