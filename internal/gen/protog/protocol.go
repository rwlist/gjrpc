package protog

import (
	astinfo2 "github.com/rwlist/gjrpc/internal/gen/astinfo"
	"strings"
)

type Protocol struct {
	Package  *astinfo2.Package
	Services []Service
	Types    map[string]ProtocolType
}

func (p Protocol) FindServiceByGoType(name string) *Service {
	name = strings.TrimPrefix(name, p.Package.PkgName+".")

	ptype, ok := p.Types[name]
	if !ok {
		return nil
	}
	return ptype.Service
}

type Service struct {
	Path      []string
	Interface *astinfo2.Type
	Methods   []Method
}

type Method struct {
	Path   []string
	Method *astinfo2.Method
}

type ProtocolType struct {
	Service *Service
}
