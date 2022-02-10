package astinfo

type Primitive struct {
	Kind string
}

var (
	PrimitiveBool        = &Primitive{Kind: "bool"}
	PrimitiveInteger     = &Primitive{Kind: "integer"}
	PrimitiveString      = &Primitive{Kind: "string"}
	PrimitiveFloat       = &Primitive{Kind: "float"}
	PrimitiveError       = &Primitive{Kind: "error"}
	PrimitiveUnsupported = &Primitive{Kind: "unsupported"}
)

func IsPrimitive(t string) *Primitive { //nolint:gocyclo
	// all primitives are located inside src/builtin/builtin.go
	switch t {
	case "bool":
		return PrimitiveBool
	case "uint8":
		return PrimitiveInteger
	case "uint16":
		return PrimitiveInteger
	case "uint32":
		return PrimitiveInteger
	case "uint64":
		return PrimitiveInteger
	case "int8":
		return PrimitiveInteger
	case "int16":
		return PrimitiveInteger
	case "int32":
		return PrimitiveInteger
	case "int64":
		return PrimitiveInteger
	case "float32":
		return PrimitiveFloat
	case "float64":
		return PrimitiveFloat
	case "complex64":
		return PrimitiveUnsupported
	case "complex128":
		return PrimitiveUnsupported
	case "string":
		return PrimitiveString
	case "int":
		return PrimitiveInteger
	case "uint":
		return PrimitiveInteger
	case "uintptr":
		return PrimitiveUnsupported
	case "byte":
		return PrimitiveInteger
	case "rune":
		return PrimitiveInteger
	case "error":
		return PrimitiveError
	}
	return nil
}
