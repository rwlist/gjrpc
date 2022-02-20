package astinfo

import "strings"

const annotationPrefix = "gjrpc:"

type Kind string

const (
	Interface Kind = "interface"
	Struct    Kind = "struct"
	Alias     Kind = "alias" // type localType AnotherType
	Func      Kind = "func"  // func (args...) (res...)
)

type TypeDecl struct {
	Name        string
	Kind        Kind
	Annotations []Annotation

	// only for Struct
	Fields []Field

	// only for Interface
	Methods []Method

	// only for Alias
	Alias *TypeRef
}

type Embedded struct {
	Kind Kind
	// TODO: support more info
}

type RefKind int

const (
	RefUnknown   RefKind = iota
	RefPrimitive         // int, string, etc
	RefEmbedded          // interface{}, struct{}, func()
	RefMap               // map[int]*Foo
	RefSlice             // []Bar
	RefRef               // *ast.AnotherType
)

type TypeRef struct {
	// *Some -> true, otherwise false
	IsPointer bool

	// If type is ***Some, then the value is 2, for example
	AdditionalPointers int

	// describes what kind of type it is, behind the pointers
	RefKind RefKind

	// if non-nil, info about the primitive
	Primitive *Primitive

	// if non-nil, info about the embedded type
	Embedded *Embedded

	// if non-nil, type is a map
	KeyType *TypeRef

	// if key is nil and this is not, type is slice
	ValueType *TypeRef

	// otherwise, ref points to this name (which is other type)
	Name string

	// if non-empty, type is found in this package
	ExternalPkg string
}

func (r *TypeRef) IsError() bool {
	return r.RefKind == RefPrimitive && r.Primitive == PrimitiveError && !r.IsPointer
}

func (r *TypeRef) KindaIs(s string) bool {
	return r.ExternalPkg+"."+r.Name == s
}

func (r *TypeRef) PackageLookSame(s string) bool {
	if r.ExternalPkg == s {
		return true
	}

	sPath := strings.Split(s, "/")
	return sPath[len(sPath)-1] == r.ExternalPkg
}

type Annotation struct {
	Key    string
	Values []string
}

type Field struct {
	Name        string
	Type        *TypeRef
	Annotations []Annotation
}

type Method struct {
	Name        string
	Params      []Field
	Results     []Field
	Annotations []Annotation
}
