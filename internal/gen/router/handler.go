package router

import (
	"fmt"
	astinfo2 "github.com/rwlist/gjrpc/internal/gen/astinfo"
	"github.com/rwlist/gjrpc/internal/gen/protog"
	"os"
)

type handler struct {
	userHandler   *userHandler
	targetService *protog.Service
	path          []string
	endpoints     []endpoint
}

func newHandlerFromAST(f astinfo2.Field, currentPkg *astinfo2.Package, proto *protog.Protocol) (*handler, error) {
	var routeAnno *astinfo2.Annotation
	for _, anno := range f.Annotations {
		switch anno.Key {
		case "gjrpc:handle-route":
			if routeAnno != nil {
				return nil, fmt.Errorf("duplicated annotation %s on field %s", anno.Key, f.Name)
			}
			anno := anno
			routeAnno = &anno
		default:
			return nil, fmt.Errorf("unknown annotation %s on field %s", anno.Key, f.Name)
		}
	}

	if routeAnno == nil {
		return nil, nil
	}

	if len(routeAnno.Values) != 1 {
		return nil, fmt.Errorf("invalid annotation %s on field %s", routeAnno.Key, f.Name)
	}
	targetServiceName := routeAnno.Values[0]

	targetService := proto.FindServiceByGoType(targetServiceName)
	if targetService == nil {
		return nil, fmt.Errorf("service named %s not found", targetServiceName)
	}

	userAST, err := lookupUserHandler(currentPkg, proto, f.Type)
	if err != nil {
		return nil, err
	}

	handlerImpl, err := prepareUserHandler(f.Name, userAST)
	if err != nil {
		return nil, err
	}

	var endpoints []endpoint
	for _, method := range targetService.Methods {
		userMethod, ok := handlerImpl.methods[method.Method.Name]
		if !ok {
			_, _ = fmt.Fprintf(
				os.Stderr,
				"WARN: method %s.%s not implemented in %s %s",
				targetService.Interface.Name,
				method.Method.Name,
				f.Name,
				f.Type,
			)
			continue
		}

		method := method
		userMethod.methodProto = &method

		var path []string
		path = append(path, targetService.Path...)
		path = append(path, method.Path...)

		endpoints = append(endpoints, endpoint{
			path:       path,
			methodImpl: userMethod,
		})
	}

	return &handler{
		userHandler:   handlerImpl,
		targetService: targetService,
		path:          targetService.Path,
		endpoints:     endpoints,
	}, nil
}

func lookupUserHandler(currentPkg *astinfo2.Package, proto *protog.Protocol, userType string) (*astinfo2.Type, error) { //nolint:unparam
	// TODO: implement real lookup, even for outer packages
	serv := proto.FindServiceByGoType(userType)
	if serv != nil {
		return serv.Interface, nil
	}

	if localType, ok := currentPkg.Types[userType]; ok {
		return localType, nil
	}

	return nil, fmt.Errorf("type %s not found", userType)
}

type userHandler struct {
	methods map[string]*methodImpl
}

type methodImpl struct {
	// field in the handlers struct, usually service name
	handler string

	methodAST   astinfo2.Method
	methodProto *protog.Method
}

func prepareUserHandler(handler string, userAST *astinfo2.Type) (*userHandler, error) { //nolint:unparam
	uh := &userHandler{
		methods: map[string]*methodImpl{},
	}

	for _, method := range userAST.Methods {
		uh.methods[method.Name] = &methodImpl{
			handler:     handler,
			methodAST:   method,
			methodProto: nil, // will be filled if matches
		}
	}
	return uh, nil
}
