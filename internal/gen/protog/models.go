package protog

import (
	"github.com/rwlist/gjrpc/internal/gen/astinfo"
)

type Model struct {
	Struct *astinfo.TypeDecl
	Fields []Field
}

type Field struct {
	Name     string
	Type     *astinfo.TypeRef
	AstField *astinfo.Field
}

func parseModel(info *astinfo.TypeDecl) (*Model, error) { //nolint:unparam
	if info.Kind != astinfo.Struct {
		return nil, nil
	}

	var fields []Field
	for i := range info.Fields {
		fields = append(fields, Field{
			Name:     info.Fields[i].Name,
			Type:     info.Fields[i].Type,
			AstField: &info.Fields[i],
		})
	}

	return &Model{
		Struct: info,
		Fields: fields,
	}, nil
}
