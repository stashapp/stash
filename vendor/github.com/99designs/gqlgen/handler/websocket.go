package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/golang-lru"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/gqlerror"
	"github.com/vektah/gqlparser/validator"
)

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

type operationMessage struct {
	Payload json.RawMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
}

type wsConnection struct {
	ctx             context.Context
	conn            *websocket.Conn
	exec            graphql.ExecutableSchema
	active          map[string]context.CancelFunc
	mu              sync.Mutex
	cfg             *Config
	cache           *lru.Cache
	keepAliveTicker *time.Ticker

	initPayload InitPayload
}

func connectWs(exec graphql.ExecutableSchema, w http.ResponseWriter, r *http.Request, cfg *Config, cache *lru.Cache) {
	ws, err := cfg.upgrader.Upgrade(w, r, http.Header{
		"Sec-Websocket-Protocol": []string{"graphql-ws"},
	})
	if err != nil {
		log.Printf("unable to upgrade %T to websocket %s: ", w, err.Error())
		sendErrorf(w, http.StatusBadRequest, "unable to upgrade")
		return
	}

	conn := wsConnection{
		active: map[string]context.CancelFunc{},
		exec:   exec,
		conn:   ws,
		ctx:    r.Context(),
		cfg:    cfg,
		cache:  cache,
	}

	if !conn.init() {
		return
	}

	conn.run()
}

func (c *wsConnection) init() bool {
	message := c.readOp()
	if message == nil {
		c.close(websocket.CloseProtocolError, "decoding error")
		return false
	}

	switch message.Type {
	case connectionInitMsg:
		if len(message.Payload) > 0 {
			c.initPayload = make(InitPayload)
			err := json.Unmarshal(message.Payload, &c.initPayload)
			if err != nil {
				return false
			}
		}

		c.write(&operationMessage{Type: connectionAckMsg})
	case connectionTerminateMsg:
		c.close(websocket.CloseNormalClosure, "terminated")
		return false
	default:
		c.sendConnectionError("unexpected message %s", message.Type)
		c.close(websocket.CloseProtocolError, "unexpected message")
		return false
	}

	return true
}

func (c *wsConnection) write(msg *operationMessage) {
	c.mu.Lock()
	c.conn.WriteJSON(msg)
	c.mu.Unlock()
}

func (c *wsConnection) run() {
	// We create a cancellation that will shutdown the keep-alive when we leave
	// this function.
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	// Create a timer that will fire every interval to keep the connection alive.
	if c.cfg.connectionKeepAlivePingInterval != 0 {
		c.mu.Lock()
		c.keepAliveTicker = time.NewTicker(c.cfg.connectionKeepAlivePingInterval)
		c.mu.Unlock()

		go c.keepAlive(ctx)
	}

	for {
		message := c.readOp()
		if message == nil {
			return
		}

		switch message.Type {
		case startMsg:
			if !c.subscribe(message) {
				return
			}
		case stopMsg:
			c.mu.Lock()
			closer := c.active[message.ID]
			c.mu.Unlock()
			if closer == nil {
				c.sendError(message.ID, gqlerror.Errorf("%s is not running, cannot stop", message.ID))
				continue
			}

			closer()
		case connectionTerminateMsg:
			c.close(websocket.CloseNormalClosure, "terminated")
			return
		default:
			c.sendConnectionError("unexpected message %s", message.Type)
			c.close(websocket.CloseProtocolError, "unexpected message")
			return
		}
	}
}

func (c *wsConnection) keepAlive(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.keepAliveTicker.Stop()
			return
		case <-c.keepAliveTicker.C:
			c.write(&operationMessage{Type: connectionKeepAliveMsg})
		}
	}
}

func (c *wsConnection) subscribe(message *operationMessage) bool {
	var reqParams params
	if err := jsonDecode(bytes.NewReader(message.Payload), &reqParams); err != nil {
		c.sendConnectionError("invalid json")
		return false
	}

	var (
		doc      *ast.QueryDocument
		cacheHit bool
	)
	if c.cache != nil {
		val, ok := c.cache.Get(reqParams.Query)
		if ok {
			doc = val.(*ast.QueryDocument)
			cacheHit = true
		}
	}
	if !cacheHit {
		var qErr gqlerror.List
		doc, qErr = gqlparser.LoadQuery(c.exec.Schema(), reqParams.Query)
		if qErr != nil {
			c.sendError(message.ID, qErr...)
			return true
		}
		if c.cache != nil {
			c.cache.Add(reqParams.Query, doc)
		}
	}

	op := doc.Operations.ForName(reqParams.OperationName)
	if op == nil {
		c.sendError(message.ID, gqlerror.Errorf("operation %s not found", reqParams.OperationName))
		return true
	}

	vars, err := validator.VariableValues(c.exec.Schema(), op, reqParams.Variables)
	if err != nil {
		c.sendError(message.ID, err)
		return true
	}
	reqCtx := c.cfg.newRequestContext(c.exec, doc, op, reqParams.Query, vars)
	ctx := graphql.WithRequestContext(c.ctx, reqCtx)

	if c.initPayload != nil {
		ctx = withInitPayload(ctx, c.initPayload)
	}

	if op.Operation != ast.Subscription {
		var result *graphql.Response
		if op.Operation == ast.Query {
			result = c.exec.Query(ctx, op)
		} else {
			result = c.exec.Mutation(ctx, op)
		}

		c.sendData(message.ID, result)
		c.write(&operationMessage{ID: message.ID, Type: completeMsg})
		return true
	}

	ctx, cancel := context.WithCancel(ctx)
	c.mu.Lock()
	c.active[message.ID] = cancel
	c.mu.Unlock()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				userErr := reqCtx.Recover(ctx, r)
				c.sendError(message.ID, &gqlerror.Error{Message: userErr.Error()})
			}
		}()
		next := c.exec.Subscription(ctx, op)
		for result := next(); result != nil; result = next() {
			c.sendData(message.ID, result)
		}

		c.write(&operationMessage{ID: message.ID, Type: completeMsg})

		c.mu.Lock()
		delete(c.active, message.ID)
		c.mu.Unlock()
		cancel()
	}()

	return true
}

func (c *wsConnection) sendData(id string, response *graphql.Response) {
	b, err := json.Marshal(response)
	if err != nil {
		c.sendError(id, gqlerror.Errorf("unable to encode json response: %s", err.Error()))
		return
	}

	c.write(&operationMessage{Type: dataMsg, ID: id, Payload: b})
}

func (c *wsConnection) sendError(id string, errors ...*gqlerror.Error) {
	var errs []error
	for _, err := range errors {
		errs = append(errs, err)
	}
	b, err := json.Marshal(errs)
	if err != nil {
		panic(err)
	}
	c.write(&operationMessage{Type: errorMsg, ID: id, Payload: b})
}

func (c *wsConnection) sendConnectionError(format string, args ...interface{}) {
	b, err := json.Marshal(&gqlerror.Error{Message: fmt.Sprintf(format, args...)})
	if err != nil {
		panic(err)
	}

	c.write(&operationMessage{Type: connectionErrorMsg, Payload: b})
}

func (c *wsConnection) readOp() *operationMessage {
	_, r, err := c.conn.NextReader()
	if err != nil {
		c.sendConnectionError("invalid json")
		return nil
	}
	message := operationMessage{}
	if err := jsonDecode(r, &message); err != nil {
		c.sendConnectionError("invalid json")
		return nil
	}

	return &message
}

func (c *wsConnection) close(closeCode int, message string) {
	c.mu.Lock()
	_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, message))
	c.mu.Unlock()
	_ = c.conn.Close()
}
