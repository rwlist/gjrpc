package protog

import (
	astinfo2 "github.com/rwlist/gjrpc/internal/gen/astinfo"
)

// Parse protocol definition located in specified directory.
func Parse(path string) (*Protocol, error) {
	pkg, err := astinfo2.ParseDir(path)
	if err != nil {
		return nil, err
	}

	ptypes := map[string]ProtocolType{}
	var services []Service
	var models []Model
	for name, t := range pkg.Types {
		serv, err := parseService(t)
		if err != nil {
			return nil, err
		}

		model, err := parseModel(t)
		if err != nil {
			return nil, err
		}

		ptype := ProtocolType{}
		if serv != nil {
			services = append(services, *serv)
			ptype.Service = serv
		} else if model != nil {
			models = append(models, *model)
			ptype.Model = model
		}

		ptypes[name] = ptype
	}

	proto := &Protocol{
		Package:  pkg,
		Services: services,
		Models:   models,
		Types:    ptypes,
	}
	return proto, nil
}
