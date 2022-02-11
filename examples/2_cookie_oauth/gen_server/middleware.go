package gen_server

import (
	"github.com/rwlist/gjrpc/pkg/gjserver"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

type AuthProvider interface {
	Check(accessToken string) error
}

func AuthMiddleware(auth AuthProvider, exceptions []string) jsonrpc.Middleware {
	exceptionMap := make(map[string]struct{})
	for _, e := range exceptions {
		exceptionMap[e] = struct{}{}
	}

	return func(next jsonrpc.Handler) jsonrpc.Handler {
		return func(req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
			if _, ok := exceptionMap[req.Method]; ok {
				return next(req)
			}

			ctx := req.Context
			httpReq := gjserver.HttpRequest(ctx)
			accessToken := AccessTokenFromRequest(httpReq)
			err := auth.Check(accessToken)
			if err != nil {
				return nil, jsonrpc.UnknownError.WithData(err.Error())
			}

			return next(req)
		}
	}
}
