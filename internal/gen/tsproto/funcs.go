package tsproto

import (
	"github.com/rwlist/gjrpc/internal/gen/astinfo"
)

var funcMap = map[string]interface{}{
	"convertGoType": convertGoType,
}

func convertGoType(t *astinfo.TypeRef) string {
	// TODO: pointers are ignored, worth an optional thing?
	switch t.RefKind {
	case astinfo.RefPrimitive:
		return convertGoPrimitive(t.Primitive)
	case astinfo.RefEmbedded:
		return "unknown"
	case astinfo.RefMap:
		return "Record<" + convertGoType(t.KeyType) + ", " + convertGoType(t.ValueType) + ">"
	case astinfo.RefSlice:
		return convertGoType(t.ValueType) + "[]"
	case astinfo.RefRef:
		if t.ExternalPkg == "" {
			return t.Name
		} else {
			// TODO: possible to embed?
			return "unknown"
		}
	}

	return "unknown"
}

func convertGoPrimitive(primitive *astinfo.Primitive) string {
	switch primitive {
	case astinfo.PrimitiveBool:
		return "boolean"
	case astinfo.PrimitiveInteger:
		return "number"
	case astinfo.PrimitiveString:
		return "string"
	case astinfo.PrimitiveFloat:
		return "number"
	}

	return "unknown"
}
