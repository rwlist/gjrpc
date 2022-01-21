//go:generate gjrpc gen:server:router --protoPkg=../proto --handlersStruct=Handlers --out=gen_router.go
package gen_server

import "2_cookie_oauth/proto"

type Handlers struct {
	//gjrpc:handle-route proto.Auth
	Auth proto.Auth
}
