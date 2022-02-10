package gen_server

import (
	"encoding/json"
	"github.com/rwlist/gjrpc/examples/1_simple_service/proto"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type inventoryImpl struct {
	foo func()
	bar func()
}

func (i *inventoryImpl) Foo() (*proto.Foo, error) {
	i.foo()
	return &proto.Foo{}, nil
}

func (i *inventoryImpl) Bar(*proto.Bar) error {
	i.bar()
	return nil
}

func TestRouter(t *testing.T) {
	impl := &inventoryImpl{}
	router := NewRouter(&Handlers{Inventory: impl}, nil)

	type testcase struct {
		method  string
		funcPtr *func()
	}

	cases := []testcase{
		{
			method:  "inventory.foo",
			funcPtr: &impl.foo,
		},
		{
			method:  "inventory.bar",
			funcPtr: &impl.bar,
		},
	}
	for _, test := range cases {
		called := false
		fun := func() { called = true }
		*test.funcPtr = fun

		router.Handle(&jsonrpc.Request{
			Version: jsonrpc.Version,
			Method:  test.method,
			Params:  json.RawMessage("{}"),
			ID:      jsonrpc.ID("0"),
		})
		assert.True(t, called)
	}
}
