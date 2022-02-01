package gjserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

type ctxKey string

const (
	ctxHttpRequest  = ctxKey("http.Request")
	ctxHttpResponse = ctxKey("http.ResponseWriter")
)

func HttpRequest(ctx context.Context) *http.Request {
	v, ok := ctx.Value(ctxHttpRequest).(*http.Request)
	if !ok {
		return nil
	}
	return v
}

func HttpResponse(ctx context.Context) http.ResponseWriter {
	v, ok := ctx.Value(ctxHttpResponse).(http.ResponseWriter)
	if !ok {
		return nil
	}
	return v
}

// https://www.jsonrpc.org/historical/json-rpc-over-http.html

type HandlerHTTP struct {
	Handler jsonrpc.Handler
}

func (h *HandlerHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxHttpRequest, r)
	ctx = context.WithValue(ctx, ctxHttpResponse, w)

	var req jsonrpc.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Context = ctx

	res, err := h.Handler(&req)
	resp := jsonrpc.NewResponse(&req, res, err)
	_ = json.NewEncoder(w).Encode(resp)
}
