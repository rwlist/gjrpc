package transport

import (
	"context"
	"net/http"
)

type ctxKey string

const (
	ctxHttpRequest  = ctxKey("http.Request")
	ctxHttpResponse = ctxKey("http.ResponseWriter")
	ctxWSConn       = ctxKey("transport.WSConn")
)

func CtxHttpRequest(ctx context.Context) *http.Request {
	v, ok := ctx.Value(ctxHttpRequest).(*http.Request)
	if !ok {
		return nil
	}
	return v
}

func CtxHttpResponse(ctx context.Context) http.ResponseWriter {
	v, ok := ctx.Value(ctxHttpResponse).(http.ResponseWriter)
	if !ok {
		return nil
	}
	return v
}

func CtxWSConn(ctx context.Context) *WSConn {
	v, ok := ctx.Value(ctxWSConn).(*WSConn)
	if !ok {
		return nil
	}
	return v
}
