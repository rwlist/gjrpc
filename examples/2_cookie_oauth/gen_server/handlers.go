//go:generate gjrpc gen:server:router --protoPkg=../proto --handlersStruct=Handlers --out=gen_router.go
package gen_server

import (
	"2_cookie_oauth/proto"
	"context"
)

type Handlers struct {
	//gjrpc:handle-route proto.Auth
	Auth AuthImpl
}

type AuthImpl interface {
	Oauth(ctx context.Context) (proto.OAuthResponse, error)
	Status(ctx context.Context) (proto.AuthStatus, error)
}
