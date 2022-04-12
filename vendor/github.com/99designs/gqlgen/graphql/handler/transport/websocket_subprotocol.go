package transport

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
)

const (
	initMessageType messageType = iota
	connectionAckMessageType
	keepAliveMessageType
	connectionErrorMessageType
	connectionCloseMessageType
	startMessageType
	stopMessageType
	dataMessageType
	completeMessageType
	errorMessageType
	pingMessageType
	pongMessageType
)

var (
	supportedSubprotocols = []string{
		graphqlwsSubprotocol,
		graphqltransportwsSubprotocol,
	}

	errWsConnClosed = errors.New("websocket connection closed")
	errInvalidMsg   = errors.New("invalid message received")
)

type (
	messageType int
	message     struct {
		payload json.RawMessage
		id      string
		t       messageType
	}
	messageExchanger interface {
		NextMessage() (message, error)
		Send(m *message) error
	}
)

func (t messageType) String() string {
	var text string
	switch t {
	default:
		text = "unknown"
	case initMessageType:
		text = "init"
	case connectionAckMessageType:
		text = "connection ack"
	case keepAliveMessageType:
		text = "keep alive"
	case connectionErrorMessageType:
		text = "connection error"
	case connectionCloseMessageType:
		text = "connection close"
	case startMessageType:
		text = "start"
	case stopMessageType:
		text = "stop subscription"
	case dataMessageType:
		text = "data"
	case completeMessageType:
		text = "complete"
	case errorMessageType:
		text = "error"
	case pingMessageType:
		text = "ping"
	case pongMessageType:
		text = "pong"
	}
	return text
}

func contains(list []string, elem string) bool {
	for _, e := range list {
		if e == elem {
			return true
		}
	}

	return false
}

func (t *Websocket) injectGraphQLWSSubprotocols() {
	// the list of subprotocols is specified by the consumer of the Websocket struct,
	// in order to preserve backward compatibility, we inject the graphql specific subprotocols
	// at runtime
	if !t.didInjectSubprotocols {
		defer func() {
			t.didInjectSubprotocols = true
		}()

		for _, subprotocol := range supportedSubprotocols {
			if !contains(t.Upgrader.Subprotocols, subprotocol) {
				t.Upgrader.Subprotocols = append(t.Upgrader.Subprotocols, subprotocol)
			}
		}
	}
}

func handleNextReaderError(err error) error {
	// TODO: should we consider all closure scenarios here for the ws connection?
	// for now we only list the error codes from the previous implementation
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
		return errWsConnClosed
	}

	return err
}
