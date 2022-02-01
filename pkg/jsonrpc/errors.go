package jsonrpc

var (
	ParseError     = Error{Code: -32700, Message: "Parse error"}
	InvalidRequest = Error{Code: -32600, Message: "Invalid Request"}
	MethodNotFound = Error{Code: -32601, Message: "Method not found"}
	InvalidParams  = Error{Code: -32602, Message: "Invalid params"}
	InternalError  = Error{Code: -32603, Message: "Internal error"}
	UnknownError   = Error{Code: -32000, Message: "Unknown error"}
)

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e Error) WithData(data interface{}) *Error {
	e.Data = data
	return &e
}

// ErrorConverter is for converting error to jsonrpc error.
// TODO: think of error -> ErrorContext, with router path, etc.
type ErrorConverter func(error) (Result, *Error)

func DefaultErrorConverter(err error) (Result, *Error) {
	return nil, UnknownError.WithData(err)
}
