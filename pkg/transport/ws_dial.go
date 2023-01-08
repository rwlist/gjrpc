package transport

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

// DialOnce tries to connect to the server and returns a connection.
// Currently it does not handle disconnections at all.
// TODO: handle reconnections and process errors better.
func DialOnce(ctx context.Context, url string, callback jsonrpc.Handler) (*WSConn, error) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	// TODO: wrap a callback in a safer middleware (verifies the protocol)

	conn := NewWSConn(ws)
	go conn.StartReadLoop(ctx, callback)
	return conn, nil
}
