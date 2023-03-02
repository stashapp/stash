package transport

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// https://github.com/enisdenjo/graphql-ws/blob/master/PROTOCOL.md
const (
	graphqltransportwsSubprotocol = "graphql-transport-ws"

	graphqltransportwsConnectionInitMsg = graphqltransportwsMessageType("connection_init")
	graphqltransportwsConnectionAckMsg  = graphqltransportwsMessageType("connection_ack")
	graphqltransportwsSubscribeMsg      = graphqltransportwsMessageType("subscribe")
	graphqltransportwsNextMsg           = graphqltransportwsMessageType("next")
	graphqltransportwsErrorMsg          = graphqltransportwsMessageType("error")
	graphqltransportwsCompleteMsg       = graphqltransportwsMessageType("complete")
	graphqltransportwsPingMsg           = graphqltransportwsMessageType("ping")
	graphqltransportwsPongMsg           = graphqltransportwsMessageType("pong")
)

var allGraphqltransportwsMessageTypes = []graphqltransportwsMessageType{
	graphqltransportwsConnectionInitMsg,
	graphqltransportwsConnectionAckMsg,
	graphqltransportwsSubscribeMsg,
	graphqltransportwsNextMsg,
	graphqltransportwsErrorMsg,
	graphqltransportwsCompleteMsg,
	graphqltransportwsPingMsg,
	graphqltransportwsPongMsg,
}

type (
	graphqltransportwsMessageExchanger struct {
		c *websocket.Conn
	}

	graphqltransportwsMessage struct {
		Payload json.RawMessage               `json:"payload,omitempty"`
		ID      string                        `json:"id,omitempty"`
		Type    graphqltransportwsMessageType `json:"type"`
		noOp    bool
	}

	graphqltransportwsMessageType string
)

func (me graphqltransportwsMessageExchanger) NextMessage() (message, error) {
	_, r, err := me.c.NextReader()
	if err != nil {
		return message{}, handleNextReaderError(err)
	}

	var graphqltransportwsMessage graphqltransportwsMessage
	if err := jsonDecode(r, &graphqltransportwsMessage); err != nil {
		return message{}, errInvalidMsg
	}

	return graphqltransportwsMessage.toMessage()
}

func (me graphqltransportwsMessageExchanger) Send(m *message) error {
	msg := &graphqltransportwsMessage{}
	if err := msg.fromMessage(m); err != nil {
		return err
	}

	if msg.noOp {
		return nil
	}

	return me.c.WriteJSON(msg)
}

func (t *graphqltransportwsMessageType) UnmarshalText(text []byte) (err error) {
	var found bool
	for _, candidate := range allGraphqltransportwsMessageTypes {
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

func (t graphqltransportwsMessageType) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (m graphqltransportwsMessage) toMessage() (message, error) {
	var t messageType
	var err error
	switch m.Type {
	default:
		err = fmt.Errorf("invalid client->server message type %s", m.Type)
	case graphqltransportwsConnectionInitMsg:
		t = initMessageType
	case graphqltransportwsSubscribeMsg:
		t = startMessageType
	case graphqltransportwsCompleteMsg:
		t = stopMessageType
	case graphqltransportwsPingMsg:
		t = pingMessageType
	case graphqltransportwsPongMsg:
		t = pongMessageType
	}

	return message{
		payload: m.Payload,
		id:      m.ID,
		t:       t,
	}, err
}

func (m *graphqltransportwsMessage) fromMessage(msg *message) (err error) {
	m.ID = msg.id
	m.Payload = msg.payload

	switch msg.t {
	default:
		err = fmt.Errorf("invalid server->client message type %s", msg.t)
	case connectionAckMessageType:
		m.Type = graphqltransportwsConnectionAckMsg
	case keepAliveMessageType:
		m.noOp = true
	case connectionErrorMessageType:
		m.noOp = true
	case dataMessageType:
		m.Type = graphqltransportwsNextMsg
	case completeMessageType:
		m.Type = graphqltransportwsCompleteMsg
	case errorMessageType:
		m.Type = graphqltransportwsErrorMsg
	case pingMessageType:
		m.Type = graphqltransportwsPingMsg
	case pongMessageType:
		m.Type = graphqltransportwsPongMsg
	}

	return err
}
