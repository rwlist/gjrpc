//go:generate gjrpc gen:server:router --protoPkg=../proto --handlersStruct=Handlers --out=gen_router.go
package gen_server

import "github.com/rwlist/gjrpc/examples/1_simple_service/proto"

type Handlers struct {
	//gjrpc:handle-route proto.Inventory
	Inventory proto.Inventory
}
