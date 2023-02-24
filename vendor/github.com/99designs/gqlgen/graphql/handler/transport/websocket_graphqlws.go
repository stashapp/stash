package transport

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// https://github.com/enisdenjo/graphql-ws/blob/master/PROTOCOL.md
const (
	graphqlwsSubprotocol = "graphql-ws"

	graphqlwsConnectionInitMsg      = graphqlwsMessageType("connection_init")
	graphqlwsConnectionTerminateMsg = graphqlwsMessageType("connection_terminate")
	graphqlwsStartMsg               = graphqlwsMessageType("start")
	graphqlwsStopMsg                = graphqlwsMessageType("stop")
	graphqlwsConnectionAckMsg       = graphqlwsMessageType("connection_ack")
	graphqlwsConnectionErrorMsg     = graphqlwsMessageType("connection_error")
	graphqlwsDataMsg                = graphqlwsMessageType("data")
	graphqlwsErrorMsg               = graphqlwsMessageType("error")
	graphqlwsCompleteMsg            = graphqlwsMessageType("complete")
	graphqlwsConnectionKeepAliveMsg = graphqlwsMessageType("ka")
)

var allGraphqlwsMessageTypes = []graphqlwsMessageType{
	graphqlwsConnectionInitMsg,
	graphqlwsConnectionTerminateMsg,
	graphqlwsStartMsg,
	graphqlwsStopMsg,
	graphqlwsConnectionAckMsg,
	graphqlwsConnectionErrorMsg,
	graphqlwsDataMsg,
	graphqlwsErrorMsg,
	graphqlwsCompleteMsg,
	graphqlwsConnectionKeepAliveMsg,
}

type (
	graphqlwsMessageExchanger struct {
		c *websocket.Conn
	}

	graphqlwsMessage struct {
		Payload json.RawMessage      `json:"payload,omitempty"`
		ID      string               `json:"id,omitempty"`
		Type    graphqlwsMessageType `json:"type"`
		noOp    bool
	}

	graphqlwsMessageType string
)

func (me graphqlwsMessageExchanger) NextMessage() (message, error) {
	_, r, err := me.c.NextReader()
	if err != nil {
		return message{}, handleNextReaderError(err)
	}

	var graphqlwsMessage graphqlwsMessage
	if err := jsonDecode(r, &graphqlwsMessage); err != nil {
		return message{}, errInvalidMsg
	}

	return graphqlwsMessage.toMessage()
}

func (me graphqlwsMessageExchanger) Send(m *message) error {
	msg := &graphqlwsMessage{}
	if err := msg.fromMessage(m); err != nil {
		return err
	}

	if msg.noOp {
		return nil
	}

	return me.c.WriteJSON(msg)
}

func (t *graphqlwsMessageType) UnmarshalText(text []byte) (err error) {
	var found bool
	for _, candidate := range allGraphqlwsMessageTypes {
		if string(candidate) == string(text) {
			*t = candidate
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("invalid message type %s", string(text))
	}

	return err
}

func (t graphqlwsMessageType) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (m graphqlwsMessage) toMessage() (message, error) {
	var t messageType
	var err error
	switch m.Type {
	default:
		err = fmt.Errorf("invalid client->server message type %s", m.Type)
	case graphqlwsConnectionInitMsg:
		t = initMessageType
	case graphqlwsConnectionTerminateMsg:
		t = connectionCloseMessageType
	case graphqlwsStartMsg:
		t = startMessageType
	case graphqlwsStopMsg:
		t = stopMessageType
	case graphqlwsConnectionAckMsg:
		t = connectionAckMessageType
	case graphqlwsConnectionErrorMsg:
		t = connectionErrorMessageType
	case graphqlwsDataMsg:
		t = dataMessageType
	case graphqlwsErrorMsg:
		t = errorMessageType
	case graphqlwsCompleteMsg:
		t = completeMessageType
	case graphqlwsConnectionKeepAliveMsg:
		t = keepAliveMessageType
	}

	return message{
		payload: m.Payload,
		id:      m.ID,
		t:       t,
	}, err
}

func (m *graphqlwsMessage) fromMessage(msg *message) (err error) {
	m.ID = msg.id
	m.Payload = msg.payload

	switch msg.t {
	default:
		err = fmt.Errorf("invalid server->client message type %s", msg.t)
	case initMessageType:
		m.Type = graphqlwsConnectionInitMsg
	case connectionAckMessageType:
		m.Type = graphqlwsConnectionAckMsg
	case keepAliveMessageType:
		m.Type = graphqlwsConnectionKeepAliveMsg
	case connectionErrorMessageType:
		m.Type = graphqlwsConnectionErrorMsg
	case connectionCloseMessageType:
		m.Type = graphqlwsConnectionTerminateMsg
	case startMessageType:
		m.Type = graphqlwsStartMsg
	case stopMessageType:
		m.Type = graphqlwsStopMsg
	case dataMessageType:
		m.Type = graphqlwsDataMsg
	case completeMessageType:
		m.Type = graphqlwsCompleteMsg
	case errorMessageType:
		m.Type = graphqlwsErrorMsg
	case pingMessageType:
		m.noOp = true
	case pongMessageType:
		m.noOp = true
	}

	return err
}
