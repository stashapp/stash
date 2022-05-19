package singlefile

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

func TestSubscriptions(t *testing.T) {
	tick := make(chan string, 1)

	resolvers := &Stub{}

	resolvers.SubscriptionResolver.InitPayload = func(ctx context.Context) (strings <-chan string, e error) {
		payload := transport.GetInitPayload(ctx)
		channel := make(chan string, len(payload)+1)

		go func() {
			<-ctx.Done()
			close(channel)
		}()

		// Test the helper function separately
		auth := payload.Authorization()
		if auth != "" {
			channel <- "AUTH:" + auth
		} else {
			channel <- "AUTH:NONE"
		}

		// Send them over the channel in alphabetic order
		keys := make([]string, 0, len(payload))
		for key := range payload {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			channel <- fmt.Sprintf("%s = %#+v", key, payload[key])
		}

		return channel, nil
	}

	resolvers.SubscriptionResolver.Updated = func(ctx context.Context) (<-chan string, error) {
		res := make(chan string, 1)

		go func() {
			for {
				select {
				case t := <-tick:
					res <- t
				case <-ctx.Done():
					close(res)
					return
				}
			}
		}()
		return res, nil
	}

	srv := handler.NewDefaultServer(
		NewExecutableSchema(Config{Resolvers: resolvers}),
	)
	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		path, _ := ctx.Value(ckey("path")).([]int)
		return next(context.WithValue(ctx, ckey("path"), append(path, 1)))
	})

	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		path, _ := ctx.Value(ckey("path")).([]int)
		return next(context.WithValue(ctx, ckey("path"), append(path, 2)))
	})

	c := client.New(srv)

	t.Run("wont leak goroutines", func(t *testing.T) {
		runtime.GC() // ensure no go-routines left from preceding tests
		initialGoroutineCount := runtime.NumGoroutine()

		sub := c.Websocket(`subscription { updated }`)

		tick <- "message"

		var msg struct {
			resp struct {
				Updated string
			}
		}

		err := sub.Next(&msg.resp)
		require.NoError(t, err)
		require.Equal(t, "message", msg.resp.Updated)
		sub.Close()

		// need a little bit of time for goroutines to settle
		start := time.Now()
		for time.Since(start).Seconds() < 2 && initialGoroutineCount != runtime.NumGoroutine() {
			time.Sleep(5 * time.Millisecond)
		}

		require.Equal(t, initialGoroutineCount, runtime.NumGoroutine())
	})

	t.Run("will parse init payload", func(t *testing.T) {
		sub := c.WebsocketWithPayload(`subscription { initPayload }`, map[string]interface{}{
			"Authorization": "Bearer of the curse",
			"number":        32,
			"strings":       []string{"hello", "world"},
		})

		var msg struct {
			resp struct {
				InitPayload string
			}
		}

		err := sub.Next(&msg.resp)
		require.NoError(t, err)
		require.Equal(t, "AUTH:Bearer of the curse", msg.resp.InitPayload)
		err = sub.Next(&msg.resp)
		require.NoError(t, err)
		require.Equal(t, "Authorization = \"Bearer of the curse\"", msg.resp.InitPayload)
		err = sub.Next(&msg.resp)
		require.NoError(t, err)
		require.Equal(t, "number = 32", msg.resp.InitPayload)
		err = sub.Next(&msg.resp)
		require.NoError(t, err)
		require.Equal(t, "strings = []interface {}{\"hello\", \"world\"}", msg.resp.InitPayload)
		sub.Close()
	})
}
