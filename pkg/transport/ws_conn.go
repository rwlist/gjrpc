package transport

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

// WSConn is a wrapper around websocket.Conn, which implements jsonrpc.Transport interface.
type WSConn struct {
	conn *websocket.Conn
	// write lock
	wlock sync.Mutex

	// outgoing requests mapping
	mu       sync.Mutex
	requests map[[maxIDLength]byte]chan<- *jsonrpc.Response
}

var _ jsonrpc.Transport = (*WSConn)(nil)

// NewWSConn creates a JSON-RPC over WebSocket connection and starts processing incoming
// requests.
func NewWSConn(ws *websocket.Conn) *WSConn {
	conn := &WSConn{
		conn:     ws,
		requests: make(map[[maxIDLength]byte]chan<- *jsonrpc.Response),
	}

	return conn
}

// StartReadLoop must be called exactly once. It starts reading messages from the connection
// and processing them. It returns when the connection is closed.
func (c *WSConn) StartReadLoop(ctx context.Context, requestsHandler jsonrpc.Handler) {
	// TODO: handle conn close, use context
	for {
		_, b, err := c.conn.ReadMessage()
		if err != nil {
			// TODO: shutdown the connection
			return
		}

		mt := jsonrpc.DetectMessageType(b)

		switch mt {
		case jsonrpc.MessageTypeRequest:
			go c.handleRequest(ctx, b, requestsHandler)
		case jsonrpc.MessageTypeResponse:
			go c.handleResponse(b)
		default:
			// TODO: failed to detect message type, what to do?
			println("failed to detect message type", string(b))
			continue
		}
	}
}

func (c *WSConn) handleRequest(ctx context.Context, b []byte, handler jsonrpc.Handler) {
	var req jsonrpc.Request
	if err := json.Unmarshal(b, &req); err != nil {
		// TODO: log error and shutdown the connection?
		return
	}
	req.Context = ctx

	if handler == nil {
		handler = jsonrpc.HandlerFunc(defaultNotFoundHandler)
	}
	resp := handler.Handle(&req)
	err := c.sendJSON(resp)
	if err != nil {
		println("failed to send response", err.Error())
		// TODO: log error
		return
	}
}

// outgoing requests

// Call implements jsonrpc.Transport interface. It sends the request and waits
// for the response if the request is not a notification.
func (c *WSConn) SendRequest(req *jsonrpc.Request) (*jsonrpc.Response, error) {
	if req.ID == nil {
		// it's a notification, just send it and return
		if err := c.sendJSON(req); err != nil {
			return nil, err
		}
		return nil, nil
	}

	id, err := transformID(req.ID)
	if err != nil {
		return nil, err
	}

	ch := make(chan *jsonrpc.Response, 1)
	ok := c.startRequest(id, ch)
	if !ok {
		return nil, errors.New("duplicate request")
	}

	if err := c.sendJSON(req); err != nil {
		return nil, err
	}

	// TODO: timeouts, use context, watch for conn close, etc.
	res := <-ch
	return res, nil
}

func (c *WSConn) startRequest(id wsID, ch chan<- *jsonrpc.Response) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.requests[id]
	if ok {
		// duplicate
		return false
	}

	c.requests[id] = ch
	return true
}

func (c *WSConn) finishRequest(id wsID) chan<- *jsonrpc.Response {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.requests[id]
	if !ok {
		return nil
	}
	delete(c.requests, id)

	return ch
}

func (c *WSConn) sendJSON(i interface{}) error {
	c.wlock.Lock()
	defer c.wlock.Unlock()
	return c.conn.WriteJSON(i)
}

func (c *WSConn) handleResponse(b []byte) {
	var resp jsonrpc.Response
	if err := json.Unmarshal(b, &resp); err != nil {
		// TODO: log error
		return
	}

	id, err := transformID(resp.ID)
	if err != nil {
		// TODO: log error
		return
	}

	ch := c.finishRequest(id)
	if ch == nil {
		// TODO: log "request not found"
		return
	}

	ch <- &resp
}
