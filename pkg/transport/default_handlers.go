package transport

import "github.com/rwlist/gjrpc/pkg/jsonrpc"

// defaultNotFoundHandler always returns MethodNotFound error.
func defaultNotFoundHandler(req *jsonrpc.Request) *jsonrpc.Response {
	return jsonrpc.NewResponse(req, nil, jsonrpc.MethodNotFound.WithBytes(nil))
}
