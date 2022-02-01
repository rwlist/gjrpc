// Code generated by gjrpc. DO NOT EDIT.

package gen_server

import (
	"encoding/json"
	proto "github.com/rwlist/gjrpc/examples/1_simple_service/proto"
	jsonrpc "github.com/rwlist/gjrpc/pkg/jsonrpc"
	"strings"
)

type Router struct {
	handlers     Handlers
	convertError jsonrpc.ErrorConverter
}

func NewRouter(handlers Handlers, convertError jsonrpc.ErrorConverter) *Router {
	if convertError == nil {
		convertError = jsonrpc.DefaultErrorConverter
	}
	return &Router{
		convertError: convertError,
		handlers:     handlers,
	}
}

func (r *Router) notFound() (jsonrpc.Result, *jsonrpc.Error) {
	return nil, &jsonrpc.MethodNotFound
}

func (r *Router) Handle(req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	path := strings.Split(req.Method, ".")
	return r.handle(path, req)
}

func (r *Router) handle(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		return r.notFound()
	}
	switch path[0] {
	case "inventory":
		return r.handleInventory(path[1:], req)
	}
	return r.notFound()
}

func (r *Router) handleInventory(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		return r.notFound()
	}
	switch path[0] {
	case "bar":
		return r.handleInventoryBar(path[1:], req)
	case "foo":
		return r.handleInventoryFoo(path[1:], req)
	}
	return r.notFound()
}

func (r *Router) handleInventoryBar(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		var request proto.Bar
		if err := json.Unmarshal(req.Params, &request); err != nil {
			return r.convertError(err)
		}
		err := r.handlers.Inventory.Bar(request)
		if err != nil {
			return r.convertError(err)
		}
		return struct{}{}, nil
	}
	return r.notFound()
}

func (r *Router) handleInventoryFoo(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		res, err := r.handlers.Inventory.Foo()
		if err != nil {
			return r.convertError(err)
		}
		return res, nil
	}
	return r.notFound()
}
