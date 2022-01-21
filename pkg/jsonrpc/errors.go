package jsonrpc

var (
	ParseError     = Error{Code: -32700, Message: "Parse error"}
	InvalidRequest = Error{Code: -32600, Message: "Invalid Request"}
	MethodNotFound = Error{Code: -32601, Message: "Method not found"}
	ForbiddenError = Error{Code: -1000, Message: "Forbidden"}
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
