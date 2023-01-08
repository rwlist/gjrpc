package agent

import (
	"context"
	proto "wsexample/proto"

	jsonrpc "github.com/rwlist/gjrpc/pkg/jsonrpc"
)

// TODO: generate this code automatically

type cliAgent struct {
	conn jsonrpc.ClientConnInterface
}

func NewAgentClient(conn jsonrpc.ClientConnInterface) *cliAgent {
	return &cliAgent{conn: conn}
}

// ExecuteV1 implements proto.Agent
func (c *cliAgent) ExecuteV1(ctx context.Context, req *proto.ExecuteV1Request) (*proto.ExecuteV1Response, error) {
	var resp proto.ExecuteV1Response
	err := c.conn.Call(ctx, "agent.executev1", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Status implements proto.Agent
func (c *cliAgent) Status(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	var resp proto.StatusResponse
	err := c.conn.Call(ctx, "agent.status", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
