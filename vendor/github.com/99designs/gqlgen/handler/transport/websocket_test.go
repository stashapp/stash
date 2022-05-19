package transport_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type ckey string

func TestWebsocket(t *testing.T) {
	handler := testserver.New()
	handler.AddTransport(transport.Websocket{})

	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("client must send valid json", func(t *testing.T) {
		c := wsConnect(srv.URL)
		defer c.Close()

		writeRaw(c, "hello")

		msg := readOp(c)
		assert.Equal(t, "connection_error", msg.Type)
		assert.Equal(t, `{"message":"invalid json"}`, string(msg.Payload))
	})

	t.Run("client can terminate before init", func(t *testing.T) {
		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionTerminateMsg}))

		_, _, err := c.ReadMessage()
		assert.Equal(t, websocket.CloseNormalClosure, err.(*websocket.CloseError).Code)
	})

	t.Run("client must send init first", func(t *testing.T) {
		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: startMsg}))

		msg := readOp(c)
		assert.Equal(t, connectionErrorMsg, msg.Type)
		assert.Equal(t, `{"message":"unexpected message start"}`, string(msg.Payload))
	})

	t.Run("server acks init", func(t *testing.T) {
		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))

		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
	})

	t.Run("client can terminate before run", func(t *testing.T) {
		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionTerminateMsg}))

		_, _, err := c.ReadMessage()
		assert.Equal(t, websocket.CloseNormalClosure, err.(*websocket.CloseError).Code)
	})

	t.Run("client gets parse errors", func(t *testing.T) {
		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)

		require.NoError(t, c.WriteJSON(&operationMessage{
			Type:    startMsg,
			ID:      "test_1",
			Payload: json.RawMessage(`{"query": "!"}`),
		}))

		msg := readOp(c)
		assert.Equal(t, errorMsg, msg.Type)
		assert.Equal(t, `[{"message":"Unexpected !","locations":[{"line":1,"column":1}],"extensions":{"code":"GRAPHQL_PARSE_FAILED"}}]`, string(msg.Payload))
	})

	t.Run("client can receive data", func(t *testing.T) {
		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)

		require.NoError(t, c.WriteJSON(&operationMessage{
			Type:    startMsg,
			ID:      "test_1",
			Payload: json.RawMessage(`{"query": "subscription { name }"}`),
		}))

		handler.SendNextSubscriptionMessage()
		msg := readOp(c)
		require.Equal(t, dataMsg, msg.Type, string(msg.Payload))
		require.Equal(t, "test_1", msg.ID, string(msg.Payload))
		require.Equal(t, `{"data":{"name":"test"}}`, string(msg.Payload))

		handler.SendNextSubscriptionMessage()
		msg = readOp(c)
		require.Equal(t, dataMsg, msg.Type, string(msg.Payload))
		require.Equal(t, "test_1", msg.ID, string(msg.Payload))
		require.Equal(t, `{"data":{"name":"test"}}`, string(msg.Payload))

		require.NoError(t, c.WriteJSON(&operationMessage{Type: stopMsg, ID: "test_1"}))

		msg = readOp(c)
		require.Equal(t, completeMsg, msg.Type)
		require.Equal(t, "test_1", msg.ID)
	})
}

func TestWebsocketWithKeepAlive(t *testing.T) {
	h := testserver.New()
	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 100 * time.Millisecond,
	})

	srv := httptest.NewServer(h)
	defer srv.Close()

	c := wsConnect(srv.URL)
	defer c.Close()

	require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
	assert.Equal(t, connectionAckMsg, readOp(c).Type)
	assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)

	require.NoError(t, c.WriteJSON(&operationMessage{
		Type:    startMsg,
		ID:      "test_1",
		Payload: json.RawMessage(`{"query": "subscription { name }"}`),
	}))

	// keepalive
	msg := readOp(c)
	assert.Equal(t, connectionKeepAliveMsg, msg.Type)

	// server message
	h.SendNextSubscriptionMessage()
	msg = readOp(c)
	assert.Equal(t, dataMsg, msg.Type)

	// keepalive
	msg = readOp(c)
	assert.Equal(t, connectionKeepAliveMsg, msg.Type)
}

func TestWebsocketInitFunc(t *testing.T) {
	t.Run("accept connection if WebsocketInitFunc is NOT provided", func(t *testing.T) {
		h := testserver.New()
		h.AddTransport(transport.Websocket{})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))

		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
	})

	t.Run("accept connection if WebsocketInitFunc is provided and is accepting connection", func(t *testing.T) {
		h := testserver.New()
		h.AddTransport(transport.Websocket{
			InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
				return context.WithValue(ctx, ckey("newkey"), "newvalue"), nil
			},
		})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))

		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
	})

	t.Run("reject connection if WebsocketInitFunc is provided and is accepting connection", func(t *testing.T) {
		h := testserver.New()
		h.AddTransport(transport.Websocket{
			InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
				return ctx, errors.New("invalid init payload")
			},
		})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))

		msg := readOp(c)
		assert.Equal(t, connectionErrorMsg, msg.Type)
		assert.Equal(t, `{"message":"invalid init payload"}`, string(msg.Payload))
	})

	t.Run("can return context for request from WebsocketInitFunc", func(t *testing.T) {
		es := &graphql.ExecutableSchemaMock{
			ExecFunc: func(ctx context.Context) graphql.ResponseHandler {
				assert.Equal(t, "newvalue", ctx.Value(ckey("newkey")))
				return graphql.OneShot(&graphql.Response{Data: []byte(`{"empty":"ok"}`)})
			},
			SchemaFunc: func() *ast.Schema {
				return gqlparser.MustLoadSchema(&ast.Source{Input: `
				schema { query: Query }
				type Query {
					empty: String
				}
			`})
			},
		}
		h := handler.New(es)

		h.AddTransport(transport.Websocket{
			InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
				return context.WithValue(ctx, ckey("newkey"), "newvalue"), nil
			},
		})

		c := client.New(h)

		socket := c.Websocket("{ empty } ")
		defer socket.Close()
		var resp struct {
			Empty string
		}
		err := socket.Next(&resp)
		require.NoError(t, err)
		assert.Equal(t, "ok", resp.Empty)
	})

	t.Run("can set a deadline on a websocket connection and close it with a reason", func(t *testing.T) {
		h := testserver.New()
		var cancel func()
		h.AddTransport(transport.Websocket{
			InitFunc: func(ctx context.Context, _ transport.InitPayload) (newCtx context.Context, _ error) {
				newCtx, cancel = context.WithTimeout(transport.AppendCloseReason(ctx, "beep boop"), time.Millisecond*5)
				return
			},
		})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)

		// Cancel should contain an actual value now, so let's call it when we exit this scope (to make the linter happy)
		defer cancel()

		time.Sleep(time.Millisecond * 10)
		m := readOp(c)
		assert.Equal(t, m.Type, connectionErrorMsg)
		assert.Equal(t, string(m.Payload), `{"message":"beep boop"}`)
	})
}

func TestWebSocketInitTimeout(t *testing.T) {
	t.Run("times out if no init message is received within the configured duration", func(t *testing.T) {
		h := testserver.New()
		h.AddTransport(transport.Websocket{
			InitTimeout: 5 * time.Millisecond,
		})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		defer c.Close()

		var msg operationMessage
		err := c.ReadJSON(&msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})

	t.Run("keeps waiting for an init message if no time out is configured", func(t *testing.T) {
		h := testserver.New()
		h.AddTransport(transport.Websocket{})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		defer c.Close()

		done := make(chan interface{}, 1)
		go func() {
			var msg operationMessage
			_ = c.ReadJSON(&msg)
			done <- 1
		}()

		select {
		case <-done:
			assert.Fail(t, "web socket read operation finished while it shouldn't have")
		case <-time.After(100 * time.Millisecond):
			// Success! I guess? Can't really wait forever to see if the read waits forever...
		}
	})
}

func TestWebSocketErrorFunc(t *testing.T) {
	t.Run("the error handler gets called when an error occurs", func(t *testing.T) {
		errFuncCalled := make(chan bool, 1)
		h := testserver.New()
		h.AddTransport(transport.Websocket{
			ErrorFunc: func(_ context.Context, err error) {
				require.Error(t, err)
				assert.Equal(t, err.Error(), "invalid message received")
				errFuncCalled <- true
			},
		})

		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
		require.NoError(t, c.WriteMessage(websocket.TextMessage, []byte("mark my words, you will regret this")))

		select {
		case res := <-errFuncCalled:
			assert.True(t, res)
		case <-time.NewTimer(time.Millisecond * 20).C:
			assert.Fail(t, "The fail handler was not called in time")
		}
	})

	t.Run("init func errors do not call the error handler", func(t *testing.T) {
		h := testserver.New()
		h.AddTransport(transport.Websocket{
			InitFunc: func(ctx context.Context, _ transport.InitPayload) (context.Context, error) {
				return ctx, errors.New("this is not what we agreed upon")
			},
			ErrorFunc: func(_ context.Context, err error) {
				assert.Fail(t, "the error handler got called when it shouldn't have", "error: "+err.Error())
			},
		})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		time.Sleep(time.Millisecond * 20)
	})

	t.Run("init func context closes do not call the error handler", func(t *testing.T) {
		h := testserver.New()
		h.AddTransport(transport.Websocket{
			InitFunc: func(ctx context.Context, _ transport.InitPayload) (context.Context, error) {
				newCtx, cancel := context.WithCancel(ctx)
				time.AfterFunc(time.Millisecond*5, cancel)
				return newCtx, nil
			},
			ErrorFunc: func(_ context.Context, err error) {
				assert.Fail(t, "the error handler got called when it shouldn't have", "error: "+err.Error())
			},
		})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
		time.Sleep(time.Millisecond * 20)
	})

	t.Run("init func context deadlines do not call the error handler", func(t *testing.T) {
		h := testserver.New()
		var cancel func()
		h.AddTransport(transport.Websocket{
			InitFunc: func(ctx context.Context, _ transport.InitPayload) (newCtx context.Context, _ error) {
				newCtx, cancel = context.WithDeadline(ctx, time.Now().Add(time.Millisecond*5))
				return newCtx, nil
			},
			ErrorFunc: func(_ context.Context, err error) {
				assert.Fail(t, "the error handler got called when it shouldn't have", "error: "+err.Error())
			},
		})
		srv := httptest.NewServer(h)
		defer srv.Close()

		c := wsConnect(srv.URL)
		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)

		// Cancel should contain an actual value now, so let's call it when we exit this scope (to make the linter happy)
		defer cancel()

		time.Sleep(time.Millisecond * 20)
	})
}

func TestWebsocketGraphqltransportwsSubprotocol(t *testing.T) {
	initialize := func(ws transport.Websocket) (*testserver.TestServer, *httptest.Server) {
		h := testserver.New()
		h.AddTransport(ws)
		return h, httptest.NewServer(h)
	}

	t.Run("server acks init", func(t *testing.T) {
		_, srv := initialize(transport.Websocket{})
		defer srv.Close()

		c := wsConnectWithSubprocotol(srv.URL, graphqltransportwsSubprotocol)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsConnectionInitMsg}))
		assert.Equal(t, graphqltransportwsConnectionAckMsg, readOp(c).Type)
	})

	t.Run("client can receive data", func(t *testing.T) {
		handler, srv := initialize(transport.Websocket{})
		defer srv.Close()

		c := wsConnectWithSubprocotol(srv.URL, graphqltransportwsSubprotocol)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsConnectionInitMsg}))
		assert.Equal(t, graphqltransportwsConnectionAckMsg, readOp(c).Type)

		require.NoError(t, c.WriteJSON(&operationMessage{
			Type:    graphqltransportwsSubscribeMsg,
			ID:      "test_1",
			Payload: json.RawMessage(`{"query": "subscription { name }"}`),
		}))

		handler.SendNextSubscriptionMessage()
		msg := readOp(c)
		require.Equal(t, graphqltransportwsNextMsg, msg.Type, string(msg.Payload))
		require.Equal(t, "test_1", msg.ID, string(msg.Payload))
		require.Equal(t, `{"data":{"name":"test"}}`, string(msg.Payload))

		handler.SendNextSubscriptionMessage()
		msg = readOp(c)
		require.Equal(t, graphqltransportwsNextMsg, msg.Type, string(msg.Payload))
		require.Equal(t, "test_1", msg.ID, string(msg.Payload))
		require.Equal(t, `{"data":{"name":"test"}}`, string(msg.Payload))

		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsCompleteMsg, ID: "test_1"}))

		msg = readOp(c)
		require.Equal(t, graphqltransportwsCompleteMsg, msg.Type)
		require.Equal(t, "test_1", msg.ID)
	})

	t.Run("receives no graphql-ws keep alive messages", func(t *testing.T) {
		_, srv := initialize(transport.Websocket{KeepAlivePingInterval: 5 * time.Millisecond})
		defer srv.Close()

		c := wsConnectWithSubprocotol(srv.URL, graphqltransportwsSubprotocol)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsConnectionInitMsg}))
		assert.Equal(t, graphqltransportwsConnectionAckMsg, readOp(c).Type)

		// If the keep-alives are sent, this deadline will not be used, and no timeout error will be found
		c.SetReadDeadline(time.Now().UTC().Add(50 * time.Millisecond))
		var msg operationMessage
		err := c.ReadJSON(&msg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})
}

func TestWebsocketWithPingPongInterval(t *testing.T) {
	initialize := func(ws transport.Websocket) (*testserver.TestServer, *httptest.Server) {
		h := testserver.New()
		h.AddTransport(ws)
		return h, httptest.NewServer(h)
	}

	t.Run("client receives ping and responds with pong", func(t *testing.T) {
		_, srv := initialize(transport.Websocket{PingPongInterval: 10 * time.Millisecond})
		defer srv.Close()

		c := wsConnectWithSubprocotol(srv.URL, graphqltransportwsSubprotocol)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsConnectionInitMsg}))
		assert.Equal(t, graphqltransportwsConnectionAckMsg, readOp(c).Type)

		assert.Equal(t, graphqltransportwsPingMsg, readOp(c).Type)
		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsPongMsg}))
		assert.Equal(t, graphqltransportwsPingMsg, readOp(c).Type)
	})

	t.Run("client sends ping and expects pong", func(t *testing.T) {
		_, srv := initialize(transport.Websocket{PingPongInterval: 10 * time.Millisecond})
		defer srv.Close()

		c := wsConnectWithSubprocotol(srv.URL, graphqltransportwsSubprotocol)
		defer c.Close()

		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsConnectionInitMsg}))
		assert.Equal(t, graphqltransportwsConnectionAckMsg, readOp(c).Type)

		require.NoError(t, c.WriteJSON(&operationMessage{Type: graphqltransportwsPingMsg}))
		assert.Equal(t, graphqltransportwsPongMsg, readOp(c).Type)
	})

	t.Run("ping-pongs are not sent when the graphql-ws sub protocol is used", func(t *testing.T) {
		// Regression test
		// ---
		// Before the refactor, the code would try to convert a ping message to a graphql-ws message type
		// But since this message type does not exist in the graphql-ws sub protocol, it would fail

		_, srv := initialize(transport.Websocket{
			PingPongInterval:      5 * time.Millisecond,
			KeepAlivePingInterval: 10 * time.Millisecond,
		})
		defer srv.Close()

		// Create connection
		c := wsConnect(srv.URL)
		defer c.Close()

		// Initialize connection
		require.NoError(t, c.WriteJSON(&operationMessage{Type: connectionInitMsg}))
		assert.Equal(t, connectionAckMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)

		// Wait for a few more keep alives to be sure nothing goes wrong
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
		assert.Equal(t, connectionKeepAliveMsg, readOp(c).Type)
	})
}

func wsConnect(url string) *websocket.Conn {
	return wsConnectWithSubprocotol(url, "")
}

func wsConnectWithSubprocotol(url, subprocotol string) *websocket.Conn {
	h := make(http.Header)
	if subprocotol != "" {
		h.Add("Sec-WebSocket-Protocol", subprocotol)
	}

	c, resp, err := websocket.DefaultDialer.Dial(strings.ReplaceAll(url, "http://", "ws://"), h)
	if err != nil {
		panic(err)
	}
	_ = resp.Body.Close()

	return c
}

func writeRaw(conn *websocket.Conn, msg string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		panic(err)
	}
}

func readOp(conn *websocket.Conn) operationMessage {
	var msg operationMessage
	if err := conn.ReadJSON(&msg); err != nil {
		panic(err)
	}
	return msg
}

// copied out from websocket_graphqlws.go to keep these private

const (
	connectionInitMsg      = "connection_init"      // Client -> Server
	connectionTerminateMsg = "connection_terminate" // Client -> Server
	startMsg               = "start"                // Client -> Server
	stopMsg                = "stop"                 // Client -> Server
	connectionAckMsg       = "connection_ack"       // Server -> Client
	connectionErrorMsg     = "connection_error"     // Server -> Client
	dataMsg                = "data"                 // Server -> Client
	errorMsg               = "error"                // Server -> Client
	completeMsg            = "complete"             // Server -> Client
	connectionKeepAliveMsg = "ka"                   // Server -> Client
)

// copied out from websocket_graphql_transport_ws.go to keep these private

const (
	graphqltransportwsSubprotocol = "graphql-transport-ws"

	graphqltransportwsConnectionInitMsg = "connection_init"
	graphqltransportwsConnectionAckMsg  = "connection_ack"
	graphqltransportwsSubscribeMsg      = "subscribe"
	graphqltransportwsNextMsg           = "next"
	graphqltransportwsCompleteMsg       = "complete"
	graphqltransportwsPingMsg           = "ping"
	graphqltransportwsPongMsg           = "pong"
)

type operationMessage struct {
	Payload json.RawMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
}
