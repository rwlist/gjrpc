package wsexample

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"wsexample/agent"
	"wsexample/manager"
	"wsexample/proto"

	"github.com/go-chi/chi/v5"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"github.com/rwlist/gjrpc/pkg/transport"
	"github.com/stretchr/testify/assert"
)

// Agent is implemented by the "client".
type Agent struct {
	executeFunc func(ctx context.Context, req *proto.ExecuteV1Request) (*proto.ExecuteV1Response, error)
}

var _ proto.Agent = (*Agent)(nil)

// ExecuteV1 implements proto.Agent
func (a *Agent) ExecuteV1(ctx context.Context, req *proto.ExecuteV1Request) (*proto.ExecuteV1Response, error) {
	return a.executeFunc(ctx, req)
}

// Status implements proto.Agent
func (a *Agent) Status(context.Context, *proto.StatusRequest) (*proto.StatusResponse, error) {
	// not used in this test
	panic("unimplemented")
}

// Manager is implemented by the "server".
type Manager struct {
	conns []*transport.WSConn
}

var _ proto.Manager = (*Manager)(nil)

// Connect implements proto.ServerAgents
func (m *Manager) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	wsConn := transport.CtxWSConn(ctx)
	m.conns = append(m.conns, wsConn)
	return &proto.ConnectResponse{}, nil
}

func createManagerServer(ctx context.Context, managerImpl proto.Manager) (string, error) {
	handlers := &manager.Handlers{
		Manager: managerImpl,
	}
	rpcHandler := manager.NewRouter(handlers, nil)
	wsHandler := transport.NewHandlerWebSocket(rpcHandler)

	r := chi.NewRouter()
	r.Handle("/ws", wsHandler)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	fmt.Println("Starting server on port", listener.Addr().(*net.TCPAddr).Port)

	go func() {
		if err := http.Serve(listener, r); err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	go func() {
		// wait for context to be canceled
		<-ctx.Done()
		// close the listener
		listener.Close()
	}()

	return listener.Addr().String(), nil
}

// TestWebsocketsWork starts a server on a random port and connects a client to it.
// It then checks that the client can subscribe and receive notifications from the server.
func TestWebsocketsWork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create and start server
	managerImpl := &Manager{}
	serverAddr, err := createManagerServer(ctx, managerImpl)
	assert.Nil(t, err)

	// create client to receive callbacks from server
	agentImpl := &Agent{
		executeFunc: nil, // will be set later
	}
	clientHandler := &agent.Handlers{
		Agent: agentImpl,
	}
	clientCallback := agent.NewRouter(clientHandler, nil)

	// connect to the server
	conn, err := transport.DialOnce(ctx, "ws://"+serverAddr+"/ws", clientCallback)
	assert.Nil(t, err)

	// create client to the manager service
	managerCli := manager.NewManagerClient(jsonrpc.NewClientConn(conn))

	// start the subscription
	_, err = managerCli.Connect(ctx, &proto.ConnectRequest{})
	assert.Nil(t, err)

	// get the server->client conn
	assert.Len(t, managerImpl.conns, 1)
	serverConn := managerImpl.conns[0]
	agentClient := agent.NewAgentClient(jsonrpc.NewClientConn(serverConn))

	var tests = []struct {
		req  *proto.ExecuteV1Request
		resp *proto.ExecuteV1Response
		err  error
	}{
		{
			req: &proto.ExecuteV1Request{
				Command: "echo hello",
			},
			resp: &proto.ExecuteV1Response{
				Stdout: "hello",
				Stderr: "echo hello",
			},
			err: nil,
		},
		{
			req: &proto.ExecuteV1Request{
				Command: "",
			},
			resp: &proto.ExecuteV1Response{
				Stdout: "",
				Stderr: "",
			},
			err: nil,
		},
		{
			req: &proto.ExecuteV1Request{
				Command: "error",
			},
			resp: nil,
			err:  fmt.Errorf("error"),
		},
	}

	for _, test := range tests {
		// update hook
		agentImpl.executeFunc = func(ctx context.Context, req *proto.ExecuteV1Request) (*proto.ExecuteV1Response, error) {
			assert.Equal(t, test.req, req)
			return test.resp, test.err
		}

		// make request
		resp, err := agentClient.ExecuteV1(ctx, test.req)
		if test.err != nil {
			assert.Nil(t, resp)
			jerr := err.(*jsonrpc.Error)
			assert.Equal(t, `"`+test.err.Error()+`"`, string(jerr.Data))
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.resp, resp)
		}
	}
}
