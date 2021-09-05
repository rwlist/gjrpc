package router

import (
	"fmt"

	"github.com/rwlist/gjrpc/pkg/gen/astinfo"
	"github.com/rwlist/gjrpc/pkg/gen/protog"
)

type Names struct {
	StructName      string
	ConstructorName string
	HandlersField   string
	HandlersType    string
	R               string
	JsonrpcPkg      string
}

type Router struct {
	proto      *protog.Protocol
	currentPkg *astinfo.Package
	handlers   []*handler
	endpoints  []endpoint
	tree       *node

	Names
}

func NewRouter(proto *protog.Protocol, currentPkg *astinfo.Package, handlersStruct *astinfo.Type, names *Names) (*Router, error) {
	if handlersStruct.Kind != astinfo.Struct {
		return nil, fmt.Errorf("%s must be struct", handlersStruct.Name)
	}

	var handlers []*handler
	var endpoints []endpoint

	for _, handlerField := range handlersStruct.Fields {
		h, err := newHandlerFromAST(handlerField, currentPkg, proto)
		if err != nil {
			return nil, err
		}
		if h == nil {
			continue
		}

		handlers = append(handlers, h)
		endpoints = append(endpoints, h.endpoints...)
	}

	root, err := newTree(handlers, endpoints)
	if err != nil {
		return nil, err
	}

	return &Router{
		proto:      proto,
		currentPkg: currentPkg,
		handlers:   handlers,
		endpoints:  endpoints,
		tree:       root,
		Names:      *names,
	}, nil
}
