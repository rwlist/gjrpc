package proto

// Inventory is a simple service with two functions.
//
//gjrpc:service inventory
type Inventory interface {
	//gjrpc:method foo
	Foo() (Foo, error)

	//gjrpc:method bar
	Bar(Bar) error
}

type Foo struct {
	Info  CommonInfo
	Index int
}

type Bar struct {
	Info CommonInfo
	Name string
}
