package transport

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

type HandlerWebSocket struct {
	Handler jsonrpc.Handler
	// TODO: make an interface for upgrader and conn
	Upgrader *websocket.Upgrader
}

func NewHandlerWebSocket(handler jsonrpc.Handler) *HandlerWebSocket {
	return &HandlerWebSocket{
		Handler:  handler,
		Upgrader: &websocket.Upgrader{},
	}
}

func (h *HandlerWebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxHttpRequest, r)
	ctx = context.WithValue(ctx, ctxHttpResponse, w)

	// TODO: call user-provided function for more context

	ws, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: log error?
		// Upgrade have already written an error response.
		return
	}
	defer ws.Close()

	// create a connection
	conn := NewWSConn(ws)
	ctx = context.WithValue(ctx, ctxWSConn, conn)

	// start reading messages from websocket
	conn.StartReadLoop(ctx, h.Handler)
}
