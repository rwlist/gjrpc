package jsonrpc

import (
	"context"
	"encoding/json"
	"strconv"
	"sync/atomic"
)

// ClientConnInterface defines the functions clients need to perform RPCs.
// It is implemented by *ClientConn, and is only intended to be referenced
// by the generated code.
type ClientConnInterface interface {
	// Call performs a unary RPC. If reply is not nil, it returns after the response
	// is received and unmarshalled into reply. If reply is nil, it returns after
	// the request was successfully sent.
	Call(ctx context.Context, method string, params interface{}, reply interface{}) error
}

type BasicClientConn struct {
	transport Transport
	seq       uint64
}

func NewClientConn(transport Transport) *BasicClientConn {
	return &BasicClientConn{
		transport: transport,
	}
}

func (c *BasicClientConn) Call(ctx context.Context, method string, params interface{}, reply interface{}) error {
	var id ID

	if reply != nil {
		// RPC call with a response.
		seqID := atomic.AddUint64(&c.seq, 1)
		id = ID(strconv.FormatUint(seqID, 10))
	} else {
		// Notification without a response.
		id = nil
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req := &Request{
		Version: Version,
		Method:  method,
		Params:  paramsJSON,
		ID:      id,
		Context: ctx,
	}

	resp, err := c.transport.SendRequest(req)
	if err != nil {
		return err
	}

	if id == nil {
		return nil
	}
	if resp.Error != nil {
		return resp.Error
	}

	err = json.Unmarshal(resp.Result, reply)
	if err != nil {
		return err
	}
	return nil
}
