# 1_simple_service

In this example, there is single service named `inventory`, with two functions:

- `inventory.foo() -> Foo` - function without params, returns Foo or error.
- `inventory.bar(Bar)` - function accepts Bar object, returns nothing or error.

All that is defined in `proto` directory.

## Code generation

There is `Handlers` struct defined in `gen_server/handlers.go`, which contains implementation of different protocol services.
If you run `go generate` in `gen_server` directory, new file `gen_router.go` will be generated, which implement `jsonrpc.Handler` type and has generated code for routing across all the provided implementations.
