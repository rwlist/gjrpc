package protog

import (
	"github.com/pkg/errors"
	"github.com/rwlist/gjrpc/internal/gen/astinfo"
)

type Service struct {
	Path      []string
	Interface *astinfo.TypeDecl
	Methods   []Method
}

func parseService(info *astinfo.TypeDecl) (*Service, error) {
	if info.Kind != astinfo.Interface {
		return nil, nil
	}

	var serviceAnno *astinfo.Annotation
	for _, anno := range info.Annotations {
		switch anno.Key {
		case "gjrpc:service":
			if serviceAnno != nil {
				return nil, errors.Errorf("duplicated annotation %s on type %s", anno.Key, info.Name)
			}
			anno := anno
			serviceAnno = &anno
		default:
			return nil, errors.Errorf("unknown annotation %s on type %s", anno.Key, info.Name)
		}
	}

	if serviceAnno == nil {
		return nil, nil
	}

	if len(serviceAnno.Values) != 1 {
		return nil, errors.Errorf("invalid annotation %s on type %s", serviceAnno.Key, info.Name)
	}

	servicePath := StringToPath(serviceAnno.Values[0])

	var methods []Method
	for _, astMethod := range info.Methods {
		astMethod := astMethod
		method, err := parseMethod(&astMethod)
		if err != nil {
			return nil, err
		}
		if method == nil {
			return nil, errors.Errorf("method %s.%s has no valid annotations", info.Name, astMethod.Name)
		}

		method.FullPath = PathToString(PathAppend(servicePath, method.Path...))
		methods = append(methods, *method)
	}

	return &Service{
		Path:      servicePath,
		Interface: info,
		Methods:   methods,
	}, nil
}

type Method struct {
	Path       []string
	FullPath   string
	Method     *astinfo.Method
	ParamsType *astinfo.TypeRef
	ResultType *astinfo.TypeRef
}

func parseMethod(method *astinfo.Method) (*Method, error) {
	var methodAnno *astinfo.Annotation
	for _, anno := range method.Annotations {
		switch anno.Key {
		case "gjrpc:method":
			if methodAnno != nil {
				return nil, errors.Errorf("duplicated annotation %s on method %s", anno.Key, method.Name)
			}
			anno := anno
			methodAnno = &anno
		default:
			return nil, errors.Errorf("unknown annotation %s on method %s", anno.Key, method.Name)
		}
	}

	if methodAnno == nil {
		return nil, nil
	}

	if len(methodAnno.Values) != 1 {
		return nil, errors.Errorf("invalid annotation %s on method %s", methodAnno.Key, method.Name)
	}

	var paramsType *astinfo.TypeRef
	if len(method.Params) > 0 && isContextType(method.Params[0].Type) {
		method.Params = method.Params[1:]
	}
	if len(method.Params) != 0 {
		if len(method.Params) != 1 {
			return nil, errors.Errorf("method %s has more than one parameter, only single params objects is supported", method.Name)
		}
		paramsType = method.Params[0].Type
	}

	var resultType *astinfo.TypeRef
	if len(method.Results) > 2 || len(method.Results) < 1 {
		return nil, errors.Errorf("method %s has more than two results, only single result object and error is supported", method.Name)
	}
	if !method.Results[len(method.Results)-1].Type.IsError() {
		return nil, errors.Errorf("method %s must have error as the last result parameter", method.Name)
	}
	if len(method.Results) == 2 {
		resultType = method.Results[0].Type
	}

	return &Method{
		Path:       StringToPath(methodAnno.Values[0]),
		Method:     method,
		ParamsType: paramsType,
		ResultType: resultType,
	}, nil
}

func isContextType(t *astinfo.TypeRef) bool {
	return t.RefKind == astinfo.RefRef && t.Name == "Context" && t.ExternalPkg == "context"
}
