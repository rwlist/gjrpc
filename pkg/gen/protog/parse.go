package protog

import (
	"fmt"

	"github.com/rwlist/gjrpc/pkg/gen/astinfo"
)

// Parse protocol definition located in specified directory.
func Parse(path string) (*Protocol, error) {
	pkg, err := astinfo.ParseDir(path)
	if err != nil {
		return nil, err
	}

	ptypes := map[string]ProtocolType{}
	var services []Service
	for name, t := range pkg.Types {
		serv, err := parseService(t)
		if err != nil {
			return nil, err
		}

		ptype := ProtocolType{}
		if serv != nil {
			services = append(services, *serv)
			ptype.Service = serv
		}

		ptypes[name] = ptype
	}

	proto := &Protocol{
		Package:  pkg,
		Services: services,
		Types:    ptypes,
	}
	return proto, nil
}

func parseService(info *astinfo.Type) (*Service, error) {
	if info.Kind != astinfo.Interface {
		return nil, nil
	}

	var serviceAnno *astinfo.Annotation
	for _, anno := range info.Annotations {
		switch anno.Key {
		case "gjrpc:service":
			if serviceAnno != nil {
				return nil, fmt.Errorf("duplicated annotation %s on type %s", anno.Key, info.Name)
			}
			anno := anno
			serviceAnno = &anno
		default:
			return nil, fmt.Errorf("unknown annotation %s on type %s", anno.Key, info.Name)
		}
	}

	if serviceAnno == nil {
		return nil, nil
	}

	if len(serviceAnno.Values) != 1 {
		return nil, fmt.Errorf("invalid annotation %s on type %s", serviceAnno.Key, info.Name)
	}

	var methods []Method
	for _, astMethod := range info.Methods {
		astMethod := astMethod
		method, err := parseMethod(&astMethod)
		if err != nil {
			return nil, err
		}
		if method == nil {
			return nil, fmt.Errorf("method %s.%s has no valid annotations", info.Name, astMethod.Name)
		}

		methods = append(methods, *method)
	}

	return &Service{
		Path:      StringToPath(serviceAnno.Values[0]),
		Interface: info,
		Methods:   methods,
	}, nil
}

func parseMethod(method *astinfo.Method) (*Method, error) {
	var methodAnno *astinfo.Annotation
	for _, anno := range method.Annotations {
		switch anno.Key {
		case "gjrpc:method":
			if methodAnno != nil {
				return nil, fmt.Errorf("duplicated annotation %s on method %s", anno.Key, method.Name)
			}
			anno := anno
			methodAnno = &anno
		default:
			return nil, fmt.Errorf("unknown annotation %s on method %s", anno.Key, method.Name)
		}
	}

	if methodAnno == nil {
		return nil, nil
	}

	if len(methodAnno.Values) != 1 {
		return nil, fmt.Errorf("invalid annotation %s on method %s", methodAnno.Key, method.Name)
	}

	return &Method{
		Path:   StringToPath(methodAnno.Values[0]),
		Method: method,
	}, nil
}
