package jsonrpc

type Transport interface {
	// SendRequest sends a request to the server. Transport should use req.Context.
	//
	// If req.ID is nil, the request is a notification. Response will always
	// be nil for notification.
	//
	// If req.ID is not nil, the request is a request with a response.
	// Response will be non-nil if err is nil.
	//
	// Transport will return an error in case of network and other issues.
	SendRequest(req *Request) (*Response, error)
}

// Handler contains a single Handle function. It should return a valid Response
// or nil if the request is a notification.
type Handler interface {
	Handle(req *Request) *Response
}

// HandlerFunc is a function that implements Handler interface.
type HandlerFunc func(req *Request) *Response

func (f HandlerFunc) Handle(req *Request) *Response {
	return f(req)
}

// Middleware is a function that takes a Handler and returns a Handler.
type Middleware func(Handler) Handler
