package jsonrpc

import (
	"encoding/json"
	"fmt"
)

var (
	ParseError     = Error{Code: -32700, Message: "Parse error"}
	InvalidRequest = Error{Code: -32600, Message: "Invalid Request"}
	MethodNotFound = Error{Code: -32601, Message: "Method not found"}
	InvalidParams  = Error{Code: -32602, Message: "Invalid params"}
	InternalError  = Error{Code: -32603, Message: "Internal error"}
	UnknownError   = Error{Code: -32000, Message: "Unknown error"}
)

type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("jsonrpc.Error(code=%d message=%s data=%s)", e.Code, e.Message, e.Data)
}

func (e Error) withJSON(data interface{}) *Error {
	b, err := json.Marshal(data)
	if err != nil {
		println("jsonrpc.Error.withJSON: ", err.Error())
		return &e
	}
	e.Data = json.RawMessage(b)
	return &e
}

func (e *Error) WithError(err error) *Error {
	return e.withJSON(err.Error())
}

func (e Error) WithBytes(data []byte) *Error {
	return e.withJSON(data)
}

func (e Error) WithString(data string) *Error {
	return e.withJSON(data)
}
