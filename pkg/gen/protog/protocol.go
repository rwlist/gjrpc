package protog

import (
	"strings"

	"github.com/rwlist/gjrpc/pkg/gen/astinfo"
)

type Protocol struct {
	Package  *astinfo.Package
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
	Interface *astinfo.Type
	Methods   []Method
}

type Method struct {
	Path   []string
	Method *astinfo.Method
}

type ProtocolType struct {
	Service *Service
}
