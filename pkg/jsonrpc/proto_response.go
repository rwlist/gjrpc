package jsonrpc

import (
	"encoding/json"
	"errors"
)

var (
	errNoID      = errors.New("response must have id")
	errEmpty     = errors.New("response must have result or error")
	errAmbiguous = errors.New("error response must not have result")
)

// Response object from specification: https://www.jsonrpc.org/specification
type Response struct {
	// MUST be exactly "2.0".
	Version string `json:"jsonrpc"`

	// This member is REQUIRED on success.
	Result json.RawMessage `json:"result,omitempty"`

	// This member is REQUIRED on error.
	Error *Error `json:"error,omitempty"`

	// ID is an identifier established by the Client that MUST contain a String, Number.
	// Cannot be NULL, because it's used only for notifications without response.
	ID ID `json:"id,omitempty"`
}

// NewResponse creates new JSON-RPC response object. If request is notification,
// response will be nil.
func NewResponse(request *Request, result json.RawMessage, e *Error) *Response {
	if request.IsNotification() {
		return nil
	}

	return &Response{
		Version: "2.0",
		Result:  result,
		Error:   e,
		ID:      request.ID,
	}
}

// Validate checks if response is valid, according to specification.
func (r *Response) Validate() error {
	if r.ID == nil {
		return errNoID
	}

	if r.Result == nil && r.Error == nil {
		return errEmpty
	}

	if r.Result != nil && r.Error != nil {
		return errAmbiguous
	}

	return nil
}
