//go:generate gjrpc gen:server:router --protoPkg=../proto --handlersStruct=Handlers --out=gen_router.go
package manager

import proto "wsexample/proto"

type Handlers struct {
	//gjrpc:handle-route proto.Manager
	Manager proto.Manager
}
