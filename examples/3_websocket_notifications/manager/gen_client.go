package manager

import (
	"context"
	"wsexample/proto"

	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

// TODO: generate this code automatically

type cliManager struct {
	conn jsonrpc.ClientConnInterface
}

func NewManagerClient(conn jsonrpc.ClientConnInterface) *cliManager {
	return &cliManager{conn: conn}
}

func (c *cliManager) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	var resp proto.ConnectResponse
	err := c.conn.Call(ctx, "manager.connect", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
