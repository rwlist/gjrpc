package jsonrpc

import (
	"context"
	"encoding/json"
)

const Version = "2.0"

type ID json.RawMessage

func (id ID) MarshalJSON() ([]byte, error) {
	return json.RawMessage(id).MarshalJSON()
}

func (id *ID) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*json.RawMessage)(id))
}

// Request object from specification: https://www.jsonrpc.org/specification
type Request struct {
	// MUST be exactly "2.0".
	Version string `json:"jsonrpc"`

	// Method names that begin with the word rpc followed by a period character (U+002E or ASCII 46)
	// are reserved for rpc-internal methods and extensions and MUST NOT be used for anything else.
	Method string `json:"method"`

	// Params is a Structured value that holds the parameter values to be used during the invocation
	// of the method. This member MAY be omitted.
	Params json.RawMessage `json:"params"`

	// ID is an identifier established by the Client that MUST contain a String, Number,
	// or NULL value if included.
	ID ID `json:"id"`

	// TODO: should we add Authorization from http header, or context is enough?

	// Context with additional metadata for the request.
	Context context.Context `json:"-"`
}
