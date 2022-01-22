package astinfo

type Kind string

const (
	Interface Kind = "interface"
	Struct    Kind = "struct"
	Alias     Kind = "alias" // type localType AnotherType
)

type Type struct {
	Name        string
	Kind        Kind
	Annotations []Annotation
	Fields      []Field
	Methods     []Method
}

type Annotation struct {
	Key    string
	Values []string
}

type Field struct {
	Name        string
	Type        string
	Annotations []Annotation
}

type Method struct {
	Name        string
	Params      []Field
	Results     []Field
	Annotations []Annotation
}
