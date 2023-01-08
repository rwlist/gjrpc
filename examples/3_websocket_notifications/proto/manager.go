package proto

import "context"

//gjrpc:service manager
type Manager interface {
	//gjrpc:method connect
	Connect(context.Context, *ConnectRequest) (*ConnectResponse, error)
}

type ConnectRequest struct {
	AgentToken string
}

type ConnectResponse struct{}
