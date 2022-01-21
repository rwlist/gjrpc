package gjserver

import (
	"context"
	"encoding/json"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"net/http"
)

type ctxKey string

const (
	ctxHttpRequest = ctxKey("http.Request")
)

func HttpRequest(ctx context.Context) *http.Request {
	v, ok := ctx.Value(ctxHttpRequest).(*http.Request)
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
	ctx := context.WithValue(r.Context(), ctxHttpRequest, r)

	var req jsonrpc.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Context = ctx

	res, err := h.Handler(&req)
	if err != nil {
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}
