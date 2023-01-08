package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

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

	resp := h.Handler.Handle(&req)
	_ = json.NewEncoder(w).Encode(resp)
}
