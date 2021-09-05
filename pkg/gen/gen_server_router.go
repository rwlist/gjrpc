package gen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/rwlist/gjrpc/pkg/gen/astinfo"
	"github.com/rwlist/gjrpc/pkg/gen/protog"
	"github.com/rwlist/gjrpc/pkg/gen/router"
)

type genServerRouterArgs struct {
	protoPkg       string
	handlersStruct string
}

func generateServerRouter(args *genServerRouterArgs) (*jen.File, error) {
	proto, err := protog.Parse(args.protoPkg)
	if err != nil {
		return nil, err
	}

	currentPkg, err := astinfo.ParseDir(".")
	if err != nil {
		return nil, err
	}

	handlers, ok := currentPkg.Types[args.handlersStruct]
	if !ok {
		return nil, fmt.Errorf("struct %s not found in current directory", args.handlersStruct)
	}

	names := router.Names{
		StructName:      "Router",
		ConstructorName: "NewRouter",
		HandlersField:   "handlers",
		HandlersType:    args.handlersStruct,
		R:               "r",
		JsonrpcPkg:      "github.com/rwlist/gjrpc/pkg/jsonrpc",
	}

	route, err := router.NewRouter(proto, currentPkg, handlers, &names)
	if err != nil {
		return nil, err
	}

	src, err := route.GenerateSrc()
	if err != nil {
		return nil, err
	}

	return src, nil
}
