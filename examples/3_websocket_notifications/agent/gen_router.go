// Code generated by gjrpc. DO NOT EDIT.

package agent

import (
	"encoding/json"
	jsonrpc "github.com/rwlist/gjrpc/pkg/jsonrpc"
	"strings"
	proto "wsexample/proto"
)

type Router struct {
	handlers     *Handlers
	convertError jsonrpc.ErrorConverter
}

func NewRouter(handlers *Handlers, convertError jsonrpc.ErrorConverter) *Router {
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

func (r *Router) Handle(req *jsonrpc.Request) *jsonrpc.Response {
	path := strings.Split(req.Method, ".")
	res, e := r.handle(path, req)
	var result json.RawMessage
	if e == nil {
		var err error
		result, err = json.Marshal(res)
		if err != nil {
			_, e = r.convertError(err)
		}
	}
	return jsonrpc.NewResponse(req, result, e)
}

func (r *Router) handle(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		return r.notFound()
	}
	switch path[0] {
	case "agent":
		return r.handleAgent(path[1:], req)
	}
	return r.notFound()
}

func (r *Router) handleAgent(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		return r.notFound()
	}
	switch path[0] {
	case "executev1":
		return r.handleAgentExecutev1(path[1:], req)
	case "status":
		return r.handleAgentStatus(path[1:], req)
	}
	return r.notFound()
}

func (r *Router) handleAgentExecutev1(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		var request proto.ExecuteV1Request
		if err := json.Unmarshal(req.Params, &request); err != nil {
			return r.convertError(err)
		}
		res, err := r.handlers.Agent.ExecuteV1(req.Context, &request)
		if err != nil {
			return r.convertError(err)
		}
		return res, nil
	}
	return r.notFound()
}

func (r *Router) handleAgentStatus(path []string, req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	if len(path) == 0 {
		var request proto.StatusRequest
		if err := json.Unmarshal(req.Params, &request); err != nil {
			return r.convertError(err)
		}
		res, err := r.handlers.Agent.Status(req.Context, &request)
		if err != nil {
			return r.convertError(err)
		}
		return res, nil
	}
	return r.notFound()
}
