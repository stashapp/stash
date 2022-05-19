package client

import "encoding/json"

// RawJsonError is a json formatted error from a GraphQL server.
type RawJsonError struct {
	json.RawMessage
}

func (r RawJsonError) Error() string {
	return string(r.RawMessage)
}
