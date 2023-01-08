package proto

import "context"

//gjrpc:service agent
type Agent interface {
	//gjrpc:method status
	Status(context.Context, *StatusRequest) (*StatusResponse, error)

	//gjrpc:method executev1
	ExecuteV1(context.Context, *ExecuteV1Request) (*ExecuteV1Response, error)
}

type StatusRequest struct{}

type StatusResponse struct {
	Version string
}

type ExecuteV1Request struct {
	Command string
}

type ExecuteV1Response struct {
	Stdout string
	Stderr string
}
