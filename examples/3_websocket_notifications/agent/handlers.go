//go:generate gjrpc gen:server:router --protoPkg=../proto --handlersStruct=Handlers --out=gen_router.go
package agent

import (
	"wsexample/proto"
)

type Handlers struct {
	//gjrpc:handle-route proto.Agent
	Agent proto.Agent
}
