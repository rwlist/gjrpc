package router

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rwlist/gjrpc/internal/gen/astinfo"
	"github.com/rwlist/gjrpc/internal/gen/protog"

	"github.com/dave/jennifer/jen"
)

type queueItem struct {
	n    *node
	path []string
}

func (r *Router) GenerateSrc() (*jen.File, error) {
	f := jen.NewFile(r.currentPkg.PkgName)

	r.genStruct(f)
	f.Line()
	r.genConstructor(f)
	f.Line()
	r.genNotFound(f)
	f.Line()
	r.genMainHandle(f)

	queue := []queueItem{{
		n:    r.tree,
		path: []string{},
	}}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		f.Line()
		err := r.genHandle(f, item)
		if err != nil {
			return nil, err
		}

		for _, child := range item.n.sortedChildren() {
			next, nextNode := child.name, child.next

			var nextPath []string
			nextPath = append(nextPath, item.path...)
			nextPath = append(nextPath, next)

			queue = append(queue, queueItem{
				n:    nextNode,
				path: nextPath,
			})
		}
	}

	return f, nil
}

func (r *Router) genStruct(f *jen.File) {
	f.Type().Id(r.StructName).Struct(
		jen.Id(r.HandlersField).Op("*").Id(r.HandlersType),
		jen.Id(r.ErrorConverterField).Qual(r.JsonrpcPkg, "ErrorConverter"),
	)
	//type Router struct {
	//	handlers *Handlers
	//	convertError jsonrpc.ErrorConverter
	//}
}

func (r *Router) genConstructor(f *jen.File) {
	f.Func().Id(r.ConstructorName).Params(
		jen.Id(r.HandlersField).Op("*").Id(r.HandlersType),
		jen.Id(r.ErrorConverterField).Qual(r.JsonrpcPkg, "ErrorConverter"),
	).Params(jen.Op("*").Id(r.StructName)).Block(
		jen.If(jen.Id(r.ErrorConverterField).Op("==").Nil()).Block(
			jen.Id(r.ErrorConverterField).Op("=").Qual(r.JsonrpcPkg, "DefaultErrorConverter"),
		),
		jen.Return().Op("&").Id(r.StructName).Values(jen.Dict{
			jen.Id(r.HandlersField):       jen.Id(r.HandlersField),
			jen.Id(r.ErrorConverterField): jen.Id(r.ErrorConverterField),
		}),
	)
	//func NewRouter(handlers Handlers, convertError jsonrpc.ErrorConverter) *Router {
	//	if convertError == nil {
	//	  convertError = jsonrpc.DefaultErrorConverter
	// }
	//	return &Router{
	//	  handlers: handlers,
	//	  convertError: convertError,
	//	}
	//}
}

func (r *Router) genMainHandle(f *jen.File) {
	r.methodHandle(f, "Handle",
		jen.Id("req").Add(r.jsonrpcRequest()),
	).Block(
		jen.Id("path").Op(":=").Qual("strings", "Split").Call(jen.Id("req").Dot("Method"), jen.Lit(".")),
		jen.Return().Id("r").Dot("handle").Call(jen.Id("path"), jen.Id("req")),
	)
	// func (r *Router) Handle(req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	//	path := strings.Split(req.Method, ".")
	//	return r.handle(path, req)
	//}
}

func (r *Router) genNotFound(f *jen.File) {
	r.methodHandle(f, "notFound").Block(jen.Return(
		jen.Nil(),
		jen.Op("&").Qual(r.JsonrpcPkg, "MethodNotFound"),
	))
	// func (r *Router) notFound() (jsonrpc.Result, *jsonrpc.Error) {
	//	return nil, jsonrpc.MethodNotFound
	//}
}

func (r *Router) genHandle(f *jen.File, item queueItem) (err error) {
	funcName := r.handleFuncName(item.path)

	r.methodHandle(f, funcName,
		jen.Id("path").Index().String(),
		jen.Id("req").Add(r.jsonrpcRequest()),
	).BlockFunc(func(fun *jen.Group) {
		fun.If().Id("len").Call(jen.Id("path")).Op("==").Lit(0).BlockFunc(func(g *jen.Group) {
			if item.n.endpoint == nil {
				g.Return().Id("r").Dot("notFound").Call()
			} else {
				err = r.genMethodCall(g, item.n.endpoint)
				if err != nil {
					return
				}
			}
		})

		r.genRouteSwitch(fun, item)

		fun.Return().Id("r").Dot("notFound").Call()
	})

	return err
}

//nolint:funlen
func (r *Router) genMethodCall(g *jen.Group, e *endpoint) error {
	reqVar := jen.Id("request")
	method := e.methodImpl

	var arguments []jen.Code
	var requestType *jen.Statement
	for _, param := range method.methodAST.Params {
		switch {
		case param.Type.KindaIs("jsonrpc.Request"):
			arguments = append(arguments, jen.Id("req"))
		case param.Type.KindaIs("context.Context"):
			arguments = append(arguments, jen.Id("req").Dot("Context"))
		default:
			// TODO: assert compatibility with methodProto
			if requestType != nil {
				return errors.Errorf("should be exactly one request object, found %s and %v", requestType.GoString(), param.Type)
			}

			if param.Type.RefKind == astinfo.RefPrimitive {
				requestType = jen.Id(param.Type.Name)
			} else if param.Type.RefKind == astinfo.RefRef &&
				(param.Type.ExternalPkg == "" || param.Type.PackageLookSame(r.proto.Package.PkgImportPath)) {
				requestType = jen.Qual(r.proto.Package.PkgImportPath, param.Type.Name)
			} else {
				return errors.Errorf("unsupported type ref %v", param.Type)
			}

			var reqArgument jen.Code = reqVar
			if param.Type.IsPointer {
				reqArgument = jen.Op("&").Add(reqArgument)
			}

			arguments = append(arguments, reqArgument)
		}
	}

	if requestType != nil {
		g.Var().Add(reqVar).Add(requestType)
		g.If(
			jen.Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(
				jen.Id("req").Dot("Params"),
				jen.Op("&").Add(reqVar),
			),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return().Id("r").Dot(r.ErrorConverterField).Call(jen.Err()),
		)
		// var req SomeRequest
		// if err = json.Unmarshal(params, &req); err != nil {
		//	return r.convertError(err)
		//}
	}

	var (
		results     []jen.Code
		resResponse *astinfo.Field
		resError    *astinfo.Field
		respVar     = "res"
		errVar      = "err"
	)

	for _, res := range method.methodAST.Results {
		switch {
		case res.Type.IsError():
			if resError != nil {
				return errors.Errorf("should be exactly one error object, found %#v and %#v", resError.Type, res.Type)
			}
			res := res
			resError = &res
			results = append(results, jen.Id(errVar))
		default:
			if resResponse != nil {
				return errors.Errorf("should be exactly one result, found %#v and %#v", resResponse.Type, res.Type)
			}
			res := res
			resResponse = &res
			results = append(results, jen.Id(respVar))
		}
	}

	if resError == nil {
		return errors.Errorf("function %s doesn't return error object", method.methodAST.Name)
	}

	// TODO: validate resResponse and resError types

	g.List(results...).Op(":=").Id("r").Dot(r.HandlersField).Dot(method.handler).Dot(method.methodAST.Name).Call(arguments...)
	// res, err := r.Handlers.Inventory.Foo(request)

	g.If(jen.Id(errVar).Op("!=").Nil()).Block(
		jen.Return(jen.Id("r").Dot(r.ErrorConverterField).Call(jen.Id(errVar))),
	)
	// if err != nil {
	//	return r.convertError(err)
	//}

	if resResponse == nil {
		g.Return(jen.Struct().Values(jen.Dict{}), jen.Nil())
		// struct{}{}
	} else {
		g.Return(jen.Id(respVar), jen.Nil())
		// return res, nil
	}

	return nil
}

func (r *Router) genRouteSwitch(g *jen.Group, item queueItem) {
	if len(item.n.children) == 0 {
		// nothing to route
		return
	}

	g.Switch(jen.Id("path").Index(jen.Lit(0))).BlockFunc(func(sw *jen.Group) {
		for _, child := range item.n.sortedChildren() {
			name := child.name
			nextFunc := r.handleFuncName(protog.PathAppend(item.path, name))

			sw.Case(jen.Lit(name)).Block(
				jen.Return().Id("r").Dot(nextFunc).Call(
					jen.Id("path").Index(jen.Lit(1), jen.Empty()),
					jen.Id("req"),
				),
			)
			// case "bar":
			// 		return r.handleBar(path[1:])
		}
	})
}

func (r *Router) method(f *jen.File, name string) *jen.Statement {
	return f.Func().Params(jen.Id("r").Op("*").Id(r.StructName)).Id(name)
}

func (r *Router) methodHandle(f *jen.File, name string, params ...jen.Code) *jen.Statement {
	return r.method(f, name).Params(params...).Params(
		jen.Qual(r.JsonrpcPkg, "Result"),
		jen.Op("*").Qual(r.JsonrpcPkg, "Error"),
	)
}

func (r *Router) handleFuncName(path []string) string {
	funcName := "handle"
	for _, el := range path {
		funcName += strings.Title(el)
	}
	return funcName
}

func (r *Router) jsonrpcRequest() jen.Code {
	return jen.Op("*").Qual(r.JsonrpcPkg, "Request")
}
